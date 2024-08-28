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

type Gym struct {
	db *sql.DB
}

func NewGym(db *sql.DB) *Gym {
	return &Gym{db: db}
}

func (s *Gym) CreateGym(gym *pb.CreateGymRequest) (*pb.CreateGymResponse, error) {
	query := `INSERT INTO sport_halls(name,owner_id,location,contact_number,latitude,longtitude,type_sport,type_gender)VALUES($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := s.db.Exec(query, gym.Name,gym.OwnerId, gym.Location, gym.ContactNumber, gym.Latitude, gym.Longtitude, gym.TypeSport, gym.TypeGender)
	if err != nil {
		return nil, errors.New("Gym was not created")
	}
	return nil, nil
}

func (s *Gym) UpdateGym(gym *pb.UpdateGymRequest) (*pb.UpdateGymResponse, error) {
	query := `UPDATE sport_halls SET `
	var condition []string
	var args []interface{}

	if gym.Name != "string" && gym.Name != "" {
		condition = append(condition, fmt.Sprintf("name = $%d", len(args)+1))
		args = append(args, gym.Name)
	}
	if gym.OwnerId != "string" && gym.OwnerId != ""{
		condition = append(condition, fmt.Sprintf("owner_id = $%d", len(args)+1))
		args = append(args, gym.OwnerId)
	}
	if gym.Location != "string" && gym.Location != "" {
		condition = append(condition, fmt.Sprintf("location = $%d", len(args)+1))
		args = append(args, gym.Location)
	}
	if gym.Longtitude != 0 {
		condition = append(condition, fmt.Sprintf("longtitude = $%d", len(args)+1))
		args = append(args, gym.Longtitude)
	}
	if gym.Latitude != 0 {
		condition = append(condition, fmt.Sprintf("latitude = $%d", len(args)+1))
		args = append(args, gym.Latitude)
	}
	if gym.ContactNumber != "string" && gym.ContactNumber != "" {
		condition = append(condition, fmt.Sprintf("contact_number = $%d", len(args)+1))
		args = append(args, gym.ContactNumber)
	}
	if gym.TypeSport != "string" && gym.TypeSport != "" {
		condition = append(condition, fmt.Sprintf("type_sport = $%d", len(args)+1))
		args = append(args, gym.TypeSport)
	}
	if gym.TypeGender != "string" && gym.TypeGender != "" {
		condition = append(condition, fmt.Sprintf("type_gender = $%d", len(args)+1))
		args = append(args, gym.TypeGender)
	}
	if len(condition) == 0 {
		return nil, errors.New("nothing to update")
	}
	condition = append(condition, "updated_at = now()")

	query += strings.Join(condition, ", ")
	query += fmt.Sprintf(" WHERE id = $%d", len(args)+1)
	args = append(args, gym.Id)

	_, err := s.db.Exec(query, args...)
	if err != nil {
		return nil, errors.New("Gym was not updated")
	}
	return nil, nil
}

func (s *Gym) DeleteGym(gym *pb.DeleteGymRequest) (*pb.DeleteGymResponse, error) {
	query := `
		UPDATE sport_halls
		SET deleted_at = $2
		WHERE id = $1 AND deleted_at = 0
	`
	_, err := s.db.Exec(query, gym.Id, time.Now().Unix())
	if err != nil {
		return nil, err
	}
	return &pb.DeleteGymResponse{}, nil
}

func (s *Gym) GetGym(gym *pb.GetGymRequest) (*pb.GetGymResponse, error) {
	query := `
		SELECT id, owner_id, name, longtitude, latitude, location, contact_number, type_sport, type_gender FROM sport_halls 
		WHERE id = $1 AND deleted_at = 0
	`
	row := s.db.QueryRow(query, gym.Id)

	user := pb.GetGymResponse{}
	err := row.Scan(
		&user.Id,
		&user.OwnerId,
		&user.Name,
		&user.Longtitude,
		&user.Latitude,
		&user.Location,
		&user.ContactNumber,
		&user.TypeSport,
		&user.TypeGender,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *Gym) ListGym(gym *pb.ListGymRequest) (*pb.ListGymResponse, error) {
	gyms := pb.ListGymResponse{}

	query := `
	SELECT
	 id,
	 owner_id,
	 name,
	 longtitude,
	 latitude,
	 location,
	 contact_number,
	 type_sport,
	 type_gender
	FROM 
	 sport_halls
	WHERE 
	 deleted_at = 0
	`

	var args []interface{}
	var conditions []string

	if gym.Name != "" && gym.Name != "string" {
		conditions = append(conditions, "LOWER(name) LIKE LOWER($"+strconv.Itoa(len(args)+1)+")")
		args = append(args, "%"+gym.Name+"%")
	}
	if gym.Location != "" && gym.Location != "string" {
		conditions = append(conditions, "LOWER(location) LIKE LOWER($"+strconv.Itoa(len(args)+1)+")")
		args = append(args, "%"+gym.Location+"%")
	}
	if gym.TypeSport != "" && gym.TypeSport != "string" {
		conditions = append(conditions, "LOWER(type_sport) LIKE LOWER($"+strconv.Itoa(len(args)+1)+")")
		args = append(args, "%"+gym.TypeSport+"%")
	}
	if gym.TypeGender != "" && gym.TypeGender != "string" {
		conditions = append(conditions, "LOWER(type_gender) LIKE LOWER($"+strconv.Itoa(len(args)+1)+")")
		args = append(args, "%"+gym.TypeGender+"%")
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	var limit int32 = 10
	var offset int32 = (gym.Page - 1) * limit

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err == sql.ErrNoRows {
		log.Println("gyms not found")
		return nil, errors.New("gym list is empty")
	}

	if err != nil {
		log.Println("Error while retrieving gyms: ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		gymResponse := pb.GetGymResponse{}

		err := rows.Scan(
			&gymResponse.Id,
			&gymResponse.OwnerId,
			&gymResponse.Name,
			&gymResponse.Longtitude,
			&gymResponse.Latitude,
			&gymResponse.Location,
			&gymResponse.ContactNumber,
			&gymResponse.TypeSport,
			&gymResponse.TypeGender,
		)

		if err != nil {
			log.Println("Error while scanning gym data: ", err)
			return nil, err
		}

		gyms.Gym = append(gyms.Gym, &gymResponse)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error after row iteration: ", err)
		return nil, err
	}

	log.Println("Successfully fetched all gyms")

	return &gyms, nil
}
