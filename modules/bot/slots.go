package tgbot

import (
	"github.com/nixys/nxs-support-bot/misc"
	"github.com/nixys/nxs-support-bot/modules/issues"
)

const (
	slotNameLang        = "lang"
	slotNameUser        = "users"
	slotNameIssueCreate = "issueCreate"
)

type slotIssueCreate struct {
	Project      misc.IDName
	Projects     []misc.IDName
	Priority     misc.IDName
	Memberships  []misc.IDName
	ProjectsPage int64
	Subject      string
	Description  string
	IssueID      int64
	Attachments  []issues.AttachmentUpload
}
