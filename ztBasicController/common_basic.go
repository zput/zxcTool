package ztBasicController

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/zput/zxcTool/ztUtil"
	"net/http"
	"strings"
	"time"
)

var DefaultCatchPanicfoo = func(){  // Handle this.Handle's panic(error)
	if e := recover(); e != nil {
		log.Error(e)
	}
}

type Basic struct {
	Ctx         *gin.Context
	RequestBody []byte

	url     string
	begin   int64
	XRealIp string
}

func (this *Basic) Prepare(ctxSource ...*gin.Context) {
	if len(ctxSource) != 0{
		this.Ctx = ctxSource[0]
	}

	var err error
	this.begin = time.Now().UnixNano()
	//this.XRealIp = this.Ctx.GetHeader("X-Client-IP")
	this.XRealIp = this.Ctx.ClientIP()

	this.url = strings.TrimSpace(this.Ctx.Request.URL.Path)

	log.Debugf("%s %s; enter", this.Ctx.Request.Method, this.url)

	if this.Ctx.Request.Body == nil {
		log.Tracef("NO body")
		return
	}
	this.RequestBody, err = this.Ctx.GetRawData()
	this.Handle(err)
	if len(this.RequestBody) > 0 {
		log.Tracef("body: %s", string(this.RequestBody))
	}
}

func (this *Basic) Finish() {
	elapse := (time.Now().UnixNano() - this.begin) / 1000000000
	log.Debugf("%s %s; go out, use %d msec", this.Ctx.Request.Method, this.url, elapse)
}

//successful
func (this *Basic) ReturnSuccJson200(data interface{}) {
	this.ReturnSuccJson(http.StatusOK, data)
}

func (this *Basic) ReturnSuccJson(httpStatus int, data interface{}) {
	if httpStatus != http.StatusNoContent {
		log.Tracef("Response: %+v", data)
		log.Tracef("Response JSON: %+v", ztUtil.FromStructToString(data))
	}
	this.Ctx.AbortWithStatusJSON(httpStatus, data)
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
	this.Ctx.AbortWithStatusJSON(httpStatus, data)
	// print
	this.Finish()
}

//failure
func (this *Basic) ReturnFailErrThroughError(err error) {
	httpStatus := http.StatusInternalServerError
	msg := ""
	if e, ok := err.(ztUtil.ErrorInterface); ok {
		log.Debugf("API Failed: %+v", e)
		httpStatus = e.Status()
		msg = e.Error()
	} else {
		log.Tracef("INTERNAL ERROR: %+v", err)
	}
	this.ReturnFailErr(httpStatus, msg)
}

func (this *Basic) Panic(err error) {
	httpStatus := http.StatusInternalServerError
	if e, ok := err.(ztUtil.ErrorInterface); ok {
		log.Debugf("API Failed: %+v", e)
		httpStatus = e.Status()
	} else {
		log.Tracef("INTERNAL ERROR: %+v", err)
	}
	this.Ctx.AbortWithStatus(httpStatus)
	// print
	this.Finish()
}

func (this *Basic) Handle(err error) {
	if err != nil {
		this.Panic(err)
		// TODO panic ?
		panic(err)
	}
}
