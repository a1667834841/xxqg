// Package web
// @Description: 封装了所以web相关的内容
package web

import (
	"embed"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/johlanse/study_xxqg/conf"
	"github.com/johlanse/study_xxqg/utils"
	"github.com/johlanse/study_xxqg/utils/update"
)

// 将静态文件嵌入到可执行程序中来
//go:embed xxqg/build
var static embed.FS

// RouterInit
// @Description:
// @return *gin.Engine
func RouterInit() *gin.Engine {
	router := gin.Default()
	router.RemoveExtraSlash = true
	router.Use(cors())

	// 挂载静态文件
	router.StaticFS("/xxqg/static", http.FS(static))
	// 访问首页时跳转到对应页面
	router.GET("/xxqg", func(ctx *gin.Context) {
		ctx.Redirect(301, "/xxqg/static/xxqg/build/home.html")
	})

	router.GET("/xxqg/about", func(context *gin.Context) {
		context.JSON(200, Resp{
			Code:    200,
			Message: "",
			Data:    utils.GetAbout(),
			Success: true,
			Error:   "",
		})
	})

	router.POST("/xxqg/restart", check(), func(ctx *gin.Context) {
		if ctx.GetInt("level") == 1 {
			ctx.JSON(200, Resp{
				Code:    200,
				Message: "",
				Data:    nil,
				Success: true,
				Error:   "",
			})
			utils.Restart()
		} else {
			ctx.JSON(200, Resp{
				Code:    401,
				Message: "",
				Data:    nil,
				Success: false,
				Error:   "",
			})
		}
	})

	router.POST("/xxqg/update", check(), func(ctx *gin.Context) {
		if ctx.GetInt("level") == 1 {
			update.SelfUpdate("", conf.GetVersion())
			ctx.JSON(200, Resp{
				Code:    200,
				Message: "",
				Data:    nil,
				Success: true,
				Error:   "",
			})
			utils.Restart()
		} else {
			ctx.JSON(200, Resp{
				Code:    401,
				Message: "",
				Data:    nil,
				Success: false,
				Error:   "",
			})
		}
	})

	if utils.FileIsExist("./config/flutter_xxqg/") {
		router.StaticFS("/xxqg/flutter_xxqg", http.Dir("./config/flutter_xxqg/"))
	}
	// 对权限的管理组
	auth := router.Group("/xxqg/auth")
	// 用户登录的接口
	auth.POST("/login", userLogin())
	// 检查登录状态的token是否正确
	auth.POST("/check/:token", checkToken())

	// 对于用户可自定义挂载文件的目录
	if utils.FileIsExist("./config/dist/") {
		router.StaticFS("/xxqg/dist", http.Dir("./config/dist/"))
	}

	config := router.Group("/xxqg/config", check())

	config.GET("", configGet())
	config.POST("", configSet())

	// 对用户管理的组
	user := router.Group("/xxqg/user", check())
	// 添加用户
	user.POST("", addUser())
	// 获取所以已登陆的用户
	user.GET("", getUsers())
	// 删除用户
	user.DELETE("", deleteUser())

	// 获取用户成绩
	router.GET("/xxqg/score", getScore())
	// 让一个用户开始学习
	router.POST("/xxqg/study", study())
	// 让一个用户停止学习
	router.POST("/xxqg/stop_study", check(), stopStudy())
	// 获取程序当天的运行日志
	router.GET("/xxqg/log", check(), getLog())

	// 登录xxqg的三个接口
	router.GET("/xxqg/sign/", sign())
	router.GET("/xxqg/login/*proxyPath", generate())
	router.POST("/xxqg/login/*proxyPath", check(), generate())
	return router
}

func check() gin.HandlerFunc {
	config := conf.GetConfig()
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		token = strings.Split(token, " ")[1]
		if token == "" {
			ctx.JSON(401, Resp{
				Code:    401,
				Message: "the auth fail",
				Data:    nil,
				Success: false,
				Error:   "",
			})
			ctx.Abort()
		} else if utils.StrMd5(config.Web.Account+config.Web.Password) == token {
			ctx.Set("level", 1)
			ctx.Set("token", token)
			ctx.Next()
		} else if checkCommonUser(token) {
			ctx.Set("level", 2)
			ctx.Set("token", token)
			ctx.Next()
		} else {
			ctx.JSON(401, Resp{
				Code:    401,
				Message: "the auth fail",
				Data:    nil,
				Success: false,
				Error:   "",
			})
			ctx.Abort()
		}
	}
}

func checkCommonUser(token string) bool {
	config := conf.GetConfig()
	for key, value := range config.Web.CommonUser {
		if token == utils.StrMd5(key+value) {
			return true
		}
	}
	return false
}
