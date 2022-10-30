package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"learn-protobuf/1/pb"
	"log"
	"time"
)

func main() {
	certFile := "/Users/koichi/Library/Application Support/mkcert/rootCA.pem"
	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		log.Fatalln(err)
	}
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))
	//conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewFileServiceClient(conn)
	//callListFiles(client)
	callDownload(client)
}

func callListFiles(client pb.FileServiceClient) {
	//md := metadata.New(map[string]string{"authorization": "Bearer test-token"})
	md := metadata.New(map[string]string{"authorization": "Bearer bad-token"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	res, err := client.ListFiles(ctx, &pb.ListFilesRequest{})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(res.GetFilenames())
}

func callDownload(client pb.FileServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.DownloadRequest{Filename: "name.txt"}
	stream, err := client.Download(ctx, req)
	if err != nil {
		log.Fatalln(err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			resErr, ok := status.FromError(err)
			if ok {
				if resErr.Code() == codes.NotFound {
					log.Fatalf("Error Code: %v, Error Message: %v", resErr.Code(), resErr.Message())
				} else if resErr.Code() == codes.DeadlineExceeded {
					log.Fatalln("deadline exceeded")
				} else {
					log.Fatalln("unknown grpc error")
				}
			} else {
				log.Fatalln(err)
			}
		}

		log.Printf("Response from Download(bytes): %v", res.GetData())
		log.Printf("Response from Download(string): %v", string(res.GetData()))
	}
}
