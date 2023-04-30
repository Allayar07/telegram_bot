package server

import (
	"github.com/zhashkevych/go-pocket-sdk"
	"log"
	"net/http"
	"strconv"
	"telegramBot/internal/repository"
)

type AuthServer struct {
	server          *http.Server
	pocketClient    *pocket.Client
	tokenRepository repository.TokenRepository
	redirectURL     string
}

func NewAuthServer(pocketClient *pocket.Client, tr repository.TokenRepository, redirectURL string) *AuthServer {
	return &AuthServer{
		pocketClient:    pocketClient,
		tokenRepository: tr,
		redirectURL:     redirectURL,
	}
}

func (s *AuthServer) Run() error {
	s.server = &http.Server{
		Addr:    ":80",
		Handler: s,
	}

	return s.server.ListenAndServe()
}

func (s *AuthServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	chatIdPrm := r.URL.Query().Get("chat_id")
	if chatIdPrm == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chatId, err := strconv.ParseInt(chatIdPrm, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestToken, err := s.tokenRepository.Get(chatId, repository.RequestToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	authResp, err := s.pocketClient.Authorize(r.Context(), requestToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := s.tokenRepository.Save(chatId, authResp.AccessToken, repository.AccessToken); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("chat_id%d\nrequest_token%s\naccess_token%s\n", chatId, requestToken, authResp.AccessToken)

	w.Header().Add("Location", s.redirectURL)
	w.WriteHeader(http.StatusMovedPermanently)
}
