package primedb

import "git.nixys.ru/apps/nxs-support-bot/misc"

const IssuesBunchTable = "issues_banch"

type IssuesBunch struct {
	ChatID    int64 `gorm:"column:tg_chat_id"`
	MessageID int64 `gorm:"column:tg_message_id"`
	IssueID   int64 `gorm:"column:rdmn_issue_id"`
}

type IssuesBunchInsertData struct {
	ChatID    int64
	MessageID int64
	IssueID   int64
}

func (IssuesBunch) TableName() string {
	return IssuesBunchTable
}

func (db *DB) IssuesBunchSave(issueBunch IssuesBunchInsertData) (IssuesBunch, error) {

	i := IssuesBunch{
		ChatID:    issueBunch.ChatID,
		MessageID: issueBunch.MessageID,
		IssueID:   issueBunch.IssueID,
	}

	r := db.client.
		Create(&i)
	if r.Error != nil {
		return IssuesBunch{}, r.Error
	}

	return i, nil
}

func (db *DB) IssuesBunchGet(chatID, messageID int64) (IssuesBunch, error) {

	b := IssuesBunch{}

	r := db.client.
		Where(
			IssuesBunch{
				ChatID:    chatID,
				MessageID: messageID,
			},
		).
		Find(&b)
	if r.Error != nil {
		return IssuesBunch{}, r.Error
	}

	if r.RowsAffected == 0 {
		return IssuesBunch{}, misc.ErrNotFound
	}

	return b, nil
}
