package impl

import (
	"dataProcess/config"
	"dataProcess/constants"
	"dataProcess/constants/atoi"
	"dataProcess/constants/compress"
	"dataProcess/constants/filej"
	"dataProcess/constants/method"
	"dataProcess/constants/plat"
	"dataProcess/constants/upload"
	"dataProcess/model"
	"dataProcess/repository/mocks"
	"dataProcess/service/remote/impl"
	"dataProcess/tools/utils/utils"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"
)

const fileSize = 512 * 1024 * 1024 // 限定方法size为512M

type DefinitionImpl struct {
	*CategoryImpl
	*MethodImpl
}

func CreateDefinitionImpl() *DefinitionImpl {
	return &DefinitionImpl{CreateCategoryImpl(), CreateMethodImpl()}
}

func getMethodsMark() (methodsMark string) {
	switch config.ApplicationConfig.SrcSource {
	case constants.Formalization:
		methodsMark = constants.Methods + constants.SubLine + constants.Formalization
	case constants.Testing:
		methodsMark = constants.Methods
	case constants.Integration:
		methodsMark = constants.Methods + constants.SubLine + constants.Integration
	case constants.Security:
		methodsMark = constants.Methods + constants.SubLine + constants.Security
	}
	return
}

func UPMethodOrmDecode(c *gin.Context) (upMethodOrm model.UPMethodOrm, err error) {
	jsonTags := utils.ReflectGetJsonTags(&upMethodOrm)
	values := GetPostFormValues(c, jsonTags)

	utils.ReflectSetStructValue(&upMethodOrm, jsonTags, values)

	err = UploadFileAssignInfo(c, &upMethodOrm, upload.FileSelf)
	if err != nil {
		return
	}

	err = UploadFileAssignInfo(c, &upMethodOrm, upload.Description)
	if err != nil {
		return
	}

	return
}

func UploadFileAssignInfo(c *gin.Context, upMethodOrm *model.UPMethodOrm, upKey string) (err error) {
	form, _ := c.MultipartForm()
	if _, ok := form.File[upKey]; !ok {
		return
	}

	fileSelf, err := c.FormFile(upKey)
	if err != nil {
		return
	}

	byFile, err := fileSelf.Open()
	if err != nil {
		return
	}

	fileNameArr := strings.Split(fileSelf.Filename, constants.Point)
	enName, extName, err := getEnNameAndExtName(fileNameArr)
	if err != nil {
		return
	}

	switch upKey {
	case upload.FileSelf:
		upMethodOrm.EnName = enName
		upMethodOrm.ExtName = extName
		upMethodOrm.FileSelf = byFile
	case upload.Description:
		upMethodOrm.DescriptionEnName = enName
		upMethodOrm.DescriptionExtName = extName
		upMethodOrm.UpDescription = byFile
	}

	return
}

func getEnNameAndExtName(fileNames []string) (enName, extName string, err error) {
	if len(fileNames) == 0 {
		err = errors.New("up file error")
		return
	}
	if len(fileNames) > 3 {
		err = errors.New("fileName can not contain point")
		return
	}

	for i := range fileNames {
		if i == 0 {
			enName = fileNames[i]
		} else {
			extName += constants.Point + fileNames[i]
		}
	}

	return
}

func AuthInfoDecode(c *gin.Context) (authInfos model.AuthInfos, err error) {
	err = json.NewDecoder(c.Request.Body).Decode(&authInfos)
	return authInfos, err
}

