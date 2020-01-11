// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package app

import (
	"encoding/json"
	"net/http"

	. "github.com/studygolang/studygolang/http"
	"github.com/studygolang/studygolang/logic"
	"github.com/studygolang/studygolang/modules/context"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/logger"
	"github.com/polaris1119/nosql"
)

const perPage = 12

func getLogger(ctx echo.Context) *logger.Logger {
	return logic.GetLogger(context.EchoContext(ctx))
}

func success(ctx echo.Context, data interface{}) error {
	result := map[string]interface{}{
		"code": 0,
		"msg":  "ok",
		"data": data,
	}

	b, err := json.Marshal(result)
	if err != nil {
		return err
	}

	go func(b []byte) {
		if cacheKey := ctx.Get(nosql.CacheKey); cacheKey != nil {
			nosql.DefaultLRUCache.CompressAndAdd(cacheKey, b, nosql.NewCacheData())
		}
	}(b)

	AccessControl(ctx)

	if ctx.Response().Committed {
		getLogger(ctx).Flush()
		return nil
	}

	return ctx.JSONBlob(http.StatusOK, b)
}

func fail(ctx echo.Context, msg string, codes ...int) error {
	AccessControl(ctx)

	if ctx.Response().Committed {
		getLogger(ctx).Flush()
		return nil
	}

	code := 1
	if len(codes) > 0 {
		code = codes[0]
	}
	result := map[string]interface{}{
		"code": code,
		"msg":  msg,
	}

	getLogger(ctx).Errorln("operate fail:", result)

	return ctx.JSON(http.StatusOK, result)
}
