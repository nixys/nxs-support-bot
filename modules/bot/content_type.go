package tgbot

import (
	"bytes"
	"image"
	"io"

	_ "image/jpeg"
	_ "image/png"

	tg "github.com/nixys/nxs-go-telegram"
)

var contentTypes = map[string]tg.FileType{
	"image/jpeg":      tg.FileTypePhoto,
	"image/bmp":       tg.FileTypePhoto,
	"image/png":       tg.FileTypePhoto,
	"application/ogg": tg.FileTypeVoice,
	"video/mp4":       tg.FileTypeVideo,
	"video/webm":      tg.FileTypeVideo,
	"video/quicktime": tg.FileTypeVideo,
	"image/webp":      tg.FileTypeSticker,
	"audio/mpeg":      tg.FileTypeAudio,
}

func fileTypeGet(contentType string, r io.Reader) (tg.FileType, io.Reader) {

	t, b := contentTypes[contentType]
	if b == false {
		t = tg.FileTypeDocument
	}

	// Prevent error from Tg API: Bad Request: PHOTO_INVALID_DIMENSIONS
	switch t {
	case tg.FileTypePhoto:

		// Download image into buffer and create a reader
		buf := &bytes.Buffer{}
		buf.ReadFrom(r)
		rr := bytes.NewReader(buf.Bytes())

		// Get image settings
		img, _, err := image.DecodeConfig(rr)
		rr.Seek(0, 0)
		if err != nil {
			return tg.FileTypeDocument, rr
		}

		// Calc pict dimensions
		if img.Width == 0 || img.Height > 0 {
			return tg.FileTypeDocument, rr
		}
		if img.Width >= img.Height {
			if img.Width/img.Height >= 20 {
				return tg.FileTypeDocument, rr
			}
		} else {
			if img.Height/img.Width >= 20 {
				return tg.FileTypeDocument, rr
			}
		}

		return t, rr
	}

	return t, r
}
