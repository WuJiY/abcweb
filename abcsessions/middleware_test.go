package abcsessions

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewareWriteHeader(t *testing.T) {

	w := httptest.NewRecorder()
	cookies := &cookiesContext{}

	cookies.cookies = make(map[string]*http.Cookie)

	cookies.cookies["lol"] = &http.Cookie{Name: "lol", Value: "test1"}
	cookies.cookies["hehe"] = &http.Cookie{Name: "hehe", Value: "test2"}
	cookies.WriteHeader(w, 200)

	cookieResults := w.Result().Cookies()
	if len(cookieResults) != 2 {
		t.Error("expected cookies len 2, got:", len(cookieResults))
	}

	for _, c := range cookieResults {
		var found bool
		for _, rc := range cookies.cookies {
			if c.Name == rc.Name {
				found = true
			}
		}
		if !found {
			t.Errorf("could not find cookie with name %s in cookies", c.Name)
		}
	}
}

func TestMiddlewareSetCookie(t *testing.T) {

	cookies := &cookiesContext{}

	cookies.SetCookie(&http.Cookie{Name: "lolcats"})
	cookies.SetCookie(&http.Cookie{Name: "catlollers"})

	if len(cookies.cookies) != 2 {
		t.Error("expected len 2, got:", len(cookies.cookies))
	}

	if _, ok := cookies.cookies["lolcats"]; !ok {
		t.Error("expected lolcats to be set")
	}
	if _, ok := cookies.cookies["catlollers"]; !ok {
		t.Error("expected catlollers to be set")
	}
}

type middlewareOverseerMock struct {
	called bool
	resetExpiryMiddleware
}

func (m *middlewareOverseerMock) ResetExpiry(w http.ResponseWriter, r *http.Request) error {
	m.called = true
	return nil
}

func TestResetMiddleware(t *testing.T) {

	o := &middlewareOverseerMock{}
	o.resetExpiryMiddleware.resetter = o

	fn := func(w http.ResponseWriter, r *http.Request) {}

	hf := o.ResetMiddleware(http.HandlerFunc(fn))

	hf.ServeHTTP(nil, nil)
	if !o.called {
		t.Error("expected called true")
	}
}
