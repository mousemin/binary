package main

var config Config

type Config struct {
	Host     string  `mapstructure:"host"`      // redis-server host
	Port     int     `mapstructure:"port"`      // redis-server port
	Password string  `mapstructure:"password"`  // redis-server password
	Database int     `mapstructure:"database"`  // redis-server database
	MemScale float64 `mapstructure:"mem_scale"` // redis-server mem scale
	Expire   int64   `mapstructure:"expire"`    // redis-server clean start expire time
	Interval int64   `mapstructure:"interval"`  // redis-server clean interval expire time
}
