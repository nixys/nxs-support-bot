package tgbot

import (
	"errors"
	"strings"

	"git.nixys.ru/apps/nxs-support-bot/misc"
	"git.nixys.ru/apps/nxs-support-bot/modules/issues"
	"git.nixys.ru/apps/nxs-support-bot/modules/localization"
	"git.nixys.ru/apps/nxs-support-bot/modules/users"
	tg "github.com/nixys/nxs-go-telegram"
)

func initRdmnAPIKeySetState(t *tg.Telegram, sess *tg.Session) (tg.StateHandlerRes, error) {

	c, err := userEnvGet(t, sess)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	m, err := c.l.MessageCreate(localization.MsgInitRdmnApiKeySet, nil)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	return tg.StateHandlerRes{
		Message:               m,
		ParseMode:             tg.ParseModeHTML,
		DisableWebPagePreview: true,
		StickMessage:          true,
	}, nil
}

func initRdmnAPIKeySetMsg(t *tg.Telegram, sess *tg.Session) (tg.MessageHandlerRes, error) {

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
			return tg.MessageHandlerRes{
				NextState: stateInitRdmnAPIKeyIncorrect,
			}, nil
		}
		return tg.MessageHandlerRes{}, err
	}

	if err = sess.SlotSave(slotNameUser, user); err != nil {
		return tg.MessageHandlerRes{}, err
	}

	// Send close message into feedback issue (if necessary)
	if err := bCtx.issues.IssueFeedbackClose(
		issues.IssueFeedbackCloseData{
			TgUserID:      sess.UserIDGet(),
			RdmnUserID:    user.RdmnID,
			RdmnLogin:     user.Login,
			RdmnFirstName: user.FirstName,
			RdmnLastName:  user.LastName,
		},
	); err != nil {
		return tg.MessageHandlerRes{}, err
	}

	return tg.MessageHandlerRes{
		NextState: stateInitEnd,
	}, nil
}
