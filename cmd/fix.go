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
	// language stores the value of the --language flag.
	language string
	// promptKey stores the value of the --prompt flag.
	promptKey string
	// englishShorthand stores the value of the --english flag.
	englishShorthand bool
	// moodKey stores the value of the --mood flag.
	moodKey string
)

// fixCmd represents the command for fixing spelling, grammar, flow, and tone of text.
var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Fix spelling, flow, and tone of text, then copy to clipboard.",
	Long: `Opens an editor for text input. The text is then sent to Gemini AI
for spelling correction, flow/tone improvement, and translation (if applicable).
The corrected text is copied to the clipboard.`,
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

		// Determine the target language for corrections.
		targetLanguage := AppConfig.DefaultLanguage
		if language != "" { // --language flag overrides the default from config.
			targetLanguage = language
		}

		// Select the appropriate AI prompt template based on flags or defaults.
		var chosenPromptTemplate string
		if promptKey != "" { // --prompt flag takes precedence if set.
			switch strings.ToLower(promptKey) {
			case "default":
				chosenPromptTemplate = AppConfig.Prompts.Default
			case "english_fix_only":
				chosenPromptTemplate = AppConfig.Prompts.EnglishFixOnly
				targetLanguage = "English" // This prompt implies English output.
			default:
				// Warn user about an unrecognized prompt key and fall back to default.
				fmt.Printf("Warning: Unknown prompt key '%s'. Using default prompt for language %s.\n", promptKey, targetLanguage)
				chosenPromptTemplate = AppConfig.Prompts.Default
			}
		} else {
			// If no specific prompt key, choose based on target language.
			if strings.EqualFold(targetLanguage, "English") && AppConfig.Prompts.EnglishFixOnly != "" {
				chosenPromptTemplate = AppConfig.Prompts.EnglishFixOnly
			} else {
				chosenPromptTemplate = AppConfig.Prompts.Default
			}
		}
		if chosenPromptTemplate == "" {
			log.Fatal("Error: No suitable prompt template could be determined. Check configuration.")
		}

		// Determine the desired mood/tone for the text.
		selectedMoodKey := AppConfig.DefaultMood
		if cmd.Flags().Changed("mood") { // --mood flag overrides the default from config.
			selectedMoodKey = moodKey
		}

		moodInstructionText := ""
		if mood, ok := AppConfig.Moods[selectedMoodKey]; ok {
			moodInstructionText = mood.Instruction
		} else {
			// Warn if a specific, non-default mood was requested but not found.
			if selectedMoodKey != "" && selectedMoodKey != AppConfig.DefaultMood {
				fmt.Printf("Warning: Mood key '%s' not found. Using default ('%s') or applying no specific mood styling if default is also misconfigured.\n", selectedMoodKey, AppConfig.DefaultMood)
			}
			// Fallback to default mood's instruction.
			if defaultMood, okDefault := AppConfig.Moods[AppConfig.DefaultMood]; okDefault {
				moodInstructionText = defaultMood.Instruction
				// selectedMoodKey is not updated here to AppConfig.DefaultMood as the warning above already informed the user.
			} else {
				fmt.Printf("Warning: Default mood '%s' also not found or has no instruction. No specific mood styling applied.\n", AppConfig.DefaultMood)
			}
		}

		// Construct the final prompt by injecting the mood instruction.
		// {LANGUAGE} and {TEXT} placeholders will be filled by ai.ProcessText.
		finalPrompt := strings.ReplaceAll(chosenPromptTemplate, "{MOOD_INSTRUCTION}", moodInstructionText)

		fmt.Println("Opening editor for input...") // User feedback
		inputText, err := editor.GetTextFromEditor(AppConfig.Editor)
		if err != nil {
			log.Fatalf("Error getting text from editor: %v", err)
		}
		if strings.TrimSpace(inputText) == "" {
			fmt.Println("No input provided. Exiting.") // User feedback
			return
		}

		fmt.Println("Processing text...") // User feedback
		printVerbose("INFO: Using Language: %s, Mood: %s, PromptKey: %s", targetLanguage, selectedMoodKey, promptKey)

		aiClient, err := ai.NewGeminiClient(cmd.Context(), apiKey, AppConfig.GeminiModel)
		if err != nil {
			log.Fatalf("Error creating AI client: %v", err)
		}

		processedText, err := aiClient.ProcessText(cmd.Context(), inputText, finalPrompt, targetLanguage)
		if err != nil {
			log.Fatalf("Error processing text with AI: %v", err)
		}

		// Attempt to copy the processed text to the clipboard.
		err = clipboard.CopyToClipboard(processedText)
		if err != nil {
			// If clipboard fails, print the output to terminal as a fallback.
			fmt.Fprintf(os.Stderr, "Error copying to clipboard: %v.\n", err)
			fmt.Println("\n--- Corrected Text (Clipboard Failed) ---")
			fmt.Println(processedText)
			fmt.Println("-----------------------------------------")
			log.Fatalf("Failed to copy to clipboard.") // Still exit with error for scripting.
		} else {
			fmt.Println("Corrected text copied to clipboard!") // User feedback
		}
	},
}

func init() {
	rootCmd.AddCommand(fixCmd)
	fixCmd.Flags().StringVarP(&language, "language", "l", "", "Target language (e.g., Norwegian, English). Overrides config default.")
	fixCmd.Flags().BoolVarP(&englishShorthand, "english", "e", false, "Shorthand for --language English and 'english_fix_only' prompt.")
	fixCmd.Flags().StringVarP(&promptKey, "prompt", "p", "", "Key of the prompt template to use (e.g., 'default', 'english_fix_only').")
	fixCmd.Flags().StringVarP(&moodKey, "mood", "m", "", "Desired mood/tone (e.g., professional, casual). Overrides config default.")

	// PersistentPreRunE is used to handle interactions between flags,
	// specifically making the --english shorthand flag work as intended.
	fixCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if englishShorthand { // If --english or -e is used
			if language != "" && !strings.EqualFold(language, "English") {
				return fmt.Errorf("--english flag is incompatible with --language specifying a different language than English")
			}
			language = "English" // Set the target language to English.
			// If --prompt was not explicitly set by the user, default to 'english_fix_only' prompt.
			if !cmd.Flags().Changed("prompt") {
				promptKey = "english_fix_only"
			}
		}
		return nil
	}
}
