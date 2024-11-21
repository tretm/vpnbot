package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	BotConfig        `yaml:"bot"`
	StorageConfig    `yaml:"storage"`
	ServerHttpConfig `yaml:"server"`
	OutlineVpn       `yaml:"outline"`
	Vray             `yaml:"vray"`
	YouMoneyConfig   `yaml:"youmoney"`
}
type Vray struct {
	UrlVray         string `yaml:"urlvray"`
	VrayLogin       string `yaml:"vraylogin"`
	VrayPassword    string `yaml:"vraypassword"`
	LimitConnection int    `yaml:"limitconnection"`
}
type OutlineVpn struct {
	UrlOutline string `yaml:"urloutline"`
	DomeName   string `yaml:"domenname"`
	Port       string `yaml:"port"`
	Method     string `yaml:"method"`
	Sscon      string `yaml:"ssconf"`
}
type BotConfig struct {
	BotToken      string   `yaml:"botToken"`
	BotPayToken   string   `yaml:"botPayToken"`
	ServerPath    string   `yaml:"serverPath"`
	ReportChanId  int64    `yaml:"reportchanid"`
	AdminId       int64    `yaml:"adminid"`
	LinkLen       int      `yaml:"linklen"`
	DefaultDays   int      `yaml:"defaultdays"`
	TrialPeriod   int      `yaml:"trialperiod"`
	PayAmount     []int    `yaml:"payamount"`
	PriceOneMonth int      `yaml:"priceonemonth"`
	Discount      int      `yaml:"discount"`
	PeriodsText   []string `yaml:"periodstext"`
	PeriodsVal    []int    `yaml:"periodsval"`
	Support       string   `yaml:"support"`
}

type StorageConfig struct {
	Host        string        `yaml:"host" env-required:"true" env:"DB_HOST" env-description:""`
	Port        string        `yaml:"port" env-required:"true" env:"DB_PORT"`
	Database    string        `yaml:"database" env-default:"tlgur" env:"DB_NAME" env-description:"Database name"`
	User        string        `yaml:"user" env-required:"true" env:"DB_USER"`
	Pass        string        `yaml:"password" env-required:"true" env:"DB_PASSWORD"`
	RowsLimit   int           `yaml:"rowslimit"`
	MaxOpenConn int           `yaml:"max_open_conn" env-description:"How many concurrent db connections are allowed"`
	MaxIdleConn int           `yaml:"max_idle_conn" env-description:"How many concurrent idle db connections are allowed"`
	ConnMaxLife time.Duration `yaml:"conn_max_life" env-description:"TTL for opened db connections"`
}
type ServerHttpConfig struct {
	Port        string        `yaml:"port" env-default:"10101" env:"HTTP_PORT"`
	Timeout     time.Duration `yaml:"timeout" env-default:"1s" env-description:"Read and Write timeouts"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"10s" env-description:"Timeout for idle connection"`
	Ttl         time.Duration `yaml:"ttl" env-default:"10s" env-description:"Timeout for idle connection"`
	Secret      string        `yaml:"secret" env-default:"123456"`
	HiddenRoute string        `yaml:"hiddenRoute" env-default:"asghdsgfhsgkfj:ehruiheur5"`
}

type YouMoneyConfig struct {
	ClientId    string `yaml:"client_id"`
	RedirectUrl string `yaml:"redirect_url"`
	Token       string `yaml:"token_ym"`
}

func PrepareConfig() (*Config, error) {
	const op = "config.prepare"
	const configPath = "config.yaml"

	var cfg Config

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &cfg, fmt.Errorf("%s: config file does not exist `%s`", op, configPath)
	}
	// read config from yaml and environment
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return &cfg, fmt.Errorf("%s: failed to read config, %w", op, err)
	}
	return &cfg, nil

}
