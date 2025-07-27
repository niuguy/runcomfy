package scanner

import "time"

type ComfyUIInstallation struct {
	BasePath     string
	CustomNodes  string
	ModelsPath   string
	CheckPoints  string
	Loras        string
	VAE          string
	ControlNet   string
	Upscale      string
	Embeddings   string
}

type FileInfo struct {
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	Size     int64     `json:"size"`
	IsDir    bool      `json:"isDir"`
	ModTime  time.Time `json:"modTime"`
	FileType string    `json:"fileType"`
}

type ScanResult struct {
	CustomNodes []string   `json:"customNodes"`
	Models      []FileInfo `json:"models"`
	TotalFiles  int        `json:"totalFiles"`
	ScanTime    time.Time  `json:"scanTime"`
	BasePath    string     `json:"basePath"`
}

type MissingDependencies struct {
	CustomNodes []string   `json:"missingCustomNodes"`
	Models      []string   `json:"missingModels"`
	Summary     string     `json:"summary"`
}