package scheme

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
	ArgContainer    map[string]interface{} `yaml:"args,omitempty"`
	BodyContainer   []interface{}          `yaml:"body,omitempty"`
	ReturnContainer []string               `yaml:"return,omitempty"`
}

type SetConfig struct {
	ShellContainer       string `yaml:"shell,omitempty"`
	InitContainer        string `yaml:"init,omitempty"`
	DefaultOptsContainer string `yaml:"opts,omitempty"`
	ShieldEnvContainer   bool   `yaml:"shield-env,omitempty"`
}
