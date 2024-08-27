package storage

import (
	pb "gym/genprotos"
)

type StorageI interface {
	Gym() GymI
	Facility() FacilityI
	Unique() UniqueI
}

type GymI interface{
	CreateGym(gym *pb.CreateGymRequest) (*pb.CreateGymResponse, error)
	UpdateGym(gym *pb.UpdateGymRequest) (*pb.UpdateGymResponse, error)
	DeleteGym(gym *pb.DeleteGymRequest) (*pb.DeleteGymResponse, error)
	GetGym(gym *pb.GetGymRequest) (*pb.GetGymResponse, error)
	ListGym(gym *pb.ListGymRequest) (*pb.ListGymResponse, error)
}

type FacilityI interface{
	CreateFacility(gym *pb.CreateFacilityRequest) (*pb.CreateFacilityResponse, error)
	UpdateFacility(gym *pb.UpdateFacilityRequest) (*pb.UpdateFacilityResponse, error)
	DeleteFacility(gym *pb.DeleteFacilityRequest) (*pb.DeleteFacilityResponse, error)
	GetFacility(gym *pb.GetFacilityRequest) (*pb.GetFacilityResponse, error)
	ListFacility(gym *pb.ListFacilityRequest) (*pb.ListFacilityResponse, error)
}

type UniqueI interface{
	CreateUnique(gym *pb.CreateUniqueRequest) (*pb.CreateUniqueResponse, error)
	UpdateUnique(gym *pb.UpdateUniqueRequest) (*pb.UpdateUniqueResponse, error)
	DeleteUnique(gym *pb.DeleteUniqueRequest) (*pb.DeleteUniqueResponse, error)
	GetUnique(gym *pb.GetUniqueRequest) (*pb.GetUniqueResponse, error)
	ListUnique(gym *pb.ListUniqueRequest) (*pb.ListUniquesResponse, error)
}