package tgbot

import (
	"git.nixys.ru/apps/nxs-support-bot/modules/localization"
	tg "github.com/nixys/nxs-go-telegram"
)

func settingsRdmnAccState(t *tg.Telegram, sess *tg.Session) (tg.StateHandlerRes, error) {

	c, err := userEnvGet(t, sess)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	m, err := c.l.MessageCreate(
		localization.MsgSettingsRdmnAcc,
		map[string]string{
			"Login":     c.user.Login,
			"FirstName": c.user.FirstName,
			"LastName":  c.user.LastName,
		},
	)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	buttons := [][]tg.Button{
		{
			{
				Text:       c.l.BotButton(localization.ButtonLink),
				Identifier: buttonIDLink,
			},
		},
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

func settingsRdmnAccCallback(t *tg.Telegram, sess *tg.Session, identifier string) (tg.CallbackHandlerRes, error) {

	switch identifier {
	case buttonIDLink:
		return tg.CallbackHandlerRes{
			NextState: stateSettingsRdmnAPIKeySet,
		}, nil
	case buttonIDFavoriteProject:
		return tg.CallbackHandlerRes{}, nil
	case buttonIDBack:
		return tg.CallbackHandlerRes{
			NextState: stateSettingsRdmn,
		}, nil
	}

	return tg.CallbackHandlerRes{}, nil
}
