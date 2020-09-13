package utils

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Cmd(name string, arg ...string) string {
	cmd := exec.Command(name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	b, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return string(b)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

const letterBytes = "1234567890"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func generateTraceId(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

var TraceId = "TraceId"

func Trace(context *gin.Context) {
	traceId := generateTraceId(24)
	context.Set(TraceId, traceId)
	context.Writer.Header().Add(TraceId, traceId)
}

func LogInfo(ctx context.Context, content string) {
	log.Infof("[%s] %s", ctx.Value(TraceId), content)
}

func LogInfof(ctx context.Context, format string, args ...interface{}) {
	log.Infof("[%s] %s", ctx.Value(TraceId), fmt.Sprintf(format, args...))
}

func LogWarn(ctx context.Context, content string) {
	log.Warnf("[%s] %s", ctx.Value(TraceId), content)
}

func LogWarnf(ctx context.Context, format string, args ...interface{}) {
	log.Warnf("[%s] %s", ctx.Value(TraceId), fmt.Sprintf(format, args...))
}

func ReadFile(filePth string) []byte {
	content, err := ioutil.ReadFile(filePth)
	if err != nil {
		panic(err)
	}
	return content
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}
