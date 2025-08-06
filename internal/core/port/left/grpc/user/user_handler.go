package grpcuserport

import (
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	userservices "github.com/giulio-alfieri/toq_server/internal/core/service/user_service"
)

type UserHandler struct {
	service userservices.UserServiceInterface
	pb.UnimplementedUserServiceServer
}

func NewUserHandler(service userservices.UserServiceInterface) *UserHandler { //userHandlerInterface {
	return &UserHandler{
		service: service,
	}

}

// type userHandlerInterface interface {
// 	pb.UserServiceServer
// }
