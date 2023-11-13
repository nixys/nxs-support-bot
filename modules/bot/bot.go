package tgbot

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	tg "github.com/nixys/nxs-go-telegram"
	"github.com/nixys/nxs-support-bot/misc"
	"github.com/nixys/nxs-support-bot/modules/cache"
	"github.com/nixys/nxs-support-bot/modules/issues"
	"github.com/nixys/nxs-support-bot/modules/localization"
	"github.com/nixys/nxs-support-bot/modules/users"
	"github.com/sirupsen/logrus"
)

type Settings struct {
	APIToken   string
	Log        *logrus.Logger
	Cache      cache.Cache
	RedisHost  string
	LangBundle localization.Bundle
	Feedback   *FeedbackSettings
	Issues     issues.Issues
	Users      users.Users
}

type FeedbackSettings struct {
	ProjectID int64
	UserID    int64
}

type Bot struct {
	bot tg.Telegram
}

type botCtx struct {
	log      *logrus.Logger
	c        cache.Cache
	users    users.Users
	issues   issues.Issues
	lb       localization.Bundle
	feedback *feedbackCtx
}

type feedbackCtx struct {
	projectID int64
	userID    int64
}

func Init(settings Settings) (*Bot, error) {

	// Setup the bot
	bot, err := tg.Init(
		tg.Settings{
			BotSettings: tg.SettingsBot{
				BotAPI: settings.APIToken,
			},
			RedisHost: settings.RedisHost,
		},

		tg.Description{

			Commands: []tg.Command{
				{
					Command:     "issue_create",
					Description: "Create issue",
					Handler:     issueCreateCmd,
				},
				{
					Command:     "settings",
					Description: "Settings",
					Handler:     settingsCmd,
				},
				{
					Command:     "help",
					Description: "Help",
					Handler:     helpCmd,
				},
			},

			InitHandler:  initHandler,
			PrimeHandler: primeHandler,
			ErrorHandler: errorHandler,

			States: map[tg.SessionState]tg.State{

				// Common
				stateHello: {
					StateHandler: helloState,
				},
				stateBye: {
					StateHandler: byeState,
				},
				stateHelp: {
					StateHandler:    helpState,
					CallbackHandler: helpCallback,
				},
				stateForbidden: {
					StateHandler: forbiddenState,
				},

				// Init settings
				stateInitLang: {
					StateHandler:    initLangState,
					CallbackHandler: initLangCallback,
				},
				stateInitMode: {
					StateHandler:    initModeState,
					CallbackHandler: initModeCallback,
				},
				stateInitAccount: {
					StateHandler:   initRdmnAPIKeySetState,
					MessageHandler: initRdmnAPIKeySetMsg,
				},
				stateInitRdmnAPIKeyIncorrect: {
					StateHandler: initRdmnAPIKeyIncorrectState,
				},
				stateInitEnd: {
					StateHandler: initEndState,
				},

				stateFeedback: {
					StateHandler: feedbackState,
				},

				// Settings
				stateSettings: {
					StateHandler:    settingsState,
					CallbackHandler: settingsCallback,
				},
				stateSettingsLangSelect: {
					StateHandler:    langSelectState,
					CallbackHandler: langSelectCallback,
				},
				stateSettingsRdmn: {
					StateHandler:    settingsRdmnState,
					CallbackHandler: settingsRdmnCallback,
				},
				stateSettingsRdmnAcc: {
					StateHandler:    settingsRdmnAccState,
					CallbackHandler: settingsRdmnAccCallback,
				},
				stateSettingsRdmnAPIKeySet: {
					StateHandler:    settingsRdmnAPIKeySetState,
					CallbackHandler: settingsRdmnAPIKeySetCallback,
					MessageHandler:  settingsRdmnAPIKeySetMsg,
				},
				stateSettingsRdmnAPIKeyIncorrect: {
					StateHandler: settingsRdmnAPIKeyIncorrectState,
				},

				// Issue create
				stateIssueCreate: {
					StateHandler:    issueCreateState,
					CallbackHandler: issueCreateCallback,
				},
				stateIssueCreateProject: {
					StateHandler:    issueCreateProjectState,
					CallbackHandler: issueCreateProjectCallback,
					MessageHandler:  issueCreateProjectMsg,
				},
				stateIssueCreatePriority: {
					StateHandler:    issueCreatePriorityState,
					CallbackHandler: issueCreatePriorityCallback,
				},
				stateIssueCreateSubject: {
					StateHandler:    issueCreateSubjectState,
					CallbackHandler: issueCreateSubjectCallback,
					MessageHandler:  issueCreateSubjectMsg,
				},
				stateIssueCreateDescription: {
					StateHandler:    issueCreateDescriptionState,
					CallbackHandler: issueCreateDescriptionCallback,
					MessageHandler:  issueCreateDescriptionMsg,
				},
				stateIssueCreateConfirm: {
					StateHandler:    issueCreateConfirmState,
					CallbackHandler: issueCreateConfirmCallback,
				},
				stateIssueCreateEnd: {
					StateHandler: issueCreateEndState,
					SentHandler:  issueCreateEndSent,
				},
			},
		},
		botCtx{
			log:    settings.Log,
			c:      settings.Cache,
			users:  settings.Users,
			issues: settings.Issues,
			lb:     settings.LangBundle,
			feedback: func() *feedbackCtx {
				if settings.Feedback == nil {
					return nil
				}
				return &feedbackCtx{
					projectID: settings.Feedback.ProjectID,
					userID:    settings.Feedback.UserID,
				}
			}(),
		})
	if err != nil {
		return nil, fmt.Errorf("bot init: %w", err)
	}

	return &Bot{
		bot: bot,
	}, nil
}

