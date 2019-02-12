package liquidweb

// EexpandSetToStrings expands the type, TypeSet into a list of strings.
func expandSetToStrings(strings []interface{}) []string {
	expandedStrings := make([]string, len(strings))
	for i, v := range strings {
		expandedStrings[i] = v.(string)
	}

	return expandedStrings
}
