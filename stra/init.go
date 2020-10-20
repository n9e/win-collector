package stra

import model "github.com/didi/nightingale/src/models"

var StraConfig StraSection
var Collect model.Collect

type StraSection struct {
	Enable   bool   `yaml:"enable"`
	Interval int    `yaml:"interval"`
	Api      string `yaml:"api"`
	Timeout  int    `yaml:"timeout"`
	PortPath string `yaml:"portPath"`
	ProcPath string `yaml:"procPath"`
	LogPath  string `yaml:"logPath"`
}

///api/portal/collects/%s

func Init(stra StraSection) {
	StraConfig = stra

	GetCollects()
}
