package dependencies

// Get returns the default implementation of CommandDependencies
var Get = func() CommandDependencies {
	return &DefaultDependencies{}
}
