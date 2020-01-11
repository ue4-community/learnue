package setting

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var config *viper.Viper

var(
	ROOT        string
	ConfigPath  string
	TemplateDir string
)

func init() {
	ROOT = GetProjectPath()
	ConfigPath = ROOT + "/config"
	TemplateDir = ROOT + "/template/"

	v := viper.New()

	envFilePath := GetProjectEnvFilePath()

	if !fileExist(envFilePath){
		fmt.Printf("默认配置不存在,新建配置并设置默认值")
		v.SetDefault("global.is_master", false)
		v.SetDefault("global.log_level", "DEBUG")
		v.SetDefault("global.pprof", "127.0.0.1:8096")
		v.SetDefault("crawl.spec", "0 0 */1 * * ?")

		v.SetDefault("listen.host", "127.0.0.1")
		v.SetDefault("listen.port", "8088")

		v.SetDefault("mysql.max_idle", 2)
		v.SetDefault("mysql.max_conn", 10)
		v.SetDefault("xorm.show_sql", false)
		v.SetDefault("xorm.log_level", 1)

		v.SetDefault("qiniu.http_domain", "http://test.static.studygolang.com/")
		v.SetDefault("qiniu.https_domain", "https://static.studygolang.com/")

		v.SetDefault("feed.day", 3)
		v.SetDefault("account.verify_email", false)

		v.SetDefault("crawl.contain_link", 10)

		v.SetDefault("feed.day", 3)
		v.SetDefault("feed.cmt_weight", 80)
		v.SetDefault("feed.view_weight", 80)
		v.SetDefault("feed.day", 3)
		v.SetDefault("feed.cmt_weight", 80)
		v.SetDefault("feed.view_weight", 80)
		v.SetDefault("feed.like_weight", 60)

		v.SetDefault("qiniu.up_host", "https://up-z2.qiniup.com")
		if err := v.SafeWriteConfigAs(envFilePath);err != nil{
			panic(fmt.Sprintf("创建默认配置失败!错误:%s",err.Error()))
		}
		config = v
	}

	v.SetConfigFile(envFilePath)

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件错误:%s \n", err))
	}

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置发生变更,已经重新加载", e.Name)
	})

	config = v
	fmt.Println("配置加载完毕!")
}

func Get() *viper.Viper {
	return config
}

func App() *viper.Viper {
	app := config.Sub("app")
	return app
}

func GetAppDomain() string {
	production := App().GetBool("production")
	if production {
		return App().GetString("domain")
	} else {
		return App().GetString("devDomain")
	}
}

func Server() *viper.Viper {
	server := config.Sub("server")
	return server
}

func Mysql() *viper.Viper {
	mysql := config.Sub("mysql")
	return mysql
}

func Redis() *viper.Viper {
	redis := config.Sub("redis")
	return redis
}

func Xorm() *viper.Viper {
	xorm := config.Sub("xorm")
	return xorm
}

func Smtp() *viper.Viper {
	smtp := config.Sub("smtp")
	return smtp
}

func Sms() *viper.Viper {
	sms := config.Sub("sms")
	return sms
}

func OAuth2() *viper.Viper {
	oauth2 := config.Sub("oauth2")
	return oauth2
}

func Jwt() *viper.Viper {
	jwt := config.Sub("jwt")
	return jwt
}

func Oss() *viper.Viper {
	oss := config.Sub("oss")
	return oss
}

func GetOssCallbackUrl() string {
	production := App().GetBool("production")
	oss := Oss()
	if production {
		return oss.GetString("callbackUrl")
	} else {
		return oss.GetString("devCallbackUrl")
	}
}
func GetOssCdnDomain(withSlash bool) string {
	domain := Oss().GetString("cdnDomain")

	if withSlash {
		domain = domain + "/"
	}
	return domain
}

const LEARNUE_HOME = "LEARNUE_HOME"

func GetProjectPath() string {
	home, ok := os.LookupEnv(LEARNUE_HOME)

	if ok {
		return filepath.ToSlash(home)
	} else {
		return "./"
	}
}

func GetConfigPath()string  {
	return filepath.Join(GetProjectPath(), "config")

}

func GetProjectEnvFilePath() string {
	return filepath.Join(GetProjectPath(), "config/env.yml")
}

func GetPublicUploadDirectory() string {
	path := filepath.Join(GetProjectPath(), "static/files")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModeDir)
	}
	return path
}

func GetPrivateUploadDirectory() string {
	path := filepath.Join(GetProjectPath(), "data/files")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModeDir)
	}
	return path
}


// fileExist 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}