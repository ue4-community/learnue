package logic_test

import (
	. "github.com/studygolang/studygolang/config"
	"github.com/polaris1119/logger"

	"github.com/studygolang/studygolang/logic"
	"testing"
)

func TestSendAuthMail(t *testing.T) {
	logger.Init(ROOT+"/log", ConfigFile.GetString("global.log_level"))

	err := logic.DefaultEmail.SendAuthMail("中文test", "内容test content，收到？", []string{"xuxinhua@zhimadj.com"})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("successful")
	}
}
