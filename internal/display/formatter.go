package display

import (
	"fmt"
	"strings"

	"inspektor/internal/models"

	"github.com/charmbracelet/lipgloss"
)

type Formatter struct{}

func NewFormatter() *Formatter {
	return &Formatter{}
}

func (f *Formatter) FormatReport(data *models.InspectionData) string {
	var output strings.Builder

	// Title with process name
	title := fmt.Sprintf("INSPEKTOR - Process %d (%s)", data.Process.PID, data.Process.Name)
	output.WriteString(titleStyle.Render(title))
	output.WriteString("\n")
	output.WriteString(separatorStyle.Render(strings.Repeat("─", 60)))
	output.WriteString("\n")

	// Process Overview - most important info first
	output.WriteString(f.formatProcessOverview(data.Process))

	// Resource Usage - key metrics
	output.WriteString(f.formatResourceMetrics(data.Process))

	// System Context
	output.WriteString(f.formatSystemContext(data.System))

	return output.String()
}

func (f *Formatter) formatProcessOverview(proc *models.ProcessInfo) string {
	var content strings.Builder

	content.WriteString(sectionStyle.Render(" PROCESS "))
	content.WriteString("\n")

	// Most important info in a clean table format
	items := []struct {
		key   string
		value string
	}{
		{"Status", f.formatStatus(proc.Status)},
		{"Command", proc.CommandLine},
		{"Executable", proc.Executable},
		{"Working Dir", proc.WorkingDir},
		{"Started", proc.CreateTime.Format("Jan 02, 15:04:05")},
	}

	for _, item := range items {
		if item.value != "" {
			content.WriteString(contentStyle.Render(
				keyStyle.Render(item.key+":") + " " + valueStyle.Render(item.value)))
			content.WriteString("\n")
		}
	}

	return content.String()
}

func (f *Formatter) formatResourceMetrics(proc *models.ProcessInfo) string {
	var content strings.Builder

	content.WriteString(sectionStyle.Render(" RESOURCES "))
	content.WriteString("\n")

	// Key metrics with visual indicators
	items := []struct {
		key   string
		value string
	}{
		{"CPU Usage", f.formatCPUUsage(proc.CPUPercent)},
		{"Memory", f.formatMemoryUsage(proc.MemoryRSS, proc.MemoryPercent)},
		{"Virtual Memory", formatBytes(proc.MemoryVMS)},
		{"Open Files", f.formatCount(proc.OpenFiles, 100)},
		{"Connections", f.formatCount(proc.Connections, 50)},
		{"Child Processes", f.formatCount(proc.Children, 10)},
	}

	for _, item := range items {
		content.WriteString(contentStyle.Render(
			keyStyle.Render(item.key+":") + " " + item.value))
		content.WriteString("\n")
	}

	return content.String()
}

func (f *Formatter) formatSystemContext(sys *models.SystemInfo) string {
	var content strings.Builder

	content.WriteString(sectionStyle.Render(" SYSTEM "))
	content.WriteString("\n")

	items := []struct {
		key   string
		value string
	}{
		{"CPU", fmt.Sprintf("%d cores, %s", sys.CPUCores, f.formatCPUUsage(sys.CPUUsage))},
		{"Memory", f.formatSystemMemory(sys.MemoryUsed, sys.MemoryTotal, sys.MemoryPercent)},
		{"CPU Model", f.truncateString(sys.CPUModel, 50)},
	}

	for _, item := range items {
		content.WriteString(contentStyle.Render(
			keyStyle.Render(item.key+":") + " " + item.value))
		content.WriteString("\n")
	}

	return content.String()
}

func (f *Formatter) FormatWarnings(warnings []string) string {
	if len(warnings) == 0 {
		return successMessageStyle.Render("✓ All systems healthy") + "\n\n"
	}

	var output strings.Builder

	// Separate warnings and recommendations
	var actualWarnings []string
	var recommendations []string

	for _, item := range warnings {
		if strings.HasPrefix(item, "⚠") {
			actualWarnings = append(actualWarnings, item)
		} else if strings.HasPrefix(item, "→") {
			recommendations = append(recommendations, item)
		} else {
			// Fallback for items without prefix
			actualWarnings = append(actualWarnings, item)
		}
	}

	// Display warnings first
	if len(actualWarnings) > 0 {
		output.WriteString(warningHeaderStyle.Render(" WARNINGS "))
		output.WriteString("\n")

		for i, warning := range actualWarnings {
			prefix := fmt.Sprintf("  %d. ", i+1)
			output.WriteString(warningItemStyle.Render(prefix + warning))
			output.WriteString("\n")
		}
		output.WriteString("\n")
	}

	// Display recommendations
	if len(recommendations) > 0 {
		recommendHeaderStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#3B82F6")).
			Background(lipgloss.Color("#1E3A8A")).
			Padding(0, 2).
			MarginTop(1).
			MarginBottom(1)

		recommendItemStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#60A5FA")).
			PaddingLeft(2)

		output.WriteString(recommendHeaderStyle.Render(" RECOMMENDATIONS "))
		output.WriteString("\n")

		for i, rec := range recommendations {
			prefix := fmt.Sprintf("  %d. ", i+1)
			output.WriteString(recommendItemStyle.Render(prefix + rec))
			output.WriteString("\n")
		}
		output.WriteString("\n")
	}

	return output.String()
}

// Helper functions for better formatting
func (f *Formatter) formatStatus(status string) string {
	switch strings.ToLower(status) {
	case "r", "running":
		return statusGoodStyle.Render("Running")
	case "s", "sleeping":
		return valueStyle.Render("Sleeping")
	case "z", "zombie":
		return statusWarningStyle.Render("Zombie")
	case "t", "stopped":
		return statusWarningStyle.Render("Stopped")
	default:
		return valueStyle.Render(status)
	}
}

func (f *Formatter) formatCPUUsage(percent float64) string {
	usage := fmt.Sprintf("%.1f%%", percent)
	if percent > 80 {
		return statusWarningStyle.Render(usage)
	} else if percent > 50 {
		return metricStyle.Render(usage)
	}
	return valueStyle.Render(usage)
}

func (f *Formatter) formatMemoryUsage(rss uint64, percent float32) string {
	memory := fmt.Sprintf("%s (%.1f%%)", formatBytes(rss), percent)
	if percent > 10 {
		return statusWarningStyle.Render(memory)
	} else if percent > 5 {
		return metricStyle.Render(memory)
	}
	return valueStyle.Render(memory)
}

func (f *Formatter) formatSystemMemory(used, total uint64, percent float64) string {
	memory := fmt.Sprintf("%s / %s (%.1f%%)", formatBytes(used), formatBytes(total), percent)
	if percent > 85 {
		return statusWarningStyle.Render(memory)
	} else if percent > 70 {
		return metricStyle.Render(memory)
	}
	return valueStyle.Render(memory)
}

func (f *Formatter) formatCount(count, threshold int) string {
	countStr := fmt.Sprintf("%d", count)
	if count > threshold {
		return statusWarningStyle.Render(countStr)
	} else if count > threshold/2 {
		return metricStyle.Render(countStr)
	}
	return valueStyle.Render(countStr)
}

func (f *Formatter) truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return valueStyle.Render(s)
	}
	return valueStyle.Render(s[:maxLen-3] + "...")
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
