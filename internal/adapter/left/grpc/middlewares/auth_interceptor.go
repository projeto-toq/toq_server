package middlewares

import (
	"context"
	"log/slog"
	"net"
	"strings"

	goroutines "github.com/giulio-alfieri/toq_server/internal/core/go_routines"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var nonAuthMethods = []string{
	"/grpc.UserService/CreateOwner",
	"/grpc.UserService/CreateRealtor",
	"/grpc.UserService/CreateAgency",
	"/grpc.UserService/SignIn",
	"/grpc.UserService/RefreshToken",
	"/grpc.UserService/RequestPasswordChange",
	"/grpc.UserService/ConfirmPasswordChange",
}

func AuthInterceptor(ctx context.Context, activityTracker *goroutines.ActivityTracker) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		for _, method := range nonAuthMethods {
			if info.FullMethod == method {
				// add toqroot as request context
				infos := usermodel.UserInfos{
					ID:            0,
					Role:          usermodel.UserRole(0),
					ProfileStatus: false,
				}
				ctx = context.WithValue(ctx, globalmodel.TokenKey, infos)
				return handler(ctx, req)
			}
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			slog.Warn("metadata is not provided on the context for the request, during AuthInterceptor")
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		tokens := md.Get("authorization")
		if len(tokens) == 0 {
			slog.Warn("authorization token is not provided, during AuthInterceptor")
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}

		// Verificar se o token está no formato "Bearer TOKEN"
		tokenParts := strings.Split(tokens[0], "Bearer ")
		if len(tokenParts) < 2 || tokenParts[1] == "" {
			slog.Warn("invalid authorization token format, expected 'Bearer TOKEN'", "token", tokens[0])
			return nil, status.Errorf(codes.Unauthenticated, "invalid authorization token format")
		}

		token := tokenParts[1]
		infos, err := validateAccessToken(token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid access token")
		}

		// Extract metadata for session (user-agent, ip)
		userAgent := ""
		if uaVals := md.Get("user-agent"); len(uaVals) > 0 {
			userAgent = uaVals[0]
		}
		// Try to derive client IP from :authority or forwarded headers
		clientIP := ""
		if ipVals := md.Get("x-forwarded-for"); len(ipVals) > 0 {
			clientIP = strings.Split(ipVals[0], ",")[0]
		}
		if clientIP == "" {
			if peerVals := md.Get(":authority"); len(peerVals) > 0 {
				host := peerVals[0]
				hostOnly, _, errHost := net.SplitHostPort(host)
				if errHost == nil {
					clientIP = hostOnly
				} else {
					clientIP = host
				}
			}
		}

		// Add user infos + meta to context (future: define structured key types)
		ctx = context.WithValue(ctx, globalmodel.TokenKey, infos)
		ctx = context.WithValue(ctx, globalmodel.UserAgentKey, userAgent)
		ctx = context.WithValue(ctx, globalmodel.ClientIPKey, clientIP)

		// Track user activity (non-blocking, fast Redis operation)
		activityTracker.TrackActivity(ctx, infos.ID)

		return handler(ctx, req)
	}
}

// Validates the given access token and updates de UserID on userDomain if valid. Otherwise err
func validateAccessToken(accessToken string) (infos usermodel.UserInfos, err error) {

	infos = usermodel.UserInfos{}

	//tenta validar o token
	token, err2 := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			slog.Warn("unexpected signing method", "method", token.Header["alg"])
			return nil, status.Error(codes.Unauthenticated, "invalid access token")
		}

		secret := globalmodel.GetJWTSecret()
		return []byte(secret), nil
	})

	if err2 != nil {
		return infos, status.Error(codes.Unauthenticated, "invalid access token")

	}

	//tenta recuperar os claims
	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		slog.Warn("cannot get claims from token")
		return infos, status.Error(codes.Unauthenticated, "invalid access token")
	}

	infosraw, ok := payload[string(globalmodel.TokenKey)].(map[string]interface{})
	if !ok {
		slog.Warn("cannot get user infos from token")
		return infos, status.Error(codes.Unauthenticated, "invalid access token")
	}
	float64ID, ok := infosraw["ID"].(float64)
	if !ok {
		slog.Warn("cannot get user ID from token")
		return infos, status.Error(codes.Unauthenticated, "invalid access token")
	}

	int64ID := int64(float64ID)
	if !ok {
		slog.Warn("cannot convert user ID to uint32")
		return infos, status.Error(codes.Unauthenticated, "invalid access token")
	}
	infos.ID = int64ID

	profileValid, ok := infosraw["ProfileStatus"].(bool)
	if !ok {
		slog.Warn("cannot get user profile status from token")
		return infos, status.Error(codes.Unauthenticated, "invalid access token")
	}
	infos.ProfileStatus = profileValid

	roleFloat64, ok := infosraw["Role"].(float64)
	if !ok {
		slog.Warn("cannot get user role from token")
		return infos, status.Error(codes.Unauthenticated, "invalid access token")
	}
	infos.Role = usermodel.UserRole(roleFloat64)

	return
}
