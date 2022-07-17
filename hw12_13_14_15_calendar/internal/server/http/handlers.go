package internalhttp

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/andrey-tushev/otus-go/hw12_13_14_15_calendar/internal/app"
)

func (s *Server) ListEvents(w http.ResponseWriter, r *http.Request) {
	list, err := s.app.ListEvents(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, _ := json.Marshal(list)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(body)
}

func (s *Server) CreateEvent(w http.ResponseWriter, r *http.Request) {
	event, err := readEvent(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newID, err := s.app.CreateEvent(context.Background(), event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	event.ID = newID

	responseBody, _ := json.Marshal(event)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(responseBody)
}

func (s *Server) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	event, err := readEvent(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.app.UpdateEvent(context.Background(), event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseBody, _ := json.Marshal(event)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(responseBody)
}

func readEvent(request *http.Request) (app.Event, error) {
	requestBody, err := io.ReadAll(request.Body)
	if err != nil {
		return app.Event{}, err
	}

	var event app.Event
	err = json.Unmarshal(requestBody, &event)
	if err != nil {
		return app.Event{}, err
	}

	return event, err
}
