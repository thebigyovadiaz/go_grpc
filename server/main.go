package main

import (
	"context"
	"log"
	"net"

	"errors"

	"github.com/thebigyovadiaz/go_grpc/proto"
	"google.golang.org/grpc"
)

type Person struct {
	ID          int32
	Name        string
	Email       string
	PhoneNumber string
}

var nextID int32 = 1
var persons = make(map[int32]Person)

type server struct {
	proto.UnimplementedPersonServiceServer
}

func (s *server) Create(ctx context.Context, in *proto.CreatePersonRequest) (*proto.PersonProfileResponse, error) {
	person := Person{
		Name:        in.GetName(),
		Email:       in.GetEmail(),
		PhoneNumber: in.GetPhoneNumber(),
	}

	if person.Name == "" || person.Email == "" || person.PhoneNumber == "" {
		return &proto.PersonProfileResponse{}, errors.New("fields missing!")
	}

	person.ID = nextID
	persons[person.ID] = person
	nextID++

	return &proto.PersonProfileResponse{Id: person.ID, Name: person.Name, Email: person.Email}, nil
}

func (s *server) Read(ctx context.Context, in *proto.SinglePersonResponse) (*proto.PersonProfileResponse, error) {
	id := in.GetId()
	person := persons[id]

	if person.ID == 0 {
		return &proto.PersonProfileResponse{}, errors.New("not found!")
	}

	return &proto.PersonProfileResponse{Id: person.ID, Name: person.Name, Email: person.Email}, nil
}

func (s *server) Update(ctx context.Context, in *proto.UpdatePersonRequest) (*proto.SuccessResponse, error) {
	id := in.GetId()
	person := persons[id]

	if person.ID == 0 {
		return &proto.SuccessResponse{Response: "not found!"}, errors.New("not found!")
	}

	person.Name = in.GetName()
	person.Email = in.GetEmail()
	person.PhoneNumber = in.GetPhoneNumber()

	if person.Name == "" || person.Email == "" || person.PhoneNumber == "" {
		return &proto.SuccessResponse{Response: "fields missing!"}, errors.New("fields missing!")
	}

	persons[person.ID] = person

	return &proto.SuccessResponse{Response: "Updated successfully!"}, nil
}

func (s *server) Delete(ctx context.Context, in *proto.SinglePersonResponse) (*proto.SuccessResponse, error) {
	id := in.GetId()
	person := persons[id]

	if person.ID == 0 {
		return &proto.SuccessResponse{Response: "not found!"}, errors.New("not found!")
	}

	delete(persons, id)

	return &proto.SuccessResponse{Response: "Deleted successfully!"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterPersonServiceServer(s, &server{})
	log.Printf("gRPC server listening at %v", lis.Addr())

	if errServe := s.Serve(lis); errServe != nil {
		log.Fatalf("failed to serve: %v", errServe)
	}
}
