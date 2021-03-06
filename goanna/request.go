package goanna

import (
	"io/ioutil"
	"net"
	"net/http"

	ghttp "github.com/99designs/goodies/http"
	"github.com/gorilla/mux"
)

// Request decorates a http.Request to add helper methods
type Request struct {
	*http.Request
	session  Session
	bodyData []byte
	bodyRead bool
}

// CookieValue returns the value of the named cookie
func (r *Request) CookieValue(name string) string {
	c, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	return c.Value
}

// BodyData returns the full request body
func (r *Request) BodyData() []byte {
	var err error
	if !r.bodyRead {
		if r.Body != nil {
			r.bodyData, err = ioutil.ReadAll(r.Body)
			if err != nil {
				// catch i/o timeout errors
				neterr, isNetError := err.(net.Error)
				if isNetError && neterr.Timeout() {
					panic(ghttp.NewHttpError(err, http.StatusRequestTimeout))
				} else {
					panic(err)
				}
			}
		}
		r.bodyRead = true
	}

	return r.bodyData
}

// QueryValue returns the value in the GET query string
func (r *Request) QueryValue(key string) string {
	return r.URL.Query().Get(key)
}

func (r *Request) QueryValueOrDefault(key string, def string) string {
	val := r.URL.Query().Get(key)
	if val == "" {
		val = def
	}

	return val
}

// FormValueOrDefault returns the result of Request.FormValue,
// and if the result is empty returns the default string
func (r *Request) FormValueOrDefault(key string, def string) string {
	val := r.FormValue(key)
	if val == "" {
		val = def
	}

	return val
}

// IsGet returns whether the request is GET
func (r *Request) IsGet() bool {
	return r.Method == "GET"
}

// IsPost returns whether the request is POST
func (r *Request) IsPost() bool {
	return r.Method == "POST"
}

// IsPost returns whether the request is HEAD
func (r *Request) IsHead() bool {
	return r.Method == "HEAD"
}

// IsPut returns whether the request is PUT
func (r *Request) IsPut() bool {
	return r.Method == "PUT"
}

// IsPatch returns whether the request is PATCH
func (r *Request) IsPatch() bool {
	return r.Method == "PATCH"
}

func (r *Request) Log(v ...string) {
	LogRequest(r, v...)
}

// UrlValue returns whether the request is PATCH
func (r *Request) UrlValue(key string) string {
	return mux.Vars(r.Request)[key]
}
