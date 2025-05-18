package cmd

import (
	"fmt"
	"sort" // Used to ensure consistent output order of moods.

	"github.com/spf13/cobra"
	// No need to import "qik/internal/config" directly,
	// as AppConfig (which holds config.Config) is a package-level variable in cmd/root.go.
)

// listMoodsCmd represents the command that displays available moods/tones from the configuration.
var listMoodsCmd = &cobra.Command{
	Use:   "list-moods",
	Short: "Lists available text moods/tones with descriptions.",
	Long: `Prints a list of moods/tones that can be applied to the text when using
commands like 'fix' or 'answer'. Each mood is listed with a short description
of how it influences the AI's output.
These moods are defined by the user in their qik configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if any moods are loaded from the configuration.
		// AppConfig is populated by initConfig in root.go.
		if len(AppConfig.Moods) == 0 {
			fmt.Println("No moods are currently defined in your configuration.")
			fmt.Println("You can define moods in the 'moods' section of your qik config file")
			fmt.Println("(e.g., ~/.config/qik/config.yaml) to customize text tones.")
			return
		}

		fmt.Println("Available text moods/tones (from your configuration):")
		fmt.Println("----------------------------------------------------")

		// Sort mood keys alphabetically to ensure a consistent and predictable output order.
		keys := make([]string, 0, len(AppConfig.Moods))
		for k := range AppConfig.Moods {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		// Iterate over the sorted keys and print details for each mood.
		for _, key := range keys {
			mood := AppConfig.Moods[key] // mood is of type config.MoodInstruction
			fmt.Printf("\nMood Key:    %s\n", key)
			fmt.Printf("  Description: %s\n", mood.Description)
			// The 'mood.Instruction' is for the AI prompt and generally not shown to the user here.
			// It could be shown in verbose mode if detailed debugging is needed in the future.
			// printVerbose("  Instruction (for prompt): %s\n", mood.Instruction)
		}

		fmt.Println("\n----------------------------------------------------")
		fmt.Printf("To use a mood, specify its 'Mood Key' with the --mood flag\n")
		fmt.Printf("(e.g., qik fix --mood professional, or qik answer -m casual).\n")
		fmt.Printf("The current default mood (used if --mood is not specified) is: '%s'.\n", AppConfig.DefaultMood)
		fmt.Println("You can change the default mood in your qik configuration file.")
	},
}

func init() {
	rootCmd.AddCommand(listMoodsCmd)
}
