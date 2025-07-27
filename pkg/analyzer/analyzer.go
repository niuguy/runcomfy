package analyzer

import (
	"fmt"
	"strings"

	"runcomfy/pkg/scanner"
	"runcomfy/pkg/workflow"
)

type Analyzer struct {
	installation *scanner.ComfyUIInstallation
}

func New(installation *scanner.ComfyUIInstallation) *Analyzer {
	return &Analyzer{
		installation: installation,
	}
}

func (a *Analyzer) AnalyzeWorkflow(w *workflow.Workflow) (*AnalysisResult, error) {
	result := &AnalysisResult{
		WorkflowPath: "",
		TotalNodes:   len(w.Nodes),
		TotalModels:  len(w.Models),
	}
	
	scanResult, err := a.installation.ScanInstallation()
	if err != nil {
		return nil, fmt.Errorf("failed to scan installation: %w", err)
	}
	result.InstalledNodes = len(scanResult.CustomNodes)
	result.InstalledModels = len(scanResult.Models)
	
	dependencies := w.ExtractDependencies()
	customNodes := w.GetCustomNodes()
	
	result.MissingNodes = a.findMissingNodes(customNodes, scanResult.CustomNodes)
	result.MissingModels = a.findMissingModels(dependencies, scanResult.Models)
	
	result.Summary = a.generateSummary(result)
	
	return result, nil
}

func (a *Analyzer) findMissingNodes(requiredNodes, installedNodes []string) []string {
	installedSet := make(map[string]bool)
	for _, node := range installedNodes {
		installedSet[node] = true
	}
	
	var missing []string
	for _, required := range requiredNodes {
		if !installedSet[required] && !isBuiltinNode(required) {
			missing = append(missing, required)
		}
	}
	
	return missing
}

func (a *Analyzer) findMissingModels(dependencies []workflow.Dependency, installedModels []scanner.FileInfo) []ModelDependency {
	installedSet := make(map[string]bool)
	for _, model := range installedModels {
		installedSet[model.Name] = true
		
		baseName := strings.TrimSuffix(model.Name, ".safetensors")
		baseName = strings.TrimSuffix(baseName, ".ckpt")
		baseName = strings.TrimSuffix(baseName, ".pt")
		baseName = strings.TrimSuffix(baseName, ".pth")
		baseName = strings.TrimSuffix(baseName, ".bin")
		installedSet[baseName] = true
	}
	
	var missing []ModelDependency
	seen := make(map[string]bool)
	
	for _, dep := range dependencies {
		if dep.Type == "model" && dep.Name != "" && !seen[dep.Name] {
			seen[dep.Name] = true
			
			baseName := strings.TrimSuffix(dep.Name, ".safetensors")
			baseName = strings.TrimSuffix(baseName, ".ckpt")
			baseName = strings.TrimSuffix(baseName, ".pt")
			baseName = strings.TrimSuffix(baseName, ".pth")
			baseName = strings.TrimSuffix(baseName, ".bin")
			
			if !installedSet[dep.Name] && !installedSet[baseName] {
				missing = append(missing, ModelDependency{
					Name:         dep.Name,
					Path:         dep.Path,
					Category:     inferModelCategory(dep.Path),
					Required:     true,
				})
			}
		}
	}
	
	return missing
}

func (a *Analyzer) generateSummary(result *AnalysisResult) string {
	var parts []string
	
	if len(result.MissingNodes) > 0 {
		parts = append(parts, fmt.Sprintf("%d missing custom nodes", len(result.MissingNodes)))
	}
	
	if len(result.MissingModels) > 0 {
		parts = append(parts, fmt.Sprintf("%d missing models", len(result.MissingModels)))
	}
	
	if len(parts) == 0 {
		return "All dependencies are satisfied ✓"
	}
	
	return "Missing: " + strings.Join(parts, ", ")
}

