package service

import (
	"context"
	"errors"

	pb "gym/genprotos"
	s "gym/storage"
)

type FacilityService struct {
	stg s.StorageI
	pb.UnimplementedFacilityServiceServer
}

func NewFacilityService(stg s.StorageI) *FacilityService {
	return &FacilityService{stg: stg}
}

func (s *FacilityService) CreateFacility(c context.Context,Facility *pb.CreateFacilityRequest) (*pb.CreateFacilityResponse, error){
	_, err := s.stg.Facility().CreateFacility(Facility)
	if err != nil {
		return nil, errors.New("facility was not created")
	}
	return nil, nil
}

func (s *FacilityService) UpdateFacility(c context.Context,Facility *pb.UpdateFacilityRequest) (*pb.UpdateFacilityResponse, error){
	_, err := s.stg.Facility().UpdateFacility(Facility)
	if err != nil {
		return nil, errors.New("facility was not updated")
	}
	return nil, nil
}

func (s *FacilityService) DeleteFacility(c context.Context,Facility *pb.DeleteFacilityRequest) (*pb.DeleteFacilityResponse, error){
	_, err := s.stg.Facility().DeleteFacility(Facility)
	if err != nil {
		return nil, errors.New("facility was not deleted")
	}
	return nil, nil
}

func (s *FacilityService) GetFacility(c context.Context,Facility *pb.GetFacilityRequest) (*pb.GetFacilityResponse, error){
	res, err := s.stg.Facility().GetFacility(Facility)
	if err != nil {
		return nil, errors.New("facility not found")
	}
	return res, nil
}

func (s *FacilityService) ListFacility(c context.Context,Facility *pb.ListFacilityRequest) (*pb.ListFacilityResponse, error){
	res, err := s.stg.Facility().ListFacility(Facility)
	if err != nil {
		return nil, errors.New("facilities not found")
	}
	return res, nil
}