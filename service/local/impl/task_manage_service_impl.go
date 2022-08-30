package impl

import (
	"bufio"
	"dataProcess/cache/redis"
	config2 "dataProcess/config"
	"dataProcess/constants"
	"dataProcess/constants/atoi"
	"dataProcess/constants/plat"
	"dataProcess/constants/source"
	"dataProcess/constants/statuses"
	"dataProcess/constants/table"
	"dataProcess/constants/train"
	"dataProcess/model"
	"dataProcess/orm"
	"dataProcess/repository/mocks"
	"dataProcess/service/remote/impl"
	dbDo "dataProcess/tcp_socket/db_do"
	"dataProcess/tools/utils/utils"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"net/url"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"
)

var headMap map[string]string

type TaskManageImpl struct {
	coll *mgo.Collection
}

func CreateTaskManageImpl() *TaskManageImpl {
	if orm.MongoDb == nil {
		mongoDb := new(dbDo.MongodbServer)
		orm.MongoDb = mongoDb.SetUp()
	}
	return &TaskManageImpl{coll: orm.MongoDb.C(table.TaskManage)}
}

func getTasksMark() (tasksMark string) {
	switch config2.ApplicationConfig.SrcSource {
	case constants.Formalization:
		tasksMark = constants.Tasks + constants.SubLine + constants.Formalization
	case constants.Testing:
		tasksMark = constants.Tasks
	case constants.Integration:
		tasksMark = constants.Tasks + constants.SubLine + constants.Integration
	case constants.Security:
		tasksMark = constants.Tasks + constants.SubLine + constants.Security
	}
	return
}

func TaskPostDecode(c *gin.Context) (taskOrm model.TaskOrm, err error) {
	err = json.NewDecoder(c.Request.Body).Decode(&taskOrm)
	return taskOrm, err
}

func TaskCompletePostDecode(c *gin.Context) (taskSaveOrm model.TaskSaveOrm, err error) {
	err = json.NewDecoder(c.Request.Body).Decode(&taskSaveOrm)
	return taskSaveOrm, err
}

func TaskGetDecode(c *gin.Context) (queryTask mocks.QueryTask, err error) {
	err = json.NewDecoder(c.Request.Body).Decode(&queryTask)
	if err != nil {
		return
	}
	return
}

func (t *TaskManageImpl) QueryTaskDetailById(id string) (taskD model.TaskSaveOrm, idInt32 int32, err error) {
	idInt32, err = utils.StringToint32(id)
	if err != nil {
		return
	}

	// 查询任务id信息
	taskD, err = t.QueryTaskDetail(idInt32)
	if err != nil {
		return
	}
	return
}

func (t *TaskManageImpl) StartTask(id string, userId string) (err error) {
	taskD, idInt32, err := t.QueryTaskDetailById(id)
	if err != nil {
		return
	}

	userIntId32, err := utils.StringToint32(userId)
	if err != nil {
		return
	}

	if taskD.CreateId != userIntId32 {
		err = errors.New("task id: " + id + " permission denied")
		return
	}

	if taskD.Status == statuses.Pending || taskD.Status == statuses.Running ||
		taskD.Status == statuses.Terminal {
		err = errors.New("this task " + id + " status: " + taskD.Status + " so you not start it")
		return
	}

	// 创建任务id的文件夹，并将外层AIFlow_config.py追加“CMD”命令，并放入任务id的文件夹下，最后再创建src文件夹
	err = PlatWriteAIFlowPy(id, GetToAppendCmd(id))
	if err != nil {
		return
	}

	// 登录
	err = GetLoginHead(constants.TrainProposer)
	if err != nil {
		return
	}

	// 写回数据库SeqId
	err = t.WriteToInputPathSeqId(idInt32, headMap)
	if err != nil {
		return
	}

	//todo 取taskD中所需GPU和CPU的最大值
	maxCpuNum, maxGpuNum, err := GetMaxTaskGpuAndCpu(taskD)
	if err != nil {
		return
	}

	// 发送启动train任务请求
	jobInfo, err := StartTrainContainer(id, taskD, headMap, maxCpuNum, maxGpuNum)
	if err != nil {
		return
	}

	if !jobInfo.Result {
		err = errors.New("start task failed")
		return
	}

	// 写回该任务的trainId
	err = t.UpdateTaskTrainId(model.CreateTaskSaveOrmTrain(idInt32, jobInfo.MissionCreatedId))
	if err != nil {
		return
	}

	// 写回任务新建完成状态
	t.UpdateTaskStatusEnter(idInt32, statuses.Pending, "", "", "")

	return
}

