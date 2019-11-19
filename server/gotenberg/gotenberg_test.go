package gotenberg

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ilgooz/mattermost-plugin-topdf/server/topdf/pdfserver"
	"github.com/stretchr/testify/require"
)

func TestStatusRunning(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/ping", r.URL.Path)
		require.Equal(t, http.MethodGet, r.Method)
	}))
	defer ts.Close()
	gt := New(ts.URL)
	err := gt.Status()
	require.NoError(t, err)
}

func TestStatusNotRunning(t *testing.T) {
	// reserve a port to dedicate to this test. otherwise, a randomly chosen port might actually
	// be belong to a running server which will brake this test.
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	require.NoError(t, err)
	ln, err := net.ListenTCP("tcp", addr)
	require.NoError(t, err)
	defer ln.Close()
	port := ln.Addr().(*net.TCPAddr).Port
	urlStr := fmt.Sprintf("http://localhost:%d", port)
	gt := New(urlStr)
	err = gt.Status()
	require.IsType(t, &pdfserver.NotReachable{}, err)
	require.True(t, err.(*pdfserver.NotReachable).Reason.(*url.Error).Timeout())
}
func TestStatusNotRunningOK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()
	gt := New(ts.URL)
	err := gt.Status()
	require.Equal(t, &pdfserver.NotReachable{"Gotenberg", errors.New("received non-OK response code")}, err)
}

func TestConvertSupportedFormat(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/convert/office", r.URL.Path)
		require.Equal(t, http.MethodPost, r.Method)
		file, _, err := r.FormFile("file")
		require.NoError(t, err)
		defer file.Close()
		data, err := ioutil.ReadAll(file)
		require.NoError(t, err)
		require.Equal(t, "docx-file", string(data))
		w.Write([]byte("pdf-file"))
	}))
	defer ts.Close()
	gt := New(ts.URL)
	pdf, err := gt.Convert("name", "docx", strings.NewReader("docx-file"))
	require.NoError(t, err)
	defer pdf.Close()
	data, err := ioutil.ReadAll(pdf)
	require.NoError(t, err)
	require.Equal(t, "pdf-file", string(data))
}

func TestConvertUnsupportedFormat(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Fail(t, "should not call server in case of unsupported file extension")
	}))
	defer ts.Close()
	gt := New(ts.URL)
	_, err := gt.Convert("name", "txt", strings.NewReader("txt-file"))
	require.Equal(t, "file extension `txt` is not supported by the PDF server", err.Error())
}
