package unit

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/kaushiksamanta/vayu"
	"net/http/httptest"
)

// TestSilentMode tests the SilentMode functionality
func TestSilentMode(t *testing.T) {
	// Save original SilentMode value and restore it after the test
	originalSilentMode := vayu.SilentMode
	defer func() { vayu.SilentMode = originalSilentMode }()

	// First test with SilentMode explicitly disabled (verbose logs)
	t.Run("VerboseLogs", func(t *testing.T) {
		// Capture log output
		var logBuffer bytes.Buffer
		log.SetOutput(&logBuffer)
		defer log.SetOutput(os.Stderr) // restore default output

		// Set SilentMode to false
		vayu.SilentMode = "false"

		// Run test that will trigger a panic
		app := vayu.New()
		app.Use(vayu.ErrorHandlerMiddleware(nil)) // Use default error handler
		app.GET("/panic", func(c *vayu.Context, next vayu.NextFunc) {
			panic("test panic with verbose logs")
		})

		req := httptest.NewRequest("GET", "/panic", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		// Check that panic was recovered (response code should be 500)
		if w.Code != vayu.StatusInternalServerError {
			t.Fatalf("Expected status %d, got %d", vayu.StatusInternalServerError, w.Code)
		}

		// Check that panic logs appear when SilentMode is false
		logOutput := logBuffer.String()
		if !strings.Contains(logOutput, "Panic recovered:") {
			t.Errorf("Expected panic logs with SilentMode=false, but found none")
			t.Logf("Log output was: %s", logOutput)
		}
	})

	// Now test with SilentMode explicitly enabled (silent logs)
	t.Run("SilentLogs", func(t *testing.T) {
		// Capture log output
		var logBuffer bytes.Buffer
		log.SetOutput(&logBuffer)
		defer log.SetOutput(os.Stderr) // restore default output

		// Set SilentMode to true
		vayu.SilentMode = "true"

		// Run test that will trigger a panic
		app := vayu.New()
		app.Use(vayu.ErrorHandlerMiddleware(nil)) // Use default error handler
		app.GET("/panic", func(c *vayu.Context, next vayu.NextFunc) {
			panic("test panic with silent logs")
		})

		req := httptest.NewRequest("GET", "/panic", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		// Check that panic was recovered (response code should be 500)
		if w.Code != vayu.StatusInternalServerError {
			t.Fatalf("Expected status %d, got %d", vayu.StatusInternalServerError, w.Code)
		}

		// Check that panic logs do NOT appear when SilentMode is true
		logOutput := logBuffer.String()
		if strings.Contains(logOutput, "Panic recovered:") {
			t.Errorf("Expected NO panic logs with SilentMode=true, but found some")
			t.Logf("Log output was: %s", logOutput)
		}

		// We should still see the regular error log
		if !strings.Contains(logOutput, "Error:") {
			t.Logf("Note: Regular error logs are still shown: %s", logOutput)
		}
	})
}
