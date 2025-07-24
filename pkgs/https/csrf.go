package https

import (
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mel0dys0ng/song/pkgs/erlogs"
	"github.com/mel0dys0ng/song/pkgs/utils/cljent"
	"github.com/mel0dys0ng/song/pkgs/utils/crypto"
	"github.com/mel0dys0ng/song/pkgs/utils/strjngs"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

const (
	DidLength                 = 32
	CSRFTokenCookieName       = "X-Song-Csrf-Token"
	CSRFTokenHeaderName       = "X-Song-Csrf-Token"
	CSRFTokenExpireHeaderName = "X-Song-Csrf-Token-Expire"
	ClientRequestIdHeaderName = "X-Song-Request-Id"
	ClientDidHeaderName       = "X-Song-Did"
	ClientCTHeaderName        = "X-Song-Ct"
	ClientCVHeaderName        = "X-Song-Cv"
)

type (
	CSRFConfig struct {
		Crypter            *crypto.AESGCMCrypter
		Verify             func(ctx *gin.Context) bool
		Erlog              erlogs.ErLogInterface
		Methods            []string
		Origins            []string
		RefererPrefixs     []string
		CookieKey          string
		HeaderKey          string
		HeaderExpireKey    string
		HeaderRequestIdKey string
		HeaderDidKey       string
		HeaderCTKey        string
		HeaderCVKey        string
		Path               string
		Domain             string
		Secure             bool
		HttpOnly           bool
		MaxAge             int
		SameSite           int64
	}

	CSRFOption struct {
		Apply func(c *CSRFConfig)
	}
)

func CSRFOptionVerify(f func(ctx *gin.Context) bool) CSRFOption {
	return CSRFOption{
		Apply: func(c *CSRFConfig) {
			c.Verify = f
		},
	}
}

func CSRFOptionCrypter(crypter *crypto.AESGCMCrypter) CSRFOption {
	return CSRFOption{
		Apply: func(c *CSRFConfig) {
			c.Crypter = crypter
		},
	}
}

func CSRFOptionErLog(e erlogs.ErLogInterface) CSRFOption {
	return CSRFOption{
		Apply: func(c *CSRFConfig) {
			c.Erlog = e
		},
	}
}

func CSRFOptionCookieKey(s string) CSRFOption {
	return CSRFOption{
		Apply: func(c *CSRFConfig) {
			c.CookieKey = s
		},
	}
}

func CSRFOptionHeaderKey(s string) CSRFOption {
	return CSRFOption{
		Apply: func(c *CSRFConfig) {
			c.HeaderKey = s
		},
	}
}

func CSRFOptionHeaderExpireKey(s string) CSRFOption {
	return CSRFOption{
		Apply: func(c *CSRFConfig) {
			c.HeaderExpireKey = s
		},
	}
}

func CSRFOptionHeaderRequestIdKey(s string) CSRFOption {
	return CSRFOption{
		Apply: func(c *CSRFConfig) {
			c.HeaderRequestIdKey = s
		},
	}
}

func CSRFOptionHeaderDidKey(s string) CSRFOption {
	return CSRFOption{
		Apply: func(c *CSRFConfig) {
			c.HeaderDidKey = s
		},
	}
}

func CSRFOptionHeaderCTKey(s string) CSRFOption {
	return CSRFOption{
		Apply: func(c *CSRFConfig) {
			c.HeaderCTKey = s
		},
	}
}

func CSRFOptionHeaderCVKey(s string) CSRFOption {
	return CSRFOption{
		Apply: func(c *CSRFConfig) {
			c.HeaderCVKey = s
		},
	}
}

func CSRFOptionMaxAge(i int) CSRFOption {
	return CSRFOption{
		Apply: func(c *CSRFConfig) {
			c.MaxAge = i
		},
	}
}

func CSRFOptionMethods(ss []string) CSRFOption {
	return CSRFOption{
		Apply: func(c *CSRFConfig) {
			c.Methods = ss
		},
	}
}

func CSRFOptionOrigins(ss []string) CSRFOption {
	return CSRFOption{
		Apply: func(c *CSRFConfig) {
			c.Origins = ss
		},
	}
}

func CSRFOptionRefererPrefixs(ss []string) CSRFOption {
	return CSRFOption{
		Apply: func(c *CSRFConfig) {
			c.RefererPrefixs = ss
		},
	}
}

func CSRFOptionDomain(s string) CSRFOption {
	return CSRFOption{
		Apply: func(c *CSRFConfig) {
			c.Domain = s
		},
	}
}

func CSRFOptionCookiePath(s string) CSRFOption {
	return CSRFOption{
		Apply: func(c *CSRFConfig) {
			c.Path = s
		},
	}
}

var sg = singleflight.Group{}

func defaultCSRFConfig() *CSRFConfig {
	return &CSRFConfig{
		Verify:             func(ctx *gin.Context) bool { return true },
		Erlog:              erlogs.New(erlogs.Log(true), erlogs.Msg("request forbidden!")),
		Methods:            []string{http.MethodPost, http.MethodDelete, http.MethodPut},
		CookieKey:          CSRFTokenCookieName,
		HeaderKey:          CSRFTokenHeaderName,
		HeaderExpireKey:    CSRFTokenExpireHeaderName,
		HeaderRequestIdKey: ClientRequestIdHeaderName,
		HeaderDidKey:       ClientDidHeaderName,
		HeaderCTKey:        ClientCTHeaderName,
		HeaderCVKey:        ClientCVHeaderName,
		Path:               "/",
		Domain:             "*",
		Secure:             true,
		HttpOnly:           true,
		MaxAge:             7200,
		SameSite:           int64(http.SameSiteStrictMode),
	}
}

// SetupVerifyCSRFToken 设置&校验CSRF-Token
func SetupVerifyCSRFToken(eng *gin.Engine, opts ...CSRFOption) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		config := defaultCSRFConfig()
		for _, opt := range opts {
			opt.Apply(config)
		}

		key := ctx.Request.Header.Get(config.HeaderRequestIdKey)
		if len(key) == 0 {
			key = lo.RandomString(32, lo.AlphanumericCharset)
		}

		reqCtx := ctx.Request.Context()
		data, err, _ := sg.Do(key, func() (res any, err error) {
			return verifyCSRFToken(ctx, config)
		})

		if err != nil {
			ResponseError(ctx, err)
			return
		}

		dtv, ok := data.(*cljent.DTV)
		if !ok {
			err = config.Erlog.Erorr(reqCtx, erlogs.Content("dtv is nil"))
			ResponseError(ctx, err)
			return
		}

		value, err := cljent.BuildClientDTV(dtv, config.Crypter)
		if err != nil {
			err = config.Erlog.Erorr(ctx.Request.Context(), erlogs.ContentError(err))
			ResponseError(ctx, err)
			return
		}

		if len(value) == 0 {
			err = config.Erlog.Erorr(ctx.Request.Context(), erlogs.Content("failed to build client dtv"))
			ResponseError(ctx, err)
			return
		}

		clientInfo, err := cljent.NewClientInfo(ctx, dtv)
		if err != nil {
			err = config.Erlog.Erorr(ctx.Request.Context(), erlogs.ContentError(err), erlogs.Fields(zap.Any("dtv", dtv)))
			ResponseError(ctx, err)
			return
		}

		if clientInfo == nil {
			err = config.Erlog.Erorr(ctx.Request.Context(), erlogs.Content("client info is nil"))
			ResponseError(ctx, err)
			return
		}

		ctx.Set(ClientInfoContextKey, clientInfo)
		ctx.Header(config.HeaderKey, value)
		ctx.Header(config.HeaderExpireKey, cast.ToString(config.MaxAge))
		ctx.Header("Access-Control-Expose-Headers", fmt.Sprintf("%s,%s", config.HeaderKey, config.HeaderExpireKey))
		ctx.SetCookie(config.CookieKey, value, config.MaxAge, config.Path, config.Domain, config.Secure, config.HttpOnly)
		ctx.SetSameSite(http.SameSite(config.SameSite))

		ctx.Next()
	}
}

