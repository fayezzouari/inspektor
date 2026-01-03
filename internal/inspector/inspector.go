package inspector

import (
	"encoding/json"
	"fmt"
	"time"

	"inspektor/internal/analyzer"
	"inspektor/internal/display"
	"inspektor/internal/models"

	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
)

type Inspector struct {
	analyzer  *analyzer.AIAnalyzer
	formatter *display.Formatter
}

func New() *Inspector {
	return &Inspector{
		analyzer:  analyzer.New(),
		formatter: display.NewFormatter(),
	}
}

func (i *Inspector) InspectWithOptions(pid int32, jsonOutput, verbose bool) error {
	// Ensure AI client is properly closed
	defer func() {
		if err := i.analyzer.Close(); err != nil {
			fmt.Printf("Warning: Failed to close AI client: %v\n", err)
		}
	}()

	// Show banner and start processing animation (skip for JSON output)
	if !jsonOutput {
		display.ShowBanner("")
		done := make(chan bool)
		go display.ShowProcessingAnimation("Analyzing process and system metrics...", done)
		defer func() {
			done <- true
			close(done)
			time.Sleep(100 * time.Millisecond) // Give time to clear the animation
		}()
	}

	// Get process information
	proc, err := process.NewProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to get process: %w", err)
	}

	// Collect process data
	processInfo, err := i.collectProcessInfo(proc)
	if err != nil {
		return fmt.Errorf("failed to collect process info: %w", err)
	}

	// Collect system data
	systemInfo, err := i.collectSystemInfo()
	if err != nil {
		return fmt.Errorf("failed to collect system info: %w", err)
	}

	// Create inspection data
	data := &models.InspectionData{
		Process: processInfo,
		System:  systemInfo,
	}

	// Generate AI analysis and warnings
	warnings := i.analyzer.AnalyzeAndWarn(data)

	if jsonOutput {
		return i.outputJSON(data, warnings)
	}

	// Display results in rich format
	fmt.Print(i.formatter.FormatReport(data))
	fmt.Print(i.formatter.FormatWarnings(warnings))

	return nil
}

func (i *Inspector) Inspect(pid int32) error {
	return i.InspectWithOptions(pid, false, false)
}

func (i *Inspector) InspectByPort(port int, jsonOutput, verbose bool) error {
	// Show banner for port lookup (skip for JSON output)
	if !jsonOutput {
		display.ShowBanner("")
		done := make(chan bool)
		go display.ShowProcessingAnimation(fmt.Sprintf("Finding process on port %d...", port), done)

		// Find the PID listening on the specified port
		pid, err := i.findProcessByPort(port)

		done <- true
		close(done)
		time.Sleep(100 * time.Millisecond)

		if err != nil {
			return fmt.Errorf("failed to find process on port %d: %w", port, err)
		}

		fmt.Printf("\n%s\n\n",
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#22C55E")).
				Bold(true).
				Render(fmt.Sprintf("✓ Found process %d listening on port %d", pid, port)))
	} else {
		// Silent lookup for JSON mode
		pid, err := i.findProcessByPort(port)
		if err != nil {
			return fmt.Errorf("failed to find process on port %d: %w", port, err)
		}
		return i.InspectWithOptions(pid, jsonOutput, verbose)
	}

	// Continue with normal inspection (which will show its own banner)
	pid, _ := i.findProcessByPort(port)
	return i.InspectWithOptions(pid, jsonOutput, verbose)
}

func (i *Inspector) findProcessByPort(port int) (int32, error) {
	// Get all network connections
	connections, err := net.Connections("all")
	if err != nil {
		return 0, fmt.Errorf("failed to get network connections: %w", err)
	}

	// Find connections matching the port
	var candidatePIDs []int32
	for _, conn := range connections {
		if conn.Laddr.Port == uint32(port) && conn.Status == "LISTEN" {
			candidatePIDs = append(candidatePIDs, conn.Pid)
		}
	}

	if len(candidatePIDs) == 0 {
		return 0, fmt.Errorf("no process found listening on port %d", port)
	}

	// Return the first valid PID
	for _, pid := range candidatePIDs {
		if pid > 0 {
			// Verify the process exists
			if _, err := process.NewProcess(pid); err == nil {
				return pid, nil
			}
		}
	}

	return 0, fmt.Errorf("no valid process found listening on port %d", port)
}

func (i *Inspector) outputJSON(data *models.InspectionData, warnings []string) error {
	output := struct {
		*models.InspectionData
		Warnings []string `json:"warnings"`
	}{
		InspectionData: data,
		Warnings:       warnings,
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	fmt.Println(string(jsonData))
	return nil
}

func (i *Inspector) collectProcessInfo(proc *process.Process) (*models.ProcessInfo, error) {
	name, _ := proc.Name()
	exe, _ := proc.Exe()
	cmdline, _ := proc.Cmdline()
	cwd, _ := proc.Cwd()
	status, _ := proc.Status()

	// CPU and Memory usage
	cpuPercent, _ := proc.CPUPercent()
	memInfo, _ := proc.MemoryInfo()
	memPercent, _ := proc.MemoryPercent()

	// Process times
	createTime, _ := proc.CreateTime()

	// Connections and open files
	connections, _ := proc.Connections()
	openFiles, _ := proc.OpenFiles()

	// Child processes
	children, _ := proc.Children()

	return &models.ProcessInfo{
		PID:           proc.Pid,
		Name:          name,
		Executable:    exe,
		CommandLine:   cmdline,
		WorkingDir:    cwd,
		Status:        status,
		CPUPercent:    cpuPercent,
		MemoryRSS:     memInfo.RSS,
		MemoryVMS:     memInfo.VMS,
		MemoryPercent: memPercent,
		CreateTime:    time.Unix(createTime/1000, 0),
		Connections:   len(connections),
		OpenFiles:     len(openFiles),
		Children:      len(children),
	}, nil
}

func (i *Inspector) collectSystemInfo() (*models.SystemInfo, error) {
	// CPU information
	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}

	// Memory information
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	return &models.SystemInfo{
		CPUCores:      len(cpuInfo),
		CPUModel:      cpuInfo[0].ModelName,
		CPUUsage:      cpuPercent[0],
		MemoryTotal:   memInfo.Total,
		MemoryUsed:    memInfo.Used,
		MemoryPercent: memInfo.UsedPercent,
		MemoryFree:    memInfo.Free,
	}, nil
}
