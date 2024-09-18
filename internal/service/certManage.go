package service

import (
	"ChainClientManageSystem/config"
	"ChainClientManageSystem/internal/dao"
	"ChainClientManageSystem/internal/model"
	"ChainClientManageSystem/pkg/constant"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"os"
	"strconv"
	"strings"
)

// CreateCert 创建证书
func CreateCert(ctx context.Context, req *CreateCertRequest) error {
	uuid := ctx.Value(constant.ReqUuid)
	log.Debugf(" %s| CreateCert access from:%s", uuid, req.UserId)
	if req.OrgId == "" || req.UserId == "" || req.UserType == "" || req.CertUsage == "" || req.Country == "" || req.Locality == "" || req.Province == "" {
		log.Errorf("create cert param invalid")
		return fmt.Errorf("create cert param invalid")
	}
	// 生成证书
	certSn, issueCertSn, cert, privateKey, err := CaGenCert(config.GetGlobalConf().ChainCAConfig.URL, req.OrgId, req.UserId, req.UserType, req.CertUsage, req.Country, req.Locality, req.Province)
	if certSn == 0 || err != nil {
		log.Errorf("CreateCert|%v", err)
		return fmt.Errorf("create cert|%v", err)
	}
	// 保存证书到数据库
	user := &model.User{
		Name: req.UserName,
		CreateModel: model.CreateModel{
			Creator: req.UserName,
		},
		ModifyModel: model.ModifyModel{
			Modifier: req.UserName,
		},
	}
	// 创建json格式的newCert
	newCert := `{"OrgId":"` + req.OrgId + `","UserId":"` + req.UserId + `","UserType":"` + req.UserType + `","CertUsage":"` + req.CertUsage + `","Country":"` + req.Country + `","Locality":"` + req.Locality + `","Province":"` + req.Province + `","CertSn":` + strconv.Itoa(certSn) + `,"IssueCertSn":` + strconv.Itoa(issueCertSn) + `,"Cert":"` + cert + `","PrivateKey":"` + privateKey + `"}`
	newCert = strings.Replace(newCert, "\n", "\\n", -1)
	if err := dao.CreateCert(newCert, user); err != nil {
		log.Errorf("CreateCert|%v", err)
		return fmt.Errorf("create cert|%v", err)
	}
	return nil
}

// QueryCert 查询证书
func QueryCert(ctx context.Context, req *QueryCertRequest) (string, error) {
	uuid := ctx.Value(constant.ReqUuid)
	log.Debugf(" %s| QueryCert access from:%s", uuid, req.UserName)
	if req.UserName == "" {
		log.Errorf("query cert param invalid")
		return "", fmt.Errorf("query cert param invalid")
	}
	// 查询证书
	println(req.UserName)
	cert, err := dao.QueryCert(req.UserName)
	println(cert)
	if err != nil {
		log.Errorf("QueryCert|%v", err)
		return "", fmt.Errorf("query cert|%v", err)
	}
	return cert, nil
}

// SaveCertToFile 保存证书到本地
func SaveCertToFile(cert string, filePath string) error {
	cert = strings.Replace(cert, "\\n", "\n", -1)
	err := os.MkdirAll(filePath[:strings.LastIndex(filePath, "\\")], 0755)
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath, []byte(cert), 0755)
	if err != nil {
		return err // 返回错误信息
	}
	return nil // 成功则返回 nil
}
