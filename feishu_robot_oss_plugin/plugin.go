package feishu_robot_oss_plugin

import (
	"fmt"
	"github.com/sinlov/drone-feishu-group-robot/feishu_plugin"
	"github.com/sinlov/drone-file-browser-plugin/file_browser_plugin"
	"github.com/sinlov/drone-info-tools/drone_info"
	"github.com/sinlov/drone-info-tools/template"
	tools "github.com/sinlov/drone-info-tools/tools/str_tools"
	"log"
	"math/rand"
	"os"
	"time"
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
			setEnvFromStr(feishu_plugin.EnvPluginFeishuOssHost, lookupStrByEnv(file_browser_plugin.EnvPluginFileBrowserResultShareHost))
			setEnvFromStr(feishu_plugin.EnvPluginFeishuOssInfoUser, lookupStrByEnv(file_browser_plugin.EnvPluginFileBrowserResultShareUser))
			setEnvFromStr(feishu_plugin.EnvPluginFeishuOssInfoPath, lookupStrByEnv(file_browser_plugin.EnvPluginFileBrowserResultShareRemotePath))
			setEnvFromStr(feishu_plugin.EnvPluginFeishuOssResourceUrl, lookupStrByEnv(file_browser_plugin.EnvPluginFileBrowserResultShareDownloadUrl))
			setEnvFromStr(feishu_plugin.EnvPluginFeishuOssPageUrl, lookupStrByEnv(file_browser_plugin.EnvPluginFileBrowserResultSharePage))
			setEnvFromStr(feishu_plugin.EnvPluginFeishuOssPagePasswd, lookupStrByEnv(file_browser_plugin.EnvPluginFileBrowserResultSharePasswd))
		}

		fileBrowserCleanResultEnvErr := fileBrowserPlugin.CleanResultEnv()
		if fileBrowserCleanResultEnvErr != nil {
			log.Fatalf("fileBrowserPlugin.CleanResultEnv() err: %v", fileBrowserCleanResultEnvErr)
		}
	}

	if ossPluginErr != nil {
		setEnvFromStr(feishu_plugin.EnvPluginFeishuOssInfoSendResult, template.RenderStatusHide)
	} else {
		setEnvFromStr(feishu_plugin.EnvPluginFeishuOssInfoSendResult, template.RenderStatusShow)
	}

	// cover by feishu env oss
	feishuCfg := p.Config.FeishuCfg
	ossHost := lookupStrCoverByEnv("", feishu_plugin.EnvPluginFeishuOssHost)
	cardOss := feishu_plugin.CardOss{}
	if ossHost == "" {
		feishuCfg.RenderOssCard = feishu_plugin.RenderStatusHide
	} else {
		feishuCfg.RenderOssCard = feishu_plugin.RenderStatusShow
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

// randomStr
// new random string by cnt
func randomStr(cnt uint) string {
	var letters = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	result := make([]byte, cnt)
	keyL := len(letters)
	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(keyL)]
	}
	return string(result)
}

// randomStr
// new random string by cnt
func randomStrBySed(cnt uint, sed string) string {
	var letters = []byte(sed)
	result := make([]byte, cnt)
	keyL := len(letters)
	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(keyL)]
	}
	return string(result)
}

func setEnvFromStr(key string, val string) {
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
