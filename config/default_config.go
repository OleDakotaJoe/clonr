package config


type config struct {
	DefaultProjectName string
}


func DefaultConfig() *config {
	this := config{
		DefaultProjectName: "clonr-app",
	}
	return &this
}