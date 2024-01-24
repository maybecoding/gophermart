package logger

import (
	"net/http"
	"time"
)

type (
	responseData struct {
		statusCode int
		contentLen int
	}

	proxyResponseWriter struct {
		wInter http.ResponseWriter
		resp   *responseData
	}
)

func newProxyResponseWriter(w http.ResponseWriter, resp *responseData) http.ResponseWriter {
	return &proxyResponseWriter{wInter: w, resp: resp}
}

func (w *proxyResponseWriter) Header() http.Header {
	return w.wInter.Header()
}
func (w *proxyResponseWriter) Write(b []byte) (int, error) {
	contentLen, err := w.wInter.Write(b)
	w.resp.contentLen = contentLen
	return contentLen, err
}
func (w *proxyResponseWriter) WriteHeader(statusCode int) {
	w.resp.statusCode = statusCode
	w.wInter.WriteHeader(statusCode)
}

func Handler(h http.Handler) http.Handler {
	handlerFn := func(w http.ResponseWriter, r *http.Request) {

		timeStart := time.Now()
		respData := responseData{}
		wProxy := newProxyResponseWriter(w, &respData)
		h.ServeHTTP(wProxy, r)
		lg.Debug().
			Str("URI", r.RequestURI).
			Dur("duration", time.Since(timeStart)).
			Str("method", r.Method).
			Int("content length", respData.contentLen).
			Int("status code", respData.statusCode).
			Msg("HTTP request handled")
	}

	return http.HandlerFunc(handlerFn)
}
