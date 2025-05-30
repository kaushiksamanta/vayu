package vayu

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// This file contains all generic functions for type-safe operations in the Vayu framework.
// These functions provide compile-time type safety through Go's generics.

// JSON sends a JSON response with the given status code and typed object.
// This provides compile-time type safety through generics.
func JSONResponse[T any](c *Context, code int, obj T) error {
	return c.JSON(code, obj)
}

// BindJSONBody binds the request body as JSON to a specific type.
// This provides compile-time type safety through generics.
// Usage: user, err := vayu.BindJSONBody[User](c)
func BindJSONBody[T any](c *Context) (T, error) {
	var result T
	err := c.BindJSON(&result)
	return result, err
}

// GetValue retrieves a typed value from the request context with the given key.
// This provides compile-time type safety through generics.
// Usage: user, ok := vayu.GetValue[User](c, "user")
func GetValue[T any](c *Context, key string) (T, bool) {
	var zero T
	val, ok := c.Get(key)
	if !ok {
		return zero, false
	}
	
	// Try type assertion
	if typed, ok := val.(T); ok {
		return typed, true
	}
	
	return zero, false
}

// SetValue stores a typed value in the request context with the given key.
// This provides compile-time type safety through generics.
// Usage: vayu.SetValue(c, "user", user)
func SetValue[T any](c *Context, key string, value T) {
	c.Set(key, value)
}

// MustBindJSONBody binds the request body as JSON to a specific type.
// This function panics if binding fails.
// Use this only when you're certain the binding will succeed.
// Usage: user := vayu.MustBindJSONBody[User](c)
func MustBindJSONBody[T any](c *Context) T {
	result, err := BindJSONBody[T](c)
	if err != nil {
		panic(err)
	}
	return result
}

// BindQueryJSON binds a JSON string from a query parameter to a specific type.
// This provides compile-time type safety through generics.
// Usage: filter, err := vayu.BindQueryJSON[Filter](c, "filter")
func BindQueryJSON[T any](c *Context, paramName string) (T, error) {
	var result T
	jsonStr := c.Query(paramName)
	if jsonStr == "" {
		return result, fmt.Errorf("query parameter %s is empty or not found", paramName)
	}
	
	// URL-decode the parameter if needed
	decodedStr, err := url.QueryUnescape(jsonStr)
	if err != nil {
		// If decoding fails, try with the original string
		decodedStr = jsonStr
	}

	// Try to unmarshal the JSON
	err = json.Unmarshal([]byte(decodedStr), &result)
	return result, err
}

// MustBindQueryJSON binds a JSON string from a query parameter to a specific type.
// This function panics if binding fails.
// Usage: filter := vayu.MustBindQueryJSON[Filter](c, "filter")
func MustBindQueryJSON[T any](c *Context, paramName string) T {
	result, err := BindQueryJSON[T](c, paramName)
	if err != nil {
		panic(err)
	}
	return result
}

// BindParamJSON binds a JSON string from a path parameter to a specific type.
// This provides compile-time type safety through generics.
// Usage: filter, err := vayu.BindParamJSON[Filter](c, "filter")
func BindParamJSON[T any](c *Context, paramName string) (T, error) {
	var result T
	jsonStr, exists := c.Params[paramName]
	if !exists || jsonStr == "" {
		return result, fmt.Errorf("path parameter %s is empty or not found", paramName)
	}
	
	// URL-decode the parameter if needed
	decodedStr, err := url.QueryUnescape(jsonStr)
	if err != nil {
		// If decoding fails, try with the original string
		decodedStr = jsonStr
	}
	
	// Try to unmarshal the JSON
	err = json.Unmarshal([]byte(decodedStr), &result)
	return result, err
}

// MustBindParamJSON binds a JSON string from a path parameter to a specific type.
// This function panics if binding fails.
// Usage: filter := vayu.MustBindParamJSON[Filter](c, "filter")
func MustBindParamJSON[T any](c *Context, paramName string) T {
	result, err := BindParamJSON[T](c, paramName)
	if err != nil {
		panic(err)
	}
	return result
}

