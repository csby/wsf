package types

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type SvcArgs struct {
	LeftWidth int
	Cfg       string
	Help      bool
	Install   bool
	Uninstall bool
	Status    bool
	Start     bool
	Stop      bool
	Restart   bool
}

func (s *SvcArgs) Parse(key, value string) {
	if key == strings.ToLower("-h") ||
		key == strings.ToLower("-help") ||
		key == strings.ToLower("--help") {
		s.Help = true
	} else if key == strings.ToLower("-cfg") {
		s.Cfg = value
	} else if key == strings.ToLower("-install") {
		s.Install = true
	} else if key == strings.ToLower("-uninstall") {
		s.Uninstall = true
	} else if key == strings.ToLower("-status") {
		s.Status = true
	} else if key == strings.ToLower("-start") {
		s.Start = true
	} else if key == strings.ToLower("-stop") {
		s.Stop = true
	} else if key == strings.ToLower("-restart") {
		s.Restart = true
	}
}

func (s *SvcArgs) ShowHelp(cfgFolder, cfgName string) {
	s.ShowLine("  -help:", "[可选]显示帮助")
	s.ShowLine("  -cfg:", fmt.Sprintf("[可选]指定配置文件路径, 默认: %s", filepath.Join(cfgFolder, cfgName)))

	s.ShowLine("  -install:", "[可选]安装服务")
	s.ShowLine("  -uninstall:", "[可选]卸载服务")
	s.ShowLine("  -status:", "[可选]查看服务状态")
	s.ShowLine("  -start:", "[可选]启动服务")
	s.ShowLine("  -stop:", "[可选]停止服务")
	s.ShowLine("  -restart:", "[可选]重启服务")
}

func (s *SvcArgs) ShowLine(label, value string) {
	leftWidth := s.LeftWidth
	if leftWidth < 1 {
		leftWidth = 15
	}
	format := fmt.Sprintf("%%-%ds %%s", leftWidth)
	fmt.Printf(format, label, value)
	fmt.Println("")
}

func (s *SvcArgs) Execute(server Server) {
	if server == nil {
		return
	}
	svcName := server.ServiceName()

	if s.Install {
		err := server.Install()
		if err != nil {
			fmt.Println("install service ", svcName, " fail: ", err)
		} else {
			fmt.Println("install service ", svcName, " success")
		}
		os.Exit(21)
	} else if s.Uninstall {
		err := server.Uninstall()
		if err != nil {
			fmt.Println("uninstall service ", svcName, " fail: ", err)
		} else {
			fmt.Println("uninstall service ", svcName, " success")
		}
		os.Exit(22)
	} else if s.Status {
		status, err := server.Status()
		if err != nil {
			fmt.Println("show status of service ", svcName, " fail: ", err)
		} else {
			if status == ServerStatusRunning {
				fmt.Print(svcName, ": ")
				fmt.Println("running")
			} else if status == ServerStatusStopped {
				fmt.Print(svcName, ": ")
				fmt.Println("stopped")
			} else {
				fmt.Print(svcName, ": ")
				fmt.Println("not installed")
			}
		}
		os.Exit(23)
	} else if s.Start {
		err := server.Start()
		if err != nil {
			fmt.Println("start service ", svcName, " fail: ", err)
		} else {
			fmt.Println("start service ", svcName, " success")
		}
		os.Exit(24)
	} else if s.Stop {
		err := server.Stop()
		if err != nil {
			fmt.Println("stop service ", svcName, " fail: ", err)
		} else {
			fmt.Println("stop service ", svcName, " success")
		}
		os.Exit(25)
	} else if s.Restart {
		err := server.Restart()
		if err != nil {
			fmt.Println("restart service ", svcName, " fail: ", err)
		} else {
			fmt.Println("restart service ", svcName, " success")
		}
		os.Exit(26)
	}
}

type SvcInfo struct {
	Name     string   `json:"name" note:"服务名称"`
	Version  string   `json:"version" note:"版本号"`
	BootTime DateTime `json:"bootTime" note:"启动时间"`
	Remark   string   `json:"remark" note:"说明"`
}
