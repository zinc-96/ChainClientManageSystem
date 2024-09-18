package v1

import (
	"ChainClientManageSystem/config"
	"ChainClientManageSystem/internal/service"
	"ChainClientManageSystem/pkg/constant"
	"ChainClientManageSystem/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// Ping 健康检查
func Ping(c *gin.Context) {
	appConfig := config.GetGlobalConf().AppConfig
	confInfo, _ := json.MarshalIndent(appConfig, "", "  ")
	appInfo := fmt.Sprintf("app_name: %s\nversion: %s\n\n%s", appConfig.AppName, appConfig.Version,
		string(confInfo))
	c.String(http.StatusOK, appInfo)
}

// Register 注册
func Register(c *gin.Context) {
	req := &service.RegisterRequest{}
	rsp := &HttpResponse{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Errorf("request json err %v", err)
		rsp.ResponseWithError(c, CodeBodyBindErr, err.Error())
		return
	}
	if err := service.Register(req); err != nil {
		rsp.ResponseWithError(c, CodeRegisterErr, err.Error())
		return
	}
	rsp.ResponseSuccess(c)
}

// Login 登录
func Login(c *gin.Context) {
	req := &service.LoginRequest{}
	rsp := &HttpResponse{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Errorf("request json err %v", err)
		rsp.ResponseWithError(c, CodeBodyBindErr, err.Error())
		return
	}

	uuid := utils.Md5String(req.UserName + time.Now().GoString())
	ctx := context.WithValue(context.Background(), "uuid", uuid)
	log.Infof("loggin start,user:%s, password:%s", req.UserName, req.PassWord)
	session, err := service.Login(ctx, req)
	if err != nil {
		rsp.ResponseWithError(c, CodeLoginErr, err.Error())
		return
	}
	// 登陆成功，设置cookie
	c.SetCookie(constant.SessionKey, session, constant.CookieExpire, "/", "", false, true)
	rsp.ResponseSuccess(c)
}

// Logout 登出
func Logout(c *gin.Context) {
	session, _ := c.Cookie(constant.SessionKey)
	ctx := context.WithValue(context.Background(), constant.SessionKey, session)
	req := &service.LogoutRequest{}
	rsp := &HttpResponse{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		log.Errorf("bind get logout request json err %v", err)
		rsp.ResponseWithError(c, CodeBodyBindErr, err.Error())
		return
	}
	uuid := utils.Md5String(req.UserName + time.Now().GoString())
	ctx = context.WithValue(ctx, "uuid", uuid)
	if err := service.Logout(ctx, req); err != nil {
		rsp.ResponseWithError(c, CodeLogoutErr, err.Error())
		return
	}
	c.SetCookie(constant.SessionKey, session, -1, "/", "", false, true)
	rsp.ResponseSuccess(c)
}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	userName := c.Query("username")
	session, _ := c.Cookie(constant.SessionKey)
	ctx := context.WithValue(context.Background(), constant.SessionKey, session)
	req := &service.GetUserInfoRequest{
		UserName: userName,
	}
	rsp := &HttpResponse{}
	uuid := utils.Md5String(req.UserName + time.Now().GoString())
	ctx = context.WithValue(ctx, "uuid", uuid)
	userInfo, err := service.GetUserInfo(ctx, req)
	if err != nil {
		log.Debug("service.GetUserInfo: ", err)
		rsp.ResponseWithError(c, CodeGetUserInfoErr, err.Error())
		return
	}
	rsp.ResponseWithData(c, userInfo)
}

