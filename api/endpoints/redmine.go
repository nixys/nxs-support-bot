package endpoints

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nixys/nxs-support-bot/ctx"
	"github.com/nixys/nxs-support-bot/misc"
	"github.com/nixys/nxs-support-bot/modules/task-handlers/rdmnhndlr"
	"github.com/sirupsen/logrus"
)

type issueActionRx struct {
	Action string `json:"action"`
}

// Issue Create Rx

type issueCreateRx struct {
	IssueCreateDataRx issueCreateDataRx `json:"data" binding:"required"`
}

type issueCreateDataRx struct {
	IssueCreateObjectRx issueCreateObjectRx `json:"issue" binding:"required"`
}

type issueCreateObjectRx struct {
	ID             int64           `json:"id"`
	Project        projectRx       `json:"project"`
	Tracker        misc.IDName     `json:"tracker"`
	Status         misc.IDName     `json:"status"`
	Priority       misc.IDName     `json:"priority"`
	Author         misc.IDName     `json:"author"`
	AssignedTo     misc.IDName     `json:"assigned_to"`
	Subject        string          `json:"subject"`
	Description    string          `json:"description"`
	StartDate      string          `json:"start_date"`
	DueDate        string          `json:"due_date"`
	DoneRatio      int64           `json:"done_ratio"`
	IsPrivate      bool            `json:"is_private"`
	EstimatedHours float64         `json:"estimated_hours"`
	SpentHours     float64         `json:"spent_hours"`
	MentionedUsers []misc.IDName   `json:"mentioned_users"`
	CustomFields   []customFieldRx `json:"custom_fields"`
	CreatedOn      string          `json:"created_on"`
	UpdatedOn      string          `json:"updated_on"`
	ClosedOn       string          `json:"closed_on"`
	Attachments    []attachmentRx  `json:"attachments"`
	Watchers       []misc.IDName   `json:"watchers"`
}

// Issue Update Rx

type issueUpdateRx struct {
	IssueUpdateDataRx issueUpdateDataRx `json:"data"`
}

type issueUpdateDataRx struct {
	IssueUpdateObjectRx issueUpdateObjectRx `json:"issue"`
}

type issueUpdateObjectRx struct {
	ID             int64          `json:"id"`
	Project        projectRx      `json:"project"`
	Tracker        misc.IDName    `json:"tracker"`
	Status         misc.IDName    `json:"status"`
	Priority       misc.IDName    `json:"priority"`
	Author         misc.IDName    `json:"author"`
	AssignedTo     misc.IDName    `json:"assigned_to"`
	Category       misc.IDName    `json:"category"`
	Subject        string         `json:"subject"`
	Description    string         `json:"description"`
	StartDate      string         `json:"start_date"`
	DueDate        string         `json:"due_date"`
	DoneRatio      int64          `json:"done_ratio"`
	IsPrivate      bool           `json:"is_private"`
	EstimatedHours float64        `json:"estimated_hours"`
	SpentHours     float64        `json:"spent_hours"`
	MentionedUsers []misc.IDName  `json:"mentioned_users"`
	CreatedOn      string         `json:"created_on"`
	UpdatedOn      string         `json:"updated_on"`
	ClosedOn       string         `json:"closed_on"`
	Attachments    []attachmentRx `json:"attachments"`
	Watchers       []misc.IDName  `json:"watchers"`
	Journals       []journalRx    `json:"journals"`
}

// Common Rx

type projectRx struct {
	ID      int64      `json:"id"`
	Name    string     `json:"name"`
	Members []memberRx `json:"members"`
}

type memberRx struct {
	ID     int64     `json:"id"`
	Name   string    `json:"name"`
	Roles  []rolesRx `json:"roles"`
	Access accessRx  `json:"access"`
}

type rolesRx struct {
	ID          int64        `json:"id"`
	Name        string       `json:"name"`
	Permissions permissionRx `json:"permissions"`
}

type accessRx struct {
	ViewCurrentIssue bool `json:"view_current_issue"`
	ViewPrivateNotes bool `json:"view_private_notes"`
}

type permissionRx struct {
	IssuesVisibility string `json:"issues_visibility"`
	ViewPrivateNotes bool   `json:"view_private_notes"`
}

type customFieldRx struct {
	ID       int64       `json:"id"`
	Name     string      `json:"name"`
	Multiple bool        `json:"multiple"`
	Value    interface{} `json:"value"`
}

type attachmentRx struct {
	ID          int         `json:"id"`
	Filename    string      `json:"filename"`
	Filesize    int64       `json:"filesize"`
	ContentType string      `json:"content_type"`
	Description string      `json:"description"`
	Author      misc.IDName `json:"author"`
	CreatedOn   string      `json:"created_on"`
}

type journalRx struct {
	ID             int64         `json:"id"`
	User           misc.IDName   `json:"user"`
	Notes          string        `json:"notes"`
	PrivateNotes   bool          `json:"private_notes"`
	CreatedOn      string        `json:"created_on"`
	Details        []detailRx    `json:"details"`
	MentionedUsers []misc.IDName `json:"mentioned_users"`
}

type detailRx struct {
	Property string `json:"property"`
	Name     string `json:"name"`
	OldValue string `json:"old_value"`
	NewValue string `json:"new_value"`
}

