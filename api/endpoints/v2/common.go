package endpointsv2

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nixys/nxs-support-bot/misc"
)

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
	ID       int64             `json:"id"`
	Name     misc.IDNameLocale `json:"name"`
	Multiple bool              `json:"multiple"`
	Value    interface{}       `json:"value"`
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

func AuthorizeRedmine(secretToken string) gin.HandlerFunc {

	return func(c *gin.Context) {

		auth := c.GetHeader("Authorization")
		if len(auth) == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		if token == auth {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if token != secretToken {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
