package topdf

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"

	sMock "github.com/ilgooz/mattermost-plugin-topdf/server/topdf/mocks"
	pMock "github.com/ilgooz/mattermost-plugin-topdf/server/x/xplugin/mocks"
	"github.com/mattermost/mattermost-server/model"
	"github.com/stretchr/testify/require"
)

func TestCheckServerStatusRunning(t *testing.T) {
	serverMock := &sMock.PDFServer{}
	serverMock.On("Status").Once().Return(true, nil)
	app := New(nil, serverMock)
	isRunning, err := app.CheckServerStatus()
	require.NoError(t, err)
	require.True(t, isRunning)
	serverMock.AssertExpectations(t)
}

func TestCheckServerStatusNotRunning(t *testing.T) {
	serverMock := &sMock.PDFServer{}
	serverMock.On("Status").Once().Return(false, nil)
	app := New(nil, serverMock)
	isRunning, err := app.CheckServerStatus()
	require.NoError(t, err)
	require.False(t, isRunning)
	serverMock.AssertExpectations(t)
}

func TestCheckServerStatusError(t *testing.T) {
	serverMock := &sMock.PDFServer{}
	serverMock.On("Status").Once().Return(false, errors.New("ops!"))
	app := New(nil, serverMock)
	_, err := app.CheckServerStatus()
	require.Equal(t, "ops!", err.Error())
	serverMock.AssertExpectations(t)
}

func TestCheckServerConvertCached(t *testing.T) {
	serverMock := &sMock.PDFServer{}
	apiMock := &pMock.API{}
	apiMock.On("KVGet", "pdf:file-id").Once().Return([]byte("1"), nil)
	apiMock.On("GetFile", "1").Once().Return([]byte{2}, nil)
	app := New(apiMock, serverMock)
	pdf, err := app.GetPDF("user-id", "file-id")
	require.NoError(t, err)
	data, err := ioutil.ReadAll(pdf)
	require.NoError(t, err)
	require.Equal(t, []byte{2}, data)
	serverMock.AssertExpectations(t)
	apiMock.AssertExpectations(t)
}

func TestCheckServerConvertNonCached(t *testing.T) {
	serverMock := &sMock.PDFServer{}
	apiMock := &pMock.API{}
	apiMock.On("KVGet", "pdf:file-id").Once().Return([]byte{}, nil)
	apiMock.On("GetFileInfo", "file-id").Once().Return(&model.FileInfo{PostId: "2", Name: "3", Extension: "4"}, nil)
	apiMock.On("GetPost", "2").Once().Return(&model.Post{ChannelId: "5"}, nil)
	apiMock.On("GetChannelMember", "5", "user-id").Once().Return(nil, nil)
	apiMock.On("GetFile", "file-id").Once().Return([]byte{3}, nil)
	serverMock.On("Convert", "3", "4", bytes.NewReader([]byte{3})).Once().Return(ioutil.NopCloser(bytes.NewReader([]byte{6})), nil)
	apiMock.On("UploadFile", []byte{6}, "5", "pdf").Once().Return(&model.FileInfo{Id: "7"}, nil)
	apiMock.On("KVSet", "pdf:file-id", []byte("7")).Once().Return(nil)
	app := New(apiMock, serverMock)
	pdf, err := app.GetPDF("user-id", "file-id")
	require.NoError(t, err)
	data, err := ioutil.ReadAll(pdf)
	require.NoError(t, err)
	require.Equal(t, []byte{6}, data)
	serverMock.AssertExpectations(t)
	apiMock.AssertExpectations(t)
}

func TestCheckServerConvertUnauthorized(t *testing.T) {
	serverMock := &sMock.PDFServer{}
	apiMock := &pMock.API{}
	apiMock.On("KVGet", "pdf:pdf-id").Once().Return(nil, &model.AppError{})
	app := New(apiMock, serverMock)
	_, err := app.GetPDF("user-id", "pdf-id")
	require.Equal(t, ErrUnauthorizedUser, err)
	apiMock.AssertExpectations(t)
}
