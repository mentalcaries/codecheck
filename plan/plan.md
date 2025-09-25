# CLI Tool - High Level Implementation Plan

## Core Components (in build order):

1. **CLI & URL Handling**
   - Parse GitHub URLs and extract repo/username info
   - Handle directory naming and conflicts

2. **Repository Cloning**
   - Execute git clone with proper error handling
   - Handle private repo errors gracefully

3. **Project Detection**
   - Check for package.json to determine project type
   - Static HTML vs Node.js projects

4. **Environment Setup**
   - Install dependencies (if needed)
   - Open VS Code
   - Start appropriate server (Go file server or npm dev server)

5. **Workflow Management**
   - Block and wait while user reviews
   - Handle cleanup when done

## Build Strategy:
- Linear progression from basic file operations to full development environment automation

## Key Technical Challenges:
- Port management across different project types
- Process lifecycle (starting/stopping servers)
- Cross-platform command execution

Simple linear progression from basic file operations to full development environment automation.


# CLI Tool Implementation Plan

## Phase 1: Core Infrastructure
**Goal**: Basic CLI structure and repository handling

### 1.1 CLI Setup
- Create Go module and main package
- Set up argument parsing (cobra or flag package)
- Accept GitHub URL as primary argument
- Add optional `-d/--dir` flag for custom directory names

### 1.2 URL Processing
- Extract repository name from GitHub URLs (both formats)
- Extract username from URL for conflict resolution
- Validate URL format (basic GitHub URL pattern matching)

### 1.3 Directory Management
- Implement directory name extraction from repo URL
- Check if target directory exists
- Build conflict resolution prompt (3 options: overwrite/rename/cancel)
- Implement username-appended naming (`repo-name-username`)

## Phase 2: Repository Operations
**Goal**: Clone repositories and handle errors

### 2.1 Git Integration
- Execute `git clone` using `os/exec`
- Handle git command failures (network, auth, not found)
- Implement error parsing for private/missing repos
- Display appropriate error messages

### 2.2 Error Handling
- Detect "repository not found" errors
- Create user-friendly error messages about private repos
- Handle network connectivity issues

## Phase 3: Project Detection
**Goal**: Identify project types and requirements

### 3.1 File System Analysis
- Check for `package.json` in root directory
- Implement recursive search for monorepos (future enhancement)
- Basic project type classification logic

### 3.2 Package Manager Detection
- Identify npm/yarn/pnpm from lock files
- Determine appropriate install command

## Phase 4: Development Environment Setup
**Goal**: Install dependencies and configure environment

### 4.1 Package Installation
- Execute package manager install commands
- Handle installation failures and errors
- Display installation progress/status

### 4.2 VS Code Integration
- Execute `code .` command to open project
- Handle cases where VS Code CLI is not available

## Phase 5: Server Implementation
**Goal**: Start appropriate development servers

### 5.1 Port Management
- Implement port availability checking (starting at 3000)
- Auto-increment port finder (3000 → 3001 → 3002...)
- Display chosen port to user

### 5.2 Static File Server (Go)
- Build HTTP file server using `net/http`
- Configure to serve `index.html` by default
- Implement graceful shutdown on Ctrl+C

### 5.3 Node.js Project Server
- Execute `npm run dev` with port override
- Handle webpack `--port` flag specifically  
- Handle Vite default behavior
- Parse and display server startup messages

## Phase 6: Workflow Management
**Goal**: Complete review lifecycle

### 6.1 Process Management
- Implement blocking execution model
- Handle Ctrl+C signal for graceful shutdown
- Ensure proper cleanup of child processes

### 6.2 Cleanup Operations
- Prompt user for directory deletion
- Implement safe directory removal
- Handle cleanup errors

## Implementation Order Priority

### Minimum Viable Product (MVP)
1. Basic CLI argument parsing
2. URL validation and repo name extraction
3. Git clone functionality with error handling
4. Simple directory conflict resolution
5. Basic static file server for HTML projects

### Iteration 2
6. Package.json detection and npm install
7. VS Code integration
8. Port management system

### Iteration 3
9. Node.js project server execution
10. Webpack port override handling
11. Cleanup workflow

### Future Enhancements
12. Monorepo support
13. Additional project type detection
14. Configuration file support

## Technical Decisions to Make

### CLI Framework
- **Decision needed**: Use `flag` package vs `cobra` vs `urfave/cli`
- **Consideration**: Cobra offers more features, flag is simpler

### Port Checking
- **Decision needed**: How to check if port is available
- **Options**: Attempt to bind, use netstat, or network scanning

### Process Management
- **Decision needed**: How to handle child process lifecycle
- **Consideration**: Ensure npm processes are properly terminated

### Error Handling Strategy
- **Decision needed**: How verbose should error messages be
- **Consideration**: Balance between helpful and overwhelming

## Potential Challenges

1. **Cross-platform compatibility**: Commands like `code .` may vary
2. **Node version compatibility**: Different projects may need different Node versions
3. **Port conflicts**: Multiple instances running simultaneously
4. **Process cleanup**: Ensuring child processes don't become orphaned
5. **Package manager variations**: Handling yarn vs npm vs pnpm consistently

## Success Criteria
- Clone any public GitHub repo successfully
- Handle directory conflicts gracefully
- Automatically detect and serve HTML or Node.js projects
- Consistent port usage across all project types
- Clean shutdown and optional cleanup