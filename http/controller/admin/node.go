// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package admin

import (
	"github.com/studygolang/studygolang/global"
	"github.com/studygolang/studygolang/logic"
	"github.com/studygolang/studygolang/model"
	"github.com/studygolang/studygolang/modules/context"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

type NodeController struct{}

// 注册路由
func (self NodeController) RegisterRoute(g *echo.Group) {
	g.GET("/community/node/list", self.List)
	g.Match([]string{"GET", "POST"}, "/community/node/modify", self.Modify)
	g.POST("/community/node/modify_seq", self.ModifySeq)
}

// List 所有主题节点
func (NodeController) List(ctx echo.Context) error {
	treeNodes := logic.DefaultNode.FindParallelTree(context.EchoContext(ctx))

	nidMap := make(map[int]int)
	keySlice := make([]int, len(treeNodes))

	for i, node := range treeNodes {
		nidMap[node.Nid] = i + 1

		if node.Parent > 0 {
			keySlice[i] = nidMap[node.Parent]
		} else {
			keySlice[i] = 0
		}
	}

	data := map[string]interface{}{
		"nodes":     treeNodes,
		"key_slice": keySlice,
	}

	return render(ctx, "topic/node.html", data)
}

func (NodeController) Modify(ctx echo.Context) error {
	if ctx.FormValue("submit") == "1" {
		forms, _ := ctx.FormParams()
		err := logic.DefaultNode.Modify(context.EchoContext(ctx), forms)
		if err != nil {
			return fail(ctx, 1, err.Error())
		}
		global.TopicNodeChan <- struct{}{}
		return success(ctx, nil)
	}

	treeNodes := logic.DefaultNode.FindParallelTree(context.EchoContext(ctx))

	data := map[string]interface{}{
		"nodes": treeNodes,
	}

	nid := goutils.MustInt(ctx.QueryParam("nid"))
	parent := goutils.MustInt(ctx.QueryParam("parent"))
	if nid == 0 && parent == 0 {
		// 新增
		data["node"] = &model.TopicNode{
			ShowIndex: true,
		}
	} else if nid > 0 {
		data["node"] = logic.DefaultNode.FindOne(nid)
	} else if parent > 0 {
		data["node"] = &model.TopicNode{
			ShowIndex: true,
		}
	}
	data["parent"] = parent

	return render(ctx, "topic/node_modify.html", data)
}

func (NodeController) ModifySeq(ctx echo.Context) error {
	nid := goutils.MustInt(ctx.FormValue("nid"))
	seq := goutils.MustInt(ctx.FormValue("seq"))
	err := logic.DefaultNode.ModifySeq(context.EchoContext(ctx), nid, seq)
	if err != nil {
		return fail(ctx, 1, err.Error())
	}
	return success(ctx, nil)

}
