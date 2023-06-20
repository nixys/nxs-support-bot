package issues

import (
	"errors"
	"fmt"

	"github.com/nixys/nxs-support-bot/ds/primedb"
	"github.com/nixys/nxs-support-bot/ds/redmine"
	"github.com/nixys/nxs-support-bot/misc"
)

type Settings struct {
	DB       primedb.DB
	Redmine  redmine.Redmine
	Feedback *FeedbackSettings
}

type FeedbackSettings struct {
	ProjectID int64
	UserID    int64
}

type Issues struct {
	d        primedb.DB
	r        redmine.Redmine
	feedback *issuesFeedback
}

type issuesFeedback struct {
	projectID int64
	userID    int64
}

type IssueCreateData struct {
	ProjectID   int64
	TrackerID   int64
	PriorityID  int64
	Subject     string
	Description string
	IsPrivate   bool
	Attachments []AttachmentUpload
}

type IssueFeedbackAddData struct {
	TgUserID    int64
	TgUsername  string
	TgFirstName string
	TgLastName  string
	Notes       string
	Attachments []AttachmentUpload
}

type IssueFeedbackCloseData struct {
	TgUserID      int64
	RdmnUserID    int64
	RdmnLogin     string
	RdmnFirstName string
	RdmnLastName  string
}

type IssueReplyData struct {
	RdmnUserID     int64
	ChatID         int64
	MessageID      int64
	ReplyMessageID int64
	Note           string
	Attachments    []AttachmentUpload
}

type AttachmentUpload redmine.AttachmentUpload
type AttachmentDownload redmine.AttachmentDownload
type UploadData redmine.UploadData

type bunch struct {
	chatID    int64
	messageID int64
	issueID   int64
}

const feedbackIssueAddTpl = `+User details:+
* *Username*: "@{{ .Username }}":https://t.me/{{ .Username }}
* *First name*: {{ .FirstName }}
* *Last name*: {{ .LastName }}

---

{{ .Description }}
`

const feedbackIssueCloseTpl = `⚠️ _User has been authorized in our Customer Care System. You can't any more write him a messages via this issue_

+New user details:+
* *Login*: {{ .Login }}
* *First name*: {{ .FirstName }}
* *Last name*: {{ .LastName }}
`

func Init(s Settings) Issues {

	return Issues{
		d: s.DB,
		r: s.Redmine,
		feedback: func() *issuesFeedback {
			if s.Feedback == nil {
				return nil
			}
			return &issuesFeedback{
				projectID: s.Feedback.ProjectID,
				userID:    s.Feedback.UserID,
			}
		}(),
	}
}

func (iss *Issues) IssueCreate(rdmnUserID int64, d IssueCreateData) (int64, error) {

	if rdmnUserID == 0 {
		return 0, misc.ErrUserNotSet
	}

	issueID, err := iss.r.IssueCreate(
		rdmnUserID,
		redmine.IssueCreateData{
			ProjectID:   d.ProjectID,
			TrackerID:   d.TrackerID,
			PriorityID:  d.PriorityID,
			Subject:     d.Subject,
			Description: d.Description,
			IsPrivate:   d.IsPrivate,
			Attachments: func() []redmine.AttachmentUpload {
				atts := []redmine.AttachmentUpload{}
				for _, a := range d.Attachments {
					atts = append(atts, redmine.AttachmentUpload(a))
				}
				return atts
			}(),
		},
	)
	if err != nil {
		return 0, fmt.Errorf("issue create: %w", err)
	}

	return issueID, nil
}

// IssueFeedbackAdd creates a new feedback issue for user (returns true)
// or add new comment into existing (returns false)
func (iss *Issues) IssueFeedbackAdd(d IssueFeedbackAddData) (bool, error) {

	if iss.feedback == nil {
		return false, nil
	}

	di, err := iss.d.FeedbackIssueGet(d.TgUserID)
	if err != nil {
		if errors.Is(err, misc.ErrNotFound) == true {

			// Feedback issue not exist for the user.
			// Create a new one

			dsc, err := misc.TemplateExec(
				feedbackIssueAddTpl,
				map[string]string{
					"Username":    d.TgUsername,
					"FirstName":   d.TgFirstName,
					"LastName":    d.TgLastName,
					"Description": d.Notes,
				},
			)
			if err != nil {
				return false, fmt.Errorf("issue feedback add: %w", err)
			}

			iID, err := iss.r.IssueCreate(
				iss.feedback.userID,
				redmine.IssueCreateData{
					ProjectID:   iss.feedback.projectID,
					Subject:     fmt.Sprintf("Feedback issue: %s %s (@%s)", d.TgFirstName, d.TgLastName, d.TgUsername),
					Description: dsc,
					Attachments: func() []redmine.AttachmentUpload {
						atts := []redmine.AttachmentUpload{}
						for _, a := range d.Attachments {
							atts = append(atts, redmine.AttachmentUpload(a))
						}
						return atts
					}(),
				},
			)
			if err != nil {
				return false, fmt.Errorf("issue feedback add: %w", err)
			}

			if _, err := iss.d.FeedbackIssueSave(
				primedb.FeedbackIssueInsertData{
					TgID:    d.TgUserID,
					IssueID: iID,
				},
			); err != nil {
				return false, fmt.Errorf("issue feedback add: %w", err)
			}

			return true, nil
		}

		return false, fmt.Errorf("issue feedback add: %w", err)
	}

	// Add new comment into feedback issue for user

	if err := iss.r.IssueUpdate(
		iss.feedback.userID,
		di.IssueID,
		redmine.IssueUpdateData{
			Notes: d.Notes,
			Attachments: func() []redmine.AttachmentUpload {
				atts := []redmine.AttachmentUpload{}
				for _, a := range d.Attachments {
					atts = append(atts, redmine.AttachmentUpload(a))
				}
				return atts
			}(),
		},
	); err != nil {
		return false, fmt.Errorf("issue feedback add: %w", err)
	}

	return false, nil
}

