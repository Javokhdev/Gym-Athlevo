package postgres_test

import (
	"errors"
	pb "gym/genprotos"
	"gym/storage/postgres"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateGym(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	gymRepo := postgres.NewGym(db)
	createGymReq := &pb.CreateGymRequest{
		Name:          "Test Gym",
		Location:      "Test Location",
		ContactNumber: "998901234567",
		Latitude:      40.7128,
		Longtitude:    74.0060,
		TypeSport:     "Basketball",
		TypeGender:    "Male",
	}

	expectedQuery := `INSERT INTO sport_halls\(name,location,contact_number,latitude,longtitude,type_sport,type_gender\)VALUES\(\$1,\$2,\$3,\$4,\$5,\$6,\$7\)`

	mock.ExpectExec(expectedQuery).
		WithArgs(createGymReq.Name, createGymReq.Location, createGymReq.ContactNumber, createGymReq.Latitude, createGymReq.Longtitude, createGymReq.TypeSport, createGymReq.TypeGender).
		WillReturnResult(sqlmock.NewResult(1, 1))

	resp, err := gymRepo.CreateGym(createGymReq)

	assert.NoError(t, err)
	assert.Nil(t, resp)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCreateGym_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gymRepo := postgres.NewGym(db)

	createGymReq := &pb.CreateGymRequest{
		Name:          "Test Gym",
		Location:      "Test Location",
		ContactNumber: "998901234567",
		Latitude:      40.7128,
		Longtitude:    74.0060,
		TypeSport:     "Basketball",
		TypeGender:    "Male",
	}

	expectedQuery := `INSERT INTO sport_halls\(name,location,contact_number,latitude,longtitude,type_sport,type_gender\)VALUES\(\$1,\$2,\$3,\$4,\$5,\$6,\$7\)`

	mock.ExpectExec(expectedQuery).
		WithArgs(createGymReq.Name, createGymReq.Location, createGymReq.ContactNumber, createGymReq.Latitude, createGymReq.Longtitude, createGymReq.TypeSport, createGymReq.TypeGender).
		WillReturnError(errors.New("database error"))

	resp, err := gymRepo.CreateGym(createGymReq)

	assert.Error(t, err)
	assert.Nil(t, resp)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
