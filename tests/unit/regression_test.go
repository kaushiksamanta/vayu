package unit

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/kaushiksamanta/vayu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Fix #1: Concurrent handler-chain race ---

func TestConcurrentHandlerChainRace(t *testing.T) {
	app := vayu.New()

	app.Use(func(c *vayu.Context, next vayu.NextFunc) {
		next()
	})

	app.GET("/a", func(c *vayu.Context, next vayu.NextFunc) {
		_, _ = c.Send(vayu.StatusOK, "a")
	})

	app.GET("/b", func(c *vayu.Context, next vayu.NextFunc) {
		_, _ = c.Send(vayu.StatusOK, "b")
	})

	var wg sync.WaitGroup
	const concurrency = 100

	errorsA := make([]string, 0)
	errorsB := make([]string, 0)
	var muA, muB sync.Mutex

	for i := 0; i < concurrency; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			req := httptest.NewRequest("GET", "/a", nil)
			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)
			if w.Body.String() != "a" {
				muA.Lock()
				errorsA = append(errorsA, w.Body.String())
				muA.Unlock()
			}
		}()
		go func() {
			defer wg.Done()
			req := httptest.NewRequest("GET", "/b", nil)
			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)
			if w.Body.String() != "b" {
				muB.Lock()
				errorsB = append(errorsB, w.Body.String())
				muB.Unlock()
			}
		}()
	}

	wg.Wait()
	assert.Empty(t, errorsA, "route /a returned wrong body in %d requests", len(errorsA))
	assert.Empty(t, errorsB, "route /b returned wrong body in %d requests", len(errorsB))
}

// --- Fix #2: Harden prefix matching ---

