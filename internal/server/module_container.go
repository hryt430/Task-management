package server

import (
	"database/sql"

	"github.com/yourorg/yourapp/internal/common/infrastructure/messaging"
	authHttp "github.com/yourorg/yourapp/internal/modules/auth/delivery/http"
	authRepo "github.com/yourorg/yourapp/internal/modules/auth/repository/mysql"
	authUC "github.com/yourorg/yourapp/internal/modules/auth/usecase"
	notifHttp "github.com/yourorg/yourapp/internal/modules/notification/delivery/http"
	notifWs "github.com/yourorg/yourapp/internal/modules/notification/delivery/websocket"
	notifRepo "github.com/yourorg/yourapp/internal/modules/notification/repository/mysql"
	notifUC "github.com/yourorg/yourapp/internal/modules/notification/usecase"
	taskHttp "github.com/yourorg/yourapp/internal/modules/task/delivery/http"
	taskRepo "github.com/yourorg/yourapp/internal/modules/task/repository/mysql"
	taskUC "github.com/yourorg/yourapp/internal/modules/task/usecase"
)

// ModuleContainer はアプリケーションの全モジュールの依存関係を管理します
type ModuleContainer struct {
	// リポジトリ
	AuthRepository  authRepo.AuthRepository
	TaskRepository  taskRepo.TaskRepository
	NotifRepository notifRepo.NotificationRepository

	// ユースケース
	AuthUsecase  authUC.AuthUsecase
	TaskUsecase  taskUC.TaskUsecase
	NotifUsecase notifUC.NotificationUsecase

	// ハンドラ
	AuthHandler      *authHttp.AuthHandler
	TaskHandler      *taskHttp.TaskHandler
	NotifHandler     *notifHttp.NotificationHandler
	WebSocketHandler *notifWs.WebSocketHandler
}

// NewModuleContainer は新しいモジュールコンテナを作成します
func NewModuleContainer(db *sql.DB, msgClient *messaging.Client, config *Config) *ModuleContainer {
	container := &ModuleContainer{}

	// リポジトリの初期化
	container.AuthRepository = authRepo.NewMySQLAuthRepository(db)
	container.TaskRepository = taskRepo.NewMySQLTaskRepository(db)
	container.NotifRepository = notifRepo.NewMySQLNotificationRepository(db)

	// ユースケースの初期化
	container.AuthUsecase = authUC.NewAuthUsecase(
		container.AuthRepository,
		config.JWTSecret,
		config.JWTExpiryHours,
		config.JWTRefreshHours,
	)
	container.TaskUsecase = taskUC.NewTaskUsecase(
		container.TaskRepository,
	)
	container.NotifUsecase = notifUC.NewNotificationUsecase(
		container.NotifRepository,
		msgClient,
	)

	// ハンドラの初期化
	container.AuthHandler = authHttp.NewAuthHandler(container.AuthUsecase)
	container.TaskHandler = taskHttp.NewTaskHandler(container.TaskUsecase)
	container.NotifHandler = notifHttp.NewNotificationHandler(container.NotifUsecase)
	container.WebSocketHandler = notifWs.NewWebSocketHandler(container.NotifUsecase)

	return container
}
