package handlers

import (
	"compress/gzip"
	"context"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"net/http"
	"strings"
	"zkKYC-backend/internal/app/config"
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

type MiddlewareHandler struct {
	cfg config.Config
}

// Create new instance of ZkKYCHandler
func NewMiddlewareHandler(c config.Config) *MiddlewareHandler {
	h := &MiddlewareHandler{
		cfg: c,
	}
	return h
}

// Unzip function
func (h *MiddlewareHandler) UnpackHandle(next http.Handler) http.Handler {
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
func (h *MiddlewareHandler) GzipHandle(next http.Handler) http.Handler {
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

// Middleware for jwt
func (h *MiddlewareHandler) JwtAuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[len("Bearer "):]
		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(h.cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*CustomClaims); ok {
			ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		}
	})
}
