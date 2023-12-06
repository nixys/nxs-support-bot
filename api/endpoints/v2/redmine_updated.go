package endpointsv2

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nixys/nxs-support-bot/api/handlers"
	"github.com/nixys/nxs-support-bot/ctx"
	"github.com/nixys/nxs-support-bot/misc"
	"github.com/nixys/nxs-support-bot/modules/task-handlers/rdmnhndlr"
	"github.com/sirupsen/logrus"
)

type issueUpdateRx struct {
	IssueUpdateDataRx issueUpdateDataRx `json:"data"`
}

type issueUpdateDataRx struct {
	IssueUpdateObjectRx issueUpdateObjectRx `json:"issue"`
}

type issueUpdateObjectRx struct {
	ID             int64             `json:"id"`
	Project        projectRx         `json:"project"`
	Tracker        misc.IDNameLocale `json:"tracker"`
	Status         misc.IDNameLocale `json:"status"`
	Priority       misc.IDNameLocale `json:"priority"`
	Author         misc.IDName       `json:"author"`
	AssignedTo     misc.IDName       `json:"assigned_to"`
	Category       misc.IDNameLocale `json:"category"`
	Subject        string            `json:"subject"`
	Description    string            `json:"description"`
	StartDate      string            `json:"start_date"`
	DueDate        string            `json:"due_date"`
	DoneRatio      int64             `json:"done_ratio"`
	IsPrivate      bool              `json:"is_private"`
	EstimatedHours float64           `json:"estimated_hours"`
	SpentHours     float64           `json:"spent_hours"`
	MentionedUsers []misc.IDName     `json:"mentioned_users"`
	CreatedOn      string            `json:"created_on"`
	UpdatedOn      string            `json:"updated_on"`
	ClosedOn       string            `json:"closed_on"`
	Attachments    []attachmentRx    `json:"attachments"`
	Watchers       []misc.IDName     `json:"watchers"`
	Journals       []journalRx       `json:"journals"`
}

func RedmineUpdated(cc *ctx.Ctx, c *gin.Context) handlers.RouteHandlerResponse {

	rx := issueUpdateRx{}

	// Fetch data from query
	if err := c.BindJSON(&rx); err != nil {

		cc.Log.WithFields(logrus.Fields{
			"details": err,
		}).Warn("redmine updated issue v2 endpoint")

		return handlers.RouteHandlerResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "can't parse body",
		}
	}

	go func() {
		if err := cc.Rdmnhndlr.IssueUpdated(
			rdmnhndlr.UpdatedData{
				ID:             rx.IssueUpdateDataRx.IssueUpdateObjectRx.ID,
				Subject:        rx.IssueUpdateDataRx.IssueUpdateObjectRx.Subject,
				Description:    rx.IssueUpdateDataRx.IssueUpdateObjectRx.Description,
				IsPrivateIssue: rx.IssueUpdateDataRx.IssueUpdateObjectRx.IsPrivate,
				Project: misc.IDName{
					ID:   rx.IssueUpdateDataRx.IssueUpdateObjectRx.Project.ID,
					Name: rx.IssueUpdateDataRx.IssueUpdateObjectRx.Project.Name,
				},
				Tracker:    rx.IssueUpdateDataRx.IssueUpdateObjectRx.Tracker,
				Status:     rx.IssueUpdateDataRx.IssueUpdateObjectRx.Status,
				Priority:   rx.IssueUpdateDataRx.IssueUpdateObjectRx.Priority,
				Category:   rx.IssueUpdateDataRx.IssueUpdateObjectRx.Category,
				Author:     rx.IssueUpdateDataRx.IssueUpdateObjectRx.Author,
				AssignedTo: rx.IssueUpdateDataRx.IssueUpdateObjectRx.AssignedTo,
				Watchers:   rx.IssueUpdateDataRx.IssueUpdateObjectRx.Watchers,
				MentionedUsers: func() []misc.IDName {
					var users []misc.IDName
					users = append(users, rx.IssueUpdateDataRx.IssueUpdateObjectRx.MentionedUsers...)
					for _, j := range rx.IssueUpdateDataRx.IssueUpdateObjectRx.Journals {
						users = append(users, j.MentionedUsers...)
					}
					return users
				}(),
				Attachments: func() []int64 {
					var atts []int64
					for _, a := range rx.IssueUpdateDataRx.IssueUpdateObjectRx.Attachments {
						atts = append(atts, int64(a.ID))
					}
					return atts
				}(),
				Journals: func() []rdmnhndlr.UpdatedJournalData {
					journals := []rdmnhndlr.UpdatedJournalData{}
					for _, j := range rx.IssueUpdateDataRx.IssueUpdateObjectRx.Journals {
						journals = append(
							journals,
							rdmnhndlr.UpdatedJournalData{
								User:           j.User,
								Notes:          j.Notes,
								IsPrivateNotes: j.PrivateNotes,
								Details: func() []rdmnhndlr.UpdatedJournalDetailData {
									details := []rdmnhndlr.UpdatedJournalDetailData{}
									for _, d := range j.Details {
										details = append(
											details,
											rdmnhndlr.UpdatedJournalDetailData{
												Property: d.Property,
												Name:     d.Name,
												OldValue: d.OldValue,
												NewValue: d.NewValue,
											},
										)
									}
									return details
								}(),
							},
						)
					}
					return journals
				}(),
				Members: func() map[int64]rdmnhndlr.PermissionsData {
					members := make(map[int64]rdmnhndlr.PermissionsData)
					for _, m := range rx.IssueUpdateDataRx.IssueUpdateObjectRx.Project.Members {
						members[m.ID] = rdmnhndlr.PermissionsData{
							ViewCurrentIssue: m.Access.ViewCurrentIssue,
							ViewPrivateNotes: m.Access.ViewPrivateNotes,
						}
					}
					return members
				}(),
			},
		); err != nil {
			cc.Log.WithFields(logrus.Fields{
				"details": err,
			}).Warn("redmine updated issue v2 endpoint")
		}
	}()

	return handlers.RouteHandlerResponse{
		StatusCode: http.StatusOK,
		Message:    "success",
	}
}