func verifyCSRFToken(ctx *gin.Context, config *CSRFConfig) (res *cljent.DTV, err error) {
	validReferer := false
	referer := ctx.Request.Referer()
	for _, v := range config.RefererPrefixs {
		if strings.HasPrefix(referer, v) {
			validReferer = true
			break
		}
	}

	origin := ctx.Request.Header.Get("Origin")
	validOrigin := slices.Contains(config.Origins, origin)

	if !validOrigin || !validReferer {
		err = config.Erlog.Erorr(ctx.Request.Context(),
			erlogs.Content("origin|referer invalid"),
			erlogs.Fields(
				zap.String("origin", origin),
				zap.Strings("origins", config.Origins),
				zap.String("referer", referer),
				zap.Strings("refererPrefixs", config.RefererPrefixs),
			),
		)
		return
	}

	if config.Crypter == nil {
		err = config.Erlog.Erorr(ctx.Request.Context(), erlogs.Content("Crypter is nil"))
		return
	}

	dtv := &cljent.DTV{
		ClientType:    cast.ToUint8(ctx.GetHeader(config.HeaderCTKey)),
		ClientVersion: ctx.GetHeader(config.HeaderCVKey),
		Did:           ctx.GetHeader(config.HeaderDidKey),
		Salt:          crypto.MD5([]any{"csrf-token-dtv", lo.RandomString(32, lo.AllCharset), time.Now().UnixNano()}),
	}

	if !cljent.IsValidClientType(dtv.ClientType) ||
		!cljent.IsValidClientVersion(dtv.ClientVersion) || len(dtv.Did) != DidLength {
		err = config.Erlog.Erorr(ctx.Request.Context(),
			erlogs.Content("ct||cv||did invalid"),
			erlogs.Fields(zap.Any("dtv", dtv)),
		)
		return
	}

	// 校验
	if config.Verify == nil || config.Verify(ctx) {
		cookieValue, _ := ctx.Cookie(config.CookieKey)
		if len(cookieValue) > 0 {
			cookieValue, err = url.QueryUnescape(cookieValue)
		}

		if err != nil {
			err = config.Erlog.Erorr(ctx.Request.Context(), erlogs.ContentError(err))
			return
		}

		headerValue := ctx.Request.Header.Get(config.HeaderKey)
		if len(cookieValue) == 0 || len(headerValue) == 0 {
			err = config.Erlog.Erorr(ctx.Request.Context(),
				erlogs.Content("cookie|header value is empty"),
				erlogs.Fields(
					zap.String("cookieValue", cookieValue),
					zap.String("headerValue", headerValue),
				),
			)
			return
		}

		// 校验value
		if !strjngs.ConstantTimeCompare(cookieValue, headerValue) {
			err = config.Erlog.Erorr(ctx.Request.Context(),
				erlogs.Content("cookieValue != headerValue"),
				erlogs.Fields(
					zap.String("cookieValue", cookieValue),
					zap.String("headerValue", headerValue),
				),
			)
			return
		}

		pr, er := cljent.ParseClientDTV(cookieValue, config.Crypter)
		if err = er; err != nil {
			err = config.Erlog.Erorr(ctx.Request.Context(),
				erlogs.ContentError(err),
				erlogs.Fields(zap.Any("cookieValue", cookieValue)),
			)
			return
		}

		// 校验dtv
		if !cljent.CompareCientDTV(dtv, pr) {
			err = config.Erlog.Erorr(ctx.Request.Context(),
				erlogs.Content("dtv is invalid"),
				erlogs.Fields(zap.Any("clientDtv", dtv), zap.Any("cookieDtv", pr)),
			)
			return
		}

		// 校验Csrf-Token
		for _, v := range config.Methods {
			if ctx.Request.Method == v {
				headerValue := ctx.Request.Header.Get(config.HeaderKey)
				if len(headerValue) == 0 || !strjngs.ConstantTimeCompare(cookieValue, headerValue) {
					err = config.Erlog.Erorr(ctx.Request.Context(),
						erlogs.Content("csrf-token invalid"),
						erlogs.Fields(
							zap.String("cookieCsrfToken", cookieValue),
							zap.String("headerCsrfToken", headerValue),
						),
					)
					return
				}
			}
		}
	}

	res = dtv
	return
}
