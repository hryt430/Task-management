package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yourorg/yourapp/internal/common/middleware"
	"github.com/yourorg/yourapp/internal/modules/auth/usecase"
)

// Router はアプリケーションのルーティング設定を担当します
type Router struct {
	router      *mux.Router
	authUsecase usecase.AuthUsecase
}

// NewRouter は新しいルーターインスタンスを作成します
func NewRouter(authUsecase usecase.AuthUsecase) *Router {
	return &Router{
		router:      mux.NewRouter(),
		authUsecase: authUsecase,
	}
}

// Setup はルーターの初期設定を行います
func (r *Router) Setup() *mux.Router {
	// 共通のミドルウェアを設定
	r.router.Use(middleware.Recovery)
	r.router.Use(middleware.Logger)
	r.router.Use(middleware.CORS)

	// APIのバージョニング
	api := r.router.PathPrefix("/api/v1").Subrouter()

	// 認証が不要なエンドポイント
	public := api.PathPrefix("").Subrouter()
	r.setupPublicRoutes(public)

	// 認証が必要なエンドポイント
	authMiddleware := middleware.NewJWTAuthMiddleware(r.authUsecase)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(authMiddleware.Middleware)
	r.setupProtectedRoutes(protected)

	// 静的ファイルの提供
	r.setupStaticFiles()

	return r.router
}

// setupPublicRoutes は認証が不要なエンドポイントを設定します
func (r *Router) setupPublicRoutes(router *mux.Router) {
	// ヘルスチェック
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// 認証関連（ログイン、登録など）
	router.PathPrefix("/auth").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 実際の実装では認証ハンドラーを呼び出します
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("認証エンドポイント"))
	})
}

// setupProtectedRoutes は認証が必要なエンドポイントを設定します
func (r *Router) setupProtectedRoutes(router *mux.Router) {
	// 管理者専用ルート
	adminRouter := router.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middleware.NewJWTAuthMiddleware(r.authUsecase).RequireRole("admin"))

	// タスク関連
	taskRouter := router.PathPrefix("/tasks").Subrouter()
	taskRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		// 実際の実装ではタスクハンドラーを呼び出します
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("タスクエンドポイント"))
	})

	// 通知関連
	notificationRouter := router.PathPrefix("/notifications").Subrouter()
	notificationRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		// 実際の実装では通知ハンドラーを呼び出します
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("通知エンドポイント"))
	})

	// ユーザー関連
	userRouter := router.PathPrefix("/users").Subrouter()
	userRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		// 実際の実装ではユーザーハンドラーを呼び出します
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ユーザーエンドポイント"))
	})
}

// setupStaticFiles は静的ファイルの提供を設定します
func (r *Router) setupStaticFiles() {
	// SPA（Single Page Application）のフロントエンドサポート
	spa := http.FileServer(http.Dir("./web/dist"))
	r.router.PathPrefix("/").Handler(spa)
}

// GetRouter はルーターインスタンスを返します
func (r *Router) GetRouter() *mux.Router {
	return r.router
}
