package analyzer

import (
	"fmt"
	"strings"
	"time"

	"inspektor/internal/models"
)

type AIAnalyzer struct {
	// For now, we'll implement rule-based analysis
}

func New() *AIAnalyzer {
	return &AIAnalyzer{}
}

// AnalyzeAndWarn generates warnings based on process and system metrics
func (a *AIAnalyzer) AnalyzeAndWarn(data *models.InspectionData) []string {
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
			"HIGH CPU USAGE: Process is consuming %.2f%% CPU - consider investigating for performance bottlenecks",
			data.Process.CPUPercent))
	} else if data.Process.CPUPercent > 50 {
		warnings = append(warnings, fmt.Sprintf(
			"MODERATE CPU USAGE: Process is using %.2f%% CPU - monitor for sustained high usage",
			data.Process.CPUPercent))
	}
	
	// High system CPU usage
	if data.System.CPUUsage > 90 {
		warnings = append(warnings, fmt.Sprintf(
			"CRITICAL SYSTEM CPU: System CPU usage at %.2f%% - immediate attention required",
			data.System.CPUUsage))
	} else if data.System.CPUUsage > 75 {
		warnings = append(warnings, fmt.Sprintf(
			"HIGH SYSTEM CPU: System CPU usage at %.2f%% - consider load balancing",
			data.System.CPUUsage))
	}
	
	return warnings
}

func (a *AIAnalyzer) analyzeMemory(data *models.InspectionData) []string {
	var warnings []string
	
	// High process memory usage
	if data.Process.MemoryPercent > 10 {
		warnings = append(warnings, fmt.Sprintf(
			"HIGH MEMORY USAGE: Process is using %.2f%% of system memory (%s RSS)",
			data.Process.MemoryPercent, formatBytes(data.Process.MemoryRSS)))
	}
	
	// Memory leak detection (simplified)
	if data.Process.MemoryVMS > data.Process.MemoryRSS*3 {
		warnings = append(warnings, fmt.Sprintf(
			"POTENTIAL MEMORY LEAK: Virtual memory (%s) significantly exceeds RSS (%s)",
			formatBytes(data.Process.MemoryVMS), formatBytes(data.Process.MemoryRSS)))
	}
	
	// System memory pressure
	if data.System.MemoryPercent > 90 {
		warnings = append(warnings, fmt.Sprintf(
			"CRITICAL MEMORY PRESSURE: System memory usage at %.2f%% - risk of OOM kills",
			data.System.MemoryPercent))
	} else if data.System.MemoryPercent > 80 {
		warnings = append(warnings, fmt.Sprintf(
			"HIGH MEMORY USAGE: System memory at %.2f%% - consider memory optimization",
			data.System.MemoryPercent))
	}
	
	return warnings
}

func (a *AIAnalyzer) analyzeProcess(data *models.InspectionData) []string {
	var warnings []string
	
	// Check process age
	processAge := time.Since(data.Process.CreateTime)
	if processAge < time.Minute {
		warnings = append(warnings, "RECENTLY STARTED: Process started less than a minute ago - monitor for stability")
	}
	
	// Check for zombie or stopped processes
	status := strings.ToLower(data.Process.Status)
	if status == "zombie" {
		warnings = append(warnings, "ZOMBIE PROCESS: Process is in zombie state - parent should reap it")
	} else if status == "stopped" {
		warnings = append(warnings, "STOPPED PROCESS: Process is currently stopped")
	}
	
	// High number of open files
	if data.Process.OpenFiles > 1000 {
		warnings = append(warnings, fmt.Sprintf(
			"HIGH FILE DESCRIPTOR USAGE: Process has %d open files - check for file descriptor leaks",
			data.Process.OpenFiles))
	}
	
	// High number of network connections
	if data.Process.Connections > 100 {
		warnings = append(warnings, fmt.Sprintf(
			"HIGH NETWORK CONNECTIONS: Process has %d active connections - monitor for connection leaks",
			data.Process.Connections))
	}
	
	// Many child processes
	if data.Process.Children > 50 {
		warnings = append(warnings, fmt.Sprintf(
			"MANY CHILD PROCESSES: Process has %d children - ensure proper process management",
			data.Process.Children))
	}
	
	return warnings
}

func (a *AIAnalyzer) analyzeSystem(data *models.InspectionData) []string {
	var warnings []string
	
	// Low core count with high usage
	if data.System.CPUCores <= 2 && data.System.CPUUsage > 60 {
		warnings = append(warnings, fmt.Sprintf(
			"LIMITED CPU RESOURCES: Only %d CPU cores available with %.2f%% usage - consider scaling up",
			data.System.CPUCores, data.System.CPUUsage))
	}
	
	// Low available memory
	freeMemoryPercent := float64(data.System.MemoryFree) / float64(data.System.MemoryTotal) * 100
	if freeMemoryPercent < 10 {
		warnings = append(warnings, fmt.Sprintf(
			"LOW FREE MEMORY: Only %.1f%% memory free (%s) - system may become unstable",
			freeMemoryPercent, formatBytes(data.System.MemoryFree)))
	}
	
	return warnings
}

func (a *AIAnalyzer) identifyServiceType(data *models.InspectionData) string {
	name := strings.ToLower(data.Process.Name)
	cmdline := strings.ToLower(data.Process.CommandLine)
	
	// Common service patterns
	if strings.Contains(name, "nginx") || strings.Contains(cmdline, "nginx") {
		return "Web Server (Nginx)"
	}
	if strings.Contains(name, "apache") || strings.Contains(cmdline, "apache") {
		return "Web Server (Apache)"
	}
	if strings.Contains(name, "mysql") || strings.Contains(cmdline, "mysql") {
		return "Database (MySQL)"
	}
	if strings.Contains(name, "postgres") || strings.Contains(cmdline, "postgres") {
		return "Database (PostgreSQL)"
	}
	if strings.Contains(name, "redis") || strings.Contains(cmdline, "redis") {
		return "Cache/Database (Redis)"
	}
	if strings.Contains(name, "docker") || strings.Contains(cmdline, "docker") {
		return "Container Runtime (Docker)"
	}
	if strings.Contains(name, "node") || strings.Contains(cmdline, "node") {
		return "Node.js Application"
	}
	if strings.Contains(name, "python") || strings.Contains(cmdline, "python") {
		return "Python Application"
	}
	if strings.Contains(name, "java") || strings.Contains(cmdline, "java") {
		return "Java Application"
	}
	
	return "Unknown Service"
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