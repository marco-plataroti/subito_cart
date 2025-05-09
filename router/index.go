package router

import (
	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	RegisterOrderRoutes(r)
}
