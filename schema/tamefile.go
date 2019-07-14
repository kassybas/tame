package schema

type Tamefile struct {
	TameVersion string                      `yaml:"tame-version,omitempty"`
	Includes    []string                    `yaml:"include,omitempty"`
	Loads       []string                    `yaml:"load,omitempty"`
	Sets        SettingsDefintion           `yaml:"settings,omitempty"`
	Globals     map[string]string           `yaml:"globals,omitempty"`
	Targets     map[string]TargetDefinition `yaml:"targets,omitempty"`

	WorkDir        string            `yaml:"work-dir,omitempty"`
	DefaultEnvVars map[string]string `yaml:"defaults,omitempty"`
}

type TargetDefinition struct {
	ArgDefinition    map[string]interface{} `yaml:"args,omitempty"`
	BodyDefinition   []StepDefinition       `yaml:"body,omitempty"`
	ReturnDefinition []string               `yaml:"return,omitempty"`
	OptsDefinition   []string               `yaml:"opts,omitempty"`
	Summary          string                 `yaml:"summary,omitempty"`
}

type StepDefinition struct {
	Shell  string                       `yaml:"shell,omitempty"`
	Call   map[string]map[string]string `yaml:"call,omitempty"`
	Result []string                     `yaml:"result,omitempty"`
	Opts   []string                     `yaml:"opts,omitempty"`
	Out    string                       `yaml:"out,omitempty"`
	Err    string                       `yaml:"err,omitempty"`
	Rc     string                       `yaml:"rc,omitempty"`
}

type SettingsDefintion struct {
	Shell      string   `yaml:"shell,omitempty"`
	Init       string   `yaml:"init,omitempty"`
	GlobalOpts []string `yaml:"opts,omitempty"`
}
