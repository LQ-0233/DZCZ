package conf

import (
	"fabric-smart-evidence-storage/fabric"
	"fabric-smart-evidence-storage/util"
	"os"

	"github.com/joho/godotenv"
)

// Init 初始化配置项
func Init() {
	// 从本地读取环境变量
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// 设置日志级别
	util.BuildLogger(os.Getenv("LOG_LEVEL"))

	fabric.Init()
	util.InitIPFSClient(os.Getenv("IPFS_NODE_URL"))
	pkBytes, err := os.ReadFile("./pk")
	if err != nil {
		panic(err)
	}
	util.InitGroth16(pkBytes)
}
