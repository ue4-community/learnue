// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"net/http"

	. "github.com/ue4-community/learnue/http"
	"github.com/ue4-community/learnue/logic"
	"github.com/ue4-community/learnue/model"
	"github.com/ue4-community/learnue/modules/context"

	echo "github.com/labstack/echo/v4"
)

type OAuthController struct{}

// 注册路由
func (self OAuthController) RegisterRoute(g *echo.Group) {
	g.GET("/oauth/github/callback", self.GithubCallback)
	g.GET("/oauth/github/login", self.GithubLogin)

	g.GET("/oauth/gitea/callback", self.GiteaCallback)
	g.GET("/oauth/gitea/login", self.GiteaLogin)
}

func (OAuthController) GithubLogin(ctx echo.Context) error {
	uri := ctx.QueryParam("uri")
	url := logic.DefaultThirdUser.GithubAuthCodeUrl(context.EchoContext(ctx), uri)
	return ctx.Redirect(http.StatusSeeOther, url)
}

func (OAuthController) GithubCallback(ctx echo.Context) error {
	code := ctx.FormValue("code")

	me, ok := ctx.Get("user").(*model.Me)
	if ok {
		// 已登录用户，绑定 github
		logic.DefaultThirdUser.BindGithub(context.EchoContext(ctx), code, me)

		redirectURL := ctx.QueryParam("redirect_url")
		if redirectURL == "" {
			redirectURL = "/account/edit#connection"
		}
		return ctx.Redirect(http.StatusSeeOther, redirectURL)
	}

	user, err := logic.DefaultThirdUser.LoginFromGithub(context.EchoContext(ctx), code)
	if err != nil || user.Uid == 0 {
		var errMsg = ""
		if err != nil {
			errMsg = err.Error()
		} else {
			errMsg = "服务内部错误"
		}

		return render(ctx, "login.html", map[string]interface{}{"error": errMsg})
	}

	// 登录成功，种cookie
	SetLoginCookie(ctx, user.Username)

	if user.Balance == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/balance")
	}

	return ctx.Redirect(http.StatusSeeOther, "/")
}

func (OAuthController) GiteaLogin(ctx echo.Context) error {
	uri := ctx.QueryParam("uri")
	url := logic.DefaultThirdUser.GiteaAuthCodeUrl(context.EchoContext(ctx), uri)
	return ctx.Redirect(http.StatusSeeOther, url)
}

func (OAuthController) GiteaCallback(ctx echo.Context) error {
	code := ctx.FormValue("code")

	me, ok := ctx.Get("user").(*model.Me)
	if ok {
		// 已登录用户，绑定 github
		logic.DefaultThirdUser.BindGitea(context.EchoContext(ctx), code, me)

		redirectURL := ctx.QueryParam("redirect_url")
		if redirectURL == "" {
			redirectURL = "/account/edit#connection"
		}
		return ctx.Redirect(http.StatusSeeOther, redirectURL)
	}

	user, err := logic.DefaultThirdUser.LoginFromGitea(context.EchoContext(ctx), code)
	if err != nil || user.Uid == 0 {
		var errMsg = ""
		if err != nil {
			errMsg = err.Error()
		} else {
			errMsg = "服务内部错误"
		}

		return render(ctx, "login.html", map[string]interface{}{"error": errMsg})
	}

	// 登录成功，种cookie
	SetLoginCookie(ctx, user.Username)

	if user.Balance == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/balance")
	}

	return ctx.Redirect(http.StatusSeeOther, "/")
}
