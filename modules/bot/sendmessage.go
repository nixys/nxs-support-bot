package tgbot

import (
	"fmt"
	"io"
	"io/ioutil"
	"sync"

	"github.com/hashicorp/go-multierror"
	tg "github.com/nixys/nxs-go-telegram"
	"github.com/nixys/nxs-support-bot/misc"
)

type SendData struct {
	Rcpts []SendRcptData
	Files []SendFileData
}

type SendRcptData struct {
	ChatID  int64
	Message string
}

type SendFileData struct {
	Reader      io.Reader
	Name        string
	Caption     string
	ContentType string
}

type SentResult struct {
	ChatID     int64
	MessageIDs []int64
}

type fupload struct {
	r io.Reader
	f tg.FileSendStream
}

const messageLengthLimit = 4096

func (b *Bot) SendMessage(messages SendData) ([]SentResult, error) {

	var (
		wg sync.WaitGroup
		sr []SentResult
	)

	// Chats count
	l := len(messages.Rcpts)

	// Sum of len chatIDs (created routine for each chat to send messages and files)
	// and files to be sent (created routine for each file to prepare readers)
	// Reason: error processing
	errChanCapacity := l + len(messages.Files)

	errch := make(chan error, errChanCapacity)
	files := make([][]fupload, l)

	// Channel to get send messages results from goroutines
	resch := make(chan SentResult, l)

	// Prepare readers for every files for all chats
	for _, f := range messages.Files {

		ft, fr := fileTypeGet(f.ContentType, f.Reader)

		// Create file readers for every Telegram receiver
		rds := misc.MultiReaderCreate(fr, l, errch)

		for j := 0; j < l; j++ {

			files[j] = append(files[j], fupload{
				r: rds[j],
				f: tg.FileSendStream{
					FileName:  f.Name,
					Caption:   f.Caption,
					ParseMode: tg.ParseModeHTML,
					FileType:  ft,
				},
			})
		}
	}

	// Send message and prepared files for all chats
	for i := range messages.Rcpts {

		wg.Add(1)

		go func(i int) {

			// Do send messages to every user we need
			r, err := b.sender(messages.Rcpts[i], files[i])

			wg.Done()

			// Send result
			resch <- r

			// Send error
			errch <- err
		}(i)
	}

	// Wait goroutines
	wg.Wait()

	// Collect an errors from routines
	var merr error
	for i := 0; i < errChanCapacity; i++ {
		if e := <-errch; e != nil {
			merr = multierror.Append(merr, e)
		}
	}

	// Collect results from senders
	for i := 0; i < l; i++ {

		r := <-resch
		sr = append(sr, r)
	}

	return sr, merr
}

func (b *Bot) sender(rcpt SendRcptData, files []fupload) (SentResult, error) {

	var (
		merr error
		r    SentResult
	)

	r.ChatID = rcpt.ChatID

	for _, msgChunk := range misc.MessageSplit(rcpt.Message, messageLengthLimit) {

		ms, err := b.bot.SendMessage(rcpt.ChatID, 0, tg.SendMessageData{
			Message:               msgChunk,
			ParseMode:             tg.ParseModeHTML,
			DisableWebPagePreview: true,
		})
		if err != nil {

			// On error discard read for all files
			// to prevent lock writer

			for _, f := range files {
				io.Copy(ioutil.Discard, f.r)
			}

			return r, fmt.Errorf("bot sender %d: %w", rcpt.ChatID, err)
		}

		for _, m := range ms {
			r.MessageIDs = append(r.MessageIDs, int64(m.MessageID))
		}
	}

	for _, f := range files {

		ms, err := b.bot.UploadFileStream(rcpt.ChatID, f.f, f.r)
		if err != nil {
			merr = multierror.Append(merr, fmt.Errorf("bot sender upload file %s: %w", f.f.FileName, err))
			continue
		}

		r.MessageIDs = append(r.MessageIDs, int64(ms.MessageID))
	}

	return r, merr
}
