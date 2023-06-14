package feishu_robot_oss_plugin

import (
	"fmt"
	"github.com/sinlov/drone-feishu-group-robot/feishu_plugin"
	"github.com/sinlov/drone-file-browser-plugin/file_browser_plugin"
	"github.com/sinlov/drone-info-tools/drone_info"
	"github.com/sinlov/drone-info-tools/template"
	tools "github.com/sinlov/drone-info-tools/tools/str_tools"
	"log"
	"os"
)

type (
	// Plugin plugin all config
	Plugin struct {
		Name    string
		Version string
		Drone   drone_info.Drone
		Config  Config
	}
)

func (p *Plugin) CleanResultEnv() error {
	for _, envItem := range cleanResultEnvList {
		err := os.Unsetenv(envItem)
		if err != nil {
			return fmt.Errorf("at FileBrowserPlugin.CleanResultEnv [ %s ], err: %v", envItem, err)
		}
	}
	return nil
}

func (p *Plugin) Exec() error {

	log.Printf("=> plugin %s version %s", p.Name, p.Version)

	if p.Config.Debug {
		for _, e := range os.Environ() {
			log.Println(e)
		}
	}

	var err error

	if !(tools.StrInArr(p.Config.OssType, supportOssType)) {
		return fmt.Errorf("-> feishu_robot_oss_type not support %s, can set %v", p.Config.OssType, supportOssType)
	}

	var ossPluginErr error
	switch p.Config.OssType {
	default:
		if p.Config.Debug {
			log.Printf("debug: now ossType is empty or not support %s\n", p.Config.OssType)
		}
	case FeishuRobotOssTypeFileBrowser:
		fileBrowserPlugin := file_browser_plugin.FileBrowserPlugin{
			Name:    p.Name,
			Version: p.Version,
			Drone:   p.Drone,
			Config:  p.Config.OssFileBrowserCfg,
		}
		ossPluginErr = fileBrowserPlugin.Exec()

		if ossPluginErr == nil {
			setEnvFromStr(*p, feishu_plugin.EnvPluginFeishuOssHost, lookupStrByEnv(file_browser_plugin.EnvPluginFileBrowserResultShareHost))
			setEnvFromStr(*p, feishu_plugin.EnvPluginFeishuOssInfoUser, lookupStrByEnv(file_browser_plugin.EnvPluginFileBrowserResultShareUser))
			setEnvFromStr(*p, feishu_plugin.EnvPluginFeishuOssInfoPath, lookupStrByEnv(file_browser_plugin.EnvPluginFileBrowserResultShareRemotePath))
			setEnvFromStr(*p, feishu_plugin.EnvPluginFeishuOssResourceUrl, lookupStrByEnv(file_browser_plugin.EnvPluginFileBrowserResultShareDownloadUrl))
			setEnvFromStr(*p, feishu_plugin.EnvPluginFeishuOssPageUrl, lookupStrByEnv(file_browser_plugin.EnvPluginFileBrowserResultSharePage))
			setEnvFromStr(*p, feishu_plugin.EnvPluginFeishuOssPagePasswd, lookupStrByEnv(file_browser_plugin.EnvPluginFileBrowserResultSharePasswd))
		}

		fileBrowserCleanResultEnvErr := fileBrowserPlugin.CleanResultEnv()
		if fileBrowserCleanResultEnvErr != nil {
			log.Fatalf("fileBrowserPlugin.CleanResultEnv() err: %v", fileBrowserCleanResultEnvErr)
		}
	}

	if ossPluginErr != nil {
		setEnvFromStr(*p, feishu_plugin.EnvPluginFeishuOssInfoSendResult, template.RenderStatusHide)
	} else {
		setEnvFromStr(*p, feishu_plugin.EnvPluginFeishuOssInfoSendResult, template.RenderStatusShow)
	}

	// cover by feishu env oss
	feishuCfg := p.Config.FeishuCfg
	ossHost := lookupStrCoverByEnv("", feishu_plugin.EnvPluginFeishuOssHost)
	cardOss := feishu_plugin.CardOss{
		Host: ossHost,
	}
	if ossHost == "" {
		feishuCfg.RenderOssCard = feishu_plugin.RenderStatusHide
	} else {
		feishuCfg.RenderOssCard = feishu_plugin.RenderStatusShow
		cardOss.InfoSendResult = lookupStrCoverByEnv(cardOss.InfoSendResult, feishu_plugin.EnvPluginFeishuOssInfoSendResult)
		cardOss.InfoUser = lookupStrCoverByEnv(cardOss.InfoUser, feishu_plugin.EnvPluginFeishuOssInfoUser)
		cardOss.InfoPath = lookupStrCoverByEnv(cardOss.InfoPath, feishu_plugin.EnvPluginFeishuOssInfoPath)
		cardOss.ResourceUrl = lookupStrCoverByEnv(cardOss.ResourceUrl, feishu_plugin.EnvPluginFeishuOssResourceUrl)
		cardOss.PageUrl = lookupStrCoverByEnv(cardOss.PageUrl, feishu_plugin.EnvPluginFeishuOssPageUrl)
		ossPagePasswd := lookupStrCoverByEnv("", feishu_plugin.EnvPluginFeishuOssPagePasswd)
		if ossPagePasswd == "" {
			cardOss.RenderResourceUrl = feishu_plugin.RenderStatusShow
		} else {
			cardOss.RenderResourceUrl = feishu_plugin.RenderStatusHide
			cardOss.PagePasswd = ossPagePasswd
		}
	}
	feishuCfg.CardOss = cardOss
	p.Config.FeishuCfg = feishuCfg

	feishuPlugin := feishu_plugin.FeishuPlugin{
		Name:    p.Name,
		Version: p.Version,
		Drone:   p.Drone,
		Config:  p.Config.FeishuCfg,
	}

	err = feishuPlugin.Exec()
	if err != nil {
		return err
	}
	if ossPluginErr != nil {
		log.Fatalf("ossPluginErr: %v", ossPluginErr)
		return ossPluginErr
	}

	log.Printf("=> plugin %s version %s", p.Name, p.Version)

	return err
}

func setEnvFromStr(p Plugin, key string, val string) {
	if p.Config.Debug {
		log.Printf("debug: setEnvFromStr key [ %s ] = %s", key, val)
	}
	err := os.Setenv(key, val)
	if err != nil {
		log.Fatalf("set env key [%v] string err: %v", key, err)
	}
}

func lookupStrByEnv(envKey string) string {
	envVal, lookupEnv := os.LookupEnv(envKey)
	if lookupEnv {
		return envVal
	}
	return ""
}

func lookupStrCoverByEnv(targetStr, envKey string) string {
	envVal, lookupEnv := os.LookupEnv(envKey)
	if lookupEnv {
		targetStr = envVal
	}
	return targetStr
}
