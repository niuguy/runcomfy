package workflow

type Workflow struct {
	Version string                 `json:"version"`
	Config  Config                 `json:"config,omitempty"`
	State   State                  `json:"state,omitempty"`
	Groups  []Group                `json:"groups,omitempty"`
	Nodes   map[string]Node        `json:"nodes"`
	Links   []Link                 `json:"links,omitempty"`
	Models  []Model                `json:"models,omitempty"`
	Extra   map[string]interface{} `json:"extra,omitempty"`
}

type Config struct {
	LinksOnTop   bool `json:"links_ontop,omitempty"`
	AlignToGrid  bool `json:"align_to_grid,omitempty"`
}

type State struct {
	GroupCount   int `json:"groups,omitempty"`
	NodeCount    int `json:"nodes,omitempty"`
	LinkCount    int `json:"links,omitempty"`
	RerouteCount int `json:"reroutes,omitempty"`
}

type Group struct {
	Title     string    `json:"title"`
	Bounding  []float64 `json:"bounding"`
	Color     string    `json:"color,omitempty"`
	FontSize  int       `json:"font_size,omitempty"`
	Locked    bool      `json:"locked,omitempty"`
}

type Node struct {
	ID         interface{}            `json:"id"`
	Type       string                 `json:"type"`
	Pos        interface{}            `json:"pos"`
	Size       interface{}            `json:"size,omitempty"`
	Flags      map[string]interface{} `json:"flags,omitempty"`
	Order      int                    `json:"order,omitempty"`
	Mode       int                    `json:"mode,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Inputs     []Input                `json:"inputs,omitempty"`
	Outputs    []Output               `json:"outputs,omitempty"`
	Widgets    []interface{}          `json:"widgets_values,omitempty"`
}

type Input struct {
	Name string      `json:"name"`
	Type string      `json:"type"`
	Link interface{} `json:"link,omitempty"`
}

type Output struct {
	Name  string   `json:"name"`
	Type  string   `json:"type"`
	Links []int    `json:"links,omitempty"`
}


type Link struct {
	ID         int         `json:"id"`
	OutputNode int         `json:"output_node"`
	OutputSlot int         `json:"output_slot"`
	InputNode  int         `json:"input_node"`
	InputSlot  int         `json:"input_slot"`
	Type       string      `json:"type,omitempty"`
	Extra      interface{} `json:"extra,omitempty"`
}

type Model struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	Directory string `json:"directory"`
	Hash      string `json:"hash,omitempty"`
	HashType  string `json:"hash_type,omitempty"`
}

type Dependency struct {
	Type string
	Name string
	Path string
}

type Analysis struct {
	CustomNodes    []string
	Models         []Model
	Dependencies   []Dependency
	MissingNodes   []string
	MissingModels  []Model
}