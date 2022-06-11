package util

import (
	"github.com/gin-gonic/gin"
	"yunyandz.com/tiktok/logger"
	"yunyandz.com/tiktok/pkg/jwtx"
)

func GetUserClaims(c *gin.Context) (*jwtx.UserClaims, bool) {
	uc, e := c.Get("claims")
	if !e {
		logger.Suger().Errorf("Get claims error: %v", e)
		return nil, false
	}
	return uc.(*jwtx.UserClaims), true
}
