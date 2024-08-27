package main

import (
	"log"
	"net"

	pb "gym/genprotos"
	"gym/service"
	"gym/storage/postgres"
	"google.golang.org/grpc"
)

func main(){
	db, err := postgres.ConnectDb()
	if err != nil {
		log.Fatal(err)
	}

	liss, err := net.Listen("tcp", ":50001")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterGymServiceServer(s, service.NewGymService(db))
	pb.RegisterFacilityServiceServer(s, service.NewFacilityService(db))
	pb.RegisterUniqueServiceServer(s, service.NewUniqueService(db))

	log.Printf("server listening at %v", liss.Addr())
	if err := s.Serve(liss); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
}
