package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/makeitmini/scaffold/layout"
)

const version = "v1.0.0"

// testApp is a wrapper around our layout that adds custom key handling
type testApp struct {
	layout      *layout.Layout
	showHelp    bool
	showPosInfo bool
}

func (t testApp) Init() tea.Cmd {
	return t.layout.Init()
}

// updateFooter updates the footer based on current state
func (t *testApp) updateFooter() {
	if t.showPosInfo {
		// Show simplified position info without indicators that need real-time updates
		t.layout.SetFooter("Press 'p' to hide position info | 'q' to quit")
	} else if t.showHelp {
		// Show help info in the footer
		t.layout.SetFooter("h: toggle help | p: show position | q: quit | ↑/↓, j/k: line | u/d: half page | page up/down | home/end, t/b: top/bottom")
	} else {
		// Default footer
		t.layout.SetFooter("Press 'h' for help | 'p' for position info | 'q' to quit")
	}
}

func (t testApp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return t, tea.Quit
		case "j", "down":
			t.layout.LineDown()
			return t, nil
		case "k", "up":
			t.layout.LineUp()
			return t, nil
		case "u":
			t.layout.HalfPageUp()
			return t, nil
		case "d":
			t.layout.HalfPageDown()
			return t, nil
		case "page_up":
			t.layout.PageUp()
			return t, nil
		case "page_down":
			t.layout.PageDown()
			return t, nil
		case "home", "t":
			t.layout.ScrollToTop()
			return t, nil
		case "end", "b":
			t.layout.ScrollToBottom()
			return t, nil
		case "h":
			// Toggle help info
			t.showHelp = !t.showHelp
			if t.showHelp {
				t.showPosInfo = false // Turn off position info if help is on
			}
			t.updateFooter()
			return t, nil
		case "p":
			// Toggle position info
			t.showPosInfo = !t.showPosInfo
			if t.showPosInfo {
				t.showHelp = false // Turn off help if position info is on
				t.updateFooter()   // Update the position immediately when toggled on
			} else {
				t.updateFooter() // Show the default footer when toggled off
			}
			return t, nil
		}
	}

	// Pass other messages to the layout and handle type assertion
	model, cmd := t.layout.Update(msg)
	if layoutModel, ok := model.(*layout.Layout); ok {
		t.layout = layoutModel
	}

	return t, cmd
}

func (t testApp) View() string {
	return t.layout.View()
}

// generateTestContent creates a rich content for testing scrolling
func generateTestContent() string {
	var content strings.Builder

	// Help section at the top
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFF00")).
		Background(lipgloss.Color("#333333")).
		Bold(true).
		Padding(1)

	help := `Navigation Controls:
↑/↓, j/k      : Scroll up/down one line
PageUp/PageDown: Scroll full page up/down
Home/End       : Go to top/bottom
u/d           : Scroll half page up/down
t/b           : Go to top/bottom
h             : Toggle help in footer
p             : Show position info (snapshot)
q             : Quit`

	content.WriteString(helpStyle.Render(help))
	content.WriteString("\n\n")

	// Section markers to help identify position
	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#8A2BE2")).
		Bold(true).
		Width(50).
		Align(lipgloss.Center).
		Padding(1, 0)

	// Add top marker
	content.WriteString(sectionStyle.Render("TOP OF CONTENT"))
	content.WriteString("\n\n")

	// Add numbered sections with distinctive styling
	sectionCount := 10
	linesPerSection := 30

	for section := 1; section <= sectionCount; section++ {
		// Section header
		sectionHeader := fmt.Sprintf("Section %d of %d", section, sectionCount)
		content.WriteString(sectionStyle.Render(sectionHeader))
		content.WriteString("\n\n")

		// Section content with alternating styles
		for line := 1; line <= linesPerSection; line++ {
			// Create line markers every 5 lines
			if line%5 == 0 {
				markerStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("#FF0000")).
					Bold(true)

				content.WriteString(markerStyle.Render(
					fmt.Sprintf("---- MARKER: Section %d, Line %d ----", section, line)))
				content.WriteString("\n")
			}

			// Regular line with section and line number
			lineContent := fmt.Sprintf("Section %d, Line %d: This is sample content for testing scrolling\n",
				section, line)

			// Alternate line colors for better visual distinction
			if line%2 == 0 {
				lineStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA"))
				content.WriteString(lineStyle.Render(lineContent))
			} else {
				content.WriteString(lineContent)
			}
		}

		content.WriteString("\n")
	}

	// Add bottom marker
	content.WriteString(sectionStyle.Render("BOTTOM OF CONTENT"))
	content.WriteString("\n")

	return content.String()
}

func main() {
	l := layout.New()

	// Customize the header with version
	l.SetHeader(fmt.Sprintf("TUI Layout Scaffold Demo %s", version))
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#8A2BE2")). // BlueViolet
		Padding(1, 1).
		Width(100).
		Align(lipgloss.Center)
	l.SetHeaderStyle(headerStyle)
	l.SetHeaderHeight(3) // Give header more space

	// Customize the footer
	footerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#333333")).
		Padding(0, 1).
		Width(100).
		Align(lipgloss.Center)
	l.SetFooterStyle(footerStyle)

	// Set rich content for testing
	l.SetContent(generateTestContent())

	// Create our test app wrapper
	app := testApp{
		layout:      l,
		showHelp:    false,
		showPosInfo: false,
	}

	// Set initial footer
	app.updateFooter()

	// Start the application
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
