package vartable

func (vt *VarTable) GetAllEnvVars(ShellFieldSeparator string) []string {
	formattedVars := []string{}
	vt.RLock()
	allVars := vt.vars
	vt.RUnlock()
	for _, v := range allVars {
		formattedVars = append(formattedVars, v.ToEnvVars(ShellFieldSeparator)...)
	}
	return formattedVars
}

func (vt *VarTable) GetAllValues() map[string]interface{} {
	vars := make(map[string]interface{})
	vt.RLock()
	allVars := vt.vars
	vt.RUnlock()
	for k, v := range allVars {
		vars[k] = v.Value()
	}
	return vars
}
