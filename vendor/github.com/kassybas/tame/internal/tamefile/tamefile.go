package tamefile

type Teafile struct {
	TameVersion string                     `yaml:"tame-version,omitempty"`
	Includes   []IncludeConfig            `yaml:"include,omitempty"`
	Loads      []LoadConfig               `yaml:"load,omitempty"`
	Sets       SetConfig                  `yaml:"set,omitempty"`
	Globals    map[string]string          `yaml:"globals,omitempty"`
	Targets    map[string]TargetContainer `yaml:"targets,omitempty"`

	WorkDir        string            `yaml:"work-dir,omitempty"`
	DefaultEnvVars map[string]string `yaml:"defaults,omitempty"`
}

type IncludeConfig struct {
	Alias string `yaml:"alias,omitempty"`
	Path  string `yaml:"path,omitempty"`
}
type LoadConfig struct {
	Alias string `yaml:"alias,omitempty"`
	Path  string `yaml:"path,omitempty"`
}
type SetConfig struct {
	ShellContainer       string `yaml:"shell,omitempty"`
	InitContainer        string `yaml:"init,omitempty"`
	DefaultOptsContainer string `yaml:"opts,omitempty"`
	ShieldEnvContainer   bool   `yaml:"shield-env,omitempty"`
}

type TargetContainer struct {
	Script       string        `yaml:"script,omitempty"`
	ArgContainer []interface{} `yaml:"args,omitempty"`
	DepContainer []interface{} `yaml:"deps,omitempty"`
}