// UpdatesGet runtimeBotUpdates checks updates at Telegram and put it into queue
func (b *Bot) UpdatesGet(ctx context.Context, ch chan error) {
	if err := b.bot.GetUpdates(ctx); err != nil {
		if err == tg.ErrUpdatesChanClosed {
			ch <- nil
		} else {
			ch <- err
		}
	} else {
		ch <- nil
	}
}

// Queue runtimeBotQueue processes an updates from queue
func (b *Bot) Queue(ctx context.Context, ch chan error) {
	timer := time.NewTimer(time.Millisecond * 200)
	for {
		select {
		case <-timer.C:
			if err := b.bot.Processing(); err != nil {
				ch <- err
			}
			timer.Reset(time.Millisecond * 200)
		case <-ctx.Done():
			return
		}
	}
}

func initHandler(t *tg.Telegram, sess *tg.Session) (tg.InitHandlerRes, error) {

	c, err := userEnvInit(t, sess)
	if err != nil {
		return tg.InitHandlerRes{}, err
	}

	switch c.user.Type {
	case users.UserTypeFeedback:

		if sess.UpdateChain().TypeGet() == tg.UpdateTypeMessage {

			msg := strings.Join(sess.UpdateChain().MessageTextGet(), "-")

			bCtx := botCtxGet(t)

			atts, err := attachmentsIssuesUpload(t, sess)
			if err != nil {
				return tg.InitHandlerRes{}, err
			}

			if _, err := bCtx.issues.IssueFeedbackAdd(
				issues.IssueFeedbackAddData{
					TgUserID:    sess.ChatIDGet(),
					TgUsername:  sess.UserNameGet(),
					TgFirstName: sess.UserFirstNameGet(),
					TgLastName:  sess.UserLastNameGet(),
					Notes:       msg,
					Attachments: atts,
				},
			); err != nil {
				return tg.InitHandlerRes{}, err
			}
		}

		return tg.InitHandlerRes{
			NextState: tg.SessStateDestroy(),
		}, nil
	}

	return tg.InitHandlerRes{
		NextState: stateHello,
	}, nil
}

