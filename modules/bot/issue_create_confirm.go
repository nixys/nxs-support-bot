package tgbot

import (
	tg "github.com/nixys/nxs-go-telegram"
	"github.com/nixys/nxs-support-bot/modules/issues"
	"github.com/nixys/nxs-support-bot/modules/localization"
)

func issueCreateConfirmState(t *tg.Telegram, sess *tg.Session) (tg.StateHandlerRes, error) {

	var issue slotIssueCreate

	c, err := userEnvGet(t, sess)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	if _, err := sess.SlotGet(slotNameIssueCreate, &issue); err != nil {
		return tg.StateHandlerRes{}, err
	}

	m, err := c.l.MessageCreate(
		localization.MsgIssueCreateConfirm,
		map[string]string{
			"Project":     issue.Project.Name,
			"Priority":    issue.Priority.Name,
			"Subject":     issue.Subject,
			"Description": issue.Description,
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
					Text:       c.l.BotButton(localization.ButtonCreateIssue),
					Identifier: buttonIDCreateIssue,
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

func issueCreateConfirmCallback(t *tg.Telegram, sess *tg.Session, identifier string) (tg.CallbackHandlerRes, error) {

	var issue slotIssueCreate

	switch identifier {
	case buttonIDCreateIssue:

		bCtx := botCtxGet(t)

		c, err := userEnvGet(t, sess)
		if err != nil {
			return tg.CallbackHandlerRes{}, err
		}

		if _, err := sess.SlotGet(slotNameIssueCreate, &issue); err != nil {
			return tg.CallbackHandlerRes{}, err
		}

		issueID, err := bCtx.issues.IssueCreate(
			c.user.RdmnID,
			issues.IssueCreateData{
				ProjectID:   issue.Project.ID,
				PriorityID:  issue.Priority.ID,
				Subject:     issue.Subject,
				Description: issue.Description,
				Attachments: issue.Attachments,
			},
		)
		if err != nil {
			return tg.CallbackHandlerRes{}, err
		}

		// Save new issue ID
		issue.IssueID = issueID

		if err := sess.SlotSave(slotNameIssueCreate, issue); err != nil {
			return tg.CallbackHandlerRes{}, err
		}

		return tg.CallbackHandlerRes{
			NextState: stateIssueCreateEnd,
		}, nil

	case buttonIDBack:
		return tg.CallbackHandlerRes{
			NextState: stateIssueCreateDescription,
		}, nil
	}

	return tg.CallbackHandlerRes{}, nil
}