func TestGroupPrefixDoesNotMatchSimilarPaths(t *testing.T) {
	app := vayu.New()

	var apiMWCalled bool
	api := app.Group("/api")
	api.Use(func(c *vayu.Context, next vayu.NextFunc) {
		apiMWCalled = true
		next()
	})
	api.GET("/users", func(c *vayu.Context, next vayu.NextFunc) {
		_, _ = c.Send(vayu.StatusOK, "api users")
	})

	app.GET("/api2/data", func(c *vayu.Context, next vayu.NextFunc) {
		_, _ = c.Send(vayu.StatusOK, "api2 data")
	})

	// Request to /api2/data should NOT trigger /api group middleware
	apiMWCalled = false
	req := httptest.NewRequest("GET", "/api2/data", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	assert.False(t, apiMWCalled, "/api group middleware should not run for /api2/data")
	assert.Equal(t, "api2 data", w.Body.String())

	// Request to /api/users SHOULD trigger /api group middleware
	apiMWCalled = false
	req2 := httptest.NewRequest("GET", "/api/users", nil)
	w2 := httptest.NewRecorder()
	app.ServeHTTP(w2, req2)

	assert.True(t, apiMWCalled, "/api group middleware should run for /api/users")
	assert.Equal(t, "api users", w2.Body.String())
}

func TestStaticPrefixDoesNotMatchSimilarPaths(t *testing.T) {
	app := vayu.New()

	// Serve static from a temp dir (we just need to verify the middleware doesn't match wrong paths)
	app.Static("/assets", ".")

	app.GET("/assets-old/file", func(c *vayu.Context, next vayu.NextFunc) {
		_, _ = c.Send(vayu.StatusOK, "assets-old file")
	})

	// /assets-old/file should NOT be handled by the /assets static handler
	req := httptest.NewRequest("GET", "/assets-old/file", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	assert.Equal(t, vayu.StatusOK, w.Code)
	assert.Equal(t, "assets-old file", w.Body.String())
}

// --- Fix #3: Prevent panics in generic query binding ---

func TestBindQueryParamsNonStructType(t *testing.T) {
	app := vayu.New()

	app.GET("/non-struct", func(c *vayu.Context, next vayu.NextFunc) {
		_, err := vayu.BindQueryParams[string](c)
		if err != nil {
			_, _ = c.Send(vayu.StatusBadRequest, err.Error())
			return
		}
		_, _ = c.Send(vayu.StatusOK, "ok")
	})

	req, _ := http.NewRequest("GET", "/non-struct?q=test", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	assert.Equal(t, vayu.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "struct")
}

func TestBindQueryParamsInt8Overflow(t *testing.T) {
	app := vayu.New()

	app.GET("/overflow", func(c *vayu.Context, next vayu.NextFunc) {
		type Params struct {
			Small int8 `query:"val"`
		}
		_, err := vayu.BindQueryParams[Params](c)
		if err != nil {
			_, _ = c.Send(vayu.StatusBadRequest, err.Error())
			return
		}
		_, _ = c.Send(vayu.StatusOK, "ok")
	})

	// int8 max is 127; 999 should overflow
	req, _ := http.NewRequest("GET", "/overflow?val=999", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	assert.Equal(t, vayu.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "int8")
}

func TestBindQueryParamsUint8Overflow(t *testing.T) {
	app := vayu.New()

	app.GET("/uint-overflow", func(c *vayu.Context, next vayu.NextFunc) {
		type Params struct {
			Small uint8 `query:"val"`
		}
		_, err := vayu.BindQueryParams[Params](c)
		if err != nil {
			_, _ = c.Send(vayu.StatusBadRequest, err.Error())
			return
		}
		_, _ = c.Send(vayu.StatusOK, "ok")
	})

	// uint8 max is 255; 999 should overflow
	req, _ := http.NewRequest("GET", "/uint-overflow?val=999", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	assert.Equal(t, vayu.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "uint8")
}

// --- Fix #4: Simplify context timeout API ---

func TestWithTimeoutReturnsCancelFunc(t *testing.T) {
	app := vayu.New()

	app.GET("/timeout", func(c *vayu.Context, next vayu.NextFunc) {
		cancel := c.WithTimeout(5 * time.Second)
		defer cancel()

		// Context should not be done yet
		select {
		case <-c.Ctx.Done():
			t.Error("context should not be done yet")
		default:
		}

		// Cancel and verify context is done
		cancel()
		assert.Error(t, c.Ctx.Err())
		assert.Equal(t, context.Canceled, c.Ctx.Err())

		_, _ = c.Send(vayu.StatusOK, "ok")
	})

	req := httptest.NewRequest("GET", "/timeout", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	assert.Equal(t, vayu.StatusOK, w.Code)
}

func TestWithTimeoutExpires(t *testing.T) {
	app := vayu.New()

	app.GET("/timeout-expire", func(c *vayu.Context, next vayu.NextFunc) {
		cancel := c.WithTimeout(10 * time.Millisecond)
		defer cancel()

		// Wait for timeout
		<-c.Ctx.Done()
		assert.ErrorIs(t, c.Ctx.Err(), context.DeadlineExceeded)

		_, _ = c.Send(vayu.StatusOK, "expired")
	})

	req := httptest.NewRequest("GET", "/timeout-expire", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	assert.Equal(t, vayu.StatusOK, w.Code)
}

// --- Fix #5: Guard recovery writes after partial response ---

func TestRecoveryAfterPartialWrite(t *testing.T) {
	app := vayu.New()
	app.Use(vayu.Recovery())

	app.GET("/partial-panic", func(c *vayu.Context, next vayu.NextFunc) {
		// Write partial response first
		_, _ = c.Send(vayu.StatusOK, "partial")
		// Then panic
		panic("after partial write")
	})

	req := httptest.NewRequest("GET", "/partial-panic", nil)
	w := httptest.NewRecorder()

	// Should not panic
	assert.NotPanics(t, func() {
		app.ServeHTTP(w, req)
	})

	// The original partial response should be intact
	assert.Equal(t, vayu.StatusOK, w.Code)
	assert.Equal(t, "partial", w.Body.String())
}

func TestRecoveryWithoutPriorWrite(t *testing.T) {
	app := vayu.New()
	app.Use(vayu.Recovery())

	app.GET("/clean-panic", func(c *vayu.Context, next vayu.NextFunc) {
		panic("clean panic")
	})

	req := httptest.NewRequest("GET", "/clean-panic", nil)
	w := httptest.NewRecorder()

	assert.NotPanics(t, func() {
		app.ServeHTTP(w, req)
	})

	// Should get the 500 JSON error response
	assert.Equal(t, 500, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Internal Server Error", response["error"])
}

// --- Fix #6: ResponseWriter compatibility ---

func TestResponseWriterImplicit200(t *testing.T) {
	app := vayu.New()

	app.GET("/implicit", func(c *vayu.Context, next vayu.NextFunc) {
		// Write body without calling WriteHeader
		_, _ = c.Writer.Write([]byte("hello"))
	})

	req := httptest.NewRequest("GET", "/implicit", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "hello", w.Body.String())
}

func TestResponseWriterStatusTracking(t *testing.T) {
	rw := vayu.NewResponseWriter(httptest.NewRecorder())

	// Before any write, status should be 0
	assert.Equal(t, 0, rw.Status())
	assert.False(t, rw.Written())

	// Write without WriteHeader
	_, _ = rw.Write([]byte("data"))

	// Should now track implicit 200
	assert.Equal(t, http.StatusOK, rw.Status())
	assert.True(t, rw.Written())
}

func TestResponseWriterDoubleWriteHeaderIgnored(t *testing.T) {
	recorder := httptest.NewRecorder()
	rw := vayu.NewResponseWriter(recorder)

	rw.WriteHeader(http.StatusCreated)
	assert.Equal(t, http.StatusCreated, rw.Status())

	// Second WriteHeader should be ignored
	rw.WriteHeader(http.StatusNotFound)
	assert.Equal(t, http.StatusCreated, rw.Status())
}

func TestResponseWriterFlushInterface(t *testing.T) {
	recorder := httptest.NewRecorder()
	rw := vayu.NewResponseWriter(recorder)

	// httptest.ResponseRecorder implements http.Flusher
	assert.NotPanics(t, func() {
		rw.Flush()
	})

	assert.True(t, recorder.Flushed)
}

func TestResponseWriterHijackUnsupported(t *testing.T) {
	recorder := httptest.NewRecorder()
	rw := vayu.NewResponseWriter(recorder)

	// httptest.ResponseRecorder does NOT implement http.Hijacker
	conn, _, err := rw.Hijack()
	assert.Nil(t, conn)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not support hijacking")
}

// --- Fix #2 additional: exact prefix match ---

func TestGroupExactPrefixMatch(t *testing.T) {
	app := vayu.New()

	var mwCalled bool
	g := app.Group("/v1")
	g.Use(func(c *vayu.Context, next vayu.NextFunc) {
		mwCalled = true
		next()
	})
	g.GET("", func(c *vayu.Context, next vayu.NextFunc) {
		_, _ = c.Send(vayu.StatusOK, "v1 root")
	})

	// Exact match on /v1 should trigger middleware
	mwCalled = false
	req := httptest.NewRequest("GET", "/v1", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	assert.True(t, mwCalled, "middleware should run for exact prefix match /v1")
	assert.Equal(t, "v1 root", w.Body.String())
}

// --- Fix #3 additional: valid int8 within range ---

func TestBindQueryParamsInt8ValidRange(t *testing.T) {
	app := vayu.New()

	app.GET("/valid-int8", func(c *vayu.Context, next vayu.NextFunc) {
		type Params struct {
			Val int8 `query:"val"`
		}
		params, err := vayu.BindQueryParams[Params](c)
		if err != nil {
			_, _ = c.Send(vayu.StatusBadRequest, err.Error())
			return
		}
		if params.Val == 42 {
			_, _ = c.Send(vayu.StatusOK, "ok")
		} else {
			_, _ = c.Send(vayu.StatusBadRequest, "wrong value")
		}
	})

	req, _ := http.NewRequest("GET", "/valid-int8?val=42", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	assert.Equal(t, vayu.StatusOK, w.Code)
	assert.Equal(t, "ok", w.Body.String())
}

// --- Fix #1 additional: race detection with middleware ---

func TestConcurrentRequestsWithMiddleware(t *testing.T) {
	app := vayu.New()

	// Add multiple middleware
	app.Use(func(c *vayu.Context, next vayu.NextFunc) {
		c.Set("mw", "executed")
		next()
	})

	app.GET("/concurrent", func(c *vayu.Context, next vayu.NextFunc) {
		val, _ := c.Get("mw")
		if val == "executed" {
			_, _ = c.Send(vayu.StatusOK, "ok")
		} else {
			_, _ = c.Send(vayu.StatusInternalServerError, "middleware not executed")
		}
	})

	var wg sync.WaitGroup
	failures := int32(0)
	var mu sync.Mutex

	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := httptest.NewRequest("GET", "/concurrent", nil)
			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)
			if w.Code != vayu.StatusOK || !strings.Contains(w.Body.String(), "ok") {
				mu.Lock()
				failures++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	assert.Equal(t, int32(0), failures, "concurrent requests should all succeed")
}
