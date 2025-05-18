package config

// MoodInstruction defines the structure for a single mood/tone configuration.
// It includes a user-facing description and the instruction text for the AI.
type MoodInstruction struct {
	// Description is a human-readable explanation of the mood,
	// used by commands like 'qik list-moods'.
	Description string `mapstructure:"description"`

	// Instruction is the specific text appended to an AI prompt
	// to guide the model towards generating output with the desired mood.
	Instruction string `mapstructure:"instruction"`
}

// Prompts defines the structure for storing various AI prompt templates
// used by different commands in the application. Each field corresponds
// to a specific task (e.g., fixing text, explaining text).
type Prompts struct {
	// Default is the general-purpose prompt for text correction and improvement.
	Default string `mapstructure:"default"`

	// EnglishFixOnly is a specialized prompt for correcting English text without translation.
	EnglishFixOnly string `mapstructure:"english_fix_only"`

	// ExplainText is the prompt used to generate explanations of provided text.
	ExplainText string `mapstructure:"explain_text"`

	// AnswerQuestion is the prompt used for generating answers to user questions.
	AnswerQuestion string `mapstructure:"answer_question"`
}

// Config is the main structure holding all application configuration settings.
// These settings are typically loaded from a YAML file (e.g., config.yaml)
// and can be overridden by environment variables.
type Config struct {
	// DefaultLanguage specifies the default target language for AI processing
	// (e.g., "Norwegian", "English") if not overridden by a command-line flag.
	DefaultLanguage string `mapstructure:"defaultLanguage"`

	// Editor defines the command-line editor to be used for text input
	// (e.g., "nvim", "vim", "nano").
	Editor string `mapstructure:"editor"`

	// GeminiAPIKey can store the Gemini API key directly in the configuration.
	// However, using environment variables (GEMINI_API_KEY) or 'pass' is recommended for security.
	GeminiAPIKey string `mapstructure:"geminiApiKey"`

	// GeminiModel specifies the particular Gemini AI model to be used for processing
	// (e.g., "gemini-1.5-flash-latest").
	GeminiModel string `mapstructure:"geminiModel"`

	// DefaultMood is the key (from the Moods map) of the mood/tone to be applied
	// by default if no specific mood is requested via a command-line flag.
	DefaultMood string `mapstructure:"defaultMood"`

	// Moods is a map where keys are mood identifiers (e.g., "professional", "casual")
	// and values are MoodInstruction structs defining the mood's description and AI instruction.
	Moods map[string]MoodInstruction `mapstructure:"moods"`

	// Prompts contains the various AI prompt templates used by the application.
	Prompts Prompts `mapstructure:"prompts"`
}
