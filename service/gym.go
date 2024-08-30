package service

import (
	"context"
	"errors"

	pb "gym/genprotos"
	s "gym/storage"
)

type GymService struct {
	stg s.StorageI
	pb.UnimplementedGymServiceServer
}

func NewGymService(stg s.StorageI) *GymService {
	return &GymService{stg: stg}
}

func (s *GymService) CreateGym(c context.Context,gym *pb.CreateGymRequest) (*pb.CreateGymResponse, error){
	_, err := s.stg.Gym().CreateGym(gym)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *GymService) UpdateGym(c context.Context,gym *pb.UpdateGymRequest) (*pb.UpdateGymResponse, error){
	_, err := s.stg.Gym().UpdateGym(gym)
	if err != nil {
		return nil, errors.New("gym was not updated")
	}
	return nil, nil
}

func (s *GymService) DeleteGym(c context.Context,gym *pb.DeleteGymRequest) (*pb.DeleteGymResponse, error){
	_, err := s.stg.Gym().DeleteGym(gym)
	if err != nil {
		return nil, errors.New("gym was not deleted")
	}
	return nil, nil
}

func (s *GymService) GetGym(c context.Context,gym *pb.GetGymRequest) (*pb.GetGymResponse, error){
	res, err := s.stg.Gym().GetGym(gym)
	if err != nil {
		return nil, errors.New("gym not found")
	}
	return res, nil
}

func (s *GymService) ListGym(c context.Context,gym *pb.ListGymRequest) (*pb.ListGymResponse, error){
	res, err := s.stg.Gym().ListGym(gym)
	if err != nil {
		return nil, errors.New("gyms not found")
	}
	return res, nil
}