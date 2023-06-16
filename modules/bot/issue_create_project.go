package tgbot

import (
	"regexp"
	"strconv"
	"strings"

	"git.nixys.ru/apps/nxs-support-bot/misc"
	"git.nixys.ru/apps/nxs-support-bot/modules/localization"
	tg "github.com/nixys/nxs-go-telegram"
)

func issueCreateProjectState(t *tg.Telegram, sess *tg.Session) (tg.StateHandlerRes, error) {

	var issue slotIssueCreate

	c, err := userEnvGet(t, sess)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	if _, err := sess.SlotGet(slotNameIssueCreate, &issue); err != nil {
		return tg.StateHandlerRes{}, err
	}

	m, err := c.l.MessageCreate(
		localization.MsgIssueCreateProject,
		map[string]string{
			"Project":  issue.Project.Name,
			"Priority": issue.Priority.Name,
		},
	)
	if err != nil {
		return tg.StateHandlerRes{}, err
	}

	// Paginate projects
	projs, isPrev, isNext := misc.IDNamePaginate(issue.Projects, issue.ProjectsPage, misc.OnPageDefault)

	// Render project buttons
	buttons := [][]tg.Button{}
	for _, p := range projs {
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
		renderCtrlButtonsPage(c.l, isPrev, isNext),
	)

	return tg.StateHandlerRes{
		Message:               m,
		ParseMode:             tg.ParseModeHTML,
		DisableWebPagePreview: true,
		Buttons:               buttons,
		StickMessage:          true,
	}, nil
}

func issueCreateProjectCallback(t *tg.Telegram, sess *tg.Session, identifier string) (tg.CallbackHandlerRes, error) {

	var issue slotIssueCreate

	switch identifier {

	case buttonIDPrevPage, buttonIDNextPage:

		if _, err := sess.SlotGet(slotNameIssueCreate, &issue); err != nil {
			return tg.CallbackHandlerRes{}, err
		}

		if identifier == buttonIDPrevPage {
			issue.ProjectsPage--
		} else {
			issue.ProjectsPage++
		}

		if err := sess.SlotSave(slotNameIssueCreate, issue); err != nil {
			return tg.CallbackHandlerRes{}, err
		}

		return tg.CallbackHandlerRes{
			NextState: stateIssueCreateProject,
		}, nil

	case buttonIDBack:
		return tg.CallbackHandlerRes{
			NextState: stateIssueCreate,
		}, nil

	default:

		if _, err := sess.SlotGet(slotNameIssueCreate, &issue); err != nil {
			return tg.CallbackHandlerRes{}, err
		}

		pID, err := strconv.ParseInt(identifier, 10, 64)
		if err != nil {
			return tg.CallbackHandlerRes{}, err
		}

		issue.Project = func() misc.IDName {
			for _, p := range issue.Memberships {
				if p.ID == pID {
					return p
				}
			}
			return misc.IDName{}
		}()

		if err := sess.SlotSave(slotNameIssueCreate, issue); err != nil {
			return tg.CallbackHandlerRes{}, err
		}

		return tg.CallbackHandlerRes{
			NextState: stateIssueCreate,
		}, nil
	}
}

func issueCreateProjectMsg(t *tg.Telegram, sess *tg.Session) (tg.MessageHandlerRes, error) {

	var issue slotIssueCreate

	if _, err := sess.SlotGet(slotNameIssueCreate, &issue); err != nil {
		return tg.MessageHandlerRes{}, err
	}

	regex := strings.Join(sess.UpdateChain().MessageTextGet(), "-")

	if regex == "*" {
		regex = ".*"
	}

	r, err := regexp.Compile("(?i)" + regex)
	if err != nil {
		return tg.MessageHandlerRes{}, err
	}

	// Filter projects
	issue.Projects = []misc.IDName{}
	for _, m := range issue.Memberships {
		if r.MatchString(m.Name) == true {
			issue.Projects = append(issue.Projects, m)
		}
	}

	// Reset page
	issue.ProjectsPage = 0

	if err = sess.SlotSave(slotNameIssueCreate, issue); err != nil {
		return tg.MessageHandlerRes{}, err
	}

	return tg.MessageHandlerRes{
		NextState: stateIssueCreateProject,
	}, nil
}
