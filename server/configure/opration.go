package configure

type Operation struct {
	Root  string `json:"root" note:"物理路径"`
	Api   Api    `json:"api" note:"接口"`
	Users []User `json:"users" note:"用户"`
}
