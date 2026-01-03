package analyzer

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"inspektor/internal/models"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

// AIAnalyzer provides intelligent analysis of system and process data using Gemini AI
type AIAnalyzer struct {
	client    *genai.Client
	model     *genai.GenerativeModel
	aiEnabled bool
}

func New() *AIAnalyzer {
	// Load environment variables
	_ = godotenv.Load()

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Println("Warning: GEMINI_API_KEY not found. AI analysis will use fallback rules.")
		return &AIAnalyzer{aiEnabled: false}
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Printf("Warning: Failed to initialize Gemini client: %v. Using fallback analysis.\n", err)
		return &AIAnalyzer{aiEnabled: false}
	}

	model := client.GenerativeModel("gemini-2.5-flash")
	model.SetTemperature(0.3) // Lower temperature for more consistent analysis

	return &AIAnalyzer{
		client:    client,
		model:     model,
		aiEnabled: true,
	}
}

// AnalyzeAndWarn generates warnings based on process and system metrics
func (a *AIAnalyzer) AnalyzeAndWarn(data *models.InspectionData) []string {
	if a.aiEnabled {
		return a.analyzeWithAI(data)
	}
	return a.analyzeWithRules(data)
}

func (a *AIAnalyzer) analyzeWithAI(data *models.InspectionData) []string {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	prompt := a.buildAnalysisPrompt(data)

	resp, err := a.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Printf("AI analysis failed: %v. Falling back to rule-based analysis.\n", err)
		return a.analyzeWithRules(data)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		log.Println("No AI response received. Falling back to rule-based analysis.")
		return a.analyzeWithRules(data)
	}

	// Parse AI response
	aiResponse := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
	return a.parseAIResponse(aiResponse)
}

func (a *AIAnalyzer) buildAnalysisPrompt(data *models.InspectionData) string {
	processAge := time.Since(data.Process.CreateTime)

	prompt := fmt.Sprintf(`You are a senior system administrator and DevOps expert analyzing a running process. Provide intelligent analysis with specific warnings and actionable recommendations.

PROCESS INFORMATION:
- PID: %d
- Name: %s
- Status: %s
- Command: %s
- Process Age: %s
- CPU Usage: %.2f%%
- Memory RSS: %s (%.2f%% of system)
- Memory VMS: %s
- Open Files: %d
- Network Connections: %d
- Child Processes: %d

SYSTEM CONTEXT:
- CPU Cores: %d
- System CPU Usage: %.2f%%
- Total Memory: %s
- Used Memory: %s (%.2f%%)
- Free Memory: %s

ANALYSIS GUIDELINES:

1. RESOURCE USAGE ASSESSMENT:
   - Evaluate if CPU/memory usage is appropriate for this process type
   - Consider normal vs abnormal patterns for system processes, web servers, databases, etc.
   - Flag resource exhaustion risks before they become critical

2. PROCESS HEALTH INDICATORS:
   - Check for zombie/stopped processes that need intervention
   - Assess if file descriptor or connection counts indicate leaks
   - Evaluate if child process count suggests fork bombs or runaway spawning

3. SYSTEM-WIDE IMPACT:
   - Consider how this process affects overall system stability
   - Flag if system resources are constrained and may cause OOM kills
   - Identify if the system needs scaling (vertical or horizontal)

4. PREVENTIVE MEASURES & BEST PRACTICES:
   - Suggest resource limits (ulimit, cgroups, systemd limits)
   - Recommend monitoring thresholds and alerting rules
   - Propose optimization strategies (memory tuning, connection pooling, etc.)
   - Advise on capacity planning if resources are trending toward limits
   - Suggest configuration improvements for common services

5. ACTIONABLE RECOMMENDATIONS:
   - Provide specific commands or configuration changes when applicable
   - Prioritize immediate actions vs long-term improvements
   - Include investigation steps for unclear issues

FORMAT YOUR RESPONSE:
- Each warning/recommendation on a separate line
- Start warnings with "WARNING:" for issues requiring attention
- Start recommendations with "RECOMMEND:" for preventive measures and best practices
- If no issues found, respond with "HEALTHY: No issues detected"
- Maximum 7 items total (warnings + recommendations)
- Order by priority: critical warnings first, then recommendations

EXAMPLES:

WARNING: High CPU usage (85%%) may indicate performance bottleneck or infinite loop
RECOMMEND: Set CPU limits using systemd (CPUQuota=80%%) to prevent system-wide impact
WARNING: Memory usage at 92%% - risk of OOM killer terminating processes
RECOMMEND: Add swap space or increase RAM; monitor with 'vmstat 1' for memory pressure
WARNING: 1500 open files detected - possible file descriptor leak
RECOMMEND: Investigate with 'lsof -p PID' and set ulimit -n to prevent exhaustion
RECOMMEND: Enable process monitoring with systemd watchdog or supervisord for auto-restart
RECOMMEND: Configure log rotation to prevent disk space exhaustion
HEALTHY: No issues detected

YOUR ANALYSIS:`,
		data.Process.PID,
		data.Process.Name,
		data.Process.Status,
		data.Process.CommandLine,
		processAge.Round(time.Second),
		data.Process.CPUPercent,
		formatBytes(data.Process.MemoryRSS),
		data.Process.MemoryPercent,
		formatBytes(data.Process.MemoryVMS),
		data.Process.OpenFiles,
		data.Process.Connections,
		data.Process.Children,
		data.System.CPUCores,
		data.System.CPUUsage,
		formatBytes(data.System.MemoryTotal),
		formatBytes(data.System.MemoryUsed),
		data.System.MemoryPercent,
		formatBytes(data.System.MemoryFree),
	)

	return prompt
}

