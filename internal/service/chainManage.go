package service

import (
	"ChainClientManageSystem/config"
	"ChainClientManageSystem/internal/dao"
	"bytes"
	"chainmaker.org/chainmaker/common/v2/random/uuid"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	. "chainmaker.org/chainmaker/sdk-go/v2"
	. "chainmaker.org/chainmaker/utils/v2"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"strings"
)

// CaTest /*
//   - @Description: 测试CA服务是否正常工作
//   - @param caServerUrl
//   - @return bool
//     */
func CaTest(caServerUrl string) bool {
	resp, err := http.Get(caServerUrl + "/test")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	println(string(body))
	return true
}

// CaQueryCerts /*
//   - @Description: 查询CA证书
//   - @param caServerUrl, orgId, userType, certUsage
//   - @return string, error
//   - @example: 获取组织TestCMorg1中用户test01的sign证书: CaQueryCerts("http://192.168.1.104:8096", "", "test01", "client", "sign")
//     */
func CaQueryCerts(caServerUrl string, orgId string, userId string, userType string, certUsage string) (string, error) {
	// 构造参数
	jsonBody := `{"orgId":"` + orgId + `","userId":"` + userId + `","userType":"` + userType + `","certUsage":"` + certUsage + `"}`
	// 发送请求
	resp, err := http.Post(caServerUrl+"/api/ca/querycerts", "application/json", bytes.NewReader([]byte(jsonBody)))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if gjson.Get(string(body), "code").Int() == 200 {
		return gjson.Get(string(body), "data").Array()[0].Get("certContent").String(), nil
	}
	return "", err
}

// CaGenCert /*
//   - @Description: 生成证书
//   - @param caServerUrl, orgId, userId, userType, certUsage, country, locality, province
//   - @return int, int, string, string, error
//   - @notice: 生成共识节点证书时，userId需要保证链上唯一；同一节点的Sign和Tls证书，userId需要保持一致。
//   - @example: 生成组织TestCMorg1中用户test01的sign证书: CaGenCert("http://192.168.1.104:8096", "TestCMorg1", "test01", "client", "sign", "CN", "Beijing", "Beijing")
//     */
func CaGenCert(caServerUrl string, orgId string, userId string, userType string, certUsage string, country string, locality string, province string) (int, int, string, string, error) {
	// 构造参数
	jsonBody := `{"orgId":"` + orgId + `","userId":"` + userId + `","userType":"` + userType + `","certUsage":"` + certUsage + `","country":"` + country + `","locality":"` + locality + `","province":"` + province + `"}`
	// 发送请求
	resp, err := http.Post(caServerUrl+"/api/ca/gencert", "application/json", bytes.NewReader([]byte(jsonBody)))
	if err != nil {
		return 0, 0, "", "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("CaGenCert|%v", err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, "", "", err
	}
	if gjson.Get(string(body), "code").Int() == 200 {
		data := gjson.Get(string(body), "data").Array()[0]
		return int(gjson.Get(data.String(), "certSn").Int()), int(gjson.Get(data.String(), "issueCertSn").Int()), gjson.Get(data.String(), "cert").String(), gjson.Get(data.String(), "privateKey").String(), nil
	}
	return 0, 0, "", "", nil
}

// CreateNode /*
//   - @Description: 创建节点
//   - @param nodeAddr, connCnt, NodeCAPath, tlsHostName
//   - @return *NodeConfig
//     */
func CreateNode(nodeAddr string, connCnt int, NodeCAPath []string, tlsHostName string) *NodeConfig {
	node := NewNodeConfig(
		// 节点地址，格式：127.0.0.1:12301
		WithNodeAddr(nodeAddr),
		// 节点连接数
		WithNodeConnCnt(connCnt),
		// 节点是否启用TLS认证
		WithNodeUseTLS(true),
		// 根证书路径
		WithNodeCAPaths(NodeCAPath),
		// TLS Hostname
		WithNodeTLSHostName(tlsHostName),
	)
	return node
}

// CreateClient /*
//   - @Description: 创建链客户端
//   - @param chainOrgId, chainId, userKeyPath, userCertPath, userSignKeyPath, userSignCrtPath, node1
//   - @return *ChainClient, error
//     */
func CreateClient(chainOrgId string, chainId string, userKeyPath string, userCertPath string, userSignKeyPath string, userSignCrtPath string, node1 *NodeConfig) (*ChainClient, error) {
	chainClient, err := NewChainClient(
		// 设置归属组织
		WithChainClientOrgId(chainOrgId),
		// 设置链ID
		WithChainClientChainId(chainId),
		// 设置logger句柄，若不设置，将采用默认日志文件输出日志
		WithChainClientLogger(log.StandardLogger()),
		// 设置客户端用户私钥路径
		WithUserKeyFilePath(userKeyPath),
		// 设置客户端用户证书路径
		WithUserCrtFilePath(userCertPath),
		// 设置客户端用户签名私钥路径
		WithUserSignKeyFilePath(userSignKeyPath),
		// 设置客户端用户签名证书路径
		WithUserSignCrtFilePath(userSignCrtPath),
		// 设置用户签名类型
		WithAuthType("permissionedwithcert"),
		// 添加节点
		AddChainClientNodeConfig(node1),
	)
	if err != nil {
		return nil, err
	}

	return chainClient, nil
}

