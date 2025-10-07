# CodeCheck

A CLI tool to streamline code review workflows for student assignments. CodeCheck automates cloning repositories, installing dependencies, opening projects in VS Code, and starting development servers.

## Features

-  Clone GitHub repositories with one command
-  Automatic dependency installation once package.json is present
-  Opens projects directly in VS Code (if installed) for viewing code
-  Auto-detects project type and starts appropriate dev server
-  Simple cleanup workflow after review
-  Configurable download directory

## Prerequisites

- [Go](https://golang.org/dl/) 1.25.1 or higher
- [Git](https://git-scm.com/)
- [Node.js & npm](https://nodejs.org/) (for Node.js projects)
- [VS Code](https://code.visualstudio.com/) with CLI (`code` command)

## Installation

```bash
go install github.com/mentalcaries/codecheck/cmd/codecheck@latest
```

## Usage

### Initial Setup

You'll be prompted to set your download directory. This can be anywhere you'd like the projects to be cloned.

If you need to change the config afterwards, you can run:

```bash
codecheck setup
```

### Review a Project

Supports GitHub URLs - HTTP and SSH. Usage:

```bash
codecheck review <github-url>
```

**Examples:**
```bash
codecheck review https://github.com/student/assignment-1
codecheck review https://github.com/student/project.git
codecheck review git@github.com:student/assignment.git
```

**What happens:**
1. Clones the repository to your configured directory
2. Detects project type (static HTML or Node.js)
3. Installs dependencies (if Node.js project)
4. Opens project in VS Code
5. Starts appropriate development server and opens your browser (for frontend projects)
6. Waits for you to complete your review

**When you're done reviewing:**
- Press `Ctrl+C` to stop the server
- Choose to delete or keep the project directory

## Supported Project Types

- **Static HTML/CSS/JS** - Serves with built-in Go file server
- **Webpack** - Runs `npm run dev` with auto port configuration
- **Vite** - Runs `npm run dev` (tested with React)
- **Node.js/Express** - Runs `npm run dev`

## Directory Conflict Resolution

If a directory with the same name already exists, you'll be prompted with options:
- **[Enter]** - Delete and overwrite existing directory
- **[n]** - Clone with modified name (appends username)
- **[q]** - Cancel operation


## Known Limitations

- Only supports GitHub URLs (HTTPS and SSH)
- Webpack projects may require pressing `Ctrl+C` twice for graceful shutdown
- Requires VS Code CLI to be configured

## Development

### Clone and Build Locally

```bash
git clone https://github.com/mentalcaries/codecheck
cd codecheck
go install ./cmd/codecheck
```

### Run Without Installing

```bash
go run ./cmd/codecheck review <github-url>
```

## Troubleshooting

**"Repository not found or is private"**
- Verify the URL is correct
- Ensure the repository is set to public


**"VS Code not available"**
- Install VS Code CLI: Open VS Code → Command Palette (`Cmd+Shift+P`) → "Shell Command: Install 'code' command in PATH"

**Port conflicts**
- CodeCheck serves static html files from Port 5543, Vite auto-assigns and WebPack uses its config
- If issues persist, manually stop any processes using the ports or manually change the config

## License

MIT


---

**Happy reviewing! 🎉**