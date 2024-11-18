package xredis

type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database int    `json:"database"`
	Alias    string `json:"alias"`
	UserName string `json:"userName"`
	Password string `json:"password"`
}
