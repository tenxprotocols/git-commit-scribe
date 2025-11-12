# Example Custom Prompts

This directory contains example custom prompt templates for git-commit-scribe.

## Files

### `custom-prompt-simple.txt`
A minimal, straightforward prompt template. Good starting point for customization.

**Usage:**
```bash
gscribe commit --prompt-file examples/custom-prompt-simple.txt
```

### `custom-prompt-detailed.txt`
A comprehensive prompt template with detailed instructions and requirements. Produces more detailed commit messages with thorough explanations.

**Usage:**
```bash
gscribe commit --prompt-file examples/custom-prompt-detailed.txt
```

## Creating Your Own

To create your own custom prompt:

1. Copy one of the example files as a starting point
2. Modify the instructions and requirements to match your preferences
3. Test with `--dry-run` first:
   ```bash
   gscribe commit --prompt-file my-prompt.txt --dry-run
   ```
4. Use it for real commits once satisfied

## Template Variables

Your prompts can use these template variables:

- `{{.Diff}}` - The git diff content
- `{{.Type}}` - Specified commit type (via `--type`)
- `{{.Scope}}` - Specified scope (via `--scope`)
- `{{.Breaking}}` - Breaking change flag
- `{{.OneLine}}` - One-line mode flag
- `{{.MaxLength}}` - Max description length
- `{{.AvailableTypes}}` - Map of available types
- `{{.Temperature}}` - AI temperature setting

## Tips

1. **Always return JSON** - The tool expects JSON output with type, scope, description, body, and breaking fields
2. **Test incrementally** - Start simple and add complexity
3. **Be specific** - Clear instructions produce better results
4. **Reference available types** - Use `{{.AvailableTypes}}` to list valid commit types

See [docs/prompts.md](../docs/prompts.md) for full documentation.
