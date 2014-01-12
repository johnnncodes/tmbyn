package tmbyn

import (
	"github.com/eknkc/amber"
	"net/http"
)

func IndexHandler() func(http.ResponseWriter, *http.Request) {
	t, err := amber.CompileFile("index.amber", amber.Options{false, false})
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		t.Execute(w, struct{}{})
	}
}
