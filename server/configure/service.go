package configure

import "time"

type Service struct {
	BootTime time.Time `json:"-" note:"启动时间"`
	Name     string    `json:"name" note:"服务名称，系统内唯一"`
}
