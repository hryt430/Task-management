package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourorg/yourapp/internal/common/infrastructure/database"
	"github.com/yourorg/yourapp/internal/common/infrastructure/messaging"
)

// Application はアプリケーション全体を表します
type Application struct {
	server    *Server
	config    *Config
	db        *database.MySQLConnection
	msgClient *messaging.Client
}

// NewApplication は新しいアプリケーションインスタンスを作成します
func NewApplication() *Application {
	// 設定の読み込み
	config := NewConfig()

	return &Application{
		config: config,
	}
}

// Initialize はアプリケーションの初期化を行います
func (a *Application) Initialize() error {
	// データベース接続の初期化
	dbConfig := database.Config{
		Host:     a.config.DBHost,
		Port:     a.config.DBPort,
		User:     a.config.DBUser,
		Password: a.config.DBPassword,
		DBName:   a.config.DBName,
	}

	db, err := database.NewMySQLConnection(dbConfig)
	if err != nil {
		return fmt.Errorf("データベース接続エラー: %w", err)
	}
	a.db = &database.MySQLConnection{DB: db}

	// メッセージングクライアントの初期化
	a.msgClient = messaging.NewClient()

	// サーバーの初期化
	a.server = NewServer()

	return nil
}

// Run はアプリケーションを実行します
func (a *Application) Run() error {
	// リソースのクリーンアップ
	defer func() {
		if a.db != nil {
			a.db.Close()
		}
		if a.msgClient != nil {
			a.msgClient.Close()
		}
	}()

	// サーバーの起動
	go func() {
		log.Printf("サーバーを %s で起動しています...", a.config.GetServerAddr())
		if err := a.server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("サーバー起動エラー: %v", err)
		}
	}()

	// シグナル待ち受け
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("サーバーのシャットダウンを開始します...")

	// タイムアウト付きでサーバーを停止
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("サーバーのシャットダウンに失敗しました: %w", err)
	}

	log.Println("サーバーが正常に停止しました")
	return nil
}
