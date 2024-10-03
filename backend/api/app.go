package api

import (
	"CengkeHelper/logger"
	. "CengkeHelper/process"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

var validHosts = []string{
	"localhost", "cengkehelper.top", "huster.pages.dev", "ursb.top", "蹭课小助手.top",
}

// IsValidReqHosts 是否是合法请求主机
func IsValidReqHosts(target string) bool {
	// 切片必须升序
	sort.Strings(validHosts)
	index := sort.SearchStrings(validHosts, target)
	//index的取值：0 ~ (len(str_array)-1)
	return index < len(validHosts) && validHosts[index] == target
}

var HelperApp *gin.Engine

func init() {

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	//gin.DefaultWriter = writer

	HelperApp = gin.Default()
	fmt.Println()
	fmt.Println()
	fmt.Println()
	logger.Debug("实例启动！")
	logger.Info("\n\n" +
		"#    ____                          _  __         _   _          _                       \n" +
		"#   / ___|   ___   _ __     __ _  | |/ /   ___  | | | |   ___  | |  _ __     ___   _ __ \n" +
		"#  | |      / _ \\ | '_ \\   / _` | | ' /   / _ \\ | |_| |  / _ \\ | | | '_ \\   / _ \\ | '__|\n" +
		"#  | |___  |  __/ | | | | | (_| | | . \\  |  __/ |  _  | |  __/ | | | |_) | |  __/ | |   \n" +
		"#   \\____|  \\___| |_| |_|  \\__, | |_|\\_\\  \\___| |_| |_|  \\___| |_| | .__/   \\___| |_|   \n" +
		"#                          |___/                                   |_|                  ")

	// 中间件
	HelperApp.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://huster.pages.dev"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if origin != "http://localhost:5173" && origin != "https://huster.pages.dev" {
				logger.Warning("存在跨域访问: ", origin)
			}
			return origin == "http://localhost:5173" || origin == "https://huster.pages.dev"
		},
		MaxAge: 12 * time.Hour,
	}))

	// 前端界面

	HelperApp.Use(static.Serve("/", static.LocalFile("dist", true)))
	// 静态文件的文件夹相对目录
	HelperApp.StaticFS("/dist", http.Dir("./dist"))
	// 第一个参数是api 第二个参数是具体的文件名字
	HelperApp.StaticFile("/favicon.ico", "./favicon.ico")
	HelperApp.NoRoute(func(c *gin.Context) {
		accept := c.Request.Header.Get("Accept")
		flag := strings.Contains(accept, "text/html")
		if flag {
			content, err := os.ReadFile("dist/index.html")
			if (err) != nil {
				c.Writer.WriteHeader(404)
				_, err := c.Writer.WriteString("Not Found")
				if err != nil {
					logger.Error(err)
					return
				}
				return
			}
			c.Writer.WriteHeader(200)
			c.Writer.Header().Add("Accept", "text/html")
			_, err = c.Writer.Write(content)
			if err != nil {
				return
			}
			c.Writer.Flush()
		}
	})

	HelperApp.POST("/teach-infos", func(c *gin.Context) {
		if !IsValidReqHosts(c.Request.Host) {
			logger.WarningF("请求的主机不合法: %v, referer为: %v", c.Request.Host, c.Request.Referer())
			c.JSON(400, gin.H{"error": "Bad Request: Invalid Host"})
			return
		}

		// Assuming GetTeachInfos(true) returns the necessary information.
		data := GetTeachInfos(true)
		c.JSON(200, data)
	})

	HelperApp.GET("/cur-time", func(c *gin.Context) {
		// Log the request for debugging purposes
		logger.InfoF("IP => %v accessing /cur-time", c.ClientIP())

		// Call to function that calculates current time details
		weekNum, weekday, lessonNum := CurCourseTime()

		// Validate the cache status
		valid := ValidCache()

		if !valid {
			logger.Warning("Cache is invalid, refreshing...")
			GetTeachInfos(false) // refresh data
			FreshCacheFlag()     // reset cache flag
		}

		// Respond with the current course time information
		c.JSON(200, gin.H{
			"isAdjust":  ChoseSpecialDate.Chose,
			"weekNum":   weekNum,
			"weekday":   weekday,
			"lessonNum": lessonNum,
			"valid":     valid,
		})
	})

	HelperApp.GET("/ping", func(c *gin.Context) {
		logger.Warning("尝试更改时间！")

		courseDate := CourseDate{}
		if err := c.ShouldBindJSON(&courseDate); err != nil {
			c.JSON(400, gin.H{
				"error": "非法参数！",
			})
			logger.Warning("非法参数！")
			return
		}
		*ChoseSpecialDate = courseDate
		c.JSON(200, gin.H{
			"res": cast.ToString(ChoseSpecialDate.Chose),
		})
	})

}

func Listen(w http.ResponseWriter, r *http.Request) {
	HelperApp.ServeHTTP(w, r)
}
