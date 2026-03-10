package metas

import (
	"github.com/mel0dys0ng/song/pkg/singleton"
	"github.com/mel0dys0ng/song/pkg/sys"
)

var singletonKeyMetadata any

func init() {
	singletonKeyMetadata = singleton.Key()
}

// Initialize 初始化元数据
func Initialize(opts *Options) MetadataInterface {
	return singleton.Once(singletonKeyMetadata, func() MetadataInterface {
		return New(opts)
	})
}

// Metadata 获取元数据
func Metadata() MetadataInterface {
	data, ok := singleton.Get[MetadataInterface](singletonKeyMetadata)
	if !ok {
		sys.Panicf("get metadata failed, please call InitializeMetadata first")
	}
	return data
}
