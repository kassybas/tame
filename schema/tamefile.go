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
	ArgDefinition  map[string]interface{}   `yaml:"args,omitempty"`
	StepDefinition []map[string]interface{} `yaml:"run,omitempty"`
	OptsDefinition []string                 `yaml:"opts,omitempty"`
	Summary        string                   `yaml:"summary,omitempty"`
}

type SettingsDefintion struct {
	Shell      string   `yaml:"shell,omitempty"`
	Init       string   `yaml:"init,omitempty"`
	GlobalOpts []string `yaml:"opts,omitempty"`
	ShieldEnv  bool     `yaml:"shieldEnv,omitempty"`
}
