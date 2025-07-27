package workflow

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func ParseWorkflow(filePath string) (*Workflow, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read workflow file: %w", err)
	}

	var workflow Workflow
	if err := json.Unmarshal(data, &workflow); err != nil {
		return nil, fmt.Errorf("failed to parse workflow JSON: %w", err)
	}

	return &workflow, nil
}

func (w *Workflow) ExtractDependencies() []Dependency {
	var deps []Dependency
	nodeTypes := make(map[string]bool)

	for _, node := range w.Nodes {
		if !nodeTypes[node.Type] {
			nodeTypes[node.Type] = true
			deps = append(deps, Dependency{
				Type: "node",
				Name: node.Type,
				Path: "",
			})
		}

		if node.Properties != nil {
			if model, ok := node.Properties["model_name"].(string); ok && model != "" {
				deps = append(deps, Dependency{
					Type: "model",
					Name: model,
					Path: inferModelPath(node.Type, model),
				})
			}
			if ckpt, ok := node.Properties["ckpt_name"].(string); ok && ckpt != "" {
				deps = append(deps, Dependency{
					Type: "model",
					Name: ckpt,
					Path: "checkpoints/" + ckpt,
				})
			}
			if lora, ok := node.Properties["lora_name"].(string); ok && lora != "" {
				deps = append(deps, Dependency{
					Type: "model",
					Name: lora,
					Path: "loras/" + lora,
				})
			}
		}

		if len(node.Widgets) > 0 {
			for _, widget := range node.Widgets {
				if str, ok := widget.(string); ok {
					if isModelFile(str) {
						deps = append(deps, Dependency{
							Type: "model",
							Name: str,
							Path: inferModelPath(node.Type, str),
						})
					}
				}
			}
		}
	}

	for _, model := range w.Models {
		deps = append(deps, Dependency{
			Type: "model",
			Name: model.Name,
			Path: model.Directory + "/" + model.Name,
		})
	}

	return deps
}

func (w *Workflow) GetCustomNodes() []string {
	var customNodes []string
	builtinNodes := map[string]bool{
		"CheckpointLoaderSimple": true,
		"CLIPTextEncode":         true,
		"KSampler":               true,
		"VAEDecode":              true,
		"SaveImage":              true,
		"LoadImage":              true,
		"EmptyLatentImage":       true,
	}

	for _, node := range w.Nodes {
		if !builtinNodes[node.Type] {
			customNodes = append(customNodes, node.Type)
		}
	}

	return removeDuplicates(customNodes)
}

func inferModelPath(nodeType, modelName string) string {
	switch {
	case strings.Contains(strings.ToLower(nodeType), "checkpoint"):
		return "checkpoints/" + modelName
	case strings.Contains(strings.ToLower(nodeType), "lora"):
		return "loras/" + modelName
	case strings.Contains(strings.ToLower(nodeType), "vae"):
		return "vae/" + modelName
	case strings.Contains(strings.ToLower(nodeType), "controlnet"):
		return "controlnet/" + modelName
	case strings.Contains(strings.ToLower(nodeType), "upscale"):
		return "upscale_models/" + modelName
	default:
		return "models/" + modelName
	}
}

func isModelFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	modelExts := map[string]bool{
		".safetensors": true,
		".ckpt":        true,
		".pt":          true,
		".pth":         true,
		".bin":         true,
	}
	return modelExts[ext]
}

func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	var result []string
	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}
	return result
}