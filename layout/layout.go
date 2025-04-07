package layout

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Layout represents our base TUI layout with header, viewport, and footer
type Layout struct {
	header          string
	footer          string
	headerHeight    int
	footerHeight    int
	viewport        viewport.Model
	windowWidth     int
	windowHeight    int
	headerStyle     lipgloss.Style
	footerStyle     lipgloss.Style
	pendingContent  string
	verticalAlign   VerticalAlignment
	viewportContent string
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
		header:        "Header",
		footer:        "Footer",
		headerHeight:  1,
		footerHeight:  1,
		headerStyle:   headerStyle,
		footerStyle:   footerStyle,
		verticalAlign: AlignTop, // Default to top alignment
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
	// Store the raw content
	l.viewportContent = content

	// Store the content temporarily if viewport isn't initialized yet
	if l.viewport.Width == 0 {
		l.pendingContent = content
		return
	}

	// Apply alignment and set the content
	l.applyAlignment()
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
				l.viewportContent = l.pendingContent
				l.pendingContent = ""
				l.applyAlignment()
			}
		} else {
			l.viewport.Width = l.windowWidth
			l.viewport.Height = viewportHeight
			// Reapply alignment with new dimensions
			l.applyAlignment()
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

func (l *Layout) SetVerticalAlignment(alignment VerticalAlignment) {
	l.verticalAlign = alignment
	// If we already have content, realign it
	if l.viewportContent != "" && l.viewport.Width > 0 {
		l.applyAlignment()
	}
}

type VerticalAlignment int

const (
	// AlignTop aligns content to the top of the viewport (default)
	AlignTop VerticalAlignment = iota
	// AlignCenter centers content vertically in the viewport
	AlignCenter
	// AlignBottom aligns content to the bottom of the viewport
	AlignBottom
)

// applyAlignment applies the vertical alignment to the viewport content
func (l *Layout) applyAlignment() {
	// If viewport isn't initialized yet, don't do anything
	if l.viewport.Width == 0 {
		return
	}

	// Count the number of lines in the content
	lines := strings.Count(l.viewportContent, "\n") + 1

	// If content is shorter than viewport height, add padding based on alignment
	if lines < l.viewport.Height {
		paddingSize := l.viewport.Height - lines

		switch l.verticalAlign {
		case AlignTop:
			// For top alignment, no padding at the top
			l.viewport.SetContent(l.viewportContent)

		case AlignCenter:
			// For center alignment, add equal padding to top and bottom
			topPadding := paddingSize / 2
			var paddedContent strings.Builder

			// Add top padding
			for i := 0; i < topPadding; i++ {
				paddedContent.WriteString("\n")
			}

			// Add content
			paddedContent.WriteString(l.viewportContent)

			l.viewport.SetContent(paddedContent.String())

		case AlignBottom:
			// For bottom alignment, add all padding to the top
			var paddedContent strings.Builder

			// Add top padding
			for i := 0; i < paddingSize; i++ {
				paddedContent.WriteString("\n")
			}

			// Add content
			paddedContent.WriteString(l.viewportContent)

			l.viewport.SetContent(paddedContent.String())
		}
	} else {
		// If content is taller than viewport, no special alignment needed
		l.viewport.SetContent(l.viewportContent)
	}
}
