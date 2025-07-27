# RunComfy

A CLI tool for analyzing ComfyUI workflow files and managing missing dependencies on RunPod instances.

## Overview

RunComfy helps you analyze ComfyUI workflow JSON files to identify missing custom nodes and models, then provides guidance for installing them. It's designed to run directly on RunPod pods where ComfyUI is installed.

## Features

- 📁 **Workflow Analysis**: Parse ComfyUI workflow JSON files to extract dependencies
- 🔍 **Dependency Scanning**: Scan local ComfyUI installation for available nodes and models
- 📊 **Missing Dependencies**: Identify missing custom nodes and models required by workflows
- 💾 **Multiple Output Formats**: Support for table and JSON output formats
- 🚀 **RunPod Optimized**: Designed for standard RunPod ComfyUI installations

## Installation

### Option 1: Build from Source

```bash
git clone <repository-url>
cd runcomfy
go build -o runcomfy .
```

### Option 2: Download Binary

Download the latest binary from the releases page and make it executable:

```bash
chmod +x runcomfy
```

## Usage

### Basic Commands

#### Analyze a Workflow

Analyze a ComfyUI workflow file to identify missing dependencies:

```bash
# Basic analysis
./runcomfy analyze workflow.json

# Verbose output
./runcomfy analyze workflow.json --verbose

# JSON output
./runcomfy analyze workflow.json --output json

# Custom ComfyUI path
./runcomfy analyze workflow.json --comfyui-path /custom/path/ComfyUI
```

#### Scan ComfyUI Installation

Scan your local ComfyUI installation to see what's available:

```bash
# Basic scan
./runcomfy scan

# Verbose output with file sizes and paths
./runcomfy scan --verbose

# JSON output
./runcomfy scan --output json
```

#### Install Dependencies (Planning Mode)

Get installation guidance for missing dependencies:

```bash
# Show installation plan
./runcomfy install workflow.json

# Dry run (show what would be installed)
./runcomfy install workflow.json --dry-run
```

### Command Options

| Flag | Description | Default |
|------|-------------|---------|
| `--comfyui-path, -p` | Path to ComfyUI installation | `/workspace/ComfyUI` |
| `--output, -o` | Output format (table, json) | `table` |
| `--verbose, -v` | Verbose output | `false` |
| `--config` | Config file path | `$HOME/.runcomfy.yaml` |

### Configuration

Create a configuration file at `~/.runcomfy.yaml`:

```yaml
comfyui-path: "/workspace/ComfyUI"
output: "table"
verbose: false
```

## Example Output

### Workflow Analysis

```bash
$ ./runcomfy analyze my_workflow.json

📁 Workflow: my_workflow.json
📊 Summary: Missing: 2 missing custom nodes, 3 missing models

Statistics:
  Nodes:  15 total, 13 installed
  Models: 5 total, 2 installed

🔴 Missing Custom Nodes (2):
  - ComfyUI-Manager
  - ComfyUI-Custom-Scripts

🔴 Missing Models (3):
  Checkpoints:
    - sd_xl_base_1.0.safetensors
  Loras:
    - detail_tweaker_xl.safetensors
    - add_detail.safetensors

💡 Tip: Use 'runcomfy install my_workflow.json' to download missing dependencies.
```

### Installation Guidance

```bash
$ ./runcomfy install my_workflow.json

📦 Installation Plan:
Summary: Missing: 2 missing custom nodes, 3 missing models

🔌 Custom Nodes to Install (2):
  - ComfyUI-Manager
  - ComfyUI-Custom-Scripts

💡 To install custom nodes:
  cd /workspace/ComfyUI/custom_nodes
  # Use ComfyUI Manager or git clone the repositories

🎨 Models to Download (3):
  Checkpoints:
    - sd_xl_base_1.0.safetensors
      Target: /workspace/ComfyUI/models/checkpoints/sd_xl_base_1.0.safetensors
  Loras:
    - detail_tweaker_xl.safetensors
      Target: /workspace/ComfyUI/models/loras/detail_tweaker_xl.safetensors
    - add_detail.safetensors
      Target: /workspace/ComfyUI/models/loras/add_detail.safetensors

💡 To download models:
  1. Use ComfyUI Manager (recommended)
  2. Download manually from:
     - HuggingFace: https://huggingface.co/models
     - Civitai: https://civitai.com/
  3. Place files in the appropriate directories shown above
```

## RunPod Integration

This tool is optimized for RunPod environments where ComfyUI is typically installed at `/workspace/ComfyUI`. It scans the following standard directories:

- **Custom Nodes**: `/workspace/ComfyUI/custom_nodes/`
- **Checkpoints**: `/workspace/ComfyUI/models/checkpoints/`
- **LoRAs**: `/workspace/ComfyUI/models/loras/`
- **VAE**: `/workspace/ComfyUI/models/vae/`
- **ControlNet**: `/workspace/ComfyUI/models/controlnet/`
- **Upscale Models**: `/workspace/ComfyUI/models/upscale_models/`
- **Embeddings**: `/workspace/ComfyUI/models/embeddings/`

## Supported File Formats

### Workflow Files
- ComfyUI JSON workflows (`.json`)
- ComfyUI embedded workflows in PNG images (planned)

### Model Files
- Safetensors (`.safetensors`)
- Checkpoint files (`.ckpt`)
- PyTorch files (`.pt`, `.pth`)
- Binary files (`.bin`)

## Development

### Project Structure

```
runcomfy/
├── cmd/                    # CLI commands
│   ├── analyze.go         # Workflow analysis command
│   ├── install.go         # Installation guidance command
│   ├── root.go            # Root command and configuration
│   ├── scan.go            # Installation scanning command
│   └── version.go         # Version command
├── pkg/
│   ├── analyzer/          # Dependency analysis logic
│   ├── scanner/           # File system scanning
│   └── workflow/          # Workflow parsing
├── main.go                # Application entry point
├── go.mod                 # Go module definition
└── README.md              # This file
```

### Adding New Features

1. **Custom Node Detection**: Add new node types to `pkg/analyzer/analyzer.go`
2. **Model Categories**: Extend model path inference in `pkg/workflow/parser.go`
3. **Output Formats**: Add new formatters in the command files

### Building

```bash
# Build for current platform
go build -o runcomfy .

# Build for Linux (RunPod)
GOOS=linux GOARCH=amd64 go build -o runcomfy-linux .

# Run tests
go test ./...
```

## Roadmap

- [ ] Automatic model downloading from HuggingFace/Civitai
- [ ] ComfyUI Manager integration
- [ ] PNG workflow extraction
- [ ] Workflow validation
- [ ] Dependency resolution optimization
- [ ] Web interface
- [ ] Docker integration

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License - see LICENSE file for details.

## Support

For issues and feature requests, please use the GitHub issue tracker.