func (d *DefinitionImpl) UploadDefinitionMethod(upMethodOrm model.UPMethodOrm) (id int32, err error) {
	// 验证方法英文名称重复与否
	err = d.ValidMethodEnName(upMethodOrm)
	if err != nil {
		return
	}

	// 获取method唯一标识id
	categoryProxy := CreateAtomicIdProxy(nil, nil, d.CategoryImpl, nil)
	id, err = categoryProxy.GetAtomicId(atoi.Method)
	if err != nil {
		return
	}

	// 读取fileSelf，写到本地method存储磁盘
	methodSavePath, mainEngine, err := CreateMethodCompressFile(upMethodOrm.CreateId, upMethodOrm.EnName, upMethodOrm.ExtName, upMethodOrm.FileSelf)
	if err != nil {
		go CashRemove(methodSavePath, upMethodOrm.EnName)
		return
	}

	// 上传算子说明文档
	methodDesc, err := CreateMethodDescriptionFile(upMethodOrm.UpDescription, upMethodOrm.DescriptionEnName, upMethodOrm.DescriptionExtName)
	if err != nil {
		return
	}

	// 更新method_manage表
	methodBd := model.MethodOrm{}
	utils.ReflectStructAssign(&upMethodOrm, &methodBd)

	// 赋值
	methodBd.Id = id
	methodBd.Type = method.MyMethod
	methodBd.MainEngine = mainEngine
	methodBd.CreateTime = utils.TimeToString(time.Now())
	methodBd.Description = methodDesc

	// 默认创建者有全部权限
	userAuths := make([]model.UserAuth, 1)
	userAuths[0] = model.UserAuth{UserId: upMethodOrm.CreateId, ControlAspects: []int32{method.Completed, method.Deleted, method.Used}}
	methodBd.AuthInfo = model.AuthInfo{UserAuth: userAuths}
	err = d.MethodImpl.CreateMethod(methodBd)
	if err != nil {
		return
	}

	// 更新method_category表
	categoryOrm := model.CategoryOrm{Id: id, Name: upMethodOrm.CnName, ParentId: upMethodOrm.CateParentId,
		CreatorId: upMethodOrm.CreateId, CreatorName: upMethodOrm.CreateName, Level: filej.MethodType}
	err = d.CategoryImpl.CreateCategory(categoryOrm)
	if err != nil {
		return
	}

	return
}

func (d *DefinitionImpl) UploadDefinitionMethodUpdate(methodId, userId int32, upMethodOrm model.UPMethodOrm) (err error) {
	super, err := CertainUserSuperManage(utils.IntToString(int(userId)))
	if err != nil {
		return
	}

	methodInfo, err := d.MethodImpl.SelectMethodById(methodId)
	if err != nil {
		return
	}

	if !super && methodInfo.CreateId != userId {
		err = errors.New("you have not permission update it")
		return
	}

	var mainEngine, methodDesc string

	// 读取fileSelf，写到本地method存储磁盘
	if upMethodOrm.FileSelf != nil {
		_, mainEngine, err = CreateMethodCompressFile(upMethodOrm.CreateId, upMethodOrm.EnName, upMethodOrm.ExtName, upMethodOrm.FileSelf)
		if err != nil {
			return
		}
	}

	// 上传算子说明文档
	if upMethodOrm.UpDescription != nil {
		methodDesc, err = CreateMethodDescriptionFile(upMethodOrm.UpDescription, upMethodOrm.DescriptionEnName, upMethodOrm.DescriptionExtName)
		if err != nil {
			return
		}
	}

	// 更新元信息
	methodBd := model.MethodOrm{}
	utils.ReflectStructAssign(&methodInfo, &methodBd)
	utils.ReflectStructAssign(&upMethodOrm, &methodBd)
	if upMethodOrm.UpDescription == nil {
		nginxHead := config.NginxConfig.Head + config.NginxConfig.Host + constants.Colon + config.NginxConfig.Port
		methodBd.Description = methodInfo.Description[len(nginxHead):]
	} else {
		methodBd.Description = methodDesc
	}

	if upMethodOrm.FileSelf == nil {
		methodBd.EnName = methodInfo.EnName
		methodBd.MainEngine = methodInfo.MainEngine
	} else {
		methodBd.EnName = upMethodOrm.EnName
		methodBd.MainEngine = mainEngine
	}

	err = d.MethodImpl.UpdateMethod(methodBd)
	if err != nil {
		return
	}

	// 更新category
	categoryOrm := model.CategoryOrm{}
	categoryOrm.Id = methodInfo.Id
	categoryOrm.Name = upMethodOrm.CnName
	categoryOrm.CreatorId = userId
	err = d.CategoryImpl.UpdateCategory(categoryOrm)
	if err != nil {
		return
	}

	return
}

func (d *DefinitionImpl) AssignDefinitionMethodAuthority(methodId, userId string, authInfos model.AuthInfo) (err error) {
	methodIdInt32, err := utils.StringToint32(methodId)
	if err != nil {
		return
	}

	userIdInt32, err := utils.StringToint32(userId)
	if err != nil {
		return
	}

	super, err := CertainUserSuperManage(userId)
	if err != nil {
		return
	}

	auth := d.AuthOrNotAccess(methodIdInt32, userIdInt32, method.Completed)

	if !auth && !super {
		err = errors.New("this method " + methodId + " permission denied")
		return
	}

	err = d.MethodImpl.UpdateMethodAuth(methodIdInt32, authInfos)
	if err != nil {
		return
	}

	return
}

