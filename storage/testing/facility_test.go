package postgres_test

import (
	"errors"
	pb "gym/genprotos"
	"gym/storage/postgres"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateFacility(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	gymRepo := postgres.NewFacility(db)
	createGymReq := &pb.CreateFacilityRequest{
		Name:        "Test Gym",
		Type:        "Test Location",
		Image:       "998901234567",
		Description: "Basketball",
	}

	expectedQuery := `INSERT INTO facility(name,type,image,description) VALUES ($1, $2, $3, $4)`

	mock.ExpectExec(expectedQuery).
		WithArgs(createGymReq.Name, createGymReq.Type, createGymReq.Image, createGymReq.Description).
		WillReturnResult(sqlmock.NewResult(1, 1))

	resp, err := gymRepo.CreateFacility(createGymReq)

	assert.NoError(t, err)
	assert.NotNil(t, resp) 

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}


func TestCreateFacility_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gymRepo := postgres.NewFacility(db)

	createGymReq := &pb.CreateFacilityRequest{
		Name:        "Test Gym",
		Type:        "Test Location",
		Image:       "998901234567",
		Description: "Basketball",
	}

	expectedQuery := `INSERT INTO facility(name,type,image,description) VALUES ($1, $2, $3, $4)`

	mock.ExpectExec(expectedQuery).
		WithArgs(createGymReq.Name, createGymReq.Type, createGymReq.Image, createGymReq.Description).
		WillReturnError(errors.New("database error"))

	resp, err := gymRepo.CreateFacility(createGymReq)

	assert.Error(t, err)
	assert.Nil(t, resp)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
