// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

// 可选择是否在启动主程序时，同时嵌入 indexer 和 crawler，减少内存占用

package main

import (
	"flag"
	"github.com/ue4-community/learnue/modules/setting"
	"time"

	"github.com/ue4-community/learnue/modules/logger"
	"github.com/robfig/cron"
	"github.com/ue4-community/learnue/logic"
)

var (
	manualIndex = flag.Bool("manual", false, "do manual index once or not")
	needAll     = flag.Bool("all", false, "是否需要全量抓取，默认否")
	whichSite   = flag.String("site", "", "抓取哪个站点（空表示所有站点）")
)

func IndexingServer() {


	if *manualIndex {
		indexing(true)
	}

	c := cron.New()
	// 构建 solr 需要的索引数据
	// 1 分钟一次增量
	c.AddFunc("@every 1m", func() {
		indexing(false)
	})
	// 一周一次全量（周六晚上2点开始）
	c.AddFunc("0 0 2 * * 6", func() {
		indexing(true)
	})

	c.Start()
}

func indexing(isAll bool) {
	logger.Infoln("indexing start...")

	start := time.Now()
	defer func() {
		logger.Infoln("indexing spend time:", time.Now().Sub(start))
	}()

	logic.DefaultSearcher.Indexing(isAll)
}

func CrawlServer() {
	if !flag.Parsed() {
		flag.Parse()
	}

	go autocrawl(*needAll, *whichSite)
}

func autocrawl(needAll bool, whichSite string) {
	if needAll {
		if whichSite != "" {
			go logic.DefaultAutoCrawl.CrawlWebsite(whichSite, needAll)
		} else {
			go logic.DefaultAutoCrawl.DoCrawl(needAll)
		}
	}

	// 定时增量
	c := cron.New()
	c.AddFunc(setting.Get().GetString("crawl.spec"), func() {
		// 抓取 reddit
		go logic.DefaultReddit.Parse("")

		projectUrl := setting.Get().GetString("crawl.project_url")
		if projectUrl != "" {
			// 抓取 project
			go logic.DefaultProject.ParseProjectList(projectUrl)
		}

		// 抓取 article
		go logic.DefaultAutoCrawl.DoCrawl(false)
	})
	c.Start()
}
