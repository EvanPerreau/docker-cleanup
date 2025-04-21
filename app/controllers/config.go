package controllers

type config struct {
	DryRun    bool
	OlderThan int
	ShowSize  bool
}

var conf = config{
	DryRun:    false,
	OlderThan: 0,
	ShowSize:  false,
}

func GetConfig() *config {
	return &conf
}
