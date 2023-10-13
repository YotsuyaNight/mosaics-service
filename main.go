package main

import (
	"context"
	"fmt"
	"gomosaics"
	"log"
	"net"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	pb "mosaics-processing-service/proto"
)

type FileProcessingServerImpl struct {
	pb.UnimplementedFileProcessingServer
}

func (s FileProcessingServerImpl) ProcessFile(c context.Context, req *pb.ProcessFileReq) (*pb.ProcessFileRepl, error) {
	inputFilepath := fmt.Sprintf("../mosaics-web/%s", req.Filepath)
	outputFilepath := fmt.Sprintf("%s.png", uuid.NewString())
	err := gomosaics.Mosaicate(inputFilepath, "../go-mosaics/icons_small/", outputFilepath, 4, 16)
	if err != nil {
		return nil, err
	}
	return &pb.ProcessFileRepl{
		Filepath: outputFilepath,
	}, nil
}

func main() {
	log.Println("Running mosaics-processing-service server")
	lis, err := net.Listen("tcp", "0.0.0.0:8081")
	if err != nil {
		log.Fatal("Could not bind network interface: ", err)
	}
	srv := grpc.NewServer()
	pb.RegisterFileProcessingServer(srv, FileProcessingServerImpl{})
	srv.Serve(lis)
}
