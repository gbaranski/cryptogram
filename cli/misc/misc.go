package misc

import "fmt"

// WithColor wraps a string with color tags for display in the messages text box.
func WithColor(color, msg string) string {
	return fmt.Sprintf("[%s]%s[-]", color, msg)
}