func (t *TaskManageImpl) CopyTask(id, userId, userName string) (newId int32, err error) {
	taskD, _, err := t.QueryTaskDetailById(id)
	if err != nil {
		return
	}

	userIdStr, err := utils.StringToint32(userId)
	if err != nil {
		return
	}

	// 查原型
	protoServer := CreateTaskPrototypeImpl()
	proto, err := protoServer.SelectProto(taskD.TaskProtoTypeId)
	if err != nil {
		return
	}

	// 查询taskId
	taskManageProxy := CreateAtomicIdProxy(t, nil, nil, nil)
	newId, err = taskManageProxy.GetAtomicId(atoi.Task)
	if err != nil {
		return
	}

	// 查询protoId
	retSaveProto, err := protoServer.SaveProto(proto)
	if err != nil {
		return
	}

	// 复制任务
	taskD.Id = newId
	taskD.Status = statuses.Newed
	taskD.UId = uuid.New().String()
	taskD.CreateId = userIdStr
	taskD.CreateName = userName
	taskD.CreateTime = utils.TimeToString(time.Now())
	taskD.TaskProtoTypeId = retSaveProto.Id
	taskD.EndTime = ""
	err = t.coll.Insert(taskD)
	if err != nil {
		return
	}

	return
}

func (t *TaskManageImpl) LogTask(id string) (taskLog string, err error) {
	taskD, _, err := t.QueryTaskDetailById(id)
	if err != nil {
		return
	}

	logFilePath, err := GetLogFilePath(id, taskD.UId, taskD.TrainRelatedId)
	if err != nil {
		return
	}

	err = GetLoginHead(constants.TrainProposer)
	if err != nil {
		return
	}

	taskTrainLogResp, err := GetTaskLog(logFilePath, headMap)
	if err != nil {
		return
	}

	if !taskTrainLogResp.Result {
		err = errors.New(taskTrainLogResp.ErrorMsg)
		return
	}

	taskLog = taskTrainLogResp.FileContent

	return
}

func (t *TaskManageImpl) DeleteTask(id, userId int32) (err error) {
	taskDetail, err := t.QueryTaskDetail(id)
	if err != nil {
		return
	}

	if taskDetail.CreateId != userId {
		err = errors.New("permission denied")
		return
	}
	err = t.coll.Remove(bson.M{"id": id})
	if err != nil {
		return
	}

	return
}

func (t *TaskManageImpl) TerminalTask(id string) (err error) {
	taskD, idInt32, err := t.QueryTaskDetailById(id)
	if err != nil {
		return
	}

	if taskD.Status == statuses.Failed || taskD.Status == statuses.Error ||
		taskD.Status == statuses.Newed || taskD.Status == statuses.Succeed ||
		taskD.Status == statuses.Terminal {
		err = errors.New("this task " + id + " status: " + taskD.Status + " so you not terminal it")
		return
	}

	err = t.DeleteRemoteTask(id)
	if err != nil {
		return
	}

	// 写回任务terminal状态
	t.UpdateTaskStatusEnter(idInt32, statuses.Terminal, "", "", "")

	return
}

func (t *TaskManageImpl) DeleteRemoteTask(id string) (err error) {
	taskD, _, err := t.QueryTaskDetailById(id)
	if err != nil {
		return
	}

	// 终止任务需要发送train_related_id、id和userId，直接删除
	err = GetLoginHead(taskD.CreateId)
	if err != nil {
		return
	}

	// 直接循环执行终止任务，避免查询
	var taskOp mocks.TaskOperateResp
	for _, rotStatue := range statuses.RoeStatues {
		taskOp, err = TaskOperateDel(config2.JinnConfig.Username, taskD.TrainRelatedId, rotStatue, headMap)
		if err != nil {
			return
		}
		if !taskOp.Result {
			continue
		} else {
			break
		}
	}
	return
}

func (t *TaskManageImpl) SaveTask(taskOrm model.TaskOrm) (id int32, err error) {
	var taskOrmSave model.TaskSaveOrm

	taskManageProxy := CreateAtomicIdProxy(t, nil, nil, nil)
	id, err = taskManageProxy.GetAtomicId(atoi.Task)
	if err != nil {
		return
	}

	if len(taskOrm.Arguments) == 0 {
		err = errors.New("current task of methods is empty")
		return
	}

	err = ValidDealMaterialType(taskOrm)
	if err != nil {
		return
	}

	err = ValidDealClassification(taskOrm)
	if err != nil {
		return
	}

	// 获取原子id后赋值, 并自动生成createTime
	taskOrmSave.Id = id
	taskOrmSave.UId = uuid.New().String()
	taskOrmSave.CreateTime = utils.TimeToString(time.Now())
	utils.ReflectStructAssign(&taskOrm, &taskOrmSave)

	if len(taskOrmSave.InputPath) != 1 {
		taskOrmSave.InputPath, err = FilterJinnPathsUsingPrefixTree(taskOrmSave.InputPath)
		if err != nil {
			return
		}
	}

	err = t.coll.Insert(&taskOrmSave)
	if err != nil {
		return
	}

	// 写回任务新建完成状态
	t.UpdateTaskStatusEnter(id, statuses.Newed, "", "", "")

	return
}

