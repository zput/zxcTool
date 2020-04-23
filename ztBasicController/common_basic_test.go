package ztBasicController

import (
	"fmt"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHandlers(t *testing.T){
	//request := httptest.NewRequest("GET", "/hello_world", strings.NewReader("my name is ..."))
	response := httptest.NewRecorder()
	response.WriteString("nice to meet you")

	basic := new(Basic)
	basic.RequestBody = response.Body.Bytes()

	basic.SetHandler(func (this *Basic){
		if reflect.DeepEqual(this.RequestBody, basic.RequestBody)==false{
			t.Fatal(string(this.RequestBody))
		}
		fmt.Println("DONE")
	})
	basic.Finish()
}