// TestUserContractClaimInvoke /*
//   - @Description: 测试调用合约
//   - @param client, contractName, method, withSyncResult
//   - @return string, error
//     */
func TestUserContractClaimInvoke(ctx context.Context, req *UserContractClaimInvokeRequest) (string, error) {
	uuid := ctx.Value("uuid")
	log.Debugf(" %s| TestUserContractClaimInvoke access from:%s", uuid, req.UserName)
	//获取证书
	certJson, err := dao.QueryCert(req.UserName)
	data := gjson.Get(certJson, "data").Array()
	if len(data) != 2 {
		log.Errorf("QueryCert|cert not found")
		return "", fmt.Errorf("query cert|cert not found")
	}
	orgId := data[0].Get("OrgId").String()
	orgId = strings.Replace(orgId, "+", "", -1)
	tlsCert := data[0].Get("Cert").String()
	signCert := data[1].Get("Cert").String()
	tlsKey := data[0].Get("PrivateKey").String()
	signKey := data[1].Get("PrivateKey").String()
	tlsKeyPath := "ChainClientManageSystem\\secret\\" + req.UserName + "\\TlsKey.key"
	tlsCertPath := "ChainClientManageSystem\\secret\\" + req.UserName + "\\TlsCert.crt"
	signKeyPath := "ChainClientManageSystem\\secret\\" + req.UserName + "\\SignKey.key"
	signCertPath := "ChainClientManageSystem\\secret\\" + req.UserName + "\\SignCert.crt"
	_ = SaveCertToFile(tlsKey, tlsKeyPath)
	_ = SaveCertToFile(tlsCert, tlsCertPath)
	_ = SaveCertToFile(signKey, signKeyPath)
	_ = SaveCertToFile(signCert, signCertPath)
	nodeCAPath := []string{"ChainClientManageSystem\\secret\\" + orgId}
	// 创建节点
	node := CreateNode(config.GetGlobalConf().ChainNodeConfig.ChainNodeURL, 10, nodeCAPath, "chainmaker.org")
	// 创建客户端
	chainClient, err := CreateClient(orgId, config.GetGlobalConf().ChainConfig.ChainID, tlsKeyPath, tlsCertPath, signKeyPath, signCertPath, node)
	if err != nil {
		log.Errorf("CreateClient|%v", err)
		return "", fmt.Errorf("create client|%v", err)
	}
	// 调用合约
	fileHash, err := UserContractClaimInvoke(chainClient, req.ContractName, req.Method, req.WithSyncResult)
	if err != nil {
		log.Errorf("UserContractClaimInvoke|%v", err)
		return "", fmt.Errorf("user contract claim invoke|%v", err)
	}
	return fileHash, nil
}

func UserContractClaimInvoke(client *ChainClient, contractName string, method string, withSyncResult bool) (string, error) {

	curTime := fmt.Sprintf("%d", CurrentTimeMillisSeconds())
	fileHash := uuid.GetUUID()
	var params []*common.KeyValuePair
	params = append(params, &common.KeyValuePair{Key: "time", Value: []byte(curTime)})
	params = append(params, &common.KeyValuePair{Key: "file_hash", Value: []byte(fileHash)})
	params = append(params, &common.KeyValuePair{Key: "file_name", Value: []byte(fmt.Sprintf("file_%s", curTime))})

	err := invokeUserContract(client, contractName, method, "", params, withSyncResult)
	if err != nil {
		return "", err
	}

	return fileHash, nil
}

func invokeUserContract(client *ChainClient, contractName, method, txId string, params []*common.KeyValuePair, withSyncResult bool) error {

	resp, err := client.InvokeContract(contractName, method, txId, params, -1, withSyncResult)
	if err != nil {
		return err
	}

	if resp.Code != common.TxStatusCode_SUCCESS {
		return fmt.Errorf("invoke contract failed, [code:%d]/[msg:%s]\n", resp.Code, resp.Message)
	}

	if !withSyncResult {
		fmt.Printf("invoke contract success, resp: [code:%d]/[msg:%s]/[txId:%s]\n", resp.Code, resp.Message, resp.ContractResult.Result)
	} else {
		fmt.Printf("invoke contract success, resp: [code:%d]/[msg:%s]/[contractResult:%s]\n", resp.Code, resp.Message, resp.ContractResult)
	}

	return nil
}
