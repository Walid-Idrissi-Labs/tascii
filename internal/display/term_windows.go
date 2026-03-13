//go:build windows

package display

// getTermWidth is not implemented on Windows; caller will use the env-var fallback.
func getTermWidth() int { return 0 }