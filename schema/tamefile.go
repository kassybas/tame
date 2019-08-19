package schema

type Tamefile struct {
	TameVersion string                      `yaml:"tameVersion,omitempty"`
	Includes    []string                    `yaml:"include,omitempty"`
	Loads       []string                    `yaml:"load,omitempty"`
	Sets        SettingsDefintion           `yaml:"settings,omitempty"`
	Globals     map[string]interface{}      `yaml:"globals,omitempty"`
	Targets     map[string]TargetDefinition `yaml:"targets,omitempty"`

	WorkDir        string            `yaml:"workDir,omitempty"`
	DefaultEnvVars map[string]string `yaml:"defaults,omitempty"`
}

type TargetDefinition struct {
	ArgDefinition    map[string]interface{} `yaml:"args,omitempty"`
	BodyDefinition   []StepDefinition       `yaml:"body,omitempty"`
	ReturnDefinition string                 `yaml:"return,omitempty"`
	OptsDefinition   []string               `yaml:"opts,omitempty"`
	Summary          string                 `yaml:"summary,omitempty"`
}

type StepDefinition struct {
	Shell  string                       `yaml:"sh,omitempty"`
	Call   map[string]map[string]string `yaml:"call,omitempty"`
	Var    map[string]interface{}       `yaml:"var,omitempty"`
	Result string                       `yaml:"result,omitempty"`
	Opts   []string                     `yaml:"opts,omitempty"`
	Out    string                       `yaml:"$,omitempty"`
	Err    string                       `yaml:"err$,omitempty"`
	Rc     string                       `yaml:"status$,omitempty"`
}

type SettingsDefintion struct {
	Shell      string   `yaml:"shell,omitempty"`
	Init       string   `yaml:"init,omitempty"`
	GlobalOpts []string `yaml:"opts,omitempty"`
	ShieldEnv  bool     `yaml:"shieldEnv,omitempty"`
}
