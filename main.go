package main

import (
	"fmt"
	"github.com/sinlov/drone-info-tools/drone_info"
	"log"
	"os"
	"time"

	"github.com/sinlov/drone-feishu-group-robot/feishu_plugin"
	"github.com/sinlov/drone-feishu-robot-oss/feishu_robot_oss_plugin"
	"github.com/sinlov/drone-file-browser-plugin/file_browser_plugin"
	"github.com/sinlov/drone-info-tools/drone_urfave_cli_v2"
	"github.com/sinlov/filebrowser-client/web_api"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

const (
	// Version of cli
	Version = "v1.0.0"
	Name    = "drone-feishu-robot-oss"
)

// action
// do cli Action before flag.
func action(c *cli.Context) error {

	isDebug := c.Bool("config.debug")

	drone := drone_urfave_cli_v2.UrfaveCliBindDroneInfo(c)

	if isDebug {
		log.Printf("debug: cli version is %s", Version)
		log.Printf("debug: load droneInfo finish at link: %v\n", drone.Build.Link)
	}

	ossType := c.String("config.feishu_robot_oss_type")
	config := feishu_robot_oss_plugin.Config{
		Debug:         c.Bool("config.debug"),
		TimeoutSecond: c.Uint("config.timeout_second"),
		OssType:       ossType,
	}

	switch ossType {
	default:
		if isDebug {
			log.Printf("debug: now ossType is empty or not support %s\n", ossType)
		}
	case feishu_robot_oss_plugin.FeishuRobotOssTypeFileBrowser:
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

	p := feishu_robot_oss_plugin.Plugin{
		Name:    Name,
		Version: Version,
		Drone:   drone,
		Config:  config,
	}
	err := p.Exec()

	if err != nil {
		log.Fatalf("err: %v", err)
		return err
	}

	return nil
}

func findStrFromCliOrCoverByEnv(c *cli.Context, ctxKey, envKey string) string {
	val := c.String(ctxKey)
	envVal, lookupEnv := os.LookupEnv(envKey)
	if lookupEnv {
		val = envVal
	}
	return val
}

// pluginOSSFileBrowser
// set plugin flag at here
func pluginOSSFileBrowser() []cli.Flag {
	return []cli.Flag{
		// file_browser_plugin start
		&cli.StringFlag{
			Name:    "config.file_browser_host,file_browser_host",
			Usage:   "must set args, file_browser host",
			EnvVars: []string{"PLUGIN_FILE_BROWSER_HOST"},
		},
		&cli.StringFlag{
			Name:    "config.file_browser_username,file_browser_username",
			Usage:   "must set args, file_browser username",
			EnvVars: []string{"PLUGIN_FILE_BROWSER_USERNAME"},
		},
		&cli.StringFlag{
			Name:    "config.file_browser_user_password,file_browser_user_password",
			Usage:   "must set args, file_browser user password",
			EnvVars: []string{"PLUGIN_FILE_BROWSER_USER_PASSWORD"},
		},
		&cli.StringFlag{
			Name:    "config.file_browser_work_space,file_browser_work_space",
			Usage:   fmt.Sprintf("file_browser work space. default will use env:%s", drone_info.EnvDroneBuildWorkSpace),
			EnvVars: []string{"PLUGIN_FILE_BROWSER_WORK_SPACE"},
		},
		&cli.UintFlag{
			Name:    "config.file_browser_timeout_push_second,file_browser_timeout_push_second",
			Usage:   "file_browser push each file timeout push second, must gather than 60",
			Value:   60,
			EnvVars: []string{"PLUGIN_FILE_BROWSER_TIMEOUT_PUSH_SECOND"},
		},
		&cli.StringFlag{
			Name:    "config.file_browser_work_mode,file_browser_work_mode",
			Usage:   "must set args, work mode only can use: send, download",
			Value:   file_browser_plugin.WorkModeSend,
			EnvVars: []string{"PLUGIN_FILE_BROWSER_WORK_MODE"},
		},
		&cli.StringFlag{
			Name:    "config.file_browser_dist_type,file_browser_dist_type",
			Usage:   "must set args, type of dist file graph only can use: git, custom",
			Value:   file_browser_plugin.DistTypeGit,
			EnvVars: []string{"PLUGIN_FILE_BROWSER_DIST_TYPE"},
		},
		&cli.StringFlag{
			Name:    "config.file_browser_dist_graph,file_browser_dist_graph",
			Usage:   "type of dist custom set as struct [ drone_info.Drone ]",
			Value:   "{{ Repo.HostName }}/{{ Repo.GroupName }}/{{ Repo.ShortName }}/s/{{ Build.Number }}/{{ Stage.Name }}-{{ Build.Number }}-{{ Stage.FinishedTime }}",
			EnvVars: []string{"PLUGIN_FILE_BROWSER_DIST_GRAPH"},
		},
		&cli.StringFlag{
			Name:    "config.file_browser_remote_root_path,file_browser_remote_root_path",
			Usage:   "must set args, this will append by file_browser_dist_type at remote",
			EnvVars: []string{"PLUGIN_FILE_BROWSER_REMOTE_ROOT_PATH"},
		},
		&cli.StringFlag{
			Name:    "config.file_browser_target_dist_root_path,file_browser_target_dist_root_path",
			Usage:   "path of file_browser local work on root, can set \"\"",
			Value:   "",
			EnvVars: []string{"PLUGIN_FILE_BROWSER_TARGET_DIST_ROOT_PATH"},
		},
		&cli.StringSliceFlag{
			Name:    "config.file_browser_target_file_globs,file_browser_target_file_globs",
			Usage:   "must set args, globs list of send to file_browser under file_browser_target_dist_root_path",
			EnvVars: []string{"PLUGIN_FILE_BROWSER_TARGET_FILE_GLOBS"},
		},
		&cli.StringFlag{
			Name:    "config.file_browser_target_file_regular,file_browser_target_file_regular",
			Usage:   "must set args, regular of send to file_browser under file_browser_target_dist_root_path",
			EnvVars: []string{"PLUGIN_FILE_BROWSER_TARGET_FILE_REGULAR"},
		},
		&cli.BoolFlag{
			Name:    "config.file_browser_share_link_enable,file_browser_share_link_enable",
			Usage:   "share dist dir as link",
			Value:   true,
			EnvVars: []string{"PLUGIN_FILE_BROWSER_SHARE_LINK_ENABLE"},
		},
		&cli.StringFlag{
			Name:    "config.file_browser_share_link_unit,file_browser_share_link_unit",
			Usage:   "take effect by open share_link, only can use as [ days hours minutes seconds ]",
			Value:   web_api.ShareUnitDays,
			EnvVars: []string{"PLUGIN_FILE_BROWSER_SHARE_LINK_UNIT"},
		},
		&cli.UintFlag{
			Name:    "config.file_browser_share_link_expires,file_browser_share_link_expires",
			Usage:   "if set 0, will allow share_link exist forever, default: 0",
			Value:   0,
			EnvVars: []string{"PLUGIN_FILE_BROWSER_SHARE_LINK_EXPIRES"},
		},
		&cli.BoolFlag{
			Name:    "config.file_browser_share_link_auto_password_enable,file_browser_share_link_auto_password_enable",
			Usage:   "password of share_link auto , if open this will cover settings.file_browser_share_link_password",
			Value:   false,
			EnvVars: []string{"PLUGIN_FILE_BROWSER_SHARE_LINK_AUTO_PASSWORD_ENABLE"},
		},
		&cli.StringFlag{
			Name:    "config.file_browser_share_link_password,file_browser_share_link_password",
			Usage:   "password of share_link, if not set will not use password, default: \"\"",
			Value:   "",
			EnvVars: []string{"PLUGIN_FILE_BROWSER_SHARE_LINK_PASSWORD"},
		},
	}
}

// pluginFeishu
// set plugin hide flag at here
func pluginFeishu() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "config.ntp_target,ntp_target",
			Usage:   "ntp target like: pool.ntp.org, time1.google.com,time.pool.aliyun.com, default not use ntpd to sync",
			EnvVars: []string{"PLUGIN_NTP_TARGET"},
		},
		&cli.StringFlag{
			Name:    "config.webhook,feishu_webhook",
			Usage:   "feishu webhook for send message",
			EnvVars: []string{feishu_plugin.EnvPluginFeishuWebhook},
		},
		&cli.StringFlag{
			Name:    "config.secret,feishu_secret",
			Usage:   "feishu secret",
			EnvVars: []string{feishu_plugin.EnvPluginFeishuSecret},
		},
		&cli.BoolFlag{
			Name:    "config.feishu_enable_forward,feishu_enable_forward",
			Usage:   "feishu message enable forward, default false",
			EnvVars: []string{feishu_plugin.EnvPluginFeishuEnableForward},
		},
		&cli.StringFlag{
			Name:    "config.msg_type,feishu_msg_type",
			Usage:   "feishu message type",
			EnvVars: []string{feishu_plugin.EnvPluginFeishuMsgType},
		},
		&cli.StringFlag{
			Name:    "config.msg_title,feishu_msg_title",
			Usage:   "feishu message title",
			EnvVars: []string{feishu_plugin.EnvPluginFeishuMsgTitle},
		},
		&cli.StringFlag{
			Name:    "config.msg_powered_by_image_key,feishu_msg_powered_by_image_key",
			Usage:   "feishu message powered by image key",
			EnvVars: []string{feishu_plugin.EnvPluginFeishuMsgPoweredByImageKey},
		},
		&cli.StringFlag{
			Name:    "config.msg_powered_by_image_alt,feishu_msg_powered_by_image_alt",
			Usage:   "feishu message powered by image alt",
			EnvVars: []string{feishu_plugin.EnvPluginFeishuMsgPoweredByImageAlt},
		},

		// oss card end
		&cli.StringFlag{
			Name:    "config.feishu_oss_host",
			Usage:   "feishu OSS host for show oss info, if empty will not show oss info",
			EnvVars: []string{feishu_plugin.EnvPluginFeishuOssHost},
		},
		&cli.StringFlag{
			Name:    "config.feishu_oss_info_user",
			Usage:   "feishu OSS user for show at card",
			EnvVars: []string{feishu_plugin.EnvPluginFeishuOssInfoUser},
		},
		&cli.StringFlag{
			Name:    "config.feishu_oss_info_path",
			Usage:   "feishu OSS path for show at card",
			EnvVars: []string{feishu_plugin.EnvPluginFeishuOssInfoPath},
		},
		&cli.StringFlag{
			Name:    "config.feishu_oss_resource_url",
			Usage:   "feishu OSS resource url",
			EnvVars: []string{feishu_plugin.EnvPluginFeishuOssResourceUrl},
		},
		&cli.StringFlag{
			Name:    "config.feishu_oss_page_url",
			Usage:   "feishu OSS page url",
			EnvVars: []string{feishu_plugin.EnvPluginFeishuOssPageUrl},
		},
		&cli.StringFlag{
			Name:    "config.feishu_oss_page_passwd",
			Usage:   "OSS password at page url, will hide PLUGIN_FEISHU_OSS_RESOURCE_URL when PAGE_PASSWD not empty",
			EnvVars: []string{feishu_plugin.EnvPluginFeishuOssPagePasswd},
		},
		// oss card end
	}
}

