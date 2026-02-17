package tools

var registry = make(map[string]Tool)

func RegisterTool(tool Tool) {
	registry[tool.Name()] = tool
}

func GetTool(name string) (Tool, bool) {
	tool, exists := registry[name]
	return tool, exists
}
