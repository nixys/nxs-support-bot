package tgbot

import (
	"git.nixys.ru/apps/nxs-support-bot/modules/localization"
	tg "github.com/nixys/nxs-go-telegram"
)

func initModeState(t *tg.Telegram, sess *tg.Session) (tg.StateHandlerRes, error) {

	bCtx := botCtxGet(t)

	// If feedback module not enabled
	if bCtx.feedback == nil {
		return tg.StateHandlerRes{
			NextState: stateInitAccount,
		}, nil
	}

	c, err := userEnvGet(t, sess)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	m, err := c.l.MessageCreate(localization.MsgInitMode, nil)
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
					Text:       c.l.BotButton(localization.ButtonAuthorize),
					Identifier: buttonIDAuthorize,
				},
			},
			{
				{
					Text:       c.l.BotButton(localization.ButtonContactUs),
					Identifier: buttonIDContactUs,
				},
			},
		},
		StickMessage: true,
	}, nil
}

func initModeCallback(t *tg.Telegram, sess *tg.Session, identifier string) (tg.CallbackHandlerRes, error) {

	switch identifier {
	case buttonIDAuthorize:
		return tg.CallbackHandlerRes{
			NextState: stateInitAccount,
		}, nil
	case buttonIDContactUs:
		return tg.CallbackHandlerRes{
			NextState: stateFeedback,
		}, nil
	}

	return tg.CallbackHandlerRes{}, nil
}