// pluginCommon
// set plugin common flag at here
func pluginCommon() []cli.Flag {
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

func main() {
	app := cli.NewApp()
	app.Version = Version
	app.Name = "Drone Plugin"
	app.Usage = ""
	year := time.Now().Year()
	app.Copyright = fmt.Sprintf("© 2022-%d sinlov", year)
	author := &cli.Author{
		Name:  "sinlov",
		Email: "sinlovgmppt@gmail.com",
	}
	app.Authors = []*cli.Author{
		author,
	}

	app.Action = action
	flags := drone_urfave_cli_v2.UrfaveCliAppendCliFlag(drone_urfave_cli_v2.DroneInfoUrfaveCliFlag(), pluginCommon())
	flags = drone_urfave_cli_v2.UrfaveCliAppendCliFlag(flags, pluginFeishu())
	flags = drone_urfave_cli_v2.UrfaveCliAppendCliFlag(flags, pluginOSSFileBrowser())
	app.Flags = flags

	// kubernetes runner patch
	if _, err := os.Stat("/run/drone/env"); err == nil {
		errDotEnv := godotenv.Overload("/run/drone/env")
		if errDotEnv != nil {
			log.Fatalf("load /run/drone/env err: %v", errDotEnv)
		}
	}

	// app run as urfave
	if err := app.Run(os.Args); nil != err {
		log.Println(err)
	}
}
