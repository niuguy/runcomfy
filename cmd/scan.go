package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"runcomfy/pkg/scanner"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan the ComfyUI installation for available nodes and models",
	Long: `Scan your local ComfyUI installation to see what custom nodes 
and models are currently available.`,
	RunE: runScan,
}

func runScan(cmd *cobra.Command, args []string) error {
	comfyUIPath := viper.GetString("comfyui-path")
	outputFormat := viper.GetString("output")
	verbose := viper.GetBool("verbose")

	if verbose {
		fmt.Printf("Scanning ComfyUI installation: %s\n", comfyUIPath)
	}

	installation := scanner.NewComfyUIInstallation(comfyUIPath)
	
	if _, err := os.Stat(installation.BasePath); os.IsNotExist(err) {
		return fmt.Errorf("ComfyUI installation not found at: %s", installation.BasePath)
	}

	result, err := installation.ScanInstallation()
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	switch outputFormat {
	case "json":
		return outputScanJSON(result)
	case "table":
		return outputScanTable(result, verbose)
	default:
		return fmt.Errorf("unsupported output format: %s", outputFormat)
	}
}

func outputScanJSON(result *scanner.ScanResult) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

func outputScanTable(result *scanner.ScanResult, verbose bool) error {
	fmt.Printf("📁 ComfyUI Installation: %s\n", result.BasePath)
	fmt.Printf("🕐 Scan Time: %s\n\n", result.ScanTime.Format("2006-01-02 15:04:05"))

	fmt.Printf("📊 Summary:\n")
	fmt.Printf("  Custom Nodes: %d\n", len(result.CustomNodes))
	fmt.Printf("  Models: %d\n", len(result.Models))
	fmt.Printf("  Total Files: %d\n\n", result.TotalFiles)

	if len(result.CustomNodes) > 0 {
		fmt.Printf("🔌 Custom Nodes (%d):\n", len(result.CustomNodes))
		for _, node := range result.CustomNodes {
			fmt.Printf("  - %s\n", node)
		}
		fmt.Println()
	}

	if len(result.Models) > 0 {
		fmt.Printf("🎨 Models (%d):\n", len(result.Models))
		
		categories := make(map[string][]scanner.FileInfo)
		for _, model := range result.Models {
			categories[model.FileType] = append(categories[model.FileType], model)
		}

		for category, models := range categories {
			fmt.Printf("  %s (%d):\n", strings.Title(category), len(models))
			for _, model := range models {
				if verbose {
					fmt.Printf("    - %s (%.2f MB, %s)\n", 
						model.Name, 
						float64(model.Size)/(1024*1024),
						model.Path)
				} else {
					fmt.Printf("    - %s\n", model.Name)
				}
			}
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(scanCmd)
}