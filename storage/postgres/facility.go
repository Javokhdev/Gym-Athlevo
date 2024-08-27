package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	pb "gym/genprotos"
	"log"
	"strconv"
	"strings"
	"time"
)

type Facility struct {
	db *sql.DB
}

func NewFacility(db *sql.DB) *Facility {
	return &Facility{db: db}
}

func (s *Facility) CreateFacility(facility *pb.CreateFacilityRequest) (*pb.CreateFacilityResponse, error) {
	query := `INSERT INTO facility(name,type,image,description)VALUES($1,$2,$3,$4)`
	_, err := s.db.Exec(query, facility.Name, facility.Type, facility.Image, facility.Description)
	if err != nil {
		return nil, errors.New("Facility was not created")
	}
	return nil, nil
}

func (s *Facility) UpdateFacility(facility *pb.UpdateFacilityRequest) (*pb.UpdateFacilityResponse, error) {
	query := `UPDATE facility SET `
	var condition []string
	var args []interface{}

	if facility.Name != "string" && facility.Name != "" {
		condition = append(condition, fmt.Sprintf("name = $%d", len(args)+1))
		args = append(args, facility.Name)
	}
	if facility.Type != "string" && facility.Type != "" {
		condition = append(condition, fmt.Sprintf("type = $%d", len(args)+1))
		args = append(args, facility.Type)
	}
	if facility.Image != "string" && facility.Image != "" {
		condition = append(condition, fmt.Sprintf("image = $%d", len(args)+1))
		args = append(args, facility.Image)
	}
	if facility.Description != "string" && facility.Description != "" {
		condition = append(condition, fmt.Sprintf("description = $%d", len(args)+1))
		args = append(args, facility.Description)
	}
	if len(condition) == 0 {
		return nil, errors.New("nothing to update")
	}
	condition = append(condition, "updated_at = now()")

	query += strings.Join(condition, ", ")
	query += fmt.Sprintf(" WHERE id = $%d", len(args)+1)
	args = append(args, facility.Id)

	_, err := s.db.Exec(query, args...)
	if err != nil {
		return nil, errors.New("Facility was not updated")
	}
	return nil, nil
}

func (s *Facility) DeleteFacility(facility *pb.DeleteFacilityRequest) (*pb.DeleteFacilityResponse, error) {
	query := `
		UPDATE facility
		SET deleted_at = $2
		WHERE id = $1 AND deleted_at = 0
	`
	_, err := s.db.Exec(query, facility.Id, time.Now().Unix())
	if err != nil {
		return nil, err
	}
	return &pb.DeleteFacilityResponse{}, nil
}

func (s *Facility) GetFacility(facility *pb.GetFacilityRequest) (*pb.GetFacilityResponse, error) {
	query := `
		SELECT id, name, type, image, description FROM facility 
		WHERE id = $1 AND deleted_at = 0
	`
	row := s.db.QueryRow(query, facility.Id)

	user := pb.GetFacilityResponse{}
	err := row.Scan(
		&user.Id,
		&user.Name,
		&user.Type,
		&user.Image,
		&user.Description,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *Facility) ListFacility(facility *pb.ListFacilityRequest) (*pb.ListFacilityResponse, error) {
	Facilitys := pb.ListFacilityResponse{}

	query := `
	SELECT
	 id,
	 name,
	 type,
	 image,
	 description
	FROM 
	 facility
	WHERE 
	 deleted_at = 0
	`

	var args []interface{}
	var conditions []string

	if facility.Name != "" && facility.Name != "string" {
		conditions = append(conditions, "LOWER(name) LIKE LOWER($"+strconv.Itoa(len(args)+1)+")")
		args = append(args, "%"+facility.Name+"%")
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	var limit int32 = 10
	var offset int32 = (facility.Page - 1) * limit

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err == sql.ErrNoRows {
		log.Println("Facilitys not found")
		return nil, errors.New("facility list is empty")
	}

	if err != nil {
		log.Println("Error while retrieving Facilitys: ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		FacilityResponse := pb.GetFacilityResponse{}

		err := rows.Scan(
			&FacilityResponse.Id,
			&FacilityResponse.Name,
			&FacilityResponse.Type,
			&FacilityResponse.Image,
			&FacilityResponse.Description,
		)

		if err != nil {
			log.Println("Error while scanning facility data: ", err)
			return nil, err
		}

		Facilitys.Facility = append(Facilitys.Facility, &FacilityResponse)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error after row iteration: ", err)
		return nil, err
	}

	log.Println("Successfully fetched all Facilitys")

	return &Facilitys, nil
}
