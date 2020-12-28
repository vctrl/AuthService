package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err *AppError) Error() string {
	return err.Message
}

//
// Middleware Error Handler in server package
//
func JSONAppErrorReporter() gin.HandlerFunc {
	return jsonAppErrorReporterT(gin.ErrorTypeAny)
}

func jsonAppErrorReporterT(errType gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		detectedErrors := c.Errors.ByType(errType)

		log.Println("Handle APP error")
		if len(detectedErrors) > 0 {
			err := detectedErrors[0].Err
			var parsedError *AppError
			switch err.(type) {
			case *AppError:
				parsedError = err.(*AppError)
			default:
				parsedError = &AppError{
					Code:    http.StatusInternalServerError,
					Message: "Internal server error",
				}
				log.Println(err.Error())
			}
			// Put the error into response
			c.IndentedJSON(parsedError.Code, parsedError)
			c.Abort()
			// or c.AbortWithStatusJSON(parsedError.Code, parsedError)
			return
		}

	}
}