func (a *AIAnalyzer) parseAIResponse(response string) []string {
	var warnings []string
	lines := strings.Split(response, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "WARNING:") {
			warning := strings.TrimSpace(strings.TrimPrefix(line, "WARNING:"))
			if warning != "" {
				warnings = append(warnings, "⚠ "+warning)
			}
		} else if strings.HasPrefix(line, "RECOMMEND:") {
			recommendation := strings.TrimSpace(strings.TrimPrefix(line, "RECOMMEND:"))
			if recommendation != "" {
				warnings = append(warnings, "→ "+recommendation)
			}
		} else if strings.HasPrefix(line, "HEALTHY:") {
			// If AI says it's healthy, return empty warnings
			return []string{}
		}
	}

	return warnings
}

// Fallback rule-based analysis (original implementation)
func (a *AIAnalyzer) analyzeWithRules(data *models.InspectionData) []string {
	var warnings []string

	// Analyze CPU usage
	warnings = append(warnings, a.analyzeCPU(data)...)

	// Analyze memory usage
	warnings = append(warnings, a.analyzeMemory(data)...)

	// Analyze process behavior
	warnings = append(warnings, a.analyzeProcess(data)...)

	// Analyze system health
	warnings = append(warnings, a.analyzeSystem(data)...)

	return warnings
}

func (a *AIAnalyzer) analyzeCPU(data *models.InspectionData) []string {
	var warnings []string

	// High process CPU usage
	if data.Process.CPUPercent > 80 {
		warnings = append(warnings, fmt.Sprintf(
			"High CPU usage detected: Process consuming %.2f%% CPU - investigate for performance bottlenecks",
			data.Process.CPUPercent))
	} else if data.Process.CPUPercent > 50 {
		warnings = append(warnings, fmt.Sprintf(
			"Moderate CPU usage: Process using %.2f%% CPU - monitor for sustained high usage",
			data.Process.CPUPercent))
	}

	// High system CPU usage
	if data.System.CPUUsage > 90 {
		warnings = append(warnings, fmt.Sprintf(
			"Critical system CPU load: %.2f%% usage - immediate attention required",
			data.System.CPUUsage))
	} else if data.System.CPUUsage > 75 {
		warnings = append(warnings, fmt.Sprintf(
			"High system CPU load: %.2f%% usage - consider load balancing",
			data.System.CPUUsage))
	}

	return warnings
}

func (a *AIAnalyzer) analyzeMemory(data *models.InspectionData) []string {
	var warnings []string

	// High process memory usage
	if data.Process.MemoryPercent > 10 {
		warnings = append(warnings, fmt.Sprintf(
			"High memory usage: Process using %.2f%% of system memory (%s RSS)",
			data.Process.MemoryPercent, formatBytes(data.Process.MemoryRSS)))
	}

	// Memory leak detection (simplified)
	if data.Process.MemoryVMS > data.Process.MemoryRSS*3 {
		warnings = append(warnings, fmt.Sprintf(
			"Potential memory leak: Virtual memory (%s) significantly exceeds RSS (%s)",
			formatBytes(data.Process.MemoryVMS), formatBytes(data.Process.MemoryRSS)))
	}

	// System memory pressure
	if data.System.MemoryPercent > 90 {
		warnings = append(warnings, fmt.Sprintf(
			"Critical memory pressure: System at %.2f%% - risk of OOM kills",
			data.System.MemoryPercent))
	} else if data.System.MemoryPercent > 80 {
		warnings = append(warnings, fmt.Sprintf(
			"High memory usage: System at %.2f%% - consider memory optimization",
			data.System.MemoryPercent))
	}

	return warnings
}

func (a *AIAnalyzer) analyzeProcess(data *models.InspectionData) []string {
	var warnings []string

	// Check process age
	processAge := time.Since(data.Process.CreateTime)
	if processAge < time.Minute {
		warnings = append(warnings, "Recently started process - monitor for stability during initialization")
	}

	// Check for zombie or stopped processes
	status := strings.ToLower(data.Process.Status)
	if status == "zombie" {
		warnings = append(warnings, "Zombie process detected - parent should reap this process")
	} else if status == "stopped" {
		warnings = append(warnings, "Process is currently stopped - may need manual intervention")
	}

	// High number of open files
	if data.Process.OpenFiles > 1000 {
		warnings = append(warnings, fmt.Sprintf(
			"High file descriptor usage: %d open files - check for file descriptor leaks",
			data.Process.OpenFiles))
	}

	// High number of network connections
	if data.Process.Connections > 100 {
		warnings = append(warnings, fmt.Sprintf(
			"High network connections: %d active connections - monitor for connection leaks",
			data.Process.Connections))
	}

	// Many child processes
	if data.Process.Children > 50 {
		warnings = append(warnings, fmt.Sprintf(
			"Many child processes: %d children - ensure proper process management",
			data.Process.Children))
	}

	return warnings
}

func (a *AIAnalyzer) analyzeSystem(data *models.InspectionData) []string {
	var warnings []string

	// Low core count with high usage
	if data.System.CPUCores <= 2 && data.System.CPUUsage > 60 {
		warnings = append(warnings, fmt.Sprintf(
			"Limited CPU resources: Only %d cores with %.2f%% usage - consider scaling up",
			data.System.CPUCores, data.System.CPUUsage))
	}

	// Low available memory
	freeMemoryPercent := float64(data.System.MemoryFree) / float64(data.System.MemoryTotal) * 100
	if freeMemoryPercent < 10 {
		warnings = append(warnings, fmt.Sprintf(
			"Low free memory: Only %.1f%% free (%s) - system may become unstable",
			freeMemoryPercent, formatBytes(data.System.MemoryFree)))
	}

	return warnings
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// Close cleans up the AI client
func (a *AIAnalyzer) Close() error {
	if a.client != nil {
		return a.client.Close()
	}
	return nil
}