func (t *TaskManageImpl) UpdateTaskStatus(taskSaveOrm model.TaskSaveOrm) (err error) {
	err = t.coll.Update(&bson.M{"id": taskSaveOrm.Id}, bson.M{"$set": bson.M{"statuses": taskSaveOrm.Status,
		"stage": taskSaveOrm.Stage, "exception": taskSaveOrm.Exception, "endTime": taskSaveOrm.EndTime}})
	if err != nil {
		return
	}
	return
}

func (t *TaskManageImpl) UpdateTaskTrainId(taskSaveOrm model.TaskSaveOrm) (err error) {
	err = t.coll.Update(&bson.M{"id": taskSaveOrm.Id}, bson.M{"$set": bson.M{"trainRelatedId": taskSaveOrm.TrainRelatedId}})
	if err != nil {
		return
	}
	return
}

func (t *TaskManageImpl) UpdateTaskInputPathSeqId(taskSaveOrm model.TaskSaveOrm) (err error) {
	err = t.coll.Update(&bson.M{"id": taskSaveOrm.Id}, bson.M{"$set": bson.M{"inputPath": taskSaveOrm.InputPath}})
	if err != nil {
		return
	}
	return
}

func (t *TaskManageImpl) UpdateTaskInputMarkVersionSeqId(taskSaveOrm model.TaskSaveOrm) (err error) {
	err = t.coll.Update(&bson.M{"id": taskSaveOrm.Id}, bson.M{"$set": bson.M{"inputMarkVersion": taskSaveOrm.InputMarkVersion}})
	if err != nil {
		return
	}
	return
}

func (t *TaskManageImpl) UpdateTask(id int32, taskOrm model.TaskOrm) (err error) {
	err = ValidDealMaterialType(taskOrm)
	if err != nil {
		return
	}

	err = ValidDealClassification(taskOrm)
	if err != nil {
		return
	}

	err = t.coll.Update(&bson.M{"id": id},
		bson.M{"$set": bson.M{
			"name":                     taskOrm.Name,
			"inputType":                taskOrm.InputType,
			"inputPath":                taskOrm.InputPath,
			"outputType":               taskOrm.OutputType,
			"outputPath":               taskOrm.OutputPath,
			"arguments":                taskOrm.Arguments,
			"downloadContent":          taskOrm.DownloadContent,
			"materialType":             taskOrm.MaterialType,
			"srcDirName":               taskOrm.SrcDirName,
			"updateTime":               utils.TimeToString(time.Now()),
			"description":              taskOrm.Description,
			"inputMarkVersion":         taskOrm.InputMarkVersion,
			"outputDataSetVersionName": taskOrm.OutputDataSetVersionName,
			"outputMarkVersion":        taskOrm.OutputMarkVersion,
		}})
	if err != nil {
		return
	}
	return
}

func (t *TaskManageImpl) QueryTaskDetail(id int32) (taskSaveOrm model.TaskSaveOrm, err error) {
	err = t.coll.Find(bson.M{"id": id}).One(&taskSaveOrm)
	return
}

func (t *TaskManageImpl) SearchTask(queryTask *mocks.QueryTask) (taskSaveOrms []model.TaskSaveOrm, err error) {
	// 如果是素材管理的超级管理员，则能查看所有任务
	super, err := CertainUserSuperManage(utils.IntToString(int(queryTask.UserId)))
	if err != nil {
		return
	}

	var query []bson.M
	if queryTask.Id != 0 {
		query = append(query, bson.M{"id": queryTask.Id})
	}
	if len(queryTask.Name) != 0 {
		query = append(query, bson.M{"name": bson.M{"$regex": bson.RegEx{Pattern: queryTask.Name, Options: "im"}}})
	}
	if len(queryTask.Statuses) != 0 {
		query = append(query, bson.M{"statuses": bson.M{"$in": queryTask.Statuses}})
	}
	if len(queryTask.BeginTime) != 0 && len(queryTask.EndTime) != 0 {
		query = append(query, bson.M{"createTime": bson.M{"$gte": queryTask.BeginTime}},
			bson.M{"createTime": bson.M{"$lte": queryTask.EndTime}})
	}

	if !super {
		if len(queryTask.Creators) != 0 {
			query = append(query, bson.M{"createId": bson.M{"$in": queryTask.Creators}})
		} else {
			query = append(query, bson.M{"createId": queryTask.UserId})
		}
	} else {
		if len(queryTask.Creators) != 0 {
			query = append(query, bson.M{"createId": bson.M{"$in": queryTask.Creators}})
		}
	}

	if len(query) == 0 {
		err = t.coll.Find(nil).Sort("-createTime", "-id").All(&taskSaveOrms)
		if err != nil {
			return
		}
		return
	}

	err = t.coll.Find(bson.M{"$and": query}).Sort("-createTime", "-id").All(&taskSaveOrms)
	if err != nil {
		return
	}

	return
}

