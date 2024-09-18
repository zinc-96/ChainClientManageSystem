package service

// RegisterRequest 注册请求
type RegisterRequest struct {
	OrgId    string `json:"org_id"`
	UserName string `json:"user_name"`
	Password string `json:"pass_word"`
	NickName string `json:"nick_name"`
}

// LoginRequest 登陆请求
type LoginRequest struct {
	UserName string `json:"user_name"`
	PassWord string `json:"pass_word"`
}

// LogoutRequest 登出请求
type LogoutRequest struct {
	UserName string `json:"user_name"`
}

// GetUserInfoRequest 获取用户信息请求
type GetUserInfoRequest struct {
	UserName string `json:"user_name"`
}

// GetUserInfoResponse 获取用户信息返回结构
type GetUserInfoResponse struct {
	UserName string `json:"user_name"`
	PassWord string `json:"pass_word"`
	NickName string `json:"nick_name"`
}

// UpdateNickNameRequest 修改用户信息返回结构
type UpdateNickNameRequest struct {
	UserName    string `json:"user_name"`
	NewNickName string `json:"new_nick_name"`
}

// DeleteUserRequest 删除用户请求
type DeleteUserRequest struct {
	UserName string `json:"user_name"`
}

// CreateCertRequest 创建证书请求
type CreateCertRequest struct {
	UserName  string `json:"user_name"`
	OrgId     string `json:"org_id"`
	UserId    string `json:"user_id"`
	UserType  string `json:"user_type"`
	CertUsage string `json:"cert_usage"`
	Country   string `json:"country"`
	Locality  string `json:"locality"`
	Province  string `json:"province"`
}

// QueryCertRequest 查询证书请求
type QueryCertRequest struct {
	UserName string `json:"user_name"`
}

// UserContractClaimInvokeRequest 用户合约声明调用请求
type UserContractClaimInvokeRequest struct {
	UserName       string `json:"user_name"`
	ContractName   string `json:"contract_name"`
	Method         string `json:"method"`
	WithSyncResult bool   `json:"with_sync_result"`
}
