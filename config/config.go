package config

type RepublishConfig struct {
	ItemsPerRequest   int
	RequestPerSecond  int
	Goroutines        int
	LogSuccessfulPush bool
	LogErrorPush      bool
	LogProgress       bool
}
