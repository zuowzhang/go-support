package router

import "net/http"

type MiddlewareFunc func(handlerFunc http.HandlerFunc) http.HandlerFunc

type App struct {
}
