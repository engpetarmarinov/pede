package preprocessor

import (
	"strings"
)

// Rule defines a function that determines if a line should be stripped.
type Rule func(line string) bool

// DefaultRules returns the default set of rules: strip comments and empty/unknown lines.
func DefaultRules() []Rule {
	return []Rule{
		// Strip comment lines
		func(line string) bool { return strings.HasPrefix(strings.TrimSpace(line), "//") },
		// Strip empty lines
		func(line string) bool { return strings.TrimSpace(line) == "" },
	}
}

// Preprocess applies the given rules to the input string, stripping lines that match any rule.
func Preprocess(input string, rules []Rule) (string, error) {
	var sb strings.Builder
	lines := strings.Split(input, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		strip := false
		for _, rule := range rules {
			if rule(trimmed) {
				strip = true
				break
			}
		}
		if strip {
			continue
		}
		// Remove inline comments (everything after //)
		if idx := strings.Index(trimmed, "//"); idx != -1 {
			trimmed = strings.TrimSpace(trimmed[:idx])
			if trimmed == "" {
				continue
			}
		}
		sb.WriteString(trimmed)
		// Only add a newline if not the last line
		if i < len(lines)-1 || trimmed != "" {
			sb.WriteString("\n")
		}
	}
	// Ensure output ends with a newline if not empty
	out := sb.String()
	if len(out) > 0 && out[len(out)-1] != '\n' {
		out += "\n"
	}
	return out, nil
}
