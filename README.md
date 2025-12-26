# jfmt

> JSON formatting from your terminal. Colors. Clipboard. Zero config.

```bash
brew install josharsh/tap/jfmt
```

---

```bash
$ echo '{"name":"john","active":true,"age":30}' | jfmt
```
```
{
  "name": "john",
  "active": true,
  "age": 30
}
```

---

### The magic trick

```bash
# Copy messy JSON from anywhere, then:
jfmt -C

# That's it. Formatted JSON is back in your clipboard.
```

---

### What it does

| | |
|---|---|
| `jfmt` | read from clipboard |
| `jfmt -C` | write back to clipboard |
| `jfmt file.json` | format a file |
| `jfmt -s` | sort keys |
| `jfmt -c` | compact/minify |
| `jfmt -f` | fix trailing commas |
| `cat x.json \| jfmt` | pipe |

---

### Install

```bash
# homebrew
brew install josharsh/tap/jfmt

# go
go install github.com/josharsh/jfmt@latest
```

---

### Web

[josharsh.github.io/jfmt](https://josharsh.github.io/jfmt)

---

MIT
