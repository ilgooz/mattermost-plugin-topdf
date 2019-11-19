package topdf

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"

	"github.com/ilgooz/mattermost-plugin-topdf/server/topdf/pdfserver"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

// toPDFPrefix used as a prefix while using KV store to access PDF files' fileID.
const toPDFPrefix = "pdf:"

// ErrUnauthorizedUser returned when user has no access to a file that requested to be converted to PDF.
var ErrUnauthorizedUser = errors.New("user is not authorized to access pdf")

// TOPDF is an application that converts files to PDFs and permanently caches them by using Mattermost APIs.
type TOPDF struct {
	// mapi is Mattermost's Plugin API.
	mapi plugin.API

	// server used to convert files to PDF.
	server pdfserver.Server
}

// New creates a new TOPDF app with mapi and PDF server.
func New(mapi plugin.API, server pdfserver.Server) *TOPDF {
	return &TOPDF{
		mapi:   mapi,
		server: server,
	}
}

// CheckServerStatus checks if underlying PDF server is running and ready to accept requests.
func (t *TOPDF) CheckServerStatus() error {
	return t.server.Status()
}

// GetPDF gets PDF for fileID that belongs userID. user has to have access to the file
// otherwise ErrUnauthorizedUser is returned.
func (t *TOPDF) GetPDF(userID, fileID string) (pdf io.ReadCloser, err error) {
	pdf, err = t.getPDF(userID, fileID)
	if err != nil {
		// return an authorization error if we got an err from Plugin's API.
		if _, ok := err.(*model.AppError); ok {
			return nil, ErrUnauthorizedUser
		}
		return nil, err
	}
	return pdf, nil
}

// getPDF gets PDF for fileID that belongs to userID.
// notes:
// - Mattermost's Plugin API does not implement io.Reader while dealing with files but this might
//   be improved in future since large files can pump memory usage. TOPDF created streams in mind,
//   this is why we use ioutil.NopCloser and bytes.Reader in the code below, to even with current Plugin API.
// - Mattermost's Plugin API does not return errors as `error`s but returns them as *model.AppError,
//   this needs to be improved since it causes issues while dealing with errors. For more info
//   please see: https://golang.org/doc/faq#nil_error
//   to workaround this, Plugin errors are normalized with normalizeAppErr().
func (t *TOPDF) getPDF(userID, fileID string) (pdf io.ReadCloser, err error) {
	// try to get id of PDF file that possibly generated and cached for fileID before.
	pid, aerr := t.mapi.KVGet(key(fileID))
	if aerr != nil {
		return nil, normalizeAppErr(aerr)
	}
	// get file's info.
	fileInfo, aerr := t.mapi.GetFileInfo(fileID)
	if aerr != nil {
		return nil, normalizeAppErr(aerr)
	}
	// get associated post for the file.
	filePost, aerr := t.mapi.GetPost(fileInfo.PostId)
	if aerr != nil {
		return nil, normalizeAppErr(aerr)
	}
	// check if the user has access to the channel where associated post submitted.
	if _, aerr := t.mapi.GetChannelMember(filePost.ChannelId, userID); aerr != nil {
		return nil, normalizeAppErr(aerr)
	}
	// if there is no PDF file cached, create it, cache and use its content.
	if len(pid) == 0 {
		data, err := t.createAndSavePDF(fileInfo, filePost)
		if err != nil {
			return nil, err
		}
		return ioutil.NopCloser(bytes.NewReader(data)), nil
	}
	// we have the PDF version in cache, directly return it back.
	data, err := t.getCachedPDF(string(pid))
	if err != nil {
		return nil, err
	}
	return ioutil.NopCloser(bytes.NewReader(data)), nil
}

// createAndSavePDF creates a PDF version of fileID and caches on Mattermost server and returns
// the pdf data back.
func (t *TOPDF) createAndSavePDF(fileInfo *model.FileInfo, filePost *model.Post) (pdf []byte, err error) {
	// get file's content by fileID.
	fileBytes, aerr := t.mapi.GetFile(fileInfo.Id)
	if aerr != nil {
		return nil, normalizeAppErr(aerr)
	}
	// convert file to PDF by using PDF server.
	r, err := t.server.Convert(fileInfo.Name, fileInfo.Extension, bytes.NewReader(fileBytes))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	// cache PDF file on Mattermost.
	inf, aerr := t.mapi.UploadFile(data, filePost.ChannelId, "pdf")
	if err != nil {
		return nil, normalizeAppErr(aerr)
	}
	// save PDF file's id by associating it with fileID.
	if aerr := t.mapi.KVSet(key(fileInfo.Id), []byte(inf.Id)); err != nil {
		return nil, normalizeAppErr(aerr)
	}
	// return PDF file's content.
	return data, nil
}

// getCachedPDF gets cached PDF data from file store.
func (t *TOPDF) getCachedPDF(fileID string) (pdf []byte, err error) {
	data, aerr := t.mapi.GetFile(string(fileID))
	if aerr != nil {
		return nil, normalizeAppErr(aerr)
	}
	return data, nil
}

// key builds a KV key for fileID.
func key(fileID string) string {
	return toPDFPrefix + fileID
}

// normalize error normalizes Plugin API's errors.
// please see this docs to know more about what this normalization do: https://golang.org/doc/faq#nil_error
func normalizeAppErr(err *model.AppError) error {
	if err == nil {
		return nil
	}
	return err
}
