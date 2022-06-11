package controller

//
//import (
//	"net/http"
//
//	"github.com/gin-gonic/gin"
//	"github.com/yunyandz/tiktok-demo-backend/internal/service"
//)
//
//type FavoriteActionRequest struct {
//	VideoId    uint64 `form:"video_id" banding:"required"`
//	ActionType string `form:"action_type" banding:"required"`
//}
//
//const (
//	FavoriteActionTypeLike    = "1"
//	FavoriteActionTypeDislike = "2"
//)
//
//// FavoriteAction 用户点赞操作
//func (ctl *Controller) FavoriteAction(c *gin.Context) {
//	// token := c.Query("token")
//	req := &FavoriteActionRequest{}
//	rsp := &service.Response{}
//	err := c.ShouldBindQuery(&req)
//	uc, _ := ctl.getUserClaims(c)
//	like := false
//	if err != nil {
//		rsp.StatusCode = -1
//		rsp.StatusMsg = err.Error()
//		c.JSON(http.StatusOK, rsp)
//		return
//	}
//	switch req.ActionType {
//	case FavoriteActionTypeLike:
//		like = true
//	case FavoriteActionTypeDislike:
//		like = false
//	default:
//		rsp.StatusCode = -1
//		rsp.StatusMsg = "action_type error"
//		c.JSON(http.StatusOK, rsp)
//		return
//	}
//	ctl.logger.Sugar().Debugf("user_id: %d, video_id: %d, action_type: %s,like: %s", uc.UserID, req.VideoId, req.ActionType, like)
//	rsp = ctl.service.LikeDisliakeVideo(uc.UserID, req.VideoId, like)
//	c.JSON(http.StatusOK, *rsp)
//}
//
//type FavoriteListResponse struct {
//	service.Response
//	service.VideoListResponse
//}
//
//// FavoriteList 用户获取点赞列表
//func (ctl *Controller) FavoriteList(c *gin.Context) {
//	req := &UserInfoRequest{}
//	rsp := &FavoriteListResponse{}
//	err := c.ShouldBindQuery(&req)
//	if err != nil {
//		c.JSON(http.StatusOK, service.Response{
//			StatusCode: -1,
//			StatusMsg:  err.Error(),
//		})
//	}
//	rsp.VideoListResponse = *ctl.service.GetLikeList(req.UserID)
//	c.JSON(http.StatusOK, rsp)
//}
