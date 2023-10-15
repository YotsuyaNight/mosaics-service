package main

import (
	"context"
	"gomosaics"
	"log"
	"net"
	"os"
	"path"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	pb "mosaics-processing-service/proto"
)

func getEnvOr(name string, def string) string {
	env := os.Getenv(name)
	if env == "" {
		return def
	}
	return env
}

func getBaseDir() string {
	return getEnvOr("BASE_DIR", ".")
}

func getResultDir() string {
	return getEnvOr("RESULT_DIR", "result")
}

func getIconsDir() string {
	return getEnvOr("ICONS_DIR", "icons")
}

func prepareDirs() {
	os.MkdirAll(path.Join(getBaseDir(), getIconsDir()), os.ModePerm)
	os.MkdirAll(path.Join(getBaseDir(), getResultDir()), os.ModePerm)
}

type FileProcessingServerImpl struct {
	pb.UnimplementedFileProcessingServer
}

func (s FileProcessingServerImpl) ProcessFile(c context.Context, req *pb.ProcessFileReq) (*pb.ProcessFileRepl, error) {
	inputFilepath := path.Join(getBaseDir(), req.Filepath)
	outputFilename := uuid.NewString() + ".png"
	outputRelFilepath := path.Join(getResultDir(), outputFilename)
	outputFilepath := path.Join(getBaseDir(), outputRelFilepath)
	iconsDir := path.Join(getBaseDir(), getIconsDir())

	err := gomosaics.Mosaicate(inputFilepath, iconsDir, outputFilepath, 4, 16)
	if err != nil {
		return nil, err
	}
	return &pb.ProcessFileRepl{
		Filepath: outputRelFilepath,
	}, nil
}

func main() {
	log.Println("Running mosaics-processing-service server")
	lis, err := net.Listen("tcp", "0.0.0.0:8081")
	if err != nil {
		log.Fatal("Could not bind network interface: ", err)
	}
	prepareDirs()
	srv := grpc.NewServer()
	pb.RegisterFileProcessingServer(srv, FileProcessingServerImpl{})
	srv.Serve(lis)
}
