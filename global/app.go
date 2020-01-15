// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

/*
sets version information for the binary where it is imported.
The version can be retrieved either from the -version command line argument.

To include in a project simply import the package.

The version and compile date is stored in App variables and
are supposed to be set during compile time. Typically this is done by the
install(bash/bat). Or date is binary modify time.

To set these manually use -ldflags together with -X, like in this example:

	go install -ldflags "-X global/Build xxxxx"

*/

package global

import (
	"fmt"
	"github.com/studygolang/studygolang/modules/setting"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/studygolang/studygolang/model"
)

type app struct {
	Name    string
	Build   string
	Version string
	Date    time.Time

	Copyright string

	// 启动时间
	LaunchTime time.Time
	Uptime     time.Duration

	Env string

	Host string
	Port string

	BaseURL string

	// CDN 资源域名
	CDNHttp  string
	CDNHttps string

	Domain string

	locker sync.Mutex
}

var App = &app{}

const (
	DEV  = "dev"
	TEST = "test"
	PRO  = "pro"
)

func init() {
	App.Name = os.Args[0]
	App.Version = "V4.0.0"
	App.LaunchTime = time.Now()

	fileInfo, err := os.Stat(os.Args[0])
	if err != nil {
		panic(err)
	}

	App.Date = fileInfo.ModTime()

	App.Env = setting.Get().GetString("global.env")

	App.CDNHttp = setting.Get().GetString("qiniu.http_domain")
	App.CDNHttps = setting.Get().GetString("qiniu.https_domain")
}

func (this *app) Init(domain string) {
	do := setting.Get().GetString("global.domain")
	if do == "" {
		do = domain
	}
	this.Domain = do
}

func (this *app) SetUptime() {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.Uptime = time.Now().Sub(this.LaunchTime)
}

func (this *app) SetCopyright() {
	curYear := time.Now().Year()
	this.locker.Lock()
	defer this.locker.Unlock()
	if curYear == model.WebsiteSetting.StartYear {
		this.Copyright = fmt.Sprintf("%d %s", model.WebsiteSetting.StartYear, model.WebsiteSetting.Domain)
	} else {
		this.Copyright = fmt.Sprintf("%d-%d %s", model.WebsiteSetting.StartYear, curYear, model.WebsiteSetting.Domain)
	}
}

func (this *app) CanonicalCDN(isHTTPS bool) string {
	cdnDomain := this.CDNHttp
	if isHTTPS {
		cdnDomain = this.CDNHttps
	}
	if !strings.HasSuffix(cdnDomain, "/") {
		cdnDomain += "/"
	}

	return cdnDomain
}

func OnlineEnv() bool {
	return App.Env == PRO
}
