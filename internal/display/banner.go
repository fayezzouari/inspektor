package display

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

const banner = ` _____  _   _  ___________ _____ _   _______ ___________  
|_   _|| \ | |/  ___| ___ \  ___| | / /_   _|  _  | ___ \ 
  | |  |  \| |\ ` + "`" + `--.| |_/ / |__ | |/ /  | | | | | | |_/ / 
  | |  | . ` + "`" + ` | ` + "`" + `--. \  __/|  __||    \  | | | | | |    /  
 _| |_ | |\  |/\__/ / |   | |___| |\  \ | | \ \_/ / |\ \  
 \___/ \_| \_/\____/\_|   \____/\_| \_/ \_/  \___/\_| \_| `

var bannerStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#0EA5E9")).
	Bold(true).
	Align(lipgloss.Center)

var processingStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#8B5CF6")).
	Italic(true).
	Align(lipgloss.Center).
	MarginTop(1)

// ShowBanner displays the INSPEKTOR banner with a processing message
func ShowBanner(message string) {
	fmt.Println()
	fmt.Println(bannerStyle.Render(banner))
	fmt.Println()
	if message != "" {
		fmt.Println(processingStyle.Render(message))
		fmt.Println()
	}
}

// ShowProcessingAnimation displays an animated processing message
func ShowProcessingAnimation(message string, done chan bool) {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0

	ticker := time.NewTicker(80 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			// Clear the line
			fmt.Print("\r\033[K")
			return
		case <-ticker.C:
			frame := frames[i%len(frames)]
			fmt.Printf("\r%s %s",
				lipgloss.NewStyle().Foreground(lipgloss.Color("#8B5CF6")).Render(frame),
				lipgloss.NewStyle().Foreground(lipgloss.Color("#64748B")).Render(message))
			i++
		}
	}
}
