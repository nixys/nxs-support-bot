package tgbot

import (
	"git.nixys.ru/apps/nxs-support-bot/modules/localization"
	tg "github.com/nixys/nxs-go-telegram"
)

func feedbackState(t *tg.Telegram, sess *tg.Session) (tg.StateHandlerRes, error) {

	c, err := userEnvGet(t, sess)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	m, err := c.l.MessageCreate(localization.MsgFeedbackGreetings, nil)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	return tg.StateHandlerRes{
		Message:               m,
		ParseMode:             tg.ParseModeHTML,
		DisableWebPagePreview: true,
		StickMessage:          true,
		NextState:             tg.SessStateDestroy(),
	}, nil
}
