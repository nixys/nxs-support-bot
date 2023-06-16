package tgbot

import (
	"strconv"

	"git.nixys.ru/apps/nxs-support-bot/misc"
	"git.nixys.ru/apps/nxs-support-bot/modules/localization"
	tg "github.com/nixys/nxs-go-telegram"
)

func issueCreatePriorityState(t *tg.Telegram, sess *tg.Session) (tg.StateHandlerRes, error) {

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
		localization.MsgIssueCreatePriority,
		map[string]string{
			"Project":  issue.Project.Name,
			"Priority": issue.Priority.Name,
		},
	)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	priorities, err := bCtx.c.PrioritiesGet()
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	buttons := [][]tg.Button{}
	for _, p := range priorities {
		buttons = append(
			buttons,
			[]tg.Button{
				{
					Text:       p.Name,
					Identifier: strconv.FormatInt(p.ID, 10),
				},
			},
		)
	}

	// Render control buttons
	buttons = append(
		buttons,
		[]tg.Button{
			{
				Text:       c.l.BotButton(localization.ButtonBack),
				Identifier: buttonIDBack,
			},
		},
	)

	return tg.StateHandlerRes{
		Message:               m,
		ParseMode:             tg.ParseModeHTML,
		DisableWebPagePreview: true,
		Buttons:               buttons,
		StickMessage:          true,
	}, nil
}

func issueCreatePriorityCallback(t *tg.Telegram, sess *tg.Session, identifier string) (tg.CallbackHandlerRes, error) {

	var issue slotIssueCreate

	switch identifier {

	case buttonIDBack:
		return tg.CallbackHandlerRes{
			NextState: stateIssueCreate,
		}, nil

	default:

		bCtx := botCtxGet(t)

		if _, err := sess.SlotGet(slotNameIssueCreate, &issue); err != nil {
			return tg.CallbackHandlerRes{}, err
		}

		pID, err := strconv.ParseInt(identifier, 10, 64)
		if err != nil {
			return tg.CallbackHandlerRes{}, err
		}

		prio, err := bCtx.c.PriorityGetByID(pID)
		if err != nil {
			return tg.CallbackHandlerRes{}, err
		}

		issue.Priority = misc.IDName{
			ID:   prio.ID,
			Name: prio.Name,
		}

		if err := sess.SlotSave(slotNameIssueCreate, issue); err != nil {
			return tg.CallbackHandlerRes{}, err
		}

		return tg.CallbackHandlerRes{
			NextState: stateIssueCreate,
		}, nil
	}
}