func (t *TaskManageImpl) UpdateTaskStatusEnter(id int32, status, stage, errMsg, endTime string) {
	taskUp := model.CreateTaskSaveOrmRecall(id, status, stage, errMsg, endTime)

	switch status {
	case statuses.Newed:
		taskUp.Status = statuses.Newed
	case statuses.Succeed:
		taskUp.Status = statuses.Succeed
	case statuses.Running:
		taskUp.Status = statuses.Running
	case statuses.Failed:
		taskUp.Status = statuses.Failed
	case statuses.Error:
		taskUp.Status = statuses.Error
	case statuses.Pending:
		taskUp.Status = statuses.Pending
	case statuses.Terminal:
		taskUp.Status = statuses.Terminal
	default:
		taskUp.Status = ""
	}

	_ = t.UpdateTaskStatus(taskUp)
}

func (t *TaskManageImpl) WriteToInputPathSeqId(IdInt32 int32, headMap map[string]string) (err error) {
	taskDetail, err := t.QueryTaskDetail(IdInt32)
	if err != nil {
		return
	}
	var dataSetResp mocks.DataSetResp

	if taskDetail.InputType == source.DataSet {
		// 发送数据集查询请求
		writeToInputPaths := make([]model.InputPaths, len(taskDetail.InputPath))
		for index, inputPath := range taskDetail.InputPath {
			dataSetResp, err = DataSetQuery(inputPath.SolidId, headMap)
			if err != nil {
				return
			}
			writeToInputPaths[index] = model.InputPaths{InPath: inputPath.InPath, SolidId: inputPath.SolidId,
				SeqId: utils.IntToString(dataSetResp.RealId), PrefixPath: inputPath.PrefixPath}
		}
		err = t.UpdateTaskInputPathSeqId(model.CreateTaskSaveOrmReplyDataSet(IdInt32, writeToInputPaths))
		if err != nil {
			return
		}
	}

	return
}

/**
前缀树方法过滤jinnPaths
*/
func FilterJinnPathsUsingPrefixTree(inputPaths []model.InputPaths) (retPaths []model.InputPaths, err error) {

	jinnPaths := GetJinnPaths(inputPaths)

	utils.CreatePrefixTreeRoot(nil, "", 0, 0)
	if utils.Root == nil {
		err = errors.New("init prefix tree root node is failed")
		return
	}

	err = utils.Root.InsertNodeToPrefixTree(jinnPaths)
	if err != nil {
		return
	}

	nodePaths := new([]string)
	initPaths := make([]string, 0)
	err = utils.Root.FindMaxPassNodeWithEnd(nodePaths, initPaths)
	if err != nil {
		return
	}

	retPaths = FilterJinnPaths(inputPaths, *nodePaths)

	return
}

/**
传统方法过滤jinnPaths
*/
func FilterJinnPathsEntry(inputPaths []model.InputPaths) (retPaths []model.InputPaths, err error) {
	jinnPaths := GetJinnPaths(inputPaths)
	newInputPath, err := ValidateAndIntegrateJinnPaths(jinnPaths)
	if err != nil {
		return
	}
	retPaths = FilterJinnPaths(inputPaths, newInputPath)
	return
}

func FilterJinnPaths(inputPaths []model.InputPaths, newInputPaths []string) (retInputPath []model.InputPaths) {
	for _, inputPath := range inputPaths {
		for _, newInputPath := range newInputPaths {
			if inputPath.InPath == newInputPath {
				retInputPath = append(retInputPath, model.InputPaths{InPath: newInputPath, SeqId: inputPath.SeqId,
					SolidId: inputPath.SolidId, PrefixPath: inputPath.PrefixPath})
				break
			}
		}
	}
	return
}