func isBuiltinNode(nodeType string) bool {
	builtinNodes := map[string]bool{
		"CheckpointLoaderSimple":  true,
		"CLIPTextEncode":          true,
		"KSampler":                true,
		"KSamplerAdvanced":        true,
		"VAEDecode":               true,
		"VAEEncode":               true,
		"SaveImage":               true,
		"LoadImage":               true,
		"EmptyLatentImage":        true,
		"LatentUpscale":           true,
		"LatentUpscaleBy":         true,
		"ControlNetLoader":        true,
		"ControlNetApply":         true,
		"LoraLoader":              true,
		"VAELoader":               true,
		"UpscaleModelLoader":      true,
		"ImageUpscaleWithModel":   true,
		"CLIPSetLastLayer":        true,
		"ConditioningCombine":     true,
		"ConditioningAverage":     true,
		"ConditioningConcat":      true,
		"ConditioningSetArea":     true,
		"ConditioningSetMask":     true,
		"ModelMergeSimple":        true,
		"ModelMergeBlocks":        true,
		"CheckpointLoader":        true,
		"DiffusionModelLoader":    true,
		"CLIPLoader":              true,
		"DualCLIPLoader":          true,
		"CLIPVisionLoader":        true,
		"StyleModelLoader":        true,
		"unCLIPCheckpointLoader":  true,
		"GLIGENLoader":            true,
		"GLIGENTextBoxApply":      true,
		"InpaintModelConditioning": true,
		"VideoLinearCFGGuidance":  true,
		"PatchModelAddDownscale":  true,
		"PatchModelAddUpscale":    true,
		"Rebatch":                 true,
		"RepeatLatentBatch":       true,
		"ImageBatch":              true,
		"ImagePadForOutpaint":     true,
		"ConditioningZeroOut":     true,
		"ConditioningSetTimestepRange": true,
		"SamplerCustom":           true,
		"BasicScheduler":          true,
		"KarrasScheduler":         true,
		"ExponentialScheduler":    true,
		"PolyexponentialScheduler": true,
		"SDTurboScheduler":        true,
		"VPScheduler":             true,
		"BetaScheduler":           true,
		"AlignYourStepsScheduler": true,
		"SplitSigmas":             true,
		"FlipSigmas":              true,
		"CFGGuider":               true,
		"DualCFGGuider":           true,
		"DisableNoise":            true,
		"RandomNoise":             true,
		"SamplerEulerAncestral":   true,
		"SamplerEuler":            true,
		"SamplerLMS":              true,
		"SamplerDPMPP_2M":         true,
		"SamplerDPMPP_2M_SDE":     true,
		"SamplerDPMPP_SDE":        true,
		"SamplerDPM_2":            true,
		"SamplerDPM_2_Ancestral":  true,
		"SamplerDPMAdaptive":      true,
		"SamplerDPMPP_3M_SDE":     true,
		"SamplerLCM":              true,
		"SamplerEulerAncestralCFGPP": true,
		"SamplerEulerCFGPP":       true,
		"SamplerDPMPP_2M_CFGpp":   true,
		"SamplerDPMPP_SDE_CFGpp":  true,
		"SamplerLCMUpscale":       true,
		"SamplerCustomAdvanced":   true,
		"ModelSamplingDiscrete":   true,
		"ModelSamplingContinuousEDM": true,
		"ModelSamplingContinuousV": true,
		"RescaleCFG":              true,
		"ImageScale":              true,
		"ImageScaleBy":            true,
		"ImageInvert":             true,
		"ImageBlend":              true,
		"ImageColorToMask":        true,
		"ImageCompositeMasked":    true,
		"MaskToImage":             true,
		"ImageToMask":             true,
		"MaskComposite":           true,
		"FeatherMask":             true,
		"GrowMask":                true,
		"ImageCrop":               true,
		"SetLatentNoiseMask":      true,
		"LatentComposite":         true,
		"LatentBlend":             true,
		"LatentRotate":            true,
		"LatentFlip":              true,
		"LatentCrop":              true,
		"PreviewImage":            true,
		"LoadImageMask":           true,
		"CropMask":                true,
		"MaskCombine":             true,
		"MaskErode":               true,
		"SolidMask":               true,
		"InvertMask":              true,
		"CropImage":               true,
		"LoadVideo":               true,
		"SaveAnimatedWEBP":        true,
		"SaveAnimatedPNG":         true,
		"ImageOnlyCheckpointLoader": true,
		"SVD_img2vid_Conditioning": true,
		"VideoTriangleCFGGuidance": true,
		"StableVideoDiffusion_Decode": true,
		"ImageOnlyCheckpointSave": true,
		"unCLIPConditioning":      true,
		"PerpNeg":                 true,
		"PerpNegGuider":           true,
		"ModelMergeSDXL":          true,
		"ModelMergeSD1":           true,
		"ModelMergeSD2":           true,
		"CLIPMergeSimple":         true,
		"CLIPMergeAdd":            true,
		"CLIPMergeSubtract":       true,
		"CLIPSubtract":            true,
		"CLIPSave":                true,
		"VAESave":                 true,
		"CheckpointSave":          true,
		"DiffusionModelSave":      true,
		"UNETLoader":              true,
		"CLIPTextEncodeSDXL":      true,
		"CLIPTextEncodeSDXLRefiner": true,
		"CLIPTextEncodeFlux":      true,
		"UnetLoaderGGUF":          true,
		"DualCLIPLoaderGGUF":      true,
		"TripleCLIPLoader":        true,
		"CLIPTextEncodeSD3":       true,
		"EmptySD3LatentImage":     true,
		"KSamplerSelect":          true,
		"SamplerDDIM":             true,
		"SamplerDDPM":             true,
		"SamplerHeun":             true,
		"SamplerHeunpp2":          true,
		"SamplerUniPC":            true,
		"SamplerIPNDM":            true,
		"SamplerIPNDMV":           true,
		"SamplerDEIS":             true,
		"PhotoMakerLoader":        true,
		"PhotoMakerEncode":        true,
	}
	
	return builtinNodes[nodeType]
}

func inferModelCategory(path string) string {
	path = strings.ToLower(path)
	switch {
	case strings.Contains(path, "checkpoint"):
		return "checkpoints"
	case strings.Contains(path, "lora"):
		return "loras"
	case strings.Contains(path, "vae"):
		return "vae"
	case strings.Contains(path, "controlnet"):
		return "controlnet"
	case strings.Contains(path, "upscale"):
		return "upscale_models"
	case strings.Contains(path, "embedding"):
		return "embeddings"
	default:
		return "models"
	}
}