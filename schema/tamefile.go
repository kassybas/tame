package schema

type Tamefile struct {
	TameVersion string                     `yaml:"tame-version,omitempty"`
	Includes    []string                   `yaml:"include,omitempty"`
	Loads       []string                   `yaml:"load,omitempty"`
	Sets        SetConfig                  `yaml:"set,omitempty"`
	Globals     map[string]string          `yaml:"globals,omitempty"`
	Targets     map[string]TargetContainer `yaml:"targets,omitempty"`

	WorkDir        string            `yaml:"work-dir,omitempty"`
	DefaultEnvVars map[string]string `yaml:"defaults,omitempty"`
}

type TargetContainer struct {
	ArgContainer    map[string]interface{}   `yaml:".args,omitempty"`
	BodyContainer   []map[string]interface{} `yaml:".body,omitempty"`
	ReturnContainer []string                 `yaml:".return,omitempty"`
	OptsContainer   string                   `yaml:".opts,omitempty"`
}

type SetConfig struct {
	Shell       string `yaml:"shell,omitempty"`
	Init        string `yaml:"init,omitempty"`
	DefaultOpts string `yaml:"opts,omitempty"`
	ShieldEnv   bool   `yaml:"shield-env,omitempty"`
}
