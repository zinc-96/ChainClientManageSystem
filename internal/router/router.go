package router

import (
	api "ChainClientManageSystem/api/http/v1"
	"ChainClientManageSystem/config"
	"ChainClientManageSystem/pkg/constant"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// InitRouterAndServe 路由配置、启动服务
func InitRouterAndServe() {

	setAppRunMode()
	r := gin.Default()

	// 健康检查
	r.GET("ping", api.Ping)
	// 用户注册
	r.POST("/user/register", api.Register)
	// 用户登录
	r.POST("/user/login", api.Login)
	// 用户登出
	r.POST("/user/logout", AuthMiddleWare(), api.Logout)
	// 获取用户信息
	r.GET("/user/get_user_info", AuthMiddleWare(), api.GetUserInfo)
	// 更新用户信息
	r.POST("/user/update_nick_name", AuthMiddleWare(), api.UpdateNickName)
	// 用户注销
	r.POST("/user/delete", AuthMiddleWare(), api.DeleteUser)
	// 创建证书
	r.POST("/cert/create", AuthMiddleWare(), api.CreateCert)
	// 查询证书
	r.POST("/cert/query", AuthMiddleWare(), api.QueryCert)
	// 测试合约
	r.POST("/contract/test", AuthMiddleWare(), api.TestContract)

	// 指定静态文件目录
	r.Static("/static/", "./web/static/")

	// 启动server
	port := config.GetGlobalConf().AppConfig.Port
	if err := r.Run(":" + strconv.Itoa(port)); err != nil {
		log.Error("start server err:" + err.Error())
	}
}

// setAppRunMode 设置运行模式
func setAppRunMode() {
	if config.GetGlobalConf().AppConfig.RunMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
}

// AuthMiddleWare 认证中间件
func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		if session, err := c.Cookie(constant.SessionKey); err == nil {
			if session != "" {
				c.Next()
				return
			}
		}
		// 返回错误
		c.JSON(http.StatusUnauthorized, gin.H{"error": "err"})
		c.Abort()
		return
	}
}
