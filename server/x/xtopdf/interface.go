package xplugin

import "io"

type TOPDF interface {
	CheckServerStatus() (err error)
	GetPDF(userID, fileID string) (pdf io.ReadCloser, err error)
}
