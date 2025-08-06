package grpclistingport

import (
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	listingservices "github.com/giulio-alfieri/toq_server/internal/core/service/listing_service"
)

type ListingHandler struct {
	service listingservices.ListingServiceInterface
	pb.UnimplementedListingServiceServer
}

func NewUserHandler(service listingservices.ListingServiceInterface) *ListingHandler {
	return &ListingHandler{
		service: service,
	}

}
