package controller

import (
	"github.com/csby/monitor"
	"github.com/csby/wsf/server/configure"
	"github.com/csby/wsf/types"
	"net/http"
	"time"
)

type Monitor struct {
	controller
}

func NewMonitor(log types.Log, cfg *configure.Configure) *Monitor {
	instance := &Monitor{}
	instance.SetLog(log)
	instance.cfg = cfg

	return instance
}

func (s *Monitor) GetHost(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	data := monitor.Host{}
	err := data.Stat()
	if err != nil {
		a.Error(types.ErrInternal, err)
		return
	}
	a.Success(data)
}

func (s *Monitor) GetHostDoc(doc types.Doc, method string, path types.HttpPath) {
	catalog := s.createCatalog(doc, "系统信息")
	function := catalog.AddFunction(method, path, "获取主机信息")
	function.SetNote("获取当前操作系统相关信息")
	function.SetOutputDataExample(&monitor.Host{
		ID:              "8f438ea2-c26b-401e-9f6b-19f2a0e4ee2e",
		Name:            "pc",
		BootTime:        monitor.DateTime(time.Now()),
		OS:              "linux",
		Platform:        "ubuntu",
		PlatformVersion: "18.04",
		KernelVersion:   "4.15.0-22-generic",
		CPU:             "Intel(R) Core(TM) i7-6700HQ CPU @ 2.60GHz x2",
		Memory:          "4GB",
		TimeZone:        "GST+08",
	})
	function.SetInputContentType("")
}

func (s *Monitor) GetNetworkInterfaces(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	data, err := monitor.Interfaces()
	if err != nil {
		a.Error(types.ErrInternal, err)
		return
	}
	a.Success(data)
}

func (s *Monitor) GetNetworkInterfacesDoc(doc types.Doc, method string, path types.HttpPath) {
	catalog := s.createCatalog(doc, "系统信息")
	function := catalog.AddFunction(method, path, "获取网卡信息")
	function.SetNote("获取主机网卡相关信息")
	function.SetOutputDataExample([]monitor.Interface{
		{
			Name:       "本地连接",
			MTU:        1500,
			MacAddress: "00:16:5d:13:b9:70",
			IPAddress: []string{
				"fe80::b1d0:ff08:1f6f:3e0b/64",
				"192.168.1.1/24",
			},
			Flags: []string{
				"up",
				"broadcast",
				"multicast",
			},
		},
	})
	function.SetInputContentType("")
}

func (s *Monitor) GetNetworkListenPorts(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	data := monitor.TcpListenPorts()
	a.Success(data)
}

func (s *Monitor) GetNetworkListenPortsDoc(doc types.Doc, method string, path types.HttpPath) {
	catalog := s.createCatalog(doc, "系统信息")
	function := catalog.AddFunction(method, path, "获取监听端口")
	function.SetNote("获取主机正在监听端口信息")
	function.SetOutputDataExample([]monitor.Listen{
		{
			Address:  "127.0.0.1",
			Port:     163,
			Protocol: "tcp",
		},
		{
			Address:  "*",
			Port:     22,
			Protocol: "tcp",
		},
	})
	function.SetInputContentType("")
}
