package api

import (
	"gin-bee/apps/tool/request"
	"github.com/gin-gonic/gin"
	qrcode "github.com/skip2/go-qrcode"
	"net/http"
)

var CQRCodeController = QRCodeController{}

type QRCodeController struct {
}

// GenerateQRCode
// @Summary
// @Schemes
// @Description 生成二维码
// @Tags
// @Security ApiKeyAuth
// @Accept json
// @Produce image/png
// @Param object body request.QrCodeReq true "请求参数"
// @Success 200 {string} string 图片
// @Failure 400 {object} response.Response
// @Router /tool/qr-code [post]
func (qc QRCodeController) GenerateQRCode(c *gin.Context) {
	var req request.QrCodeReq
	err := c.BindJSON(&req)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "请求参数不正确"})
		return
	}
	img, err := qrcode.Encode(req.Url, qrcode.Medium, req.Size)

	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "生成二维码失败"})
		return
	}

	c.Data(http.StatusOK, "image/png", img)
}
