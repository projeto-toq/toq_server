package goroutines

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	userservices "github.com/giulio-alfieri/toq_server/internal/core/service/user_service"
	"github.com/google/uuid"
)

func GoUpdateLastActivity(wg *sync.WaitGroup, service userservices.UserServiceInterface, activityChannel chan int64, ctx context.Context) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Goroutine canceled")
			return
		case userID := <-activityChannel:
			infos := usermodel.UserInfos{} //todo verifique se é necessário manter pois já está incluido no ctx de config
			infos.ID = userID
			ctx := context.WithValue(ctx, globalmodel.TokenKey, infos)
			ctx = context.WithValue(ctx, globalmodel.RequestIDKey, uuid.New().String())
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := service.UpdateLastActivity(ctx, userID)
				if err != nil {
					slog.Error("Error updating last activity", "error", err)
				}
			}()
		}
	}
}
