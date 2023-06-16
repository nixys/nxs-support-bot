package primedb

import "git.nixys.ru/apps/nxs-support-bot/misc"

const FeedbackIssueTable = "feedback_issues"

type FeedbackIssue struct {
	TgID    int64 `gorm:"column:tlgrm_userid"`
	IssueID int64 `gorm:"column:rdmn_issue_id"`
}

type FeedbackIssueInsertData struct {
	TgID    int64
	IssueID int64
}

func (FeedbackIssue) TableName() string {
	return FeedbackIssueTable
}

func (db *DB) FeedbackIssueSave(feedbackIssue FeedbackIssueInsertData) (FeedbackIssue, error) {

	p := FeedbackIssue{
		TgID:    feedbackIssue.TgID,
		IssueID: feedbackIssue.IssueID,
	}

	r := db.client.
		Create(&p)
	if r.Error != nil {
		return FeedbackIssue{}, r.Error
	}

	return p, nil
}

func (db *DB) FeedbackIssueGet(tgID int64) (FeedbackIssue, error) {

	p := FeedbackIssue{}

	r := db.client.
		Where(
			FeedbackIssue{
				TgID: tgID,
			},
		).
		Find(&p)
	if r.Error != nil {
		return FeedbackIssue{}, r.Error
	}

	if r.RowsAffected == 0 {
		return FeedbackIssue{}, misc.ErrNotFound
	}

	return p, nil
}

func (db *DB) FeedbackIssueGetByIssueID(issueID int64) (FeedbackIssue, error) {

	p := FeedbackIssue{}

	r := db.client.
		Where(
			FeedbackIssue{
				IssueID: issueID,
			},
		).
		Find(&p)
	if r.Error != nil {
		return FeedbackIssue{}, r.Error
	}

	if r.RowsAffected == 0 {
		return FeedbackIssue{}, misc.ErrNotFound
	}

	return p, nil
}
