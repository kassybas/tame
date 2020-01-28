package vartable

func (vt *VarTable) GetAllEnvVars(ShellFieldSeparator string) []string {
	formattedVars := []string{}
	vt.RLock()
	allVars := vt.vars
	for _, v := range allVars {
		formattedVars = append(formattedVars, v.ToEnvVars(ShellFieldSeparator)...)
	}
	vt.RUnlock()
	return formattedVars
}

func (vt *VarTable) GetAllValues() map[string]interface{} {
	vars := make(map[string]interface{})
	vt.RLock()
	allVars := vt.vars
	for k, v := range allVars {
		vars[k] = v.Value()
	}
	vt.RUnlock()
	return vars
}
