package tgbot

import (
	tg "github.com/nixys/nxs-go-telegram"
	"github.com/nixys/nxs-support-bot/misc"
	"github.com/nixys/nxs-support-bot/modules/localization"
	"github.com/nixys/nxs-support-bot/modules/users"
)

func issueCreateCmd(t *tg.Telegram, sess *tg.Session, cmd string, args string) (tg.CommandHandlerRes, error) {

	c, err := userEnvGet(t, sess)
	if err != nil {
		return tg.CommandHandlerRes{}, err
	}

	switch c.user.Type {
	case users.UserTypeUnauthorized:
		return tg.CommandHandlerRes{
			NextState: stateInitLang,
		}, nil
	case users.UserTypeFeedback:
		return tg.CommandHandlerRes{
			NextState: stateFeedback,
		}, nil
	case users.UserTypeInternal:
		return tg.CommandHandlerRes{
			NextState: stateIssueCreate,
		}, nil
	}

	return tg.CommandHandlerRes{
		NextState: tg.SessStateBreak(),
	}, nil
}

func issueCreateState(t *tg.Telegram, sess *tg.Session) (tg.StateHandlerRes, error) {

	var issue slotIssueCreate

	bCtx := botCtxGet(t)

	c, err := userEnvGet(t, sess)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	b, err := sess.SlotGet(slotNameIssueCreate, &issue)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	if b == false {

		mms, err := bCtx.users.MembershipsGet(sess.UserIDGet())
		if err != nil {
			return tg.StateHandlerRes{}, err
		}

		// TODO: check if membership projects empty
		proj := misc.IDName{}
		if len(mms) > 0 {
			proj = mms[0]
		}

		prio, err := bCtx.c.PriorityGetDefaultLocale(c.l.GetTag())
		if err != nil {
			return tg.StateHandlerRes{}, err
		}

		issue = slotIssueCreate{
			Project:     proj,
			Projects:    mms,
			Memberships: mms,
			Priority: misc.IDName{
				ID:   prio.ID,
				Name: prio.Name,
			},
		}

		if err := sess.SlotSave(slotNameIssueCreate, issue); err != nil {
			return tg.StateHandlerRes{}, err
		}
	}

	m, err := c.l.MessageCreate(
		localization.MsgIssueCreate,
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
					Text:       c.l.BotButton(localization.ButtonProject),
					Identifier: buttonIDProject,
				},
				{
					Text:       c.l.BotButton(localization.ButtonPriority),
					Identifier: buttonIDPriority,
				},
			},
			{
				{
					Text:       c.l.BotButton(localization.ButtonCreateIssue),
					Identifier: buttonIDCreateIssue,
				},
			},
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

func issueCreateCallback(t *tg.Telegram, sess *tg.Session, identifier string) (tg.CallbackHandlerRes, error) {

	switch identifier {
	case buttonIDProject:
		return tg.CallbackHandlerRes{
			NextState: stateIssueCreateProject,
		}, nil
	case buttonIDPriority:
		return tg.CallbackHandlerRes{
			NextState: stateIssueCreatePriority,
		}, nil
	case buttonIDCreateIssue:
		return tg.CallbackHandlerRes{
			NextState: stateIssueCreateSubject,
		}, nil
	case buttonIDCancel:
		return tg.CallbackHandlerRes{
			NextState: stateBye,
		}, nil
	}

	return tg.CallbackHandlerRes{}, nil
}
