package ctx

import (
	"fmt"
	"os"
	"time"

	"github.com/docker/go-units"
	appctx "github.com/nixys/nxs-go-appctx/v3"
	"github.com/nixys/nxs-support-bot/ds/primedb"
	"github.com/nixys/nxs-support-bot/ds/redmine"
	tgbot "github.com/nixys/nxs-support-bot/modules/bot"
	"github.com/nixys/nxs-support-bot/modules/cache"
	"github.com/nixys/nxs-support-bot/modules/issues"
	"github.com/nixys/nxs-support-bot/modules/localization"
	"github.com/nixys/nxs-support-bot/modules/task-handlers/rdmnhndlr"
	"github.com/nixys/nxs-support-bot/modules/users"
	"github.com/sirupsen/logrus"
)

// Ctx defines application custom context
type Ctx struct {
	Conf      confOpts
	Cache     cacheSettings
	Bot       *tgbot.Bot
	API       apiSettings
	Rdmnhndlr rdmnhndlr.RdmnHndlr
	Log       *logrus.Logger
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

func AppCtxInit() (any, error) {

	c := &Ctx{}

	args, err := ArgsRead()
	if err != nil {
		return nil, err
	}

	conf, err := confRead(args.ConfigPath)
	if err != nil {
		tmpLogError("ctx init", err)
		return nil, err
	}

	c.Log, err = logInit(conf.LogFile, conf.LogLevel)
	if err != nil {
		tmpLogError("ctx init", err)
		return nil, err
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
		c.Log.WithFields(logrus.Fields{
			"details": err,
		}).Errorf("ctx init")
		return nil, err
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
		c.Log.WithFields(logrus.Fields{
			"details": err,
		}).Errorf("ctx init")
		return nil, err
	}

	// Localization init
	lb, err := localization.Init(c.Conf.Localization.Path)
	if err != nil {
		c.Log.WithFields(logrus.Fields{
			"details": err,
		}).Errorf("ctx init")
		return nil, err
	}

	var feedback *feedbackSettings
	if c.Conf.Redmine.Feedback != nil {

		proj, err := rdmn.ProjectGetByIdentifier(c.Conf.Redmine.Feedback.ProjectIdentifier)
		if err != nil {
			c.Log.WithFields(logrus.Fields{
				"details": err,
			}).Errorf("ctx init")
			return nil, err
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
		Log:        c.Log,
		Cache:      c.Cache.C,
		RedisHost:  redisHost,
		LangBundle: lb,
		Issues:     iss,
		Users:      usrs,
		Feedback:   (*tgbot.FeedbackSettings)(feedback),
	})
	if err != nil {
		c.Log.WithFields(logrus.Fields{
			"details": err,
		}).Errorf("ctx init")
		return nil, err
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
		c.Log.WithFields(logrus.Fields{
			"details": err,
		}).Errorf("ctx init")
		return nil, err
	}

	c.Cache.TTL, err = time.ParseDuration(c.Conf.Cache.TTL)
	if err != nil {
		c.Log.WithFields(logrus.Fields{
			"details": err,
		}).Errorf("ctx init")
		return nil, err
	}

	return c, nil
}

func tmpLogError(msg string, err error) {
	l, _ := appctx.DefaultLogInit(os.Stderr, logrus.InfoLevel, &logrus.JSONFormatter{})
	l.WithFields(logrus.Fields{
		"details": err,
	}).Errorf(msg)
}

func logInit(file, level string) (*logrus.Logger, error) {

	var (
		f   *os.File
		err error
	)

	switch file {
	case "stdout":
		f = os.Stdout
	case "stderr":
		f = os.Stderr
	default:
		f, err = os.OpenFile(file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
		if err != nil {
			return nil, fmt.Errorf("log init: %w", err)
		}
	}

	// Validate log level
	l, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, fmt.Errorf("log init: %w", err)
	}

	return appctx.DefaultLogInit(f, l, &logrus.JSONFormatter{})
}