func GetJinnPaths(inputPaths []model.InputPaths) (jinnPaths []string) {
	for _, inputPath := range inputPaths {
		jinnPaths = append(jinnPaths, inputPath.InPath)
	}
	return
}

func ValidateAndIntegrateJinnPaths(jinnPaths []string) (newJinnPaths []string, err error) {
	if len(jinnPaths) == 0 || len(jinnPaths) > 20 {
		err = errors.New("inputPath length is not less 0 and not over 20")
		return
	}
	sort.Strings(jinnPaths)
	sentiSpace, sentiContext, err := utils.SplitJinnString(jinnPaths[0])
	if err != nil {
		return
	}
	contexts := make([]string, 0)
	contexts = append(contexts, sentiContext)
	var space, context string
	for index, jinnPath := range jinnPaths {
		if index == 0 {
			continue
		}
		space, context, err = utils.SplitJinnString(jinnPath)
		if err != nil {
			return
		}
		if sentiSpace != space {
			err = errors.New("space is not different")
			return
		}
		contexts = append(contexts, context)
	}
	newJinnPaths = IntegrateJinnPath(contexts, jinnPaths)
	return
}

/**
如果复杂度过高可以利用preTree解决
*/
func IntegrateJinnPath(contexts, jinnPaths []string) (newJinnPaths []string) {
	for i := 0; i < len(contexts)-1; i += 2 {
		var saveJinnPathI, saveJinnPathJ string
		j := i + 1
		equal, leer := JudgeTwoContextPrefix(contexts[i], contexts[j])
		if !equal { // 全部保留
			saveJinnPathI = jinnPaths[i]
			saveJinnPathJ = jinnPaths[j]
		} else { // 保留短的
			switch leer {
			case 1:
				saveJinnPathI = jinnPaths[j]
			case 2:
				saveJinnPathI = jinnPaths[i]
			}
		}
		if saveJinnPathI != "" {
			existJ := PathExistOrNot(newJinnPaths, saveJinnPathI)
			if !existJ {
				newJinnPaths = append(newJinnPaths, saveJinnPathI)
			}
		}
		if saveJinnPathJ != "" {
			existJ := PathExistOrNot(newJinnPaths, saveJinnPathJ)
			if !existJ {
				newJinnPaths = append(newJinnPaths, saveJinnPathJ)
			}
		}
	}
	return
}

func PathExistOrNot(newJinnPaths []string, comparePath string) (exist bool) {
	if len(newJinnPaths) == 0 {
		return false
	}
	for _, validPath := range newJinnPaths {
		equal, _ := JudgeTwoContextPrefix(comparePath, validPath)
		if equal { // 已经存在更短的前缀
			return true
		}
	}
	return false
}

func JudgeTwoContextPrefix(preContext, postContext string) (equal bool, leer int) {
	var lenContext, shoContext string
	if len(preContext) > len(postContext) {
		lenContext = preContext
		shoContext = postContext
		leer = 1
	} else if len(preContext) < len(postContext) {
		lenContext = postContext
		shoContext = preContext
		leer = 2
	} else {
		if preContext != postContext {
			return false, 0
		} else {
			return true, 0
		}
	}
	lenArr := strings.Split(lenContext, plat.LinuxSpiltRex)
	shoArr := strings.Split(shoContext, plat.LinuxSpiltRex)
	var prefix string
	for index, _ := range shoArr {
		if index == 0 {
			continue
		}
		prefix += plat.LinuxSpiltRex + lenArr[index]
	}
	if shoContext == prefix {
		return true, leer
	} else {
		return false, leer
	}
}

func GetToAppendCmd(taskId string) (Cmd string) {
	srcMark := getSrcMark()
	Cmd = train.CCMD + constants.Quotes + constants.CD + constants.WhiteSpace + config2.TrainPlatConfig.MountTrainContext + plat.LinuxAIConfigContext +
		srcMark + constants.PreFixCMD + constants.WhiteSpace + plat.LinuxClientAppConfigureContext +
		constants.WhiteSpace + taskId + constants.Quotes
	return
}

func getSrcMark() (srcMark string) {
	switch config2.ApplicationConfig.SrcSource {
	case constants.Formalization:
		srcMark = constants.SRC + constants.SubLine + constants.Formalization
	case constants.Testing:
		srcMark = constants.SRC
	case constants.Integration:
		srcMark = constants.SRC + constants.SubLine + constants.Integration
	case constants.Security:
		srcMark = constants.SRC + constants.SubLine + constants.Security
	}
	return
}

