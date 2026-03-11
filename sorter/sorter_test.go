package sorter

import (
	"testing"
)

func TestSortJsonnet(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "BasicSorting",
			input: `local items = [
  third,
  first,
  second,
];`,
			expected: `local items = [
  first,
  second,
  third,
];`,
		},
		{
			name: "SortingWithComments",
			input: `local animals = [
  tiger,  # Big cat
  elephant,  # Large mammal
  antelope,  # Swift runner
];`,
			expected: `local animals = [
  antelope,  # Swift runner
  elephant,  # Large mammal
  tiger,  # Big cat
];`,
		},
		{
			name: "AlreadySorted",
			input: `local items = [
  alpha,
  bravo,
  charlie,
];`,
			expected: `local items = [
  alpha,
  bravo,
  charlie,
];`,
		},
		{
			name: "SingleElement",
			input: `local items = [
  single,
];`,
			expected: `local items = [
  single,
];`,
		},
		{
			name: "MultipleArrays",
			input: `local first = [
  zebra,
  alpha,
];

local second = [
  beta,
  gamma,
  alpha,
];`,
			expected: `local first = [
  alpha,
  zebra,
];

local second = [
  alpha,
  beta,
  gamma,
];`,
		},
		{
			name: "PreserveNonArrayLines",
			input: `local config = {
  name: "test",
  items: [
    zebra,
    alpha,
  ],
  other: "value",
};`,
			expected: `local config = {
  name: "test",
  items: [
    alpha,
    zebra,
  ],
  other: "value",
};`,
		},
		{
			name: "CaseInsensitiveSorting",
			input: `local items = [
  Zebra,
  alpha,
  Charlie,
];`,
			expected: `local items = [
  alpha,
  Charlie,
  Zebra,
];`,
		},
		{
			name: "SortingWithSlashComments",
			input: `local items = [
  third,  // Third item
  first,  // First item
  second,  // Second item
];`,
			expected: `local items = [
  first,  // First item
  second,  // Second item
  third,  // Third item
];`,
		},
		{
			name: "EmptyArray",
			input: `local items = [
];`,
			expected: `local items = [
];`,
		},
		{
			name: "NestedPathsSorting",
			input: `local animals = [
  zoo.africa.zebra,
  zoo.africa.elephant,
  zoo.africa.lion,
];`,
			expected: `local animals = [
  zoo.africa.elephant,
  zoo.africa.lion,
  zoo.africa.zebra,
];`,
		},
		{
			name: "ComplexMultiLineObjects",
			input: `local species = {
  mammals: [
    species.cat,  // Feline species
    species.elephant,  // Largest land animal
    species.bat,  // Flying mammal
  ],
  birds: [
    birds.penguin,  // Flightless bird
    birds.eagle,  // Bird of prey
    birds.canary,  // Song bird
  ],
};`,
			expected: `local species = {
  mammals: [
    species.bat,  // Flying mammal
    species.cat,  // Feline species
    species.elephant,  // Largest land animal
  ],
  birds: [
    birds.canary,  // Song bird
    birds.eagle,  // Bird of prey
    birds.penguin,  // Flightless bird
  ],
};`,
		},
		{
			name: "FunctionCallParameters",
			input: `{
  TEAM_A: {
    name: 'team-a',
    description: 'Team A description',
    include: combineMembers(
      team.resources.team_a_members.include,
      team.resources.team_a_engineers.include,
    ),
  },
}`,
			expected: `{
  TEAM_A: {
    name: 'team-a',
    description: 'Team A description',
    include: combineMembers(
      team.resources.team_a_members.include,
      team.resources.team_a_engineers.include,
    ),
  },
}`,
		},
		{
			name: "RealWorldComplexExample",
			input: `{
  TEAM_ZEBRA: {
    azuread: { name: 'TEAM_ZEBRA' },
    github: { name: 'team-zebra' },
    description: 'Zebra team members',
    include: combineMembers(
      org.resources.team_zebra_members.include,
      org.resources.team_zebra_engineers.include,
    ),
  },

  TEAM_ALPHA: {
    azuread: { name: 'TEAM_ALPHA' },
    github: { name: 'team-alpha' },
    description: 'Alpha team members',
    include: combineMembers(
      org.resources.team_alpha_members.include,
      org.resources.team_alpha_engineers.include,
    ),
  },

  TEAM_BETA: {
    azuread: { name: 'TEAM_BETA' },
    github: { name: 'team-beta' },
    description: 'Beta team members',
    include: [
      org.resources.team_beta_member_3,
      org.resources.team_beta_member_1,
      org.resources.team_beta_member_2,
    ],
  },
}`,
			expected: `{
  TEAM_ZEBRA: {
    azuread: { name: 'TEAM_ZEBRA' },
    github: { name: 'team-zebra' },
    description: 'Zebra team members',
    include: combineMembers(
      org.resources.team_zebra_members.include,
      org.resources.team_zebra_engineers.include,
    ),
  },

  TEAM_ALPHA: {
    azuread: { name: 'TEAM_ALPHA' },
    github: { name: 'team-alpha' },
    description: 'Alpha team members',
    include: combineMembers(
      org.resources.team_alpha_members.include,
      org.resources.team_alpha_engineers.include,
    ),
  },

  TEAM_BETA: {
    azuread: { name: 'TEAM_BETA' },
    github: { name: 'team-beta' },
    description: 'Beta team members',
    include: [
      org.resources.team_beta_member_1,
      org.resources.team_beta_member_2,
      org.resources.team_beta_member_3,
    ],
  },
}`,
		},
		{
			name: "InlineArrayStart",
			input: `local viewers = {
    PR: [ resource.zebra,
          resource.elephant,
          resource.antelope,
    ] + extra.PR,
    NP: [ resource.monkey,
          resource.giraffe,
          resource.bear,
    ] + extra.NP,
};`,
			expected: `local viewers = {
    PR: [ resource.antelope,
          resource.elephant,
          resource.zebra,
    ] + extra.PR,
    NP: [ resource.bear,
          resource.giraffe,
          resource.monkey,
    ] + extra.NP,
};`,
		},
		{
			name: "ObjectWithInlineArray",
			input: `{
  TEAM_PLATFORM_CODEOWNERS: {
    azuread: { name: 'TEAM_PLATFORM_CODEOWNERS' },
    github: { name: 'team-platform-services' },
    description: 'Platform Engineers who are codeowners for specific areas of the code',
    include: [ experience.resources.platform_team ],
  },
}`,
			expected: `{
  TEAM_PLATFORM_CODEOWNERS: {
    azuread: { name: 'TEAM_PLATFORM_CODEOWNERS' },
    github: { name: 'team-platform-services' },
    description: 'Platform Engineers who are codeowners for specific areas of the code',
    include: [ experience.resources.platform_team ],
  },
}`,
		},
		{
			name: "MultipleObjectsWithInlineArrays",
			input: `{
  OUTER_OBJECT: {
    TEAM_ONE: {
      azuread: { name: 'TEAM_ONE_AD' },
      github: { name: 'team-one-gh' },
      description: 'Team one description',
      include: [ resources.team_one ],
    },
    TEAM_TWO: {
      azuread: { name: 'TEAM_TWO_AD' },
      github: { name: 'team-two-gh' },
      description: 'Team two description',
      include: [ resources.team_two ],
    },
  },
}`,
			expected: `{
  OUTER_OBJECT: {
    TEAM_ONE: {
      azuread: { name: 'TEAM_ONE_AD' },
      github: { name: 'team-one-gh' },
      description: 'Team one description',
      include: [ resources.team_one ],
    },
    TEAM_TWO: {
      azuread: { name: 'TEAM_TWO_AD' },
      github: { name: 'team-two-gh' },
      description: 'Team two description',
      include: [ resources.team_two ],
    },
  },
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SortJsonnet(tt.input)
			if err != nil {
				t.Fatalf("SortJsonnet() error = %v", err)
			}
			if result != tt.expected {
				t.Errorf("SortJsonnet() =\n%q\n\nwant:\n%q", result, tt.expected)
			}
		})
	}
}

func TestExtractSortKey(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "SimpleIdentifier",
			input:    "  alpha,",
			expected: "alpha",
		},
		{
			name:     "WithHashComment",
			input:    "  alpha,  # Comment here",
			expected: "alpha",
		},
		{
			name:     "WithSlashComment",
			input:    "  alpha,  // Comment here",
			expected: "alpha",
		},
		{
			name:     "NestedPath",
			input:    "  zoo.mammals.elephant,",
			expected: "zoo.mammals.elephant",
		},
		{
			name:     "UppercaseConvertedToLowercase",
			input:    "  ALPHA,",
			expected: "alpha",
		},
		{
			name:     "MixedCase",
			input:    "  AlPhA,",
			expected: "alpha",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractSortKey(tt.input)
			if result != tt.expected {
				t.Errorf("extractSortKey(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCleanupWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "TrailingSpaces",
			input:    "line one  \nline two\t\nline three   ",
			expected: "line one\nline two\nline three\n",
		},
		{
			name:     "MultipleBlankLines",
			input:    "line one\n\n\n\nline two",
			expected: "line one\n\n\nline two\n",
		},
		{
			name:     "NoTrailingNewline",
			input:    "line one\nline two",
			expected: "line one\nline two\n",
		},
		{
			name:     "AlreadyHasTrailingNewline",
			input:    "line one\nline two\n",
			expected: "line one\nline two\n",
		},
		{
			name:     "EmptyString",
			input:    "",
			expected: "",
		},
		{
			name:     "OnlyWhitespace",
			input:    "   \n\t\n  ",
			expected: "\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanupWhitespace(tt.input)
			if result != tt.expected {
				t.Errorf("CleanupWhitespace() =\n%q\n\nwant:\n%q", result, tt.expected)
			}
		})
	}
}

func TestSortBlock(t *testing.T) {
	tests := []struct {
		name     string
		input    []arrayElement
		expected []string
	}{
		{
			name: "SimpleSort",
			input: []arrayElement{
				{original: "  charlie,", sortKey: "charlie"},
				{original: "  alpha,", sortKey: "alpha"},
				{original: "  bravo,", sortKey: "bravo"},
			},
			expected: []string{
				"  alpha,",
				"  bravo,",
				"  charlie,",
			},
		},
		{
			name: "PreserveOriginalFormatting",
			input: []arrayElement{
				{original: "    zebra,  # Z comment", sortKey: "zebra"},
				{original: "    alpha,  # A comment", sortKey: "alpha"},
			},
			expected: []string{
				"    alpha,  # A comment",
				"    zebra,  # Z comment",
			},
		},
		{
			name: "SingleElement",
			input: []arrayElement{
				{original: "  single,", sortKey: "single"},
			},
			expected: []string{
				"  single,",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sortBlock(tt.input)
			if len(result) != len(tt.expected) {
				t.Fatalf("sortBlock() returned %d lines, want %d", len(result), len(tt.expected))
			}
			for i, line := range result {
				if line != tt.expected[i] {
					t.Errorf("sortBlock()[%d] = %q, want %q", i, line, tt.expected[i])
				}
			}
		})
	}
}
