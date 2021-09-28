package main

import (
	"flag"
	"fmt"
	"net"
	"path/filepath"
	"strconv"
	"strings"

	"git.mousemin.com/mousemin/binary/pkg/osutil"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
)

func cleanRedis() error {
	// 链接redis
	options := make([]redis.DialOption, 0, 8)
	if config.Database != 0 {
		options = append(options, redis.DialDatabase(config.Database))
	}
	if len(config.Password) != 0 {
		options = append(options, redis.DialPassword(config.Password))
	}
	redisFp, err := redis.Dial("tcp", net.JoinHostPort(config.Host, strconv.Itoa(config.Port)), options...)
	if err != nil {
		return err
	}
	defer redisFp.Close()
	minExpire := config.Expire
	for {
		infos, err := redis.String(redisFp.Do("INFO"))
		if err != nil {
			return err
		}
		infoSlice := strings.Split(infos, "\r\n")
		infoMap := make(map[string]string, len(infoSlice))
		for _, info := range infoSlice {
			if len(info) == 0 || info[0] == '#' {
				continue
			}
			tmpSlice := strings.Split(info, ":")
			if len(tmpSlice) != 2 {
				continue
			}
			infoMap[tmpSlice[0]] = tmpSlice[1]
		}
		maxMem, err := strconv.ParseFloat(infoMap["maxmemory"], 10)
		if err != nil {
			return err
		}
		useMem, err := strconv.ParseFloat(infoMap["used_memory"], 10)
		if err != nil {
			return err
		}
		fmt.Println("redis mem:", useMem/maxMem)
		if (useMem / maxMem) < config.MemScale {
			return nil
		}
		cursor := int64(0)
		for {
			cursorNew, keys, err := ScanValues(redisFp.Do("SCAN", cursor, "COUNT", 10000))
			if err != nil {
				return err
			}
			if cursorNew == 0 {
				break
			}
			cursor = cursorNew
			for _, key := range keys {
				expire, err := redis.Int64(redisFp.Do("TTL", key))
				if err != nil {
					continue
				}
				if expire < minExpire {
					_, _ = redisFp.Do("DEL", key)
				}
			}
		}
		if config.Interval == 0 {
			break
		}
		minExpire = minExpire + config.Interval
	}
	return nil
}

func ScanValues(v interface{}, err error) (int64, []string, error) {
	if err != nil {
		return 0, nil, err
	}
	ret, err := redis.Values(v, err)
	if err != nil || len(ret) != 2 {
		return 0, nil, err
	}
	cursorNew, err := redis.Int64(ret[0], nil)
	if err != nil {
		return 0, nil, err
	}
	data, err := redis.Strings(ret[1], nil)
	return cursorNew, data, err
}

func main() {
	var file string
	flag.StringVar(&file, "f", ".cache.yml", "配置文件")
	flag.Parse()
	// 获取当前执行文件夹
	f, err := filepath.Abs(file)
	if err != nil {
		panic("获取配置文件路径失败, err: " + err.Error())
	}
	if !osutil.IsFile(f) {
		panic("配置 " + file + " 文件不存在")
	}

	// 加载配置文件
	viper.SetConfigFile(f)
	if err := viper.ReadInConfig(); err != nil {
		panic("配置文件加载失败, err: " + err.Error())
	}

	if err := viper.Unmarshal(&config); err != nil {
		panic("配置文件加载失败, err: " + err.Error())
	}
	if err := cleanRedis(); err != nil {
		panic("redis清理失败, err: " + err.Error())
	}
}
