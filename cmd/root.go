package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"qik/internal/config" // Local package for application configuration structures.

	"github.com/spf13/cobra" // CLI framework.
	"github.com/spf13/viper"  // Configuration management.
	"gopkg.in/yaml.v3"        // YAML marshalling for default config.
)

var (
	// cfgFile holds the path to the config file specified by the user via a flag.
	cfgFile string
	// AppConfig holds the loaded application configuration. It's populated by initConfig.
	AppConfig config.Config
	// verbose controls whether verbose logging is enabled. Set by a persistent flag.
	verbose bool
)

// defaultPromptsConfig stores the application's built-in default prompt templates.
// It's initialized in the package init() function and used if prompts are missing
// from the user's configuration or when generating a new default config file.
var defaultPromptsConfig config.Prompts

// rootCmd represents the base command when qik is called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "qik",
	Short: "qik: A CLI tool to fix, explain, answer, and adjust text using Gemini AI.",
	Long: `qik allows you to fix spelling, flow, and tone of text,
get a simple explanation of a given text, or get an answer to a question.
It uses Gemini AI and by default prints explanations/answers to the terminal.

qik primarily targets Norwegian but can be configured for English or other languages.
Your Gemini API key should be stored securely, for example, in the 'pass' password manager
under the name 'gemini_api_key' or as an environment variable 'GEMINI_API_KEY'.`,
}

// printVerbose prints the formatted message to stdout only if the verbose flag is true.
// It ensures a newline is appended if not already present in the format string.
func printVerbose(format string, a ...interface{}) {
	if verbose {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		fmt.Printf(format, a...)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is the main entry point called by main.main(). It only needs to happen once.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// Critical errors during command execution are printed to stderr.
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// init is a Go special function called when the package is initialized.
// It sets up default prompt configurations and registers Cobra's initialization function.
func init() {
	// Define the application's built-in default AI prompt templates.
	// These are used as fallbacks or for generating a new configuration file.
	defaultPromptsConfig = config.Prompts{
		Default: `You are an expert proofreader and linguistic assistant.
Your primary task is to meticulously review the following text.
Correct all spelling and grammatical errors.
Improve the flow and clarity of the text, rephrasing sentences or restructuring paragraphs if necessary to make it sound natural and well-written.
The final output should be in {LANGUAGE}.
{MOOD_INSTRUCTION}
Do NOT include any preambles, apologies, or explanations in your response. Only return the corrected and refined text.

Original text to process:
---
{TEXT}
---`,
		EnglishFixOnly: `You are an expert English proofreader.
Your primary task is to meticulously review the following English text.
Correct all spelling and grammatical errors.
Improve the flow and clarity of the text, rephrasing sentences or restructuring paragraphs if necessary to make it sound natural and well-written.
The text is already in English, so no translation is needed.
{MOOD_INSTRUCTION}
Do NOT include any preambles, apologies, or explanations in your response. Only return the corrected and refined text.

Original text to process:
---
{TEXT}
---`,
		ExplainText: `You are an expert at simplifying complex topics.
The user will provide a piece of text. Your task is to explain the main concepts or what the person is talking about in that text.
The explanation should be:
1. Simple and easy to understand, even for someone not familiar with the topic.
2. Concise and to the point. Aim for a short summary.
3. If a specific output language is requested via the {LANGUAGE} placeholder, use that language. Otherwise, attempt to provide the explanation in the SAME language as the input text ({TEXT}).
Do NOT include any preambles, apologies, or phrases like "This text is about...". Just provide the explanation directly.

Text to explain:
---
{TEXT}
---`,
		AnswerQuestion: `You are an intelligent and helpful assistant.
The user will provide a question. Your task is to provide a clear, concise, and accurate answer to that question.
Consider the following when formulating your response:
1. Directly address the question asked.
2. Provide the answer in the {LANGUAGE} language.
3. Adjust the tone of your answer according to the specified mood: {MOOD_INSTRUCTION}
   If no specific mood instruction is given for "neutral", answer in a helpful and informative default tone.
Do NOT include any preambles like "Here is the answer to your question:" or "The answer is:". Just provide the answer directly.

Question to answer:
---
{TEXT}
---`,
	}

	// cobra.OnInitialize registers functions to be called when Cobra initializes.
	// initConfig will be called to load application configuration.
	cobra.OnInitialize(initConfig)

	// Define persistent flags, available to the root command and all subcommands.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/qik/config.yaml or ./config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output for detailed logging.")
}

// getDefaultConfigPath determines the default expected path for the qik configuration file.
func getDefaultConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get user home directory: %w", err)
	}
	configDir := filepath.Join(home, ".config", "qik")
	return filepath.Join(configDir, "config.yaml"), nil
}

