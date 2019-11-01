package gotenberg

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStatusRunning(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/ping", r.URL.Path)
		require.Equal(t, http.MethodGet, r.Method)
	}))
	defer ts.Close()
	gt := New(ts.URL)
	isRunning, err := gt.Status()
	require.NoError(t, err)
	require.True(t, isRunning)
}

func TestStatusNotRunning(t *testing.T) {
	gt := New("http://not-existent-host")
	isRunning, err := gt.Status()
	require.NoError(t, err)
	require.False(t, isRunning)
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
