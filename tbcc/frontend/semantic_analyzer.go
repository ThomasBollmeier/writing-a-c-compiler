package frontend

func AnalyzeSemantics(program *Program, nameCreator NameCreator) (*Program, error) {
	checker := newLabelChecker()
	err := checker.check(program)
	if err != nil {
		return nil, err
	}

	labeler := newLoopLabeler(nameCreator)
	err = labeler.addLabels(program)
	if err != nil {
		return nil, err
	}

	resolver := newVariableResolver(nameCreator)
	return resolver.resolve(program)
}
