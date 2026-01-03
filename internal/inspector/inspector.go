package inspector

import (
	"encoding/json"
	"fmt"
	"time"

	"inspektor/internal/analyzer"
	"inspektor/internal/display"
	"inspektor/internal/models"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
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
