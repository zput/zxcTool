package zxcToolBasicController

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
	"github.com/zput/zxcTool/zxcUtil"
)

type Basic struct{
	ctx *gin.Context
	RequestBody []byte

	url     string
	begin   int64
	XRealIp string
}

func (this *Basic) Prepare(){
	var err error
	this.begin = time.Now().UnixNano()
	//this.XRealIp = this.ctx.GetHeader("X-Client-IP")
	this.XRealIp = this.ctx.ClientIP()

	this.url = strings.TrimSpace(this.ctx.Request.URL.Path)

	log.Debugf("%s %s; enter", this.ctx.Request.Method, this.url)

	this.RequestBody, err = this.ctx.GetRawData()
	this.handle(err)
	if len(this.RequestBody) > 0 {
		log.Tracef("body: %s", string(this.RequestBody))
	}
}

func (this *Basic) Finish() {
	elapse := (time.Now().UnixNano() - this.begin) / 1000000000
	log.Debugf("%s %s; go out, use %d msec", this.ctx.Request.Method, this.url, elapse)
}

//successful
func (this *Basic) ReturnSuccJson200(data interface{}) {
	this.ReturnSuccJson(http.StatusOK, data)
}

func (this *Basic) ReturnSuccJson(httpStatus int, data interface{}) {
	if httpStatus != http.StatusNoContent {
		log.Tracef("Response: %+v", data)
		log.Tracef("Response JSON: %+v", zxcUtil.FromStructToString(data))
	}
	this.ctx.AbortWithStatusJSON(httpStatus, data)
	// print
	this.Finish()
}

//failure
func (this *Basic) ReturnFailErr(httpStatus int, msg ...string) {
	tempMessage := http.StatusText(httpStatus)
	if len(msg) != 0 {
		tempMessage = msg[0]
	}
	data := make(map[string]string)
	data["message"] = tempMessage
	this.ctx.AbortWithStatusJSON(httpStatus, data)
	// print
	this.Finish()
}

func (this *Basic) panic(err error) {
	httpStatus := http.StatusInternalServerError
	if e, ok := err.(ErrorInterface); ok {
		log.Debugf("API Failed: %+v", e)
		httpStatus = e.Status()
	} else {
		log.Tracef("INTERNAL ERROR: %+v", err)
	}
	this.ctx.AbortWithStatus(httpStatus)
	// print
	this.Finish()
}

func (this *Basic) handle(err error) {
	if err != nil {
		this.panic(err)
		// TODO panic ?
		panic(err)
	}
}

// interface
type ErrorInterface interface {
	Status() int
}

//struct
type Error struct {
	status  int    `json:"-"`
	Code    *int   `json:"code,omitempty"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func (e Error) Status() int {
	return e.status
}

//new a above error
func CodedErrorf(status int, code *int, format string, args ...interface{}) error {
	return &Error{status: status, Code: code, Message: fmt.Sprintf(format, args...)}
}

func Errorf(status int, format string, args ...interface{}) error {
	return &Error{status: status, Message: fmt.Sprintf(format, args...)}
}
