package userservices

import (
	"context"
	"database/sql"
	"time"

	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	cnpjport "github.com/projeto-toq/toq_server/internal/core/port/right/cnpj"
	cpfport "github.com/projeto-toq/toq_server/internal/core/port/right/cpf"

	// creciport "github.com/projeto-toq/toq_server/internal/core/port/right/creci"
	devicetokenrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/device_token_repository"
	sessionrepoport "github.com/projeto-toq/toq_server/internal/core/port/right/repository/session_repository"
	userrepoport "github.com/projeto-toq/toq_server/internal/core/port/right/repository/user_repository"
	storageport "github.com/projeto-toq/toq_server/internal/core/port/right/storage"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	listingservices "github.com/projeto-toq/toq_server/internal/core/service/listing_service"
	permissionservices "github.com/projeto-toq/toq_server/internal/core/service/permission_service"
	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
)

type userService struct {
	repo                userrepoport.UserRepoPortInterface
	sessionRepo         sessionrepoport.SessionRepoPortInterface
	deviceTokenRepo     devicetokenrepository.DeviceTokenRepoPortInterface
	globalService       globalservice.GlobalServiceInterface
	listingService      listingservices.ListingServiceInterface
	photoSessionService photosessionservices.PhotoSessionServiceInterface
	cpf                 cpfport.CPFPortInterface
	cnpj                cnpjport.CNPJPortInterface
	// creci               creciport.CreciPortInterface
	cloudStorageService storageport.CloudStoragePortInterface
	permissionService   permissionservices.PermissionServiceInterface // NOVO
	cfg                 Config
	emailRenderer       *systemUserWelcomeEmailRenderer
}

func NewUserService(
	ur userrepoport.UserRepoPortInterface,
	sr sessionrepoport.SessionRepoPortInterface,
	deviceTokenRepo devicetokenrepository.DeviceTokenRepoPortInterface,
	gsi globalservice.GlobalServiceInterface,
	listingService listingservices.ListingServiceInterface,
	photoSessionService photosessionservices.PhotoSessionServiceInterface,
	cpf cpfport.CPFPortInterface,
	cnpj cnpjport.CNPJPortInterface,
	// creci creciport.CreciPortInterface, // Pode ser nil temporariamente
	cloudStorage storageport.CloudStoragePortInterface,
	permissionService permissionservices.PermissionServiceInterface, // NOVO
	cfg Config,

) UserServiceInterface {
	config := normalizeConfig(cfg)
	return &userService{
		repo:                ur,
		sessionRepo:         sr,
		deviceTokenRepo:     deviceTokenRepo,
		globalService:       gsi,
		listingService:      listingService,
		photoSessionService: photoSessionService,
		cpf:                 cpf,
		cnpj:                cnpj,
		// creci:               creci, // Pode ser nil
		cloudStorageService: cloudStorage,
		permissionService:   permissionService, // NOVO
		cfg:                 config,
		emailRenderer:       newSystemUserWelcomeEmailRenderer(),
	}
}

