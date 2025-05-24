package unit

import (
	"encoding/json"
	"github.com/kaushiksamanta/vayu"
	"net/http/httptest"
	"testing"
)

func TestContextSend(t *testing.T) {
	app := vayu.New()

	app.GET("/send", func(c *vayu.Context, next vayu.NextFunc) {
		_, err := c.Send(vayu.StatusOK, "Hello, World!")
		if err != nil {
			t.Fatalf("Failed to send response: %v", err)
		}
	})

	req := httptest.NewRequest("GET", "/send", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != vayu.StatusOK {
		t.Errorf("Expected status code %d, got %d", vayu.StatusOK, w.Code)
	}

	if w.Body.String() != "Hello, World!" {
		t.Errorf("Expected body 'Hello, World!', got '%s'", w.Body.String())
	}

	if w.Header().Get("Content-Type") != "text/plain" {
		t.Errorf("Expected Content-Type 'text/plain', got '%s'", w.Header().Get("Content-Type"))
	}
}

func TestContextJSON(t *testing.T) {
	app := vayu.New()

	type TestData struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
	}

	testData := TestData{
		Message: "Hello, JSON!",
		Status:  200,
	}

	app.GET("/json", func(c *vayu.Context, next vayu.NextFunc) {
		err := c.JSON(vayu.StatusOK, testData)
		if err != nil {
			t.Fatalf("Failed to send JSON response: %v", err)
		}
	})

	req := httptest.NewRequest("GET", "/json", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != vayu.StatusOK {
		t.Errorf("Expected status code %d, got %d", vayu.StatusOK, w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", w.Header().Get("Content-Type"))
	}

	// Decode the response body
	var response TestData
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	if response.Message != testData.Message {
		t.Errorf("Expected message '%s', got '%s'", testData.Message, response.Message)
	}

	if response.Status != testData.Status {
		t.Errorf("Expected status %d, got %d", testData.Status, response.Status)
	}
}

func TestContextHTML(t *testing.T) {
	app := vayu.New()

	htmlContent := "<h1>Hello, HTML!</h1>"

	app.GET("/html", func(c *vayu.Context, next vayu.NextFunc) {
		_, err := c.HTML(vayu.StatusOK, htmlContent)
		if err != nil {
			t.Fatalf("Failed to send HTML response: %v", err)
		}
	})

	req := httptest.NewRequest("GET", "/html", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != vayu.StatusOK {
		t.Errorf("Expected status code %d, got %d", vayu.StatusOK, w.Code)
	}

	if w.Body.String() != htmlContent {
		t.Errorf("Expected body '%s', got '%s'", htmlContent, w.Body.String())
	}

	if w.Header().Get("Content-Type") != "text/html" {
		t.Errorf("Expected Content-Type 'text/html', got '%s'", w.Header().Get("Content-Type"))
	}
}

func TestContextQuery(t *testing.T) {
	app := vayu.New()

	app.GET("/query", func(c *vayu.Context, next vayu.NextFunc) {
		name := c.Query("name")
		_, err := c.Send(vayu.StatusOK, "Hello, "+name+"!")
		if err != nil {
			t.Fatalf("Failed to send response: %v", err)
		}
	})

	req := httptest.NewRequest("GET", "/query?name=Vayu", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != vayu.StatusOK {
		t.Errorf("Expected status code %d, got %d", vayu.StatusOK, w.Code)
	}

	if w.Body.String() != "Hello, Vayu!" {
		t.Errorf("Expected body 'Hello, Vayu!', got '%s'", w.Body.String())
	}
}

func TestContextSetGet(t *testing.T) {
	app := vayu.New()

	app.GET("/store", func(c *vayu.Context, next vayu.NextFunc) {
		c.Set("key", "value")
		val, ok := c.Get("key")
		if !ok {
			t.Error("Expected key to be found in context store")
		}

		strVal, ok := val.(string)
		if !ok {
			t.Error("Expected value to be a string")
		}

		_, err := c.Send(vayu.StatusOK, strVal)
		if err != nil {
			t.Fatalf("Failed to send response: %v", err)
		}
	})

	req := httptest.NewRequest("GET", "/store", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != vayu.StatusOK {
		t.Errorf("Expected status code %d, got %d", vayu.StatusOK, w.Code)
	}

	if w.Body.String() != "value" {
		t.Errorf("Expected body 'value', got '%s'", w.Body.String())
	}
}
