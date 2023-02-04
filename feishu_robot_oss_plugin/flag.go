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
		fileBrowserConfig := file_browser_plugin.Config{

			Debug:         c.Bool("config.debug"),
			TimeoutSecond: c.Uint("config.timeout_second"),

			FileBrowserBaseConfig: file_browser_plugin.FileBrowserBaseConfig{
				FileBrowserHost:              c.String("config.file_browser_host"),
				FileBrowserUsername:          c.String("config.file_browser_username"),
				FileBrowserUserPassword:      c.String("config.file_browser_user_password"),
				FileBrowserTimeoutPushSecond: c.Uint("config.file_browser_timeout_push_second"),
				FileBrowserWorkSpace:         c.String("config.file_browser_work_space"),
			},

			FileBrowserWorkMode: c.String("config.file_browser_work_mode"),

			FileBrowserSendModeConfig: file_browser_plugin.FileBrowserSendModeConfig{
				FileBrowserDistType:           c.String("config.file_browser_dist_type"),
				FileBrowserDistGraph:          c.String("config.file_browser_dist_graph"),
				FileBrowserRemoteRootPath:     c.String("config.file_browser_remote_root_path"),
				FileBrowserTargetDistRootPath: c.String("config.file_browser_target_dist_root_path"),
				FileBrowserTargetFileGlob:     c.StringSlice("config.file_browser_target_file_globs"),
				FileBrowserTargetFileRegular:  c.String("config.file_browser_target_file_regular"),

				FileBrowserShareLinkEnable:             c.Bool("config.file_browser_share_link_enable"),
				FileBrowserShareLinkUnit:               c.String("config.file_browser_share_link_unit"),
				FileBrowserShareLinkExpires:            c.Uint("config.file_browser_share_link_expires"),
				FileBrowserShareLinkAutoPasswordEnable: c.Bool("config.file_browser_share_link_auto_password_enable"),
				FileBrowserShareLinkPassword:           c.String("config.file_browser_share_link_password"),
			},

			FileBrowserDownloadModeConfig: file_browser_plugin.FileBrowserDownloadModeConfig{
				FileBrowserDownloadEnable:    c.Bool("config.file_browser_download_enable"),
				FileBrowserDownloadPath:      c.String("config.file_browser_download_remote_path"),
				FileBrowserDownloadLocalPath: c.String("config.file_browser_download_local_path"),
			},
		}
		config.OssFileBrowserCfg = fileBrowserConfig
	}

	feishuCfg := feishu_plugin.Config{
		Debug:               c.Bool("config.debug"),
		TimeoutSecond:       c.Int("config.timeout_second"),
		NtpTarget:           c.String("config.ntp_target"),
		Webhook:             c.String("config.webhook"),
		Secret:              c.String("config.secret"),
		FeishuEnableForward: c.Bool("config.feishu_enable_forward"),
		MsgType:             c.String("config.msg_type"),
		Title:               c.String("config.msg_title"),
		PoweredByImageKey:   c.String("config.msg_powered_by_image_key"),
		PoweredByImageAlt:   c.String("config.msg_powered_by_image_alt"),
	}

	ossHost := findStrFromCliOrCoverByEnv(c, "config.feishu_oss_host", feishu_plugin.EnvPluginFeishuOssHost)
	cardOss := feishu_plugin.CardOss{}
	if ossHost == "" {
		feishuCfg.RenderOssCard = feishu_plugin.RenderStatusHide
	} else {
		feishuCfg.RenderOssCard = feishu_plugin.RenderStatusShow
		cardOss.InfoSendResult = findStrFromCliOrCoverByEnv(c, "config.feishu_oss_info_send_result", feishu_plugin.EnvPluginFeishuOssInfoSendResult)
		cardOss.InfoUser = findStrFromCliOrCoverByEnv(c, "config.feishu_oss_info_user", feishu_plugin.EnvPluginFeishuOssInfoUser)
		cardOss.InfoPath = findStrFromCliOrCoverByEnv(c, "config.feishu_oss_info_path", feishu_plugin.EnvPluginFeishuOssInfoPath)
		cardOss.ResourceUrl = findStrFromCliOrCoverByEnv(c, "config.feishu_oss_resource_url", feishu_plugin.EnvPluginFeishuOssResourceUrl)
		cardOss.PageUrl = findStrFromCliOrCoverByEnv(c, "config.feishu_oss_page_url", feishu_plugin.EnvPluginFeishuOssPageUrl)
		ossPagePasswd := findStrFromCliOrCoverByEnv(c, "config.feishu_oss_page_passwd", feishu_plugin.EnvPluginFeishuOssPagePasswd)
		if ossPagePasswd == "" {
			cardOss.RenderResourceUrl = feishu_plugin.RenderStatusShow
		} else {
			cardOss.RenderResourceUrl = feishu_plugin.RenderStatusHide
			cardOss.PagePasswd = ossPagePasswd
		}
	}
	feishuCfg.CardOss = cardOss

	config.FeishuCfg = feishuCfg

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
		&cli.StringFlag{
			Name:    "config.feishu_robot_oss_type,feishu_robot_oss_type",
			Usage:   "choose oss type, if type is \"\" or not set, will use feishu robot send no oss message",
			EnvVars: []string{"PLUGIN_FEISHU_ROBOT_OSS_TYPE"},
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
