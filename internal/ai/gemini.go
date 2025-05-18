package ai

import (
	"context"
	"fmt"
	"log" // Used for warnings or unexpected API responses.
	"strings"

	"github.com/google/generative-ai-go/genai" // Official Google Gemini Go SDK.
	"google.golang.org/api/option"             // Used for API client options, like setting the API key.
)

// GeminiClient provides an interface for interacting with the Google Gemini API.
// It encapsulates a generative model client configured for a specific model.
type GeminiClient struct {
	model *genai.GenerativeModel
}

// NewGeminiClient initializes and returns a new GeminiClient.
// It requires a context, the API key, and the name of the Gemini model to use (e.g., "gemini-1.5-flash-latest").
// If modelName is empty, it defaults to "gemini-1.5-flash-latest" and logs a warning.
func NewGeminiClient(ctx context.Context, apiKey string, modelName string) (*GeminiClient, error) {
	// Create a new client with the provided API key.
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	effectiveModelName := modelName
	if effectiveModelName == "" {
		// This warning is for developers/maintainers if the model name isn't passed correctly from config.
		log.Println("Warning: Gemini model name was empty in NewGeminiClient, defaulting to 'gemini-1.5-flash-latest'.")
		effectiveModelName = "gemini-1.5-flash-latest"
	}

	// Get a specific generative model instance (e.g., for text generation).
	model := client.GenerativeModel(effectiveModelName)

	// Optional: Advanced model configuration for generation parameters or safety settings.
	// These can be uncommented and adjusted if specific behaviors are required.
	// For example, to control creativity (temperature) or response length (MaxOutputTokens).
	// Default safety settings are generally recommended unless specific needs arise.
	/*
		model.GenerationConfig = genai.GenerationConfig{
			// Temperature:     refFloat32(0.7), // Helper func needed: func refFloat32(f float32) *float32 { return &f }
			// TopP:            refFloat32(1.0),
			// TopK:            refInt32(40),    // Helper func needed: func refInt32(i int32) *int32 { return &i }
			MaxOutputTokens: 2048,
		}
		model.SafetySettings = []*genai.SafetySetting{
			{Category: genai.HarmCategoryHarassment, Threshold: genai.HarmBlockMediumAndAbove},
			{Category: genai.HarmCategoryHateSpeech, Threshold: genai.HarmBlockMediumAndAbove},
			// ... (configure other categories and thresholds as needed)
		}
	*/

	return &GeminiClient{model: model}, nil
}

// ProcessText sends the given text to the configured Gemini model for processing
// based on the provided prompt template and target language.
// The promptTemplate argument is expected to have placeholders {TEXT} and {LANGUAGE},
// which this function will replace. Other placeholders (e.g., {MOOD_INSTRUCTION})
// should be resolved by the caller before this function is invoked.
func (c *GeminiClient) ProcessText(ctx context.Context, textToProcess string, promptTemplate string, targetLanguage string) (string, error) {
	// Substitute placeholders in the prompt template with actual content.
	promptWithText := strings.ReplaceAll(promptTemplate, "{TEXT}", textToProcess)
	finalPrompt := strings.ReplaceAll(promptWithText, "{LANGUAGE}", targetLanguage)

	// For debugging: Uncomment to log the exact prompt being sent to the AI.
	// log.Printf("DEBUG: Sending prompt to Gemini:\n---\n%s\n---\n", finalPrompt)

	// Generate content using the Gemini model.
	resp, err := c.model.GenerateContent(ctx, genai.Text(finalPrompt))
	if err != nil {
		return "", fmt.Errorf("Gemini API call failed to generate content: %w", err)
	}

	// Basic validation of the API response.
	if resp == nil {
		return "", fmt.Errorf("Gemini API returned a nil response, which is unexpected")
	}

	// Check for any explicit blocking reasons from the API due to safety filters or other issues.
	if resp.PromptFeedback != nil {
		if resp.PromptFeedback.BlockReason != genai.BlockReasonUnspecified {
			return "", fmt.Errorf("content generation blocked by Gemini. Reason: %s. Review input or adjust safety settings if appropriate.", resp.PromptFeedback.BlockReason.String())
		}
		// Also check individual safety ratings if a block reason isn't specified but content might still be affected.
		for _, rating := range resp.PromptFeedback.SafetyRatings {
			if rating.Blocked {
				return "", fmt.Errorf("content generation blocked by Gemini due to safety rating. Category: %s, Probability: %s. Review input or adjust safety settings.", rating.Category, rating.Probability)
			}
		}
	}

	// Ensure the response contains valid candidates and content parts.
	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		// Log the full response if it's unexpectedly empty, for debugging.
		log.Printf("Warning: Gemini API returned no candidates or content parts. This might indicate an issue with the prompt, model configuration, or an unexpected safety filter. Full response: %+v", resp)
		return "", fmt.Errorf("AI returned no processable content. Please try rephrasing your input or check the model's status.")
	}

	// Extract the text content from the first candidate's first part.
	// This assumes the model's response for these tasks is primarily text.
	part := resp.Candidates[0].Content.Parts[0]
	if textPart, ok := part.(genai.Text); ok {
		return string(textPart), nil
	}

	// If the content part is not of the expected type.
	return "", fmt.Errorf("unexpected content part type from AI: %T. Expected genai.Text. Content: %+v", part, part)
}

// Helper functions for setting optional pointer fields in genai.GenerationConfig, if used.
// Example: func refFloat32(f float32) *float32 { return &f }
// Example: func refInt32(i int32) *int32 { return &i }
