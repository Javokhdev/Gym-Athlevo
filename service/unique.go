package service

import (
	"context"
	"errors"

	pb "gym/genprotos"
	s "gym/storage"
)

type UniqueService struct {
	stg s.StorageI
	pb.UnimplementedUniqueServiceServer
}

func NewUniqueService(stg s.StorageI) *UniqueService {
	return &UniqueService{stg: stg}
}

func (s *UniqueService) CreateUnique(c context.Context,Unique *pb.CreateUniqueRequest) (*pb.CreateUniqueResponse, error){
	_, err := s.stg.Unique().CreateUnique(Unique)
	if err != nil {
		return nil, errors.New("unique was not created")
	}
	return nil, nil
}

func (s *UniqueService) UpdateUnique(c context.Context,Unique *pb.UpdateUniqueRequest) (*pb.UpdateUniqueResponse, error){
	_, err := s.stg.Unique().UpdateUnique(Unique)
	if err != nil {
		return nil, errors.New("unique was not updated")
	}
	return nil, nil
}

func (s *UniqueService) DeleteUnique(c context.Context,Unique *pb.DeleteUniqueRequest) (*pb.DeleteUniqueResponse, error){
	_, err := s.stg.Unique().DeleteUnique(Unique)
	if err != nil {
		return nil, errors.New("unique was not deleted")
	}
	return nil, nil
}

func (s *UniqueService) GetUnique(c context.Context,Unique *pb.GetUniqueRequest) (*pb.GetUniqueResponse, error){
	res, err := s.stg.Unique().GetUnique(Unique)
	if err != nil {
		return nil, errors.New("unique not found")
	}
	return res, nil
}

func (s *UniqueService) ListUnique(c context.Context,Unique *pb.ListUniqueRequest) (*pb.ListUniquesResponse, error){
	res, err := s.stg.Unique().ListUnique(Unique)
	if err != nil {
		return nil, errors.New("uniques not found")
	}
	return res, nil
}