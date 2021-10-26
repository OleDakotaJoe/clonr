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
	Args                []string
	AddFlag             bool
	UpdateFlag          bool
	DeleteFlag          bool
	IsLocalFlag         bool
	AliasNameFlag       string
	ActualAliasName     string
	AliasLocationFlag   string
	ActualAliasLocation string
	ConfirmFunction     func(string)
	StringInputReader   InputReader
}
