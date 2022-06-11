package controller

import (
	"fmt"
	"yunyandz.com/tiktok/user-part/settings"
	"yunyandz.com/tiktok/user-part/web/common"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net/http"
	"yunyandz.com/tiktok/pkg/constant"
	"yunyandz.com/tiktok/proto/pb"

	"github.com/gin-gonic/gin"
)

var UserServiceClient pb.UserServiceClient

func InitUserServiceClient() {

	addr := fmt.Sprintf("consul://%s:%d/%s?wait=14s",
		settings.ServiceConf.ConsulConf.Host,
		settings.ServiceConf.ConsulConf.Port,
		settings.ServiceConf.UserWebServerConf.Name)
	//"consul://127.0.0.1:8500/user_web_server?wait=14s"

	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	if err != nil {
		zap.S().Errorw("grpc Dial failed", "err", err.Error())

		panic(err)
	}

	UserServiceClient = pb.NewUserServiceClient(conn)

}

type RegisterRequest struct {
	Username string `form:"username" binding:"required,min=3,max=32"`
	Password string `form:"password" binding:"required,max=32"`
}

type RegisterResponse struct {
	common.Response
	UserID uint64 `json:"user_id"`
	Token  string `json:"token"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	var rsp RegisterResponse
	err := c.ShouldBindQuery(&req)
	if err != nil {
		rsp.Response = common.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		}
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	var reqPB pb.RegisterReq
	reqPB.Username = req.Username
	reqPB.Password = req.Password

	r, err := UserServiceClient.Register(c.Copy(), &reqPB)
	if err != nil {
		rsp.Response = common.Response{
			StatusCode: constant.CodeFail,
			StatusMsg:  err.Error(),
		}
		c.JSON(http.StatusInternalServerError, rsp)
		return
	}

	rsp.Response = common.Response{
		StatusCode: r.StatusCode,
		StatusMsg:  r.StatusMsg,
	}
	rsp.UserID = uint64(r.UserId)
	rsp.Token = r.Token
	c.JSON(http.StatusOK, rsp)
}

type LoginRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type LoginResponse struct {
	common.Response
	UserID uint64 `json:"user_id"`
	Token  string `json:"token"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	var rsp LoginResponse
	err := c.ShouldBindQuery(&req)
	if err != nil {
		rsp.Response = common.Response{
			StatusCode: constant.CodeFail,
			StatusMsg:  err.Error(),
		}
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	var reqPB pb.LoginReq
	reqPB.Username = req.Username
	reqPB.Password = req.Password

	r, err := UserServiceClient.Login(c.Copy(), &reqPB)
	if err != nil {
		rsp.Response = common.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		}
		c.JSON(http.StatusInternalServerError, rsp)
		return
	}

	rsp.Response = common.Response{
		StatusCode: r.StatusCode,
		StatusMsg:  r.StatusMsg,
	}
	rsp.UserID = uint64(r.UserId)
	rsp.Token = r.Token
	c.JSON(http.StatusOK, rsp)
}

type UserInfoRequest struct {
	UserID uint64 `form:"user_id" binding:"required"`
	Token  string `form:"token" binding:"required"`
}

type UserInfoResponse struct {
	common.Response
	*pb.User `json:"user"`
}

func UserInfo(c *gin.Context) {
	var req UserInfoRequest
	var rsp UserInfoResponse
	err := c.ShouldBindQuery(&req)
	if err != nil {
		rsp.Response = common.Response{
			StatusCode: constant.CodeFail,
			StatusMsg:  err.Error(),
		}
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	var reqPB pb.InfoReq
	reqPB.UserId = int64(req.UserID)
	reqPB.Token = req.Token

	r, err := UserServiceClient.Info(c.Copy(), &reqPB)
	if err != nil {
		rsp.Response = common.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		}
		c.JSON(http.StatusInternalServerError, rsp)
		return
	}

	rsp.Response = common.Response{
		StatusCode: r.StatusCode,
		StatusMsg:  r.StatusMsg,
	}
	rsp.User = r.User
	c.JSON(http.StatusOK, rsp)

}
