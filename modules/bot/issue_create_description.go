package tgbot

import (
	"strings"

	tg "github.com/nixys/nxs-go-telegram"
	"github.com/nixys/nxs-support-bot/modules/localization"
)

func issueCreateDescriptionState(t *tg.Telegram, sess *tg.Session) (tg.StateHandlerRes, error) {

	var issue slotIssueCreate

	c, err := userEnvGet(t, sess)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	if _, err := sess.SlotGet(slotNameIssueCreate, &issue); err != nil {
		return tg.StateHandlerRes{}, err
	}

	m, err := c.l.MessageCreate(
		localization.MsgIssueCreateDescription,
		map[string]string{
			"Project":  issue.Project.Name,
			"Priority": issue.Priority.Name,
			"Subject":  issue.Subject,
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
					Text:       c.l.BotButton(localization.ButtonLeaveEmpty),
					Identifier: buttonIDLeaveEmpty,
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

func issueCreateDescriptionCallback(t *tg.Telegram, sess *tg.Session, identifier string) (tg.CallbackHandlerRes, error) {

	switch identifier {
	case buttonIDLeaveEmpty:
		return tg.CallbackHandlerRes{
			NextState: stateIssueCreateConfirm,
		}, nil
	case buttonIDBack:
		return tg.CallbackHandlerRes{
			NextState: stateIssueCreateSubject,
		}, nil
	}

	return tg.CallbackHandlerRes{}, nil
}

func issueCreateDescriptionMsg(t *tg.Telegram, sess *tg.Session) (tg.MessageHandlerRes, error) {

	var issue slotIssueCreate

	_, err := sess.SlotGet(slotNameIssueCreate, &issue)
	if err != nil {
		return tg.MessageHandlerRes{}, err
	}

	issue.Description = strings.Join(sess.UpdateChain().MessageTextGet(), "-")

	issue.Attachments, err = attachmentsIssuesUpload(t, sess)
	if err != nil {
		return tg.MessageHandlerRes{}, err
	}

	if err := sess.SlotSave(slotNameIssueCreate, issue); err != nil {
		return tg.MessageHandlerRes{}, err
	}

	return tg.MessageHandlerRes{
		NextState: stateIssueCreateConfirm,
	}, nil
}
