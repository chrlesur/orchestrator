# Orchestrator

Orchestrator is a job and pipeline management and execution software developed in Go. It provides a text-based user interface and a REST API for remote control.

**ALPHA VERSION UNDER HEAVY DEVELOPMENT.**

## Features

- Job management: execution of system commands with parameters
- Pipeline management: definition and execution of job sequences
- Text-based User Interface (TUI)
- REST API for remote control
- Persistent data storage with BoltDB
- Plugin system for extending functionality
- Specific support for rtmscli and aiyou.cli tools

## Prerequisites

- Go 1.16 or higher
- BoltDB

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/your-username/orchestrator.git
   ```

2. Navigate to the project folder:
   ```
   cd orchestrator
   ```

3. Build the project:
   ```
   go build ./cmd/orchestrator
   ```

## Usage

Launch the application by running the compiled binary:

```
./orchestrator
```

### Available Commands

- `help`: Displays the list of available commands
- `addjob <name> <command> <arg1> <arg2> ...`: Adds a new job
- `addpipeline <id> <name> <job1> <job2> ...`: Adds a new pipeline
- `executeplugin <plugin_name> <arg1> <arg2> ...`: Executes a plugin
- `setloglevel <DEBUG|INFO|WARNING|ERROR>`: Sets the log level

## Configuration

The configuration file is located at `configs/config.yaml`. You can adjust parameters such as server port, database path, and default job parameters.

## Development

### Project Structure

- `cmd/`: Application entry point
- `internal/`: Internal packages (api, job, pipeline, ui, etc.)
- `pkg/`: Reusable packages (logger, utils, etc.)
- `plugins/`: Extensible plugins

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.

## License

GNU GPL 3.0

## Contact

[Christophe Lesur] - [christophe.lesur@cloud-temple.com]

Project Link: [https://github.com/chrlesur/orchestrator](https://github.com/chrlesur/orchestrator)
```