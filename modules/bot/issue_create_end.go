package tgbot

import (
	"strconv"

	tg "github.com/nixys/nxs-go-telegram"
	"github.com/nixys/nxs-support-bot/misc"
	"github.com/nixys/nxs-support-bot/modules/localization"
)

func issueCreateEndState(t *tg.Telegram, sess *tg.Session) (tg.StateHandlerRes, error) {

	var issue slotIssueCreate

	bCtx := botCtxGet(t)

	c, err := userEnvGet(t, sess)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	if _, err := sess.SlotGet(slotNameIssueCreate, &issue); err != nil {
		return tg.StateHandlerRes{}, err
	}

	m, err := c.l.MessageCreate(
		localization.MsgIssueCreateEnd,
		map[string]string{
			"Project":      issue.Project.Name,
			"IssueURL":     bCtx.issues.IssueURL(issue.IssueID),
			"IssueID":      strconv.FormatInt(issue.IssueID, 10),
			"IssueSubject": issue.Subject,
		},
	)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	return tg.StateHandlerRes{
		Message:               m,
		ParseMode:             tg.ParseModeHTML,
		DisableWebPagePreview: true,
		StickMessage:          true,
		NextState:             tg.SessStateDestroy(),
	}, nil
}

func issueCreateEndSent(t *tg.Telegram, sess *tg.Session, messages []tg.MessageSent) error {

	var issue slotIssueCreate

	if len(messages) == 0 {
		return misc.ErrZeroLen
	}

	bCtx := botCtxGet(t)

	if _, err := sess.SlotGet(slotNameIssueCreate, &issue); err != nil {
		return err
	}

	// Using only first message
	mID := int64(messages[0].MessageID)

	return bCtx.issues.BunchAdd(sess.ChatIDGet(), mID, issue.IssueID)
}
