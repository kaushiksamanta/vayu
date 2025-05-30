package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kaushiksamanta/vayu"
	"github.com/stretchr/testify/assert"
)

func TestBindQueryJSON(t *testing.T) {
	app := vayu.New()
	results := make(chan bool, 3)
	
	// Setup a test route for valid JSON
	app.GET("/test-valid", func(c *vayu.Context, next vayu.NextFunc) {
		type Filter struct {
			Name  string `json:"name"`
			Age   int    `json:"age"`
			Admin bool   `json:"admin"`
		}
		
		filter, err := vayu.BindQueryJSON[Filter](c, "filter")
		
		if err == nil && 
		   filter.Name == "John" && 
		   filter.Age == 30 && 
		   filter.Admin == true {
			c.Writer.WriteHeader(200)
			results <- true
		} else {
			c.Writer.WriteHeader(400)
			results <- false
		}
	})
	
	// Setup a route for invalid JSON
	app.GET("/test-invalid", func(c *vayu.Context, next vayu.NextFunc) {
		type Filter struct {
			Name string `json:"name"`
		}
		
		_, err := vayu.BindQueryJSON[Filter](c, "filter")
		
		if err != nil {
			c.Writer.WriteHeader(400)
			results <- true // Expected to get an error
		} else {
			c.Writer.WriteHeader(200)
			results <- false // Shouldn't succeed
		}
	})
	
	// Setup a route for missing parameter
	app.GET("/test-missing", func(c *vayu.Context, next vayu.NextFunc) {
		type Filter struct {
			Name string `json:"name"`
		}
		
		_, err := vayu.BindQueryJSON[Filter](c, "filter")
		
		if err != nil && err.Error() != "" { // Should contain error message
			c.Writer.WriteHeader(400)
			results <- true // Expected to get an error
		} else {
			c.Writer.WriteHeader(200)
			results <- false // Shouldn't succeed
		}
	})
	
	// Test 1: Valid JSON
	reqValid, _ := http.NewRequest("GET", "/test-valid?filter=%7B%22name%22%3A%22John%22%2C%22age%22%3A30%2C%22admin%22%3Atrue%7D", nil)
	respValid := httptest.NewRecorder()
	app.ServeHTTP(respValid, reqValid)
	assert.Equal(t, 200, respValid.Code)
	assert.True(t, <-results)
	
	// Test 2: Invalid JSON
	reqInvalid, _ := http.NewRequest("GET", "/test-invalid?filter=invalid-json", nil)
	respInvalid := httptest.NewRecorder()
	app.ServeHTTP(respInvalid, reqInvalid)
	assert.Equal(t, 400, respInvalid.Code)
	assert.True(t, <-results)
	
	// Test 3: Missing parameter
	reqMissing, _ := http.NewRequest("GET", "/test-missing", nil)
	respMissing := httptest.NewRecorder()
	app.ServeHTTP(respMissing, reqMissing)
	assert.Equal(t, 400, respMissing.Code)
	assert.True(t, <-results)
}

func TestBindParamJSON(t *testing.T) {
	app := vayu.New()
	results := make(chan bool, 2)
	
	// Setup a test route with path parameter
	app.GET("/test/:config", func(c *vayu.Context, next vayu.NextFunc) {
		type Config struct {
			Theme  string   `json:"theme"`
			Colors []string `json:"colors"`
		}
		
		config, err := vayu.BindParamJSON[Config](c, "config")
		
		if err == nil && 
		   config.Theme == "dark" && 
		   len(config.Colors) == 2 && 
		   config.Colors[0] == "black" && 
		   config.Colors[1] == "blue" {
			c.Writer.WriteHeader(200)
			results <- true
		} else {
			c.Writer.WriteHeader(400)
			results <- false
		}
	})
	
	// Setup a route for invalid JSON path parameter
	app.GET("/invalid/:config", func(c *vayu.Context, next vayu.NextFunc) {
		type Config struct {
			Theme string `json:"theme"`
		}
		
		_, err := vayu.BindParamJSON[Config](c, "config")
		
		if err != nil {
			c.Writer.WriteHeader(400)
			results <- true // Expected to get an error
		} else {
			c.Writer.WriteHeader(200)
			results <- false // Shouldn't succeed
		}
	})
	
	// Test 1: Valid JSON in path parameter
	// We're using a pre-encoded JSON here to avoid issues with URL escaping
	reqValid, _ := http.NewRequest("GET", "/test/%7B%22theme%22%3A%22dark%22%2C%22colors%22%3A%5B%22black%22%2C%22blue%22%5D%7D", nil)
	respValid := httptest.NewRecorder()
	app.ServeHTTP(respValid, reqValid)
	assert.Equal(t, 200, respValid.Code)
	assert.True(t, <-results)
	
	// Test 2: Invalid JSON in path parameter
	reqInvalid, _ := http.NewRequest("GET", "/invalid/not-json", nil)
	respInvalid := httptest.NewRecorder()
	app.ServeHTTP(respInvalid, reqInvalid)
	assert.Equal(t, 400, respInvalid.Code)
	assert.True(t, <-results)
}