func (d *DefinitionImpl) SelectDefinitionMethodAuthority(methodId, userId string) (authInfos model.AuthInfo, err error) {
	methodIdInt32, err := utils.StringToint32(methodId)
	if err != nil {
		return
	}

	userIdInt32, err := utils.StringToint32(userId)
	if err != nil {
		return
	}

	super, err := CertainUserSuperManage(userId)
	if err != nil {
		return
	}

	auth := d.AuthOrNotAccess(methodIdInt32, userIdInt32, method.Completed)

	if !auth && !super {
		err = errors.New("this method " + methodId + " permission denied")
		return
	}

	methodInfo, err := d.MethodImpl.SelectMethodById(methodIdInt32)
	if err != nil {
		return
	}

	authInfos = methodInfo.AuthInfo
	return
}

func (d *DefinitionImpl) DeleteDefinitionMethod(methodId, userId string) (err error) {
	methodIdInt32, err := utils.StringToint32(methodId)
	if err != nil {
		return
	}

	userIdInt32, err := utils.StringToint32(userId)
	if err != nil {
		return
	}

	super, err := CertainUserSuperManage(userId)
	if err != nil {
		return
	}

	auth := d.AuthOrNotAccess(methodIdInt32, userIdInt32, method.Deleted) ||
		d.AuthOrNotAccess(methodIdInt32, userIdInt32, method.Completed)

	if !auth && !super {
		err = errors.New("this method " + methodId + " is permission denied")
		return
	}

	methodInfo, err := d.MethodImpl.SelectMethodById(methodIdInt32)
	if err != nil {
		return
	}
	methodsMark := getMethodsMark()
	userIdString := utils.IntToString(int(methodInfo.CreateId))
	switch runtime.GOOS {
	case plat.Windows:
		go CashRemove(plat.WindowsAIConfigContext+methodsMark+plat.WindowsSpiltRex+userIdString, methodInfo.EnName)
	case plat.Linux:
		go CashRemove(config.TrainPlatConfig.MountContext+plat.LinuxAIConfigContext+methodsMark+plat.LinuxSpiltRex+userIdString, methodInfo.EnName)
	}

	err = d.CategoryImpl.DeleteCategory(methodIdInt32, userIdInt32)
	if err != nil {
		return
	}

	err = d.MethodImpl.DeleteMethod(methodIdInt32)
	if err != nil {
		return
	}

	return
}

func (d *DefinitionImpl) AuthOrNotAccess(methodId, userId, authType int32) (auth bool) {
	// 首先要查询userId是否可操作methodId
	methodInfo, err := d.MethodImpl.SelectMethodById(methodId)
	if err != nil {
		return
	}

	if methodInfo.DeleteMark {
		return
	}

	if methodInfo.Type == method.PubMethod {
		err = errors.New("this method is public, so you can not operate it")
		return
	}

	if methodInfo.CreateId == userId { // 创建者有完全控制其权限
		auth = true
		return
	}

	userAuth := d.UserAuthValid(userId, authType, methodInfo)
	userGroupAuth := d.UserGroupAuthValid(userId, authType, methodInfo)

	return userAuth || userGroupAuth
}

func (d *DefinitionImpl) UserAuthValid(userId, authType int32, methodInfo model.MethodOrm) (userAuth bool) {
	for _, info := range methodInfo.AuthInfo.UserAuth { // todo 用户组未校验
		tempAuth := false

		switch authType {
		case method.Completed:
			tempAuth = utils.Int32ContainKeyOrNot(method.Completed, info.ControlAspects)
		case method.Deleted:
			tempAuth = utils.Int32ContainKeyOrNot(method.Deleted, info.ControlAspects)
		case method.Used:
			tempAuth = utils.Int32ContainKeyOrNot(method.Used, info.ControlAspects)
		}

		if !tempAuth {
			continue
		}

		if info.UserId == userId {
			userAuth = true
			break
		}
	}
	return
}

func (d *DefinitionImpl) UserGroupAuthValid(userId, authType int32, methodInfo model.MethodOrm) (userGroupAuth bool) {
	for _, info := range methodInfo.AuthInfo.UserGroupsAuth {
		userGroupReq := impl.UserGroupsReq{GroupId: info.UserGroupId, PlatformCode: constants.JinnPlatformCode}
		query := make(url.Values)
		resp, err := impl.CreateRemoteReq(&userGroupReq).RequestRemote(query, nil)
		if err != nil {
			return
		}

		groupTempAuth := false

		switch authType {
		case method.Completed:
			groupTempAuth = utils.Int32ContainKeyOrNot(method.Completed, info.ControlAspects)
		case method.Deleted:
			groupTempAuth = utils.Int32ContainKeyOrNot(method.Deleted, info.ControlAspects)
		case method.Used:
			groupTempAuth = utils.Int32ContainKeyOrNot(method.Used, info.ControlAspects)
		}

		if !groupTempAuth {
			continue
		}

		UserGroupContext := resp.(mocks.UserGroupContext)
		for _, userGroupData := range UserGroupContext.Data {
			if userGroupData.UserAccount == utils.IntToString(int(userId)) {
				userGroupAuth = true
				break
			}
		}

		if userGroupAuth {
			break
		}
	}
	return
}

