package analyzer

type AnalysisResult struct {
	WorkflowPath     string            `json:"workflowPath"`
	TotalNodes       int               `json:"totalNodes"`
	TotalModels      int               `json:"totalModels"`
	InstalledNodes   int               `json:"installedNodes"`
	InstalledModels  int               `json:"installedModels"`
	MissingNodes     []string          `json:"missingNodes"`
	MissingModels    []ModelDependency `json:"missingModels"`
	Summary          string            `json:"summary"`
}

type ModelDependency struct {
	Name         string `json:"name"`
	Path         string `json:"path"`
	Category     string `json:"category"`
	Required     bool   `json:"required"`
	DownloadURL  string `json:"downloadUrl,omitempty"`
	Size         int64  `json:"size,omitempty"`
}