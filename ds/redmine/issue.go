package redmine

import (
	"fmt"

	"git.nixys.ru/apps/nxs-support-bot/misc"
	rdmn "github.com/nixys/nxs-go-redmine/v4"
)

type IssueCreateData struct {
	ProjectID   int64
	TrackerID   int64
	PriorityID  int64
	Subject     string
	Description string
	IsPrivate   bool
	Attachments []AttachmentUpload
}

type IssueUpdateData struct {
	Notes       string
	Attachments []AttachmentUpload
}

// IssueCreate creates new issue in Redmine.
// New issue ID will be returned
func (r *Redmine) IssueCreate(userID int64, issue IssueCreateData) (int64, error) {

	if userID == 0 {
		return 0, misc.ErrUserNotSet
	}

	c, err := r.ctxGetByUserID(userID)
	if err != nil {
		return 0, fmt.Errorf("create redmine issue: %w", err)
	}

	ic, _, err := c.IssueCreate(rdmn.IssueCreateObject{
		ProjectID:   int(issue.ProjectID),
		TrackerID:   int(issue.TrackerID),
		PriorityID:  int(issue.PriorityID),
		Subject:     issue.Subject,
		Description: issue.Description,
		IsPrivate:   issue.IsPrivate,
		Uploads: func() []rdmn.AttachmentUploadObject {
			uploads := []rdmn.AttachmentUploadObject{}
			for _, a := range issue.Attachments {
				uploads = append(uploads, rdmn.AttachmentUploadObject(a))
			}
			return uploads
		}(),
	})
	if err != nil {
		return 0, fmt.Errorf("create redmine issue: %w", err)
	}

	return int64(ic.ID), nil
}

func (r *Redmine) IssueUpdate(userID, issueID int64, d IssueUpdateData) error {

	if userID == 0 {
		return misc.ErrUserNotSet
	}

	c, err := r.ctxGetByUserID(userID)
	if err != nil {
		return err
	}

	_, err = c.IssueUpdate(
		int(issueID),
		rdmn.IssueUpdateObject{
			Notes: d.Notes,
			Uploads: func() []rdmn.AttachmentUploadObject {
				uploads := []rdmn.AttachmentUploadObject{}
				for _, a := range d.Attachments {
					uploads = append(uploads, rdmn.AttachmentUploadObject(a))
				}
				return uploads
			}(),
		},
	)
	if err != nil {
		return fmt.Errorf("update redmine issue: %w", err)
	}

	return nil
}

func (r *Redmine) IssueURL(issueID int64) string {
	return fmt.Sprintf("%s/issues/%d", r.host, issueID)
}

func (r *Redmine) ctxGetByUserID(userID int64) (rdmn.Context, error) {

	var c rdmn.Context

	if userID == 0 {
		return r.c, nil
	}

	u, _, err := r.c.UserSingleGet(
		int(userID),
		rdmn.UserSingleGetRequest{},
	)
	if err != nil {
		return c, fmt.Errorf("create redmine issue: %w", err)
	}

	c.SetEndpoint(r.host)
	c.SetAPIKey(u.APIKey)

	return c, nil
}