func (d *DefinitionImpl) ValidMethodEnName(upMethodOrm model.UPMethodOrm) (err error) {
	methodInfos, err := d.SelectMyMethodByCreateId(upMethodOrm.CreateId)
	if err != nil {
		return
	}
	for _, methodInfo := range methodInfos {
		if methodInfo.EnName == upMethodOrm.EnName && !methodInfo.DeleteMark {
			err = errors.New("current method compress package name is repeated")
			return
		}
	}
	return
}

func CreateMethodCompressFile(userId int32, fileEnName, extName string, srcFile multipart.File) (methodSavePath, mainEngine string, err error) {
	if extName != compress.Tar && extName != compress.TarGz && extName != compress.Zip {
		err = errors.New("only support .tar .tar.gz and .zip compressed package")
		return
	}

	defer srcFile.Close()

	user := utils.IntToString(int(userId))

	var fileStorePath string

	methodsMark := getMethodsMark()

	switch runtime.GOOS {
	case plat.Windows:
		methodSavePath = plat.WindowsAIConfigContext + methodsMark + plat.WindowsSpiltRex + user
		mainEngine = methodSavePath + plat.WindowsSpiltRex + fileEnName + plat.WindowsSpiltRex + constants.RunEntry
		fileStorePath = methodSavePath + plat.WindowsSpiltRex + fileEnName + extName
	case plat.Linux:
		methodSavePath = config.TrainPlatConfig.MountContext + plat.LinuxAIConfigContext + methodsMark + plat.LinuxSpiltRex + user
		containerSavePath := config.TrainPlatConfig.MountTrainContext + plat.LinuxAIConfigContext + methodsMark + plat.LinuxSpiltRex + user
		mainEngine = containerSavePath + plat.LinuxSpiltRex + fileEnName + plat.LinuxSpiltRex + constants.RunEntry
		fileStorePath = methodSavePath + plat.LinuxSpiltRex + fileEnName + extName
	}

	err = os.MkdirAll(methodSavePath, os.ModePerm)
	if err != nil {
		return
	}

	destFile, err := os.Create(fileStorePath)
	if err != nil {
		return
	}
	defer destFile.Close()

	// 写入
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return
	}

	// 解压
	err = UnCompress(extName, methodSavePath, fileEnName)
	if err != nil { // todo 上传失败可能需要清除该文件夹
		return
	}

	return
}

func CreateMethodDescriptionFile(srcFile multipart.File, fileEnName, extName string) (descPos string, err error) {
	if extName != constants.PdfSuffix {
		err = errors.New("upload description file must be .pdf format")
		return
	}

	descFilePos := ""
	switch runtime.GOOS {
	case plat.Windows:
		descPos = plat.WindowsAIConfigContext + fileEnName + extName
		descFilePos = plat.WindowsAIConfigContext + fileEnName + extName
	case plat.Linux:
		descPos = config.NginxConfig.DescContext + plat.LinuxSpiltRex + fileEnName + extName
		descFilePos = config.NginxConfig.SaveContext + plat.LinuxSpiltRex + fileEnName + extName
	}

	destFile, err := os.Create(descFilePos)
	if err != nil {
		return
	}
	defer destFile.Close()

	// 写入
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return
	}

	return
}

func UnCompress(extName, filePath, fileName string) (err error) {
	switch extName {
	case compress.Zip:
		err = utils.UnZipFile(filePath, fileName+extName)
	case compress.Tar:
		err = utils.UnTarFile(filePath, fileName+extName)
	case compress.TarGz:
		err = utils.UnTarGzFile(filePath, fileName+extName)
	default:
		err = errors.New("not support this unCompress format")
		return
	}
	return
}

func CashRemove(methodSavePath, fileEnName string) {
	fileInfos, err := ioutil.ReadDir(methodSavePath)
	if err != nil {
		return
	}
	for _, fileInfo := range fileInfos {
		if fileInfo.Name() == fileEnName ||
			fileInfo.Name() == fileEnName+compress.Zip ||
			fileInfo.Name() == fileEnName+compress.TarGz ||
			fileInfo.Name() == fileEnName+compress.Tar {
			_ = os.RemoveAll(methodSavePath + string(os.PathSeparator) + fileInfo.Name())
		}
	}
	return
}
