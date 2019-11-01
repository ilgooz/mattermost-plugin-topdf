package xplugin

import "io"

type TOPDF interface {
	CheckServerStatus() (running bool, err error)
	GetPDF(userID, fileID string) (pdf io.ReadCloser, err error)
}
