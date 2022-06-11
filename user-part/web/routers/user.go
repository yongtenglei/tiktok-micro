package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"yunyandz.com/tiktok/logger"
	"yunyandz.com/tiktok/user-part/web/controller"
	"yunyandz.com/tiktok/user-part/web/middlewares"
)

func SetUp() *gin.Engine {

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	//apiRouter.GET("/feed/", middleware.JWTAuth(logger.Logger(), false), ctl.Feed)

	// using jwt auth
	apiRouter.GET("/user/", middleware.JWTAuth(logger.Logger(), true), controller.UserInfo)

	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)

	//apiRouter.POST("/publish/action/", middleware.JWTAuth(logger.Logger(), true), ctl.Publish)
	//apiRouter.GET("/publish/list/", middleware.JWTAuth(logger.Logger(), true), ctl.PublishList)

	// extra apis - I
	//apiRouter.POST("/favorite/action/", middleware.JWTAuth(logger.Logger(), true), ctl.FavoriteAction)
	//apiRouter.GET("/favorite/list/", middleware.JWTAuth(logger.Logger(), true), ctl.FavoriteList)

	//apiRouter.POST("/comment/action/", middleware.JWTAuth(logger.Logger(), true), ctl.CommentAction)
	//apiRouter.GET("/comment/list/", ctl.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", middleware.JWTAuth(logger.Logger(), true), controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", middleware.JWTAuth(logger.Logger(), true), controller.FollowList)
	apiRouter.GET("/relation/follower/list/", controller.FollowerList)

	return r
}
