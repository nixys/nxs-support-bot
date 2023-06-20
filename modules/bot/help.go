package tgbot

import (
	tg "github.com/nixys/nxs-go-telegram"
	"github.com/nixys/nxs-support-bot/modules/localization"
)

func helpCmd(t *tg.Telegram, sess *tg.Session, cmd string, args string) (tg.CommandHandlerRes, error) {
	return tg.CommandHandlerRes{
		NextState: stateHelp,
	}, nil
}

func helpState(t *tg.Telegram, sess *tg.Session) (tg.StateHandlerRes, error) {

	c, err := userEnvGet(t, sess)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	m, err := c.l.MessageCreate(localization.MsgHelp, nil)
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
					Text:       c.l.BotButton(localization.ButtonQuit),
					Identifier: buttonIDQuit,
				},
			},
		},
		StickMessage: true,
	}, nil
}

func helpCallback(t *tg.Telegram, sess *tg.Session, identifier string) (tg.CallbackHandlerRes, error) {

	switch identifier {
	case buttonIDQuit:
		return tg.CallbackHandlerRes{
			NextState: stateBye,
		}, nil
	}

	return tg.CallbackHandlerRes{}, nil
}
