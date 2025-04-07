# Scaffold

A lightweight, reusable layout library for terminal user interfaces in Go, built with Bubble Tea.

## Why We Built This

Building terminal user interfaces often involves recreating the same basic layout patterns. We created this scaffold to eliminate boilerplate code and provide a clean foundation for TUI applications. Now you can focus on your application logic rather than implementing header, footer, and scrolling viewport components from scratch.

## Features

- **Three-Panel Layout**: Fixed header, scrollable content viewport, and fixed footer
- **Flexible Styling**: Comprehensive styling for headers and footers using lipgloss
- **Viewport Navigation**: Built-in scrolling with multiple navigation methods
- **Content Management**: Simple API for setting and manipulating content
- **Window Resizing**: Automatic handling of terminal window resize events
- **Minimal Dependencies**: Built on proven libraries (Bubble Tea, Bubbles, Lipgloss)

## Installation

```bash
go get github.com/makeitmini/scaffold
```

## Usage Example

Here's a simple example demonstrating how to use the layout:

```go
package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/makeitmini/scaffold/layout"
)

func main() {
	// Initialize our layout
	l := layout.New()

	// Customize the header
	l.SetHeader("My Application")
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#0000FF")).
		Padding(1, 0).
		Width(100).
		Align(lipgloss.Center)
	l.SetHeaderStyle(headerStyle)
	l.SetHeaderHeight(3)

	// Customize the footer
	l.SetFooter("Press 'q' to quit | Use arrow keys to navigate")
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#333333")).
		Padding(0, 1).
		Width(100).
		Align(lipgloss.Center)
	l.SetFooterStyle(footerStyle)

	// Generate some content for the viewport
	var content strings.Builder
	for i := 1; i <= 100; i++ {
		fmt.Fprintf(&content, "Line %d: This is some sample content.\n", i)
	}

	// Set the content
	l.SetContent(content.String())

	// Create a custom application model that embeds our layout
	app := myApp{layout: l}

	// Start the application with our custom update function
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

// A simple application model that embeds our layout
type myApp struct {
	layout *layout.Layout
}

func (m myApp) Init() tea.Cmd {
	return m.layout.Init()
}

func (m myApp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Add custom key handlers
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "home":
			m.layout.ScrollToTop()
			return m, nil
		case "end":
			m.layout.ScrollToBottom()
			return m, nil
		}
	}

	// Delegate other messages to the layout
	var cmd tea.Cmd
	model, cmd := m.layout.Update(msg)
	// Handle type assertion
	if layoutModel, ok := model.(*layout.Layout); ok {
		m.layout = layoutModel
	}
	return m, cmd
}

func (m myApp) View() string {
	return m.layout.View()
}
```

## Navigation Methods

The layout provides several methods for navigating the viewport content:

```go
// Single line scrolling
layout.LineDown()
layout.LineUp()

// Page scrolling
layout.PageDown()
layout.PageUp()

// Half page scrolling
layout.HalfPageDown()
layout.HalfPageUp()

// Jump to extremes
layout.ScrollToTop()
layout.ScrollToBottom()
```

## Customizing and Getting Viewport Information

You can also access viewport information:

```go
// Get viewport dimensions
height := layout.GetViewportHeight()
width := layout.GetViewportWidth()

// Get position information
position := layout.GetViewportYPosition()
atTop := layout.ViewportAtTop()
atBottom := layout.ViewportAtBottom()
```

## License

MIT

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.
