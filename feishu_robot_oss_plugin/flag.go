package feishu_robot_oss_plugin

import (
	"github.com/sinlov/drone-feishu-group-robot/feishu_plugin"
	"github.com/sinlov/drone-file-browser-plugin/file_browser_plugin"
	"github.com/sinlov/drone-info-tools/drone_info"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func BindFlag(c *cli.Context, isDebug bool, cliVersion, cliName string, drone drone_info.Drone) Plugin {
	ossType := c.String("config.feishu_robot_oss_type")
	config := Config{
		Debug:         c.Bool("config.debug"),
		TimeoutSecond: c.Uint("config.timeout_second"),
		OssType:       ossType,
	}

	switch ossType {
	default:
		if isDebug {
			log.Printf("debug: now ossType is empty or not support %s\n", ossType)
		}
	case FeishuRobotOssTypeFileBrowser:
		// append filebrowser oss config
		fileBrowserPlugin := file_browser_plugin.BindFlag(c, cliVersion, cliName, drone)
		config.OssFileBrowserCfg = fileBrowserPlugin.Config
	}

	feishuPlugin := feishu_plugin.BindFlag(c, cliVersion, cliName, drone)

	config.FeishuCfg = feishuPlugin.Config

	if isDebug {
		log.Printf("config.timeout_second: %v", config.TimeoutSecond)
	}

	p := Plugin{
		Name:    cliName,
		Version: cliVersion,
		Drone:   drone,
		Config:  config,
	}
	return p
}

// Flag
// set plugin flag at here
func Flag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "config.feishu_robot_oss_type,feishu_robot_oss_type",
			Usage:   "choose oss type, if type is \"\" or not set, will use feishu robot send no oss message",
			EnvVars: []string{EnvPluginFeishuRobotOssType},
		},
	}
}

// CommonFlag
// set plugin common flag at here
func CommonFlag() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:    "config.debug,debug",
			Usage:   "debug mode",
			EnvVars: []string{"PLUGIN_DEBUG"},
		},
		&cli.IntFlag{
			Name:    "config.timeout_second,timeout_second",
			Usage:   "do request timeout setting second",
			EnvVars: []string{"PLUGIN_TIMEOUT_SECOND"},
		},
	}
}

func findStrFromCliOrCoverByEnv(c *cli.Context, ctxKey, envKey string) string {
	val := c.String(ctxKey)
	envVal, lookupEnv := os.LookupEnv(envKey)
	if lookupEnv {
		val = envVal
	}
	return val
}
