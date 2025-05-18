package editor

import (
	"fmt"
	"io/ioutil" // Used for TempFile and ReadFile. Consider os.CreateTemp and os.ReadFile for Go >= 1.16.
	"os"
	"os/exec" // Used to run the external editor command.
)

// GetTextFromEditor facilitates user text input by:
// 1. Creating a temporary file.
// 2. Opening this file in the user-specified command-line editor (editorCmd).
// 3. Waiting for the editor process to terminate (indicating the user has finished editing).
// 4. Reading the content of the temporary file.
// 5. Deleting the temporary file.
// It returns the content of the file as a string or an error if any step fails.
// The editorCmd should be a command that blocks until the user saves and closes the file (e.g., "nvim", "vim", "nano", "code --wait").
func GetTextFromEditor(editorCmd string) (string, error) {
	// Create a temporary file with a "qik-" prefix for easy identification.
	// An empty first argument to TempFile means it will use the default directory for temporary files (e.g., /tmp).
	tempFile, err := ioutil.TempFile("", "qik-*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file for editor input: %w", err)
	}
	tempFilePath := tempFile.Name()

	// Close the file handle immediately after creation.
	// This is crucial because some editors (especially on some OSes)
	// might require exclusive access or fail if the file is still held open by this program.
	if err := tempFile.Close(); err != nil {
		// If closing fails, attempt to clean up the created file before returning the error.
		os.Remove(tempFilePath)
		return "", fmt.Errorf("failed to close temporary file handle (%s) before editing: %w", tempFilePath, err)
	}

	// Ensure the temporary file is deleted when this function returns, regardless of success or failure.
	// This defer statement is placed after the initial file creation and closing checks
	// to avoid attempting to remove a file that wasn't successfully created or properly closed.
	defer os.Remove(tempFilePath)

	// Prepare the command to run the specified editor with the temporary file.
	cmd := exec.Command(editorCmd, tempFilePath)

	// Connect the editor's standard input, output, and error streams to the current process's streams.
	// This allows the editor to interact directly with the user's terminal.
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the editor command and wait for it to complete.
	// The .Run() method blocks until the editor process exits.
	if err := cmd.Run(); err != nil {
		// Provide a more helpful error message if the editor command fails.
		return "", fmt.Errorf("error running editor command '%s' on file '%s': %w. Ensure the editor is in your PATH and configured to block until closed (e.g., 'code --wait' for VS Code).", editorCmd, tempFilePath, err)
	}

	// Read the content from the temporary file after the editor has been closed.
	// For Go 1.16+, os.ReadFile(tempFilePath) is preferred over ioutil.ReadFile.
	content, err := ioutil.ReadFile(tempFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read content from temporary file '%s' after editing: %w", tempFilePath, err)
	}

	return string(content), nil
}
