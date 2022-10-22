package config

type Config struct {
	Mode       string `json:"Mode"`
	ServerPort int    `json:"ServerPort"`
	ServerHost string `json:"ServerHost"`
	BindPort   int    `json:"BindPort"`
	BindHost   string `json:"BindHost"`
	Password   string `json:"Password"`
}
