package service

import (
	"context"
	"errors"
	"log/slog"

	"github/Igo87/crypt/models"
)

var ErrNilPointerRecevied = errors.New("nil pointer received")
var ErrFailedToFetchData = errors.New("failed to fetch data")
var ErrFailedToFetchDataFromDB = errors.New("failed to fetch data by today from DB")
var ErrFailedToFetchDataFromDBByLastDay = errors.New("failed to fetch data by last day from DB")
var ErrFailedToInsertData = errors.New("failed to insert data to DB")
var ErrNameNotFound = errors.New("name not found")

//go:generate mockgen -source=service.go -destination=mocks/mock_service.go -package=mocks
type repositoryManager interface {
	AddCurrencies(ctx context.Context, data models.Crypto) error
	SelectCurrenciesByToday(ctx context.Context) (models.Currencies, error)
	SelectCurrenciesByLastDay(ctx context.Context) (models.Currencies, error)
	Fetch(ctx context.Context, l slog.Logger) (models.Crypto, error)
	Run(ctx context.Context, l slog.Logger) error
	SelectCurrenciesByName(ctx context.Context, name string) (models.Currencies, error)
}

type Service struct {
	repo repositoryManager
}

// NewService returns a new instance of the Service struct.
// It takes a repositoryManager as an argument and returns a pointer to a Service.
func NewService(repo repositoryManager) *Service {
	return &Service{repo: repo}
}

func (s *Service) Run(ctx context.Context, l *slog.Logger) error {
	if s == nil || s.repo == nil {
		return ErrNilPointerRecevied
	}

	data, err := s.repo.Fetch(ctx, *l)
	if err != nil {
		l.Error(ErrFailedToFetchData.Error())
		return ErrFailedToFetchData
	}

	err = s.InsertData(ctx, data)
	if err != nil {
		l.Error(ErrFailedToInsertData.Error())
		return ErrFailedToInsertData
	}

	return nil
}

func (s *Service) GetData(ctx context.Context, log slog.Logger) (models.Crypto, error) {
	if s == nil {
		return models.Crypto{}, ErrNilPointerRecevied
	}

	if s.repo == nil {
		log.Error(ErrNilPointerRecevied.Error())
		return models.Crypto{}, ErrNilPointerRecevied
	}

	data, err := s.repo.Fetch(ctx, log)
	if err != nil {
		log.Error(ErrFailedToFetchData.Error())
		return models.Crypto{}, ErrFailedToFetchData
	}
	return data, nil
}

func (s *Service) InsertData(ctx context.Context, data models.Crypto) error {
	if s.repo == nil {
		return ErrNilPointerRecevied
	}

	return s.repo.AddCurrencies(ctx, data)
}

func (s *Service) GetDataByToday(ctx context.Context) (models.Currencies, error) {
	if s == nil || s.repo == nil {
		return models.Currencies{}, ErrNilPointerRecevied
	}

	data, err := s.repo.SelectCurrenciesByToday(ctx)
	if err != nil {
		return models.Currencies{}, ErrFailedToFetchDataFromDB
	}

	return data, nil
}
func (s *Service) GetDataByLastDay(ctx context.Context) (models.Currencies, error) {
	if s == nil {
		return models.Currencies{}, ErrNilPointerRecevied
	}

	if s.repo == nil {
		return models.Currencies{}, ErrNilPointerRecevied
	}

	data, err := s.repo.SelectCurrenciesByLastDay(ctx)
	if err != nil {
		return models.Currencies{}, ErrFailedToFetchDataFromDBByLastDay
	}
	return data, nil
}

func (s *Service) GetDataByName(ctx context.Context, name string) (models.Currencies, error) {
	if s == nil || s.repo == nil || name == "" {
		return models.Currencies{}, ErrNilPointerRecevied
	}

	data, err := s.repo.SelectCurrenciesByName(ctx, name)
	if err != nil {
		return models.Currencies{}, ErrFailedToFetchData
	}

	return data, nil
}
