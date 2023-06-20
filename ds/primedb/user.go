package primedb

import "github.com/nixys/nxs-support-bot/misc"

const UsersTableName = "users"

type User struct {
	TgID   int64  `gorm:"column:tlgrm_userid"`
	RdmnID int64  `gorm:"column:rdmn_userid"`
	Lang   string `gorm:"column:lang"`
}

type UserUpdateData struct {
	TgID   int64   `gorm:"column:tlgrm_userid"`
	RdmnID *int64  `gorm:"column:rdmn_userid"`
	Lang   *string `gorm:"column:lang"`
}

func (User) TableName() string {
	return UsersTableName
}

func (UserUpdateData) TableName() string {
	return UsersTableName
}

func (db *DB) UserGet(tgID int64) (User, error) {

	user := User{}

	r := db.client.
		Where(
			User{
				TgID: tgID,
			},
		).
		Find(&user)
	if r.Error != nil {
		return User{}, r.Error
	}

	if r.RowsAffected == 0 {
		return User{}, misc.ErrNotFound
	}

	return user, nil
}

func (db *DB) UserGetByRdmnID(rdmnID int64) (User, error) {

	user := User{}

	r := db.client.
		Where(
			User{
				RdmnID: rdmnID,
			},
		).
		Find(&user)
	if r.Error != nil {
		return User{}, r.Error
	}

	if r.RowsAffected == 0 {
		return User{}, misc.ErrNotFound
	}

	return user, nil
}

func (db *DB) UserUpdate(u UserUpdateData) (User, error) {

	// Revoke other users with same Redmine
	// ID (previous owners this Redmine account)
	if u.RdmnID != nil && *u.RdmnID != 0 {

		r := db.client.
			Where(
				UserUpdateData{
					RdmnID: u.RdmnID,
				},
			).
			Updates(UserUpdateData{
				RdmnID: func() *int64 {
					i := int64(0)
					return &i
				}(),
			})
		if r.Error != nil {
			return User{}, r.Error
		}
	}

	user := User{}

	r := db.client.
		Where(
			UserUpdateData{
				TgID: u.TgID,
			},
		).
		Assign(u).
		FirstOrCreate(&user)
	if r.Error != nil {
		return User{}, r.Error
	}

	return user, nil
}
