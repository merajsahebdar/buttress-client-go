package auth

import (
	"context"
	"crypto/rsa"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// AuthOptions
type AuthOptions struct {
	Claims     *jwt.StandardClaims
	PrivateKey []byte
}

// AuthInterceptor
type AuthInterceptor struct {
	token string
}

// signJwt
func (ai *AuthInterceptor) signJwt(privateKey *rsa.PrivateKey, c *jwt.StandardClaims) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("RS256"))
	t.Claims = c
	return t.SignedString(privateKey)
}

// attachAuthorization
func (ai *AuthInterceptor) attachAuthorization(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", ai.token)
}

// Unary
func (ai *AuthInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		return invoker(ai.attachAuthorization(ctx), method, req, reply, cc, opts...)
	}
}

// Stream
func (ai *AuthInterceptor) Stream() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		return streamer(ai.attachAuthorization(ctx), desc, cc, method, opts...)
	}
}

// NewAuthInterceptor
func NewAuthInterceptor(uuid string, pem []byte) (*AuthInterceptor, error) {
	ai := &AuthInterceptor{}

	var err error
	var key *rsa.PrivateKey

	key, err = jwt.ParseRSAPrivateKeyFromPEM(pem)
	if err != nil {
		return nil, err
	}

	ai.token, err = ai.signJwt(
		key,
		&jwt.StandardClaims{
			Subject: uuid,
		},
	)
	if err != nil {
		return nil, err
	}

	return ai, nil
}
