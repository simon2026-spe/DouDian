package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"doudian/internal/config"
)

// JWTClaims JWT 载荷
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Exp      int64  `json:"exp"` // 过期时间戳
	Iat      int64  `json:"iat"` // 签发时间戳
}

// jwtHeader JWT 头部
type jwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

// base64URLEncode base64url 编码（无填充）
func base64URLEncode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

// base64URLDecode base64url 解码（无填充）
func base64URLDecode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}

// GenerateToken 生成 JWT token，有效期 24 小时
func GenerateToken(userID uint, username string) (string, error) {
	cfg := config.Get()
	secret := []byte(cfg.JWTSecret)

	now := time.Now()
	exp := now.Add(24 * time.Hour).Unix()

	// 构造 header
	header := jwtHeader{
		Alg: "HS256",
		Typ: "JWT",
	}
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	headerEncoded := base64URLEncode(headerJSON)

	// 构造 payload
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Exp:      exp,
		Iat:      now.Unix(),
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	payloadEncoded := base64URLEncode(claimsJSON)

	// 构造签名消息：header.payload
	signingMessage := headerEncoded + "." + payloadEncoded

	// 计算 HMAC-SHA256 签名
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(signingMessage))
	signature := mac.Sum(nil)
	signatureEncoded := base64URLEncode(signature)

	// 组合完整 token
	token := signingMessage + "." + signatureEncoded
	return token, nil
}

// ParseToken 解析验证 token，返回 claims
func ParseToken(tokenString string) (*JWTClaims, error) {
	cfg := config.Get()
	secret := []byte(cfg.JWTSecret)

	// 分割 token 为三部分
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("无效的 token 格式")
	}

	headerEncoded := parts[0]
	payloadEncoded := parts[1]
	signatureEncoded := parts[2]

	// 验证签名
	signingMessage := headerEncoded + "." + payloadEncoded
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(signingMessage))
	expectedSignature := mac.Sum(nil)

	actualSignature, err := base64URLDecode(signatureEncoded)
	if err != nil {
		return nil, errors.New("签名解码失败")
	}

	if !hmac.Equal(expectedSignature, actualSignature) {
		return nil, errors.New("token 签名无效")
	}

	// 解码 header（可选验证 alg）
	headerJSON, err := base64URLDecode(headerEncoded)
	if err != nil {
		return nil, errors.New("header 解码失败")
	}
	var header jwtHeader
	if err := json.Unmarshal(headerJSON, &header); err != nil {
		return nil, errors.New("header 解析失败")
	}
	if header.Alg != "HS256" {
		return nil, errors.New("不支持的签名算法: " + header.Alg)
	}

	// 解码 payload
	claimsJSON, err := base64URLDecode(payloadEncoded)
	if err != nil {
		return nil, errors.New("payload 解码失败")
	}

	var claims JWTClaims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, errors.New("payload 解析失败")
	}

	// 验证过期时间
	if claims.Exp < time.Now().Unix() {
		return nil, errors.New("token 已过期")
	}

	return &claims, nil
}
