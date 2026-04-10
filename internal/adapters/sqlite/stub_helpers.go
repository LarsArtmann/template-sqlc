// Package sqlite provides SQLite-specific database adapters.
// This file contains shared helper functions for stub implementations.
package sqlite

// stubPanic panics with a message indicating the method needs implementation.
func stubPanic() {
	panic("implement me: use actual sqlc generated code")
}

// stubNotImplemented returns an error indicating the method is not yet implemented.
func stubNotImplemented(method string) error {
	return nil // Placeholder - will be replaced with actual error
}
