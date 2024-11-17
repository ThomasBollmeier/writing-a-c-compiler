package frontend

func AnalyzeSemantics(program *Program, nameCreator NameCreator) (*Program, error) {

	labeler := newLoopLabeler(nameCreator)
	err := labeler.addLabels(program)
	if err != nil {
		return nil, err
	}

	checker := newLabelChecker()
	err = checker.check(program)
	if err != nil {
		return nil, err
	}

	globalEnv := newEnvironment(nil)
	errorList := newTypeChecker(globalEnv).check(program)
	if len(errorList) > 0 {
		return nil, errorList[0]
	}

	resolver := newIdentifierResolver(nameCreator)
	return resolver.resolve(program)
}
