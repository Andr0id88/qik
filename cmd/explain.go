package cmd

import (
	"fmt"
	"log" // Used for log.Fatal and log.Fatalf
	"strings"
	"os"

	"qik/internal/ai"
	"qik/internal/clipboard"
	"qik/internal/editor"
	"qik/internal/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// explainLanguage stores the value of the --language flag for the explain command.
	explainLanguage string
	// explainCopyToClipboard stores the value of the --copy flag for the explain command.
	explainCopyToClipboard bool
)

// explainCmd represents the command to generate a simple explanation for a given text.
var explainCmd = &cobra.Command{
	Use:   "explain",
	Short: "Explain a given text in simple terms, output to terminal.",
	Long: `Opens an editor for text input. The text is then sent to Gemini AI
to generate a simple and concise explanation.
The explanation is printed to the terminal by default.
Use --copy to also copy it to the clipboard.
Use --language to specify the desired language of the explanation;
otherwise, the AI will attempt to match the input text's language or use the default.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve API key, respecting verbosity for messages about key source.
		// 'verbose' is a package-level variable from root.go.
		apiKey, err := utils.GetGeminiAPIKey(viper.GetString("geminiApiKey"), verbose)
		if err != nil {
			log.Fatalf("Error getting API key: %v", err)
		}
		if apiKey == "" {
			log.Fatal("Gemini API key not found. Set GEMINI_API_KEY, use 'pass gemini_api_key', or set 'geminiApiKey' in config.")
		}

		// Determine target language for the explanation.
		// The 'explain_text' prompt is designed to infer input language if no override is given.
		targetLanguageForPrompt := AppConfig.DefaultLanguage // Fallback or base for prompt
		if cmd.Flags().Changed("language") {
			targetLanguageForPrompt = explainLanguage
			printVerbose("INFO: Explanation language explicitly set to: %s", targetLanguageForPrompt)
		} else {
			printVerbose("INFO: Explanation language not specified. AI will attempt to match input or use default (%s).", targetLanguageForPrompt)
		}

		// Retrieve the appropriate prompt template for explaining text.
		explainPromptTemplate := AppConfig.Prompts.ExplainText
		if explainPromptTemplate == "" {
			log.Fatal("Error: 'explain_text' prompt not defined in configuration. Check your config file.")
		}

		// Mood is not explicitly used by the 'explain' command's prompt,
		// as the 'explain_text' prompt itself dictates the desired simple and concise tone.

		fmt.Println("Opening editor for text to explain...") // User feedback
		inputText, err := editor.GetTextFromEditor(AppConfig.Editor)
		if err != nil {
			log.Fatalf("Error getting text from editor: %v", err)
		}
		if strings.TrimSpace(inputText) == "" {
			fmt.Println("No input provided. Exiting.") // User feedback
			return
		}

		fmt.Println("Generating explanation...") // User feedback
		// No detailed printVerbose here as language is already covered.

		aiClient, err := ai.NewGeminiClient(cmd.Context(), apiKey, AppConfig.GeminiModel)
		if err != nil {
			log.Fatalf("Error creating AI client: %v", err)
		}

		// The ProcessText function will replace {TEXT} and {LANGUAGE} in the explainPromptTemplate.
		explanation, err := aiClient.ProcessText(cmd.Context(), inputText, explainPromptTemplate, targetLanguageForPrompt)
		if err != nil {
			log.Fatalf("Error generating explanation with AI: %v", err)
		}

		// Display the generated explanation in the terminal.
		fmt.Println("\n--- Explanation ---")
		fmt.Println(strings.TrimSpace(explanation)) // Trim whitespace for cleaner output
		fmt.Println("-------------------")

		// Optionally copy the explanation to the clipboard.
		if explainCopyToClipboard {
			err = clipboard.CopyToClipboard(explanation)
			if err != nil {
				// Non-fatal warning if clipboard operation fails but terminal output succeeded.
				fmt.Fprintf(os.Stderr, "\nWarning: Error copying explanation to clipboard: %v.\n", err)
			} else {
				fmt.Println("\nExplanation also copied to clipboard!")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(explainCmd)
	explainCmd.Flags().StringVarP(&explainLanguage, "language", "l", "", "Language for the explanation. Overrides AI's attempt to match input language.")
	explainCmd.Flags().BoolVarP(&explainCopyToClipboard, "copy", "c", false, "Copy the explanation to the clipboard in addition to printing it.")
}
