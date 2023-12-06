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

type issueCreateRx struct {
	IssueCreateDataRx issueCreateDataRx `json:"data" binding:"required"`
}

type issueCreateDataRx struct {
	IssueCreateObjectRx issueCreateObjectRx `json:"issue" binding:"required"`
}

type issueCreateObjectRx struct {
	ID             int64             `json:"id"`
	Project        projectRx         `json:"project"`
	Tracker        misc.IDNameLocale `json:"tracker"`
	Status         misc.IDNameLocale `json:"status"`
	Priority       misc.IDNameLocale `json:"priority"`
	Author         misc.IDName       `json:"author"`
	AssignedTo     misc.IDName       `json:"assigned_to"`
	Subject        string            `json:"subject"`
	Description    string            `json:"description"`
	StartDate      string            `json:"start_date"`
	DueDate        string            `json:"due_date"`
	DoneRatio      int64             `json:"done_ratio"`
	IsPrivate      bool              `json:"is_private"`
	EstimatedHours float64           `json:"estimated_hours"`
	SpentHours     float64           `json:"spent_hours"`
	MentionedUsers []misc.IDName     `json:"mentioned_users"`
	CustomFields   []customFieldRx   `json:"custom_fields"`
	CreatedOn      string            `json:"created_on"`
	UpdatedOn      string            `json:"updated_on"`
	ClosedOn       string            `json:"closed_on"`
	Attachments    []attachmentRx    `json:"attachments"`
	Watchers       []misc.IDName     `json:"watchers"`
}

func RedmineCreated(cc *ctx.Ctx, c *gin.Context) handlers.RouteHandlerResponse {

	rx := issueCreateRx{}

	// Fetch data from query
	if err := c.BindJSON(&rx); err != nil {

		cc.Log.WithFields(logrus.Fields{
			"details": err,
		}).Warn("redmine created issue v2 endpoint")

		return handlers.RouteHandlerResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "can't parse body",
		}
	}

	go func() {
		if err := cc.Rdmnhndlr.IssueCreated(
			rdmnhndlr.CreatedData{
				ID:             rx.IssueCreateDataRx.IssueCreateObjectRx.ID,
				Subject:        rx.IssueCreateDataRx.IssueCreateObjectRx.Subject,
				Description:    rx.IssueCreateDataRx.IssueCreateObjectRx.Description,
				IsPrivateIssue: rx.IssueCreateDataRx.IssueCreateObjectRx.IsPrivate,
				Project: misc.IDName{
					ID:   rx.IssueCreateDataRx.IssueCreateObjectRx.Project.ID,
					Name: rx.IssueCreateDataRx.IssueCreateObjectRx.Project.Name,
				},
				Tracker:        rx.IssueCreateDataRx.IssueCreateObjectRx.Tracker,
				Status:         rx.IssueCreateDataRx.IssueCreateObjectRx.Status,
				Priority:       rx.IssueCreateDataRx.IssueCreateObjectRx.Priority,
				Author:         rx.IssueCreateDataRx.IssueCreateObjectRx.Author,
				AssignedTo:     rx.IssueCreateDataRx.IssueCreateObjectRx.AssignedTo,
				Watchers:       rx.IssueCreateDataRx.IssueCreateObjectRx.Watchers,
				MentionedUsers: rx.IssueCreateDataRx.IssueCreateObjectRx.MentionedUsers,
				Attachments: func() []int64 {
					var atts []int64
					for _, a := range rx.IssueCreateDataRx.IssueCreateObjectRx.Attachments {
						atts = append(atts, int64(a.ID))
					}
					return atts
				}(),
				Members: func() map[int64]rdmnhndlr.PermissionsData {
					members := make(map[int64]rdmnhndlr.PermissionsData)
					for _, m := range rx.IssueCreateDataRx.IssueCreateObjectRx.Project.Members {
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
			}).Warn("redmine created issue v2 endpoint")
		}
	}()

	return handlers.RouteHandlerResponse{
		StatusCode: http.StatusOK,
		Message:    "success",
	}
}
