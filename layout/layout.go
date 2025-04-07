package layout

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Layout represents our base TUI layout with header, viewport, and footer
type Layout struct {
	header         string
	footer         string
	headerHeight   int
	footerHeight   int
	viewport       viewport.Model
	windowWidth    int
	windowHeight   int
	headerStyle    lipgloss.Style
	footerStyle    lipgloss.Style
	pendingContent string // Add this field
}

// New creates a new Layout with default styles
func New() *Layout {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#0000FF")).
		Width(100).
		Align(lipgloss.Center)

	footerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#333333")).
		Width(100).
		Align(lipgloss.Center)

	return &Layout{
		header:       "Header",
		footer:       "Footer",
		headerHeight: 1,
		footerHeight: 1,
		headerStyle:  headerStyle,
		footerStyle:  footerStyle,
	}
}

// SetHeader sets the header content
func (l *Layout) SetHeader(header string) {
	l.header = header
}

// SetHeaderHeight sets the header height
func (l *Layout) SetHeaderHeight(height int) {
	l.headerHeight = height
}

// SetHeaderStyle sets the header style
func (l *Layout) SetHeaderStyle(style lipgloss.Style) {
	l.headerStyle = style
}

// SetFooter sets the footer content
func (l *Layout) SetFooter(footer string) {
	l.footer = footer
}

// SetFooterHeight sets the footer height
func (l *Layout) SetFooterHeight(height int) {
	l.footerHeight = height
}

// SetFooterStyle sets the footer style
func (l *Layout) SetFooterStyle(style lipgloss.Style) {
	l.footerStyle = style
}

// SetContent sets the viewport content
func (l *Layout) SetContent(content string) {
	// Store the content temporarily if viewport isn't initialized yet
	if l.viewport.Width == 0 {
		l.pendingContent = content
		return
	}
	l.viewport.SetContent(content)
}

// Init initializes the layout model
func (l *Layout) Init() tea.Cmd {
	return nil
}

// Update handles events and updates the layout state
func (l *Layout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		l.windowWidth = msg.Width
		l.windowHeight = msg.Height

		// Adjust styles based on window width
		l.headerStyle = l.headerStyle.Width(l.windowWidth)
		l.footerStyle = l.footerStyle.Width(l.windowWidth)

		// Calculate viewport height (window height - header - footer)
		viewportHeight := l.windowHeight - l.headerHeight - l.footerHeight

		// Initialize or resize the viewport
		if l.viewport.Width == 0 {
			l.viewport = viewport.New(l.windowWidth, viewportHeight)
			// Apply any pending content if it exists
			if l.pendingContent != "" {
				l.viewport.SetContent(l.pendingContent)
				l.pendingContent = ""
			}
		} else {
			l.viewport.Width = l.windowWidth
			l.viewport.Height = viewportHeight
		}

	// Handle keyboard inputs
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return l, tea.Quit
		}

		// Delegate other keyboard events to the viewport
		l.viewport, cmd = l.viewport.Update(msg)
		return l, cmd
	}

	// Handle viewport events (like scrolling)
	l.viewport, cmd = l.viewport.Update(msg)
	return l, cmd
}

// View renders the layout
func (l *Layout) View() string {
	if l.windowHeight == 0 {
		return "Initializing..."
	}

	// Render the header
	header := l.headerStyle.Render(l.header)

	// Render the viewport
	viewportContent := l.viewport.View()

	// Render the footer
	footer := l.footerStyle.Render(l.footer)

	// Combine all elements
	return lipgloss.JoinVertical(lipgloss.Left, header, viewportContent, footer)
}
