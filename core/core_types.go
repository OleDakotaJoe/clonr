package core

type FileRenderer struct {
	filePath   string
	renderFile func(string, map[string]string)
}
