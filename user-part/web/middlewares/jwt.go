package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"yunyandz.com/tiktok/pkg/constant"
	"yunyandz.com/tiktok/pkg/jwtx"
	"yunyandz.com/tiktok/user-part/web/common"

	"go.uber.org/zap"
)

func JWTAuth(logger *zap.Logger, strict bool) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 获取token
		token := c.Query("token")
		if token == "" || len(token) == 0 {
			token = c.PostForm("token")
			if token == "" || len(token) == 0 {
				if strict {
					rsp := common.Response{
						StatusCode: constant.CodeFail,
						StatusMsg:  "Invalid token",
					}
					logger.Sugar().Errorf("Invalid token: %v", rsp)
					c.JSON(http.StatusUnauthorized, rsp)
					c.Abort()
				} else {
					c.Next()
				}
				return
			}
		}

		// 解析token

		parsedToken, err := jwtx.ParseUserClaims(token)
		if err != nil {
			if strict {
				rsp := common.Response{
					StatusCode: constant.CodeFail,
					StatusMsg:  "Parse token failed",
				}

				//logger.Sugar().Errorf("Parse token failed: %v", rsp)
				c.JSON(http.StatusNonAuthoritativeInfo, rsp)
				c.Abort()
			} else {
				c.Next()
			}
			return
		}

		c.Set("claims", parsedToken)
		c.Next()
	}
}
