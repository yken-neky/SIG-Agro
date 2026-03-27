package handler

import (
	"context"
	"log"

	"github.com/sig-agro/services/parcel-service/internal/repository"

	pb "github.com/sig-agro/api/proto/parcel"
)

type ParcelService struct {
	pb.UnimplementedParcelServiceServer
	repo *repository.Repository
}

func NewParcelService(repo *repository.Repository) *ParcelService {
	return &ParcelService{repo: repo}
}

func (s *ParcelService) CreateParcel(ctx context.Context, req *pb.CreateParcelRequest) (*pb.CreateParcelResponse, error) {
	log.Printf("CreateParcel request for producer ID: %d\n", req.ProducerId)

	parcelID, err := s.repo.CreateParcel(ctx, req.ProducerId, req.Name, req.Description, req.GeometryWkt, float64(req.AreaHectares), req.CropType)
	if err != nil {
		log.Printf("Error creating parcel: %v\n", err)
		return &pb.CreateParcelResponse{Message: "Error creating parcel"}, err
	}

	// Publish event: parcel created (TODO: RabbitMQ integration)
	log.Printf("Event: Parcel created (ID: %d)\n", parcelID)

	return &pb.CreateParcelResponse{
		ParcelId: parcelID,
		Name:     req.Name,
		Message:  "Parcel created successfully",
	}, nil
}

func (s *ParcelService) GetParcel(ctx context.Context, req *pb.GetParcelRequest) (*pb.GetParcelResponse, error) {
	log.Printf("GetParcel request for parcel ID: %d\n", req.ParcelId)

	producerID, name, description, geometryWKT, areaHectares, cropType, err := s.repo.GetParcel(ctx, req.ParcelId)
	if err != nil {
		log.Printf("Error getting parcel: %v\n", err)
		return nil, err
	}

	return &pb.GetParcelResponse{
		ParcelId:     req.ParcelId,
		ProducerId:   producerID,
		Name:         name,
		Description:  description,
		GeometryWkt:  geometryWKT,
		AreaHectares: areaHectares,
		CropType:     cropType,
	}, nil
}

func (s *ParcelService) ListParcels(ctx context.Context, req *pb.ListParcelsRequest) (*pb.ListParcelsResponse, error) {
	log.Printf("ListParcels request for producer ID: %d\n", req.ProducerId)

	parcels, err := s.repo.ListParcels(ctx, req.ProducerId, req.Limit, req.Offset)
	if err != nil {
		log.Printf("Error listing parcels: %v\n", err)
		return nil, err
	}

	resp := &pb.ListParcelsResponse{Total: int32(len(parcels))}
	for _, parcel := range parcels {
		pbParcel := &pb.GetParcelResponse{
			ParcelId:     parcel["id"].(int64),
			ProducerId:   parcel["producer_id"].(int64),
			Name:         parcel["name"].(string),
			Description:  parcel["description"].(string),
			GeometryWkt:  parcel["geometry_wkt"].(string),
			AreaHectares: parcel["area_hectares"].(float64),
			CropType:     parcel["crop_type"].(string),
			CreatedAt:    parcel["created_at"].(int64),
		}
		resp.Parcels = append(resp.Parcels, pbParcel)
	}

	return resp, nil
}

func (s *ParcelService) UpdateParcel(ctx context.Context, req *pb.UpdateParcelRequest) (*pb.UpdateParcelResponse, error) {
	log.Printf("UpdateParcel request for parcel ID: %d\n", req.ParcelId)

	err := s.repo.UpdateParcel(ctx, req.ParcelId, req.Name, req.Description, req.CropType, req.GeometryWkt)
	if err != nil {
		log.Printf("Error updating parcel: %v\n", err)
		return &pb.UpdateParcelResponse{Message: "Error updating parcel"}, err
	}

	return &pb.UpdateParcelResponse{
		ParcelId: req.ParcelId,
		Message:  "Parcel updated successfully",
	}, nil
}

func (s *ParcelService) DeleteParcel(ctx context.Context, req *pb.DeleteParcelRequest) (*pb.DeleteParcelResponse, error) {
	log.Printf("DeleteParcel request for parcel ID: %d\n", req.ParcelId)

	err := s.repo.DeleteParcel(ctx, req.ParcelId)
	if err != nil {
		log.Printf("Error deleting parcel: %v\n", err)
		return &pb.DeleteParcelResponse{Message: "Error deleting parcel"}, err
	}

	return &pb.DeleteParcelResponse{
		Success: true,
		Message: "Parcel deleted successfully",
	}, nil
}

func (s *ParcelService) QueryByGeometry(ctx context.Context, req *pb.QueryByGeometryRequest) (*pb.QueryByGeometryResponse, error) {
	log.Printf("QueryByGeometry request with type: %s\n", req.QueryType)

	parcels, err := s.repo.QueryByGeometry(ctx, req.GeometryWkt)
	if err != nil {
		log.Printf("Error querying by geometry: %v\n", err)
		return nil, err
	}

	resp := &pb.QueryByGeometryResponse{Total: int32(len(parcels))}
	for _, parcel := range parcels {
		pbParcel := &pb.GetParcelResponse{
			ParcelId:     parcel["id"].(int64),
			ProducerId:   parcel["producer_id"].(int64),
			Name:         parcel["name"].(string),
			Description:  parcel["description"].(string),
			GeometryWkt:  parcel["geometry_wkt"].(string),
			AreaHectares: parcel["area_hectares"].(float64),
			CropType:     parcel["crop_type"].(string),
			CreatedAt:    parcel["created_at"].(int64),
		}
		resp.Parcels = append(resp.Parcels, pbParcel)
	}

	return resp, nil
}
