package http

import (
	"bytes"
	"context"
	"fmt"
	"image/jpeg"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/disintegration/imaging"

	"github.com/andrey-tushev/otus-go/previewer/internal/cache"
	"github.com/andrey-tushev/otus-go/previewer/internal/preview"
)

type Server struct {
	logger       Logger
	httpServer   *http.Server
	cache        cache.Cache
	targetPrefix string
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

func New(logger Logger, targetPrefix string) *Server {
	return &Server{
		logger:       logger,
		cache:        cache.New(),
		targetPrefix: targetPrefix,
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
	// Получам параметры требуемой превьюшки
	requestedPreview, err := preview.NewFromURL(r.RequestURI)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	// Попробуем найти в кэше и вернуть из кэша
	cachedContent := s.cache.Get(requestedPreview)
	if cachedContent != nil {
		w.Header().Set("Content-Length", strconv.Itoa(len(cachedContent)))
		w.Header().Set("X-Proxy", "proxy-resizer")
		io.Copy(w, bytes.NewReader(cachedContent))

		s.logger.Info("taken from cache")
		return
	}

	// Запросим оригинальную картинку из целевого сервера
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	targetUrl := s.targetPrefix + requestedPreview.Path
	s.logger.Info("target: " + targetUrl)
	targetResp, err := client.Get(targetUrl)

	if targetResp.StatusCode < 200 || targetResp.StatusCode > 299 {
		text := fmt.Sprintf("bad target server response %d", targetResp.StatusCode)
		s.logger.Error(text)
		http.Error(w, text, http.StatusBadGateway)
		return
	}

	// Копируем заголовки
	for name, _ := range targetResp.Header {
		w.Header().Set(name, targetResp.Header.Get(name))
	}

	// Ресайзим
	resizedContent, err := resize(targetResp.Body, requestedPreview.Width, requestedPreview.Height)
	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	// Записываем уменьшенную картику в ответ и в кэш
	w.WriteHeader(targetResp.StatusCode)
	w.Header().Set("Content-Length", strconv.Itoa(len(resizedContent)))
	w.Header().Set("X-Proxy", "proxy-resizer")

	io.Copy(w, bytes.NewReader(resizedContent))
	s.cache.Set(requestedPreview, resizedContent)
}

func resize(content io.ReadCloser, width, height int) ([]byte, error) {
	targetImage, err := imaging.Decode(content)
	if err != nil {
		return nil, err
	}
	resizedImage := imaging.Fit(targetImage, width, height, imaging.Lanczos)
	buf := bytes.Buffer{}
	err = jpeg.Encode(&buf, resizedImage, nil)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
