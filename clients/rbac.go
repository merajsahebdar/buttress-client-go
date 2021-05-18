package clients

import (
	"context"

	interceptors "github.com/merajsahebdar/buttress-client-go/internal/interceptors"
	pb "github.com/merajsahebdar/buttress-implementation-go/rbac"
	"google.golang.org/grpc"
)

// RbacAc
type RbacAc struct {
	ai     *interceptors.AuthInterceptor
	ctx    context.Context
	conn   *grpc.ClientConn
	client pb.RbacServiceClient
}

func NewRbacAc(addr string, uuid string, pem []byte) (*RbacAc, *AcError) {
	ai, err := interceptors.NewAuthInterceptor(uuid, pem)
	if err != nil {
		return nil, &AcError{Type: TokenGenerationError, Err: err}
	}

	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(ai.Unary()),
		grpc.WithStreamInterceptor(ai.Stream()),
	)
	if err != nil {
		return nil, &AcError{Type: ConnectionError, Err: err}
	}

	ctx := context.Background()
	client := pb.NewRbacServiceClient(conn)

	_, err = client.CreateRbacInstance(ctx, &pb.EmptyRequest{})
	if err != nil {
		return nil, &AcError{Type: InstanceCreationError, Err: err}
	}

	ac := &RbacAc{
		ai:     ai,
		ctx:    ctx,
		conn:   conn,
		client: client,
	}

	return ac, nil
}

func (ac *RbacAc) HasPermission(subject string, object string, action string) (*pb.HasPermissionResponse, error) {
	res, err := ac.client.HasPermission(
		ac.ctx,
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

func (ac *RbacAc) GrantPermissionToSubject(subject string, object string, action string) (*pb.EmptyResponse, error) {
	res, err := ac.client.GrantPermissionToSubject(
		ac.ctx,
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
