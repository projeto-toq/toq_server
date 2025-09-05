package validationservice

import (
	userrepo "github.com/giulio-alfieri/toq_server/internal/core/port/right/repository/user_repository"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
)

type service struct {
	repo          userrepo.UserRepoPortInterface
	globalService globalservice.GlobalServiceInterface
}

func New(repo userrepo.UserRepoPortInterface, gs globalservice.GlobalServiceInterface) Service {
	return &service{repo: repo, globalService: gs}
}
