package vayu

// Response helpers for common HTTP status code responses

// OK sends a 200 OK response with the given content.
func (c *Context) OK(content interface{}) error {
	return c.JSON(StatusOK, content)
}

// Created sends a 201 Created response with the given content.
func (c *Context) Created(content interface{}) error {
	return c.JSON(StatusCreated, content)
}

// NoContent sends a 204 No Content response.
func (c *Context) NoContent() (int, error) {
	c.Writer.WriteHeader(StatusNoContent)
	return 0, nil
}

// BadRequest sends a 400 Bad Request response with the given error message.
func (c *Context) BadRequest(message string) error {
	return c.JSON(StatusBadRequest, map[string]string{"error": message})
}

// Unauthorized sends a 401 Unauthorized response with the given error message.
func (c *Context) Unauthorized(message string) error {
	return c.JSON(StatusUnauthorized, map[string]string{"error": message})
}

// Forbidden sends a 403 Forbidden response with the given error message.
func (c *Context) Forbidden(message string) error {
	return c.JSON(StatusForbidden, map[string]string{"error": message})
}

// NotFound sends a 404 Not Found response with the given error message.
func (c *Context) NotFound(message string) error {
	return c.JSON(StatusNotFound, map[string]string{"error": message})
}

// InternalServerError sends a 500 Internal Server Error response with the given error message.
func (c *Context) InternalServerError(message string) error {
	return c.JSON(StatusInternalServerError, map[string]string{"error": message})
}
