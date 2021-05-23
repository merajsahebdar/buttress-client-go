package client

import (
	"context"

	"github.com/merajsahebdar/buttress-client-go/internal/app/auth"
	pb "github.com/merajsahebdar/buttress-implementation-go/rbac"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// RbacClient
type RbacClient struct {
	ai   *auth.AuthInterceptor
	ctx  context.Context
	conn *grpc.ClientConn
	svc  pb.RbacServiceClient
}

// NewAuthorizedRbacClient
func NewAuthorizedRbacClient(addr string, app string, pem []byte) (*RbacClient, *ClientError) {
	ai, err := auth.NewAuthInterceptor(app, pem)
	if err != nil {
		return nil, &ClientError{Type: TokenGenerationError, Err: err}
	}

	opts := []grpc.DialOption{
		grpc.WithUnaryInterceptor(ai.Unary()),
		grpc.WithStreamInterceptor(ai.Stream()),
	}

	return newRbacClient(addr, opts, app)
}

// NewRbacClient
func NewRbacClient(addr string, app string) (*RbacClient, *ClientError) {
	return newRbacClient(addr, []grpc.DialOption{}, app)
}

// newRbacClient
func newRbacClient(addr string, opts []grpc.DialOption, app string) (*RbacClient, *ClientError) {
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, &ClientError{Type: ConnectionError, Err: err}
	}

	md := metadata.New(
		map[string]string{
			"x-app": app,
		},
	)
	ctx := metadata.NewOutgoingContext(
		context.Background(),
		md,
	)
	svc := pb.NewRbacServiceClient(conn)

	_, err = svc.CreateRbacInstance(ctx, &pb.EmptyRequest{})
	if err != nil {
		return nil, &ClientError{Type: InstanceCreationError, Err: err}
	}

	return &RbacClient{
		ctx:  ctx,
		svc:  svc,
		conn: conn,
	}, nil
}

// HasPermission
func (c *RbacClient) HasPermission(subject string, object string, action string) (bool, error) {
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
		return false, err
	}

	return res.Has, nil
}

// GrantPermissionToSubject
func (c *RbacClient) GrantPermission(subject string, object string, action string) error {
	_, err := c.svc.GrantPermission(
		c.ctx,
		&pb.GrantPermissionRequest{
			Subject: subject,
			Permission: &pb.PermissionDefinition{
				Object: object,
				Action: action,
			},
		})

	return err
}

// RevokePermission
func (c *RbacClient) RevokePermission(subject string, object string, action string) error {
	_, err := c.svc.RevokePermission(
		c.ctx,
		&pb.RevokePermissionRequest{
			Subject: subject,
			Permission: &pb.PermissionDefinition{
				Object: object,
				Action: action,
			},
		})

	return err
}

// AddChildSubjectToParentSubject
func (c *RbacClient) AddChildSubjectToParentSubject(child string, parent string) error {
	_, err := c.svc.AddChildSubjectToParentSubject(
		c.ctx,
		&pb.AddChildSubjectToParentSubjectRequest{
			Child:  child,
			Parent: parent,
		},
	)

	return err
}

// RemoveChildSubjectFromParentSubject
func (c *RbacClient) RemoveChildSubjectFromParentSubject(child string, parent string) error {
	_, err := c.svc.RemoveChildSubjectFromParentSubject(
		c.ctx,
		&pb.RemoveChildSubjectFromParentSubjectRequest{
			Child:  child,
			Parent: parent,
		},
	)

	return err
}
