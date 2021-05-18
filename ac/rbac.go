package buttressac

import (
	"context"

	buttresserror "github.com/merajsahebdar/buttress-client-go/error"
	buttressinterceptor "github.com/merajsahebdar/buttress-client-go/interceptor"
	pb "github.com/merajsahebdar/buttress-implementation-go/rbac"
	"google.golang.org/grpc"
)

// RbacAc
type RbacAc struct {
	ai     *buttressinterceptor.AuthInterceptor
	ctx    context.Context
	conn   *grpc.ClientConn
	client pb.RbacServiceClient
}

func NewRbacAc(addr string, uuid string, pem []byte) (*RbacAc, *buttresserror.AcError) {
	ai, err := buttressinterceptor.NewAuthInterceptor(uuid, pem)
	if err != nil {
		return nil, &buttresserror.AcError{Type: buttresserror.TokenGenerationError, Err: err}
	}

	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(ai.Unary()),
		grpc.WithStreamInterceptor(ai.Stream()),
	)
	if err != nil {
		return nil, &buttresserror.AcError{Type: buttresserror.ConnectionError, Err: err}
	}

	ctx := context.Background()
	client := pb.NewRbacServiceClient(conn)

	_, err = client.CreateRbacInstance(ctx, &pb.EmptyRequest{})
	if err != nil {
		return nil, &buttresserror.AcError{Type: buttresserror.InstanceCreationError, Err: err}
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
