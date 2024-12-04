package aws

func NewTags(name string) map[string]string {
	return map[string]string{
		"Name":    name,
		"Project": "flowlogs-cli",
	}
}
