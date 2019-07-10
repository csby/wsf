package configure

type User struct {
	Account  string `json:"account" note:"账号"`
	Password string `json:"password" note:"密码"`
}
