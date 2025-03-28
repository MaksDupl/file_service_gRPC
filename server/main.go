package main

import (
	"context"
	"log"
	"net"

	"file_service/proto"
	"file_service/storage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fileServer struct {
	proto.UnimplementedFileServiceServer
	storage      *storage.FileStorage
	uploadSem    chan struct{} 
	downloadSem  chan struct{} 
	listFilesSem chan struct{} 
}

func NewFileServer() *fileServer {
	return &fileServer{
		storage:      storage.NewFileStorage("uploads"),
		uploadSem:    make(chan struct{}, 10), 
		downloadSem:  make(chan struct{}, 10), 
		listFilesSem: make(chan struct{}, 100), 
	}
}

func (s *fileServer) UploadFile(ctx context.Context, req *proto.UploadFileRequest) (*proto.UploadFileResponse, error) {
	// Ограничение concurrent запросов
	select {
	case s.uploadSem <- struct{}{}:
		defer func() { <-s.uploadSem }()
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "context canceled")
	}

	// Сохраняем файл
	fileInfo, err := s.storage.Save(req.Name, req.Content)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save file: %v", err)
	}

	return &proto.UploadFileResponse{
		Id:   fileInfo.Name,
		Name: fileInfo.Name,
		Size: fileInfo.Size,
	}, nil
}

func (s *fileServer) ListFiles(ctx context.Context, req *proto.ListFilesRequest) (*proto.ListFilesResponse, error) {
	// Ограничение concurrent запросов
	select {
	case s.listFilesSem <- struct{}{}:
		defer func() { <-s.listFilesSem }()
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "context canceled")
	}

	files, err := s.storage.List()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list files: %v", err)
	}

	var protoFiles []*proto.FileInfo
	for _, f := range files {
		protoFiles = append(protoFiles, &proto.FileInfo{
			Name:      f.Name,
			CreatedAt: f.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: f.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &proto.ListFilesResponse{Files: protoFiles}, nil
}

func (s *fileServer) DownloadFile(ctx context.Context, req *proto.DownloadFileRequest) (*proto.DownloadFileResponse, error) {
	// Ограничение concurrent запросов
	select {
	case s.downloadSem <- struct{}{}:
		defer func() { <-s.downloadSem }()
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "context canceled")
	}

	content, fileInfo, err := s.storage.Get(req.Name)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "file not found: %v", err)
	}

	return &proto.DownloadFileResponse{
		Name:    fileInfo.Name,
		Content: content,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterFileServiceServer(s, NewFileServer())

	log.Println("Server started on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
