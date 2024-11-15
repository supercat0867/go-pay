package handler

import (
	"github.com/gin-gonic/gin"
	"go-pay/internal/model"
	"go-pay/pkg/wechatpay/jsapi"
	"log"
	"net/http"
	"strconv"
)

// DealWechatPayNotify 处理微信支付
func (h *Handler) DealWechatPayNotify(g *gin.Context) {
	mchId := g.Param("id")
	mchIdInt, err := strconv.Atoi(mchId)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"code": "FAIL"})
		return
	}
	// 查询商户信息
	merchant, err := h.Service.MerchantRepo.FindByID(uint(mchIdInt))
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"code": "FAIL"})
		return
	}
	client := jsapi.NewClient(merchant.AppID, merchant.MchID, merchant.Secret, merchant.Cert, merchant.CertNum)
	resp, err := client.DealPayNotify(g.Request)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"code": "FAIL"})
		return
	}

	// 查询支付记录
	pay, err := h.Service.PayRepo.FindByMchIDAndOrderID(resp.MchID, resp.OutTradeNo)
	if err != nil {
		g.JSON(http.StatusOK, gin.H{"code": "SUCCESS"})
		return
	}

	if pay.PayState != model.PayStatePending {
		g.JSON(http.StatusOK, gin.H{"code": "SUCCESS"})
		return
	}

	if resp.TradeState == "SUCCESS" {
		// 更新支付记录状态
		pay.PayState = model.PayStateSuccess
		pay.TransactionID = resp.TransactionID
		if err = h.Service.PayRepo.Update(pay); err != nil {
			log.Println(err)
			g.JSON(http.StatusInternalServerError, gin.H{"code": "FAIL"})
			return
		}
	}
	g.JSON(http.StatusOK, gin.H{"code": "SUCCESS"})
}
