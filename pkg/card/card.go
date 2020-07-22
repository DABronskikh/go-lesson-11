package card

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrMapTypeNotFound   = errors.New("map type not found")
	ErrMapSystemNotFound = errors.New("map system not found")
	ErrUserIdNotFound    = errors.New("user ID not found")
	ErrNoBaseCard        = errors.New("no base card")
)

type Card struct {
	Id     int64
	UserId string
	Number int64
	Type   string
	System string
}

type Service struct {
	mu     sync.RWMutex
	cards  []*Card
	mun    sync.Mutex
	number int64
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) All(context.Context) []*Card {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cards
}

func (s *Service) Add(userId string, typeCard string, systemCard string) (*Card, error) {
	number := s.getNumber()
	err := getTypeCard(typeCard)
	if err != nil {
		return &Card{}, err
	}
	err = getSystemCard(systemCard)
	if err != nil {
		return &Card{}, err
	}
	err = s.isGotBaseCard(userId)
	if err != nil && typeCard != "basic" {
		return &Card{}, err
	}
	card := &Card{
		Id:     number,
		UserId: userId,
		Number: number,
		Type:   typeCard,
		System: systemCard,
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cards = append(s.cards, card)
	return card, nil
}

func getTypeCard(typeCard string) error {
	typesCard := []string{"basic", "additional", "virtual"}
	for _, value := range typesCard {
		if value == typeCard {
			return nil
		}
	}
	return ErrMapTypeNotFound
}

func getSystemCard(systemCard string) error {
	systemsCard := []string{"Visa", "MasterCard", "Mir"}
	for _, value := range systemsCard {
		if value == systemCard {
			return nil
		}
	}
	return ErrMapSystemNotFound
}

func (s *Service) getNumber() int64 {
	s.mun.Lock()
	defer s.mun.Unlock()
	s.number += 1
	return s.number
}

func (s *Service) isGotBaseCard(userId string) error {
	for _, value := range s.cards {
		if value.UserId == userId {
			return nil
		}
	}
	return ErrNoBaseCard
}
