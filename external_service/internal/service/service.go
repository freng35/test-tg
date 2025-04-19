package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"ext_service/internal/config"
	tgservice "ext_service/internal/interconnect/tg_service"
)

type RestService struct {
	cfg             config.Rest
	mux             *chi.Mux
	tgServiceClient *tgservice.Client
}

type memberRequest struct {
	BotToken   string `json:"bot_token"`
	ChannelURL string `json:"channel_url"`
	UserID     int    `json:"user_id"`
}

func NewRestService(ctx context.Context, cfg config.Rest, tgServiceClient *tgservice.Client) (*RestService, error) {
	mux := chi.NewRouter()
	return &RestService{
		cfg:             cfg,
		mux:             mux,
		tgServiceClient: tgServiceClient,
	}, nil
}

func (s *RestService) Run(ctx context.Context) error {
	s.mux.HandleFunc("/api/check_user", s.checkUser)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port),
		Handler: s.mux,
	}

	serveChan := make(chan error, 1)
	go func() {
		serveChan <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		ctxT, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		return srv.Shutdown(ctxT)

	case err := <-serveChan:
		return err
	}
}

func (s *RestService) checkUser(w http.ResponseWriter, r *http.Request) {
	var req memberRequest
	dec := json.NewDecoder(r.Body)

	err := dec.Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	if req.BotToken == "" || req.ChannelURL == "" || req.UserID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := s.tgServiceClient.ChekMembership(r.Context(), req.BotToken, req.ChannelURL, int64(req.UserID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	res := fmt.Sprintf("is_member %v, added_to_db %v", resp.IsMember, resp.AddedToDb)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(res))
}
