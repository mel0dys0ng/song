package erlogs

var (
	// BaseBiz 基础业务线
	BaseBiz          = OptionBiz(0, "base")
	ErLogsBiz        = OptionBiz(1, "erlogs")
	MetasBiz         = OptionBiz(2, "metas")
	VipersBiz        = OptionBiz(3, "vipers")
	CobrasBiz        = OptionBiz(4, "cobras")
	HttpsBiz         = OptionBiz(5, "https")
	ClientsMySQLBiz  = OptionBiz(6, "clients_mysql")
	ClientsRedisBiz  = OptionBiz(7, "clients_redis")
	ClientsRestyBiz  = OptionBiz(8, "clients_resty")
	ClientsPubSubBiz = OptionBiz(9, "clients_pubsub")

	// BaseEL 基础日志记录器
	BaseEL          = WithOptions(BaseBiz)
	ErLogsEL        = WithOptions(ErLogsBiz)
	MetasEL         = WithOptions(MetasBiz)
	VipersEL        = WithOptions(VipersBiz)
	CobrasEL        = WithOptions(CobrasBiz)
	HttpsEL         = WithOptions(HttpsBiz)
	ClientsMySQLEL  = WithOptions(ClientsMySQLBiz)
	ClientsRedisEL  = WithOptions(ClientsRedisBiz)
	ClientsRestyEL  = WithOptions(ClientsRestyBiz)
	ClientsPubSubEL = WithOptions(ClientsPubSubBiz)

	// 成功状态码 0
	Ok = BaseEL.Status(0, "ok")

	// 客户端错误 40000 ～ 49999
	BadRequest       = BaseEL.Status(40000, "请求错误")
	InvalidArguments = BaseEL.Status(40001, "请求参数错误")
	Unauthorized     = BaseEL.Status(40002, "请登录授权")
	FrequencyLimit   = BaseEL.Status(40003, "操作频率过快")
	TooManyRequests  = BaseEL.Status(40004, "请求次数过多")
	InvalidCSRFToken = BaseEL.Status(40005, "请求非法")
	InvalidSign      = BaseEL.Status(40006, "签名无效")

	// 服务端错误，50000 ～ 59999
	ServerError   = BaseEL.Status(50000, "服务错误，请稍后重试")
	InvalidParams = BaseEL.Status(50001, "服务参数错误，请检查参数")
	MySQLError    = BaseEL.Status(50002, "数据存储异常，请稍后重试")
	CacheError    = BaseEL.Status(50003, "缓存异常，请稍后重试")
	ClientError   = BaseEL.Status(50004, "组件异常，请稍后重试")
)
