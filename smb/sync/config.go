package main

var config Config

type Config struct {
	Path string `mapstructure:"path"` // 移动的路径
	Delay int64 `mapstructure:"delay"` // 每次操作的时间间隔(单位: 秒)
	IncludeExt []string `mapstructure:"include_ext"` // 包含后缀 默认是*
	ExcludeDir []string `mapstructure:"exclude_dir"` // 不包含路径
	ExcludeExt []string `mapstructure:"exclude_ext"` // 不包含的尾缀
	ExcludeFile []string `mapstructure:"exclude_file"` // 不包含文件
	Log string `mapstructure:"log"` // 日志文件
}