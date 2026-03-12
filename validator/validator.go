package validator

import (
	"fmt"

	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
)

// ValidateJsonnet validates that the given jsonnet content is syntactically correct
// using the go-jsonnet parser
func ValidateJsonnet(content, filename string) error {
	// Create a new VM for parsing
	vm := jsonnet.MakeVM()

	// Parse the jsonnet content
	node, err := jsonnet.SnippetToAST(filename, content)
	if err != nil {
		return fmt.Errorf("jsonnet parse error: %w", err)
	}

	// Verify we got a valid AST node
	if node == nil {
		return fmt.Errorf("jsonnet parsing produced nil AST")
	}

	// Additional validation: ensure the AST is not empty
	if _, ok := node.(*ast.Error); ok {
		return fmt.Errorf("jsonnet parsing produced error node")
	}

	// VM is unused but kept for potential future use (e.g., evaluation)
	_ = vm

	return nil
}
