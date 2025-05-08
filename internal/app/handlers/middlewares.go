package handlers

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// Wrapper for io.Writer
func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

type gzipReader struct {
	Body   io.ReadCloser
	Reader *gzip.Reader
}

// Wrapper for io.Reader
func (r gzipReader) Read(p []byte) (n int, err error) {
	return r.Reader.Read(p)
}

// Wrapper for body closing
func (r gzipReader) Close() error {
	if err := r.Body.Close(); err != nil {
		return err
	}

	if err := r.Reader.Close(); err != nil {
		return err
	}

	return nil
}

// Unzip function
func UnpackHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(r.Body)
			gzr := gzipReader{
				Body:   r.Body,
				Reader: gz,
			}
			if err != nil {
				io.WriteString(w, err.Error())
				return
			}
			r.Body = gzr
			defer gzr.Close()
		}
		next.ServeHTTP(w, r)
	})
}

// Zip function
func GzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")

		next.ServeHTTP(gzipWriter{
			ResponseWriter: w,
			Writer:         gz,
		}, r)
	})
}
