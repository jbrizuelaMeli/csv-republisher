package config

type RepublishConfig struct {
	ItemsPerRequest   int
	RequestPerSecond  int
	LogSuccessfulPush bool
	LogErrorPush      bool
	LogProgress       bool
}
