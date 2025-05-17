package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

// DecodeJSON はリクエストボディからJSONをデコードするヘルパー関数です
func DecodeJSON(r *http.Request, v interface{}) error {
	// リクエストボディが空でないことを確認
	if r.Body == nil {
		return errors.New("リクエストボディが空です")
	}

	// Content-Typeヘッダーの確認
	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return errors.New("Content-Typeが'application/json'ではありません")
	}

	// JSONのデコード
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // 未知のフィールドを許可しない

	if err := decoder.Decode(v); err != nil {
		// デコードエラーの詳細を提供
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return errors.New("JSONの構文エラーです")
		case errors.As(err, &unmarshalTypeError):
			return errors.New("JSONの型が不正です")
		case errors.As(err, &invalidUnmarshalError):
			return errors.New("デコード先のオブジェクトが不正です")
		case errors.Is(err, io.EOF):
			return errors.New("リクエストボディが空です")
		default:
			return err
		}
	}

	// 追加のJSONがないことを確認
	if decoder.More() {
		return errors.New("リクエストボディに複数のJSONオブジェクトが含まれています")
	}

	return nil
}

// ValidateRequiredFields は必須フィールドの存在を確認するヘルパー関数です
func ValidateRequiredFields(data map[string]interface{}, requiredFields []string) error {
	for _, field := range requiredFields {
		value, exists := data[field]
		if !exists || isEmptyValue(value) {
			return errors.New("フィールド '" + field + "' は必須です")
		}
	}
	return nil
}

// isEmptyValue は値が空かどうかを判定するヘルパー関数です
func isEmptyValue(value interface{}) bool {
	if value == nil {
		return true
	}

	switch v := value.(type) {
	case string:
		return v == ""
	case int, int8, int16, int32, int64:
		return v == 0
	case float32, float64:
		return v == 0
	case bool:
		return !v
	case []interface{}:
		return len(v) == 0
	case map[string]interface{}:
		return len(v) == 0
	default:
		return false
	}
}

// ValidateEmail はメールアドレスの形式を検証するヘルパー関数です
func ValidateEmail(email string) error {
	// 簡易的な検証
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return errors.New("メールアドレスの形式が不正です")
	}
	return nil
}

// ValidatePassword はパスワードの強度を検証するヘルパー関数です
func ValidatePassword(password string) error {
	// 簡易的な検証
	if len(password) < 8 {
		return errors.New("パスワードは8文字以上必要です")
	}
	return nil
}
