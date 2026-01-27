package prefs

type Prefs struct {
	Editor      Editor
	Controls    Controls
	Interface   Interface
	Application Application
}

type Interface struct {
	Scale int
	Fps   int
	Theme string
}

type Controls struct {
	AltScrollBehaviour   bool
	QuickEditContextMenu bool
	QuickEditMapPane     bool
}

type Editor struct {
	SaveFormat        string
	CodeEditor        string
	NudgeMode         string
	SanitizeVariables bool
	// Obsolete object replacement paths
	ObsoleteObjectPath string
	ObsoleteTurfPath   string
	ObsoleteAreaPath   string
	// Randomize direction when placing objects
	RandomizeDirection bool
	// Suppress warning when transferring vars from obsolete objects
	SuppressObsoleteVarWarning bool
}

type Application struct {
	CheckForUpdates bool
	AutoUpdate      bool
}
