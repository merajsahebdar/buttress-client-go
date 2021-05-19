package client

import (
	"context"

	"github.com/merajsahebdar/buttress-client-go/internal/app/auth"
	pb "github.com/merajsahebdar/buttress-implementation-go/rbac"
	"google.golang.org/grpc"
)

// RbacClient
type RbacClient struct {
	ai   *auth.AuthInterceptor
	ctx  context.Context
	conn *grpc.ClientConn
	svc  pb.RbacServiceClient
}

// NewRbacClient
func NewRbacClient(addr string, uuid string, pem []byte) (*RbacClient, *ClientError) {
	ai, err := auth.NewAuthInterceptor(uuid, pem)
	if err != nil {
		return nil, &ClientError{Type: TokenGenerationError, Err: err}
	}

	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(ai.Unary()),
		grpc.WithStreamInterceptor(ai.Stream()),
	)
	if err != nil {
		return nil, &ClientError{Type: ConnectionError, Err: err}
	}

	ctx := context.Background()
	svc := pb.NewRbacServiceClient(conn)

	_, err = svc.CreateRbacInstance(ctx, &pb.EmptyRequest{})
	if err != nil {
		return nil, &ClientError{Type: InstanceCreationError, Err: err}
	}

	ac := &RbacClient{
		ai:   ai,
		ctx:  ctx,
		svc:  svc,
		conn: conn,
	}

	return ac, nil
}

// HasPermission
func (c *RbacClient) HasPermission(subject string, object string, action string) (*pb.HasPermissionResponse, error) {
	res, err := c.svc.HasPermission(
		c.ctx,
		&pb.HasPermissionRequest{
			Subject: subject,
			Permission: &pb.PermissionDefinition{
				Object: object,
				Action: action,
			},
		})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GrantPermissionToSubject
func (c *RbacClient) GrantPermissionToSubject(subject string, object string, action string) (*pb.EmptyResponse, error) {
	res, err := c.svc.GrantPermissionToSubject(
		c.ctx,
		&pb.GrantPermissionToSubjectRequest{
			Subject: subject,
			Permission: &pb.PermissionDefinition{
				Object: object,
				Action: action,
			},
		})
	if err != nil {
		return nil, err
	}

	return res, nil
}
