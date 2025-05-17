package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/yourorg/yourapp/internal/common/middleware"
	authHttp "github.com/yourorg/yourapp/internal/modules/auth/delivery/http"
	notifHttp "github.com/yourorg/yourapp/internal/modules/notification/delivery/http"
	notifWs "github.com/yourorg/yourapp/internal/modules/notification/delivery/websocket"
	taskHttp "github.com/yourorg/yourapp/internal/modules/task/delivery/http"
)

// Server はHTTPサーバーを表します
type Server struct {
	router *mux.Router
	server *http.Server
}

// NewServer は新しいサーバーインスタンスを作成します
func NewServer() *Server {
	router := mux.NewRouter()

	// 共通のミドルウェアを設定
	router.Use(middleware.Recovery)
	router.Use(middleware.Logger)
	router.Use(middleware.CORS)

	// ヘルスチェックなど基本的なエンドポイント
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	return &Server{
		router: router,
		server: &http.Server{
			Addr:         getEnv("SERVER_ADDR", ":8080"),
			Handler:      router,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

// SetupRoutes は各モジュールのルートを設定します
func (s *Server) SetupRoutes(
	authHandler *authHttp.AuthHandler,
	taskHandler *taskHttp.TaskHandler,
	notifHandler *notifHttp.NotificationHandler,
	wsHandler *notifWs.WebSocketHandler,
) {
	// 認証モジュールのルート
	authRouter := s.router.PathPrefix("/api/auth").Subrouter()
	authHandler.RegisterRoutes(authRouter)

	// タスク管理モジュールのルート
	taskRouter := s.router.PathPrefix("/api/tasks").Subrouter()
	taskHandler.RegisterRoutes(taskRouter)

	// 通知モジュールのルート
	notifRouter := s.router.PathPrefix("/api/notifications").Subrouter()
	notifHandler.RegisterRoutes(notifRouter)

	// WebSocketハンドラのルート
	s.router.HandleFunc("/ws/notifications", wsHandler.HandleWebSocket)

	// 静的ファイルの提供（フロントエンドアプリケーション用）
	fsHandler := http.FileServer(http.Dir("./static"))
	s.router.PathPrefix("/").Handler(http.StripPrefix("/", fsHandler))
}

// Start はサーバーを起動します
func (s *Server) Start() error {
	log.Printf("サーバーを %s で起動しています...", s.server.Addr)
	return s.server.ListenAndServe()
}

// Shutdown はサーバーを停止します
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