func (iss *Issues) IssueFeedbackClose(d IssueFeedbackCloseData) error {

	if iss.feedback == nil {
		return nil
	}

	di, err := iss.d.FeedbackIssueGet(d.TgUserID)
	if err != nil {
		if errors.Is(err, misc.ErrNotFound) == true {
			return nil
		}
		return fmt.Errorf("issue feedback close: %w", err)
	}

	dsc, err := misc.TemplateExec(
		feedbackIssueCloseTpl,
		map[string]string{
			"Login":     d.RdmnLogin,
			"FirstName": d.RdmnFirstName,
			"LastName":  d.RdmnLastName,
		},
	)
	if err != nil {
		return fmt.Errorf("issue feedback close: %w", err)
	}

	if err := iss.r.IssueUpdate(
		iss.feedback.userID,
		di.IssueID,
		redmine.IssueUpdateData{
			Notes: dsc,
		},
	); err != nil {
		return fmt.Errorf("issue feedback close: %w", err)
	}

	return nil
}

// IssueFeedbackUserGet gets user for specified feedback issue
func (iss *Issues) IssueFeedbackUserGet(issueID int64) (int64, error) {

	i, err := iss.d.FeedbackIssueGetByIssueID(issueID)
	if err != nil {
		return 0, fmt.Errorf("issue feedback user get: %w", err)
	}

	return i.TgID, nil
}

func (iss *Issues) IssueReply(d IssueReplyData) error {

	// Get bunch for source message
	b, err := iss.bunchGet(d.ChatID, d.MessageID)
	if err != nil {
		return fmt.Errorf("issue reply: %w", err)
	}

	// Send new comment to Redmine
	if err := iss.r.IssueUpdate(
		d.RdmnUserID,
		b.issueID,
		redmine.IssueUpdateData{
			Notes: d.Note,
			Attachments: func() []redmine.AttachmentUpload {
				atts := []redmine.AttachmentUpload{}
				for _, a := range d.Attachments {
					atts = append(atts, redmine.AttachmentUpload(a))
				}
				return atts
			}(),
		},
	); err != nil {
		return fmt.Errorf("issue reply: %w", err)
	}

	// Save bunch for new user messagee
	if err := iss.BunchAdd(d.ChatID, d.ReplyMessageID, b.issueID); err != nil {
		return fmt.Errorf("issue reply: %w", err)
	}

	return nil
}

func (iss *Issues) BunchAdd(chatID, msgID, issueID int64) error {

	if _, err := iss.d.IssuesBunchSave(
		primedb.IssuesBunchInsertData{
			ChatID:    chatID,
			MessageID: msgID,
			IssueID:   issueID,
		},
	); err != nil {
		return fmt.Errorf("issue bunch add: %w", err)
	}

	return nil
}

func (iss *Issues) AttachmentsUpload(rdmnUserID int64, uploads []UploadData) ([]AttachmentUpload, error) {

	var (
		ruds []redmine.UploadData
		atts []AttachmentUpload
	)

	if rdmnUserID == 0 {
		return nil, misc.ErrUserNotSet
	}

	for _, u := range uploads {
		ruds = append(ruds, redmine.UploadData(u))
	}

	ratts, err := iss.r.AttachmensUpload(rdmnUserID, ruds)
	if err != nil {
		return nil, fmt.Errorf("issue attachments upload: %w", err)
	}

	for _, a := range ratts {
		atts = append(atts, AttachmentUpload(a))
	}

	return atts, nil
}

func (iss *Issues) AttachmentsDownload(attachments []int64) ([]AttachmentDownload, error) {

	ratts, err := iss.r.AttachmentsDownload(attachments)
	if err != nil {
		return nil, fmt.Errorf("issue bunch get: %w", err)
	}

	atts := []AttachmentDownload{}
	for _, a := range ratts {
		atts = append(atts, AttachmentDownload(a))
	}

	return atts, nil
}

func (iss *Issues) IssueURL(issueID int64) string {
	return iss.r.IssueURL(issueID)
}

func (iss *Issues) bunchGet(chatID, msgID int64) (bunch, error) {

	b, err := iss.d.IssuesBunchGet(chatID, msgID)
	if err != nil {
		return bunch{}, fmt.Errorf("issue bunch get: %w", err)
	}

	return bunch{
		chatID:    b.ChatID,
		messageID: b.MessageID,
		issueID:   b.IssueID,
	}, nil
}
