package jsapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"net/http"
	"strings"
	"time"
)

// GetPrepayInfo 生成预支付信息
// description: 订单描述
// outTradeNo: 商户订单号
// expireAt: 订单过期时间
// amount: 订单金额
// openId: 支付者openId
// notifyUrl: 支付结果通知地址
func (c *Client) GetPrepayInfo(description, outTradeNo, openId, notifyUrl string,
	expireAt time.Time, amount float32) (*PrePayInfo, error) {
	privateKeyStr := c.Cert
	privateKeyStr = strings.ReplaceAll(privateKeyStr, "\r\n", "\n")
	privateKeyStr = strings.TrimSpace(privateKeyStr)

	// 加载临时文件中的私钥
	mchPrivateKey, err := utils.LoadPrivateKey(privateKeyStr)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(c.MchId, c.CertNum, mchPrivateKey, c.Secret),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		return nil, err
	}
	svc := jsapi.JsapiApiService{Client: client}
	resp, _, err := svc.Prepay(ctx,
		jsapi.PrepayRequest{
			Appid:       core.String(c.AppId),
			Mchid:       core.String(c.MchId),
			Description: core.String(description),
			OutTradeNo:  core.String(outTradeNo),
			TimeExpire:  core.Time(expireAt),
			Amount: &jsapi.Amount{
				Currency: core.String("CNY"),
				Total:    core.Int64(int64(amount * 100)),
			},
			Payer: &jsapi.Payer{
				Openid: core.String(openId),
			},
			NotifyUrl: core.String(notifyUrl),
		},
	)
	if err != nil {
		return nil, err
	}

	// 生成支付参数
	timeStamp := fmt.Sprintf("%d", time.Now().Unix())
	nonceStr := generateNonceStr()
	pkg := "prepay_id=" + *resp.PrepayId

	// 计算签名
	paySign, err := calculatePaySign(c.AppId, timeStamp, nonceStr, pkg, mchPrivateKey)
	if err != nil {
		return nil, err
	}

	return &PrePayInfo{
		AppId:     c.AppId,
		Timestamp: timeStamp,
		NonceStr:  nonceStr,
		Package:   pkg,
		SignType:  "RSA",
		PaySign:   paySign,
	}, nil
}

// DealPayNotify 处理支付结果通知
func (c *Client) DealPayNotify(req *http.Request) (*PayNotifyResponse, error) {
	privateKeyStr := c.Cert
	privateKeyStr = strings.ReplaceAll(privateKeyStr, "\r\n", "\n")
	privateKeyStr = strings.TrimSpace(privateKeyStr)
	mchPrivateKey, err := utils.LoadPrivateKey(privateKeyStr)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	err = downloader.MgrInstance().RegisterDownloaderWithPrivateKey(ctx, mchPrivateKey, c.CertNum,
		c.MchId, c.Secret)
	certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(c.MchId)
	handler, err := notify.NewRSANotifyHandler(c.Secret, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))
	if err != nil {
		return nil, err
	}

	content := make(map[string]interface{})
	_, err = handler.ParseNotifyRequest(context.Background(), req, &content)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}
	var resp PayNotifyResponse
	err = json.Unmarshal(jsonData, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Refund 退款申请
// reason: 退款原因
// outTradeNo: 商户订单号
// outRefundNo: 商户退款单号
// refund: 退款金额
// total: 订单金额
func (c *Client) Refund(outTradeNo, outRefundNo, reason, notifyUrl string,
	refund, total float32) (*refunddomestic.Refund, error) {
	privateKeyStr := c.Cert
	privateKeyStr = strings.ReplaceAll(privateKeyStr, "\r\n", "\n")
	privateKeyStr = strings.TrimSpace(privateKeyStr)
	mchPrivateKey, err := utils.LoadPrivateKey(privateKeyStr)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(c.MchId, c.CertNum, mchPrivateKey, c.Secret),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		return nil, err
	}

	svc := refunddomestic.RefundsApiService{Client: client}
	resp, _, err := svc.Create(ctx,
		refunddomestic.CreateRequest{
			OutTradeNo:   core.String(outTradeNo),
			OutRefundNo:  core.String(outRefundNo),
			Reason:       core.String(reason),
			FundsAccount: refunddomestic.REQFUNDSACCOUNT_AVAILABLE.Ptr(),
			Amount: &refunddomestic.AmountReq{
				Currency: core.String("CNY"),
				Refund:   core.Int64(int64(refund * 100)),
				Total:    core.Int64(int64(total * 100)),
			},
			NotifyUrl: core.String(notifyUrl),
		},
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
