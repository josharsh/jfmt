# jfmt

**JSON formatting for humans.** Fast, colorful, zero config.

```bash
echo '{"name":"john","age":30}' | jfmt
```

```json
{
  "name": "john",
  "age": 30
}
```

## Why jfmt?

You format JSON 10 times a day. Your current workflow:

1. Google "json formatter"
2. Find a sketchy site with ads
3. Paste your JSON
4. Copy the result
5. Close tab

**jfmt** makes it instant:

```bash
jfmt                      # Format from clipboard
jfmt -C                   # Format and copy back
jfmt data.json            # Format a file
curl api.io | jfmt        # Format from pipe
jfmt https://api.io/data  # Fetch URL and format
```

## Features

- **Colored output** - Keys, strings, numbers, booleans in distinct colors
- **Clipboard integration** - Read from clipboard, copy result back
- **Auto-fix mode** - Fix trailing commas, single quotes
- **Sort keys** - Alphabetically sort object keys
- **Compact mode** - Minify instead of prettify
- **URL fetching** - Format JSON from any URL
- **Zero dependencies** - Single binary, no runtime needed

## Installation

### Homebrew (macOS/Linux)

```bash
brew install josharsh/tap/jfmt
```

### Direct Download

```bash
curl -fsSL https://raw.githubusercontent.com/josharsh/jfmt/main/install.sh | bash
```

### Go Install

```bash
go install github.com/josharsh/jfmt@latest
```

### Build from Source

```bash
git clone https://github.com/josharsh/jfmt.git
cd jfmt
make install
```

## Usage

```bash
# Basic formatting
echo '{"a":1}' | jfmt

# Format and copy to clipboard
jfmt -C < data.json

# No input? Reads from clipboard automatically
jfmt

# Compact output (minify)
jfmt -c data.json

# Sort keys alphabetically
jfmt -s data.json

# Fix common JSON issues (trailing commas, single quotes)
echo "{'name': 'test',}" | jfmt -f

# Fetch and format a URL
jfmt https://api.github.com/users/octocat

# Combine flags
jfmt -csC data.json   # compact, sorted, copy to clipboard
```

## Options

| Flag | Description |
|------|-------------|
| `-c` | Compact output (single line) |
| `-s` | Sort object keys alphabetically |
| `-C` | Copy result to clipboard |
| `-f` | Fix common JSON issues |
| `-m` | Monochrome (no colors) |
| `-h` | Show help |

## The Killer Workflow

Copy malformed JSON from anywhere, then:

```bash
jfmt -C
```

Done. Formatted JSON is now in your clipboard.

## Environment Variables

| Variable | Description |
|----------|-------------|
| `NO_COLOR` | Disable colored output |

## License

MIT

---

**jfmt** - because life's too short for browser-based JSON formatters.
