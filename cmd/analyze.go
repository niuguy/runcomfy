package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"runcomfy/pkg/analyzer"
	"runcomfy/pkg/scanner"
	"runcomfy/pkg/workflow"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze <workflow.json>",
	Short: "Analyze a ComfyUI workflow for missing dependencies",
	Long: `Analyze a ComfyUI workflow file to identify missing custom nodes and models.

This command scans your local ComfyUI installation and compares it with
the requirements from the workflow file to show what's missing.`,
	Args: cobra.ExactArgs(1),
	RunE: runAnalyze,
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	workflowPath := args[0]
	comfyUIPath := viper.GetString("comfyui-path")
	outputFormat := viper.GetString("output")
	verbose := viper.GetBool("verbose")

	if verbose {
		fmt.Printf("Analyzing workflow: %s\n", workflowPath)
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

	result.WorkflowPath = workflowPath

	switch outputFormat {
	case "json":
		return outputJSON(result)
	case "table":
		return outputTable(result, verbose)
	default:
		return fmt.Errorf("unsupported output format: %s", outputFormat)
	}
}

func outputJSON(result *analyzer.AnalysisResult) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

func outputTable(result *analyzer.AnalysisResult, verbose bool) error {
	fmt.Printf("📁 Workflow: %s\n", filepath.Base(result.WorkflowPath))
	fmt.Printf("📊 Summary: %s\n\n", result.Summary)

	fmt.Printf("Statistics:\n")
	fmt.Printf("  Nodes:  %d total, %d installed\n", result.TotalNodes, result.InstalledNodes)
	fmt.Printf("  Models: %d total, %d installed\n\n", result.TotalModels, result.InstalledModels)

	if len(result.MissingNodes) > 0 {
		fmt.Printf("🔴 Missing Custom Nodes (%d):\n", len(result.MissingNodes))
		for _, node := range result.MissingNodes {
			fmt.Printf("  - %s\n", node)
		}
		fmt.Println()
	}

	if len(result.MissingModels) > 0 {
		fmt.Printf("🔴 Missing Models (%d):\n", len(result.MissingModels))
		
		categories := make(map[string][]analyzer.ModelDependency)
		for _, model := range result.MissingModels {
			categories[model.Category] = append(categories[model.Category], model)
		}

		for category, models := range categories {
			fmt.Printf("  %s:\n", strings.Title(category))
			for _, model := range models {
				if verbose {
					fmt.Printf("    - %s (path: %s)\n", model.Name, model.Path)
				} else {
					fmt.Printf("    - %s\n", model.Name)
				}
			}
		}
		fmt.Println()
	}

	if len(result.MissingNodes) == 0 && len(result.MissingModels) == 0 {
		fmt.Println("✅ All dependencies are satisfied! You can run this workflow.")
	} else {
		fmt.Println("💡 Tip: Use 'runcomfy install <workflow.json>' to download missing dependencies.")
	}

	return nil
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
}