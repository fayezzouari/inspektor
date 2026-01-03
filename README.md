# Inspektor 🔍

An AI-powered CLI tool for process inspection and system monitoring. Inspektor analyzes running processes, collects system metrics, and generates intelligent warnings about potential issues using Google's Gemini AI.

## Features

- **Process Analysis**: Detailed information about any running process by PID
- **Resource Monitoring**: CPU, memory, file descriptors, and network connections
- **System Health**: Overall system resource usage and health metrics
- **AI-Powered Analysis**: Intelligent warnings and recommendations using Gemini AI
- **Fallback Analysis**: Rule-based analysis when AI is unavailable
- **Rich Terminal Output**: Beautiful, color-coded display with visual indicators
- **JSON Output**: Machine-readable format for automation and integration

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

## Configuration

### AI Features (Optional)

To enable AI-powered analysis, you'll need a Gemini API key:

1. Get your API key from [Google AI Studio](https://makersuite.google.com/app/apikey)
2. Create a `.env` file in the project root:

```bash
cp .env.example .env
# Edit .env and add your API key
GEMINI_API_KEY=your_gemini_api_key_here
```

**Note**: If no API key is provided, Inspektor will automatically fall back to rule-based analysis.

## Usage

```bash
# Inspect a process by PID
./inspektor 1234

# With verbose output
./inspektor -v 1234

# JSON output format
./inspektor -j 1234

# Get help
./inspektor --help
```

## Example Output

### Rich Terminal Display
```
INSPEKTOR - Process 1234 (nginx)

────────────────────────────────────────────────────────────

 PROCESS 
        Status: Running
       Command: nginx: master process /usr/sbin/nginx
    Executable: /usr/sbin/nginx
   Working Dir: /
       Started: Jan 15, 10:30:45

 RESOURCES 
     CPU Usage: 2.5%
        Memory: 45.2 MB (0.8%)
Virtual Memory: 123.4 MB
    Open Files: 12
   Connections: 8
Child Processes: 4

 SYSTEM 
           CPU: 8 cores, 15.3%
        Memory: 8.2 GB / 16.0 GB (51.2%)
     CPU Model: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz

✓ All systems healthy
```

### AI-Powered Warnings
When issues are detected, Inspektor provides intelligent analysis with warnings and recommendations:

```
 WARNINGS 

  1. ⚠ High memory usage detected - process consuming 15% of system memory
  2. ⚠ System memory at 85% - risk of OOM killer terminating processes

 RECOMMENDATIONS 

  1. → Set memory limits using systemd (MemoryMax=2G) to prevent system-wide impact
  2. → Monitor with 'vmstat 1' to track memory pressure patterns
  3. → Consider adding swap space or increasing RAM for better stability
  4. → Enable process monitoring with systemd watchdog for auto-restart capability
```

The AI analysis includes:
- **Warnings**: Critical issues requiring immediate attention
- **Recommendations**: Preventive measures and best practices to avoid future problems
- **Specific Commands**: Actionable steps with exact commands to run
- **Context-Aware**: Considers process type and system patterns for intelligent analysis

## AI vs Rule-Based Analysis

- **With AI (Gemini)**: Context-aware analysis that considers process type, system patterns, and provides nuanced recommendations
- **Without AI**: Fast rule-based analysis using predefined thresholds and heuristics
- **Automatic Fallback**: Seamlessly switches to rule-based analysis if AI is unavailable

## Dependencies

- [gopsutil](https://github.com/shirou/gopsutil) - Cross-platform system and process utilities
- [cobra](https://github.com/spf13/cobra) - CLI framework
- [lipgloss](https://github.com/charmbracelet/lipgloss) - Rich terminal styling
- [generative-ai-go](https://github.com/google/generative-ai-go) - Google Gemini AI SDK
- [godotenv](https://github.com/joho/godotenv) - Environment variable loading

## License

MIT License