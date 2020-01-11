// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"github.com/studygolang/studygolang/logic"
	"github.com/studygolang/studygolang/modules/context"

	echo "github.com/labstack/echo/v4"
)

type LinkController struct{}

// 注册路由
func (self LinkController) RegisterRoute(g *echo.Group) {
	g.GET("/links", self.FindLinks)
}

// FindLinks 友情链接
func (LinkController) FindLinks(ctx echo.Context) error {

	friendLinks := logic.DefaultFriendLink.FindAll(context.EchoContext(ctx))

	return render(ctx, "link.html", map[string]interface{}{"links": friendLinks})
}