func PlatWriteAIFlowPy(taskId, appendCmd string) (err error) {
	var aiContext string
	tasksMark := getTasksMark()
	// 创建src文件夹并写入文件AIFlowPy文件
	switch runtime.GOOS {
	case plat.Windows:
		aiContext = plat.WindowsAIConfigContext + tasksMark + plat.WindowsSpiltRex + taskId
		err = CreateAiConfigPyContext(plat.WindowsAIConfigContext, aiContext,
			aiContext+plat.WindowsSpiltRex+constants.SRC, appendCmd)
	case plat.Linux:
		aiContext = config2.TrainPlatConfig.MountContext + plat.LinuxAIConfigContext + tasksMark + plat.LinuxSpiltRex + taskId
		err = CreateAiConfigPyContext(config2.TrainPlatConfig.MountContext+plat.LinuxAIConfigContext, aiContext,
			aiContext+plat.LinuxSpiltRex+constants.SRC, appendCmd)
	}

	if err != nil {
		return
	}
	return
}

func CreateAiConfigPyContext(sourcePath, aiContextPath, srcContextPath, appendCmd string) (err error) {
	uid, err := utils.StringToint32(config2.JinnConfig.Username)
	if err != nil {
		return
	}
	gid, err := utils.StringToint32(config2.JinnConfig.GroupName)
	if err != nil {
		return
	}

	back, err := utils.PathExists(aiContextPath)
	if err != nil {
		return
	}
	if !back && runtime.GOOS == plat.Linux {
		err = os.MkdirAll(aiContextPath, os.ModePerm)
		if err != nil {
			return
		}
		err = os.Chown(aiContextPath, int(uid), int(gid))
		if err != nil {
			return
		}
		err = os.MkdirAll(srcContextPath, os.ModePerm)
		if err != nil {
			return
		}
		err = os.Chown(srcContextPath, int(uid), int(gid))
		if err != nil {
			return
		}
	}
	switch runtime.GOOS {
	case plat.Windows:
		err = WriterAIFlowPy(sourcePath+train.AIFlowConfigPy, aiContextPath+plat.WindowsSpiltRex+train.AIFlowConfigPy, appendCmd, uid, gid)
	case plat.Linux:
		err = WriterAIFlowPy(sourcePath+train.AIFlowConfigPy, aiContextPath+plat.LinuxSpiltRex+train.AIFlowConfigPy, appendCmd, uid, gid)
	}
	if err != nil {
		return
	}
	return
}

func WriterAIFlowPy(sourceFile, targetPath, appendCmd string, uid, gid int32) (err error) {
	srcFile, err := os.OpenFile(sourceFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return
	}
	defer srcFile.Close()

	taskFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return
	}
	defer taskFile.Close()

	readBuf := bufio.NewReader(srcFile)

	var line []byte
	for {
		line, _, err = readBuf.ReadLine()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		_, err = taskFile.WriteString(string(line) + "\n")
		if err != nil {
			return err
		}
	}

	_, err = taskFile.WriteString(appendCmd + "\n")
	if err != nil {
		return
	}

	if runtime.GOOS == plat.Linux {
		err = os.Chown(sourceFile, int(uid), int(gid))
		if err != nil {
			return
		}

		err = os.Chown(targetPath, int(uid), int(gid))
		if err != nil {
			return
		}
	}

	return
}

func GetLoginHead(createId int32) (err error) {
	if headMap == nil {
		headMap = make(map[string]string)
	}

	headMap[train.ContentKey] = train.ContentValue

	token, err := redis.GetOperate(constants.JinnToken + constants.SubLine + utils.IntToString(int(createId)))
	if err != nil {
		return
	}

	newToken, err := directRefreshToken(token)
	if err != nil {
		return
	}

	// 加入登录的token
	headMap[train.CookieKey] = train.TrainCookiePreValue + newToken
	return
}

func directRefreshToken(oldToken string) (newToken string, err error) {
	headMap = make(map[string]string)
	headMap[train.ContentKey] = train.ContentValue
	newToken, err = RefreshToken(oldToken, headMap)
	if err != nil {
		return
	}
	return
}

func RefreshToken(token string, headMap map[string]string) (newToken string, err error) {
	resp, err := impl.CreateOldTokenContext(token).RequestRemote(nil, headMap)
	if err != nil {
		return
	}

	// assert -> mocks.TokenUpdateContext
	tokenBodyResp := resp.(mocks.TokenUpdateContext)

	if !tokenBodyResp.Result {
		err = errors.New(tokenBodyResp.ErrorMsg)
		return
	}

	newToken = tokenBodyResp.Token
	return
}

