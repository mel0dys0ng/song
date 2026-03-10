package metas

import (
	"github.com/mel0dys0ng/song/internal/core/metas"
)

type (
	ModeType          = metas.ModeType
	KindType          = metas.KindType
	Options           = metas.Options
	MetadataInterface = metas.MetadataInterface
)

const (
	KindAPI       KindType = metas.KindAPI
	KindJob       KindType = metas.KindJob
	KindTool      KindType = metas.KindTool
	KindMessaging KindType = metas.KindMessaging

	ModeLocal      ModeType = metas.ModeLocal
	ModeTest       ModeType = metas.ModeTest
	ModeStaging    ModeType = metas.ModeStaging
	ModeProduction ModeType = metas.ModeProduction
)

// Initialize 初始化元数据
// 注意：此函数只应该被调用一次，在应用程序启动时
func Initialize(opts *Options) MetadataInterface {
	return metas.Initialize(opts)
}

// Metadata 获取元数据实例
// 在调用此函数之前，请确保已经调用了 Initialize
func Metadata() MetadataInterface {
	return metas.Metadata()
}
