package types

type CloneCmdArgs struct {
	Args        []string
	NameFlag    string
	IsLocalPath bool
	IsAlias     bool
}

type DocsCmdArgs struct {
	Args      []string
	OutputDir string
}

type AliasCmdArgs struct {
	Args          []string
	AddFlag       bool
	UpdateFlag    bool
	DeleteFlag    bool
	IsLocalFlag   bool
	AliasNameFlag string
	AliasLocation string
}
