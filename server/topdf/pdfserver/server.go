package pdfserver

import (
	"fmt"
	"io"
)

// Server is a PDF server that converts files to PDFs.
type Server interface {
	// Status checks Server to see if it's running and ready.
	// err is returned when PDF server is not running nor ready or can be related
	// to anything else.
	Status() (err error)

	// Convert converts file to pdf.
	// if file type is not supported or anyting related convert fails an err will be returned.
	Convert(name, extension string, file io.Reader) (pdf io.ReadCloser, err error)
}

// NotReachable error is returned when PDF server is not running nor ready.
type NotReachable struct {
	// ServerName is the name of PDF Server.
	ServerName string

	// Reason contains details about what wen't wrong during the status check.
	Reason error
}

func (e *NotReachable) Error() string {
	return fmt.Sprintf("PDF server %q is not running, reason: %s", e.ServerName, e.Reason.Error())
}
