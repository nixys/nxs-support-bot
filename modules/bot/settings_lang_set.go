package tgbot

import (
	tg "github.com/nixys/nxs-go-telegram"
	"github.com/nixys/nxs-support-bot/modules/localization"
	"github.com/nixys/nxs-support-bot/modules/users"
)

func langSelectState(t *tg.Telegram, sess *tg.Session) (tg.StateHandlerRes, error) {

	c, err := userEnvGet(t, sess)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	m, err := c.l.MessageCreate(localization.MsgSettingsLang, nil)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	return tg.StateHandlerRes{
		Message:               m,
		ParseMode:             tg.ParseModeHTML,
		DisableWebPagePreview: true,
		Buttons: [][]tg.Button{
			{
				{
					Text:       c.l.BotButton(localization.ButtonEN),
					Identifier: buttonIDEN,
				},
			},
			{
				{
					Text:       c.l.BotButton(localization.ButtonRU),
					Identifier: buttonIDRU,
				},
			},
			{
				{
					Text:       c.l.BotButton(localization.ButtonBack),
					Identifier: buttonIDBack,
				},
			},
		},
		StickMessage: true,
	}, nil
}

func langSelectCallback(t *tg.Telegram, sess *tg.Session, identifier string) (tg.CallbackHandlerRes, error) {

	var l string

	bCtx := botCtxGet(t)

	switch identifier {
	case buttonIDBack:
		return tg.CallbackHandlerRes{
			NextState: stateSettings,
		}, nil
	case buttonIDEN, buttonIDRU:
		l = identifier
	default:
		l = buttonIDEN
	}

	user, err := bCtx.users.UserUpdate(
		sess.UserIDGet(),
		users.UserUpdateData{
			Lang: &l,
		},
	)
	if err != nil {
		return tg.CallbackHandlerRes{}, err
	}

	err = sess.SlotSave(slotNameUser, user)
	if err != nil {
		return tg.CallbackHandlerRes{}, err
	}

	return tg.CallbackHandlerRes{
		NextState: stateSettings,
	}, nil
}
