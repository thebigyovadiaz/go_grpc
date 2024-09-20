package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/thebigyovadiaz/go_grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func executeClient() {
	// Connect to the gRPC server
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewPersonServiceClient(conn)

	// Timeout for context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	// Example 1: Create a new person
	fmt.Println("Creating a new person")
	createReq := &proto.CreatePersonRequest{
		Name:        "Pepe Pepito",
		Email:       "pepe-pepito@sydesign.com",
		PhoneNumber: "123-456-7890",
	}

	createResp, err := client.Create(ctx, createReq)
	if err != nil {
		log.Fatalf("Error during create: %v", err)
	}

	fmt.Printf("Created response: %+v\n\n", createResp)

	// Example 2: Read a person
	fmt.Println("Reading person by ID...")
	readReq := &proto.SinglePersonResponse{
		Id: createResp.Id,
	}

	readResp, err := client.Read(ctx, readReq)
	if err != nil {
		log.Fatalf("Error during get: %v", err)
	}

	fmt.Printf("Read response: %+v\n\n", readResp)

	// Example 3: Update a person
	fmt.Println("Updating a person's details")
	updateReq := &proto.UpdatePersonRequest{
		Id:          createResp.Id,
		Name:        "Updated Name",
		Email:       "updated-name@sydesign.com",
		PhoneNumber: "123-222",
	}

	updateResp, err := client.Update(ctx, updateReq)
	if err != nil {
		log.Fatalf("Error during update: %v", err)
	}

	fmt.Printf("Updated response: %+v\n\n", updateResp)

	// Example 4: Delete a person
	fmt.Println("Deleting person by ID...")
	deleteReq := &proto.SinglePersonResponse{
		Id: createResp.Id,
	}

	deleteResp, err := client.Delete(ctx, deleteReq)
	if err != nil {
		log.Fatalf("Error during delete: %v", err)
	}

	fmt.Printf("Deleted response: %+v\n\n", deleteResp)
}

func main() {
	fmt.Printf("Executing the client gRPC...\n\n")
	executeClient()
}