func Redmine(cc *ctx.Ctx, c *gin.Context) RouteHandlerResponse {

	rx := issueActionRx{}

	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return RouteHandlerResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		}
	}

	if err := json.Unmarshal(b, &rx); err != nil {
		cc.Log.WithFields(logrus.Fields{
			"details": err,
		}).Warn("redmine endpoint")
		return RouteHandlerResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "incorrect body json",
		}
	}

	switch rx.Action {
	case "issue_create":

		data := issueCreateRx{}

		if err := json.Unmarshal(b, &data); err != nil {
			return RouteHandlerResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "incorrect body json",
			}
		}

		go func() {
			if err := cc.Rdmnhndlr.IssueCreated(
				rdmnhndlr.CreatedData{
					ID:             data.IssueCreateDataRx.IssueCreateObjectRx.ID,
					Subject:        data.IssueCreateDataRx.IssueCreateObjectRx.Subject,
					Description:    data.IssueCreateDataRx.IssueCreateObjectRx.Description,
					IsPrivateIssue: data.IssueCreateDataRx.IssueCreateObjectRx.IsPrivate,
					Project: misc.IDName{
						ID:   data.IssueCreateDataRx.IssueCreateObjectRx.Project.ID,
						Name: data.IssueCreateDataRx.IssueCreateObjectRx.Project.Name,
					},
					Tracker:        data.IssueCreateDataRx.IssueCreateObjectRx.Tracker,
					Status:         data.IssueCreateDataRx.IssueCreateObjectRx.Status,
					Priority:       data.IssueCreateDataRx.IssueCreateObjectRx.Priority,
					Author:         data.IssueCreateDataRx.IssueCreateObjectRx.Author,
					AssignedTo:     data.IssueCreateDataRx.IssueCreateObjectRx.AssignedTo,
					Watchers:       data.IssueCreateDataRx.IssueCreateObjectRx.Watchers,
					MentionedUsers: data.IssueCreateDataRx.IssueCreateObjectRx.MentionedUsers,
					Attachments: func() []int64 {
						var atts []int64
						for _, a := range data.IssueCreateDataRx.IssueCreateObjectRx.Attachments {
							atts = append(atts, int64(a.ID))
						}
						return atts
					}(),
					Members: func() map[int64]rdmnhndlr.PermissionsData {
						members := make(map[int64]rdmnhndlr.PermissionsData)
						for _, m := range data.IssueCreateDataRx.IssueCreateObjectRx.Project.Members {
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
				}).Warn("redmine created issue endpoint")
			}
		}()

		return RouteHandlerResponse{
			StatusCode: http.StatusOK,
			Message:    "success",
		}
	case "issue_edit":

		data := issueUpdateRx{}

		if err := json.Unmarshal(b, &data); err != nil {
			cc.Log.WithFields(logrus.Fields{
				"details": err,
			}).Warn("redmine updated issue endpoint")
			return RouteHandlerResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "incorrect body json",
			}
		}

		go func() {
			if err := cc.Rdmnhndlr.IssueUpdated(
				rdmnhndlr.UpdatedData{
					ID:             data.IssueUpdateDataRx.IssueUpdateObjectRx.ID,
					Subject:        data.IssueUpdateDataRx.IssueUpdateObjectRx.Subject,
					Description:    data.IssueUpdateDataRx.IssueUpdateObjectRx.Description,
					IsPrivateIssue: data.IssueUpdateDataRx.IssueUpdateObjectRx.IsPrivate,
					Project: misc.IDName{
						ID:   data.IssueUpdateDataRx.IssueUpdateObjectRx.Project.ID,
						Name: data.IssueUpdateDataRx.IssueUpdateObjectRx.Project.Name,
					},
					Tracker:    data.IssueUpdateDataRx.IssueUpdateObjectRx.Tracker,
					Status:     data.IssueUpdateDataRx.IssueUpdateObjectRx.Status,
					Priority:   data.IssueUpdateDataRx.IssueUpdateObjectRx.Priority,
					Category:   data.IssueUpdateDataRx.IssueUpdateObjectRx.Category,
					Author:     data.IssueUpdateDataRx.IssueUpdateObjectRx.Author,
					AssignedTo: data.IssueUpdateDataRx.IssueUpdateObjectRx.AssignedTo,
					Watchers:   data.IssueUpdateDataRx.IssueUpdateObjectRx.Watchers,
					MentionedUsers: func() []misc.IDName {
						var users []misc.IDName
						users = append(users, data.IssueUpdateDataRx.IssueUpdateObjectRx.MentionedUsers...)
						for _, j := range data.IssueUpdateDataRx.IssueUpdateObjectRx.Journals {
							users = append(users, j.MentionedUsers...)
						}
						return users
					}(),
					Attachments: func() []int64 {
						var atts []int64
						for _, a := range data.IssueUpdateDataRx.IssueUpdateObjectRx.Attachments {
							atts = append(atts, int64(a.ID))
						}
						return atts
					}(),
					Journals: func() []rdmnhndlr.UpdatedJournalData {
						journals := []rdmnhndlr.UpdatedJournalData{}
						for _, j := range data.IssueUpdateDataRx.IssueUpdateObjectRx.Journals {
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
						for _, m := range data.IssueUpdateDataRx.IssueUpdateObjectRx.Project.Members {
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
				}).Warn("redmine updated issue endpoint")
			}
		}()

		return RouteHandlerResponse{
			StatusCode: http.StatusOK,
			Message:    "success",
		}
	}

	return RouteHandlerResponse{
		StatusCode: http.StatusBadRequest,
		Message:    "unknown issue action",
	}
}
