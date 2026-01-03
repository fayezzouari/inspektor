# Inspektor 🔍

An AI-powered CLI tool for process inspection and system monitoring. Inspektor analyzes running processes, collects system metrics, and generates intelligent warnings about potential issues.

## Features

- **Process Analysis**: Detailed information about any running process by PID
- **Resource Monitoring**: CPU, memory, file descriptors, and network connections
- **System Health**: Overall system resource usage and health metrics
- **AI-Powered Warnings**: Intelligent analysis with actionable recommendations
- **Service Detection**: Automatic identification of common services and applications

## Installation

```bash
# Clone the repository
git clone <repository-url>
cd inspektor

# Install dependencies
go mod tidy

# Build the binary
go build -o inspektor

# Or install globally
go install
```

## Usage

```bash
# Inspect a process by PID
./inspektor 1234

# With verbose output
./inspektor -v 1234

# JSON output format
./inspektor -j 1234
```

## Example Output

```
🔍 Process Inspection Report
═══════════════════════════════════════════════════════════════

📋 Process Information:
  PID: 1234
  Name: nginx
  Executable: /usr/sbin/nginx
  Status: running
  Working Directory: /
  Created: 2024-01-15 10:30:45
  Command Line: nginx: master process /usr/sbin/nginx

📊 Resource Usage:
  CPU Usage: 2.5%
  Memory RSS: 45.2 MB
  Memory VMS: 123.4 MB
  Memory Usage: 0.8%
  Open Files: 12
  Network Connections: 8
  Child Processes: 4

🖥️  System Information:
  CPU Model: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
  CPU Cores: 8
  System CPU Usage: 15.3%
  Total Memory: 16.0 GB
  Used Memory: 8.2 GB (51.2%)
  Free Memory: 7.8 GB

⚠️  AI Analysis & Warnings:
═══════════════════════════════════════════════════════════════
• No warnings detected - system appears healthy
```

## Dependencies

- [gopsutil](https://github.com/shirou/gopsutil) - Cross-platform system and process utilities
- [cobra](https://github.com/spf13/cobra) - CLI framework

## License

MIT License