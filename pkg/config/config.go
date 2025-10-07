package config

import "github.com/BurntSushi/toml"

const (
	version string = "0.1"
)

type Config struct {
	ConfFile  string
	Cabinet   string `toml:"cabinet"`
	Database  string `toml:"database"`
	QRPort    string `toml:"qr_port"`
	RFIDPort  string `toml:"rfid_port"`
	RFIDPort2 string `toml:"rfid_port2"`
	QRPass    string `toml:"qr_pass"`
	Version   string
}

func (c *Config) GetConfig() {
	_, err := toml.DecodeFile(c.ConfFile, &c)
	if err != nil {
		panic(err)
	}
	c.Version = version
}
