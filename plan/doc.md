# CLI Tool for Code Review - Project Summary

## Purpose
A CLI tool to streamline the code review process for student assignments by automating repository cloning, dependency installation, and development environment setup.

## Core Functionality

### Repository Handling
- **Input**: GitHub URLs in either format:
  - `https://github.com/username/repo-name`
  - `https://github.com/username/repo-name.git`
- **Assumption**: All repositories are public (no authentication needed)
- **Cloning**: Uses `git clone` directly (both URL formats work)

### Directory Management
- **Default directory**: Extract repo name from URL (e.g., `assignment-1`)
- **Conflict resolution**: When directory already exists, present 3 options:
  1. Overwrite existing directory
  2. Clone to `{repo-name}-{username}` (e.g., `assignment-1-alice123`)
  3. Cancel operation
- **Optional flag**: Allow custom directory specification (e.g., `-d custom-name`)

### Error Handling
- **Private/Non-existent repos**: When `git clone` fails (usually with a 404-like error)
  - Display clear error message: "Repository not found or is private"
  - Suggest student should:
    - Check the URL is correct
    - Make sure the repository is set to public
    - Verify they've shared the right link
- **Network issues**: Handle cases where git clone fails due to connectivity

### Project Detection & Setup
- **Detection logic**: Check for `package.json` in root directory
  - **No `package.json`**: Vanilla HTML/CSS/JS project → Use Go's `http.FileServer`
  - **Has `package.json`**: Node.js-based project → Install dependencies and run appropriate dev server
- **Monorepo support**: Recursively search subdirectories for `package.json` files
- **Package installation**: Detect and run appropriate package manager (npm/yarn/pnpm)
- **Supported project types**:
  - HTML/CSS (Go file server)
  - HTML/CSS/JS (Go file server)
  - HTML/CSS/JS with Webpack (npm run dev with --port flag)
  - React with Vite (npm run dev)
  - Node.js with Nodemon

### Port Management
- **Consistent port**: All projects serve on port 3000 by default
- **Auto-increment**: If port 3000 is busy, automatically try 3001, 3002, etc.
- **Port override**: For webpack projects, use `--port` flag to override config
- **User feedback**: Display chosen port (e.g., "Starting dev server on http://localhost:3001")

### Workflow
**Sequential execution (blocking):**
1. Clone repository to appropriate directory
2. Install packages (if `package.json` exists)
3. Open project in VS Code (`code .` command)
4. Start appropriate server:
   - **Static projects**: Go's `http.FileServer` serving `index.html`
   - **Node projects**: Execute `npm run dev` (or equivalent)
5. **Tool waits** while server runs and user reviews code
6. When user stops server (Ctrl+C), tool continues
7. Prompt to delete project directory
8. Clean up and exit

### Development Environment
- **Editor**: Open project in VS Code
- **Server types**:
  - Go-based file server for static HTML projects
  - NPM scripts for Node.js-based projects
- **Cleanup**: Optional directory deletion after review completion

## Use Case Context
- **Primary user**: Instructor reviewing student code submissions
- **Workflow**: Clone → inspect/run → provide feedback → delete
- **Benefits**: Single terminal manages entire review lifecycle with natural stopping points