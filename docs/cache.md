# Cache System

## Overview

The git-commit-scribe cache system stores AI-generated commit messages to avoid redundant API calls and speed up repeated operations. The cache uses a **multi-level architecture** with both memory and disk layers for optimal performance.

## Architecture

### Multi-Level Caching

The cache system employs a two-tier approach:

1. **Memory Cache (L1)** - Fast, in-memory storage with LRU eviction
2. **Disk Cache (L2)** - Persistent, compressed file storage

When retrieving cached data:
- Check memory cache first (fastest)
- If not found, check disk cache
- If found on disk, promote to memory cache for future speed
- If not found anywhere, generate with AI and store in both layers

### Cache Key Generation

Cache keys are generated using SHA256 hashing of:
- Git diff content
- Commit type (if specified)
- Commit scope (if specified)
- Breaking change flag
- One-line format flag
- Description length limit
- AI model name

This ensures that identical changes with identical options will retrieve the same cached result.

## Configuration

Configure the cache in `~/.config/gscribe/config.yaml`:

```yaml
cache:
  enabled: true           # Enable/disable caching (default: true)
  ttl: 86400             # Time-to-live in seconds (default: 24 hours)
  max_memory_mb: 100     # Max memory cache size in MB (default: 100)
  max_disk_mb: 1000      # Max disk cache size in MB (default: 1000)
```

### Configuration Options

- **enabled**: Toggle caching on/off
- **ttl**: How long cached entries remain valid (in seconds)
  - Default: 86400 (24 hours)
  - Set to 0 for no expiration
- **max_memory_mb**: Maximum size of in-memory cache
  - Older entries are evicted when limit is reached
- **max_disk_mb**: Maximum size of disk cache
  - Older files are deleted when limit is reached

## Cache Storage

### Location

Cache files are stored in the XDG-compliant cache directory:
- **Linux/macOS**: `~/.cache/gscribe/`
- **Custom**: Set `$XDG_CACHE_HOME` environment variable

### File Format

- Each cached entry is stored as a compressed gzip file: `{sha256-hash}.gz`
- Files contain JSON-encoded commit results with timestamps
- Compression typically reduces file size by 60-80%

### Data Structure

Each cache entry stores:
```json
{
  "value": "{\"type\":\"feat\",\"scope\":\"cache\",\"description\":\"...\",\"body\":\"...\",\"breaking\":false}",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## CLI Commands

### View Cache Statistics

```bash
gscribe cache stats
```

Shows:
- Cache directory location
- Cache configuration settings
- Number of cached entries
- Total disk size
- Recent cache entries with timestamps and sizes

Example output:
```
Cache Statistics
═══════════════════════════════════════

Cache directory: /Users/username/.cache/gscribe

Cache settings:
  Enabled: true
  TTL: 86400 seconds
  Max memory: 100 MB
  Max disk: 1000 MB

Cache entries: 42
Total size: 2.34 MB

Recent cache entries:
  • a3f2b1c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0.gz (12.45 KB, modified 2024-01-15 10:30:00)
  • b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2c3.gz (8.92 KB, modified 2024-01-15 09:15:00)
  ...
```

### Clear Cache

```bash
# Clear disk cache
gscribe cache clear

# Clear both memory and disk cache
gscribe cache clear --all
```

## Runtime Behavior

### Cache Hits

When a cache hit occurs:
1. Message retrieval is nearly instantaneous
2. No API call is made (saves API quota and costs)
3. In verbose mode (`-v`), shows "✓ Using cached commit message"

Example:
```bash
$ gscribe commit -v
Getting staged changes...
Analyzing 3 staged file(s)...
Cache initialized
✓ Using cached commit message

Generated commit message:
─────────────────────────────────────────
feat(cache): integrate cache system

Implement multi-level caching with memory and disk layers
to avoid redundant AI API calls and improve performance.
─────────────────────────────────────────
```

### Cache Misses

When no cached entry is found:
1. AI provider is called to generate the message
2. Result is stored in both memory and disk cache
3. Subsequent identical diffs will use the cached result

### Expiration

Cached entries expire based on the TTL setting:
- Memory cache checks expiration on every `Get()` call
- Disk cache checks expiration and deletes expired files
- Expired entries are silently removed and regenerated as needed

### Eviction

When cache limits are reached:
- **Memory**: LRU (Least Recently Used) eviction removes oldest entries
- **Disk**: Oldest files are deleted based on modification time

## Disabling Cache

### Temporary (per-command)

Use the `--no-cache` flag:
```bash
gscribe commit --no-cache
```

### Permanent

Set `cache.enabled: false` in config:
```bash
gscribe config set cache.enabled false
```

Or edit `~/.config/gscribe/config.yaml`:
```yaml
cache:
  enabled: false
```

## Performance Impact

### With Cache Enabled

- **First run**: ~2-5 seconds (AI API call)
- **Subsequent runs**: <100ms (cache retrieval)
- **API savings**: ~100% for identical diffs

### Disk Space

- Typical cache entry: 5-20 KB (compressed)
- 1000 cached entries: ~10-20 MB
- Configurable max size prevents unbounded growth

## Best Practices

1. **Keep cache enabled** - Dramatically improves performance for repeated operations
2. **Set appropriate TTL** - Balance freshness vs. cache effectiveness
   - Development: 3600 seconds (1 hour)
   - Production: 86400 seconds (24 hours)
3. **Monitor cache size** - Use `gscribe cache stats` to check usage
4. **Clear stale data** - Run `gscribe cache clear` periodically if needed
5. **Adjust memory limits** - Increase for better performance on machines with available RAM

## Troubleshooting

### Cache not working

1. Check if cache is enabled:
   ```bash
   gscribe config show | grep cache
   ```

2. Verify cache directory exists and is writable:
   ```bash
   ls -la ~/.cache/gscribe/
   ```

3. Check cache statistics:
   ```bash
   gscribe cache stats
   ```

4. Try with verbose mode to see cache behavior:
   ```bash
   gscribe commit -v
   ```

### High disk usage

Reduce max disk size:
```bash
gscribe config set cache.max_disk_mb 500
```

Or clear the cache:
```bash
gscribe cache clear
```

### Stale cached results

Reduce TTL:
```bash
gscribe config set cache.ttl 3600  # 1 hour
```

Or clear the cache:
```bash
gscribe cache clear
