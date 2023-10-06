package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"

	pb "mosaicsgrpc/services"
)

type HelloWorldServerImpl struct {
	pb.UnimplementedHelloWorldServer
}

func (s HelloWorldServerImpl) SayHello(c context.Context, req *pb.HelloReq) (*pb.HelloRepl, error) {
	return &pb.HelloRepl{Msg: fmt.Sprintf("Hello there, %s!!!!", req.GetName())}, nil
}

func askForHello(client pb.HelloWorldClient) {
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	val, _ := client.SayHello(c, &pb.HelloReq{Name: "Ytsy"})
	log.Println("Received a response: ", val)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("No args provided")
	}

	mode := os.Args[1]
	if mode == "client" {
		log.Println("Running as client")
		conn, err := grpc.Dial("localhost:8081")
		if err != nil {
			log.Fatal("An error occured: ", err)
		}
		defer conn.Close()
		client := pb.NewHelloWorldClient(conn)
		askForHello(client)
	} else if mode == "server" {
		log.Println("Running as server")
		lis, _ := net.Listen("tcp", "localhost:8081")
		srv := grpc.NewServer()
		pb.RegisterHelloWorldServer(srv, HelloWorldServerImpl{})
		srv.Serve(lis)
	} else {
		log.Fatal("Must run either as client or as server")
	}
}
