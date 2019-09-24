package ztUtil

import "encoding/json"

func FromStructToString(source interface{}) string {
	bReq, _ := json.Marshal(source)
	return string(bReq)
}
