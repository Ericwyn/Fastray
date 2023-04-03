package conf

import (
	"encoding/json"
	"github.com/Ericwyn/GoTools/file"
	"github.com/Ericwyn/fastTray/log"
	"os"
	"path/filepath"
	"strings"
)

// ConfigJson 配置文件的结构体
type ConfigJson struct {
	IconPath string
	Commands []Command
}

// Command 命令
type Command struct {
	ItemName    string
	ItemID      string
	ItemCommand string
}

func defaultConf() ConfigJson {
	return ConfigJson{
		Commands: []Command{
			Command{
				ItemName:    "截屏",
				ItemID:      "screenshot",
				ItemCommand: "gnome-screenshot  -a",
			},
			Command{
				ItemName:    "文件管理器",
				ItemID:      "open download dir",
				ItemCommand: "nautilus",
			},
		},
	}
}

func GetConfigJson() ConfigJson {
	configDir := file.OpenFile(GetRunnerPath() + "/.conf/")
	if !configDir.Exits() {
		log.I("create config dir")
		configDir.Mkdirs()
	}

	configFile := file.OpenFile(GetRunnerPath() + "/.conf/config.json")
	if configFile.Exits() {
		var conf ConfigJson
		read, err := configFile.Read()
		if err != nil {
			log.E("can't read config file", err)
			return defaultConf()
		}
		err = json.Unmarshal(read, &conf)
		if err != nil {
			log.E("can't parse config file", err)
			return defaultConf()
		}
		return conf
	} else {
		jsonStr, _ := json.MarshalIndent(defaultConf(), "", "	")

		err := configFile.WriteDate(file.W_NEW, string(jsonStr))
		if err != nil {
			log.E("can't write config file", err)
			return defaultConf()
		}
		log.I("create default config file in ", configFile.AbsPath())
	}

	return defaultConf()
}

var runnerPath = ""

func GetRunnerPath() string {
	if runnerPath == "" {
		//返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径

		log.D("os.Args[0]:" + os.Args[0])

		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.E("can't get runner path")
			log.E(err)
		}

		//将\替换成/
		runnerPath = strings.Replace(dir, "\\", "/", -1)

		log.D("programme run in path :" + runnerPath)

		// 如果运行的目录是在 Temp 下面的话, 那么看看 ./ 目录是什么
		if strings.Contains(runnerPath, "AppData/Local/Temp") ||
			strings.HasPrefix(runnerPath, "/tmp") {
			log.D("the runner path is in Temp dir")
			dir, err := filepath.Abs(filepath.Dir("./"))
			if err != nil {
				log.E("get ./ path failed")
				log.E(err)
			}

			runnerPath = strings.Replace(dir, "\\", "/", -1)
			log.D("the runner path update to : " + runnerPath)
		}
	}

	return runnerPath
}
