package schema

import "github.com/kassybas/tame/types/steptype"

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

// MergedStepSchema is the base format of step
type MergedStepSchema struct {
	ForLoop          *map[string]interface{} `mapstructure:"for"`
	ForRawSteps      []interface{}           `mapstructure:"do"`
	Return           *[]interface{}          `mapstructure:"return"` // single interface is allowed due to weak decode
	Opts             *[]string               `mapstructure:"opts"`   // string is allowed due to weak decode
	ResultContainers *[]string               `mapstructure:"="`      // string is allowed due to weak decode
	Script           *[]string               `mapstructure:"sh"`     // string is allowed due to weak decode
	Expr             *string                 `mapstructure:"expr"`
	Wait             *interface{}            `mapstructure:"wait"`

	// loaded dynamically since the yaml key defines the step type or step data
	ForSteps            []MergedStepSchema     `mapstructure:"-"` // do: [ForSteps]
	IfCondition         string                 `mapstructure:"-"` // if IfCondition
	IfSteps             []MergedStepSchema     `mapstructure:"-"` // if IfCondition: IfSteps
	ElseSteps           []MergedStepSchema     `mapstructure:"-"` // else: ElseSteps
	CalledTargetName    string                 `mapstructure:"-"` // CalledTargetName: {}
	CallArgumentsPassed map[string]interface{} `mapstructure:"-"` // CalledTargetName: {CallArgumentsPassed}
	VarName             string                 `mapstructure:"-"` // $VarName
	VarValue            interface{}            `mapstructure:"-"` // $VarName: VarValue
	StepType            steptype.Steptype      `mapstructure:"-"` // type of step
}
