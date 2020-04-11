package webengine

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type Param struct {
	Key string
	Value string
}
type Params []Param

type DispatchHandler struct{
	pool sync.Pool
	Params	Params
	MethodPathMappings 	methodPathMappings
}

type HandlerFunc func(*Context)
type HandlersChain []HandlerFunc

type Context struct{
	Request	*http.Request
	Writer http.ResponseWriter
	Params Params
	handlers HandlersChain
	bodyArray []map[string]interface{}
}

// ServeHTTP conforms to the http.Handler interface.
func (handler *DispatchHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := handler.pool.Get().(*Context)
	c.Writer = w
	c.Request = req
	c.reset()
	handler.handleHTTPRequest(c)

	handler.pool.Put(c)
}

func (handler *DispatchHandler) handleHTTPRequest(c *Context) {
	httpMethod := c.Request.Method
	rPath := c.Request.URL.Path
	pathMapping := handler.MethodPathMappings.get(httpMethod)

	index := strings.Index(rPath, "?")
	if index <= 0{
		index = len(rPath)
	}
	c.setParams()

	handlersChain := pathMapping[rPath[:index]]

	for _, handlerFunc := range handlersChain{
		handlerFunc(c)
	}
}

func (handler *DispatchHandler) InitPool() {
	handler.pool.New = func() interface{}{
		return handler.allocateContext()
	}
}

func (handler *DispatchHandler) allocateContext() *Context {
	return &Context{}
}

func (c *Context) reset() {
	c.Params = c.Params[0:0]
	c.handlers = nil
	c.bodyArray = c.bodyArray[0:0]
}

// Status sets the HTTP response code.
func (c *Context) Status(code int) {
	c.Writer.WriteHeader(code)
}

func (c *Context) GetJsonObject() *map[string]interface{}{
	if len(c.bodyArray) == 0{
		var jsonObject map[string]interface{}
		body,_ := ioutil.ReadAll(c.Request.Body)
		json.Unmarshal(body,&jsonObject)
		c.bodyArray = append(c.bodyArray,jsonObject)
	}
	return &c.bodyArray[0]
}

func (c *Context) GetJsonArray() *[]map[string]interface{}{
	if len(c.bodyArray) == 0{
		var jsonArray []map[string]interface{}
		body,_ := ioutil.ReadAll(c.Request.Body)
		json.Unmarshal(body,&jsonArray)
		c.bodyArray = jsonArray
	}
	return &c.bodyArray
}

func (c *Context) JSON(code int, obj interface{}) {
	c.Status(code)
	c.Writer.Header().Set("Content-Type","application/json;charset=utf-8")
	json.NewEncoder(c.Writer).Encode(obj)
}

func (c *Context) setParams() {
	vars := c.Request.URL.Query();
	for k, arr := range vars{
		for _,v := range arr{
			c.Params = append(c.Params,Param{
				Key:   k,
				Value: v,
			})
		}
	}
}

func (c *Context) GetHeaderValueByName(name string) string{
	return c.Request.Header.Get(name)
}

func (c *Context) String(code int, s string) {
	c.JSON(code,s)
}

func (ps Params) Get(name string)(string,bool){
	for _,entry := range ps{
		if entry.Key == name{
			return entry.Value,true
		}
	}
	return "",false
}

func (ps Params) ByName(name string)(va string){
	va, _ = ps.Get(name)
	return
}

