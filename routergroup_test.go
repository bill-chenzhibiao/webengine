package webengine

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestRouterGroupBasic(t *testing.T) {
	router := New()
	group := router.Group("/hola", func(c *Context) {})
	group.Use(func(c *Context) {})

	assert.Len(t, group.Handlers, 2)
	assert.Equal(t, "/hola", group.BasePath())
	assert.Equal(t, router, group.engine)

	group2 := group.Group("manu")
	group2.Use(func(c *Context) {}, func(c *Context) {})

	assert.Len(t, group2.Handlers, 4)
	assert.Equal(t, "/hola/manu", group2.BasePath())
	assert.Equal(t, router, group2.engine)
}


func TestRouterGroupBadMethod(t *testing.T) {
	router := New()
	assert.Panics(t, func() {
		router.Handle(http.MethodGet, "/")
	})
	assert.Panics(t, func() {
		router.Handle(" GET", "/")
	})
	assert.Panics(t, func() {
		router.Handle("GET ", "/")
	})
	assert.Panics(t, func() {
		router.Handle("", "/")
	})
	assert.Panics(t, func() {
		router.Handle("PO ST", "/")
	})
	assert.Panics(t, func() {
		router.Handle("1GET", "/")
	})
	assert.Panics(t, func() {
		router.Handle("PATCh", "/")
	})
}

func TestRouterGroupPipeline(t *testing.T) {
	router := New()
	testRoutesInterface(t, router)

	v1 := router.Group("/v1")
	testRoutesInterface(t, v1)
}

func testRoutesInterface(t *testing.T, r IRoutes) {
	handler := func(c *Context) {}
	assert.Equal(t, r, r.Use(handler))

	assert.Equal(t, r, r.Handle(http.MethodGet, "/handler", handler))
	assert.Equal(t, r, r.Any("/any", handler))
	assert.Equal(t, r, r.GET("/", handler))
	assert.Equal(t, r, r.POST("/", handler))
	assert.Equal(t, r, r.DELETE("/", handler))
	assert.Equal(t, r, r.PATCH("/", handler))
	assert.Equal(t, r, r.PUT("/", handler))
	assert.Equal(t, r, r.OPTIONS("/", handler))
	assert.Equal(t, r, r.HEAD("/", handler))
}

