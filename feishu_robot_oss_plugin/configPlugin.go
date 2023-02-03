package feishu_robot_oss_plugin

import (
	"github.com/sinlov/drone-feishu-group-robot/feishu_plugin"
	"github.com/sinlov/drone-file-browser-plugin/file_browser_plugin"
)

const (
	EnvPluginResultShareHost = "PLUGIN_RESULT_SHARE_HOST"

	FeishuRobotOssTypeFileBrowser = "filebrowser"
)

var (
	// supportMsgType
	supportOssType = []string{
		"",
		FeishuRobotOssTypeFileBrowser,
	}

	cleanResultEnvList = []string{
		EnvPluginResultShareHost,
	}
)

type (

	// Config plugin private config
	Config struct {
		Debug         bool
		TimeoutSecond uint

		FeishuCfg feishu_plugin.Config

		// OssType
		// just use var:supportOssType
		OssType string

		OssFileBrowserCfg file_browser_plugin.Config
	}
)
