package main

type Config struct {
	ConsumerKey    string `yaml:"ConsumerKey"`
	ConsumerSecret string `yaml:"ConsumerSecret"`
	SeedString     string `yaml:"SeedString"`
	URL            string `yaml:"URL"`
	Port           string `yaml:"Port"`
	DBUser           string `yaml:"DBUser"`
	DBPass          string `yaml:"DBPass"`
	DBName           string `yaml:"DBName"`
}
