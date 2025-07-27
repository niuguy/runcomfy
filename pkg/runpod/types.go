package runpod

import "time"

type Pod struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Status          string    `json:"status"`
	Runtime         Runtime   `json:"runtime,omitempty"`
	Machine         Machine   `json:"machine,omitempty"`
	DesiredStatus   string    `json:"desiredStatus,omitempty"`
	ImageName       string    `json:"imageName,omitempty"`
	ContainerDiskInGb int     `json:"containerDiskInGb,omitempty"`
	VolumeInGb      int       `json:"volumeInGb,omitempty"`
	VolumeMountPath string    `json:"volumeMountPath,omitempty"`
	CreatedAt       time.Time `json:"createdAt,omitempty"`
	LastStatusChange time.Time `json:"lastStatusChange,omitempty"`
}

type Runtime struct {
	UptimeInSeconds int       `json:"uptimeInSeconds,omitempty"`
	Ports           []Port    `json:"ports,omitempty"`
	GPUs            []GPU     `json:"gpus,omitempty"`
	CPUs            int       `json:"cpus,omitempty"`
	MemoryInGb      int       `json:"memoryInGb,omitempty"`
	LastStartedAt   time.Time `json:"lastStartedAt,omitempty"`
}

type Port struct {
	IP           string `json:"ip,omitempty"`
	PrivatePort  int    `json:"privatePort,omitempty"`
	PublicPort   int    `json:"publicPort,omitempty"`
	Type         string `json:"type,omitempty"`
}

type GPU struct {
	ID               string `json:"id,omitempty"`
	GPUTypeID        string `json:"gpuTypeId,omitempty"`
	GPUTypeName      string `json:"gpuTypeName,omitempty"`
	MemoryInGb       int    `json:"memoryInGb,omitempty"`
}

type Machine struct {
	PodHostId string `json:"podHostId,omitempty"`
}

type VolumeInfo struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	DataCenterId    string    `json:"dataCenterId,omitempty"`
	Size            int       `json:"size,omitempty"`
	CreatedAt       time.Time `json:"createdAt,omitempty"`
}

type FileInfo struct {
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	IsDir     bool      `json:"isDir"`
	ModTime   time.Time `json:"modTime"`
	FileType  string    `json:"fileType,omitempty"`
}

type ScanResult struct {
	CustomNodes []string   `json:"customNodes"`
	Models      []FileInfo `json:"models"`
	TotalFiles  int        `json:"totalFiles"`
	ScanTime    time.Time  `json:"scanTime"`
}