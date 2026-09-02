package middleware

import (
	"bytes"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/royanqodri/Login-Gateway-API/util"
	"github.com/royanqodri/Login-Gateway-API/util/logging"
	"github.com/sirupsen/logrus"
)

type responseCapturer struct {
	gin.ResponseWriter
	BodyBuffer *bytes.Buffer
}

func (c *responseCapturer) Write(data []byte) (int, error) {
	c.BodyBuffer.Write(data)
	return c.ResponseWriter.Write(data)
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// init vars
		var requestBodyStr string
		var responseBodyStr string

		// set trace id
		traceId := util.GetTimeNowMillisInStr() + util.ToString(rand.Intn(20-1)+1)

		// start timer
		start := time.Now()

		// capture request body
		if ctx.Request.Body != nil {
			requestBody, err := ctx.GetRawData()
			if err != nil {
				util.WriteResponse(ctx, util.JSON, http.StatusInternalServerError, []any{}, []any{}, "request body conversion is error", 0, 0)
				ctx.Abort()
				return
			}
			requestBodyStr = util.CleanString(string(requestBody), "\n", "\t", "\\", "  ")

			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// capture response body
		capturingWriter := &responseCapturer{ResponseWriter: ctx.Writer, BodyBuffer: bytes.NewBuffer(nil)}
		ctx.Writer = capturingWriter

		// process request
		ctx.Set("trace_id", traceId)
		ctx.Next()

		// calculate latency
		latency := time.Since(start)

		// get response body from capturer
		responseBody := capturingWriter.BodyBuffer.Bytes()
		responseBodyStr = util.CleanString(string(responseBody), "\n", "\t", "\\", "  ")

		// log request details
		logging.LogWithFields(logging.HTTP_REQUEST, logging.INFO, logrus.Fields{
			"trace_id":      traceId,
			"status":        ctx.Writer.Status(),
			"method":        ctx.Request.Method,
			"path":          ctx.Request.URL.String(),
			"timestamp":     util.GetFormattedDateTimeMsInStr(start),
			"latency":       latency,
			"client_ip":     ctx.ClientIP(),
			"user_agent":    ctx.Request.UserAgent(),
			"request_body":  requestBodyStr,
			"response_body": responseBodyStr,
		})
	}
}
