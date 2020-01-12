package schema

type Tamefile struct {
	TameVersion string                  `yaml:"tameVersion,omitempty"`
	Includes    []IncludeSchema         `yaml:"include,omitempty"`
	Loads       []string                `yaml:"load,omitempty"`
	Sets        SettingsShema           `yaml:"settings,omitempty"`
	Globals     map[string]interface{}  `yaml:"globals,omitempty"`
	Targets     map[string]TargetSchema `yaml:"targets,omitempty"`
	Commands    map[string]string       `yaml:"cmds,omitempty"`

	WorkDir        string            `yaml:"workDir,omitempty"`
	DefaultEnvVars map[string]string `yaml:"defaults,omitempty"`
}

type IncludeSchema struct {
	Path  string `yaml:"path,omitempty"`
	Alias string `yaml:"as,omitempty"`
}

type TargetSchema struct {
	ArgDefinition  map[string]interface{}   `yaml:"args,omitempty"`
	StepDefinition []map[string]interface{} `yaml:"run,omitempty"`
	OptsDefinition []string                 `yaml:"opts,omitempty"`
	Summary        string                   `yaml:"summary,omitempty"`
}

type SettingsShema struct {
	Shell               string   `yaml:"shell,omitempty"`
	Init                string   `yaml:"init,omitempty"`
	GlobalOpts          []string `yaml:"opts,omitempty"`
	ShieldEnv           bool     `yaml:"shieldEnv,omitempty"`
	ShellFieldSeparator string   `yaml:"shellFieldSeparator,omitempty"`
}

type ForLoopSchema struct {
	Iterator string      `mapstructure:"$"`
	Iterable interface{} `mapstructure:"in"`
}

// MergedStepSchema is the base format of step
type MergedStepSchema struct {
	ForLoop          ForLoopSchema  `mapstructure:"for"`
	Condition        *string        `mapstructure:"if"`
	Return           *[]interface{} `mapstructure:"return"` // string is allowed due to weak decode
	Opts             *[]string      `mapstructure:"opts"`   // string is allowed due to weak decode
	ResultContainers *[]string      `mapstructure:"$"`      // string is allowed due to weak decode
	Script           *[]string      `mapstructure:"sh"`     // string is allowed due to weak decode
	Expr             *string        `mapstructure:"expr"`
	Wait             *interface{}   `mapstructure:"wait"`

	// Name is a dynamic key can be either (but only one of):
	CalledTargetName    *string                `mapstructure:"-"` // loaded dynamically since the key is the called target
	CallArgumentsPassed map[string]interface{} `mapstructure:"-"` // loaded dynamically since the key is the called target so arguments are unkown
	VarName             *string                `mapstructure:"-"` // loaded dynamically since the key is the variable name
	VarValue            interface{}            `mapstructure:"-"` // loaded dynamically since the key is the variable name so value is unknown
}