func LoginCall(headMap map[string]string) (loginResp mocks.LoginResp, err error) {
	encryptBd := CreateEncryptionBody(SetEncodeText(config2.JinnConfig.Password))
	encryptPasswd, err := encryptBd.Encrypt()
	if err != nil {
		return
	}

	loginBody := impl.LoginReqContext{
		UserAccount: config2.JinnConfig.Username,
		Password:    encryptPasswd,
	}
	resp, err := impl.CreateRemoteReq(&loginBody).RequestRemote(nil, headMap)
	if err != nil {
		return
	}

	// assert -> mocks.LoginResp
	loginResp = resp.(mocks.LoginResp)
	if !loginResp.Result {
		err = errors.New(loginResp.ErrorMsg)
		return
	}
	return
}

/**
获取某个任务中最大所需的GPU数量和CPU数量
*/
func GetMaxTaskGpuAndCpu(taskD model.TaskSaveOrm) (maxCpuNum, maxGpuNum int32, err error) {
	var methodInfo model.MethodOrm

	for _, argument := range taskD.Arguments {
		methodInfo, err = CreateMethodImpl().SelectMethodById(argument.MethodId)
		if err != nil {
			return
		}
		if methodInfo.CpuNums >= maxCpuNum {
			maxCpuNum = methodInfo.CpuNums
		}
		if methodInfo.GpuNums >= maxGpuNum {
			maxGpuNum = methodInfo.GpuNums
		}
	}

	return
}

func CertainUserSuperManage(userId string) (superManage bool, err error) {
	userDetail, err := UserDetailCall(userId)
	if err != nil {
		return
	}

	roleInfos := userDetail.Data.Data.RoleInfo

	if len(roleInfos) == 0 {
		return
	}

	for _, info := range roleInfos {
		if info.RoleName == constants.SuperManager || info.RoleName == constants.NormalManager {
			superManage = true
			break
		}
	}

	return
}

func UserDetailCall(userId string) (userDetailContext mocks.UserDetailContext, err error) {
	userDetailBody := impl.UserDetailReq{UserAccount: userId, PlatFormCode: constants.JinnPlatformCode}
	query := make(url.Values)
	resp, err := impl.CreateRemoteReq(&userDetailBody).RequestRemote(query, headMap)
	if err != nil {
		return
	}

	userDetailContext = resp.(mocks.UserDetailContext)
	if !userDetailContext.Result {
		err = errors.New(userDetailContext.ErrorMsg)
		return
	}
	return
}

func StartTrainContainer(taskId string, taskD model.TaskSaveOrm, headMap map[string]string,
	maxCpuNum, maxGpuNum int32) (trainResp mocks.TrainResp, err error) {

	resourceGroup := impl.ResourceGroup{
		ResGroupName:           config2.TrainPlatConfig.ResGroupName,
		ResGroupLabel:          config2.TrainPlatConfig.ResGroupLabel,
		ResourceGroupCloudId:   config2.TrainPlatConfig.ResourceGroupCloudId,
		ResourceGroupServerCpu: config2.TrainPlatConfig.ResourceGroupServerCpu,
		ResourceCpuOverRation:  config2.TrainPlatConfig.ResourceCpuOverRation,
		ResourceGroupType:      config2.TrainPlatConfig.ResourceGroupType,
	}

	resourceGroups := make([]impl.ResourceGroup, 0)
	resourceGroups = append(resourceGroups, resourceGroup)

	secondStorePathList := make([]string, 0)

	trainContainer := impl.TrainContainerReqContext{
		UserName:             config2.JinnConfig.Username,
		ImageName:            config2.TrainPlatConfig.ImageName,
		ImageTag:             config2.TrainPlatConfig.ImageTag,
		WorkPath:             config2.TrainPlatConfig.MountTrainContext + plat.LinuxAIConfigContext + getTasksMark() + plat.LinuxSpiltRex + taskId,
		ContainerConfig:      &impl.ContainerConfig{GPU: uint32(maxGpuNum), CPU: uint32(maxCpuNum)},
		Dataset:              make([]interface{}, 0),
		UserGroup:            config2.TrainPlatConfig.UserGroup,
		ResourceGroup:        &resourceGroups,
		IsMulti:              config2.TrainPlatConfig.IsMulti,
		IsLoop:               config2.TrainPlatConfig.IsLoop, // optional: false表示开启正常的训练任务，true表示开启debug任务
		NodeToUse:            config2.TrainPlatConfig.NodeToUse,
		Priority:             config2.TrainPlatConfig.Priority,
		SecondStorePathList:  secondStorePathList,
		GpuType:              config2.TrainPlatConfig.GpuType,
		HardwareType:         config2.TrainPlatConfig.HardwareType,
		Tag:                  config2.TrainPlatConfig.Tag,
		MissionName:          taskD.UId, // 这里防止任务名称出现中文
		FrameType:            config2.TrainPlatConfig.FrameType,
		DatasetVersionIdList: make([]interface{}, 0),
		ResourceModel:        config2.TrainPlatConfig.ResourceModel,
	}

	resp, err := impl.CreateRemoteReq(&trainContainer).RequestRemote(nil, headMap)
	if err != nil {
		return
	}
	trainResp = resp.(mocks.TrainResp)
	if !trainResp.Result {
		err = errors.New(trainResp.ErrorMsg)
		return
	}
	return
}

