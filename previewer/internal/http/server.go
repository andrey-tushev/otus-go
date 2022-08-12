package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type Server struct {
	logger     Logger
	httpServer *http.Server
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

func New(logger Logger) *Server {
	return &Server{
		logger: logger,
	}
}

func (s *Server) Start(ctx context.Context, host string, port int) error {
	s.logger.Info("web-server start")

	s.httpServer = &http.Server{
		Addr:         net.JoinHostPort(host, strconv.Itoa(port)),
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		Handler:      s,
	}
	err := s.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.httpServer.Shutdown(ctx)
	s.logger.Info("web-server stopped")

	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.logger.Info(r.RequestURI)

	img, err := parse(r.RequestURI)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	targetUrl := "http://localhost:8082/" + img.Path
	s.logger.Info("target: " + targetUrl)
	resp, err := client.Get(targetUrl)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		text := fmt.Sprintf("bad target server response %d", resp.StatusCode)
		s.logger.Error(text)
		http.Error(w, text, http.StatusBadGateway)
		return
	}

	for name, _ := range resp.Header {
		w.Header().Set(name, resp.Header.Get(name))
	}
	w.Header().Set("X-Proxy", "proxy-resizer")
	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
}

type PreviewImage struct {
	Path   string
	Width  int
	Height int
}

var ErrBadImageRequestURL = errors.New("bad image request url")

const (
	MaxWidth  = 1024
	MaxHeight = 1024
)

var urlRexExp = regexp.MustCompile(`^\/fill\/(\d+)\/(\d+)/((?:[\/a-z\d\-\._])+\.jpe?g)$`)

func parse(uri string) (PreviewImage, error) {
	parts := urlRexExp.FindStringSubmatch(uri)
	if len(parts) != 3+1 {
		return PreviewImage{}, ErrBadImageRequestURL
	}

	w, _ := strconv.Atoi(parts[1])
	if w < 1 || w > MaxWidth {
		return PreviewImage{}, ErrBadImageRequestURL
	}

	h, _ := strconv.Atoi(parts[2])
	if h < 1 || h > MaxHeight {
		return PreviewImage{}, ErrBadImageRequestURL
	}

	return PreviewImage{
		Path:   parts[3],
		Width:  w,
		Height: h,
	}, nil
}
