package router

import (
	"net/http"
)

type Filter func(http.HandlerFunc) http.HandlerFunc
