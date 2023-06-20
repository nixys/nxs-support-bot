package redmine

import (
	"fmt"
	"io"

	rdmn "github.com/nixys/nxs-go-redmine/v4"
	"github.com/nixys/nxs-support-bot/misc"
)

type UploadData struct {
	File io.Reader
	Name string
}

type AttachmentUpload rdmn.AttachmentUploadObject

type AttachmentDownload struct {
	Reader      io.Reader
	Name        string
	Caption     string
	ContentType string
}

func (r *Redmine) AttachmensUpload(userID int64, uploads []UploadData) ([]AttachmentUpload, error) {

	var atts []AttachmentUpload

	if userID == 0 {
		return nil, misc.ErrUserNotSet
	}

	c, err := r.ctxGetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("redmine attachments upload: %w", err)
	}

	for _, u := range uploads {

		a, _, err := c.AttachmentUploadStream(u.File, u.Name)
		if err != nil {
			return nil, fmt.Errorf("redmine attachments upload: %w", err)
		}

		atts = append(atts, AttachmentUpload(a))
	}

	return atts, nil
}

func (r *Redmine) AttachmentsDownload(downloads []int64) ([]AttachmentDownload, error) {

	atts := []AttachmentDownload{}

	for _, d := range downloads {

		s, o, _, err := r.c.AttachmentDownloadStream(int(d))
		if err != nil {
			return nil, fmt.Errorf("redmine attachments download: %w", err)
		}

		atts = append(
			atts,
			AttachmentDownload{
				Reader:      s,
				Name:        o.FileName,
				Caption:     o.Description,
				ContentType: o.ContentType,
			},
		)
	}

	return atts, nil
}
