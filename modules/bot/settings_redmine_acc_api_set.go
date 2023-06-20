package tgbot

import (
	"errors"
	"strings"

	tg "github.com/nixys/nxs-go-telegram"
	"github.com/nixys/nxs-support-bot/misc"
	"github.com/nixys/nxs-support-bot/modules/localization"
	"github.com/nixys/nxs-support-bot/modules/users"
)

func settingsRdmnAPIKeySetState(t *tg.Telegram, sess *tg.Session) (tg.StateHandlerRes, error) {

	c, err := userEnvGet(t, sess)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	m, err := c.l.MessageCreate(localization.MsgSettingsRdmnApiKeySet, nil)
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
					Text:       c.l.BotButton(localization.ButtonCancel),
					Identifier: buttonIDCancel,
				},
			},
		},
		StickMessage: true,
	}, nil
}

func settingsRdmnAPIKeySetCallback(t *tg.Telegram, sess *tg.Session, identifier string) (tg.CallbackHandlerRes, error) {

	switch identifier {
	case buttonIDCancel:
		return tg.CallbackHandlerRes{
			NextState: stateSettingsRdmnAcc,
		}, nil
	}

	return tg.CallbackHandlerRes{}, nil
}

func settingsRdmnAPIKeySetMsg(t *tg.Telegram, sess *tg.Session) (tg.MessageHandlerRes, error) {

	key := strings.Join(sess.UpdateChain().MessageTextGet(), "-")

	bCtx := botCtxGet(t)

	user, err := bCtx.users.UserUpdate(
		sess.UserIDGet(),
		users.UserUpdateData{
			RedmineKey: &key,
		},
	)
	if err != nil {
		if errors.Is(err, misc.ErrAPIKey) == true {
			// If api key incorrect
			return tg.MessageHandlerRes{
				NextState: stateSettingsRdmnAPIKeyIncorrect,
			}, nil
		}
		return tg.MessageHandlerRes{}, err
	}

	if err = sess.SlotSave(slotNameUser, user); err != nil {
		return tg.MessageHandlerRes{}, err
	}

	return tg.MessageHandlerRes{
		NextState: stateSettingsRdmnAcc,
	}, nil
}
