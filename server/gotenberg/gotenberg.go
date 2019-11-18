// Package gotenberg is a Gotenberg client with stream (io.Reader) support.
package gotenberg

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

const (
	// pingTimeout defines the timeout for ping requests to the Gotenberg server.
	pingTimeout = time.Second * 5

	// defaultConvertTimeout defines the timeout for file conversion requests to the Gotenberg server
	// to Gotenberg server.
	defaultConvertTimeout = time.Minute * 10
)

const (
	// pingEndpoint used to status check Gotenberg to see if it's running and ready.
	pingEndpoint = "/ping"

	// convertEndpoint used to convert files to PDFs.
	convertEndpoint = "/convert/office"
)

const multipartField = "file"

// supportedFormats are the supported file formats  that can be converted to PDF by Gotenberg.
var supportedFormats = []string{"doc", "docx", "odt", "xls", "xlsx", "ods", "ppt", "pptx", "odp"}

// Gotenberg is a client for Gotenberg server.
// for more info see: https://github.com/thecodingmachine/gotenberg
// warning! don't forget to set proper timeouts as you need in Gotenberg server because default ones are too low. ex:
// - docker run -d -p 4798:3000 --env DEFAULT_WAIT_TIMEOUT=600 --env MAXIMUM_WAIT_TIMEOUT=600 thecodingmachine/gotenberg:6
type Gotenberg struct {
	// addr is network address of Gotenberg server.
	addr string
	// convertTimout is used during sending file convert requests.
	convertTimeout time.Duration
}

// New creates new Gotenberg client with given Gotenberg server addr and options.
func New(addr string, options ...Option) *Gotenberg {
	g := &Gotenberg{addr: addr}
	g.applyOptions(options...)
	return g
}

// applyOptions applies user given options to Gotenberg configuration.
func (g *Gotenberg) applyOptions(options ...Option) {
	for _, o := range options {
		o(g)
	}
	if g.convertTimeout == 0 {
		g.convertTimeout = defaultConvertTimeout
	}
}

// Option used to customize Gotenberg defaults.
type Option func(*Gotenberg)

// ConversionTimeoutOption sets the timeout for conversion requests.
// to Gotenberg server.
func ConvertTimeoutOption(convertTimeout time.Duration) Option {
	return func(g *Gotenberg) {
		g.convertTimeout = convertTimeout
	}
}

// Status checks if Gotenberg server is running and ready to accept connections.
func (g *Gotenberg) Status() (running bool, err error) {
	c := &http.Client{Timeout: pingTimeout}
	url, err := buildGotenbergURL(g.addr, pingEndpoint)
	if err != nil {
		return false, err
	}
	if _, err := c.Get(url); err != nil {
		return false, nil
	}
	return true, nil
}

// Convert converts file with given name and extension to PDF.
// caller is responsible to Close() PDF stream after done.
func (g *Gotenberg) Convert(name, extension string, file io.Reader) (pdf io.ReadCloser, err error) {
	// check to see if given file extension is supported.
	if !isSupported(extension) {
		return nil, fmt.Errorf("file extension `%s` is not supported by the PDF server", extension)
	}
	// create a pipe and:
	// - give the pw to multipart writer so it can start writing multipart data back while reading
	//   the contents of file(io.Reader).
	// - give the pr to HTTP request that will be made to Gotenberg. this way, as Gotenberg server
	//   do reads, multipart writer can continue to write in sync.
	//
	// with this approach we have data cycled as a stream which reduces memory usage.
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	// closer finishes multipart writer's work and closes the pipe with a regular EOF or an error.
	// in both cases HTTP request to Gotenberg will be completed resulting without an error (in case of EOF)
	// or cancelling with an error (in case of non-EOF).
	closer := func(err error) {
		writer.Close()
		pw.CloseWithError(err)
	}
	// create a 'multipart file' and copy whole content of file as Gotenberg server continues to read.
	go func() {
		part, err := writer.CreateFormFile(multipartField, fmt.Sprintf("%s.%s", name, extension))
		if err != nil {
			closer(err)
			return
		}
		_, err = io.Copy(part, file)
		closer(err)
	}()
	url, err := buildGotenbergURL(g.addr, convertEndpoint)
	if err != nil {
		return nil, err
	}
	// make an HTTP request to Gotenberg to initialize process.
	req, err := http.NewRequest("POST", url, pr)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	c := &http.Client{Timeout: g.convertTimeout}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	// check if Gotenberg is cool with the file we sent to see if it's gonna response back with a PDF data.
	if res.StatusCode != http.StatusOK &&
		res.StatusCode != http.StatusCreated {
		defer res.Body.Close()
		// request is somehow not successful:
		// file content can be invalid or some timeout might be hitting set by Gotenberg's end or here in the request.
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, errors.Wrap(err, "error while reading error message from Gotenberg")
		}
		return nil, fmt.Errorf("error from Gotenberg with '%d' code: %s", res.StatusCode, string(data))
	}
	// we have Gotenberg willing to stream PDF data, give it to the caller so it can start reading.
	return res.Body, nil
}

// buildGotenbergURL generates a Gotenberg API URL from given addr for endpoint.
func buildGotenbergURL(addr, endpoint string) (string, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return "", err
	}
	u.Path = endpoint
	return u.String(), nil
}

// isSupported check if file extension is supported.
func isSupported(extension string) (ok bool) {
	for _, supext := range supportedFormats {
		if supext == extension {
			return true
		}
	}
	return false
}
