# Prompt Generator for AI

A command-line tool that scans a project folder and creates a single text file with your entire project's structure and contents. You can then paste that file directly into an AI assistant (like ChatGPT or Claude) to provide context.

This tool is written in Go and runs on Windows, Linux, and macOS. It requires no dependencies—just a single executable file.

## How It Works

When you run the program:

- It asks you which folder to scan (you can press Enter for the current directory, provide an absolute path, or type a folder name that sits next to the executable).
- It walks through all folders and files, showing the progress in color on the terminal.
- Files that are binary, hidden, too large, or match your ignore rules are skipped (their content is not included in the output).
- The result is a `.txt` file that contains two parts:
  1. A tree view of your folder structure.
  2. The full text of every included file, each preceded by a header with its relative path.

The output file is named automatically with a timestamp (like `project_prompt_2026-05-22_14-30-00.txt`), but you can specify a custom name if you prefer.

## Key Features

- **Smart skipping**: The program skips files that are obviously binary (like images, fonts, or executables) based on their extension, and also uses a content check to catch binary files with unexpected extensions.
- **Configurable ignores**: An automatically generated `config.json` lets you define patterns for folders and files that should be completely ignored (for example, `.git`, `node_modules`, or custom patterns).
- **Hidden files and size limits**: By default, files starting with a dot (`.`) are skipped, and files larger than 500 KB are left out. Both settings can be changed in the configuration file.
- **Cross-platform**: Works the same way on Windows, Linux, and macOS. Paths are normalized to forward slashes in the output.

## Getting Started

### 1. Download or Build

** Build from source**  
Make sure you have Go 1.18 or later installed, then run:

```bash
git clone https://github.com/MehdiShekari/prompt-generator-for-ai.git
cd prompt-generator-for-ai
go build -o prompt-generator .
```

### 2. Run the Tool

Execute the binary from a terminal:

```bash
./prompt-generator
```

You will see a banner and then a prompt:

```
Enter the target directory path (press Enter for current dir):
```

Type a path and hit Enter, or just press Enter to scan the directory where you are currently located.

### 3. Find Your Output

After scanning, the program creates a file named something like `project_prompt_2026-05-22_14-30-00.txt` in the current working directory. That file is ready to be shared with an AI.

## Customizing Behavior

When you first run the tool, it creates a `config.json` file next to the executable. You can open this file and adjust the following settings:

- `binary_extensions`: A list of file extensions that should be treated as binary (their content is never read).
- `ignore_patterns`: Glob patterns for files or folders to ignore completely (they won't appear in the structure or content). For example, `"*.log"` or `"temp/*"`.
- `ignore_hidden`: Set to `true` to skip all files and folders whose name begins with a dot.
- `max_file_size_kb`: Files larger than this number of kilobytes are skipped. Set to `0` to disable the size check.

After editing, simply run the program again; it will use your updated settings immediately.

## Command Line Flags (Optional)

If you prefer not to use the interactive prompt, you can give the folder as a command-line argument:

```bash
./prompt-generator /path/to/my/project
```

You can also specify an output file path and a custom configuration file:

```bash
./prompt-generator -output my_context.txt -config /path/to/my_config.json
```

For help, run:

```bash
./prompt-generator -h
```

## Output Example

Suppose you scan a tiny Go project. The generated file might look like this:

```
=== Project Folder Structure ===

myproject:
--------main.go
--------utils
------------helper.go

=== File Contents ===

--- main.go ---
package main

import "fmt"

func main() {
    fmt.Println("Hello, world!")
}

--- utils/helper.go ---
package utils

func Greet() string {
    return "Hi from utils"
}
```

## Troubleshooting

- **"directory not found"**: If you type a folder name that doesn't exist, the tool tries to find it relative to the executable and then relative to your current working directory. Use an absolute path if you're unsure.
- **Permission denied on push (for repository maintainers)**: Make sure your Git remote uses your own credentials (SSH or a personal access token) and not another user's cached login. See GitHub's documentation for setting up authentication.
- **Colors not showing**: This program uses ANSI escape codes. They work in most modern terminals (Windows Terminal, PowerShell, Linux/macOS terminals). Older command prompts like `cmd.exe` on Windows may not display them.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

## Author

Maintained by Mehdi Shekari  
GitHub: https://github.com/MehdiShekari/prompt-generator-for-ai
```
