package webengine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddRoute(t *testing.T){
	engine := Default()
	method := "POST"
	pathA := "/data/summary/visit/set"
	var handlerA HandlerFunc = func(context *Context) {
	}

	var handlerChainA HandlersChain
	handlerChainA = append(handlerChainA,handlerA)
	assert.Equal(t,0,len(engine.DispatchHandler.MethodPathMappings),"should be 0")
	engine.addRoute(method,pathA,handlerChainA)
	assert.Equal(t,1,len(engine.DispatchHandler.MethodPathMappings),"set handlerChain error")
}