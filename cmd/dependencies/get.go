package dependencies

// Get returns the default implementation of CommandDependencies
func Get() CommandDependencies {
	return &DefaultDependencies{}
}
