package webengine

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReset(t *testing.T){
	c := &Context{
		Params:make(Params,2),
		handlers:make(HandlersChain,5),
		bodyArray:make([]map[string]interface{},5),
	}
	c.reset()
	assert.Equal(t,0,len(c.Params))
	assert.Equal(t,true,nil == c.handlers)
	assert.Equal(t,0,len(c.bodyArray))
}

func TestDispatchHandler_ServeHTTP(t *testing.T) {
	flag := true
	r := Default()
	testhandler := func(c *Context) {
		flag = false
		assert.Equal(t,"name",c.Params[0].Key)
		assert.Equal(t,"zhangsan",c.Params[0].Value)

		object := c.GetJsonObject()
		assert.Equal(t,(*object)["data"],float64(-1))

		assert.Equal(t,c.GetHeaderValueByName("Content-Type"),"application/json")
	}
	r.POST("/test",testhandler)

	requestBody := make(map[string]interface{})
	requestBody["data"] = -1
	jsons,_ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/test?name=zhangsan", strings.NewReader(string(jsons)))
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.DispatchHandler.ServeHTTP(w, req)

	assert.Equal(t,false,flag)
}