package config

type Config struct {
	Port       int    `json:"port"`
	Host       string `json:"host"`
	HostString string `json:"host_str"`
}
