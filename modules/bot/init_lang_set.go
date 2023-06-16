package tgbot

import (
	"git.nixys.ru/apps/nxs-support-bot/modules/localization"
	"git.nixys.ru/apps/nxs-support-bot/modules/users"
	tg "github.com/nixys/nxs-go-telegram"
)

func initLangState(t *tg.Telegram, sess *tg.Session) (tg.StateHandlerRes, error) {

	c, err := userEnvGet(t, sess)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	m, err := c.l.MessageCreate(localization.MsgInitLang, nil)
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
		},
		StickMessage: true,
	}, nil
}

func initLangCallback(t *tg.Telegram, sess *tg.Session, identifier string) (tg.CallbackHandlerRes, error) {

	var l string

	switch identifier {
	case buttonIDEN, buttonIDRU:
		l = identifier
	default:
		l = buttonIDEN
	}

	bCtx := botCtxGet(t)

	u, err := bCtx.users.UserUpdate(
		sess.UserIDGet(),
		users.UserUpdateData{
			Lang: &l,
		},
	)
	if err != nil {
		return tg.CallbackHandlerRes{}, err
	}

	if err = sess.SlotSave(slotNameUser, u); err != nil {
		return tg.CallbackHandlerRes{}, err
	}

	return tg.CallbackHandlerRes{
		NextState: stateInitMode,
	}, nil
}
