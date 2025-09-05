package sessionservice

import (
	sessionrepoport "github.com/giulio-alfieri/toq_server/internal/core/port/right/repository/session_repository"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
)

// service is the concrete implementation of Service.
type service struct {
	repo          sessionrepoport.SessionRepoPortInterface
	globalService globalservice.GlobalServiceInterface
}

// New constructs a Service with its dependencies.
func New(repo sessionrepoport.SessionRepoPortInterface, gs globalservice.GlobalServiceInterface) Service {
	return &service{repo: repo, globalService: gs}
}