// UpdateNickName 更新用户昵称
func UpdateNickName(c *gin.Context) {
	req := &service.UpdateNickNameRequest{}
	rsp := &HttpResponse{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		log.Errorf("bind update user info request json err %v", err)
		rsp.ResponseWithError(c, CodeBodyBindErr, err.Error())
		return
	}
	session, _ := c.Cookie(constant.SessionKey)
	log.Infof("UpdateNickName|session=%s", session)
	ctx := context.WithValue(context.Background(), constant.SessionKey, session)
	uuid := utils.Md5String(req.UserName + time.Now().GoString())
	ctx = context.WithValue(ctx, "uuid", uuid)
	if err := service.UpdateUserNickName(ctx, req); err != nil {
		rsp.ResponseWithError(c, CodeUpdateUserInfoErr, err.Error())
		return
	}
	rsp.ResponseSuccess(c)
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	req := &service.DeleteUserRequest{}
	rsp := &HttpResponse{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		log.Errorf("bind delete user request json err %v", err)
		rsp.ResponseWithError(c, CodeBodyBindErr, err.Error())
		return
	}
	session, _ := c.Cookie(constant.SessionKey)
	log.Infof("DeleteUser|session=%s", session)
	ctx := context.WithValue(context.Background(), constant.SessionKey, session)
	uuid := utils.Md5String(req.UserName + time.Now().GoString())
	ctx = context.WithValue(ctx, "uuid", uuid)
	if err := service.DeleteUser(ctx, req); err != nil {
		rsp.ResponseWithError(c, CodeDeleteUserErr, err.Error())
		return
	}
	rsp.ResponseSuccess(c)
}

// CreateCert 创建证书
func CreateCert(c *gin.Context) {
	req := &service.CreateCertRequest{}
	rsp := &HttpResponse{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		log.Errorf("bind create cert request json err %v", err)
		rsp.ResponseWithError(c, CodeBodyBindErr, err.Error())
		return
	}
	session, _ := c.Cookie(constant.SessionKey)
	log.Infof("CreateCert|session=%s", session)
	ctx := context.WithValue(context.Background(), constant.SessionKey, session)
	uuid := utils.Md5String(req.UserName + time.Now().GoString())
	ctx = context.WithValue(ctx, "uuid", uuid)
	if err := service.CreateCert(ctx, req); err != nil {
		rsp.ResponseWithError(c, CodeCreateCertErr, err.Error())
		return
	}
	rsp.ResponseSuccess(c)
}

// QueryCert 查询证书
func QueryCert(c *gin.Context) {
	req := &service.QueryCertRequest{}
	rsp := &HttpResponse{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		log.Errorf("bind query cert request json err %v", err)
		rsp.ResponseWithError(c, CodeBodyBindErr, err.Error())
		return
	}
	session, _ := c.Cookie(constant.SessionKey)
	log.Infof("QueryCert|session=%s", session)
	ctx := context.WithValue(context.Background(), constant.SessionKey, session)
	uuid := utils.Md5String(req.UserName + time.Now().GoString())
	ctx = context.WithValue(ctx, "uuid", uuid)
	cert, err := service.QueryCert(ctx, req)
	if err != nil {
		rsp.ResponseWithError(c, CodeQueryCertErr, err.Error())
		return
	}
	rsp.ResponseWithData(c, cert)
}

// TestContract 测试合约
func TestContract(c *gin.Context) {
	req := &service.UserContractClaimInvokeRequest{}
	rsp := &HttpResponse{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		log.Errorf("bind test contract request json err %v", err)
		rsp.ResponseWithError(c, CodeBodyBindErr, err.Error())
		return
	}
	session, _ := c.Cookie(constant.SessionKey)
	log.Infof("TestContract|session=%s", session)
	ctx := context.WithValue(context.Background(), constant.SessionKey, session)
	uuid := utils.Md5String(req.UserName + time.Now().GoString())
	ctx = context.WithValue(ctx, "uuid", uuid)
	filehash, err := service.TestUserContractClaimInvoke(ctx, req)
	if err != nil {
		log.Errorf("TestContract|TestUserContractClaimInvoke err %v", err)
		rsp.ResponseWithError(c, CodeTestContractErr, err.Error())
		return
	}
	rsp.ResponseWithData(c, filehash)
}
