package ctx

import (
	"fmt"
	"time"

	"github.com/docker/go-units"
	appctx "github.com/nixys/nxs-go-appctx/v2"
	"github.com/nixys/nxs-support-bot/ds/primedb"
	"github.com/nixys/nxs-support-bot/ds/redmine"
	tgbot "github.com/nixys/nxs-support-bot/modules/bot"
	"github.com/nixys/nxs-support-bot/modules/cache"
	"github.com/nixys/nxs-support-bot/modules/issues"
	"github.com/nixys/nxs-support-bot/modules/localization"
	"github.com/nixys/nxs-support-bot/modules/task-handlers/rdmnhndlr"
	"github.com/nixys/nxs-support-bot/modules/users"
)

// Ctx defines application custom context
type Ctx struct {
	Conf      confOpts
	Cache     cacheSettings
	Bot       *tgbot.Bot
	API       apiSettings
	Rdmnhndlr rdmnhndlr.RdmnHndlr
}

type apiSettings struct {
	ClientMaxBodySizeBytes int64
}

type cacheSettings struct {
	C   cache.Cache
	TTL time.Duration
}

type feedbackSettings struct {
	ProjectID int64
	UserID    int64
}

// Init initiates application custom context
func (c *Ctx) Init(opts appctx.CustomContextFuncOpts) (appctx.CfgData, error) {

	// Read config file
	conf, err := confRead(opts.Config)
	if err != nil {
		return appctx.CfgData{}, err
	}

	// Set application context
	c.Conf = conf

	// Connect to MySQL
	primeDB, err := primedb.Connect(primedb.Settings{
		Host:     c.Conf.MySQL.Host,
		Port:     c.Conf.MySQL.Port,
		Database: c.Conf.MySQL.DB,
		User:     c.Conf.MySQL.User,
		Password: c.Conf.MySQL.Password,
	})
	if err != nil {
		return appctx.CfgData{}, err
	}

	// Set redmine context
	rdmn := redmine.Init(c.Conf.Redmine.Host, c.Conf.Redmine.Key)

	redisHost := fmt.Sprintf("%s:%d", c.Conf.Redis.Host, c.Conf.Redis.Port)

	// Set cache
	c.Cache.C, err = cache.Init(cache.Settings{
		Redmine:   rdmn,
		RedisHost: redisHost,
	})
	if err != nil {
		return appctx.CfgData{}, err
	}

	// Localization init
	lb, err := localization.Init(c.Conf.Localization.Path)
	if err != nil {
		return appctx.CfgData{}, err
	}

	var feedback *feedbackSettings
	if c.Conf.Redmine.Feedback != nil {

		proj, err := rdmn.ProjectGetByIdentifier(c.Conf.Redmine.Feedback.ProjectIdentifier)
		if err != nil {
			return appctx.CfgData{}, err
		}

		feedback = &feedbackSettings{
			ProjectID: proj.ID,
			UserID:    c.Conf.Redmine.Feedback.UserID,
		}
	}

	iss := issues.Init(
		issues.Settings{
			DB:       primeDB,
			Redmine:  rdmn,
			Feedback: (*issues.FeedbackSettings)(feedback),
		},
	)

	usrs := users.Init(
		users.Settings{
			DB:       primeDB,
			Cache:    c.Cache.C,
			Redmine:  rdmn,
			Feedback: (*users.FeedbackSettings)(feedback),
		},
	)

	// Set bot context
	c.Bot, err = tgbot.Init(tgbot.Settings{
		APIToken:   c.Conf.Telegram.APIToken,
		Log:        opts.Log,
		Cache:      c.Cache.C,
		RedisHost:  redisHost,
		LangBundle: lb,
		Issues:     iss,
		Users:      usrs,
		Feedback:   (*tgbot.FeedbackSettings)(feedback),
	})
	if err != nil {
		return appctx.CfgData{}, err
	}

	c.Rdmnhndlr = rdmnhndlr.Init(
		rdmnhndlr.Settings{
			Bot:        c.Bot,
			LangBundle: lb,
			Users:      usrs,
			Issues:     iss,
			Feedback:   (*rdmnhndlr.FeedbackSettings)(feedback),
		},
	)

	c.API.ClientMaxBodySizeBytes, err = units.RAMInBytes(c.Conf.API.ClientMaxBodySize)
	if err != nil {
		return appctx.CfgData{}, err
	}

	c.Cache.TTL, err = time.ParseDuration(c.Conf.Cache.TTL)
	if err != nil {
		return appctx.CfgData{}, err
	}

	return appctx.CfgData{
		LogFile:  c.Conf.LogFile,
		LogLevel: c.Conf.LogLevel,
		PidFile:  c.Conf.PidFile,
	}, nil
}

// Reload reloads application custom context
func (c *Ctx) Reload(opts appctx.CustomContextFuncOpts) (appctx.CfgData, error) {

	opts.Log.Debug("reloading context")

	return c.Init(opts)
}

// Free frees application custom context
func (c *Ctx) Free(opts appctx.CustomContextFuncOpts) int {

	opts.Log.Debug("freeing context")

	return 0
}
