package ui

import (
	"github.com/Ericwyn/GoTools/file"
	"github.com/Ericwyn/GoTools/shell"
	"github.com/Ericwyn/fastTray/conf"
	"github.com/Ericwyn/fastTray/log"
	"github.com/getlantern/systray"
	"reflect"
	"strings"
)

func ShowAppTray() {
	systray.Run(onReady, onExit)
}

func onReady() {

	confJson := conf.GetConfigJson()

	if confJson.IconPath == "" {
		systray.SetIcon(defaultIcon)
	} else {
		iconFile := file.OpenFile(confJson.IconPath)
		iconBytes, err := iconFile.Read()
		if err != nil {
			log.E("can't read icon file")
			return
		}

		//systray.SetIcon(icon.Data)
		systray.SetIcon(iconBytes)
	}

	//systray.SetTitle("E")
	systray.SetTooltip("fastray tools")

	itemArr := make([]*systray.MenuItem, len(confJson.Commands))
	chans := make([]chan struct{}, len(confJson.Commands))
	cmdMap := make(map[chan struct{}]string)

	for _, command := range confJson.Commands {
		item := systray.AddMenuItem(command.ItemName, command.ItemName)
		itemArr = append(itemArr, item)
		chans = append(chans, item.ClickedCh)
		cmdMap[item.ClickedCh] = command.ItemCommand
	}

	//refreshItem := systray.AddMenuItem("刷新配置", "刷新配置")
	systray.AddSeparator()
	quitItem := systray.AddMenuItem("退出", "退出程序")

	go func() {
		for {
			select {
			case <-quitItem.ClickedCh:
				log.I("Requesting quit")
				systray.Quit()
				log.I("Finished quitting")
				//case <-refreshItem.ClickedCh:
				//	log.I("refresh item")
				//	systray.Quit()
				//	go func() { ShowAppTray() }()
			}
		}
	}()

	// 另一个 goroutine 来处理其他 item
	go func() {
		cases := make([]reflect.SelectCase, len(chans))
		for i, ch := range chans {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
		}

		for {
			//chosen, value, ok := reflect.Select(cases)
			//log.I("chosen:, ", chosen, ", value: ", value, ", ok: ", ok)
			// ok will be true if the channel has not been closed.

			//msg := value.String()
			//log.I("Read from channel ", ch, " and received", msg)

			chosen, _, _ := reflect.Select(cases)
			ch := chans[chosen]
			if cmdMap[ch] != "" {
				go runCmd(cmdMap[ch])
			}
		}
	}()
}

func onExit() {
	// clean up here
}

func runCmd(cmd string) {
	log.I("run cmd: ", cmd)
	split := strings.Split(cmd, " ")
	shell.RunShellRes(split[0], split[1:]...)
}
