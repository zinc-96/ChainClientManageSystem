package dao

import (
	"ChainClientManageSystem/config"
	"ChainClientManageSystem/internal/model"
	"ChainClientManageSystem/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
)

// CreateCert 创建证书
func CreateCert(cert string, user *model.User) error {
	existedCerts, _ := QueryCert(user.Name)
	newCerts := ""
	if existedCerts != "" {
		newCerts = existedCerts[:len(existedCerts)-2] + "," + cert + "]}"
	} else {
		newCerts = `{"data":[` + cert + "]}"
	}
	encNewCerts, _ := utils.EncryptAES(newCerts, config.GetGlobalConf().DbConfig.AesKey, config.GetGlobalConf().DbConfig.AesIv)
	if err := utils.GetDB().Model(&model.User{}).Where("`name` = ?", user.Name).Update("cert", encNewCerts).Error; err != nil {
		log.Errorf("CreateCert fail: %v", err)
		return fmt.Errorf("CreateCert fail: %v", err)
	}
	log.Infof("insert success")
	return nil
}

// QueryCert 查询证书
func QueryCert(userName string) (string, error) {
	var user model.User
	if err := utils.GetDB().Where("`name` = ?", userName).First(&user).Error; err != nil {
		log.Errorf("QueryCert fail: %v", err)
		return "", fmt.Errorf("QueryCert fail: %v", err)
	}
	log.Infof("query success")
	cert, _ := utils.DecryptAES(user.Cert, config.GetGlobalConf().DbConfig.AesKey, config.GetGlobalConf().DbConfig.AesIv)
	return cert, nil
}
