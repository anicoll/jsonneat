package sorter

import (
	"sort"
	"strings"
)

type arrayElement struct {
	original string
	sortKey  string
}

// SortJsonnet sorts array elements in jsonnet files alphabetically
func SortJsonnet(content string) (string, error) {
	lines := strings.Split(content, "\n")
	result := make([]string, 0, len(lines))

	var currentBlock []arrayElement
	inArray := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check if this line is an array element (ends with comma, possibly with comment)
		isArrayElement := isArrayLine(trimmed)

		if isArrayElement {
			// Extract the sortable part (before any comment)
			sortKey := extractSortKey(trimmed)
			currentBlock = append(currentBlock, arrayElement{
				original: line,
				sortKey:  sortKey,
			})
			inArray = true
		} else {
			// If we were in an array block, sort and flush it
			if inArray && len(currentBlock) > 0 {
				sortedBlock := sortBlock(currentBlock)
				result = append(result, sortedBlock...)
				currentBlock = nil
				inArray = false
			}

			// Add the current non-array line
			result = append(result, line)
		}
	}

	// Flush any remaining array block
	if len(currentBlock) > 0 {
		sortedBlock := sortBlock(currentBlock)
		result = append(result, sortedBlock...)
	}

	return strings.Join(result, "\n"), nil
}

func isArrayLine(line string) bool {
	if line == "" {
		return false
	}

	// Skip comment-only lines
	if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") || strings.HasPrefix(line, "*") {
		return false
	}

	// Skip lines that are structural (brackets, braces)
	if strings.HasPrefix(line, "[") || strings.HasPrefix(line, "{") ||
		strings.HasPrefix(line, "]") || strings.HasPrefix(line, "}") {
		return false
	}

	// Check if line contains a comma (before or after any comment)
	hasComma := strings.Contains(line, ",")
	if !hasComma {
		return false
	}

	// Make sure the comma comes before any comment
	hashIdx := strings.Index(line, "#")
	slashIdx := strings.Index(line, "//")
	commaIdx := strings.Index(line, ",")

	// If there's a hash comment, comma must come before it
	if hashIdx != -1 && commaIdx > hashIdx {
		return false
	}

	// If there's a slash comment, comma must come before it
	if slashIdx != -1 && commaIdx > slashIdx {
		return false
	}

	// Skip object/array structure lines like "],", "},", etc.
	codeBeforeComment := line
	if hashIdx != -1 {
		codeBeforeComment = line[:hashIdx]
	}
	if slashIdx != -1 && (hashIdx == -1 || slashIdx < hashIdx) {
		codeBeforeComment = line[:slashIdx]
	}
	codeBeforeComment = strings.TrimSpace(codeBeforeComment)

	if codeBeforeComment == "]," || codeBeforeComment == "}," ||
		codeBeforeComment == "]" || codeBeforeComment == "}" {
		return false
	}

	return true
}

func extractSortKey(line string) string {
	// Remove leading/trailing whitespace
	line = strings.TrimSpace(line)

	// Split on comment markers to get just the code part
	if idx := strings.Index(line, "//"); idx != -1 {
		line = line[:idx]
	}
	if idx := strings.Index(line, "#"); idx != -1 {
		line = line[:idx]
	}

	// Trim again to remove trailing spaces before comment
	line = strings.TrimSpace(line)

	// Remove trailing comma
	line = strings.TrimSuffix(line, ",")

	return strings.ToLower(line)
}

func sortBlock(block []arrayElement) []string {
	// Sort by the sort key
	sort.SliceStable(block, func(i, j int) bool {
		return block[i].sortKey < block[j].sortKey
	})

	// Extract the original lines in sorted order
	result := make([]string, len(block))
	for i, elem := range block {
		result[i] = elem.original
	}

	return result
}

// CleanupWhitespace removes trailing whitespace and ensures consistent formatting
func CleanupWhitespace(content string) string {
	lines := strings.Split(content, "\n")
	result := make([]string, 0, len(lines))

	for _, line := range lines {
		// Remove trailing whitespace
		cleaned := strings.TrimRight(line, " \t")
		result = append(result, cleaned)
	}

	output := strings.Join(result, "\n")

	// Ensure file ends with a single newline if it had content
	if len(output) > 0 && !strings.HasSuffix(output, "\n") {
		output += "\n"
	}

	// Remove multiple consecutive blank lines (more than 2)
	output = strings.ReplaceAll(output, "\n\n\n\n", "\n\n\n")
	output = strings.ReplaceAll(output, "\n\n\n\n", "\n\n\n")

	return output
}
