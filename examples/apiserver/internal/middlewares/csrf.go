package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mel0dys0ng/song/core/https"
	"github.com/mel0dys0ng/song/examples/apiserver/internal/status"
	"github.com/mel0dys0ng/song/examples/apiserver/internal/utils"
)

// SetupVerifyCSRFToken 校验CSRF-Token
func SetupVerifyCSRFToken(eng *gin.Engine) gin.HandlerFunc {
	return https.SetupVerifyCSRFToken(eng,
		https.CSRFOptionCrypter(utils.NewAESGCMCrypter()),
		https.CSRFOptionErLog(status.Forbidden),
		https.CSRFOptionMethods(http.MethodPost, http.MethodDelete, http.MethodPut),
		https.CSRFOptionOrigins("http://localhost:8080"),
		https.CSRFOptionRefererPrefixs("http://localhost:8080"),
		https.CSRFOptionDomain("localhost"),
		https.CSRFOptionMaxAge(3600),
	)
}
