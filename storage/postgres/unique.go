package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	pb "gym/genprotos"
	"log"
	"strings"
	"time"
)

type Unique struct {
	db *sql.DB
}

func NewUnique(db *sql.DB) *Unique {
	return &Unique{db: db}
}

func (s *Unique) CreateUnique(unique *pb.CreateUniqueRequest) (*pb.CreateUniqueResponse, error) {
	query := `INSERT INTO gym_facility(sport_halls_id,facility_id,count)VALUES($1,$2,$3)`
	_, err := s.db.Exec(query, unique.SportHallsId, unique.FacilityId, unique.Count)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
func (s *Unique) UpdateUnique(unique *pb.UpdateUniqueRequest) (*pb.UpdateUniqueResponse, error) {
	query := `UPDATE gym_facility SET `
	var condition []string
	var args []interface{}

	if unique.SportHallsId != "string" && unique.SportHallsId != "" {
		condition = append(condition, fmt.Sprintf("sport_halls_id = $%d", len(args)+1))
		args = append(args, unique.SportHallsId)
	}
	if unique.FacilityId != "string" && unique.FacilityId != "" {
		condition = append(condition, fmt.Sprintf("facility_id = $%d", len(args)+1))
		args = append(args, unique.FacilityId)
	}
	if unique.Count != 0 {
		condition = append(condition, fmt.Sprintf("count = $%d", len(args)+1))
		args = append(args, unique.Count)
	}
	if len(condition) == 0 {
		return nil, errors.New("nothing to update")
	}
	condition = append(condition, "updated_at = now()")

	query += strings.Join(condition, ", ")
	query += fmt.Sprintf(" WHERE id = $%d", len(args)+1)
	args = append(args, unique.Id)

	_, err := s.db.Exec(query, args...)
	if err != nil {
		return nil, errors.New("Unique was not updated")
	}
	return nil, nil
}

func (s *Unique) DeleteUnique(unique *pb.DeleteUniqueRequest) (*pb.DeleteUniqueResponse, error) {
	query := `
		UPDATE gym_facility
		SET deleted_at = $2
		WHERE id = $1 AND deleted_at = 0
	`
	_, err := s.db.Exec(query, unique.Id, time.Now().Unix())
	if err != nil {
		return nil, err
	}
	return &pb.DeleteUniqueResponse{}, nil
}

func (s *Unique) GetUnique(unique *pb.GetUniqueRequest) (*pb.GetUniqueResponse, error) {
	query := `
		SELECT id, sport_halls_id, facility_id, count FROM gym_facility 
		WHERE id = $1 AND deleted_at = 0
	`
	row := s.db.QueryRow(query, unique.Id)

	user := pb.GetUniqueResponse{}
	err := row.Scan(
		&user.Id,
		&user.SportHallsId,
		&user.FacilityId,
		&user.Count,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("no rows in result set")
		}
		return nil, err
	}

	return &user, nil
}

func (s *Unique) ListUnique(unique *pb.ListUniqueRequest) (*pb.ListUniquesResponse, error) {
	const limit = 10
	offset := int32(0)

	if unique.Page > 0 {
		offset = (unique.Page - 1) * limit
	}

	query := `
		SELECT id, sport_halls_id, facility_id, count
		FROM gym_facility
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		log.Printf("Failed to list uniques: %v", err)
		return nil, fmt.Errorf("failed to list uniques: %v", err)
	}
	defer rows.Close()

	var uniques []*pb.GetUniqueResponse
	for rows.Next() {
		var res pb.GetUniqueResponse
		if err := rows.Scan(
			&res.Id,
			&res.SportHallsId,
			&res.FacilityId,
			&res.Count,
		); err != nil {
			log.Printf("Failed to scan unique row: %v", err)
			return nil, fmt.Errorf("failed to scan unique row: %v", err)
		}
		uniques = append(uniques, &res)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Failed to iterate over unique rows: %v", err)
		return nil, fmt.Errorf("failed to iterate over unique rows: %v", err)
	}

	return &pb.ListUniquesResponse{
		Uniques: uniques,
	}, nil
}
