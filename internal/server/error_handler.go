package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// ErrorResponse はエラーレスポンスの構造を定義します
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// Common error types
var (
	ErrBadRequest          = errors.New("不正なリクエストです")
	ErrUnauthorized        = errors.New("認証が必要です")
	ErrForbidden           = errors.New("アクセス権限がありません")
	ErrNotFound            = errors.New("リソースが見つかりません")
	ErrInternalServerError = errors.New("サーバー内部エラーが発生しました")
)

// HTTPErrorHandler はエラーハンドリングを行う関数です
func HTTPErrorHandler(w http.ResponseWriter, err error, statusCode int) {
	// ステータスコードの決定
	if statusCode == 0 {
		statusCode = determineStatusCode(err)
	}

	// エラーレスポンスの作成
	errorResponse := ErrorResponse{
		Status:  statusCode,
		Message: err.Error(),
	}

	// 内部エラーの場合はログを残す
	if statusCode == http.StatusInternalServerError {
		log.Printf("Internal server error: %v", err)
		errorResponse.Message = "サーバー内部エラーが発生しました"
		errorResponse.Error = err.Error()
	}

	// JSONレスポンスの送信
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		log.Printf("Failed to encode error response: %v", err)
	}
}

// determineStatusCode はエラーに応じたHTTPステータスコードを決定します
func determineStatusCode(err error) int {
	switch {
	case errors.Is(err, ErrBadRequest):
		return http.StatusBadRequest
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

// ResponseJSON はJSONレスポンスを送信するヘルパー関数です
func ResponseJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	// ステータスコードのデフォルト値
	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	// JSONレスポンスの送信
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode response: %v", err)
		HTTPErrorHandler(w, err, http.StatusInternalServerError)
	}
}

// SuccessResponse は成功レスポンスを表します
type SuccessResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// RespondWithSuccess は成功レスポンスを送信するヘルパー関数です
func RespondWithSuccess(w http.ResponseWriter, data interface{}, message string, statusCode int) {
	// ステータスコードのデフォルト値
	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	// 成功レスポンスの作成
	successResponse := SuccessResponse{
		Status:  statusCode,
		Message: message,
		Data:    data,
	}

	// JSONレスポンスの送信
	ResponseJSON(w, successResponse, statusCode)
}
