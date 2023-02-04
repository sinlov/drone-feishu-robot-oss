package feishu_robot_oss_plugin_test

import (
	"github.com/sinlov/drone-feishu-group-robot/feishu_plugin"
	"github.com/sinlov/drone-feishu-robot-oss/feishu_robot_oss_plugin"
	"github.com/sinlov/drone-info-tools/drone_info"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestPlugin(t *testing.T) {
	// mock Plugin
	t.Logf("~> mock Plugin")
	feishuCfg := feishu_plugin.Config{}
	p := feishu_robot_oss_plugin.Plugin{
		Name:    mockName,
		Version: mockVersion,
		Config: feishu_robot_oss_plugin.Config{
			FeishuCfg: feishuCfg,
		},
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
	if nil == err {
		t.Fatal("args [missing feishu webhook] error should be catch!")
	}
	t.Logf("~> mock FeishuPlugin")
	p.Config.FeishuCfg.Webhook = envFeishuWebHook
	p.Config.FeishuCfg.Secret = envFeishuSecret
	p.Config.FeishuCfg.FeishuEnableForward = false
	p.Config.FeishuCfg.CardOss.InfoSendResult = ""
	pagePasswd := mockOssPagePasswd

	p.Drone = *drone_info.MockDroneInfo(drone_info.DroneBuildStatusSuccess)
	checkCardOssRenderByPlugin(&p.Config.FeishuCfg, pagePasswd, true)
	// verify Plugin
	err = p.Exec()
	if err != nil {
		t.Fatalf("send failure error at %v", err)
	}
	p.Drone = *drone_info.MockDroneInfo(drone_info.DroneBuildStatusFailure)
	checkCardOssRenderByPlugin(&p.Config.FeishuCfg, pagePasswd, true)
	// verify Plugin
	err = p.Exec()
	if err != nil {
		t.Fatalf("send failure error at %v", err)
	}

	// open oss FeishuRobotOssTypeFileBrowser if open this must send file
	//p.Config.OssType = feishu_robot_oss_plugin.FeishuRobotOssTypeFileBrowser
	//p.Drone = *drone_info.MockDroneInfo("success")
	//p.Config.FeishuCfg.RenderOssCard = feishu_plugin.RenderStatusShow
	//p.Drone.Commit.Message = "build success but oss send failure and render RenderOssCard show"
	//p.Config.FeishuCfg.FeishuEnableForward = true
	//checkCardOssRenderByPlugin(&p.Config.FeishuCfg, pagePasswd, false)
	//// verify Plugin
	//err = p.Exec()
	//if err != nil {
	//	t.Fatalf("send failure error at %v", err)
	//}
	//p.Drone = *drone_info.MockDroneInfo("failure")
	//p.Drone.Commit.Message = "build failure and hide Oss settings and render OssStatus"
	//p.Config.FeishuCfg.RenderOssCard = feishu_plugin.RenderStatusShow
	//checkCardOssRenderByPlugin(&p.Config.FeishuCfg, pagePasswd, true)
	//// verify Plugin
	//err = p.Exec()
	//if err != nil {
	//	t.Fatalf("send failure error at %v", err)
	//}

	err = p.CleanResultEnv()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "", os.Getenv(feishu_robot_oss_plugin.EnvPluginResultShareHost))
}

func checkCardOssRenderByPlugin(cfg *feishu_plugin.Config, pagePasswd string, sendOssSucc bool) {
	cfg.CardOss.PagePasswd = pagePasswd
	if cfg.CardOss.PagePasswd == "" {
		cfg.CardOss.RenderResourceUrl = feishu_plugin.RenderStatusShow
	} else {
		cfg.CardOss.RenderResourceUrl = feishu_plugin.RenderStatusHide
	}
	if sendOssSucc {
		cfg.CardOss.InfoSendResult = feishu_plugin.RenderStatusShow
	} else {
		cfg.CardOss.InfoSendResult = feishu_plugin.RenderStatusHide
	}
	cfg.CardOss.Host = mockOssHost
	cfg.CardOss.InfoUser = mockOssUser
	cfg.CardOss.InfoPath = mockOssPath
	cfg.CardOss.ResourceUrl = mockOssResourceUrl
	cfg.CardOss.PageUrl = mockOssPageUrl
}
