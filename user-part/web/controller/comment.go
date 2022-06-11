package controller

//
//import (
//	"net/http"
//
//	"github.com/gin-gonic/gin"
//	"github.com/yunyandz/tiktok-demo-backend/internal/errorx"
//	"github.com/yunyandz/tiktok-demo-backend/internal/jwtx"
//	"github.com/yunyandz/tiktok-demo-backend/internal/service"
//)
//
//type CommentActionRequest struct {
//	Token       string `form:"token" binding:"required"`
//	VideoID     uint64 `form:"video_id" binding:"required"`
//	ActionType  uint64 `form:"action_type" binding:"required"` //
//	CommentText string `form:"comment_text"`
//	CommentID   uint64 `form:"comment_id"`
//}
//type CommentActionResponse struct {
//	service.Response
//	Comment service.Comment `json:"comment"`
//}
//
//const (
//	CommentActionTypePublish = 1
//	CommentActionTypeDelete  = 2
//)
//
//func (ctl *Controller) CommentAction(c *gin.Context) {
//	// token := c.Query("token")
//	var req CommentActionRequest
//	var rsp CommentActionResponse
//	err := c.ShouldBindQuery(&req)
//	if err != nil {
//		rsp.Response = service.Response{
//			StatusCode: -1,
//			StatusMsg:  err.Error(),
//		}
//		ctl.logger.Sugar().Errorf("bind query failed: %s", err.Error())
//		c.JSON(http.StatusBadRequest, rsp)
//		return
//	}
//	uc, e := c.Get("claims")
//	if !e {
//		ctl.logger.Sugar().Errorf("Get claims error: %v", e)
//		c.JSON(http.StatusUnauthorized, service.Response{
//			StatusCode: -1,
//			StatusMsg:  errorx.ErrInvalidToken.Error(),
//		})
//		return
//	}
//	claims := uc.(*jwtx.UserClaims)
//	switch req.ActionType {
//	case CommentActionTypePublish:
//		r, err := ctl.service.PublishComment(claims.UserID, req.VideoID, req.CommentText)
//		if err != nil {
//			rsp.Response = service.Response{
//				StatusCode: -1,
//				StatusMsg:  err.Error(),
//			}
//			ctl.logger.Sugar().Errorf("publish comment failed: %s", err.Error())
//			c.JSON(http.StatusBadRequest, rsp)
//			return
//		}
//		rsp.Response = r.Response
//		rsp.Comment = r.Comment
//		c.JSON(http.StatusOK, rsp)
//	case CommentActionTypeDelete:
//		r, err := ctl.service.DeleteComment(req.CommentID)
//		if err != nil {
//			rsp.Response = service.Response{
//				StatusCode: -1,
//				StatusMsg:  err.Error(),
//			}
//			ctl.logger.Sugar().Errorf("delete comment failed: %s", err.Error())
//			c.JSON(http.StatusBadRequest, rsp)
//			return
//		}
//		rsp.Response = r.Response
//		rsp.Comment = r.Comment
//		c.JSON(http.StatusOK, rsp)
//	default:
//		rsp.Response = service.Response{
//			StatusCode: -1,
//			StatusMsg:  "Bad ActionType",
//		}
//		ctl.logger.Sugar().Errorf("Bad ActionType: %d", req.ActionType)
//		c.JSON(http.StatusBadRequest, rsp)
//		return
//	}
//}
//
//type CommentListRequest struct {
//	Token   string `form:"token" binding:"required"`
//	VideoID uint64 `form:"video_id" binding:"required"`
//}
//type CommentListResponse struct {
//	service.Response
//	CommentList []service.Comment `json:"comment_list"`
//}
//
//func (ctl *Controller) CommentList(c *gin.Context) {
//	var req CommentListRequest
//	var rsp CommentListResponse
//	err := c.ShouldBindQuery(&req)
//	if err != nil {
//		rsp.Response = service.Response{
//			StatusCode: -1,
//			StatusMsg:  err.Error(),
//		}
//		c.JSON(http.StatusBadRequest, rsp)
//	}
//	r, err := ctl.service.GetCommentList(req.VideoID)
//	if err != nil {
//		rsp.Response = service.Response{
//			StatusCode: -1,
//			StatusMsg:  err.Error(),
//		}
//		c.JSON(http.StatusInternalServerError, rsp)
//		return
//	}
//	rsp.Response = r.Response
//	rsp.CommentList = r.CommentList
//	c.JSON(http.StatusOK, rsp)
//}
