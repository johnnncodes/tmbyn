package tmbyn

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/eknkc/amber"
)

func IndexHandler() func(http.ResponseWriter, *http.Request) {
	tf, err := sfs.Open("/index.amber")
	if err != nil {
		log.Fatal(err.Error())
	}
	tb, err := ioutil.ReadAll(tf)
	if err != nil {
		log.Fatal(err.Error())
	}
	t, err := amber.Compile(string(tb), amber.Options{false, false})
	if err != nil {
		log.Fatal(err.Error())
	}
	return func(w http.ResponseWriter, r *http.Request) {
		t.Execute(w, struct{}{})
	}
}
