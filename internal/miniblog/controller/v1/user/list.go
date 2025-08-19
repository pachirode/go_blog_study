package user

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/pachirode/go_blog_study/internal/pkg/log"
	pb "github.com/pachirode/go_blog_study/pkg/proto/miniblog/v1"
)

func (ctrl *UserController) ListUser(ctx context.Context, r *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	log.C(ctx).Infow("List user function called")

	resp, err := ctrl.b.Users().List(ctx, int(r.Offset), int(r.Limit))
	if err != nil {
		return nil, err
	}

	users := make([]*pb.UserInfo, len(resp.Users))
	for _, user := range resp.Users {
		createdAt, _ := time.Parse("2001-01-01 01:00:00", user.CreatedAt)
		updatedAt, _ := time.Parse("2001-01-01 01:00:00", user.UpdatedAt)
		users = append(users, &pb.UserInfo{
			Username:  user.Username,
			Nickname:  user.Nickname,
			Email:     user.Email,
			Phone:     user.Phone,
			PostCount: user.PostCount,
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: timestamppb.New(updatedAt),
		})
	}

	ret := &pb.ListUserResponse{TotalCount: resp.TotalCount, Users: users}
	return ret, nil
}
