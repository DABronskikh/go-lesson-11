package app

import (
	"encoding/json"
	"fmt"
	"github.com/DABronskikh/go-lesson-11/cmd/bank/app/dto"
	"github.com/DABronskikh/go-lesson-11/pkg/card"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Server struct {
	cardSvc *card.Service
	mux     *http.ServeMux
}

func NewServer(cardSvc *card.Service, mux *http.ServeMux) *Server {
	return &Server{cardSvc: cardSvc, mux: mux}
}

func (s *Server) Init() {
	s.mux.HandleFunc("/getCards", s.getCards)
	s.mux.HandleFunc("/addCard", s.addCard)
	s.mux.HandleFunc("/editCard", s.editCard)
	s.mux.HandleFunc("/removeCard", s.removeCard)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) getCards(w http.ResponseWriter, r *http.Request) {
	userId, ok := url.Parse(r.URL.Query().Get("userId"))
	userIdstr := fmt.Sprintf("%v", userId)
	if ok != nil || userIdstr == "" {
		dtos := dto.CardErrDTO{Err: card.ErrUserIdNotFound.Error()}
		prepareResponseErr(w, r, dtos)
		return
	}

	cards := s.cardSvc.All(r.Context())
	dtos := []*dto.CardDTO{}

	for _, c := range cards {
		if userIdstr == c.UserId {
			dtos = append(
				dtos,
				&dto.CardDTO{
					Id:     c.Id,
					UserId: c.UserId,
					Number: c.Number,
					Type:   c.Type,
					System: c.System,
				})
		}
	}
	prepareResponse(w, r, dtos)
}

func (s *Server) addCard(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
	}

	params := &dto.CardDTO{}
	err = json.Unmarshal(body, params)
	if err != nil {
		log.Print(err)
	}
	card, err := s.cardSvc.Add(params.UserId, params.Type, params.System)

	if err != nil {
		dtos := dto.CardErrDTO{Err: err.Error()}
		prepareResponseErr(w, r, dtos)
		return
	}

	dtos := []*dto.CardDTO{}
	dtos = append(dtos,
		&dto.CardDTO{
			Id:     card.Id,
			UserId: card.UserId,
			Number: card.Number,
			Type:   card.Type,
			System: card.System,
		})

	prepareResponse(w, r, dtos)
}

func (s *Server) editCard(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (s *Server) removeCard(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func prepareResponse(w http.ResponseWriter, r *http.Request, dtos []*dto.CardDTO) {
	respBody, err := json.Marshal(dtos)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	// по умолчанию статус 200 Ok
	_, err = w.Write(respBody)
	if err != nil {
		log.Println(err)
	}
}

func prepareResponseErr(w http.ResponseWriter, r *http.Request, dtos dto.CardErrDTO) {
	respBody, err := json.Marshal(dtos)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(respBody)
	if err != nil {
		log.Println(err)
	}
}
