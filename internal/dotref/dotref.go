package dotref

// func Empty() []RefField {
// 	return []RefField{}
// }

// func checkIndexSeparators(field string) (bool, error) {
// 	if !strings.Contains(field, keywords.IndexingSeparatorL) && !strings.Contains(field, keywords.IndexingSeparatorR) {
// 		// no indexing in field
// 		return false, nil
// 	}
// 	if strings.Contains(field, keywords.IndexingSeparatorL) != strings.Contains(field, keywords.IndexingSeparatorR) {
// 		// one of them is not in the string
// 		return false, fmt.Errorf("missing bracket: %s", field)
// 	}
// 	return true, nil
// }

// func IsDotRef(dotName string) bool {
// 	if strings.Contains(dotName, keywords.IndexingSeparatorL) || strings.Contains(dotName, keywords.IndexingSeparatorR) ||
// 		strings.Contains(dotName, keywords.TameFieldSeparator) {
// 		return true
// 	}
// 	return false
// }

// func getSingleRef(fieldStr string) (*RefField, error) {
// 	isLDot := strings.HasPrefix(fieldStr, keywords.TameFieldSeparator)
// 	isRDot := strings.HasSuffix(fieldStr, keywords.TameFieldSeparator)
// 	isLBr := strings.HasPrefix(fieldStr, keywords.IndexingSeparatorL)
// 	isRBr := strings.HasSuffix(fieldStr, keywords.IndexingSeparatorR)

// 	if len(fieldStr) == 0 {
// 		return nil, nil
// 	}
// 	// variable reference
// 	if isLDot {
// 		s := strings.TrimPrefix(fieldStr, keywords.TameFieldSeparator)
// 		if strings.HasPrefix(s, keywords.PrefixReference) || strings.HasPrefix(s, `"`) {
// 			return nil, fmt.Errorf("dotted field reference cannot be variable or quoted: %s", fieldStr)
// 		}
// 		return &RefField{
// 			FieldName: s,
// 			Type:      Literal,
// 		}, nil
// 	}
// 	// index reference
// 	index, err := strconv.Atoi(fieldStr)
// 	if err == nil {
// 		return &RefField{
// 			Index: index,
// 			Type:  Index,
// 		}, nil
// 	}
// 	// if (!strings.HasPrefix(fieldStr, `"`) || !strings.HasSuffix(fieldStr, `"`) {
// 	// 	return nil, fmt.Errorf(`literal references in brackets must be qoted, got: [%s] (correct: ["%s"]`, fieldStr, fieldStr)
// 	// }
// 	// literal reference
// 	return &RefField{
// 		FieldName: strings.Trim(fieldStr, `"`),
// 		Type:      Literal,
// 	}, nil
// 	// TODO: check for quoting
// }

// func findClosingBracket(dotName string) int {
// 	bracketsStarted := 0
// 	pos := 0
// 	// the first iterator in range of string is the byte position
// 	for _, ch := range dotName {
// 		if ch == '[' {
// 			bracketsStarted++
// 		} else if ch == ']' {
// 			bracketsStarted--
// 			if bracketsStarted == 0 {
// 				return pos
// 			}
// 			if bracketsStarted < 0 {
// 				return -1
// 			}
// 		}
// 		pos++
// 	}
// 	return -1
// }

// func parseUntilDot(dotName string, firstDot int) ([]RefField, error) {
// 	result := []RefField{}
// 	field, err := getSingleRef(dotName[0:firstDot], false)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if field != nil {
// 		result = []RefField{*field}
// 	}
// 	if len(dotName) == firstDot+1 {
// 		return nil, fmt.Errorf("reference ending with '.': %s", dotName)
// 	}
// 	remainingFields, err := ParseFieldsRec(dotName[firstDot:])
// 	result = append(result, remainingFields...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return result, nil
// }

// func parseUntilLBr(dotName string, firstLBr int) ([]RefField, error) {
// 	result := []RefField{}
// 	if len(dotName) == firstLBr+1 {
// 		return nil, fmt.Errorf("reference ending with '[': %s", dotName)
// 	}

// 	field, err := getSingleRef(dotName[0:firstLBr], true)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if field != nil {
// 		result = []RefField{*field}
// 	}
// 	closingRBr := findClosingBracket(dotName)
// 	if closingRBr == -1 {
// 		return nil, fmt.Errorf("could not find closing bracket in reference: %s", dotName)
// 	}
// 	// inner fields
// 	innerStr := dotName[firstLBr+1 : closingRBr]
// 	innerFields, err := ParseFieldsRec(innerStr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	result = append(result, RefField{InnerRef: innerFields, Type: InnerRef})
// 	//remaining fields

// 	remainingFields, err := ParseFieldsRec(dotName[closingRBr+1:])
// 	if err != nil {
// 		return nil, err
// 	}
// 	result = append(result, remainingFields...)
// 	return result, nil
// }

// func ParseFieldsRec(dotName string) ([]RefField, error) {
// 	fmt.Println("STARTING...", dotName)
// 	if len(dotName) == 0 {
// 		return []RefField{}, nil
// 	}
// 	firstDot := strings.Index(dotName, keywords.TameFieldSeparator)
// 	firstLBr := strings.Index(dotName, keywords.IndexingSeparatorL)
// 	if firstDot == -1 && firstLBr == -1 {
// 		// no field reference
// 		field, err := getSingleRef(dotName, false)
// 		if field == nil {
// 			return []RefField{}, err
// 		}
// 		return []RefField{*field}, err
// 	}

// 	// first is dot
// 	if firstLBr == -1 || (firstDot != -1 && firstDot < firstLBr) {
// 		return parseUntilDot(dotName, firstDot)
// 	}
// 	// first is LBr
// 	return parseUntilLBr(dotName, firstLBr)

// }

// func ParseFields(dotName string) ([]RefField, error) {
// 	tmp_fields := strings.Split(dotName, keywords.TameFieldSeparator)
// 	fields := []RefField{}

// 	for _, field := range tmp_fields {
// 		hasIndex, err := checkIndexSeparators(field)
// 		if err != nil {
// 			return fields, err
// 		}
// 		if hasIndex {
// 			index, f, indexRef, err := parseBrackets(field)
// 			if err != nil {
// 				return fields, err
// 			}
// 			fields = append(fields, RefField{FieldName: f})
// 			fields = append(fields, RefField{Index: index, FieldName: indexRef})
// 			continue
// 		}
// 		fields = append(fields, RefField{FieldName: field})
// 	}
// 	return fields, nil
// }

// // returns the index, variable name and index variable if exists
// // if index is a variable reference, -1 is returned
// func parseBrackets(name string) (int, string, string, error) {
// 	lBr := strings.Index(name, keywords.IndexingSeparatorL) + 1
// 	rBr := strings.Index(name, keywords.IndexingSeparatorR)
// 	index, err := strconv.Atoi(name[lBr:rBr])
// 	if err != nil {
// 		// variable index or bracket ref
// 		return -1, name[0 : lBr-1], name[lBr:rBr], nil
// 	}
// 	return index, name[0 : lBr-1], "", nil
// }
