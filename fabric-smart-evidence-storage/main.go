package main

import (
	"fabric-smart-evidence-storage/conf"
	"fabric-smart-evidence-storage/server"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 从配置文件读取配置
	conf.Init()

	// 装载路由
	gin.SetMode(os.Getenv("GIN_MODE"))
	r := server.NewRouter()
	r.Run(":10010")
}
