// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package main

import (
	"github.com/polaris1119/keyword"
	"github.com/polaris1119/logger"
	"github.com/studygolang/studygolang/db"
)

func Indexer() {
	logger.Init(db.ROOT+"/log", db.ConfigFile.GetString("global.log_level"))
	go keyword.Extractor.Init(keyword.DefaultProps, true, db.ROOT+"/data/programming.txt,"+db.ROOT+"/data/dictionary.txt")

	IndexingServer()

	select {}
}