func TestBindQueryParams(t *testing.T) {
	app := vayu.New()

	// Setup a test route
	app.GET("/search", func(c *vayu.Context, next vayu.NextFunc) {
		type SearchParams struct {
			Term       string        `query:"q"`
			Page       int           `query:"page"`
			PerPage    int           `query:"per_page"`
			Filter     string        `query:"filter"`
			SortBy     string        `query:"sort"`
			Descending bool          `query:"desc"`
			MinPrice   float64       `query:"min_price"`
			Tags       []string      `query:"tags"`
			Ratings    []int         `query:"ratings"`
			Timeout    time.Duration `query:"timeout"`
		}

		params, err := vayu.BindQueryParams[SearchParams](c)

		if assert.NoError(t, err) {
			assert.Equal(t, "golang", params.Term)
			assert.Equal(t, 2, params.Page)
			assert.Equal(t, 20, params.PerPage)
			assert.Equal(t, "books", params.Filter)
			assert.Equal(t, "price", params.SortBy)
			assert.Equal(t, true, params.Descending)
			assert.Equal(t, 10.5, params.MinPrice)
			assert.Equal(t, []string{"programming", "backend"}, params.Tags)
			assert.Equal(t, []int{4, 5}, params.Ratings)
			assert.Equal(t, 30*time.Second, params.Timeout)
			c.Writer.WriteHeader(200)
		} else {
			c.Writer.WriteHeader(400)
		}
	})

	// Create a request with multiple query params
	req, _ := http.NewRequest("GET", "/search?q=golang&page=2&per_page=20&filter=books&sort=price&desc=true&min_price=10.5&tags=programming,backend&ratings=4,5&timeout=30s", nil)
	resp := httptest.NewRecorder()

	// Serve the request
	app.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	// Test with required parameter
	app.GET("/required-test", func(c *vayu.Context, next vayu.NextFunc) {
		type RequiredParams struct {
			Term string `query:"q" required:"true"`
			Page int    `query:"page"`
		}

		_, err := vayu.BindQueryParams[RequiredParams](c)

		if err != nil {
			c.Writer.WriteHeader(400)
		} else {
			c.Writer.WriteHeader(200)
		}
	})

	reqMissing, _ := http.NewRequest("GET", "/required-test?page=2", nil)
	respMissing := httptest.NewRecorder()

	app.ServeHTTP(respMissing, reqMissing)

	assert.Equal(t, 400, respMissing.Code)
}

func TestBindQueryParamErrors(t *testing.T) {
	app := vayu.New()

	// Setup a test route for type errors
	app.GET("/type-errors", func(c *vayu.Context, next vayu.NextFunc) {
		type TypeErrorParams struct {
			Age    int     `query:"age"`
			Amount float64 `query:"amount"`
			Active bool    `query:"active"`
		}

		_, err := vayu.BindQueryParams[TypeErrorParams](c)
		if err != nil {
			c.Writer.WriteHeader(400)
			return
		}
		c.Writer.WriteHeader(200)
	})

	// Test with invalid types
	reqInvalid, _ := http.NewRequest("GET", "/type-errors?age=notanumber&amount=alsowrong&active=notbool", nil)
	respInvalid := httptest.NewRecorder()

	app.ServeHTTP(respInvalid, reqInvalid)

	assert.Equal(t, 400, respInvalid.Code)
}
