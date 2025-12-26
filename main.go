package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"golang.org/x/term"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorKey    = "\033[38;5;75m"  // Blue
	colorString = "\033[38;5;114m" // Green
	colorNumber = "\033[38;5;209m" // Orange
	colorBool   = "\033[38;5;170m" // Purple
	colorNull   = "\033[38;5;245m" // Gray
	colorBrace  = "\033[38;5;248m" // Light gray
)

var (
	flagCompact   bool
	flagSortKeys  bool
	flagCopy      bool
	flagFix       bool
	flagNoColor   bool
	flagHelp      bool
	flagMonochrome bool
)

func main() {
	args := parseArgs(os.Args[1:])

	if flagHelp {
		printUsage()
		os.Exit(0)
	}

	var input []byte
	var err error

	switch {
	case len(args) > 0 && strings.HasPrefix(args[0], "http"):
		// URL mode: fetch and format
		input, err = fetchURL(args[0])
	case len(args) > 0 && args[0] != "-":
		// File mode
		input, err = os.ReadFile(args[0])
	case !term.IsTerminal(int(os.Stdin.Fd())):
		// Pipe mode
		input, err = io.ReadAll(os.Stdin)
	default:
		// No input: read from clipboard
		input, err = readClipboard()
		if err != nil || len(bytes.TrimSpace(input)) == 0 {
			printUsage()
			os.Exit(1)
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	// Try to fix common JSON issues if -f flag is set
	if flagFix {
		input = fixJSON(input)
	}

	// Parse JSON
	var data any
	if err := json.Unmarshal(input, &data); err != nil {
		printJSONError(input, err)
		os.Exit(1)
	}

	// Sort keys if requested
	if flagSortKeys {
		data = sortKeys(data)
	}

	// Format output
	var output []byte
	if flagCompact {
		output, err = json.Marshal(data)
	} else {
		output, err = json.MarshalIndent(data, "", "  ")
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting JSON: %v\n", err)
		os.Exit(1)
	}

	// Colorize if outputting to terminal and colors not disabled
	outputStr := string(output)
	if !flagNoColor && !flagMonochrome && term.IsTerminal(int(os.Stdout.Fd())) {
		outputStr = colorize(outputStr)
	}

	fmt.Println(outputStr)

	// Copy to clipboard if requested
	if flagCopy {
		if err := writeClipboard(output); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not copy to clipboard: %v\n", err)
		} else {
			fmt.Fprintln(os.Stderr, "\033[38;5;245m(copied to clipboard)\033[0m")
		}
	}
}

func parseArgs(args []string) []string {
	var remaining []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "http") {
			for _, c := range arg[1:] {
				switch c {
				case 'c':
					flagCompact = true
				case 's':
					flagSortKeys = true
				case 'C':
					flagCopy = true
				case 'f':
					flagFix = true
				case 'm':
					flagMonochrome = true
				case 'h':
					flagHelp = true
				}
			}
		} else {
			remaining = append(remaining, arg)
		}
	}

	// Check NO_COLOR environment variable
	if os.Getenv("NO_COLOR") != "" {
		flagNoColor = true
	}

	return remaining
}

func printUsage() {
	usage := `jfmt - JSON formatter in a flash

Usage:
  jfmt [options] [file|url]
  echo '{"a":1}' | jfmt
  jfmt                      # read from clipboard

Options:
  -c    Compact output (single line)
  -s    Sort object keys alphabetically
  -C    Copy result to clipboard
  -f    Fix common JSON issues (trailing commas, single quotes)
  -m    Monochrome output (no colors)
  -h    Show this help

Examples:
  jfmt data.json                    # format file
  jfmt -s data.json                 # format with sorted keys
  curl api.io/data | jfmt           # format from pipe
  jfmt https://api.github.com/zen   # fetch and format URL
  jfmt -C                           # clipboard in, formatted + copied out
`
	fmt.Print(usage)
}

func fetchURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func readClipboard() ([]byte, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbpaste")
	case "linux":
		// Try xclip first, fall back to xsel
		cmd = exec.Command("xclip", "-selection", "clipboard", "-o")
	case "windows":
		cmd = exec.Command("powershell", "-command", "Get-Clipboard")
	default:
		return nil, fmt.Errorf("clipboard not supported on %s", runtime.GOOS)
	}
	return cmd.Output()
}

