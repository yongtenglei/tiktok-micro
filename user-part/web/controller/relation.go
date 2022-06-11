package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"yunyandz.com/tiktok/pkg/constant"
	"yunyandz.com/tiktok/pkg/util"
	"yunyandz.com/tiktok/proto/pb"
	"yunyandz.com/tiktok/user-part/web/common"
)

type RelationRequest struct {
	Token      string `form:"token" binding:"required"`
	ToUserId   uint64 `form:"to_user_id" binding:"required"`
	ActionType int32  `form:"action_type" binding:"required"`
}

type RelationResponse struct {
	common.Response
}

func RelationAction(c *gin.Context) {
	var req RelationRequest
	var rsp RelationResponse

	err := c.ShouldBind(&req)
	uc, _ := util.GetUserClaims(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Response{StatusCode: constant.CodeFail, StatusMsg: err.Error()})
		return
	}
	if req.ToUserId == uc.UserID {
		c.JSON(http.StatusOK, common.Response{StatusCode: constant.CodeFail, StatusMsg: "不能对自己进行操作"})
		return
	}

	var reqPB pb.FollowReq
	reqPB.UserId = int64(uc.UserID)
	reqPB.ToUserId = int64(req.ToUserId)
	reqPB.ActionType = req.ActionType
	reqPB.Token = req.Token

	res, err := UserServiceClient.Follow(c.Copy(), &reqPB)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Response{StatusCode: constant.CodeFail, StatusMsg: err.Error()})
		return
	}

	rsp.Response = common.Response{
		StatusCode: res.StatusCode,
		StatusMsg:  res.StatusMsg,
	}

	c.JSON(http.StatusOK, rsp)
}

type FollowListRequest struct {
	UserId uint64 `form:"user_id" binding:"required"`
	Token  string `form:"token" binding:"required"`
}

type FollowListResponse struct {
	common.Response
	UserList []*pb.User `json:"user_list"`
}

func FollowList(c *gin.Context) {
	var req FollowListRequest
	var rsp FollowListResponse
	err := c.ShouldBindQuery(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Response{StatusCode: constant.CodeFail, StatusMsg: err.Error()})
		return
	}

	var reqPB pb.FollowListReq
	reqPB.UserId = int64(req.UserId)
	reqPB.Token = req.Token

	res, err := UserServiceClient.FollowList(c.Copy(), &reqPB)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Response{StatusCode: constant.CodeFail, StatusMsg: err.Error()})
		return
	}

	rsp.Response = common.Response{
		StatusCode: res.StatusCode,
		StatusMsg:  res.StatusMsg,
	}
	rsp.UserList = res.UserList

	c.JSON(http.StatusOK, rsp)
}

type FollowerListRequest struct {
	UserId uint64 `form:"user_id" binding:"required"`
	Token  string `form:"token" binding:"required"`
}

type FollowerListResponse struct {
	common.Response
	UserList []*pb.User
}

func FollowerList(c *gin.Context) {
	var req FollowerListRequest
	var rsp FollowerListResponse
	err := c.ShouldBindQuery(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Response{StatusCode: constant.CodeFail, StatusMsg: err.Error()})
		return
	}

	var reqPb pb.FollowerListReq
	reqPb.UserId = int64(req.UserId)
	reqPb.Token = req.Token

	res, err := UserServiceClient.FollowerList(c.Copy(), &reqPb)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Response{StatusCode: constant.CodeFail, StatusMsg: err.Error()})
		return
	}

	rsp.Response = common.Response{
		StatusCode: res.StatusCode,
		StatusMsg:  res.StatusMsg,
	}
	rsp.UserList = res.UserList

	c.JSON(http.StatusOK, rsp)
}
