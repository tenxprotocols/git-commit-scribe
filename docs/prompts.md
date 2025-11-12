# Custom Prompts

git-commit-scribe allows you to customize the AI prompt used to generate commit messages. This gives you control over the style, format, and requirements of your generated commits.

## Overview

The tool uses a prompt template system with Go's `text/template` syntax. You can provide custom prompts in three ways:

1. **Default embedded prompt** - Used automatically if no custom prompt is specified
2. **Inline prompt** - Pass prompt text directly on the command line
3. **Prompt file** - Load prompt from a file

## Default Prompt

The default prompt is embedded in the binary at `internal/ai/prompts/default.txt`. It includes:

- Instructions for following Conventional Commits v1.0.0
- Available commit types from your configuration
- Requirements for structure and formatting
- JSON response format specification

## Template Variables

Your prompt template has access to the following variables:

- `.Diff` - The git diff content
- `.Type` - Commit type (if specified via `--type`)
- `.Scope` - Commit scope (if specified via `--scope`)
- `.Breaking` - Boolean indicating breaking change
- `.OneLine` - Boolean for one-line commits
- `.MaxLength` - Maximum description length
- `.AvailableTypes` - Map of available commit types and descriptions
- `.Temperature` - AI temperature setting

## Using Custom Prompts

### Method 1: Inline Prompt

Pass the prompt directly on the command line:

```bash
gscribe commit --prompt "Generate a conventional commit message for this diff: {{.Diff}}"
```

### Method 2: Prompt File

Create a prompt template file and reference it:

```bash
# Create custom prompt file
cat > my-prompt.txt << 'EOF'
You are an expert at writing commit messages.

{{if .AvailableTypes}}Available types:
{{range $type, $desc := .AvailableTypes}}- {{$type}}: {{$desc}}
{{end}}
{{end}}

Generate a commit message in JSON format:
{
  "type": "feat",
  "scope": "",
  "description": "brief description",
  "body": "",
  "breaking": false
}

Diff:
```
{{.Diff}}
```
EOF

# Use the custom prompt
gscribe commit --prompt-file my-prompt.txt
```

### Priority

If both `--prompt` and `--prompt-file` are specified, `--prompt-file` takes precedence.

## Template Syntax

The prompt uses Go's `text/template` syntax:

### Conditionals

```
{{if .OneLine}}
- Generate ONLY a one-line commit message
{{else}}
- Include a body if needed
{{end}}
```

### Loops

```
{{range $type, $desc := .AvailableTypes}}
- {{$type}}: {{$desc}}
{{end}}
```

### Variables

```
{{.Diff}}
{{.Type}}
{{.MaxLength}}
```

## Example: Custom Style Prompt

Here's an example of a custom prompt that enforces a specific style:

```
You are an expert developer writing commit messages.

RULES:
- Use conventional commits format
- Keep descriptions under {{.MaxLength}} characters
- Use present tense, imperative mood
{{if .Type}}- Type must be: {{.Type}}{{end}}
{{if .Breaking}}- Mark as BREAKING CHANGE{{end}}

Available types:
{{range $type, $desc := .AvailableTypes}}- {{$type}}: {{$desc}}
{{end}}

Return JSON:
{
  "type": "type",
  "scope": "optional-scope",
  "description": "what this commit does",
  "body": "detailed explanation (optional)",
  "breaking": false
}

Changes:
```
{{.Diff}}
```

IMPORTANT: Be specific and technical in your description.
```

## Example: Simplified Prompt

For a more minimal approach:

```
Analyze this diff and create a conventional commit message.

Types: {{range $type, $desc := .AvailableTypes}}{{$type}}, {{end}}

Return JSON with type, scope, description, body, and breaking fields.

Diff:
```
{{.Diff}}
```
```

## Viewing the Default Prompt

The default prompt is embedded in the binary. To view it, you can extract it from the source:

```bash
cat internal/ai/prompts/default.txt
```

## Best Practices

1. **Always return JSON** - The tool expects a JSON response with the required fields
2. **Include available types** - Reference `{{.AvailableTypes}}` to ensure valid types
3. **Respect constraints** - Use the template variables to enforce settings like `{{.MaxLength}}`
4. **Be specific** - Provide clear instructions for the AI model
5. **Test your prompts** - Use `--dry-run` to test without creating commits

## Troubleshooting

### Prompt not being used

- Verify the file path is correct (use absolute paths if needed)
- Check that the file is readable

### Invalid template syntax

- Ensure you're using valid Go template syntax
- Test with a simple prompt first

### Poor results

- Try adjusting the prompt instructions
- Experiment with different levels of detail
- Consider the model's capabilities and limitations

## Environment Variables

While there's no environment variable for custom prompts directly, you can:

1. Create a wrapper script that calls `gscribe` with your custom prompt file
2. Use shell aliases for frequently-used custom prompts

```bash
# In your .bashrc or .zshrc
alias gscribe-custom='gscribe commit --prompt-file ~/.config/gscribe/my-prompt.txt'
