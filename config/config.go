package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	ConfigFile *viper.Viper
	ConfigPath string
	ROOT       string

	TemplateDir string
)

const mainIniPath = "/config/env.yml"

func init() {
	curFilename := os.Args[0]

	binaryPath, err := exec.LookPath(curFilename)
	if err != nil {
		panic(err)
	}

	binaryPath, err = filepath.Abs(binaryPath)
	binaryPath = filepath.ToSlash(binaryPath)

	if err != nil {
		panic(err)
	}

	ROOT = filepath.ToSlash(filepath.Dir(filepath.Dir(binaryPath)))

	ConfigPath = ROOT + mainIniPath

	ConfigFile = viper.New()

	if !fileExist(ConfigPath) {
		curDir, _ := os.Getwd()
		pos := strings.LastIndex(curDir, "src")
		if pos == -1 {
			// panic("can't find " + mainIniPath)
			fmt.Println("can't find " + mainIniPath)
		} else {
			ROOT = curDir[:pos]

			ConfigPath = ROOT + mainIniPath
		}
	} else {
		ConfigFile.SetConfigFile(ConfigPath)
		err = ConfigFile.ReadInConfig()
		if err != nil {
			// panic(err)
			fmt.Println("load config file error:", err)
		}
	}

	ConfigFile.SetDefault("global.is_master", false)
	ConfigFile.SetDefault("global.log_level", "DEBUG")
	ConfigFile.SetDefault("global.pprof", "127.0.0.1:8096")
	ConfigFile.SetDefault("crawl.spec", "0 0 */1 * * ?")

	ConfigFile.SetDefault("listen.host", "127.0.0.1")
	ConfigFile.SetDefault("listen.port", "8088")

	ConfigFile.SetDefault("mysql.max_idle", 2)
	ConfigFile.SetDefault("mysql.max_conn", 10)
	ConfigFile.SetDefault("xorm.show_sql", false)
	ConfigFile.SetDefault("xorm.log_level", 1)

	ConfigFile.SetDefault("qiniu.http_domain", "http://test.static.studygolang.com/")
	ConfigFile.SetDefault("qiniu.https_domain", "https://static.studygolang.com/")

	ConfigFile.SetDefault("feed.day", 3)
	ConfigFile.SetDefault("account.verify_email", false)

	ConfigFile.SetDefault("crawl.contain_link", 10)

	ConfigFile.SetDefault("feed.day", 3)
	ConfigFile.SetDefault("feed.cmt_weight", 80)
	ConfigFile.SetDefault("feed.view_weight", 80)
	ConfigFile.SetDefault("feed.day", 3)
	ConfigFile.SetDefault("feed.cmt_weight", 80)
	ConfigFile.SetDefault("feed.view_weight", 80)
	ConfigFile.SetDefault("feed.like_weight", 60)

	ConfigFile.SetDefault("qiniu.up_host", "https://up-z2.qiniup.com")

	TemplateDir = ROOT + "/template/"

	//TODO 在这里统一设置配置默认值

}

// fileExist 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
