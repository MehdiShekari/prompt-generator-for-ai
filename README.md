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

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

## Author

Maintained by Mehdi Shekari  
GitHub: https://github.com/MehdiShekari/prompt-generator-for-ai

