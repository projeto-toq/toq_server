package userservices

import (
	"context"
	"database/sql"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	cnpjport "github.com/giulio-alfieri/toq_server/internal/core/port/right/cnpj"
	cpfport "github.com/giulio-alfieri/toq_server/internal/core/port/right/cpf"

	// creciport "github.com/giulio-alfieri/toq_server/internal/core/port/right/creci"
	sessionrepoport "github.com/giulio-alfieri/toq_server/internal/core/port/right/repository/session_repository"
	userrepoport "github.com/giulio-alfieri/toq_server/internal/core/port/right/repository/user_repository"
	storageport "github.com/giulio-alfieri/toq_server/internal/core/port/right/storage"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
	listingservices "github.com/giulio-alfieri/toq_server/internal/core/service/listing_service"
	permissionservices "github.com/giulio-alfieri/toq_server/internal/core/service/permission_service"
)

type userService struct {
	repo           userrepoport.UserRepoPortInterface
	sessionRepo    sessionrepoport.SessionRepoPortInterface
	globalService  globalservice.GlobalServiceInterface
	listingService listingservices.ListingServiceInterface
	cpf            cpfport.CPFPortInterface
	cnpj           cnpjport.CNPJPortInterface
	// creci               creciport.CreciPortInterface
	cloudStorageService storageport.CloudStoragePortInterface
	permissionService   permissionservices.PermissionServiceInterface // NOVO
	securityLogger      SecurityEventLoggerInterface                  // NOVO - Logger de eventos de segurança
}

func NewUserService(
	ur userrepoport.UserRepoPortInterface,
	sr sessionrepoport.SessionRepoPortInterface,
	gsi globalservice.GlobalServiceInterface,
	listingService listingservices.ListingServiceInterface,
	cpf cpfport.CPFPortInterface,
	cnpj cnpjport.CNPJPortInterface,
	// creci creciport.CreciPortInterface, // Pode ser nil temporariamente
	cloudStorage storageport.CloudStoragePortInterface,
	permissionService permissionservices.PermissionServiceInterface, // NOVO

) UserServiceInterface {
	return &userService{
		repo:           ur,
		sessionRepo:    sr,
		globalService:  gsi,
		listingService: listingService,
		cpf:            cpf,
		cnpj:           cnpj,
		// creci:               creci, // Pode ser nil
		cloudStorageService: cloudStorage,
		permissionService:   permissionService,        // NOVO
		securityLogger:      NewSecurityEventLogger(), // NOVO - Inicializa o logger
	}
}

type UserServiceInterface interface {
	AcceptInvitation(ctx context.Context, userID int64) (err error)
	AddAlternativeRole(ctx context.Context, userID int64, roleSlug permissionmodel.RoleSlug, creciInfo ...string) (err error)
	ConfirmEmailChange(ctx context.Context, userID int64, code string) (tokens usermodel.Tokens, err error)
	ConfirmPasswordChange(ctx context.Context, nationalID string, password string, code string) (err error)
	ConfirmPhoneChange(ctx context.Context, userID int64, code string) (tokens usermodel.Tokens, err error)
	// Fluxos de criação de conta que retornam tokens via SignIn padrão
	CreateAgency(ctx context.Context, agency usermodel.UserInterface, plainPassword string, deviceToken string, ipAddress string, userAgent string) (tokens usermodel.Tokens, err error)
	CreateOwner(ctx context.Context, owner usermodel.UserInterface, plainPassword string, deviceToken string, ipAddress string, userAgent string) (tokens usermodel.Tokens, err error)
	CreateRealtor(ctx context.Context, realtor usermodel.UserInterface, plainPassword string, deviceToken string, ipAddress string, userAgent string) (tokens usermodel.Tokens, err error)
	CreateTokens(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface, expired bool) (tokens usermodel.Tokens, err error)
	DeleteAccount(ctx context.Context, userID int64) (tokens usermodel.Tokens, err error)
	DeleteAgencyOfRealtor(ctx context.Context, realtorID int64) (err error)
	DeleteRealtorOfAgency(ctx context.Context, agencyID int64, realtorID int64) (err error)
	GetAgencyOfRealtor(ctx context.Context, realtorID int64) (agency usermodel.UserInterface, err error)
	GetOnboardingStatus(ctx context.Context, userID int64) (UserRoleStatus string, reason string, err error)
	GetProfile(ctx context.Context, userID int64) (user usermodel.UserInterface, err error)
	GetRealtorsByAgency(ctx context.Context, agencyID int64) (realtors []usermodel.UserInterface, err error)
	GetUsers(ctx context.Context) (users []usermodel.UserInterface, err error)
	Home(ctx context.Context, userID int64) (user usermodel.UserInterface, err error)
	InviteRealtor(ctx context.Context, phoneNumber string) (err error)
	RefreshTokens(ctx context.Context, refresh string) (tokens usermodel.Tokens, err error)
	RejectInvitation(ctx context.Context, realtorID int64) (err error)
	RequestEmailChange(ctx context.Context, userID int64, newEmail string) (err error)
	RequestPasswordChange(ctx context.Context, nationalID string) (err error)
	RequestPhoneChange(ctx context.Context, userID int64, newPhone string) (err error)
	ResendEmailChangeCode(ctx context.Context, userID int64) (err error)
	ResendPhoneChangeCode(ctx context.Context, userID int64) (err error)
	SignIn(ctx context.Context, nationalID string, password string, deviceToken string) (tokens usermodel.Tokens, err error)
	SignInWithContext(ctx context.Context, nationalID string, password string, deviceToken string, ipAddress string, userAgent string) (tokens usermodel.Tokens, err error)
	SignOut(ctx context.Context, userID int64, deviceToken, refreshToken string) (err error)
	SwitchUserRole(ctx context.Context, userID int64, roleSlug permissionmodel.RoleSlug) (tokens usermodel.Tokens, err error)
	BatchUpdateLastActivity(ctx context.Context, userIDs []int64, timestamps []int64) (err error)
	BlockUserTemporarily(ctx context.Context, userID int64) (err error)
	UnblockUserTemporarily(ctx context.Context, userID int64) (err error)
	UpdateProfile(ctx context.Context, user usermodel.UserInterface) (err error)
	UpdateOptStatus(ctx context.Context, optIn bool) (err error)
	GetPhotoUploadURL(ctx context.Context, objectName, contentType string) (signedURL string, err error)
	GeneratePhotoDownloadURL(ctx context.Context, userID int64, photoType string) (signedURL string, err error)
	GetProfileThumbnails(ctx context.Context, userID int64) (thumbnails usermodel.ProfileThumbnails, err error)
	CreateUserFolder(ctx context.Context, userID int64) (err error)
	DeleteUserFolder(ctx context.Context, userID int64) (err error)
}
