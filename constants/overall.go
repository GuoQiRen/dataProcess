package constants

const (
	Root                         = "tmp"               // 预处理磁盘
	Preprocessing                = "preprocessing"     // 预处理根目录
	InputData                    = "input_data"        // 输入数据目录
	Resource                     = "resource"          // 输入数据目录原层级
	Result                       = "result"            // 输入数据目录原层级
	OutPutData                   = "output_data"       // 输出数据目录
	JsonSuffix                   = ".json"             // 保存json后缀
	PdfSuffix                    = ".pdf"              // 保存pdf后缀
	JsonFile                     = "jsonFile"          // 保存jsonFile
	Colon                        = ":"                 // 通用冒号
	Point                        = "."                 // 通用点号
	WhiteSpace                   = " "                 // 空格符号
	ParamsPrefix                 = "--"                // 参数前缀符号
	Sub                          = "-"                 // 通用减号
	Quotes                       = "'"                 // 单引号
	CD                           = "cd"                // cd符号
	SRC                          = "src"               // src文件夹代表源文件
	PreFixCMD                    = " && ./init_client" // 启动任务的命令
	Tasks                        = "tasks"             // 存放任务的文件夹
	Methods                      = "methods"           // 存放上传算子的文件夹
	MethodExec                   = "MethodExec"        // 方法执行入口
	RunEntry                     = "run.sh"            // 方法执行入口命令
	JinnPlatformCode             = "JinnData_material" // jinn用户组代码
	SuperManager                 = "SuperManager"      // jinn用户组的超级管理员
	NormalManager                = "normalManager"     // jinn用户组的超级管理员
	JinnToken                    = "token"             // jinn_token
	SubLine                      = "_"                 // 下划线
	Material                     = "material"          // 材料
	MarkData                     = "markData"          // 标注数据
	MarkDataDataProcessingPreFix = "数据处理"              // 标注数据上传前缀
	Formalization                = "formalization"     // 正式环境的标识
	Testing                      = "testing"           // 测试环境的标识
	Integration                  = "integration"       // 集成环境的标识
	Security                     = "security"          // 集成环境的标识
	TrainProposer                = 9933914             // 提交训练任务的用户id
	SchemaVer                    = 2                   // dahua2标注文件唯一id标识
)
