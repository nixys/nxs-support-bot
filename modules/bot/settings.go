package tgbot

import (
	tg "github.com/nixys/nxs-go-telegram"
	"github.com/nixys/nxs-support-bot/modules/localization"
	"github.com/nixys/nxs-support-bot/modules/users"
)

func settingsCmd(t *tg.Telegram, sess *tg.Session, cmd string, args string) (tg.CommandHandlerRes, error) {

	c, err := userEnvGet(t, sess)
	if err != nil {
		return tg.CommandHandlerRes{}, err
	}

	switch c.user.Type {
	case users.UserTypeUnauthorized:
		return tg.CommandHandlerRes{
			NextState: stateInitLang,
		}, nil
	case users.UserTypeFeedback:
		return tg.CommandHandlerRes{
			NextState: stateInitAccount,
		}, nil
	case users.UserTypeInternal:
		return tg.CommandHandlerRes{
			NextState: stateSettings,
		}, nil
	}

	return tg.CommandHandlerRes{
		NextState: tg.SessStateBreak(),
	}, nil
}

func settingsState(t *tg.Telegram, sess *tg.Session) (tg.StateHandlerRes, error) {

	c, err := userEnvGet(t, sess)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	m, err := c.l.MessageCreate(localization.MsgSettings, nil)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	buttons := [][]tg.Button{
		{
			{
				Text:       c.l.BotButton(localization.ButtonRdmn),
				Identifier: buttonIDRdmn,
			},
		},
		{
			{
				Text:       c.l.BotButton(localization.ButtonLang),
				Identifier: buttonIDLang,
			},
		},
		{
			{
				Text:       c.l.BotButton(localization.ButtonCancel),
				Identifier: buttonIDCancel,
			},
		},
	}

	return tg.StateHandlerRes{
		Message:               m,
		ParseMode:             tg.ParseModeHTML,
		DisableWebPagePreview: true,
		Buttons:               buttons,
		StickMessage:          true,
	}, nil
}

func settingsCallback(t *tg.Telegram, sess *tg.Session, identifier string) (tg.CallbackHandlerRes, error) {

	switch identifier {
	case buttonIDRdmn:
		return tg.CallbackHandlerRes{
			NextState: stateSettingsRdmn,
		}, nil
	case buttonIDLang:
		return tg.CallbackHandlerRes{
			NextState: stateSettingsLangSelect,
		}, nil
	case buttonIDCancel:
		return tg.CallbackHandlerRes{
			NextState: stateBye,
		}, nil
	}

	return tg.CallbackHandlerRes{}, nil
}
