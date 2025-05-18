package clipboard

import (
	"fmt"
	"github.com/atotto/clipboard" // Cross-platform clipboard package.
)

// CopyToClipboard writes the provided text string to the system clipboard.
// It returns an error if the write operation fails.
func CopyToClipboard(text string) error {
	// clipboard.WriteAll attempts to write the entire string to the clipboard.
	if err := clipboard.WriteAll(text); err != nil {
		// Wrap the error from the clipboard library for more context.
		return fmt.Errorf("failed to write to system clipboard: %w", err)
	}
	return nil
}
