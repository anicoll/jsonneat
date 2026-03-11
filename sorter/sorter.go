package sorter

import (
	"sort"
	"strings"
)

type arrayElement struct {
	original string
	sortKey  string
}

type inlineArrayInfo struct {
	hasInlineArray bool
	prefix         string
	firstElement   string
	indent         string
}

// detectInlineArray checks if a line contains an inline array start like "PR: [ element,"
// Returns info about the inline array, or nil if not detected or if array closes on same line
func detectInlineArray(line, trimmed string) *inlineArrayInfo {
	// Must contain [ and comma, but not start with [
	if !strings.Contains(line, "[") || !strings.Contains(trimmed, ",") || strings.HasPrefix(trimmed, "[") {
		return nil
	}

	bracketIdx := strings.Index(line, "[")
	afterBracket := line[bracketIdx+1:]
	trimmedAfter := strings.TrimSpace(afterBracket)

	// Check if array also closes on same line - if so, don't treat as multi-line
	if strings.Contains(afterBracket, "]") {
		return nil
	}

	// Check if there's an array element after the [
	if trimmedAfter != "" && !strings.HasPrefix(trimmedAfter, "]") && strings.Contains(trimmedAfter, ",") {
		indent := strings.Repeat(" ", len(line)-len(trimmed))
		return &inlineArrayInfo{
			hasInlineArray: true,
			prefix:         line[:bracketIdx+1],
			firstElement:   trimmedAfter,
			indent:         indent,
		}
	}

	return nil
}

// SortJsonnet sorts array elements in jsonnet files alphabetically
func SortJsonnet(content string) (string, error) {
	lines := strings.Split(content, "\n")
	result := make([]string, 0, len(lines))

	var currentBlock []arrayElement
	var arrayStartPrefix string // Stores the prefix like "PR: [" for inline arrays
	inArray := false
	parenDepth := 0
	bracketDepth := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check for inline array start like "PR: [ element,"
		if info := detectInlineArray(line, trimmed); info != nil {
			arrayStartPrefix = info.prefix
			afterBracket := line[len(info.prefix):]
			parenDepth += countChar(afterBracket, '(') - countChar(afterBracket, ')')
			bracketDepth = 1

			currentBlock = append(currentBlock, arrayElement{
				original: info.indent + "      " + info.firstElement,
				sortKey:  extractSortKey(info.firstElement),
			})
			inArray = true
			continue
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
					result = append(result, arrayStartPrefix+" "+firstElem)
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
			result = append(result, arrayStartPrefix+" "+firstElem)
			sortedBlock = sortedBlock[1:]
		}

		result = append(result, sortedBlock...)
	}

	return strings.Join(result, "\n"), nil
}

func countChar(s string, c rune) int {
	return strings.Count(s, string(c))
}

func isArrayLine(line string) bool {
	if line == "" {
		return false
	}

	// Skip comment-only lines
	commentPrefixes := []string{"//", "/*", "*"}
	for _, prefix := range commentPrefixes {
		if strings.HasPrefix(line, prefix) {
			return false
		}
	}

	// Skip structural lines
	structuralPrefixes := []string{"[", "{", "]", "}"}
	for _, prefix := range structuralPrefixes {
		if strings.HasPrefix(line, prefix) {
			return false
		}
	}

	// Must contain comma before any comment
	commaIdx := strings.Index(line, ",")
	if commaIdx == -1 {
		return false
	}

	hashIdx := strings.Index(line, "#")
	slashIdx := strings.Index(line, "//")

	if (hashIdx != -1 && commaIdx > hashIdx) || (slashIdx != -1 && commaIdx > slashIdx) {
		return false
	}

	// Skip structural closing lines like "],", "},", etc.
	codeOnly := stripComment(line)
	structuralClosings := []string{"],", "},", "]", "}"}
	for _, closing := range structuralClosings {
		if codeOnly == closing {
			return false
		}
	}

	return true
}

// stripComment removes comments and trailing commas from a line
func stripComment(line string) string {
	// Remove comments
	if idx := strings.Index(line, "//"); idx != -1 {
		line = line[:idx]
	}
	if idx := strings.Index(line, "#"); idx != -1 {
		line = line[:idx]
	}
	return strings.TrimSpace(line)
}

func extractSortKey(line string) string {
	line = stripComment(strings.TrimSpace(line))
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
