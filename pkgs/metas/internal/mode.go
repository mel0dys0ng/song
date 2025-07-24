package internal

const (
	// ModeDebug 调试环境
	ModeDebug ModeType = "debug"
	// ModeTest 测试环境
	ModeTest ModeType = "test"
	// ModePre  预发布环境
	ModePre ModeType = "pre"
	// ModeGray 灰度环境
	ModeGray ModeType = "gray"
	// ModeProduction 生产环境
	ModeProduction ModeType = "production"
)

// ModeType debug|test|pre|gray|production
type ModeType string

func (mt ModeType) Validate() bool {
	switch mt {
	case ModeDebug, ModeTest, ModePre, ModeGray, ModeProduction:
		return true
	default:
		return false
	}
}

func (mt ModeType) IsModeDebug() bool {
	return mt == ModeDebug
}

func (mt ModeType) IsModeTest() bool {
	return mt == ModeTest
}

func (mt ModeType) IsModePre() bool {
	return mt == ModePre
}

func (mt ModeType) IsModeGray() bool {
	return mt == ModeGray
}

func (mt ModeType) IsModeProduction() bool {
	return mt == ModeProduction
}

func (mt ModeType) IsModeGrayOrProduction() bool {
	return mt == ModeProduction || mt == ModeGray
}

func (mt ModeType) IsModeTestOrPre() bool {
	return mt == ModeTest || mt == ModePre
}

func (mt ModeType) IsModeTestOrPreOrDebug() bool {
	return mt == ModeTest || mt == ModePre || mt == ModeDebug
}

func (mt ModeType) String() string {
	return string(mt)
}
