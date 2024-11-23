package frontend

func AnalyzeSemantics(program *Program, nameCreator NameCreator) (*Program, *Environment, error) {

	labeler := newLoopLabeler(nameCreator)
	err := labeler.addLabels(program)
	if err != nil {
		return nil, nil, err
	}

	checker := newLabelChecker()
	err = checker.check(program)
	if err != nil {
		return nil, nil, err
	}

	globalEnv := NewEnvironment(nil)
	errorList := newTypeChecker(globalEnv).check(program)
	if len(errorList) > 0 {
		return nil, nil, errorList[0]
	}

	resolver := newIdentifierResolver(nameCreator)
	program, err = resolver.resolve(program)
	if err != nil {
		return nil, nil, err
	}

	return program, globalEnv, nil
}
