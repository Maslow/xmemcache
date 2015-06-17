package config

func GetServers() []string {
	config := []string{"192.168.33.60:11211", "192.168.33.61:11211", "192.168.33.62:11211", "192.168.33.63:11211"}
	return config
}

func GetListenAddr() string {
	return ":11211"
}

func GetCopyCount() int {
	return 500
}

func GetConcurrencyCount() int {
	return 100
}
