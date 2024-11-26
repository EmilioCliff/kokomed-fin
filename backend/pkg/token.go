package pkg

import (
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Payload struct {
	ID       uuid.UUID `json:"id"`
	UserID   uint32    `json:"user_id"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
	BranchID uint32    `json:"branch_id"`
	jwt.RegisteredClaims
}

type JWTMaker struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

func NewJWTMaker(privateKeyPEM string, publicKeyPEM string) (*JWTMaker, error) {
	if privateKeyPEM == "" || publicKeyPEM == "" {
		return nil, Errorf(INTERNAL_ERROR, "private key or public key is empty")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyPEM))
	if err != nil {
		return nil, Errorf(INTERNAL_ERROR, "failed to parse private key: %v", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyPEM))
	if err != nil {
		return nil, Errorf(INTERNAL_ERROR, "failed to parse public key: %v", err)
	}

	maker := &JWTMaker{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

	return maker, nil
}

func (maker *JWTMaker) CreateToken(email string, userID, branchID uint32, role string, duration time.Duration) (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", Errorf(INTERNAL_ERROR, "failed to create uuid: %v", err)
	}

	claims := Payload{
		id,
		userID,
		email,
		role,
		branchID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "kokomedLoanApp",
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	token, err := jwtToken.SignedString(maker.PrivateKey)
	if err != nil {
		return "", Errorf(INTERNAL_ERROR, "failed to create token: %v", err)
	}

	return token, nil
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodRSA)
		if !ok {
			return nil, Errorf(INTERNAL_ERROR, "unexpected signing method")
		}

		return maker.PublicKey, nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		return nil, Errorf(INTERNAL_ERROR, "failed to parse token: %v", err)
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, Errorf(INTERNAL_ERROR, "failed to parse token is invalid")
	}

	if payload.RegisteredClaims.Issuer != "kokomedLoanApp" {
		return nil, Errorf(INTERNAL_ERROR, "invalid issuer")
	}

	if payload.RegisteredClaims.ExpiresAt.Time.Before(time.Now()) {
		return nil, Errorf(INTERNAL_ERROR, "token is expired")
	}

	return payload, nil
}

func (maker *JWTMaker) GetPayload(token string) (*Payload, error) {
	parser := jwt.NewParser()

	jwtToken, _, err := parser.ParseUnverified(token, &Payload{})
	if err != nil {
		return nil, Errorf(INTERNAL_ERROR, "failed to parse token: %v", err)
	}

	_, ok := jwtToken.Method.(*jwt.SigningMethodRSA)
	if !ok {
		return nil, Errorf(INTERNAL_ERROR, "unexpected signing method")
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, Errorf(INTERNAL_ERROR, "failed to parse token is invalid")
	}

	if payload.RegisteredClaims.Issuer != "kokomedLoanApp" {
		return nil, Errorf(INTERNAL_ERROR, "invalid issuer")
	}

	return payload, nil
}
