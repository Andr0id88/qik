package utils

import (
	"fmt"
	"os"
	"os/exec" // Used to execute the 'pass' command-line utility.
	"strings"
)

const (
	// PassEntryName defines the expected name of the API key entry
	// within the 'pass' password manager store.
	PassEntryName = "gemini_api_key"

	// EnvVarName defines the name of the environment variable
	// that can be used to supply the Gemini API key.
	EnvVarName = "GEMINI_API_KEY"
)

// GetGeminiAPIKey retrieves the Gemini API key by checking various sources in a specific order:
// 1. Environment variable (GEMINI_API_KEY).
// 2. 'pass' password manager (entry named 'gemini_api_key').
// 3. A key provided directly from the application's configuration (configKey).
//
// The 'verbose' flag controls whether informational messages about the source of the API key
// and security warnings are printed to the console.
// It returns the API key string or an error if the key cannot be found.
func GetGeminiAPIKey(configKey string, verbose bool) (string, error) {
	// printV is a local helper for conditional verbose printing.
	printV := func(format string, a ...interface{}) {
		if verbose {
			if !strings.HasSuffix(format, "\n") {
				format += "\n"
			}
			fmt.Printf(format, a...)
		}
	}

	// 1. Attempt to retrieve the API key from the environment variable.
	apiKey := os.Getenv(EnvVarName)
	if apiKey != "" {
		printV("Using Gemini API key from environment variable %s.", EnvVarName)
		return apiKey, nil
	}

	// 2. Attempt to retrieve the API key using the 'pass' password manager.
	// First, check if the 'pass' command is available in the system's PATH.
	if _, err := exec.LookPath("pass"); err == nil {
		cmd := exec.Command("pass", PassEntryName)
		output, errPass := cmd.Output() // Execute 'pass gemini_api_key'.
		if errPass == nil {
			// 'pass' command executed successfully.
			key := strings.TrimSpace(string(output))
			if key != "" {
				// The output from 'pass' might contain multiple lines if the entry does.
				// We assume the API key is on the first line.
				key = strings.SplitN(key, "\n", 2)[0]
				printV("Using Gemini API key from 'pass %s'.", PassEntryName)
				return key, nil
			}
			// If 'pass' returns an empty string, it's treated as key not found.
		}
		// If 'pass <entry>' fails (e.g., entry not found), we don't log it verbosely by default,
		// as this is an expected part of the fallback mechanism. A specific error (errPass) occurred.
	} else {
		// 'pass' command itself is not installed or not in PATH.
		printV("'pass' command not found in PATH. Skipping 'pass' lookup for API key.")
	}

	// 3. Attempt to use the API key provided from the application's configuration.
	// This is generally the least secure method and should be used with caution.
	if configKey != "" {
		// This warning is important for security awareness.
		// It's tied to the 'verbose' flag; consider if it should always be shown
		// or controlled by a more specific "security_warnings" flag in the future.
		printV("Warning: Using Gemini API key from the application configuration file. " +
			"For better security, prefer using environment variables or a password manager like 'pass'.")
		return configKey, nil
	}

	// If the API key is not found in any of the checked sources.
	return "", fmt.Errorf("Gemini API key not found. Please set the %s environment variable, "+
		"store it in 'pass' as '%s', or add 'geminiApiKey' to your qik configuration file.",
		EnvVarName, PassEntryName)
}
