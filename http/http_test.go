package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Success")
	}))
	defer ts.Close()
	res, _ := Get(ts.URL, ts.Client())
	data, _ := ioutil.ReadAll(res.Body)

	if data == nil {
		t.Fail()
	}
}
