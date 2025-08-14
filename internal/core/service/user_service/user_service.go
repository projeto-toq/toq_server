package userservices

import (
	"context"
	"database/sql"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	cnpjport "github.com/giulio-alfieri/toq_server/internal/core/port/right/cnpj"
	cpfport "github.com/giulio-alfieri/toq_server/internal/core/port/right/cpf"
	creciport "github.com/giulio-alfieri/toq_server/internal/core/port/right/creci"
	gcsport "github.com/giulio-alfieri/toq_server/internal/core/port/right/gcs"
	sessionrepoport "github.com/giulio-alfieri/toq_server/internal/core/port/right/repository/session_repository"
	userrepoport "github.com/giulio-alfieri/toq_server/internal/core/port/right/repository/user_repository"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
	listingservices "github.com/giulio-alfieri/toq_server/internal/core/service/listing_service"
)

type userService struct {
	repo               userrepoport.UserRepoPortInterface
	sessionRepo        sessionrepoport.SessionRepoPortInterface
	globalService      globalservice.GlobalServiceInterface
	listingService     listingservices.ListingServiceInterface
	cpf                cpfport.CPFPortInterface
	cnpj               cnpjport.CNPJPortInterface
	creci              creciport.CreciPortInterface
	googleCloudService gcsport.GCSPortInterface
}

func NewUserService(
	ur userrepoport.UserRepoPortInterface,
	sr sessionrepoport.SessionRepoPortInterface,
	gsi globalservice.GlobalServiceInterface,
	listingService listingservices.ListingServiceInterface,
	cpf cpfport.CPFPortInterface,
	cnpj cnpjport.CNPJPortInterface,
	creci creciport.CreciPortInterface,
	gcs gcsport.GCSPortInterface,

) UserServiceInterface {
	return &userService{
		repo:               ur,
		sessionRepo:        sr,
		globalService:      gsi,
		listingService:     listingService,
		cpf:                cpf,
		cnpj:               cnpj,
		creci:              creci,
		googleCloudService: gcs,
	}
}

type UserServiceInterface interface {
	AcceptInvitation(ctx context.Context, userID int64) (err error)
	AddAlternativeRole(ctx context.Context, userID int64, role usermodel.UserRole, creciInfo ...string) (err error)
	ConfirmEmailChange(ctx context.Context, userID int64, code string) (tokens usermodel.Tokens, err error)
	ConfirmPasswordChange(ctx context.Context, nationalID string, password string, code string) (err error)
	ConfirmPhoneChange(ctx context.Context, userID int64, code string) (tokens usermodel.Tokens, err error)
	CreateAgency(ctx context.Context, agency usermodel.UserInterface) (tokens usermodel.Tokens, err error)
	CreateBaseRole(ctx context.Context, role usermodel.UserRole, name string) (err error)
	CreateOwner(ctx context.Context, owner usermodel.UserInterface) (tokens usermodel.Tokens, err error)
	CreateRealtor(ctx context.Context, realtor usermodel.UserInterface) (tokens usermodel.Tokens, err error)
	CreateRoot(ctx context.Context, root usermodel.UserInterface) (err error)
	CreateTokens(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface, expired bool) (tokens usermodel.Tokens, err error)
	DeleteAccount(ctx context.Context, userID int64) (tokens usermodel.Tokens, err error)
	DeleteAgencyOfRealtor(ctx context.Context, realtorID int64) (err error)
	DeleteRealtorOfAgency(ctx context.Context, agencyID int64, realtorID int64) (err error)
	GetAgencyOfRealtor(ctx context.Context, realtorID int64) (agency usermodel.UserInterface, err error)
	GetBaseRoles(ctx context.Context) (roles []usermodel.BaseRoleInterface, err error)
	GetCrecisToValidateByStatus(ctx context.Context, UserRoleStatus usermodel.UserRoleStatus) (realtors []usermodel.UserInterface, err error)
	GetOnboardingStatus(ctx context.Context, userID int64) (UserRoleStatus string, reason string, err error)
	GetProfile(ctx context.Context, userID int64) (user usermodel.UserInterface, err error)
	GetRealtorsByAgency(ctx context.Context, agencyID int64) (realtors []usermodel.UserInterface, err error)
	GetUserRolesByUser(ctx context.Context, userID int64) (roles []usermodel.UserRoleInterface, err error)
	GetUsers(ctx context.Context) (users []usermodel.UserInterface, err error)
	Home(ctx context.Context, userID int64) (user usermodel.UserInterface, err error)
	InviteRealtor(ctx context.Context, phoneNumber string) (err error)
	RefreshTokens(ctx context.Context, refresh string) (tokens usermodel.Tokens, err error)
	RejectInvitation(ctx context.Context, realtorID int64) (err error)
	RequestEmailChange(ctx context.Context, userID int64, newEmail string) (err error)
	RequestPasswordChange(ctx context.Context, nationalID string) (err error)
	RequestPhoneChange(ctx context.Context, userID int64, newPhone string) (err error)
	SignIn(ctx context.Context, nationalID string, password string, deviceToken string) (tokens usermodel.Tokens, err error)
	SignOut(ctx context.Context, userID int64, deviceToken, refreshToken string) (err error)
	SwitchUserRole(ctx context.Context, userID int64, userRoleID int64) (tokens usermodel.Tokens, err error)
	BatchUpdateLastActivity(ctx context.Context, userIDs []int64, timestamps []int64) (err error)
	UpdateProfile(ctx context.Context, user usermodel.UserInterface) (err error)
	ValidateCreciData(ctx context.Context, realtors []usermodel.UserInterface)
	ValidateCreciFace(ctx context.Context, realtors []usermodel.UserInterface)
	VerifyCreciImages(ctx context.Context, realtorID int64) (err error)
	UpdateOptStatus(ctx context.Context, optIn bool) (err error)
	GenerateGCSUploadURL(ctx context.Context, objectName, contentType string) (signedURL string, err error)
	GeneratePhotoDownloadURL(ctx context.Context, userID int64, photoType string) (signedURL string, err error)
	GetProfileThumbnails(ctx context.Context, userID int64) (thumbnails usermodel.ProfileThumbnails, err error)
	CreateUserFolder(ctx context.Context, userID int64) (err error)
	DeleteUserFolder(ctx context.Context, userID int64) (err error)
}
