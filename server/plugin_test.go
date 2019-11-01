package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ilgooz/mattermost-plugin-topdf/server/topdf"

	pMock "github.com/ilgooz/mattermost-plugin-topdf/server/x/xplugin/mocks"
	tMock "github.com/ilgooz/mattermost-plugin-topdf/server/x/xtopdf/mocks"
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/stretchr/testify/require"
)

func TestHandleStatusRunning(t *testing.T) {
	topdfMock := &tMock.TOPDF{}
	p := &Plugin{app: topdfMock}
	req := httptest.NewRequest("GET", "http://localhost.com/status", nil)
	w := httptest.NewRecorder()
	topdfMock.On("CheckServerStatus").Once().Return(true, nil)
	p.ServeHTTP(nil, w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	require.Equal(t, `{"isGotenbergRunning":true}`, string(body))
	topdfMock.AssertExpectations(t)
}

func TestHandleStatusNotRunning(t *testing.T) {
	topdfMock := &tMock.TOPDF{}
	p := &Plugin{app: topdfMock}
	req := httptest.NewRequest("GET", "http://localhost.com/status", nil)
	w := httptest.NewRecorder()
	topdfMock.On("CheckServerStatus").Once().Return(false, nil)
	p.ServeHTTP(nil, w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	require.Equal(t, `{"isGotenbergRunning":false}`, string(body))
	topdfMock.AssertExpectations(t)
}

func TestHandleStatusInternalServerError(t *testing.T) {
	topdfMock := &tMock.TOPDF{}
	apiMock := &pMock.API{}
	p := &Plugin{MattermostPlugin: plugin.MattermostPlugin{API: apiMock}, app: topdfMock}
	req := httptest.NewRequest("GET", "http://localhost.com/status", nil)
	w := httptest.NewRecorder()
	topdfMock.On("CheckServerStatus").Once().Return(false, errors.New("a failure"))
	apiMock.On("LogError", "a failure").Once()
	p.ServeHTTP(nil, w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	require.Equal(t, `{"error":{"message":"a failure"}}`, string(body))
	topdfMock.AssertExpectations(t)
	apiMock.AssertExpectations(t)
}

func TestHandleConvert(t *testing.T) {
	topdfMock := &tMock.TOPDF{}
	p := &Plugin{app: topdfMock}
	req := httptest.NewRequest("GET", "http://localhost.com/files/1", nil)
	req.Header.Set("Mattermost-User-Id", "2")
	w := httptest.NewRecorder()
	topdfMock.On("GetPDF", "2", "1").Once().Return(ioutil.NopCloser(bytes.NewReader([]byte{3})), nil)
	p.ServeHTTP(nil, w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "application/pdf", resp.Header.Get("Content-Type"))
	require.Equal(t, []byte{3}, body)
	topdfMock.AssertExpectations(t)
}

func TestHandleConvertInternalError(t *testing.T) {
	topdfMock := &tMock.TOPDF{}
	apiMock := &pMock.API{}
	p := &Plugin{MattermostPlugin: plugin.MattermostPlugin{API: apiMock}, app: topdfMock}
	req := httptest.NewRequest("GET", "http://localhost.com/files/1", nil)
	req.Header.Set("Mattermost-User-Id", "2")
	w := httptest.NewRecorder()
	topdfMock.On("GetPDF", "2", "1").Once().Return(nil, errors.New("internal"))
	apiMock.On("LogError", "internal").Once()
	p.ServeHTTP(nil, w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	require.Equal(t, `{"error":{"message":"internal"}}`, string(body))
	topdfMock.AssertExpectations(t)
}

func TestHandleConvertAuthorizationError(t *testing.T) {
	topdfMock := &tMock.TOPDF{}
	apiMock := &pMock.API{}
	p := &Plugin{MattermostPlugin: plugin.MattermostPlugin{API: apiMock}, app: topdfMock}
	req := httptest.NewRequest("GET", "http://localhost.com/files/1", nil)
	req.Header.Set("Mattermost-User-Id", "2")
	w := httptest.NewRecorder()
	topdfMock.On("GetPDF", "2", "1").Once().Return(nil, topdf.ErrUnauthorizedUser)
	apiMock.On("LogError", "user is not authorized to access pdf").Once()
	p.ServeHTTP(nil, w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	require.Equal(t, `{"error":{"message":"user is not authorized to access pdf"}}`, string(body))
	topdfMock.AssertExpectations(t)
	apiMock.AssertExpectations(t)
}

func TestHandleConvertWithoutAuthentication(t *testing.T) {
	apiMock := &pMock.API{}
	p := &Plugin{MattermostPlugin: plugin.MattermostPlugin{API: apiMock}}
	req := httptest.NewRequest("GET", "http://localhost.com/files/1", nil)
	w := httptest.NewRecorder()
	apiMock.On("LogError", "user is not authorized to access pdf").Once()
	p.ServeHTTP(nil, w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	require.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	require.Equal(t, `{"error":{"message":"user is not authorized to access pdf"}}`, string(body))
	apiMock.AssertExpectations(t)
}
