package feishu_robot_oss_plugin

import (
	"github.com/sinlov/drone-feishu-group-robot/feishu_plugin"
	"github.com/sinlov/drone-file-browser-plugin/file_browser_plugin"
	"github.com/sinlov/drone-info-tools/drone_info"
	"github.com/sinlov/drone-info-tools/drone_log"
	"github.com/urfave/cli/v2"
	"os"
)

// IsBuildDebugOpen
// when config or drone build open debug will open debug
func IsBuildDebugOpen(c *cli.Context) bool {
	return c.Bool(NamePluginDebug) || c.Bool(drone_info.NameCliStepsDebug)
}

// BindCliFlag
// check args here
func BindCliFlag(c *cli.Context, cliVersion, cliName string, drone drone_info.Drone) (*Plugin, error) {
	debug := IsBuildDebugOpen(c)
	p := BindFlag(c, debug, cliVersion, cliName, drone)

	return &p, nil
}

func BindFlag(c *cli.Context, isDebug bool, cliVersion, cliName string, drone drone_info.Drone) Plugin {
	ossType := c.String("config.feishu_robot_oss_type")
	config := Config{
		Debug:         isDebug,
		TimeoutSecond: c.Uint(NamePluginTimeOut),
		OssType:       ossType,
	}
	drone_log.Debugf("args %s: %v", NamePluginTimeOut, config.TimeoutSecond)

	switch ossType {
	default:
		if isDebug {
			drone_log.Warnf("debug: now ossType is empty or not support %s\n", ossType)
		}
	case FeishuRobotOssTypeFileBrowser:
		// append filebrowser oss config
		fileBrowserPlugin := file_browser_plugin.BindFlag(c, cliVersion, cliName, drone)
		config.OssFileBrowserCfg = fileBrowserPlugin.Config
	}

	feishuPlugin := feishu_plugin.BindFlag(c, cliVersion, cliName, drone)

	config.FeishuCfg = feishuPlugin.Config

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
		&cli.UintFlag{
			Name:    NamePluginTimeOut,
			Usage:   "do request timeout setting second.",
			Hidden:  true,
			Value:   10,
			EnvVars: []string{EnvPluginTimeOut},
		},
		&cli.BoolFlag{
			Name:    NamePluginDebug,
			Usage:   "debug mode",
			Value:   false,
			EnvVars: []string{drone_info.EnvKeyPluginDebug},
		},
	}
}

//nolint:golint,unused
func findStrFromCliOrCoverByEnv(c *cli.Context, ctxKey, envKey string) string {
	val := c.String(ctxKey)
	envVal, lookupEnv := os.LookupEnv(envKey)
	if lookupEnv {
		val = envVal
	}
	return val
}
