package server

import (
	"fmt"
	"net/http"
)

type MiscController struct {
	// No dependencies necessary here
}

func (c *MiscController) HandleEmptyURI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Ymir is running!")
}

func (c *MiscController) RegisterRoutes(r muxRouter) {
	r.HandleFunc("/", c.HandleEmptyURI)
}
