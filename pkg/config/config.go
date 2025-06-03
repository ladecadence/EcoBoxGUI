package config

import "github.com/BurntSushi/toml"

const (
	version string = "0.1"
)

type Config struct {
	ConfFile string
	ID       string `toml:"id"`
	Database string `toml:"database"`
	QRPort   string `toml:"qr_port"`
	RFIDPort string `toml:"rfid_port"`
	Version  string
}

func (c *Config) GetConfig() {
	_, err := toml.DecodeFile(c.ConfFile, &c)
	if err != nil {
		panic(err)
	}
	c.Version = version
}
