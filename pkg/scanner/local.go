package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func NewComfyUIInstallation(basePath string) *ComfyUIInstallation {
	if basePath == "" {
		basePath = "/workspace/ComfyUI"
		
		if _, err := os.Stat(basePath); os.IsNotExist(err) {
			basePath = "./ComfyUI"
		}
		
		if _, err := os.Stat(basePath); os.IsNotExist(err) {
			basePath = "."
		}
	}
	
	return &ComfyUIInstallation{
		BasePath:    basePath,
		CustomNodes: filepath.Join(basePath, "custom_nodes"),
		ModelsPath:  filepath.Join(basePath, "models"),
		CheckPoints: filepath.Join(basePath, "models", "checkpoints"),
		Loras:       filepath.Join(basePath, "models", "loras"),
		VAE:         filepath.Join(basePath, "models", "vae"),
		ControlNet:  filepath.Join(basePath, "models", "controlnet"),
		Upscale:     filepath.Join(basePath, "models", "upscale_models"),
		Embeddings:  filepath.Join(basePath, "models", "embeddings"),
	}
}

func (c *ComfyUIInstallation) ScanInstallation() (*ScanResult, error) {
	result := &ScanResult{
		ScanTime: time.Now(),
		BasePath: c.BasePath,
	}
	
	customNodes, err := c.scanCustomNodes()
	if err != nil {
		return nil, fmt.Errorf("failed to scan custom nodes: %w", err)
	}
	result.CustomNodes = customNodes
	
	models, err := c.scanModels()
	if err != nil {
		return nil, fmt.Errorf("failed to scan models: %w", err)
	}
	result.Models = models
	result.TotalFiles = len(models)
	
	return result, nil
}

func (c *ComfyUIInstallation) scanCustomNodes() ([]string, error) {
	var nodes []string
	
	if _, err := os.Stat(c.CustomNodes); os.IsNotExist(err) {
		return nodes, nil
	}
	
	entries, err := os.ReadDir(c.CustomNodes)
	if err != nil {
		return nil, fmt.Errorf("failed to read custom_nodes directory: %w", err)
	}
	
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			nodes = append(nodes, entry.Name())
		}
	}
	
	return nodes, nil
}

func (c *ComfyUIInstallation) scanModels() ([]FileInfo, error) {
	var models []FileInfo
	
	modelDirs := map[string]string{
		"checkpoints":    c.CheckPoints,
		"loras":         c.Loras,
		"vae":           c.VAE,
		"controlnet":    c.ControlNet,
		"upscale_models": c.Upscale,
		"embeddings":    c.Embeddings,
	}
	
	for category, dirPath := range modelDirs {
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			continue
		}
		
		err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			
			if !info.IsDir() && isModelFile(info.Name()) {
				relPath, _ := filepath.Rel(c.BasePath, path)
				models = append(models, FileInfo{
					Name:     info.Name(),
					Path:     relPath,
					Size:     info.Size(),
					IsDir:    false,
					ModTime:  info.ModTime(),
					FileType: category,
				})
			}
			
			return nil
		})
		
		if err != nil {
			return nil, fmt.Errorf("failed to walk %s directory: %w", category, err)
		}
	}
	
	return models, nil
}

func (c *ComfyUIInstallation) HasCustomNode(nodeName string) bool {
	nodePath := filepath.Join(c.CustomNodes, nodeName)
	_, err := os.Stat(nodePath)
	return !os.IsNotExist(err)
}

func (c *ComfyUIInstallation) HasModel(modelName string) bool {
	modelDirs := []string{
		c.CheckPoints,
		c.Loras,
		c.VAE,
		c.ControlNet,
		c.Upscale,
		c.Embeddings,
	}
	
	for _, dir := range modelDirs {
		modelPath := filepath.Join(dir, modelName)
		if _, err := os.Stat(modelPath); !os.IsNotExist(err) {
			return true
		}
	}
	
	return false
}

func (c *ComfyUIInstallation) GetModelPath(modelName string) (string, bool) {
	modelDirs := []string{
		c.CheckPoints,
		c.Loras,
		c.VAE,
		c.ControlNet,
		c.Upscale,
		c.Embeddings,
	}
	
	for _, dirPath := range modelDirs {
		modelPath := filepath.Join(dirPath, modelName)
		if _, err := os.Stat(modelPath); !os.IsNotExist(err) {
			relPath, _ := filepath.Rel(c.BasePath, modelPath)
			return relPath, true
		}
	}
	
	return "", false
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