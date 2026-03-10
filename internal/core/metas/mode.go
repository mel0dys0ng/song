package metas

const (
	// ModeLocal 本地环境
	ModeLocal ModeType = "local"
	// ModeTest 测试环境
	ModeTest ModeType = "test"
	// ModeStaging  预发布环境
	ModeStaging ModeType = "staging"
	// ModeProduction 生产环境
	ModeProduction ModeType = "prod"
)

// ModeType 运行模式，如：local, test, staging, prod
type ModeType string

func (mt ModeType) Validate() bool {
	switch mt {
	case ModeLocal, ModeTest, ModeStaging, ModeProduction:
		return true
	default:
		return false
	}
}

func (mt ModeType) IsModeLocal() bool {
	return mt == ModeLocal
}

func (mt ModeType) IsModeTest() bool {
	return mt == ModeTest
}

func (mt ModeType) IsModeStaging() bool {
	return mt == ModeStaging
}

func (mt ModeType) IsModeProduction() bool {
	return mt == ModeProduction
}

func (mt ModeType) IsModeDebug() bool {
	return mt == ModeLocal || mt == ModeTest
}

func (mt ModeType) IsModeOnline() bool {
	return mt == ModeStaging || mt == ModeProduction
}

func (mt ModeType) String() string {
	return string(mt)
}