func writeClipboard(data []byte) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		cmd = exec.Command("xclip", "-selection", "clipboard")
	case "windows":
		cmd = exec.Command("powershell", "-command", "Set-Clipboard")
	default:
		return fmt.Errorf("clipboard not supported on %s", runtime.GOOS)
	}
	cmd.Stdin = bytes.NewReader(data)
	return cmd.Run()
}

func fixJSON(input []byte) []byte {
	s := string(input)

	// Remove trailing commas before } or ]
	re := regexp.MustCompile(`,(\s*[}\]])`)
	s = re.ReplaceAllString(s, "$1")

	// Replace single quotes with double quotes (simple cases)
	// This is naive but handles common cases
	re2 := regexp.MustCompile(`'([^']*)'(\s*:)`)
	s = re2.ReplaceAllString(s, `"$1"$2`)
	re3 := regexp.MustCompile(`:\s*'([^']*)'`)
	s = re3.ReplaceAllString(s, `: "$1"`)

	return []byte(s)
}

func sortKeys(v any) any {
	switch vv := v.(type) {
	case map[string]any:
		sorted := make(map[string]any)
		for k, val := range vv {
			sorted[k] = sortKeys(val)
		}
		return sorted
	case []any:
		for i, val := range vv {
			vv[i] = sortKeys(val)
		}
		return vv
	default:
		return v
	}
}

func colorize(s string) string {
	var result strings.Builder
	inString := false
	isKey := false
	i := 0

	for i < len(s) {
		c := s[i]

		switch {
		case c == '"':
			if inString {
				result.WriteString(string(c))
				result.WriteString(colorReset)
				inString = false
			} else {
				inString = true
				// Look ahead to determine if this is a key
				isKey = false
				for j := i + 1; j < len(s); j++ {
					if s[j] == '"' {
						// Check what comes after the closing quote
						for k := j + 1; k < len(s); k++ {
							if s[k] == ':' {
								isKey = true
								break
							} else if s[k] != ' ' && s[k] != '\t' && s[k] != '\n' {
								break
							}
						}
						break
					}
				}
				if isKey {
					result.WriteString(colorKey)
				} else {
					result.WriteString(colorString)
				}
				result.WriteString(string(c))
			}
		case inString:
			result.WriteString(string(c))
		case c == '{' || c == '}' || c == '[' || c == ']':
			result.WriteString(colorBrace)
			result.WriteString(string(c))
			result.WriteString(colorReset)
		case c >= '0' && c <= '9' || c == '-' || c == '.':
			result.WriteString(colorNumber)
			for i < len(s) && (s[i] >= '0' && s[i] <= '9' || s[i] == '.' || s[i] == '-' || s[i] == 'e' || s[i] == 'E' || s[i] == '+') {
				result.WriteString(string(s[i]))
				i++
			}
			result.WriteString(colorReset)
			continue
		case strings.HasPrefix(s[i:], "true"):
			result.WriteString(colorBool)
			result.WriteString("true")
			result.WriteString(colorReset)
			i += 4
			continue
		case strings.HasPrefix(s[i:], "false"):
			result.WriteString(colorBool)
			result.WriteString("false")
			result.WriteString(colorReset)
			i += 5
			continue
		case strings.HasPrefix(s[i:], "null"):
			result.WriteString(colorNull)
			result.WriteString("null")
			result.WriteString(colorReset)
			i += 4
			continue
		default:
			result.WriteString(string(c))
		}
		i++
	}

	return result.String()
}

func printJSONError(input []byte, err error) {
	fmt.Fprintf(os.Stderr, "\033[31mError:\033[0m Invalid JSON\n")

	// Try to extract position from error
	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		offset := int(syntaxErr.Offset)
		line, col := findPosition(input, offset)
		fmt.Fprintf(os.Stderr, "  at line %d, column %d\n\n", line, col)
		printErrorContext(input, offset)
	} else {
		fmt.Fprintf(os.Stderr, "  %v\n", err)
	}
}

func findPosition(input []byte, offset int) (line, col int) {
	line = 1
	col = 1
	for i := 0; i < offset && i < len(input); i++ {
		if input[i] == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}
	return
}

func printErrorContext(input []byte, offset int) {
	lines := strings.Split(string(input), "\n")
	line, col := findPosition(input, offset)

	// Print the relevant line
	if line <= len(lines) {
		lineContent := lines[line-1]
		fmt.Fprintf(os.Stderr, "  %s\n", lineContent)
		// Print the caret
		fmt.Fprintf(os.Stderr, "  %s\033[31m^\033[0m\n", strings.Repeat(" ", col-1))
	}
}
