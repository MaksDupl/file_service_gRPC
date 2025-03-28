package main

import (
	"context"
	"file_service/proto"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewFileServiceClient(conn)

	if len(os.Args) > 1 && os.Args[1] == "upload" {
		fileContent, err := os.ReadFile(os.Args[2])
		if err != nil {
			log.Fatalf("could not read file: %v", err)
		}

		resp, err := client.UploadFile(context.Background(), &proto.UploadFileRequest{
			Name:    filepath.Base(os.Args[2]),
			Content: fileContent,
		})
		if err != nil {
			log.Fatalf("could not upload file: %v", err)
		}
		fmt.Printf("File uploaded: %s (%d bytes)\n", resp.Name, resp.Size)
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "list" {
		resp, err := client.ListFiles(context.Background(), &proto.ListFilesRequest{})
		if err != nil {
			log.Fatalf("could not list files: %v", err)
		}

		fmt.Println("Files:")
		for _, f := range resp.Files {
			fmt.Printf("- %s (created: %s, updated: %s)\n", f.Name, f.CreatedAt, f.UpdatedAt)
		}
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "download" {
		resp, err := client.DownloadFile(context.Background(), &proto.DownloadFileRequest{
			Name: os.Args[2],
		})
		if err != nil {
			log.Fatalf("could not download file: %v", err)
		}

		err = os.WriteFile(os.Args[2], resp.Content, 0644)
		if err != nil {
			log.Fatalf("could not save file: %v", err)
		}
		fmt.Printf("File downloaded: %s (%d bytes)\n", resp.Name, len(resp.Content))
		return
	}

	fmt.Println("Usage:")
	fmt.Println("  upload <file> - Upload a file")
	fmt.Println("  list - List all files")
	fmt.Println("  download <name> - Download a file")
}
