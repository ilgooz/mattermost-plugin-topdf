package topdf

import "io"

// PDFServer is a server that converts files to PDFs.
type PDFServer interface {
	// Status checks PDFServer to see if it's running and ready.
	Status() (running bool, err error)

	// Convert converts file to pdf.
	// if file type is not supported or anyting related convert fails an err will be returned.
	Convert(name, extension string, file io.Reader) (pdf io.ReadCloser, err error)
}
