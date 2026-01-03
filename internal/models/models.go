package models

import "time"

// ProcessInfo contains detailed information about a specific process
type ProcessInfo struct {
	PID           int32     `json:"pid"`
	Name          string    `json:"name"`
	Executable    string    `json:"executable"`
	CommandLine   string    `json:"command_line"`
	WorkingDir    string    `json:"working_dir"`
	Status        string    `json:"status"`
	CPUPercent    float64   `json:"cpu_percent"`
	MemoryRSS     uint64    `json:"memory_rss"`
	MemoryVMS     uint64    `json:"memory_vms"`
	MemoryPercent float32   `json:"memory_percent"`
	CreateTime    time.Time `json:"create_time"`
	Connections   int       `json:"connections"`
	OpenFiles     int       `json:"open_files"`
	Children      int       `json:"children"`
}

// SystemInfo contains system-wide resource information
type SystemInfo struct {
	CPUCores      int     `json:"cpu_cores"`
	CPUModel      string  `json:"cpu_model"`
	CPUUsage      float64 `json:"cpu_usage"`
	MemoryTotal   uint64  `json:"memory_total"`
	MemoryUsed    uint64  `json:"memory_used"`
	MemoryPercent float64 `json:"memory_percent"`
	MemoryFree    uint64  `json:"memory_free"`
}

// InspectionData combines process and system information
type InspectionData struct {
	Process *ProcessInfo `json:"process"`
	System  *SystemInfo  `json:"system"`
}