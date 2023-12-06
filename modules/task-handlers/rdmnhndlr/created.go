package rdmnhndlr

import (
	"fmt"

	"github.com/nixys/nxs-support-bot/misc"
	"github.com/nixys/nxs-support-bot/modules/localization"
)

type CreatedData struct {
	ID             int64
	Subject        string
	Description    string
	IsPrivateIssue bool
	Project        misc.IDName
	Tracker        misc.IDNameLocale
	Status         misc.IDNameLocale
	Priority       misc.IDNameLocale
	Author         misc.IDName
	AssignedTo     misc.IDName
	MentionedUsers []misc.IDName
	Watchers       []misc.IDName
	Attachments    []int64
	Members        map[int64]PermissionsData
}

func (rh *RdmnHndlr) IssueCreated(data CreatedData) error {

	// Prepare data to send users
	sd, err := rh.sendMessagesPrep(
		sendMessagesPrepData{
			users: func() []misc.IDName {
				elts := []misc.IDName{}
				elts = append(elts, data.AssignedTo)
				elts = append(elts, data.Watchers...)
				elts = append(elts, data.MentionedUsers...)
				return elts
			}(),
			author: data.Author,
			permissions: permissionsCheck{
				members: data.Members,
			},
			formatter: issueCreatedMessage,
			data:      data,
		},
	)
	if err != nil {
		return fmt.Errorf("redmine created handler: %w", err)
	}

	// Send data to users
	if err := rh.send(data.ID, sd, data.Attachments); err != nil {
		return fmt.Errorf("redmine created handler: %w", err)
	}

	return nil
}

func issueCreatedMessage(rh *RdmnHndlr, lang string, data any) (string, error) {

	d, b := data.(CreatedData)
	if b == false {
		return "", misc.ErrMalformedData
	}

	l, err := rh.lb.LangSwitch(lang)
	if err != nil {
		return "", fmt.Errorf("message prep created: %w", err)
	}

	m, err := l.MessageCreate(
		localization.MsgIssueCreated,
		map[string]any{
			"Project":        d.Project.Name,
			"IssueID":        d.ID,
			"IssueSubject":   d.Subject,
			"IssueURL":       rh.iss.IssueURL(d.ID),
			"Author":         d.Author.Name,
			"Status":         d.Status.ValueGet(lang),
			"Priority":       d.Priority.ValueGet(lang),
			"AssignedTo":     d.AssignedTo.Name,
			"IsPrivateIssue": d.IsPrivateIssue,
			"Description":    d.Description,
		})
	if err != nil {
		return "", fmt.Errorf("send message prep created: %w", err)
	}

	return m, nil
}
