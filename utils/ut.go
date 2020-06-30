package utils

import (
	"bytes"
	"io/ioutil"
	"encoding/json"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)
var r *gin.Engine


func UTInit(){
	r = gin.Default()
}

func Router() *gin.Engine {
	return r
}

// ParseToStr :parse the JSON params to be HTTP query postfix format
func ParseToStr(mp map[string]string) string {
	values := ""
	for key, val := range mp {
		   values += "&" + key + "=" + val
	}
	temp := values[1:]
	values = "?" + temp
	return values
}

// PostJSON :simulate a HTTP POST req with JSON params, return the rsp body and code
func PostJSON(uri string, param map[string]interface{}, router *gin.Engine) ([]byte, int) {
	jsonByte,_ := json.Marshal(param)
	req := httptest.NewRequest("POST", uri, bytes.NewReader(jsonByte))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	body,_ := ioutil.ReadAll(result.Body)
	return body, result.StatusCode
}

// Delete :simulate a HTTP DELETE req, return the rsp body and code
func Delete(uri string, router *gin.Engine) ([]byte, int) {
	req := httptest.NewRequest("DELETE", uri, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	body,_ := ioutil.ReadAll(result.Body)   
	return body, result.StatusCode
}

// PutJSON :simulate a HTTP PUT req with JSON params, return the rsp body and code
func PutJSON(uri string, param map[string]interface{}, router *gin.Engine) ([]byte, int) {
	jsonByte,_ := json.Marshal(param)
	req := httptest.NewRequest("PUT", uri, bytes.NewReader(jsonByte))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	body,_ := ioutil.ReadAll(result.Body)
	return body, result.StatusCode
}

// Get :simulate a HTTP GET req, return the rsp body and code
func Get(uri string, router *gin.Engine) ([]byte, int) {
	req := httptest.NewRequest("GET", uri, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	body,_ := ioutil.ReadAll(result.Body)   
	return body, result.StatusCode
}