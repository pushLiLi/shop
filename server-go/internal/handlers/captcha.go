package handlers

import (
	"image/color"
	"net/http"

	"bycigar-server/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

var captchaInstance *base64Captcha.Captcha

func init() {
	driver := &base64Captcha.DriverString{
		Height:          80,
		Width:           240,
		NoiseCount:      20,
		ShowLineOptions: base64Captcha.OptionShowHollowLine | base64Captcha.OptionShowSlimeLine,
		Length:          4,
		Source:          "0123456789",
		BgColor:         &color.RGBA{R: 245, G: 245, B: 245, A: 255},
	}
	driver.ConvertFonts()
	store := base64Captcha.NewMemoryStore(1024, 300)
	captchaInstance = base64Captcha.NewCaptcha(driver, store)
}

func GetCaptcha(c *gin.Context) {
	id, b64s, _, err := captchaInstance.Generate()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "生成验证码失败")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"captchaId":    id,
		"captchaImage": b64s,
	})
}