func primeHandler(t *tg.Telegram, sess *tg.Session, hs tg.HandlerSource) (tg.PrimeHandlerRes, error) {

	c, err := userEnvInit(t, sess)
	if err != nil {
		return tg.PrimeHandlerRes{}, err
	}

	state, isSessExist, err := sess.StateGet()
	if err != nil {
		return tg.PrimeHandlerRes{}, err
	}

	switch c.user.Type {
	case users.UserTypeUnauthorized:
		// User account not created in DB

		if isSessExist == true {
			switch state {
			case
				// If user account not created in DB
				// allowed only this states to interact with bot
				stateInitLang:
				return tg.PrimeHandlerRes{
					NextState: tg.SessStateContinue(),
				}, nil
			}
		}

		// For all other states (including no state) bot will be switch to
		// specified state
		return tg.PrimeHandlerRes{
			NextState: stateInitLang,
		}, nil

	case users.UserTypeFeedback:
		// User account created in DB but not linked with Redmine

		if isSessExist == true {
			switch state {
			case
				stateFeedback,
				stateInitLang,
				stateInitMode,
				stateInitAccount,
				stateInitRdmnAPIKeyIncorrect,
				stateInitEnd:
				return tg.PrimeHandlerRes{
					NextState: tg.SessStateContinue(),
				}, nil
			}

			return tg.PrimeHandlerRes{
				NextState: stateFeedback,
			}, nil
		}

		// Prevent all callbacks at this state of bot
		if sess.UpdateChain().TypeGet() == tg.UpdateTypeCallback {
			return tg.PrimeHandlerRes{
				NextState: tg.SessStateDestroy(),
			}, nil
		}

		return tg.PrimeHandlerRes{
			NextState: tg.SessStateContinue(),
		}, nil

	case users.UserTypeInternal:

		// Check processing update as a reply
		b, err := replyMessage(t, sess)
		if err != nil {
			// If error while processing message
			return tg.PrimeHandlerRes{}, err
		}
		if b == true {
			// If message was processed
			return tg.PrimeHandlerRes{}, nil
		}

		return tg.PrimeHandlerRes{
			NextState: tg.SessStateContinue(),
		}, nil
	}

	return tg.PrimeHandlerRes{}, fmt.Errorf("unknown user type")
}

func errorHandler(t *tg.Telegram, s *tg.Session, e error) (tg.ErrorHandlerRes, error) {

	ss, _, _ := s.StateGet()

	_, err := t.SendMessage(s.UserIDGet(), 0, tg.SendMessageData{
		Message: "bot error: `" + e.Error() + "` (state: `" + ss.String() + "`)",
	})
	if err != nil {
		return tg.ErrorHandlerRes{}, err
	}

	return tg.ErrorHandlerRes{}, nil
}

// replyMessage checks user message it's a reply to bunch message.
// If so message will be sent into appropriate issue
// Returning bool value indicates whether user updates was processing as a
// reply or not
func replyMessage(t *tg.Telegram, sess *tg.Session) (bool, error) {

	// Continue if updates has message type
	uc := sess.UpdateChain()
	if uc.TypeGet() != tg.UpdateTypeMessage {
		return false, nil
	}

	// Continue if updates has non-zero len
	updates := uc.Get()
	if len(updates) == 0 {
		return false, nil
	}

	upd := updates[0]

	// Continue if message replies to other message
	if upd.Message.ReplyToMessage == nil {
		return false, nil
	}

	bCtx := botCtxGet(t)

	c, err := userEnvGet(t, sess)
	if err != nil {
		return false, err
	}

	atts, err := attachmentsIssuesUpload(t, sess)
	if err != nil {
		return false, err
	}

	// Trying to reply issue
	if err := bCtx.issues.IssueReply(
		issues.IssueReplyData{
			RdmnUserID:     c.user.RdmnID,
			ChatID:         sess.ChatIDGet(),
			MessageID:      int64(upd.Message.ReplyToMessage.MessageID),
			ReplyMessageID: int64(upd.Message.MessageID),
			Note:           strings.Join(uc.MessageTextGet(), "\n"),
			Attachments:    atts,
		},
	); err != nil {
		if errors.Is(err, misc.ErrNotFound) == true {

			m, err := c.l.MessageCreate(localization.MsgOrphanedReply, nil)
			if err != nil {
				return false, err
			}

			// Send message
			if _, err := t.SendMessage(
				sess.ChatIDGet(),
				0,
				tg.SendMessageData{
					Message: m,
				},
			); err != nil {
				return false, err
			}

			return true, nil
		}
		return false, err
	}

	return true, nil
}

func botCtxGet(t *tg.Telegram) botCtx {
	return t.UsrCtxGet().(botCtx)
}
