package webengine

import (
	"net/http"
)

type Engine struct{
	RouterGroup
	DispatchHandler *DispatchHandler
	methodPathMappings 	methodPathMappings
}

var _ IRouter = &Engine{}

func New() *Engine{
	engine := &Engine{
		RouterGroup:     RouterGroup{
			Handlers: nil,
			basePath: "/",
			root:     true,
		},
		methodPathMappings: make(methodPathMappings,0,9),
		DispatchHandler:    &DispatchHandler{},
	}
	engine.RouterGroup.engine = engine
	engine.DispatchHandler.MethodPathMappings = engine.methodPathMappings
	engine.DispatchHandler.InitPool()
	return engine
}

func Default() *Engine {
	engine := New()
	return engine
}

func (engine *Engine) Use(middleware ...HandlerFunc) IRoutes {
	engine.RouterGroup.Use(middleware...)
	return engine
}

func (engine *Engine) addRoute(method, path string, handlers HandlersChain) {
	assert1(path[0] == '/', "path must begin with '/'")
	assert1(method != "", "HTTP method can not be empty")
	assert1(len(handlers) > 0, "there must be at least one handler")

	pathMapping := engine.methodPathMappings.get(method)
	if pathMapping == nil{
		pathMapping = make(map[string]HandlersChain)
		engine.methodPathMappings = append(engine.methodPathMappings,methodPathMapping{method:method,pathMapping:pathMapping})
		engine.DispatchHandler.MethodPathMappings = engine.methodPathMappings
	}
	pathMapping[path] = handlers
}

func (engine *Engine) Run(addr ...string) (err error) {
	address := resolveAddress(addr)
	err = http.ListenAndServe(address, engine.DispatchHandler)
	return
}
