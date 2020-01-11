// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package app

import (
	. "github.com/studygolang/studygolang/http"
	"github.com/studygolang/studygolang/logic"
	"github.com/studygolang/studygolang/model"
	"github.com/studygolang/studygolang/modules/context"

	echo "github.com/labstack/echo/v4"
	"github.com/studygolang/studygolang/modules/goutils"
)

type ProjectController struct{}

// 注册路由
func (self ProjectController) RegisterRoute(g *echo.Group) {
	g.GET("/projects", self.ReadList)
	g.GET("/project/detail", self.Detail)
}

// ReadList 开源项目列表页
func (ProjectController) ReadList(ctx echo.Context) error {
	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	paginator := logic.NewPaginatorWithPerPage(curPage, perPage)

	projects := logic.DefaultProject.FindAll(context.EchoContext(ctx), paginator, "id DESC", "")

	total := logic.DefaultProject.Count(context.EchoContext(ctx), "")
	hasMore := paginator.SetTotal(total).HasMorePage()

	data := map[string]interface{}{
		"projects": projects,
		"has_more": hasMore,
	}

	return success(ctx, data)
}

// Detail 项目详情
func (ProjectController) Detail(ctx echo.Context) error {
	id := goutils.MustInt(ctx.QueryParam("id"))
	project := logic.DefaultProject.FindOne(context.EchoContext(ctx), id)
	if project == nil || project.Id == 0 {
		return fail(ctx, "获取失败或已下线")
	}

	logic.Views.Incr(Request(ctx), model.TypeProject, project.Id)

	// 为了阅读数即时看到
	project.Viewnum++

	// 回复信息（评论）
	replies, _, lastReplyUser := logic.DefaultComment.FindObjComments(context.EchoContext(ctx), project.Id, model.TypeProject, 0, project.Lastreplyuid)
	// 有人回复
	if project.Lastreplyuid != 0 {
		project.LastReplyUser = lastReplyUser
	}

	return success(ctx, map[string]interface{}{
		"project": project,
		"replies": replies,
	})
}
