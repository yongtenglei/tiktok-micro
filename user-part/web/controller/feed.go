package controller

//
//import (
//	"context"
//	"net/http"
//	"strconv"
//	"time"
//
//	"github.com/gin-gonic/gin"
//	"github.com/yunyandz/tiktok-demo-backend/internal/service"
//)
//
//// Feed same demo video list for every request
//func (ctl *Controller) Feed(c *gin.Context) {
//	isnew := true
//	lasttimestring := c.Query("latest_time")
//	lasttimestamp := 0
//	var isTourist bool
//	userid := uint64(0)
//	uc, e := ctl.getUserClaims(c)
//	if !e {
//		isTourist = true
//	} else {
//		userid = uc.UserID
//	}
//	var err error
//	var lastestTime time.Time
//	if lasttimestring != "" {
//		lasttimestamp, err = strconv.Atoi(lasttimestring)
//		if err != nil {
//			ctl.logger.Sugar().Errorf("Atoi error: %v", err)
//			c.JSON(http.StatusBadRequest, gin.H{
//				"status": -1,
//				"msg":    err.Error(),
//			})
//			return
//		}
//		ctl.logger.Sugar().Debugf("lasttimestamp: %d", lasttimestamp)
//		lastestTime = time.UnixMilli(int64(lasttimestamp))
//		isnew = false
//	}
//
//	var rsp service.FeedResponse
//	r, err := ctl.service.GetFeed(context.Background(), userid, isnew, isTourist, lastestTime)
//	if err != nil {
//		rsp.Response = service.Response{
//			StatusCode: -1,
//			StatusMsg:  err.Error(),
//		}
//		c.JSON(http.StatusInternalServerError, rsp)
//		return
//	}
//	c.JSON(http.StatusOK, r)
//}
