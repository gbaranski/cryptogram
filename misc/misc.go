package misc

import (
	"fmt"
	"runtime"
)

// MemStats represents memory stats, all values are in MiB
type MemStats struct {
	Alloc float64
}

// GetMemStats returns MemStats
func GetMemStats() *MemStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return &MemStats{
		Alloc: ByteToMiB(m.HeapAlloc),
	}
}

// ByteToMiB converts bytes to MiB
func ByteToMiB(b uint64) float64 {
	return float64(b) / 1024 / 1024
}

// WithColor wraps a string with color tags for display in the messages text box. Only for UI
func WithColor(color, msg string) string {
	return fmt.Sprintf("[%s]%s[-]", color, msg)
}
