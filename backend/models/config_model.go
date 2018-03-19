package models

type Config struct {
	DB   DBConfig   `json:"db"`
	JWT  JWTConfig  `json:"jwt"`
	Hash HashConfig `json:"hash"`
}

type DBConfig struct {
	DBUser     string `json:"dbuser"`
	DBPassword string `json:"dbpassword"`
	Host       string `json:"host"`
	Port       string `json:"port"`
	DBName     string `json:"dbname"`
}

type JWTConfig struct {
	Secret string `json:"secret"`
}

type HashConfig struct {
	Secret string `json:"secret"`
}
