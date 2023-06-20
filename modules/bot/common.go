package tgbot

import (
	"errors"
	"fmt"

	tg "github.com/nixys/nxs-go-telegram"
	"github.com/nixys/nxs-support-bot/modules/issues"
	"github.com/nixys/nxs-support-bot/modules/localization"
	"github.com/nixys/nxs-support-bot/modules/users"
)

type userEnv struct {
	user users.User
	l    localization.Lang
}

// userEnvInit initiates a user from datasource and try to save to slot
func userEnvInit(t *tg.Telegram, sess *tg.Session) (userEnv, error) {

	// Get bot context
	bCtx := botCtxGet(t)

	// Get user from datasources
	u, err := bCtx.users.Get(sess.UserIDGet())
	if err != nil {
		return userEnv{
			user: u,
			l:    bCtx.lb.DefaultLanguage(),
		}, fmt.Errorf("user env init: %w", err)
	}

	// Set user language
	l, err := bCtx.lb.LangSwitch(u.Lang)
	if err != nil {
		return userEnv{
			user: u,
			l:    bCtx.lb.DefaultLanguage(),
		}, fmt.Errorf("user env init: %w", err)
	}

	// Save user to slot
	if err = sess.SlotSave(slotNameUser, u); err != nil && errors.Is(err, tg.ErrSessionNotExist) == false {

		// Ignoring `ErrSessionNotExist` error

		return userEnv{
			user: u,
			l:    l,
		}, fmt.Errorf("user env init: %w", err)
	}

	return userEnv{
		user: u,
		l:    l,
	}, nil
}

// userEnvGet tries to get user from slot and retrieves info from datasource if session or slot are not exist
func userEnvGet(t *tg.Telegram, sess *tg.Session) (userEnv, error) {

	var u users.User

	// Get bot context
	bCtx := botCtxGet(t)

	b, err := sess.SlotGet(slotNameUser, &u)
	if err != nil && errors.Is(err, tg.ErrSessionNotExist) == false {
		return userEnv{
			user: users.User{},
			l:    bCtx.lb.DefaultLanguage(),
		}, fmt.Errorf("user env get: %w", err)
	}

	if (err != nil && errors.Is(err, tg.ErrSessionNotExist) == false) || b == false {

		// If session or slot not exist, get user from datasource

		u, err = bCtx.users.Get(sess.UserIDGet())
		if err != nil {
			return userEnv{
				user: u,
				l:    bCtx.lb.DefaultLanguage(),
			}, fmt.Errorf("user env get: %w", err)
		}
	}

	// Set user language
	l, err := bCtx.lb.LangSwitch(u.Lang)
	if err != nil {
		return userEnv{
			user: u,
			l:    bCtx.lb.DefaultLanguage(),
		}, fmt.Errorf("user env get: %w", err)
	}

	return userEnv{
		user: u,
		l:    l,
	}, nil
}

func renderCtrlButtonsPage(l localization.Lang, isPrev, isBack bool) []tg.Button {

	var ctrlButtons []tg.Button

	// Render prev button
	if isPrev == true {
		ctrlButtons = append(
			ctrlButtons,
			tg.Button{
				Text:       l.BotButton(localization.ButtonPrevPage),
				Identifier: buttonIDPrevPage,
			},
		)
	}

	// Render back button
	ctrlButtons = append(
		ctrlButtons,
		tg.Button{
			Text:       l.BotButton(localization.ButtonBack),
			Identifier: buttonIDBack,
		},
	)

	// Render next button
	if isBack == true {
		ctrlButtons = append(
			ctrlButtons,
			tg.Button{
				Text:       l.BotButton(localization.ButtonNextPage),
				Identifier: buttonIDNextPage,
			},
		)
	}

	return ctrlButtons
}

func attachmentsIssuesUpload(t *tg.Telegram, sess *tg.Session) ([]issues.AttachmentUpload, error) {

	bCtx := botCtxGet(t)

	c, err := userEnvGet(t, sess)
	if err != nil {
		return nil, err
	}

	// Extract files from Telegram updates
	files, err := sess.UpdateChain().FilesGet(*t)
	if err != nil {
		return nil, err
	}

	uds := []issues.UploadData{}

	for _, f := range files {

		// Get Telegram stream for file
		s, err := t.DownloadFileStream(f)
		if err != nil {
			return nil, err
		}

		// Add issue upload data
		uds = append(
			uds,
			issues.UploadData{
				File: s,
				Name: f.FileName,
			},
		)
	}

	// Upload files to issue tracker
	atts, err := bCtx.issues.AttachmentsUpload(c.user.RdmnID, uds)
	if err != nil {
		return nil, err
	}

	return atts, nil
}
