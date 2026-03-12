package validator

import (
	"testing"
)

func TestValidateJsonnet_ValidSyntax(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name: "simple object",
			content: `{
  name: "test",
  value: 123,
}`,
		},
		{
			name: "array",
			content: `[
  "item1",
  "item2",
  "item3",
]`,
		},
		{
			name: "local variable",
			content: `local animals = [
  "zebra",
  "elephant",
  "antelope",
];

animals`,
		},
		{
			name: "object with inline array",
			content: `{
  PR: [ "reviewer1", "reviewer2", ],
  name: "test",
}`,
		},
		{
			name: "nested objects",
			content: `{
  outer: {
    inner: {
      value: 42,
    },
  },
}`,
		},
		{
			name: "with comments",
			content: `{
  // This is a comment
  name: "test",  # inline comment
  /* block comment */
  value: 123,
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateJsonnet(tt.content, "test.jsonnet")
			if err != nil {
				t.Errorf("ValidateJsonnet() expected valid jsonnet, got error: %v", err)
			}
		})
	}
}

func TestValidateJsonnet_InvalidSyntax(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "unclosed brace",
			content: `{ name: "test"`,
		},
		{
			name:    "unclosed bracket",
			content: `[ "item1", "item2"`,
		},
		{
			name:    "missing comma",
			content: `{ name: "test" value: 123 }`,
		},
		{
			name:    "invalid syntax",
			content: `this is not valid jsonnet`,
		},
		{
			name:    "mismatched brackets",
			content: `[ "item" }`,
		},
		{
			name:    "extra closing brace",
			content: `{ name: "test" }}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateJsonnet(tt.content, "test.jsonnet")
			if err == nil {
				t.Errorf("ValidateJsonnet() expected error for invalid jsonnet, got nil")
			}
		})
	}
}

func TestValidateJsonnet_EmptyContent(t *testing.T) {
	// Empty content is invalid in jsonnet - it needs at least some expression
	err := ValidateJsonnet("", "test.jsonnet")
	if err == nil {
		t.Errorf("ValidateJsonnet() expected error for empty content, got nil")
	}
}
