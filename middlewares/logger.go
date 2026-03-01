// package middlewares

// // importing necessary packages
// //  go get gitlab.com/pragmaticreviews/golang-gin-poc/middlewares
// import ( 
// 	"github.com/gin-gonic/gin"
// 	"fmt"
// 	"time"
// )

// func Logger() gin.HandlerFunc {
// 	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
// 		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \n",
// 		//  client IP, timestamp, HTTP method, path, protocol, status code, latency, user agent, error message
// 		//  clinet IP 
// 			param.ClientIP,
// 			param.TimeStamp.Format(time.RFC1123),
// 			param.Method,
// 			param.Path,
// 			param.Request.Proto,
// 			param.StatusCode,
// 			param.Latency,
// 		)
// 	})
// }