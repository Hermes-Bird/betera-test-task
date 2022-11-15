package repos

import (
	"github.com/Hermes-Bird/betera-test-task/app/domain"
	"io"
	"time"
)

type ApodRepo interface {
	GetApodsList() ([]domain.Apod, error)
	GetApodsForDate(date time.Time) ([]domain.Apod, error)
	SaveApod(apods domain.Apod) error
	SaveApods(apods []domain.Apod) error
	io.Closer
}
