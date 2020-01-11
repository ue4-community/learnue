// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"github.com/studygolang/studygolang/model"

	"golang.org/x/net/context"

	. "github.com/studygolang/studygolang/db"

	"github.com/studygolang/studygolang/modules/logger"
)

type DynamicLogic struct{}

var DefaultDynamic = DynamicLogic{}

// FindBy 获取动态列表（分页）
func (DynamicLogic) FindBy(ctx context.Context, lastId int, limit int) []*model.Dynamic {
	dynamicList := make([]*model.Dynamic, 0)
	err := MasterDB.Where("id>?", lastId).OrderBy("seq DESC").Limit(limit).Find(&dynamicList)
	if err != nil {
		logger.Errorln("DynamicLogic FindBy Error:", err)
	}

	return dynamicList
}
