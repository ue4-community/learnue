// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"runtime"
	"strconv"

	"github.com/studygolang/studygolang/db"
	"github.com/studygolang/studygolang/global"
	"github.com/studygolang/studygolang/logic"
	"github.com/studygolang/studygolang/model"
	"github.com/studygolang/studygolang/modules/context"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

type InstallController struct{}

// 注册路由
func (self InstallController) RegisterRoute(g *echo.Group) {
	g.GET("/install", self.SetupConfig)
	g.Match([]string{"GET", "POST"}, "/install/setup-config", self.SetupConfig)
	g.Match([]string{"GET", "POST"}, "/install/do", self.DoInstall)
	g.Match([]string{"GET", "POST"}, "/install/options", self.SetupOptions)
}

func (self InstallController) SetupConfig(ctx echo.Context) error {
	// config/env.ini 存在
	if db.MasterDB != nil {
		if logic.DefaultInstall.IsTableExist(context.EchoContext(ctx)) {
			return ctx.Redirect(http.StatusSeeOther, "/")
		}
		return ctx.Redirect(http.StatusSeeOther, "/install/do")
	}

	step := goutils.MustInt(ctx.QueryParam("step"))
	if step == 2 {
		err := self.genConfig(ctx)
		if err != nil {
			data := map[string]interface{}{
				"dbhost":   ctx.FormValue("dbhost"),
				"dbport":   ctx.FormValue("dbport"),
				"dbname":   ctx.FormValue("dbname"),
				"uname":    ctx.FormValue("uname"),
				"err_type": 1,
			}

			if err == db.ConnectDBErr {
				data["err_type"] = 1
			} else if err == db.UseDBErr {
				data["err_type"] = 2
			}

			return renderInstall(ctx, "install/setup-err.html", data)
		}
	}
	return renderInstall(ctx, "install/setup-config.html", map[string]interface{}{"step": step})
}

// DoInstall 执行安装，包括站点简单配置，安装数据库（创建数据库、表，填充基本数据）等
func (self InstallController) DoInstall(ctx echo.Context) error {
	if db.MasterDB == nil {
		return ctx.Redirect(http.StatusSeeOther, "/install")
	}

	if logic.DefaultInstall.IsTableExist(context.EchoContext(ctx)) {
		if logic.DefaultInstall.HadRootUser(context.EchoContext(ctx)) {
			return ctx.Redirect(http.StatusSeeOther, "/")
		}
	}

	step := goutils.MustInt(ctx.QueryParam("step"), 1)
	data := map[string]interface{}{
		"user_name":   "admin",
		"admin_email": "",
		"step":        step,
	}

	if step == 2 {
		username := ctx.FormValue("user_name")
		email := ctx.FormValue("admin_email")
		password1 := ctx.FormValue("admin_password")
		password2 := ctx.FormValue("admin_password2")

		if username == "" || email == "" {
			data["err"] = "用户名和邮箱不能留空"
			return renderInstall(ctx, "install/install.html", data)
		}

		data["user_name"] = username
		data["admin_email"] = email

		if password1 != password2 {
			data["err"] = "两次输入的密码不一致"
			return renderInstall(ctx, "install/install.html", data)
		}

		err := logic.DefaultInstall.CreateTable(context.EchoContext(ctx))
		if err != nil {
			data["err"] = "创建数据表失败！"
			return renderInstall(ctx, "install/install.html", data)
		}

		err = logic.DefaultInstall.InitTable(context.EchoContext(ctx))
		if err != nil {
			data["err"] = "初始化数据表失败！"
			return renderInstall(ctx, "install/install.html", data)
		}

		if password1 == "" {
			password1 = goutils.RandString(12)
			data["passwd"] = password1
		}

		// 创建管理员
		form := url.Values{
			"username": {username},
			"email":    {email},
			"passwd":   {password1},
			"is_root":  {"true"},
			"status":   {strconv.Itoa(model.UserStatusAudit)},
		}
		errMsg, err := logic.DefaultUser.CreateUser(context.EchoContext(ctx), form)
		if err != nil {
			data["err"] = errMsg
			return renderInstall(ctx, "install/install.html", data)
		}

		data["step"] = 3

		data["os"] = runtime.GOOS

		// 为了保证程序正常，需要重启
		go self.reload()
	}
	return renderInstall(ctx, "install/install.html", data)
}

