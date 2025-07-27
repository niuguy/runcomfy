package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"runcomfy/pkg/analyzer"
	"runcomfy/pkg/scanner"
	"runcomfy/pkg/workflow"
)

var installCmd = &cobra.Command{
	Use:   "install <workflow.json>",
	Short: "Install missing dependencies for a ComfyUI workflow",
	Long: `Install missing custom nodes and models required by a ComfyUI workflow.

This command analyzes the workflow and provides instructions for installing
missing dependencies.`,
	Args: cobra.ExactArgs(1),
	RunE: runInstall,
}

var (
	dryRun    bool
	autoYes   bool
)

func runInstall(cmd *cobra.Command, args []string) error {
	workflowPath := args[0]
	comfyUIPath := viper.GetString("comfyui-path")
	verbose := viper.GetBool("verbose")

	if verbose {
		fmt.Printf("Installing dependencies for workflow: %s\n", workflowPath)
		fmt.Printf("ComfyUI path: %s\n", comfyUIPath)
	}

	if _, err := os.Stat(workflowPath); os.IsNotExist(err) {
		return fmt.Errorf("workflow file not found: %s", workflowPath)
	}

	w, err := workflow.ParseWorkflow(workflowPath)
	if err != nil {
		return fmt.Errorf("failed to parse workflow: %w", err)
	}

	installation := scanner.NewComfyUIInstallation(comfyUIPath)
	
	if _, err := os.Stat(installation.BasePath); os.IsNotExist(err) {
		return fmt.Errorf("ComfyUI installation not found at: %s", installation.BasePath)
	}

	a := analyzer.New(installation)
	result, err := a.AnalyzeWorkflow(w)
	if err != nil {
		return fmt.Errorf("analysis failed: %w", err)
	}

	if len(result.MissingNodes) == 0 && len(result.MissingModels) == 0 {
		fmt.Println("✅ All dependencies are already satisfied!")
		return nil
	}

	fmt.Printf("📦 Installation Plan:\n")
	fmt.Printf("Summary: %s\n\n", result.Summary)

	if dryRun {
		fmt.Println("🔍 Dry run mode - showing what would be installed:")
	}

	if len(result.MissingNodes) > 0 {
		fmt.Printf("🔌 Custom Nodes to Install (%d):\n", len(result.MissingNodes))
		for _, node := range result.MissingNodes {
			fmt.Printf("  - %s\n", node)
		}
		fmt.Println("\n💡 To install custom nodes:")
		fmt.Printf("  cd %s/custom_nodes\n", comfyUIPath)
		fmt.Println("  # Use ComfyUI Manager or git clone the repositories")
		fmt.Println()
	}

	if len(result.MissingModels) > 0 {
		fmt.Printf("🎨 Models to Download (%d):\n", len(result.MissingModels))
		
		categories := make(map[string][]analyzer.ModelDependency)
		for _, model := range result.MissingModels {
			categories[model.Category] = append(categories[model.Category], model)
		}

		for category, models := range categories {
			fmt.Printf("  %s:\n", strings.Title(category))
			for _, model := range models {
				fmt.Printf("    - %s\n", model.Name)
				if model.Path != "" {
					fmt.Printf("      Target: %s/%s\n", comfyUIPath, model.Path)
				}
			}
		}
		
		fmt.Println("\n💡 To download models:")
		fmt.Println("  1. Use ComfyUI Manager (recommended)")
		fmt.Println("  2. Download manually from:")
		fmt.Println("     - HuggingFace: https://huggingface.co/models")
		fmt.Println("     - Civitai: https://civitai.com/")
		fmt.Println("  3. Place files in the appropriate directories shown above")
		fmt.Println()
	}

	if !dryRun {
		fmt.Println("⚠️  Automatic installation not yet implemented.")
		fmt.Println("   Please install dependencies manually using the guidance above.")
		fmt.Println("   Future versions will support automatic downloads.")
	}

	return nil
}

func init() {
	installCmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would be installed without actually installing")
	installCmd.Flags().BoolVar(&autoYes, "yes", false, "automatically answer yes to all prompts")
	
	rootCmd.AddCommand(installCmd)
}