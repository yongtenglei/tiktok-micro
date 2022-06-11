package biz

import (
	"context"
	"errors"
	"fmt"
	"yunyandz.com/tiktok/dao/mysql"
	"yunyandz.com/tiktok/pkg/constant"
	"yunyandz.com/tiktok/user-part/model"

	"yunyandz.com/tiktok/logger"
	"yunyandz.com/tiktok/pkg/errorx"
	"yunyandz.com/tiktok/pkg/jwtx"
	"yunyandz.com/tiktok/pkg/scryptx"
	"yunyandz.com/tiktok/proto/pb"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
}

func (us *UserServiceServer) Register(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterRes, error) {
	var res pb.RegisterRes
	res.StatusCode = constant.CodeFail

	user, err := mysql.GetUserByName(req.Username)
	if err == nil {
		if user.ID > 0 {
			res.StatusMsg = errorx.ErrUserAlreadyExists.Error()
			return &res, errorx.ErrUserAlreadyExists
		}
		return &res, err
	}

	u := model.User{
		Username: req.Username,
		Password: scryptx.PasswordEncrypt(req.Password),
	}

	id, err := mysql.CreateUser(&u)
	if err != nil {
		logger.Suger().Errorw("Register save failed", "err", err.Error())
		res.StatusMsg = err.Error()
		return &res, errorx.ErrInternalBusy
	}

	token, err := jwtx.CreateUserClaims(jwtx.UserInfo{
		Username: req.Username,
		UserID:   id,
	})
	if err != nil {
		logger.Suger().Errorw("Register CreateUserClaims failed", "err", err.Error())
		res.StatusMsg = err.Error()
		return &res, errorx.ErrInternalBusy
	}

	res.StatusCode = constant.CodeSuccess
	res.StatusMsg = constant.MsgSuccess
	res.UserId = int64(id)
	res.Token = token

	return &res, nil
}

func (us *UserServiceServer) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginRes, error) {
	var res pb.LoginRes
	res.StatusCode = constant.CodeSuccess

	user, err := mysql.GetUserByName(req.Username)
	if err != nil {
		res.StatusMsg = err.Error()
		return &res, errorx.ErrUserDoesNotExists
	}

	if !scryptx.PasswordValidate(req.Password, user.Password) {
		res.StatusMsg = err.Error()
		return &res, errorx.ErrUserPassword
	}

	token, err := jwtx.CreateUserClaims(jwtx.UserInfo{
		Username: req.Username,
		UserID:   uint64(user.ID),
	})
	if err != nil {
		logger.Suger().Errorw("Login CreateUserClaims failed", "err", err.Error())
		res.StatusMsg = err.Error()
		return &res, errorx.ErrInternalBusy
	}

	res.StatusCode = constant.CodeSuccess
	res.StatusMsg = constant.MsgSuccess
	res.UserId = int64(user.ID)
	res.Token = token

	return &res, nil
}

func (us *UserServiceServer) Info(ctx context.Context, req *pb.InfoReq) (*pb.InfoRes, error) {
	var res pb.InfoRes
	res.StatusCode = constant.CodeSuccess

	user, err := mysql.GetUser(uint64(req.UserId))
	if err != nil {
		res.StatusMsg = err.Error()
		return &res, errorx.ErrUserDoesNotExists
	}

	u, err := convertUserToPB(uint64(req.UserId), user)
	if err != nil {
		res.StatusMsg = err.Error()
		return &res, errorx.ErrInternalBusy
	}

	res.StatusCode = constant.CodeSuccess
	res.StatusMsg = constant.MsgSuccess
	res.User = u

	return &res, nil
}
func (us *UserServiceServer) Follow(ctx context.Context, req *pb.FollowReq) (*pb.FollowRes, error) {
	var isOk = 0 // 0 is ok, -1 is reverse
	var err error

	switch {
	case req.ActionType == constant.FollowAction:
		err = mysql.CreateFollow(uint64(req.UserId), uint64(req.ToUserId))
		if err != nil {
			isOk = -1

		}

	case req.ActionType == constant.UnFollowAction:
		err = mysql.DeleteFollow(uint64(req.UserId), uint64(req.ToUserId))
		if err != nil {
			isOk = -1
		}
	default:
		isOk = -1
		err = errors.New("invalid action")
	}

	var res pb.FollowRes
	if isOk == 0 {
		res.StatusCode = constant.CodeSuccess
		res.StatusMsg = constant.MsgSuccess

	} else if isOk == -1 {
		res.StatusCode = constant.CodeFail
		res.StatusMsg = err.Error()
	}

	return &res, err

}

func (us *UserServiceServer) FollowList(ctx context.Context, req *pb.FollowListReq) (*pb.FollowListRes, error) {
	var res pb.FollowListRes

	followList, err := mysql.GetFollowList(uint64(req.UserId))
	if err != nil {
		res.StatusCode = constant.CodeFail
		res.StatusMsg = err.Error()
		return &res, nil
	}
	fmt.Printf("list: %#v\n", followList)
	users := convertManyUserModelToUser(uint64(req.UserId), followList)

	res.StatusCode = constant.CodeSuccess
	res.StatusMsg = constant.MsgSuccess
	res.UserList = users

	fmt.Printf("list: %#v\n", users)

	return &res, nil

}

func (us *UserServiceServer) FollowerList(ctx context.Context, req *pb.FollowerListReq) (*pb.FollowerListRes, error) {
	var res pb.FollowerListRes
	followList, err := mysql.GetFollowerList(uint64(req.UserId))
	if err != nil {
		res.StatusCode = constant.CodeFail
		res.StatusMsg = err.Error()
		return &res, nil
	}

	users := convertManyUserModelToUser(uint64(req.UserId), followList)

	res.StatusCode = constant.CodeSuccess
	res.StatusMsg = constant.MsgSuccess
	res.UserList = users

	return &res, nil
}

func convertManyUserModelToUser(selfId uint64, userList []*model.User) []*pb.User {
	var users = make([]*pb.User, 0, len(userList))
	for _, item := range userList {
		user, err := convertUserToPB(selfId, item)
		if err != nil {
			continue
		}
		users = append(users, user)
	}
	return users
}

func convertUserToPB(selfId uint64, user *model.User) (*pb.User, error) {
	isFollow, err := mysql.IsFollow(selfId, uint64(user.Model.ID))
	if err != nil {
		logger.Suger().Debugf("convertUserToPB isFollow err: %s", err.Error())
		return &pb.User{}, err
	}
	//logger.Suger().Debugf("convertUserToPB isFollow: %d->%d", selfId, user.ID)
	u := pb.User{
		Id:            int64(user.ID),
		Name:          user.Username,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      isFollow,
	}
	return &u, nil
}
