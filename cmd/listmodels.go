package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listModelsCmd represents the command that displays information about available Gemini models.
var listModelsCmd = &cobra.Command{
	Use:   "list-models",
	Short: "Lists common Gemini models with brief descriptions.",
	Long: `Prints a list of commonly used Gemini models suitable for text processing tasks,
along with a short explanation of their typical use cases and strengths.
This list is curated and intended as a helpful starting point; it may not be exhaustive.
For the most up-to-date information, always refer to official Google Gemini documentation.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Commonly available Gemini models for text processing tasks:")
		fmt.Println("The model names ending with '-latest' generally point to the most recent stable version of that model series.")
		fmt.Println("You might also be able to use specific versioned model names (e.g., 'gemini-1.5-flash-001').")
		fmt.Println("------------------------------------------------------------------------------------")

		// Model: gemini-1.5-flash-latest
		fmt.Printf("\nModel Name: gemini-1.5-flash-latest\n")
		fmt.Println("  Generation:  1.5")
		fmt.Println("  Description: The fastest and most cost-effective model in the 1.5 generation,")
		fmt.Println("               optimized for speed and efficiency with strong multimodal capabilities and a large context window.")
		fmt.Println("  Strengths:   High speed, lower cost, 1M token context window, multimodal understanding (text, image, audio, video).")
		fmt.Println("               Excellent for summarization, chat, captioning, data extraction, and high-throughput tasks.")
		fmt.Println("  Use Case:    Recommended for this tool (qik) for most tasks, especially quick corrections, explanations, and answers")
		fmt.Println("               where a good balance of performance, capability, and cost is desired.")

		// Model: gemini-1.5-pro-latest
		fmt.Printf("\nModel Name: gemini-1.5-pro-latest\n")
		fmt.Println("  Generation:  1.5")
		fmt.Println("  Description: Google's most capable 1.5 generation model, designed for highly complex reasoning and multimodal tasks.")
		fmt.Println("  Strengths:   State-of-the-art performance, advanced reasoning, superior quality on complex tasks, very large context window (1M tokens),")
		fmt.Println("               strong multimodal capabilities (text, image, audio, video analysis).")
		fmt.Println("  Use Case:    When maximum quality or complex reasoning is paramount. Suitable for in-depth analysis, creative content generation,")
		fmt.Println("               or if 'gemini-1.5-flash-latest' doesn't provide sufficient quality for a specific demanding task.")
		fmt.Println("               Typically has higher latency and cost compared to 'flash' models.")

		// Model: gemini-pro (often an alias for gemini-1.0-pro)
		fmt.Printf("\nModel Name: gemini-pro (typically refers to the latest stable gemini-1.0-pro version)\n")
		fmt.Println("  Generation:  1.0")
		fmt.Println("  Description: A capable text-optimized model from the previous (1.0) generation.")
		fmt.Println("  Strengths:   Solid performance for a wide range of natural language tasks, including text generation and multi-turn chat.")
		fmt.Println("               Has a smaller context window compared to 1.5 models.")
		fmt.Println("  Use Case:    Can be an alternative if 1.5 models are not available or if specific behavior of the 1.0 generation is needed.")
		fmt.Println("               Generally, 'gemini-1.5-flash-latest' is recommended over this for new use cases due to improvements in")
		fmt.Println("               cost, speed, context window, and overall capabilities.")

		fmt.Println("\n------------------------------------------------------------------------------------")
		fmt.Println("Important Notes:")
		fmt.Println("  - Model availability can depend on your Google Cloud project, region, and API access.")
		fmt.Println("  - For the absolute latest list, specific version identifiers, and detailed capabilities,")
		fmt.Println("    always consult the official Google Cloud Vertex AI or Google AI Studio documentation.")
		fmt.Println("  - You can set your preferred model in the qik configuration file")
		fmt.Println("    (e.g., ~/.config/qik/config.yaml) under the 'geminiModel' key.")
		fmt.Println("    This tool will use whatever valid model identifier string you provide there.")
	},
}

func init() {
	rootCmd.AddCommand(listModelsCmd)
}
