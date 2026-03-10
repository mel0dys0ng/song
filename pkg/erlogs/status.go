package erlogs

import "github.com/mel0dys0ng/song/internal/core/erlogs"

var (
	BaseBiz          = erlogs.BaseBiz
	MetasBiz         = erlogs.MetasBiz
	VipersBiz        = erlogs.VipersBiz
	CobrasBiz        = erlogs.CobrasBiz
	HttpsBiz         = erlogs.HttpsBiz
	ClientsMySQLBiz  = erlogs.ClientsMySQLBiz
	ClientsRedisBiz  = erlogs.ClientsRedisBiz
	ClientsRestyBiz  = erlogs.ClientsRestyBiz
	ClientsPubSubBiz = erlogs.ClientsPubSubBiz

	BaseEL          = erlogs.BaseEL
	ErLogsEL        = erlogs.ErLogsEL
	MetasEL         = erlogs.MetasEL
	VipersEL        = erlogs.VipersEL
	CobrasEL        = erlogs.CobrasEL
	HttpsEL         = erlogs.HttpsEL
	ClientsMySQLEL  = erlogs.ClientsMySQLEL
	ClientsRedisEL  = erlogs.ClientsRedisEL
	ClientsRestyEL  = erlogs.ClientsRestyEL
	ClientsPubSubEL = erlogs.ClientsPubSubEL

	Ok = erlogs.Ok

	BadRequest       = erlogs.BadRequest
	InvalidArguments = erlogs.InvalidArguments
	Unauthorized     = erlogs.Unauthorized
	FrequencyLimit   = erlogs.FrequencyLimit
	TooManyRequests  = erlogs.TooManyRequests
	InvalidCSRFToken = erlogs.InvalidCSRFToken
	InvalidSign      = erlogs.InvalidSign

	ServerError   = erlogs.ServerError
	InvalidParams = erlogs.InvalidParams
	MySQLError    = erlogs.MySQLError
	CacheError    = erlogs.CacheError
	ClientError   = erlogs.ClientError
)
