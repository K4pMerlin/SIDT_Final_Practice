package api

import (
	"CengkeHelper/Backend/process"
	"CengkeHelper/logger"
	"net/http"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
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
		//AllowOrigins: []string{"https://foo.com"},
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{"*"},
		//ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if origin != "https://huster.pages.dev" {
				logger.Warning("存在跨域访问: ", origin)
			}
			return true
			//return origin == "http://localhost"
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
			logger.WarningF("请求的主机不合法: %v, refer为: %v", c.Request.Host, c.Request.Referer())
			c.JSON(400, "bad request")
		} else {
			//logger.Warning(GetTeachInfos(true))
			c.JSON(200, GetTeachInfos(true))
		}

	})

	HelperApp.GET("/cur-time", func(c *gin.Context) {
		//logger.InfoF("ip => %v 「%v」 ==> website", c.ClientIP(),
		//	location.IpToLocation(c.ClientIP()))
		weekNum, weekday, lessonNum := CurCourseTime()

		valid := true
		if !ValidCache() {
			logger.Warning("缓存失效")
			GetTeachInfos(false)
			FreshCacheFlag()
			valid = false
		}

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

	
	HelperApp.GET("/search", func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			c.JSON(400, gin.H{"error": "Query parameter 'q' is required"})
			return
		}

		results := process.SearchCourses(query)
		c.JSON(200, results)
	})

}

func Listen(w http.ResponseWriter, r *http.Request) {
	HelperApp.ServeHTTP(w, r)
}
