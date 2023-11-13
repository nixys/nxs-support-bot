package ctx

import (
	"fmt"

	conf "github.com/nixys/nxs-go-conf"
	"github.com/nixys/nxs-support-bot/misc"
)

type confOpts struct {
	LogFile  string `conf:"logfile" conf_extraopts:"default=stdout"`
	LogLevel string `conf:"loglevel" conf_extraopts:"default=info"`
	PidFile  string `conf:"pidfile"`

	Redis        redisConf        `conf:"redis"`
	MySQL        mysqlConf        `conf:"mysql" conf_extraopts:"required"`
	Telegram     telegramConf     `conf:"telegram" conf_extraopts:"required"`
	Redmine      rdmnConf         `conf:"redmine" conf_extraopts:"required"`
	Localization localizationConf `conf:"localization"`
	Cache        cacheConf        `conf:"cache"`

	API apiConf `conf:"api" conf_extraopts:"required"`
}

type redisConf struct {
	Host string `conf:"host" conf_extraopts:"default=127.0.0.1"`
	Port int    `conf:"port" conf_extraopts:"default=6379"`
}

type mysqlConf struct {
	Host     string `conf:"host" conf_extraopts:"default=127.0.0.1"`
	Port     int    `conf:"port" conf_extraopts:"default=3306"`
	DB       string `conf:"db" conf_extraopts:"required"`
	User     string `conf:"user" conf_extraopts:"required"`
	Password string `conf:"password" conf_extraopts:"required"`
}

type telegramConf struct {
	APIToken string `conf:"apiToken" conf_extraopts:"required"`
}

type rdmnConf struct {
	Host     string            `conf:"host" conf_extraopts:"required"`
	Key      string            `conf:"key" conf_extraopts:"required"`
	Feedback *rdmnFeedbackConf `conf:"feedback"`
}

type rdmnFeedbackConf struct {
	ProjectIdentifier string `conf:"projectIdentifier" conf_extraopts:"required"`
	UserID            int64  `conf:"userID" conf_extraopts:"required"`
}

type apiConf struct {
	Bind               string   `conf:"bind" conf_extraopts:"default=0.0.0.0:80"`
	TLS                *tlsConf `conf:"tls"`
	ClientMaxBodySize  string   `conf:"clientMaxBodySize" conf_extraopts:"default=36m"`
	RedmineSecretToken string   `conf:"secretToken" conf_extraopts:"required"`
}

type tlsConf struct {
	CertFile string `conf:"certfile" conf_extraopts:"required"`
	KeyFie   string `conf:"keyfile" conf_extraopts:"required"`
}

type localizationConf struct {
	Path string `conf:"path" conf_extraopts:"default=/localization"`
}

type cacheConf struct {
	TTL string `conf:"ttl" conf_extraopts:"default=5m"`
}

func confRead(confPath string) (confOpts, error) {

	var c confOpts

	err := conf.Load(&c, conf.Settings{
		ConfPath:    confPath,
		ConfType:    conf.ConfigTypeYAML,
		UnknownDeny: true,
	})
	if err != nil {
		return c, err
	}

	c.Localization.Path, err = misc.DirNormalize(c.Localization.Path)
	if err != nil {
		return c, fmt.Errorf("conf init: %w", err)
	}

	return c, err
}
