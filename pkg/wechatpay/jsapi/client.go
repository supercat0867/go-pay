package jsapi

type Client struct {
	AppId   string // 应用Id
	MchId   string // 商户Id
	Secret  string // 密钥
	Cert    string // 证书
	CertNum string // 证书序列号
}

func NewClient(appid, mchid, secret, cert, certNum string) *Client {
	return &Client{
		AppId:   appid,
		MchId:   mchid,
		Secret:  secret,
		Cert:    cert,
		CertNum: certNum,
	}
}
