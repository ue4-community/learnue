// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package app

import (
	"github.com/studygolang/studygolang/http/middleware"
	"github.com/studygolang/studygolang/logic"
	"github.com/studygolang/studygolang/model"
	"github.com/studygolang/studygolang/modules/context"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

type CommentController struct{}

func (self CommentController) RegisterRoute(g *echo.Group) {
	g.POST("/comment/:objid", self.Create, middleware.NeedLogin(), middleware.Sensivite(), middleware.PublishNotice())
}

// Create 评论（或回复）
func (CommentController) Create(ctx echo.Context) error {
	user := ctx.Get("user").(*model.Me)

	// 入库
	objid := goutils.MustInt(ctx.Param("objid"))
	if objid == 0 {
		return fail(ctx, "参数有误，请刷新后重试！", 1)
	}
	forms, _ := ctx.FormParams()
	comment, err := logic.DefaultComment.Publish(context.EchoContext(ctx), user.Uid, objid, forms)
	if err != nil {
		return fail(ctx, "服务器内部错误", 2)
	}

	return success(ctx, map[string]interface{}{"comment": comment})
}
