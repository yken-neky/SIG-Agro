package handler

import (
	"context"
	"log"

	"github.com/sig-agro/services/producer-service/internal/repository"

	pb "github.com/sig-agro/api/proto/producer"
)

type ProducerService struct {
	pb.UnimplementedProducerServiceServer
	repo *repository.Repository
}

func NewProducerService(repo *repository.Repository) *ProducerService {
	return &ProducerService{repo: repo}
}

func (s *ProducerService) CreateProducer(ctx context.Context, req *pb.CreateProducerRequest) (*pb.CreateProducerResponse, error) {
	log.Printf("CreateProducer request for user ID: %d\n", req.UserId)

	producerID, err := s.repo.CreateProducer(ctx, req.UserId, req.Name, req.DocumentId, req.Phone, req.Email, req.Address)
	if err != nil {
		log.Printf("Error creating producer: %v\n", err)
		return &pb.CreateProducerResponse{Message: "Error creating producer"}, err
	}

	return &pb.CreateProducerResponse{
		ProducerId: producerID,
		Name:       req.Name,
		Message:    "Producer created successfully",
	}, nil
}

func (s *ProducerService) GetProducer(ctx context.Context, req *pb.GetProducerRequest) (*pb.GetProducerResponse, error) {
	log.Printf("GetProducer request for producer ID: %d\n", req.ProducerId)

	userID, name, docID, phone, email, address, err := s.repo.GetProducer(ctx, req.ProducerId)
	if err != nil {
		log.Printf("Error getting producer: %v\n", err)
		return nil, err
	}

	return &pb.GetProducerResponse{
		ProducerId: req.ProducerId,
		UserId:     userID,
		Name:       name,
		DocumentId: docID,
		Phone:      phone,
		Email:      email,
		Address:    address,
	}, nil
}

func (s *ProducerService) ListProducers(ctx context.Context, req *pb.ListProducersRequest) (*pb.ListProducersResponse, error) {
	log.Printf("ListProducers request for user ID: %d\n", req.UserId)

	producers, err := s.repo.ListProducers(ctx, req.UserId, req.Limit, req.Offset)
	if err != nil {
		log.Printf("Error listing producers: %v\n", err)
		return nil, err
	}

	resp := &pb.ListProducersResponse{Total: int32(len(producers))}
	for _, producer := range producers {
		pbProducer := &pb.GetProducerResponse{
			ProducerId: producer["id"].(int64),
			UserId:     producer["user_id"].(int64),
			Name:       producer["name"].(string),
			DocumentId: producer["document_id"].(string),
			Phone:      producer["phone"].(string),
			Email:      producer["email"].(string),
			Address:    producer["address"].(string),
			CreatedAt:  producer["created_at"].(int64),
		}
		resp.Producers = append(resp.Producers, pbProducer)
	}

	return resp, nil
}

func (s *ProducerService) UpdateProducer(ctx context.Context, req *pb.UpdateProducerRequest) (*pb.UpdateProducerResponse, error) {
	log.Printf("UpdateProducer request for producer ID: %d\n", req.ProducerId)

	err := s.repo.UpdateProducer(ctx, req.ProducerId, req.Name, req.Phone, req.Email, req.Address)
	if err != nil {
		log.Printf("Error updating producer: %v\n", err)
		return &pb.UpdateProducerResponse{Message: "Error updating producer"}, err
	}

	return &pb.UpdateProducerResponse{
		ProducerId: req.ProducerId,
		Message:    "Producer updated successfully",
	}, nil
}

func (s *ProducerService) DeleteProducer(ctx context.Context, req *pb.DeleteProducerRequest) (*pb.DeleteProducerResponse, error) {
	log.Printf("DeleteProducer request for producer ID: %d\n", req.ProducerId)

	err := s.repo.DeleteProducer(ctx, req.ProducerId)
	if err != nil {
		log.Printf("Error deleting producer: %v\n", err)
		return &pb.DeleteProducerResponse{Message: "Error deleting producer"}, err
	}

	return &pb.DeleteProducerResponse{
		Success: true,
		Message: "Producer deleted successfully",
	}, nil
}
