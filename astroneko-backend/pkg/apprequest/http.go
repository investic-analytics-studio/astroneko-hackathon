package apprequest

import "github.com/valyala/fasthttp"

var (
	// POST Method
	POST = []byte(fasthttp.MethodPost)
	// GET Method
	GET = []byte(fasthttp.MethodGet)
	// PUT Method
	PUT = []byte(fasthttp.MethodPut)
	// PATCH Method
	PATCH = []byte(fasthttp.MethodPatch)
	// DELETE Method
	DELETE = []byte(fasthttp.MethodDelete)
	// ApplicationJSON header
	ApplicationJSON = []byte("application/json")
)

// HTTPRequest interface
type HTTPRequest interface {
	NewRequest(body []byte, method []byte, url string) (*fasthttp.Request, *fasthttp.Response)
	FastSetHeaderAuthorizationBearer(req *fasthttp.Request, token string)
}

// FastHTTP struct
type FastHTTP struct {
}

// NewRequester creates a new instance of fastHTTP
func NewRequester() *FastHTTP { return &FastHTTP{} }

// NewRequest creates a new request with the given body, method, and URL
func (u *FastHTTP) NewRequest(body []byte, method []byte, url string) (*fasthttp.Request, *fasthttp.Response) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	req.SetBody(body)
	req.Header.SetMethodBytes(method)
	req.SetRequestURIBytes([]byte(url))

	return req, resp
}

// FastSetHeaderAuthorizationBearer sets the Authorization header with the given token
func (u *FastHTTP) FastSetHeaderAuthorizationBearer(req *fasthttp.Request, token string) {
	req.Header.SetBytesKV([]byte("Authorization"), []byte("Bearer "+token))
}
