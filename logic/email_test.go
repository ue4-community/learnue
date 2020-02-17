package logic_test

import (
	"github.com/ue4-community/learnue/modules/logger"
	"github.com/ue4-community/learnue/db"

	"github.com/ue4-community/learnue/logic"
	"testing"
)

func TestSendAuthMail(t *testing.T) {
	logger.Init(setting.ROOT+"/log", setting.Get().GetString("global.log_level"))

	err := logic.DefaultEmail.SendAuthMail("中文test", "内容test content，收到？", []string{"xuxinhua@zhimadj.com"})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("successful")
	}
}
