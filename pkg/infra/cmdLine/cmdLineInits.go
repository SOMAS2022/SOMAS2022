package cmdline

type CmdLine struct {
	FixedSanctionDuration int
	DynamicSanctions      bool
	GraduatedSanctions    bool
}

var CmdLineInits CmdLine = CmdLine{}