// getDefaultMoods returns a map of predefined mood configurations for the application.
// These are used when generating a default config file or if moods are missing from user's config.
func getDefaultMoods() map[string]config.MoodInstruction {
	// Descriptions are for the 'list-moods' command.
	// Instructions are injected into AI prompts.
	return map[string]config.MoodInstruction{
		"neutral":    {"Standard correction without specific tone alteration. Relies on base prompt's natural styling.", ""},
		"professional": {"Refine text to be formal, objective, and suitable for business or academic contexts.", "Additionally, adjust the tone of the text to be highly professional, formal, and objective. Avoid colloquialisms and ensure a polished, business-like feel."},
		"casual":     {"Make the text sound more relaxed, friendly, and conversational.", "Additionally, adjust the tone of the text to be more casual, friendly, and conversational. Use simpler language and a more relaxed style where appropriate."},
		"funny":      {"Inject humor, wit, or lightheartedness into the text (use with care).", "Additionally, try to inject appropriate and subtle humor or a lighthearted tone into the text. Make it engaging and amusing without undermining the core message, if applicable."},
		"persuasive": {"Make the text more convincing, confident, and impactful.", "Additionally, refine the text to be more persuasive and impactful. Strengthen arguments, use confident language, and aim to convince the reader."},
		"empathetic": {"Adjust the text to convey understanding, support, and compassion.", "Additionally, adjust the tone to be empathetic and supportive. Use language that conveys understanding and compassion, suitable for sensitive topics."},
		"concise":    {"Make the text as brief and to-the-point as possible, removing fluff.", "Additionally, ensure the text is extremely concise and to-the-point. Remove any redundant words or phrases and focus on conveying the core message with maximum brevity."},
	}
}

// createDefaultConfig generates a default configuration file at the specified path.
// This is typically called if no existing configuration file is found.
func createDefaultConfig(configPath string) error {
	configDir := filepath.Dir(configPath)
	// Ensure the configuration directory exists.
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil { // 0755 standard directory permissions.
			return fmt.Errorf("could not create config directory %s: %w", configDir, err)
		}
		printVerbose("Created config directory: %s", configDir)
	}

	// Populate the default configuration structure.
	defaultCfg := config.Config{
		DefaultLanguage: "Norwegian",
		Editor:          "nvim",
		GeminiModel:     "gemini-1.5-flash-latest",
		DefaultMood:     "neutral",
		Prompts:         defaultPromptsConfig, // Use the globally defined default prompts.
		Moods:           getDefaultMoods(),
	}

	// Marshal the default configuration to YAML.
	yamlData, err := yaml.Marshal(&defaultCfg)
	if err != nil {
		return fmt.Errorf("could not marshal default config to YAML: %w", err)
	}
	// Write the YAML data to the configuration file.
	if err := os.WriteFile(configPath, yamlData, 0644); err != nil { // 0644 standard file permissions.
		return fmt.Errorf("could not write default config file %s: %w", configPath, err)
	}

	// Provide feedback to the user.
	fmt.Printf("Created default config file: %s\n", configPath) // Always show this important message.
	printVerbose("You might want to review it. Key settings include 'defaultLanguage', 'editor', 'geminiModel', 'defaultMood'.")
	printVerbose("Run 'qik list-models' and 'qik list-moods' for available options.")
	printVerbose("For API key, use GEMINI_API_KEY env var or 'pass gemini_api_key'.")
	return nil
}

