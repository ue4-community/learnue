// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package db

import (
"database/sql"
"errors"
"fmt"
"github.com/spf13/viper"
"github.com/ue4-community/learnue/modules/setting"
"time"

_ "github.com/go-sql-driver/mysql"
"xorm.io/core"
"xorm.io/xorm"
)

var MasterDB *xorm.Engine

var dns string

func init() {
	mysqlConfig := setting.Get().Sub("mysql")
	if mysqlConfig == nil {
		fmt.Println("get mysql config error:")
		return
	}

	fillDns(mysqlConfig)

	// 启动时就打开数据库连接
	if err := initEngine(); err != nil {
		panic(err)
	}

	// 测试数据库连接是否 OK
	if err := MasterDB.Ping(); err != nil {
		panic(err)
	}
}

var (
	ConnectDBErr = errors.New("connect db error")
	UseDBErr     = errors.New("use db error")
)

// TestDB 测试数据库
func TestDB() error {
	mysqlConfig := setting.Get().Sub("mysql")
	if mysqlConfig == nil {
		fmt.Println("get mysql config error:")
		return errors.New("get mysql config error")
	}

	tmpDns := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=%s&parseTime=True&loc=Local",
		mysqlConfig.GetString("user"),
		mysqlConfig.GetString("password"),
		mysqlConfig.GetString("host"),
		mysqlConfig.GetString("port"),
		mysqlConfig.GetString("charset"))
	egnine, err := xorm.NewEngine("mysql", tmpDns)
	if err != nil {
		fmt.Println("new engine error:", err)
		return err
	}
	defer egnine.Close()

	// 测试数据库连接是否 OK
	if err = egnine.Ping(); err != nil {
		fmt.Println("ping db error:", err)
		return ConnectDBErr
	}

	_, err = egnine.Exec("use " + mysqlConfig.GetString("dbname"))
	if err != nil {
		fmt.Println("use db error:", err)
		_, err = egnine.Exec("CREATE DATABASE " + mysqlConfig.GetString("dbname") + " DEFAULT CHARACTER SET " + mysqlConfig.GetString("charset"))
		if err != nil {
			fmt.Println("create database error:", err)

			return UseDBErr
		}

		fmt.Println("create database successfully!")
	}

	// 初始化 MasterDB
	return Init()
}

func Init() error {
	mysqlConfig := setting.Get().Sub("mysql")
	if mysqlConfig == nil {
		fmt.Println("get mysql config error")
		return errors.New("get mysql config error")
	}

	fillDns(mysqlConfig)

	// 启动时就打开数据库连接
	if err := initEngine(); err != nil {
		fmt.Println("mysql is not open:", err)
		return err
	}

	return nil
}

func fillDns(mysqlConfig *viper.Viper) {
	dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		mysqlConfig.GetString("user"),
		mysqlConfig.GetString("password"),
		mysqlConfig.GetString("host"),
		mysqlConfig.GetString("port"),
		mysqlConfig.GetString("dbname"),
		mysqlConfig.GetString("charset"))
}

func initEngine() error {
	var err error

	MasterDB, err = xorm.NewEngine("mysql", dns)
	if err != nil {
		return err
	}

	maxIdle := setting.Get().GetInt("mysql.max_idle")
	maxConn := setting.Get().GetInt("mysql.max_conn")

	MasterDB.SetMaxIdleConns(maxIdle)
	MasterDB.SetMaxOpenConns(maxConn)

	showSQL := setting.Get().GetBool("xorm.show_sql")
	logLevel := setting.Get().GetInt("xorm.log_level")

	MasterDB.ShowSQL(showSQL)
	MasterDB.Logger().SetLevel(core.LogLevel(logLevel))
	MasterDB.TZLocation, _ = time.LoadLocation("Asia/Shanghai")

	// 启用缓存
	// cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
	// MasterDB.SetDefaultCacher(cacher)

	return nil
}

func StdMasterDB() *sql.DB {
	return MasterDB.DB().DB
}
