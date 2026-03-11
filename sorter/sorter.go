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
	var arrayStartPrefix string  // Stores the prefix like "PR: [" for inline arrays
	inArray := false
	parenDepth := 0
	bracketDepth := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Special handling for inline array start: "PR: [ element,"
		// Split this into prefix and first element
		if strings.Contains(line, "[") && strings.Contains(trimmed, ",") && !strings.HasPrefix(trimmed, "[") {
			bracketIdx := strings.Index(line, "[")
			afterBracket := line[bracketIdx+1:]
			trimmedAfter := strings.TrimSpace(afterBracket)

			// Check if there's an array element after the [
			if trimmedAfter != "" && !strings.HasPrefix(trimmedAfter, "]") && strings.Contains(trimmedAfter, ",") {
				// Save the prefix (everything up to and including [)
				arrayStartPrefix = line[:bracketIdx+1]

				// Process the rest as an array element
				parenDepth += countChar(afterBracket, '(') - countChar(afterBracket, ')')
				bracketDepth = 1

				sortKey := extractSortKey(trimmedAfter)
				// Create the element with proper indentation (match continuation lines)
				indent := strings.Repeat(" ", len(line) - len(trimmed))
				currentBlock = append(currentBlock, arrayElement{
					original: indent + "      " + trimmedAfter,  // Add extra indent for continuation
					sortKey:  sortKey,
				})
				inArray = true
				continue
			}
		}

		// Update parentheses and bracket depth for this line
		parenDepth += countChar(line, '(') - countChar(line, ')')
		bracketDepth += countChar(line, '[') - countChar(line, ']')

		// Check if this line is an array element
		isArrayElement := bracketDepth > 0 && parenDepth == 0 && isArrayLine(trimmed)

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

				// If we have an inline array prefix, reconstruct the first line
				if arrayStartPrefix != "" {
					// First sorted element goes on the same line as the prefix
					firstElem := strings.TrimSpace(sortedBlock[0])
					result = append(result, arrayStartPrefix + " " + firstElem)
					sortedBlock = sortedBlock[1:]
					arrayStartPrefix = ""
				}

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

		if arrayStartPrefix != "" {
			firstElem := strings.TrimSpace(sortedBlock[0])
			result = append(result, arrayStartPrefix + " " + firstElem)
			sortedBlock = sortedBlock[1:]
		}

		result = append(result, sortedBlock...)
	}

	return strings.Join(result, "\n"), nil
}

func countChar(s string, c rune) int {
	count := 0
	for _, ch := range s {
		if ch == c {
			count++
		}
	}
	return count
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
