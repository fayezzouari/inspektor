package display

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Color palette - more subtle and professional
	primaryColor   = lipgloss.Color("#0EA5E9")  // Sky blue
	secondaryColor = lipgloss.Color("#8B5CF6")  // Purple
	accentColor    = lipgloss.Color("#F59E0B")  // Amber
	warningColor   = lipgloss.Color("#EF4444")  // Red
	successColor   = lipgloss.Color("#22C55E")  // Green
	mutedColor     = lipgloss.Color("#64748B")  // Slate
	textColor      = lipgloss.Color("#F8FAFC")  // Light text
	
	// Main title
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(primaryColor).
		Align(lipgloss.Center).
		MarginBottom(1).
		PaddingTop(1)
	
	// Section headers - cleaner look
	sectionStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(secondaryColor).
		Background(lipgloss.Color("#1E293B")).
		Padding(0, 2).
		MarginTop(1).
		MarginBottom(1)
	
	// Key-value pairs
	keyStyle = lipgloss.NewStyle().
		Foreground(mutedColor).
		Width(20).
		Align(lipgloss.Right)
	
	valueStyle = lipgloss.NewStyle().
		Foreground(textColor).
		Bold(false)
	
	// Important values (metrics)
	metricStyle = lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true)
	
	// Status indicators
	statusGoodStyle = lipgloss.NewStyle().
		Foreground(successColor).
		Bold(true)
	
	statusWarningStyle = lipgloss.NewStyle().
		Foreground(warningColor).
		Bold(true)
	
	// Warning section
	warningHeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(warningColor).
		Background(lipgloss.Color("#7F1D1D")).
		Padding(0, 2).
		MarginTop(2).
		MarginBottom(1)
	
	warningItemStyle = lipgloss.NewStyle().
		Foreground(warningColor).
		PaddingLeft(2)
	
	// Success message
	successMessageStyle = lipgloss.NewStyle().
		Foreground(successColor).
		Bold(true).
		MarginTop(2).
		Align(lipgloss.Center)
	
	// Container styles
	contentStyle = lipgloss.NewStyle().
		PaddingLeft(2).
		MarginBottom(1)
	
	separatorStyle = lipgloss.NewStyle().
		Foreground(mutedColor).
		MarginTop(1).
		MarginBottom(1)
)