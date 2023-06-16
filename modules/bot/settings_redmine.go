package tgbot

import (
	"git.nixys.ru/apps/nxs-support-bot/modules/localization"
	tg "github.com/nixys/nxs-go-telegram"
)

func settingsRdmnState(t *tg.Telegram, sess *tg.Session) (tg.StateHandlerRes, error) {

	c, err := userEnvGet(t, sess)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	m, err := c.l.MessageCreate(localization.MsgSettingsRdmn, nil)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	buttons := [][]tg.Button{
		{
			{
				Text:       c.l.BotButton(localization.ButtonAccount),
				Identifier: buttonIDAccount,
			},
		},
		/*
			{
				{
					Text:       l.BotButton(localization.ButtonFavoriteProjects),
					Identifier: buttonIDFavoriteProject,
				},
			},
		*/
		{
			{
				Text:       c.l.BotButton(localization.ButtonBack),
				Identifier: buttonIDBack,
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

func settingsRdmnCallback(t *tg.Telegram, sess *tg.Session, identifier string) (tg.CallbackHandlerRes, error) {

	switch identifier {
	case buttonIDAccount:
		return tg.CallbackHandlerRes{
			NextState: stateSettingsRdmnAcc,
		}, nil
	case buttonIDFavoriteProject:
		return tg.CallbackHandlerRes{}, nil
	case buttonIDBack:
		return tg.CallbackHandlerRes{
			NextState: stateSettings,
		}, nil
	}

	return tg.CallbackHandlerRes{}, nil
}
