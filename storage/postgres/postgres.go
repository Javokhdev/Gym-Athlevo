package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"gym/config"
	postgres "gym/storage"

	_ "github.com/lib/pq"
)

type Storage struct {
	db        *sql.DB
	GymC      postgres.GymI
	FacilityC postgres.FacilityI
	UniqueC   postgres.UniqueI
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func ConnectDb() (*Storage, error) {
	cfg := config.Load()
	dbCon := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDatabase)
	db, err := sql.Open("postgres", dbCon)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	log.Println("Postgres connected succesfully")
	gymS := NewGym(db)
	FacilityS := NewFacility(db)
	UniqueS := NewUnique(db)

	return &Storage{
		GymC:      gymS,
		FacilityC: FacilityS,
		UniqueC:   UniqueS,
	}, nil
}

func (s *Storage) Gym() postgres.GymI {
	if s.GymC == nil {
		s.GymC = NewGym(s.db)
	}
	return s.GymC
}

func (s *Storage) Facility() postgres.FacilityI {
	if s.FacilityC == nil {
		s.FacilityC = NewFacility(s.db)
	}
	return s.FacilityC
}

func (s *Storage) Unique() postgres.UniqueI {
	if s.UniqueC == nil {
		s.UniqueC = NewUnique(s.db)
	}
	return s.UniqueC
}
