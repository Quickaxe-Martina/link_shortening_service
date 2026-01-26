package grpcserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/auth"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/config"
	pb "github.com/Quickaxe-Martina/link_shortening_service/internal/grpc"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/service"
	"github.com/Quickaxe-Martina/link_shortening_service/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServer struct {
	pb.UnimplementedShortenerServiceServer
	service *service.ShortenerService
	storage storage.Storage
	cfg     *config.Config
}

func NewGRPCServer(svc *service.ShortenerService, cfg *config.Config, storage storage.Storage) *GRPCServer {
	return &GRPCServer{service: svc, cfg: cfg, storage: storage}
}

func (s *GRPCServer) ShortenURL(ctx context.Context, req *pb.URLShortenRequest) (*pb.URLShortenResponse, error) {
	tokenExp := time.Hour * time.Duration(s.cfg.TokenExp)
	user, err := auth.GetOrCreateUserFromGRPCContext(ctx, s.storage, s.cfg.SecretKey, tokenExp)
	if err != nil {
		return nil, err
	}

	shortURL, err := s.service.Shorten(ctx, user.ID, req.Url)
	if err != nil && !errors.Is(err, storage.ErrURLAlreadyExists) {
		return nil, err
	}

	return &pb.URLShortenResponse{Result: shortURL}, nil
}

func (s *GRPCServer) ExpandURL(ctx context.Context, req *pb.URLExpandRequest) (*pb.URLExpandResponse, error) {
	longURL, err := s.service.RedirectURL(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.URLExpandResponse{Result: longURL}, nil
}

func (s *GRPCServer) ListUserURLs(ctx context.Context, _ *emptypb.Empty) (*pb.UserURLsResponse, error) {
	tokenExp := time.Hour * time.Duration(s.cfg.TokenExp)
	user, err := auth.GetOrCreateUserFromGRPCContext(ctx, s.storage, s.cfg.SecretKey, tokenExp)
	if err != nil {
		return nil, err
	}

	urls, err := s.service.GetUserURLs(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	resp := &pb.UserURLsResponse{}
	for _, u := range urls {
		resp.Url = append(resp.Url, &pb.URLData{
			ShortUrl:    s.cfg.ServerAddr + u.Code,
			OriginalUrl: u.URL,
		})
	}

	return resp, nil
}

func RunGRPCServer(ctx context.Context, svc *service.ShortenerService, cfg *config.Config, storage storage.Storage) error {
	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterShortenerServiceServer(grpcServer, NewGRPCServer(svc, cfg, storage))
	reflection.Register(grpcServer)

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	fmt.Println("gRPC server listening on", cfg.GRPCAddr)
	return grpcServer.Serve(lis)
}
