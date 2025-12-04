package models

// Environment represents where a mod can run
type Environment string

const (
	EnvironmentClient Environment = "client"
	EnvironmentServer Environment = "server"
	EnvironmentBoth   Environment = "both"
)

// String returns the string representation of the environment
func (e Environment) String() string {
	return string(e)
}

// IsCompatibleWith checks if this environment is compatible with another
func (e Environment) IsCompatibleWith(other Environment) bool {
	if e == EnvironmentBoth || other == EnvironmentBoth {
		return true
	}
	return e == other
}
