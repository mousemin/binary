package main

import (
	"flag"
	"git.mousemin.com/mousemin/binary/pkg/osutil"
	"git.mousemin.com/mousemin/binary/pkg/shell"
	"git.mousemin.com/mousemin/binary/pkg/slices"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)


func parseFile(pathStr, joinPath string) error {
	files, err := ioutil.ReadDir(pathStr)
	if err != nil {
		return err
	}

	for _, file := range files {
		filename := file.Name()
		// 忽略文件
		if slices.InSliceStr(filename, config.ExcludeFile) {
			continue
		}
		ext := path.Ext(filename)
		// 忽略尾缀
		if slices.InSliceStr(ext, config.ExcludeExt) {
			continue
		}
		if len(config.IncludeExt) != 0 && !slices.InSliceStr(ext, config.IncludeExt) {
			continue
		}
		if file.IsDir() {
			dirname := filepath.Join(pathStr, filename)
			if !slices.InSliceStr(dirname, config.ExcludeDir) {
				if err := parseFile(dirname, filepath.Join(joinPath, filename)); err != nil {
					return err
				}
			}
			continue
		}
		newpath := filepath.Join(config.Path, joinPath)
		if !osutil.IsDir(newpath) {
			if err := os.MkdirAll(newpath, os.ModeDevice); err != nil {
				return err
			}
		}
		filepath1 := filepath.Join(pathStr, filename)
		filepath2 := filepath.Join(newpath, filename)
		if err := shell.Mv(filepath1, filepath2); err != nil {
			log.Println(filepath1, " => ", filepath2, " err: ", err)
			return err
		}
		log.Println(filepath1, " => ", filepath2, " success")
	}

	if len(joinPath) != 0 {
		fps, _ := ioutil.ReadDir(pathStr)
		if len(fps) == 0 {
			if err := os.Remove(pathStr); err != nil {
				log.Println(pathStr, " delete failed, err: ", err)
			}
		}
	}
	return nil
}



func main() {
	var file string
	flag.StringVar(&file, "f", ".sync.yml", "配置文件")
	flag.Parse()
	pwd, err := os.Getwd()
	if err != nil {
		panic("获取当前文件失败, err: " + err.Error())
	}
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

	if config.Delay == 0 {
		config.Delay = 300
	}

	for i, dir := range config.ExcludeDir {
		abs, err := filepath.Abs(dir)
		if err != nil {
			continue
		}
		config.ExcludeDir[i] = abs
	}
	if len(config.Log) != 0 {
		logFile, err := filepath.Abs(config.Log)
		if err != nil {
			panic("日志文件获取文件路径失败, err: " + err.Error())
		}
		fp, err := os.OpenFile(logFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			panic("日志文件初始化失败")
		}
		defer fp.Close()
		log.SetOutput(fp)
		log.SetFlags(log.Ldate| log.Ltime | log.Lshortfile)
	}
	for {
		if err := parseFile(pwd, ""); err != nil {
			panic(err)
		}
		log.Println("休眠", config.Delay, "秒")
		time.Sleep(time.Duration(config.Delay) * time.Second)
	}
}
