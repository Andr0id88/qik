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
	// answerLanguage stores the value of the --language flag for the answer command.
	answerLanguage string
	// answerMoodKey stores the value of the --mood flag for the answer command.
	answerMoodKey string
	// answerCopyToClipboard stores the value of the --copy flag for the answer command.
	answerCopyToClipboard bool
)

// answerCmd represents the command to get an answer to a user's question.
var answerCmd = &cobra.Command{
	Use:   "answer",
	Short: "Answer a given question, output to terminal.",
	Long: `Opens an editor for question input. The question is then sent to Gemini AI
to generate an answer. The answer can be adjusted for language and mood.
The answer is printed to the terminal by default.
Use --copy to also copy it to the clipboard.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve API key, respecting verbosity for messages about key source.
		apiKey, err := utils.GetGeminiAPIKey(viper.GetString("geminiApiKey"), verbose)
		if err != nil {
			log.Fatalf("Error getting API key: %v", err)
		}
		if apiKey == "" {
			log.Fatal("Gemini API key not found. Set GEMINI_API_KEY, use 'pass gemini_api_key', or set 'geminiApiKey' in config.")
		}

		// Determine target language for the answer.
		targetLanguage := AppConfig.DefaultLanguage
		if cmd.Flags().Changed("language") {
			targetLanguage = answerLanguage
		}

		// Determine mood for the answer.
		selectedMoodKey := AppConfig.DefaultMood
		if cmd.Flags().Changed("mood") {
			selectedMoodKey = answerMoodKey
		}

		moodInstructionText := ""
		if mood, ok := AppConfig.Moods[selectedMoodKey]; ok {
			moodInstructionText = mood.Instruction
		} else {
			// Warn if a specific, non-default mood was requested but not found.
			if selectedMoodKey != "" && selectedMoodKey != AppConfig.DefaultMood {
				fmt.Printf("Warning: Mood key '%s' not found in configuration. Using default mood ('%s').\n", selectedMoodKey, AppConfig.DefaultMood)
			}
			// Fallback to default mood's instruction.
			if defaultMood, okDefault := AppConfig.Moods[AppConfig.DefaultMood]; okDefault {
				moodInstructionText = defaultMood.Instruction
				selectedMoodKey = AppConfig.DefaultMood // Ensure selectedMoodKey reflects the actual mood used.
			} else {
				// This should be rare if AppConfig.DefaultMood is always valid.
				fmt.Printf("Warning: Default mood '%s' not found or has no instruction. Applying no specific mood styling.\n", AppConfig.DefaultMood)
			}
		}

		// Provide a default instruction for "neutral" if its configured instruction is empty,
		// to ensure the prompt to the AI is well-formed.
		if moodInstructionText == "" && selectedMoodKey == "neutral" {
			moodInstructionText = "Answer in a standard, helpful, and informative tone."
		}

		// Retrieve the appropriate prompt template for answering questions.
		answerPromptTemplate := AppConfig.Prompts.AnswerQuestion
		if answerPromptTemplate == "" {
			log.Fatal("Error: 'answer_question' prompt not defined in configuration. Check your config file.")
		}

		// Prepare the final prompt by injecting the mood instruction.
		// {LANGUAGE} and {TEXT} placeholders will be filled by ai.ProcessText.
		promptWithMood := strings.ReplaceAll(answerPromptTemplate, "{MOOD_INSTRUCTION}", moodInstructionText)

		fmt.Println("Opening editor for your question...")
		inputText, err := editor.GetTextFromEditor(AppConfig.Editor)
		if err != nil {
			log.Fatalf("Error getting text from editor: %v", err)
		}
		if strings.TrimSpace(inputText) == "" {
			fmt.Println("No question provided. Exiting.")
			return
		}

		fmt.Println("Generating answer...") // User feedback indicating AI call
		printVerbose("INFO: Answering with Language: %s, Mood: %s", targetLanguage, selectedMoodKey)

		aiClient, err := ai.NewGeminiClient(cmd.Context(), apiKey, AppConfig.GeminiModel)
		if err != nil {
			log.Fatalf("Error creating AI client: %v", err)
		}

		answer, err := aiClient.ProcessText(cmd.Context(), inputText, promptWithMood, targetLanguage)
		if err != nil {
			log.Fatalf("Error generating answer with AI: %v", err)
		}

		// Display the generated answer in the terminal.
		fmt.Println("\n--- Answer ---")
		fmt.Println(strings.TrimSpace(answer)) // Trim whitespace for cleaner output
		fmt.Println("--------------")

		// Optionally copy the answer to the clipboard.
		if answerCopyToClipboard {
			err = clipboard.CopyToClipboard(answer)
			if err != nil {
				// Non-fatal warning if clipboard operation fails but terminal output succeeded.
				fmt.Fprintf(os.Stderr, "\nWarning: Error copying answer to clipboard: %v.\n", err)
			} else {
				fmt.Println("\nAnswer also copied to clipboard!")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(answerCmd)
	answerCmd.Flags().StringVarP(&answerLanguage, "language", "l", "", "Language for the answer (e.g., Norwegian, English). Overrides config default language.")
	answerCmd.Flags().StringVarP(&answerMoodKey, "mood", "m", "", "Desired mood/tone for the answer (e.g., professional, neutral). Overrides config default mood.")
	answerCmd.Flags().BoolVarP(&answerCopyToClipboard, "copy", "c", false, "Copy the answer to the clipboard in addition to printing it.")
}
