package main

type (
	// Type for custom context key.
	contextKey string
)

const (
	// Context key for user authentication state.
	isAuthenticatedContextKey = contextKey("isAuthenticated")
)
