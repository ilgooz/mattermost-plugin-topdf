package main

import (
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/ilgooz/mattermost-plugin-topdf/server/gotenberg"
	"github.com/ilgooz/mattermost-plugin-topdf/server/topdf"
	"github.com/ilgooz/mattermost-plugin-topdf/server/topdf/pdfserver"
	"github.com/ilgooz/mattermost-plugin-topdf/server/x/xhttp"
	"github.com/ilgooz/mattermost-plugin-topdf/server/x/xtime"
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/rs/cors"
)

// Plugin is a TOPDF app that implements plugin.MattermostPlugin.
type Plugin struct {
	plugin.MattermostPlugin

	// app is the actual, underlying TOPDF app and its features exposed to
	// network via Plugin's HTTP API.
	app interface {
		CheckServerStatus() (err error)
		GetPDF(userID, fileID string) (pdf io.ReadCloser, err error)
	} // *topdf.TOPDF
}

// configuration holds Plugin's config.
// see plugin.json at root for more info about all configurations.
type configuration struct {
	GotenbergAddress        string
	GotenbergConvertTimeout xtime.Duration
}

func main() {
	plugin.ClientMain(&Plugin{})
}

// OnConfigurationChange hook updates underlying Gotenberg configs.
func (p *Plugin) OnConfigurationChange() error {
	var conf configuration
	if err := p.API.LoadPluginConfiguration(&conf); err != nil {
		return err
	}
	p.init(conf)
	return nil
}

// init initializes a new topdf with given c.
func (p *Plugin) init(c configuration) {
	gt := gotenberg.New(c.GotenbergAddress, []gotenberg.Option{
		gotenberg.ConvertTimeoutOption(time.Duration(c.GotenbergConvertTimeout)),
	}...)
	p.app = topdf.New(p.MattermostPlugin.API, gt)
}

// ServeHTTP hook exposes a RESTful API for `topdf` plugin.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()
	// GET /status gives status info about underlying(Gotenberg) PDF server.
	router.HandleFunc("/status", p.handleStatus).Methods("GET")
	// GET /files/{id} responses with a PDF version of file that attached to a Mattermost Post.
	// it caches PDF files that requested for same files.
	router.HandleFunc("/files/{id}", p.handleConvert).Methods("GET")
	// allow CORS for the API.
	handler := cors.AllowAll().Handler(router)
	// serve request.
	handler.ServeHTTP(w, r)
}

// handleStatus handles Plugin's status check requests.
func (p *Plugin) handleStatus(w http.ResponseWriter, r *http.Request) {
	err := p.app.CheckServerStatus()
	if err != nil {
		if _, ok := err.(*pdfserver.NotReachable); !ok {
			xhttp.ResponseJSON(w, http.StatusInternalServerError, createErrorResponse(err))
			p.logError(err)
			return
		}
		xhttp.ResponseJSON(w, http.StatusOK, statusResponse{IsGotenbergRunning: false})
		return
	}
	xhttp.ResponseJSON(w, http.StatusOK, statusResponse{IsGotenbergRunning: true})
}

// handlePDF handles file to PDF convert requests.
func (p *Plugin) handleConvert(w http.ResponseWriter, r *http.Request) {
	fileID := mux.Vars(r)["id"]
	userID := r.Header.Get("Mattermost-User-Id")
	// check if there is authenticated user, otherwise fail request since accessing files always
	// requires a user.
	if userID == "" {
		xhttp.ResponseJSON(w, http.StatusUnauthorized, createErrorResponse(topdf.ErrUnauthorizedUser))
		p.logError(topdf.ErrUnauthorizedUser)
		return
	}
	// get pdf for fileID with userID.
	// if user does not have access to file, requester will be responded with authorization error.
	pdf, err := p.app.GetPDF(userID, fileID)
	if err != nil {
		code := http.StatusInternalServerError
		if err == topdf.ErrUnauthorizedUser {
			code = http.StatusUnauthorized
		}
		xhttp.ResponseJSON(w, code, createErrorResponse(err))
		p.logError(err)
		return
	}
	defer pdf.Close()
	w.Header().Set("Content-Type", "application/pdf")
	// stream PDF content to requester.
	io.Copy(w, pdf)
}

// logError logs errors with Plugin API.
func (p *Plugin) logError(err error) {
	p.API.LogError(err.Error())
}

// createErrorResponse creates a new error response from err to be sent HTTP client.
func createErrorResponse(err error) errorResponse {
	return errorResponse{
		errorResponseBody{Message: err.Error()},
	}
}

// statusResponse is status response sent to client.
type statusResponse struct {
	IsGotenbergRunning bool `json:"isGotenbergRunning"`
}

// statusResponse is error response sent to client.
type errorResponse struct {
	Error errorResponseBody `json:"error"`
}

type errorResponseBody struct {
	Message string `json:"message"`
}
