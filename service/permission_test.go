package service_test

import (
	"context"
	"fmt"
	"net"
	"strings"
	"testing"

	"io/ioutil"

	pb "github.com/meateam/permission-service/proto"

	"github.com/meateam/permission-service/internal/test"
	"github.com/meateam/permission-service/server"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

// Global variables
var (
	s         *server.PermissionServer
	logger    = logrus.New()
	lis       *bufconn.Listener
	authToken string
)

// Constants
const bufSize = 1024 * 1024
const fileID = "mongoID"
const userID = "userID"
const creator = "creatorID"
const appID = "appID"

func init() {
	lis = bufconn.Listen(bufSize)

	// Disable log output.
	logger.SetOutput(ioutil.Discard)

	s = server.NewServer(logger)
	go func() {
		s.Serve(lis)
	}()

	var err error
	authToken, err = test.GenerateJwtToken()
	fmt.Printf("jwt token: %s \n", authToken)
	if err != nil {
		fmt.Printf("Error signing jwt token: %s \n", err)
	}

}

func Test_CreatePermission(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.CreatePermissionRequest
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    pb.PermissionObject
	}{
		{
			name: "create Role_NONE",
			args: args{
				ctx: context.Background(),
				req: &pb.CreatePermissionRequest{
					FileID:   fileID + ":Role_NONE",
					UserID:   userID + ":Role_NONE",
					Role:     pb.Role_NONE,
					Creator:  "creator",
					Override: false,
					AppID:    appID,
				},
			},
			wantErr: false,
		},
		{
			name: "create Role_WRITE",
			args: args{
				ctx: context.Background(),
				req: &pb.CreatePermissionRequest{
					FileID:   fileID + ":Role_WRITE",
					UserID:   userID + ":Role_WRITE",
					Role:     pb.Role_WRITE,
					Creator:  "creator",
					Override: false,
					AppID:    appID,
				},
			},
			wantErr: false,
		},
		{
			name: "create Role_READ",
			args: args{
				ctx: context.Background(),
				req: &pb.CreatePermissionRequest{
					FileID:   fileID + ":Role_READ",
					UserID:   userID + ":Role_READ",
					Role:     pb.Role_READ,
					Creator:  "creator",
					Override: false,
					AppID:    appID,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
			if err != nil {
				t.Fatalf("failed to dial bufnet: %v", err)
			}
			defer conn.Close()

			t.Parallel()

			client := pb.NewPermissionClient(conn)
			permissionObject, err := client.CreatePermission(tt.args.ctx, tt.args.req)

			// Unanticipated error - isn't related to tt.wantErr
			if err != nil {
				t.Fatalf("DownloadService.Download() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !comperePermissionObject(permissionObject, tt.args.req) {
				t.Errorf(
					"PermissionService.CreatePermission() Fails to create permission, wantErr %v", tt.wantErr,
				)
			}

		})
	}
}

func Test_GetFilePermissions(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.GetFilePermissionsRequest
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    pb.PermissionObject
	}{
		{
			name: "read file Role_NONE",
			args: args{
				ctx: context.Background(),
				req: &pb.GetFilePermissionsRequest{
					FileID: fileID + ":Role_NONE",
				},
			},
			wantErr: false,
		},
		{
			name: "read file Role_WRITE",
			args: args{
				ctx: context.Background(),
				req: &pb.GetFilePermissionsRequest{
					FileID: fileID + ":Role_WRITE",
				},
			},
			wantErr: false,
		},
		{
			name: "read file Role_READ",
			args: args{
				ctx: context.Background(),
				req: &pb.GetFilePermissionsRequest{
					FileID: fileID + ":Role_READ",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
			if err != nil {
				t.Fatalf("failed to dial bufnet: %v", err)
			}
			defer conn.Close()

			t.Parallel()

			client := pb.NewPermissionClient(conn)
			permissionObject, err := client.GetFilePermissions(tt.args.ctx, tt.args.req)

			// Unanticipated error - isn't related to tt.wantErr
			if err != nil {
				t.Fatalf("DownloadService.Download() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !compereFilePermissions(permissionObject, tt.args.req.FileID) {
				t.Errorf(
					"PermissionService.CreatePermission() Fails to create permission, wantErr %v", tt.wantErr,
				)
			}
		})
	}
}

func comperePermissionObject(compere *pb.PermissionObject, to *pb.CreatePermissionRequest) bool {

	if compere == nil {
		return false
	}

	if compere.Id == "" {
		return false
	}

	if compere.FileID == to.FileID &&
		compere.UserID == to.UserID &&
		compere.Role == to.Role &&
		compere.Creator == to.Creator &&
		compere.AppID == to.AppID {
		return true
	}

	return false
}

func compereFilePermissions(compere *pb.GetFilePermissionsResponse, to string) bool {

	if compere == nil {
		return false
	}

	filePermissions := compere.Permissions[0]

	if filePermissions == nil {
		return false
	}

	s := strings.Split(to, ":")

	if filePermissions.UserID == userID+":"+s[1] {
		return true
	}

	return false
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}
