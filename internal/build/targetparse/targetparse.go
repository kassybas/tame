package targetparse

import (
	"fmt"
	"strings"

	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/schema"
	"github.com/mitchellh/mapstructure"
)

func ParseTargets(dynamicKeys map[string]interface{}) (map[string]schema.TargetSchema, error) {
	targets := map[string]schema.TargetSchema{}
	for k, v := range dynamicKeys {
		if strings.HasPrefix(k, keywords.PrefixReference) {
			continue
		}
		var newTargetSch schema.TargetSchema
		switch v := v.(type) {
		// simple targets with single commands result in a single shell stepped target
		case string:
			{
				newTargetSch.StepDefinition = []map[string]interface{}{
					map[string]interface{}{
						keywords.ShellStep: v,
					},
				}
			}
		// complex targets
		default:
			{
				var md mapstructure.Metadata
				err := mapstructure.DecodeMetadata(v, &newTargetSch, &md)
				if err != nil {
					return nil, fmt.Errorf("failed to parse target, incorrect yaml format: %s\n\t%s", k, err.Error())
				}
				if len(md.Unused) != 0 {
					return nil, fmt.Errorf("unknown keys in target: %s\n\t%s", k, md.Unused)
				}
			}
		}
		targets[k] = newTargetSch
		delete(dynamicKeys, k)
	}
	return targets, nil
}