type UserServiceInterface interface {
	AcceptInvitation(ctx context.Context, userID int64) (err error)
	AddAlternativeRole(ctx context.Context, userID int64, roleSlug permissionmodel.RoleSlug, creciInfo ...string) (err error)
	ConfirmEmailChange(ctx context.Context, code string) (err error)
	ConfirmPasswordChange(ctx context.Context, nationalID string, password string, code string) (err error)
	ConfirmPhoneChange(ctx context.Context, code string) (err error)
	// Fluxos de criação de conta que retornam tokens via SignIn padrão
	CreateAgency(ctx context.Context, agency usermodel.UserInterface, plainPassword string, deviceToken string, ipAddress string, userAgent string) (tokens usermodel.Tokens, err error)
	CreateOwner(ctx context.Context, owner usermodel.UserInterface, plainPassword string, deviceToken string, ipAddress string, userAgent string) (tokens usermodel.Tokens, err error)
	CreateRealtor(ctx context.Context, realtor usermodel.UserInterface, plainPassword string, deviceToken string, ipAddress string, userAgent string) (tokens usermodel.Tokens, err error)
	CreateTokens(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface, expired bool) (tokens usermodel.Tokens, err error)
	DeleteAccount(ctx context.Context) (tokens usermodel.Tokens, err error)
	DeleteAgencyOfRealtor(ctx context.Context, realtorID int64) (err error)
	DeleteRealtorOfAgency(ctx context.Context, agencyID int64, realtorID int64) (err error)
	GetAgencyOfRealtor(ctx context.Context, realtorID int64) (agency usermodel.UserInterface, err error)
	GetProfile(ctx context.Context) (user usermodel.UserInterface, err error)
	GetRealtorsByAgency(ctx context.Context, agencyID int64) (realtors []usermodel.UserInterface, err error)
	GetUsers(ctx context.Context) (users []usermodel.UserInterface, err error)
	Home(ctx context.Context, userID int64) (user usermodel.UserInterface, err error)
	InviteRealtor(ctx context.Context, phoneNumber string) (err error)
	RefreshTokens(ctx context.Context, refresh string) (tokens usermodel.Tokens, err error)
	RejectInvitation(ctx context.Context, realtorID int64) (err error)
	RequestEmailChange(ctx context.Context, newEmail string) (err error)
	RequestPasswordChange(ctx context.Context, nationalID string) (err error)
	RequestPhoneChange(ctx context.Context, newPhone string) (err error)
	ResendEmailChangeCode(ctx context.Context) (err error)
	ResendPhoneChangeCode(ctx context.Context) (err error)
	SignIn(ctx context.Context, nationalID string, password string, deviceToken string, deviceID string) (tokens usermodel.Tokens, err error)
	SignInWithContext(ctx context.Context, nationalID string, password string, deviceToken string, deviceID string, ipAddress string, userAgent string) (tokens usermodel.Tokens, err error)
	SignOut(ctx context.Context, deviceToken, refreshToken, deviceID string) (err error)
	SwitchUserRole(ctx context.Context) (tokens usermodel.Tokens, err error)
	BatchUpdateLastActivity(ctx context.Context, userIDs []int64, timestamps []int64) (err error)
	BlockUserTemporarily(ctx context.Context, userID int64) (err error)
	UnblockUserTemporarily(ctx context.Context, userID int64) (err error)
	// UpdateProfile updates allowed user profile fields using a typed input contract.
	// It must not change email, phone or password; those have dedicated flows.
	UpdateProfile(ctx context.Context, in UpdateProfileInput) (err error)
	UpdateOptStatus(ctx context.Context, optIn bool) (err error)
	GetPhotoUploadURL(ctx context.Context, variant, contentType string) (signedURL string, err error)
	GetPhotoDownloadURL(ctx context.Context, variant string) (signedURL string, err error)
	CreateUserFolder(ctx context.Context, userID int64) (err error)
	DeleteUserFolder(ctx context.Context, userID int64) (err error)
	// GetCreciUploadURL generates a signed upload URL for realtor CRECI documents
	GetCreciUploadURL(ctx context.Context, documentType, contentType string) (signedURL string, err error)
	// GetCreciDownloadURLs generates signed download URLs for CRECI documents of a target user
	GetCreciDownloadURLs(ctx context.Context, userID int64) (CreciDocumentDownloadURLs, error)
	// VerifyCreciDocuments checks S3 for required CRECI images and sets status to PendingManual
	VerifyCreciDocuments(ctx context.Context) (err error)
	// GetUserByID returns the user with the active role eagerly loaded (read-only tx)
	GetUserByID(ctx context.Context, id int64) (usermodel.UserInterface, error)
	// GetUserByIDWithTx returns the user with the active role using the provided transaction
	GetUserByIDWithTx(ctx context.Context, tx *sql.Tx, id int64) (usermodel.UserInterface, error)
	// GetActiveRoleStatus returns only the status of the active user role
	GetActiveRoleStatus(ctx context.Context) (status permissionmodel.UserRoleStatus, err error)
	// GetCrecisToValidateByStatus returns realtors filtered by active role status
	GetCrecisToValidateByStatus(ctx context.Context, status permissionmodel.UserRoleStatus) ([]usermodel.UserInterface, error)

	// ApproveCreciManual updates realtor status from pending manual to approved/refused and sends notification
	ApproveCreciManual(ctx context.Context, userID int64, status permissionmodel.UserRoleStatus) error

	// Admin system user management
	ListUsers(ctx context.Context, input ListUsersInput) (ListUsersOutput, error)
	ListPendingRealtors(ctx context.Context, page, limit int) (ListPendingRealtorsOutput, error)
	CreateSystemUser(ctx context.Context, input CreateSystemUserInput) (SystemUserResult, error)
	UpdateSystemUser(ctx context.Context, input UpdateSystemUserInput) (SystemUserResult, error)
	DeleteSystemUser(ctx context.Context, input DeleteSystemUserInput) error

	ValidateCPF(ctx context.Context, nationalID string, bornAt time.Time) error
	ValidateCNPJ(ctx context.Context, nationalID string) error
}

// CreciDocumentDownloadURLs encapsula as URLs assinadas geradas pelo serviço para os documentos CRECI
type CreciDocumentDownloadURLs struct {
	Selfie           string
	Front            string
	Back             string
	ExpiresInMinutes int
}