// initConfig is called by Cobra during initialization. It reads the configuration file
// (or creates a default one), binds environment variables, and unmarshals the
// configuration into the AppConfig struct. It also applies programmatic defaults
// if certain configuration values are missing.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag if provided.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search for config file in default locations.
		home, err := os.UserHomeDir()
		if err != nil {
			// This is an important operational warning, show even if not verbose.
			fmt.Fprintf(os.Stderr, "Warning: Could not get user home directory: %v. Will only check current directory for config.\n", err)
			viper.AddConfigPath(".") // Fallback to current directory.
		} else {
			primaryConfigDir := filepath.Join(home, ".config", "qik")
			viper.AddConfigPath(primaryConfigDir) // Preferred location.
			viper.AddConfigPath(".")              // Also check current directory.
		}
		viper.SetConfigName("config") // Name of config file (without extension).
		viper.SetConfigType("yaml")   // File type.
	}

	viper.AutomaticEnv()                               // Read matching environment variables.
	viper.SetEnvPrefix("QIK")                          // E.g., QIK_DEFAULTLANGUAGE. (Changed from SF)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // For nested keys like prompts.default -> QIK_PROMPTS_DEFAULT.

	// Attempt to read the configuration file.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; attempt to create a default one.
			fmt.Fprintln(os.Stderr, "Config file not found.") // Always inform user.
			defaultPath, pathErr := getDefaultConfigPath()
			if pathErr == nil {
				if createErr := createDefaultConfig(defaultPath); createErr != nil {
					fmt.Fprintf(os.Stderr, "Warning: Could not create default config: %v\n", createErr)
				} else {
					// If created successfully, try to read it again.
					if errReadAgain := viper.ReadInConfig(); errReadAgain != nil {
						fmt.Fprintf(os.Stderr, "Warning: Could not read newly created config: %v\n", errReadAgain)
					} else {
						printVerbose("Successfully read newly created config.")
					}
				}
			} else {
				fmt.Fprintf(os.Stderr, "Warning: Could not determine default config path: %v\n", pathErr)
			}
		} else {
			// Other errors while reading the config file.
			fmt.Fprintf(os.Stderr, "Error reading config file: %s\n", err)
		}
	} else {
		printVerbose("Successfully read config file: %s", viper.ConfigFileUsed())
	}

	// Unmarshal the loaded configuration (from file or env vars) into AppConfig.
	if err := viper.Unmarshal(&AppConfig); err != nil {
		// This is a critical error if the config is malformed.
		fmt.Fprintf(os.Stderr, "Unable to decode config into struct: %v. Please check your configuration file format.\n", err)
		// Depending on how critical a valid config is, you might os.Exit(1) here.
		// For now, we proceed and rely on programmatic defaults.
	}

	// Apply programmatic defaults if specific values are still missing after loading config.
	// This ensures the application has sensible fallbacks.
	if AppConfig.DefaultLanguage == "" {
		printVerbose("DefaultLanguage not set in config, using program default: Norwegian")
		AppConfig.DefaultLanguage = "Norwegian"
	}
	if AppConfig.Editor == "" {
		printVerbose("Editor not set in config, using program default: nvim")
		AppConfig.Editor = "nvim"
	}
	if AppConfig.GeminiModel == "" {
		printVerbose("GeminiModel not set in config, using program default: gemini-1.5-flash-latest")
		AppConfig.GeminiModel = "gemini-1.5-flash-latest"
	}
	if AppConfig.DefaultMood == "" {
		printVerbose("DefaultMood not set in config, using program default: neutral")
		AppConfig.DefaultMood = "neutral"
	}

	// Ensure prompt templates are populated, especially if loaded from an older config.
	if AppConfig.Prompts.Default == "" || !strings.Contains(AppConfig.Prompts.Default, "{MOOD_INSTRUCTION}") {
		printVerbose("Default prompt missing or outdated, setting to program default.")
		AppConfig.Prompts.Default = defaultPromptsConfig.Default
	}
	if AppConfig.Prompts.EnglishFixOnly == "" || !strings.Contains(AppConfig.Prompts.EnglishFixOnly, "{MOOD_INSTRUCTION}") {
		printVerbose("EnglishFixOnly prompt missing or outdated, setting to program default.")
		AppConfig.Prompts.EnglishFixOnly = defaultPromptsConfig.EnglishFixOnly
	}
	if AppConfig.Prompts.ExplainText == "" { // Also ensure ExplainText is present
		printVerbose("ExplainText prompt missing, setting to program default.")
		AppConfig.Prompts.ExplainText = defaultPromptsConfig.ExplainText
	}
	if AppConfig.Prompts.AnswerQuestion == "" { // And AnswerQuestion
		printVerbose("AnswerQuestion prompt missing, setting to program default.")
		AppConfig.Prompts.AnswerQuestion = defaultPromptsConfig.AnswerQuestion
	}

	// Ensure moods map is populated if missing.
	if AppConfig.Moods == nil || len(AppConfig.Moods) == 0 {
		printVerbose("Moods not set in config, using program defaults.")
		AppConfig.Moods = getDefaultMoods()
	}
}
