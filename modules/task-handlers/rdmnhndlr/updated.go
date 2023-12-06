package rdmnhndlr

import (
	"fmt"
	"strconv"

	"github.com/nixys/nxs-support-bot/misc"
	"github.com/nixys/nxs-support-bot/modules/localization"
)

type UpdatedData struct {
	ID             int64
	Subject        string
	Description    string
	IsPrivateIssue bool
	Project        misc.IDName
	Tracker        misc.IDNameLocale
	Category       misc.IDNameLocale
	Status         misc.IDNameLocale
	Priority       misc.IDNameLocale
	Author         misc.IDName
	AssignedTo     misc.IDName
	MentionedUsers []misc.IDName
	Watchers       []misc.IDName
	Attachments    []int64
	Journals       []UpdatedJournalData
	Members        map[int64]PermissionsData
}

type UpdatedJournalData struct {
	User           misc.IDName
	Notes          string
	IsPrivateNotes bool
	Details        []UpdatedJournalDetailData
}

type UpdatedJournalDetailData struct {
	Property string
	Name     string
	OldValue string
	NewValue string
}

type renderData struct {
	ID             int64
	Subject        string
	Project        misc.IDName
	Author         misc.IDName
	IsPrivateIssue bool
	Notes          string
	IsPrivateNotes bool
	Description    *string
	AssignedTo     *string
	Status         *misc.IDNameLocale
	Priority       *misc.IDNameLocale
	Tracker        *misc.IDNameLocale
	Category       *misc.IDNameLocale
	Attachments    []int64
}

func (rh *RdmnHndlr) IssueUpdated(data UpdatedData) error {

	var pu *feedbackUser

	if len(data.Journals) == 0 {
		return fmt.Errorf("redmine updated handler: %w", misc.ErrMalformedData)
	}

	journal := data.Journals[0]

	rd := renderData{
		ID:             data.ID,
		Subject:        data.Subject,
		Project:        data.Project,
		Author:         journal.User,
		IsPrivateIssue: data.IsPrivateIssue,
		Notes:          journal.Notes,
		IsPrivateNotes: journal.IsPrivateNotes,
	}

	isUpdated := false
	if len(journal.Notes) > 0 {
		isUpdated = true
	}

	for _, d := range journal.Details {
		switch d.Property {
		case "attr":
			switch d.Name {
			case "assigned_to_id":
				isUpdated = true
				rd.AssignedTo = &data.AssignedTo.Name
			case "status_id":
				isUpdated = true
				rd.Status = &data.Status
			case "priority_id":
				isUpdated = true
				rd.Priority = &data.Priority
			case "tracker_id":
				isUpdated = true
				rd.Tracker = &data.Tracker
			case "category_id":
				isUpdated = true
				rd.Category = &data.Category
			case "description":
				isUpdated = true
				s := d.NewValue
				rd.Description = &s
			}
		case "attachment":
			isUpdated = true
			att, err := strconv.ParseInt(d.Name, 10, 64)
			if err != nil {
				return fmt.Errorf("redmine updated handler: %w", err)
			}
			rd.Attachments = append(rd.Attachments, att)
		}
	}

	// If necessary parameters does not changed in update
	if isUpdated == false {
		return nil
	}

	// Prepare data for feedback user (if necessary).
	// Check update conditions for feedback user
	if rd.IsPrivateNotes == false && (len(rd.Notes) > 0 || len(rd.Attachments) > 0) {

		u, err := rh.feedbackUserGet(data.ID)
		if err != nil {
			return fmt.Errorf("redmine updated handler: %w", err)
		}

		pu = u
	}

	// Prepare data to send users
	sd, err := rh.sendMessagesPrep(
		sendMessagesPrepData{
			users: func() []misc.IDName {
				elts := []misc.IDName{}
				elts = append(elts, data.Author)
				elts = append(elts, data.AssignedTo)
				elts = append(elts, data.Watchers...)
				elts = append(elts, data.MentionedUsers...)
				return elts
			}(),
			author: rd.Author,
			permissions: permissionsCheck{
				members:        data.Members,
				isPrivateNotes: rd.IsPrivateNotes,
			},
			formatter:         issueUpdatedMessage,
			formatterFeedback: issueUpdatedFeedbackMessage,
			data:              rd,
			pu:                pu,
		},
	)
	if err != nil {
		return fmt.Errorf("redmine updated handler: %w", err)
	}

	// Send data to users
	if err := rh.send(data.ID, sd, rd.Attachments); err != nil {
		return fmt.Errorf("redmine updated handler: %w", err)
	}

	return nil
}

func issueUpdatedMessage(rh *RdmnHndlr, lang string, data any) (string, error) {

	d, b := data.(renderData)
	if b == false {
		return "", misc.ErrMalformedData
	}

	l, err := rh.lb.LangSwitch(lang)
	if err != nil {
		return "", fmt.Errorf("message prep updated: %w", err)
	}

	m, err := l.MessageCreate(
		localization.MsgIssueUpdated,
		map[string]any{
			"Project":        d.Project.Name,
			"IssueID":        d.ID,
			"IssueSubject":   d.Subject,
			"IssueURL":       rh.iss.IssueURL(d.ID),
			"Author":         d.Author.Name,
			"IsPrivateIssue": d.IsPrivateIssue,
			"AssignedTo":     d.AssignedTo,
			"Status":         d.Status.ValueGet(lang),
			"Tracker":        d.Tracker.ValueGet(lang),
			"Category":       d.Category.ValueGet(lang),
			"Notes":          d.Notes,
			"IsPrivateNotes": d.IsPrivateNotes,
			"Description":    d.Description,
		})
	if err != nil {
		return "", fmt.Errorf("send message prep updated: %w", err)
	}

	return m, nil
}

func issueUpdatedFeedbackMessage(rh *RdmnHndlr, lang string, data any) (string, error) {

	d, b := data.(renderData)
	if b == false {
		return "", misc.ErrMalformedData
	}

	l, err := rh.lb.LangSwitch(lang)
	if err != nil {
		return "", fmt.Errorf("message feedback prep updated: %w", err)
	}

	m, err := l.MessageCreate(
		localization.MsgIssueUpdated,
		map[string]any{
			"Project":        d.Project.Name,
			"IssueID":        d.ID,
			"IssueSubject":   d.Subject,
			"IssueURL":       rh.iss.IssueURL(d.ID),
			"Author":         d.Author.Name,
			"IsPrivateIssue": d.IsPrivateIssue,
			"AssignedTo":     d.AssignedTo,
			"Status":         d.Status.ValueGet(lang),
			"Priority":       d.Priority.ValueGet(lang),
			"Tracker":        d.Tracker.ValueGet(lang),
			"Category":       d.Category.ValueGet(lang),
			"Notes":          d.Notes,
			"IsPrivateNotes": d.IsPrivateNotes,
			"Description":    d.Description,
		})
	if err != nil {
		return "", fmt.Errorf("send message prep updated: %w", err)
	}

	return m, nil
}
