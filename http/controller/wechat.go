// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of self source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"io/ioutil"
	"net/http"

	"github.com/studygolang/studygolang/logic"
	"github.com/studygolang/studygolang/modules/context"

	echo "github.com/labstack/echo/v4"
)

type WechatController struct{}

// 注册路由
func (self WechatController) RegisterRoute(g *echo.Group) {
	g.Any("/wechat/autoreply", self.AutoReply)
}

func (self WechatController) AutoReply(ctx echo.Context) error {
	// 配置微信（不校验，直接返回成功）
	if ctx.QueryParam("echostr") != "" {
		return ctx.String(http.StatusOK, ctx.QueryParam("echostr"))
	}

	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return ctx.String(http.StatusOK, "")
	}

	if len(body) == 0 {
		return ctx.String(http.StatusOK, "")
	}

	wechatReply, err := logic.DefaultWechat.AutoReply(context.EchoContext(ctx), body)
	if err != nil {
		return ctx.String(http.StatusOK, "")
	}

	return ctx.XML(http.StatusOK, wechatReply)
}
