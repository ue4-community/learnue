package logic_test

import (
	"github.com/polaris1119/logger"
	"github.com/studygolang/studygolang/db"

	"github.com/studygolang/studygolang/logic"
	"testing"
)

func TestSendAuthMail(t *testing.T) {
	logger.Init(db.ROOT+"/log", db.ConfigFile.GetString("global.log_level"))

	err := logic.DefaultEmail.SendAuthMail("中文test", "内容test content，收到？", []string{"xuxinhua@zhimadj.com"})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("successful")
	}
}