// BindQueryParams binds query parameters to a struct based on struct tags.
// This provides compile-time type safety through generics.
// Usage: params := vayu.BindQueryParams[SearchParams](c)
// Define your struct with `query` tags: type SearchParams struct { Term string `query:"q"` }
func BindQueryParams[T any](c *Context) (T, error) {
	var result T
	val := reflect.ValueOf(&result).Elem()
	typ := val.Type()
	
	errs := make([]string, 0)
	processed := false
	
	// Process each field with a `query` tag
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		queryTag := field.Tag.Get("query")
		if queryTag == "" {
			continue
		}
		
		processed = true
		queryValue := c.Query(queryTag)
		if queryValue == "" {
			// Check if the field is required
			if requiredTag, ok := field.Tag.Lookup("required"); ok && requiredTag == "true" {
				errs = append(errs, fmt.Sprintf("required query parameter %s missing", queryTag))
			}
			continue
		}
		
		fieldValue := val.Field(i)
		if !fieldValue.CanSet() {
			errs = append(errs, fmt.Sprintf("field %s cannot be set (is it unexported?)", field.Name))
			continue
		}
		
		// Convert the string value to the appropriate field type
		if err := setFieldFromString(fieldValue, queryValue); err != nil {
			errs = append(errs, fmt.Sprintf("parameter %s: %v", queryTag, err))
		}
	}
	
	if !processed {
		return result, fmt.Errorf("no fields with 'query' tag found in type %s", typ.Name())
	}
	
	if len(errs) > 0 {
		return result, fmt.Errorf("binding query parameters: %s", strings.Join(errs, "; "))
	}
	
	return result, nil
}

// MustBindQueryParams binds query parameters to a struct and panics if binding fails.
// Usage: params := vayu.MustBindQueryParams[SearchParams](c)
func MustBindQueryParams[T any](c *Context) T {
	result, err := BindQueryParams[T](c)
	if err != nil {
		panic(err)
	}
	return result
}

// setFieldFromString converts a string value to the appropriate field type and sets it
func setFieldFromString(fieldValue reflect.Value, value string) error {
	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.SetString(value)
		return nil
		
	case reflect.Bool:
		v, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("cannot convert '%s' to bool: %w", value, err)
		}
		fieldValue.SetBool(v)
		return nil
		
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if fieldValue.Type().String() == "time.Duration" {
			v, err := time.ParseDuration(value)
			if err != nil {
				return fmt.Errorf("cannot convert '%s' to duration: %w", value, err)
			}
			fieldValue.Set(reflect.ValueOf(v))
			return nil
		}
		
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("cannot convert '%s' to int: %w", value, err)
		}
		fieldValue.SetInt(v)
		return nil
		
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("cannot convert '%s' to uint: %w", value, err)
		}
		fieldValue.SetUint(v)
		return nil
		
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("cannot convert '%s' to float: %w", value, err)
		}
		fieldValue.SetFloat(v)
		return nil
		
	case reflect.Slice:
		// Handle comma-separated values for slices
		values := strings.Split(value, ",")
		sliceType := fieldValue.Type().Elem().Kind()
		slice := reflect.MakeSlice(fieldValue.Type(), len(values), len(values))
		
		for i, v := range values {
			elemValue := slice.Index(i)
			v = strings.TrimSpace(v)
			
			switch sliceType {
			case reflect.String:
				elemValue.SetString(v)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				intVal, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					return fmt.Errorf("cannot convert '%s' to int in slice: %w", v, err)
				}
				elemValue.SetInt(intVal)
			case reflect.Float32, reflect.Float64:
				floatVal, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return fmt.Errorf("cannot convert '%s' to float in slice: %w", v, err)
				}
				elemValue.SetFloat(floatVal)
			case reflect.Bool:
				boolVal, err := strconv.ParseBool(v)
				if err != nil {
					return fmt.Errorf("cannot convert '%s' to bool in slice: %w", v, err)
				}
				elemValue.SetBool(boolVal)
			default:
				return fmt.Errorf("unsupported slice element type: %s", sliceType)
			}
		}
		
		fieldValue.Set(slice)
		return nil
		
	default:
		return fmt.Errorf("unsupported field type: %s", fieldValue.Kind())
	}
}