func (InstallController) SetupOptions(ctx echo.Context) error {
	var (
		noEmailConf = false
		noQiniuConf = false
	)

	if db.ConfigFile.GetString("email.smtp_username") == "" {
		noEmailConf = true
	}

	if db.ConfigFile.GetString("qiniu.access_key") == "" {
		noQiniuConf = true
	}

	if !noEmailConf && !noQiniuConf {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	if ctx.Request().Method == "POST" {
		//config.ConfigFile.SetSectionComments("email", "用于注册发送激活码等")
		emailFields := []string{"smtp_host", "smtp_port", "smtp_username", "smtp_password", "from_email"}
		for _, field := range emailFields {
			if field == "smtp_port" && ctx.FormValue("smtp_port") == "" {
				//config.ConfigFile.SetValue("email", field, "25")
				db.ConfigFile.Set(fmt.Sprintf("email.%s", field), "25")
			} else {
				db.ConfigFile.Set(fmt.Sprintf("email.%s", field), ctx.FormValue(field))
			}
		}

		//config.ConfigFile.SetSectionComments("qiniu", "图片存储在七牛云，如果没有可以通过 https://portal.qiniu.com/signup?code=3lfz4at7pxfma 免费申请")
		qiniuFields := []string{"access_key", "secret_key", "bucket_name", "http_domain", "https_domain"}
		for _, field := range qiniuFields {
			db.ConfigFile.Set(fmt.Sprintf("qiniu.%s", field), ctx.FormValue(field))
		}
		if ctx.FormValue("https_domain") == "" {
			db.ConfigFile.Set("qiniu.https_domain", ctx.FormValue("http_domain"))
		}

		if err := db.ConfigFile.SafeWriteConfigAs(db.ConfigPath); err != nil {
			fmt.Println("写入配置失败:" + err.Error())
		}

		return renderInstall(ctx, "install/setup-options.html", map[string]interface{}{"success": true})
	}

	data := map[string]interface{}{
		"no_email_conf": noEmailConf,
		"no_qiniu_conf": noQiniuConf,
	}
	return renderInstall(ctx, "install/setup-options.html", data)
}

func (InstallController) genConfig(ctx echo.Context) error {
	env := ctx.FormValue("env")

	//config.ConfigFile.SetSectionComments("global", "")
	db.ConfigFile.Set("global.env", env)

	var (
		logLevel = "DEBUG"
		// domain       = global.App.Host + ":" + global.App.Port
		xormLogLevel = "0"
		xormShowSql  = "true"
	)
	if env == "pro" {
		logLevel = "INFO"
		xormLogLevel = "1"
		xormShowSql = "false"
	}

	db.ConfigFile.Set("global.log_level", logLevel)
	db.ConfigFile.Set("global.cookie_secret", goutils.RandString(10))
	db.ConfigFile.Set("global.data_path", "data/max_online_num")

	//config.ConfigFile.SetSectionComments("listen", "")
	db.ConfigFile.Set("listen.host", "")
	db.ConfigFile.Set("listen.port", global.App.Port)

	dbname := ctx.FormValue("dbname")
	uname := ctx.FormValue("uname")
	pwd := ctx.FormValue("pwd")
	dbhost := ctx.FormValue("dbhost")
	dbport := ctx.FormValue("dbport")

	//config.ConfigFile.SetSectionComments("mysql", "")
	db.ConfigFile.Set("mysql.host", dbhost)
	db.ConfigFile.Set("mysql.port", dbport)
	db.ConfigFile.Set("mysql.user", uname)
	db.ConfigFile.Set("mysql.password", pwd)
	db.ConfigFile.Set("mysql.dbname", dbname)
	db.ConfigFile.Set("mysql.charset", "utf8")
	//config.ConfigFile.SetKeyComments("mysql.max_idle", "最大空闲连接数")
	db.ConfigFile.Set("mysql.max_idle", "2")
	//config.ConfigFile.SetKeyComments("mysql.max_conn", "最大打开连接数")
	db.ConfigFile.Set("mysql.max_conn", "10")

	//config.ConfigFile.SetSectionComments("xorm.")
	db.ConfigFile.Set("xorm.show_sql", xormShowSql)
	//config.ConfigFile.SetKeyComments("xorm.log_level", "0-debug, 1-info, 2-warning, 3-error, 4-off, 5-unknow")
	db.ConfigFile.Set("xorm.log_level", xormLogLevel)

	//config.ConfigFile.SetSectionComments("security.")
	//config.ConfigFile.SetKeyComments("security.unsubscribe_token_key", "退订邮件使用的 token key")
	db.ConfigFile.Set("security.unsubscribe_token_key", goutils.RandString(18))
	//config.ConfigFile.SetKeyComments("security.activate_sign_salt", "注册激活邮件使用的 sign salt")
	db.ConfigFile.Set("security.activate_sign_salt", goutils.RandString(18))

	//config.ConfigFile.SetSectionComments("sensitive", "过滤广告")
	//config.ConfigFile.SetKeyComments("sensitive.title", "标题关键词")
	db.ConfigFile.Set("sensitive.title", "")
	//config.ConfigFile.SetKeyComments("sensitive.content", "内容关键词")
	db.ConfigFile.Set("sensitive.content", "")

	//config.ConfigFile.SetSectionComments("search.搜索配置")
	db.ConfigFile.Set("search.engine_url", "")

	// 校验数据库配置是否正确有效
	if err := db.TestDB(); err != nil {
		return err
	}
	if err := db.ConfigFile.SafeWriteConfigAs(db.ConfigPath); err != nil {
		fmt.Println("写入配置失败:" + err.Error())
	}
	return nil
}

func renderInstall(ctx echo.Context, filename string, data map[string]interface{}) error {
	objLog := getLogger(ctx)

	if data == nil {
		data = make(map[string]interface{})
	}

	filename = db.TemplateDir + filename

	requestURI := ctx.Request().RequestURI
	tpl, err := template.ParseFiles(filename)
	if err != nil {
		objLog.Errorf("解析模板出错（ParseFiles）：[%q] %s\n", requestURI, err)
		return err
	}

	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, data)
	if err != nil {
		objLog.Errorf("执行模板出错（Execute）：[%q] %s\n", requestURI, err)
		return err
	}

	return ctx.HTML(http.StatusOK, buf.String())
}
