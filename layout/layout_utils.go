package layout

// GetViewportHeight returns the current height of the viewport
func (l *Layout) GetViewportHeight() int {
	return l.viewport.Height
}

// GetViewportWidth returns the current width of the viewport
func (l *Layout) GetViewportWidth() int {
	return l.viewport.Width
}

// SetViewportYPosition sets the vertical scroll position of the viewport
func (l *Layout) SetViewportYPosition(y int) {
	l.viewport.YPosition = y
}

// GetViewportYPosition gets the current vertical scroll position
func (l *Layout) GetViewportYPosition() int {
	return l.viewport.YPosition
}

// ViewportAtTop returns whether the viewport is scrolled to the top
func (l *Layout) ViewportAtTop() bool {
	return l.viewport.AtTop()
}

// ViewportAtBottom returns whether the viewport is scrolled to the bottom
func (l *Layout) ViewportAtBottom() bool {
	return l.viewport.AtBottom()
}

// LineDown scrolls the viewport down one line
func (l *Layout) LineDown() {
	l.viewport.LineDown(1)
}

// LineUp scrolls the viewport up one line
func (l *Layout) LineUp() {
	l.viewport.LineUp(1)
}

// HalfPageDown scrolls the viewport down half a page
func (l *Layout) HalfPageDown() {
	// Calculate half page height
	halfPageHeight := l.viewport.Height / 2
	l.viewport.LineDown(halfPageHeight)
}

// HalfPageUp scrolls the viewport up half a page
func (l *Layout) HalfPageUp() {
	// Calculate half page height
	halfPageHeight := l.viewport.Height / 2
	l.viewport.LineUp(halfPageHeight)
}

// PageDown scrolls the viewport down a full page
func (l *Layout) PageDown() {
	l.viewport.LineDown(l.viewport.Height)
}

// PageUp scrolls the viewport up a full page
func (l *Layout) PageUp() {
	l.viewport.LineUp(l.viewport.Height)
}

// ScrollToTop scrolls to the top of the viewport
func (l *Layout) ScrollToTop() {
	l.viewport.GotoTop()
}

// ScrollToBottom scrolls to the bottom of the viewport
func (l *Layout) ScrollToBottom() {
	l.viewport.GotoBottom()
}
