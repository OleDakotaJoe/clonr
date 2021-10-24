package types

type CloneCmdArgs struct {
	Args        []string
	NameFlag    string
	IsLocalPath bool
}

type DocsCmdArgs struct {
	Args      []string
	OutputDir string
}
