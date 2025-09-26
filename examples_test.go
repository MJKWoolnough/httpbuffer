package httpbuffer_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"vimagination.zapto.org/httpbuffer"
)

func handler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, World!")
}

func Example() {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	handler(w, r)

	fmt.Println(w.Result().ContentLength)

	w = httptest.NewRecorder()
	buf := httpbuffer.Handler{Handler: http.HandlerFunc(handler)}
	buf.ServeHTTP(w, r)

	fmt.Println(w.Result().ContentLength)
	fmt.Println(w.Body)

	// Output:
	// -1
	// 13
	// Hello, World!
}
