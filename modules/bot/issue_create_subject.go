package tgbot

import (
	"strings"

	tg "github.com/nixys/nxs-go-telegram"
	"github.com/nixys/nxs-support-bot/modules/localization"
)

func issueCreateSubjectState(t *tg.Telegram, sess *tg.Session) (tg.StateHandlerRes, error) {

	var issue slotIssueCreate

	c, err := userEnvGet(t, sess)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	if _, err := sess.SlotGet(slotNameIssueCreate, &issue); err != nil {
		return tg.StateHandlerRes{}, err
	}

	m, err := c.l.MessageCreate(
		localization.MsgIssueCreateSubject,
		map[string]string{
			"Project":  issue.Project.Name,
			"Priority": issue.Priority.Name,
		},
	)
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
					Text:       c.l.BotButton(localization.ButtonBack),
					Identifier: buttonIDBack,
				},
			},
		},
		StickMessage: true,
	}, nil
}

func issueCreateSubjectCallback(t *tg.Telegram, sess *tg.Session, identifier string) (tg.CallbackHandlerRes, error) {

	switch identifier {
	case buttonIDBack:
		return tg.CallbackHandlerRes{
			NextState: stateIssueCreate,
		}, nil
	}

	return tg.CallbackHandlerRes{}, nil
}

func issueCreateSubjectMsg(t *tg.Telegram, sess *tg.Session) (tg.MessageHandlerRes, error) {

	var issue slotIssueCreate

	if _, err := sess.SlotGet(slotNameIssueCreate, &issue); err != nil {
		return tg.MessageHandlerRes{}, err
	}

	issue.Subject = strings.Join(sess.UpdateChain().MessageTextGet(), "-")

	if err := sess.SlotSave(slotNameIssueCreate, issue); err != nil {
		return tg.MessageHandlerRes{}, err
	}

	return tg.MessageHandlerRes{
		NextState: stateIssueCreateDescription,
	}, nil
}
