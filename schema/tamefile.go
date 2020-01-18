package schema

import "github.com/kassybas/tame/types/steptype"

type Tamefile struct {
	TameVersion string                  `mapstructure:"tameVersion,omitempty"`
	Includes    []IncludeSchema         `mapstructure:"include,omitempty"`
	Loads       []string                `mapstructure:"load,omitempty"`
	Sets        SettingsShema           `mapstructure:"settings,omitempty"`
	Summary     string                  `mapstructure:"summary,omitempty"`
	Globals     map[string]interface{}  `mapstructure:"-,omitempty"`
	Targets     map[string]TargetSchema `mapstructure:"-,omitempty"`

	WorkDir        string            `mapstructure:"workDir,omitempty"`
	DefaultEnvVars map[string]string `mapstructure:"defaults,omitempty"`
}

type IncludeSchema struct {
	Path  string `mapstructure:"path,omitempty"`
	Alias string `mapstructure:"as,omitempty"`
}

type TargetSchema struct {
	ArgDefinition  map[string]interface{}   `mapstructure:"args,omitempty"`
	StepDefinition []map[string]interface{} `mapstructure:"run,omitempty"`
	OptsDefinition []string                 `mapstructure:"opts,omitempty"`
	Summary        string                   `mapstructure:"summary,omitempty"`
}

type SettingsShema struct {
	Shell               string   `mapstructure:"shell,omitempty"`
	Init                string   `mapstructure:"init,omitempty"`
	GlobalOpts          []string `mapstructure:"opts,omitempty"`
	ShieldEnv           bool     `mapstructure:"shieldEnv,omitempty"`
	ShellFieldSeparator string   `mapstructure:"shellFieldSeparator,omitempty"`
}

type DumpSchema struct {
	SourceVarName string `mapstructure:"var,omitempty"`
	Path          string `mapstructure:"path,omitempty"`
	Format        string `mapstructure:"format,omitempty"`
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
	Dump             *DumpSchema             `mapstructure:"dump"`

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
