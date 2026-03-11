# jsonneat

A Go CLI tool that parses and sorts jsonnet files alphabetically.

## Installation

```bash
go install github.com/anicoll/jsonneat@latest
```

Or build from source:

```bash
make build
```

Or using go directly:

```bash
go build -o jsonneat
```

## Usage

```bash
jsonneat [-w] <path> [paths...]
```

### Options

- `-w` - Write result to source file instead of stdout

### Path Arguments

The tool accepts various path formats:

- **Specific file**: `jsonneat file.jsonnet`
- **Directory** (non-recursive): `jsonneat ./configs`
- **Recursive search**: `jsonneat ./...` or `jsonneat ./configs/...`

### Examples

Sort and print to stdout:
```bash
jsonneat example.jsonnet
```

Sort in-place:
```bash
jsonneat -w example.jsonnet
```

Sort all jsonnet files in a directory:
```bash
jsonneat -w ./configs
```

Recursively sort all jsonnet files:
```bash
jsonneat -w ./...
```

## What it does

The tool:
1. Sorts array elements in jsonnet files alphabetically
2. Preserves inline comments (e.g., `# Andrew Nicoll`)
3. Cleans up trailing whitespace
4. Ensures consistent file formatting

### Example

Input:
```jsonnet
local animals = [
  zebra,
  elephant,
  antelope,
];
```

Output:
```jsonnet
local animals = [
  antelope,
  elephant,
  zebra,
];
```

With comments:
```jsonnet
local animals = [
  tiger,  # Big cat
  elephant,  # Large mammal
  antelope,  # Swift runner
];
```

Becomes:
```jsonnet
local animals = [
  antelope,  # Swift runner
  elephant,  # Large mammal
  tiger,  # Big cat
];
```
