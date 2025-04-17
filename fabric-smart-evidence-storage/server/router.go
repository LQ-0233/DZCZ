package server

import (
	"fabric-smart-evidence-storage/api"
	"fabric-smart-evidence-storage/middleware"

	"github.com/gin-gonic/gin"
)

// NewRouter 路由配置
func NewRouter() *gin.Engine {
	r := gin.Default()

	// 中间件, 顺序不能改
	r.Use(middleware.Cors())

	// 路由
	v1 := r.Group("/api/v1")
	{
		v1.POST("common/upload", api.UploadFile)
		v1.GET("common/download", api.DownloadFile)
		// 用户登录
		v1.POST("user/register", api.UserRegister)

		// 用户登录
		v1.POST("user/login", api.UserLogin)

		// 用户列表
		auth := v1.Group("")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.POST("user/changePwd", api.UserChangePwd)

			input := auth.Group("")
			input.Use(middleware.InputAuthRequired())
			{
				input.GET("evidence/authorizeUserList", api.GetAuthorizeUserList)
				// 证据相关路由
				input.POST("evidence/create", api.EvidenceCreate)
				input.GET("evidence/userEvidence", api.GetUserEvidence)
				input.POST("evidence/authorize", api.EvidenceAuthorize)
				input.POST("evidence/cancelAuthorize", api.EvidenceCancelAuthorize)
				input.GET("evidence/authorizedList", api.EvidenceAuthorizedList)
				input.GET("evidence/viewRecordList", api.EvidenceViewRecordList)
			}

			auth.GET("evidence/receivedAuthorizedList", api.EvidenceReceivedAuthorizedList)
			auth.POST("evidence/view", api.EvidenceView)

			// 管理员
			admin := auth.Group("")
			admin.Use(middleware.AdminAuthRequired())
			{
				admin.POST("user/update", api.UserUpdate)
				admin.GET("user/list", api.UserList)
			}
		}

	}
	return r
}