func DataSetQuery(solidId int32, headMap map[string]string) (dataSetResp mocks.DataSetResp, err error) {
	dataSetReq := impl.CreateRemoteReq(&impl.DataSetReqContext{SolidId: solidId})
	query := make(url.Values)
	resp, err := dataSetReq.RequestRemote(query, headMap)
	if err != nil {
		return
	}
	dataSetResp = resp.(mocks.DataSetResp)
	if !dataSetResp.Result {
		err = errors.New(dataSetResp.ErrorMsg)
		return
	}
	return
}

func TaskOperateDel(userId, missionCreatedId, missionStatus string, headMap map[string]string) (taskOperateResp mocks.TaskOperateResp, err error) {
	taskOperateReq := impl.TaskOperateReq{UserName: userId, Operation: "DELETE", MissionCreatedId: missionCreatedId, MissionStatus: missionStatus}
	resp, err := impl.CreateRemoteReq(&taskOperateReq).RequestRemote(nil, headMap)
	if err != nil {
		return
	}

	taskOperateResp = resp.(mocks.TaskOperateResp)
	if !taskOperateResp.Result {
		err = errors.New(taskOperateResp.ErrorMsg)
		return
	}
	return
}

func GetTaskLog(filePosition string, headMap map[string]string) (taskTrainLogResp mocks.TaskTrainLogResp, err error) {
	taskLogReq := impl.TrainLogReq{FilePath: filePosition}
	query := make(url.Values)
	resp, err := impl.CreateRemoteReq(&taskLogReq).RequestRemote(query, headMap)
	if err != nil {
		return
	}

	taskTrainLogResp = resp.(mocks.TaskTrainLogResp)
	if !taskTrainLogResp.Result {
		err = errors.New(taskTrainLogResp.ErrorMsg)
		return
	}
	return
}

func GetLogFilePath(taskId string, taskName string, trainId string) (logFilePath string, err error) {
	trainIdArr := strings.Split(trainId, constants.Sub)
	if len(trainIdArr) == 0 || len(trainIdArr) < 2 {
		err = errors.New("trainId is invalid")
		return
	}

	taskLogPath := taskName + constants.Sub + trainIdArr[1]

	logFilePath = config2.TrainPlatConfig.MountTrainContext + plat.LinuxAIConfigContext + getTasksMark() + plat.LinuxSpiltRex + taskId + plat.LinuxSpiltRex +
		train.JinnTrainResult + plat.LinuxSpiltRex + taskLogPath + plat.LinuxSpiltRex + train.TrainLog

	return
}

func ValidDealMaterialType(taskOrm model.TaskOrm) (err error) {
	var methodOrm model.MethodOrm

	dealMaterialTypes := make(map[string]bool)

	for ind, method := range taskOrm.Arguments {
		methodOrm, err = CreateMethodImpl().SelectMethodById(method.MethodId)
		if err != nil {
			return
		}

		if ind == 0 {
			dealMaterialTypes[methodOrm.DealMaterialType] = true
			continue
		}

		if _, ok := dealMaterialTypes[methodOrm.DealMaterialType]; !ok {
			err = errors.New("deal data type is not consistent")
			return
		}
	}
	return
}

func ValidDealClassification(taskOrm model.TaskOrm) (err error) {
	var methodOrm model.MethodOrm
	var categoryOrm model.CategoryOrm

	dealClassificationTypes := make(map[int32]bool)

	for ind, method := range taskOrm.Arguments {
		methodOrm, err = CreateMethodImpl().SelectMethodById(method.MethodId)
		if err != nil {
			return
		}

		categoryOrm, err = CreateCategoryImpl().SelectCategoryById(methodOrm.Id)
		if err != nil {
			return
		}

		if ind == 0 {
			dealClassificationTypes[categoryOrm.ParentId] = true
			continue
		}

		if _, ok := dealClassificationTypes[categoryOrm.ParentId]; !ok {
			err = errors.New("deal algorithm classification is not consistent")
			return
		}
	}
	return
}
