package config

type RepublishConfig struct {
	ItemsPerRequest  int
	Goroutines       int
	RequestPerSecond int
}
