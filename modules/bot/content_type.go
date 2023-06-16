package tgbot

import tg "github.com/nixys/nxs-go-telegram"

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

func fileTypeGet(contentType string) tg.FileType {
	t, b := contentTypes[contentType]
	if b == true {
		return t
	}
	return tg.FileTypeDocument
}
