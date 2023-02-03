package feishu_robot_oss_plugin_test

import (
	"github.com/sinlov/drone-feishu-robot-oss/feishu_robot_oss_plugin"
	"github.com/sinlov/drone-info-tools/drone_info"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestPlugin(t *testing.T) {
	// mock Plugin
	t.Logf("~> mock Plugin")
	p := feishu_robot_oss_plugin.Plugin{
		Name:    mockName,
		Version: mockVersion,
	}
	// do Plugin
	t.Logf("~> do Plugin")
	if envCheck(t) {
		return
	}

	// use env:ENV_DEBUG
	p.Config.Debug = envDebug

	p.Config.OssType = mockOssTypeOther
	err := p.Exec()
	if nil == err {
		t.Fatal("args [ feishu_robot_oss_type ] error should be catch!")
	}

	// close oss
	p.Config.OssType = ""

	err = p.Exec()
	if nil != err {
		t.Fatal(err)
	}

	p.Drone = *drone_info.MockDroneInfo("success")
	// verify Plugin

	assert.Equal(t, "sinlov", p.Drone.Repo.OwnerName)

	err = p.CleanResultEnv()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "", os.Getenv(feishu_robot_oss_plugin.EnvPluginResultShareHost))
}
