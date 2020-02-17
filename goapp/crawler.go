// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package main

import (
	"github.com/ue4-community/learnue/modules/keyword"
	"github.com/ue4-community/learnue/modules/logger"
	"github.com/ue4-community/learnue/modules/setting"
)

func Crawler() {
	logger.Init(setting.ROOT+"/log", setting.Get().GetString("global.log_level"), "crawl")
	go keyword.Extractor.Init(keyword.DefaultProps, true, setting.ROOT+"/data/programming.txt,"+setting.ROOT+"/data/dictionary.txt")

	CrawlServer()

	select {}
}
