package repos

import (
	"database/sql"
	"fmt"
	"github.com/Hermes-Bird/betera-test-task/app/config"
	"github.com/Hermes-Bird/betera-test-task/app/domain"
	_ "github.com/lib/pq"
	"log"
	"strings"
	"time"
)

type ApodPostgresRepo struct {
	Db *sql.DB
}

func NewApodPostgresRepo(cfg *config.Config) ApodRepo {
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgesDBName)
	conn, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = conn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Exec(`
		CREATE table IF NOT exists apods (
			apod_id serial PRIMARY KEY,
			apod_date DATE NOT NULL,
			apod_url TEXT NOT NULL
		)`)

	if err != nil {
		log.Fatal(err)
	}

	return &ApodPostgresRepo{
		Db: conn,
	}
}

func (a *ApodPostgresRepo) GetApodsList() ([]domain.Apod, error) {
	res, err := a.Db.Query("select apod_id, apod_date, apod_url from apods")
	if err != nil {
		return nil, err
	}
	defer res.Close()

	return getApodsFromRows(res)
}

func (a *ApodPostgresRepo) GetApodsForDate(date time.Time) ([]domain.Apod, error) {
	res, err := a.Db.Query("select apod_id, apod_date, apod_url from apods where apod_date=$1", date)
	if err != nil {
		return nil, err
	}

	return getApodsFromRows(res)
}

func (a *ApodPostgresRepo) SaveApods(apods []domain.Apod) error {
	apodValues := getInsertApodData(apods)

	_, err := a.Db.Exec(fmt.Sprintf("insert into apods(apod_date, apod_url) values %s", apodValues))
	return err
}

func (r ApodPostgresRepo) SaveApod(apod domain.Apod) error {
	_, err := r.Db.Exec("insert into apods(apod_date, apod_url) values ($1, $2)", apod.Date, apod.Url)

	return err
}

func getApodsFromRows(rows *sql.Rows) ([]domain.Apod, error) {
	apods := make([]domain.Apod, 0)
	for rows.Next() {
		var apod domain.Apod

		err := rows.Scan(&apod.Id, &apod.Date, &apod.Url)
		if err != nil {
			return nil, err
		}

		apods = append(apods, apod)
	}

	return apods, nil
}

func getInsertApodData(apods []domain.Apod) string {
	b := strings.Builder{}
	for i, apod := range apods {
		date := apod.Date.Format("2006-01-02")
		b.WriteString(fmt.Sprintf("('%s', '%s')", date, apod.Url))
		if i != len(apods)-1 {
			b.WriteString(",")
		}
	}

	return b.String()
}

func (a *ApodPostgresRepo) Close() error {
	return a.Db.Close()
}
