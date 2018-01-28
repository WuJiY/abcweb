package abcsessions

import "net/http"

type cookieWriter interface {
	SetCookie(cookie *http.Cookie)
	GetCookie(name string) *http.Cookie
}

// sessionsResponseWriter is a wrapper of the ResponseWriter object used to
// buffer the cookies across session API calls, so that they can be written
// at the very end of the response workflow (opposed to written on every operation)
type cookiesContext struct {
	wroteHeader  bool
	wroteCookies bool
	cookies      map[string]*http.Cookie
}

// WriteHeader sets all cookies in the buffer on the underlying ResponseWriter's
// headers and calls the underlying ResponseWriter WriteHeader func
func (c *cookiesContext) WriteHeader(w http.ResponseWriter, code int) {
	c.wroteHeader = true

	// Set all the cookies in the cookie buffer
	if !c.wroteCookies {
		c.wroteCookies = true
		for _, c := range c.cookies {
			http.SetCookie(w, c)
		}
	}

	w.WriteHeader(code)
}

func (c *cookiesContext) SetCookie(cookie *http.Cookie) {
	if c.cookies == nil {
		c.cookies = make(map[string]*http.Cookie)
	}

	if len(cookie.Name) == 0 {
		panic("cookie name cannot be empty")
	}

	c.cookies[cookie.Name] = cookie
}

func (c *cookiesContext) GetCookie(name string) *http.Cookie {
	return c.cookies[name]
}

type resetExpiryMiddleware struct {
	resetter Resetter
}

// ResetMiddleware resets the users session expiry on each request.
//
// Note: Generally you will want to use Middleware or MiddlewareWithReset instead
// of this middleware, but if you have a requirement to insert a middleware
// between the sessions Middleware and the sessions ResetMiddleware then you can
// use the Middleware first, and this one second, as a two-step process, instead
// of the combined MiddlewareWithReset middleware.
//
// It's also important to note that the sessions Middleware must come BEFORE
// this middleware in the chain, or you will get a panic.
func (m resetExpiryMiddleware) ResetMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := m.resetter.ResetExpiry(w, r)
		// It's possible that the session hasn't been created yet
		// so there's nothing to reset. In that case, do not explode.
		if err != nil && !IsNoSessionError(err) {
			panic(err)
		}

		next.ServeHTTP(w, r)
	})
}
