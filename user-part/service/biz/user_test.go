package biz

import (
	"context"
	"flag"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"yunyandz.com/tiktok/dao/mysql"
	_ "yunyandz.com/tiktok/dao/mysql"
	"yunyandz.com/tiktok/pkg/constant"
	"yunyandz.com/tiktok/proto/pb"
	"yunyandz.com/tiktok/user-part/settings"
)

var userTestServer UserServiceServer

func init() {
	configPath := flag.String("f", "setting.yaml", "config file path")
	settings.ParseConfig(*configPath)
	settings.ParseConfig("setting.yaml")

	mysql.Init()

}

func TestService_Register(t *testing.T) {
	for i := 0; i < 10; i++ {
		s := strconv.Itoa(i)
		response, err := userTestServer.Register(context.Background(), &pb.RegisterReq{
			Username: s,
			Password: s,
		})
		require.NoError(t, err)
		require.NotEmpty(t, response)
		require.NotEmpty(t, response.Token)
		require.Equal(t, int64(i+1), response.UserId)
		require.Equal(t, int32(constant.CodeSuccess), response.StatusCode)
		require.Equal(t, constant.MsgSuccess, response.StatusMsg)
	}

}

func TestService_Login(t *testing.T) {
	for i := 0; i < 10; i++ {
		s := strconv.Itoa(i)
		response, err := userTestServer.Login(context.Background(), &pb.LoginReq{
			Username: s,
			Password: s,
		})
		require.NoError(t, err)
		require.NotEmpty(t, response)
		require.NotEmpty(t, response.Token)
		require.Equal(t, int64(i+1), response.UserId)
		require.Equal(t, int32(constant.CodeSuccess), response.StatusCode)
		require.Equal(t, constant.MsgSuccess, response.StatusMsg)
	}

}

func TestService_GetUserInfo(t *testing.T) {
	for i := 0; i < 10; i++ {
		s := strconv.Itoa(i)
		response, err := userTestServer.Info(context.Background(), &pb.InfoReq{
			UserId: int64(i + 1),
			Token:  "",
		})
		require.NoError(t, err)
		require.NotEmpty(t, response)
		require.Equal(t, int64(i+1), response.User.Id)
		require.Equal(t, s, response.User.Name)

		// have not give real value which is init value now
		require.Zero(t, response.User.FollowCount)
		require.Zero(t, response.User.FollowerCount)
		require.Equal(t, false, response.User.IsFollow)

		require.Equal(t, int32(constant.CodeSuccess), response.StatusCode)
		require.Equal(t, constant.MsgSuccess, response.StatusMsg)
	}
}

func TestUserServiceServer_Follow(t *testing.T) {
	for i := 0; i < 9; i++ {
		response, err := userTestServer.Follow(context.Background(), &pb.FollowReq{
			UserId:     int64(i + 1),
			Token:      "",
			ToUserId:   int64(i + 2),
			ActionType: constant.FollowAction,
		})
		require.NoError(t, err)
		require.NotEmpty(t, response)

		require.Equal(t, int32(constant.CodeSuccess), response.StatusCode)
		require.Equal(t, constant.MsgSuccess, response.StatusMsg)
	}
}

func TestUserServiceServer_Follow_Case_OneFollowMany(t *testing.T) {
	for i := 3; i < 9; i++ {
		response, err := userTestServer.Follow(context.Background(), &pb.FollowReq{
			UserId:     1,
			Token:      "",
			ToUserId:   int64(i),
			ActionType: constant.FollowAction,
		})
		require.NoError(t, err)
		require.NotEmpty(t, response)

		require.Equal(t, int32(constant.CodeSuccess), response.StatusCode)
		require.Equal(t, constant.MsgSuccess, response.StatusMsg)
	}
}

func TestUserServiceServer_UnFollow(t *testing.T) {
	for i := 0; i < 9; i++ {
		response, err := userTestServer.Follow(context.Background(), &pb.FollowReq{
			UserId:     int64(i + 1),
			Token:      "",
			ToUserId:   int64(i + 2),
			ActionType: constant.UnFollowAction,
		})
		require.NoError(t, err)
		require.NotEmpty(t, response)

		require.Equal(t, int32(constant.CodeSuccess), response.StatusCode)
		require.Equal(t, constant.MsgSuccess, response.StatusMsg)
	}
}

func TestUserServiceServer_FollowList(t *testing.T) {
	for i := 1; i < 9; i++ {
		response, err := userTestServer.FollowList(context.Background(), &pb.FollowListReq{
			UserId: int64(i + 1),
			Token:  "",
		})

		require.NoError(t, err)
		require.NotEmpty(t, response)
		require.NotEmpty(t, response.UserList)
		require.Len(t, response.UserList, 1)

		require.Equal(t, int32(constant.CodeSuccess), response.StatusCode)
		require.Equal(t, constant.MsgSuccess, response.StatusMsg)

		for j := 0; j < len(response.UserList); j++ {
			require.Equal(t, int64(i+2), response.UserList[j].Id)
		}

	}

}

func TestUserServiceServer_FollowerList(t *testing.T) {
	for i := 1; i < 10; i++ {
		response, err := userTestServer.FollowerList(context.Background(), &pb.FollowerListReq{
			UserId: int64(i + 1),
			Token:  "",
		})

		require.NoError(t, err)
		require.NotEmpty(t, response)
		require.NotEmpty(t, response.UserList)
		//require.Len(t, response.UserList, 1)

		require.Equal(t, int32(constant.CodeSuccess), response.StatusCode)
		require.Equal(t, constant.MsgSuccess, response.StatusMsg)

		ids := make([]int64, 0)
		for j := 0; j < len(response.UserList); j++ {
			ids = append(ids, response.UserList[j].Id)
		}

		require.Contains(t, ids, int64(i))
	}

}
