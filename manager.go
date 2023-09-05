package main

import (
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type modifiers struct {
	alt   bool
	ctrl  bool
	shift bool
	super bool
}

type Manager struct {
	text      string
	x, y      int
	modifiers modifiers
}

func NewManager() *Manager {
	manager := &Manager{
		text: "",
		x:    0,
		modifiers: modifiers{
			alt:   false,
			ctrl:  false,
			shift: false,
			super: false,
		},
		y: 0,
	}
	return manager
}
func (m *Manager) modifiersToString(key string) string {
	if m.modifiers.alt {
		return "Alt+" + key
	}
	if m.modifiers.ctrl {
		return "Ctrl+" + key
	}
	if m.modifiers.shift {
		return "Shift+" + key
	}
	if m.modifiers.super {
		return "Super+" + key
	}
	return key
}
func (m *Manager) checkShortcut(win *gtk.Window, key string) {
	// TODO: Fix the shortcuts system
	// switch m.modifiersToString(key) {
	// case "Ctrl+x":
	// 	if m.y != 0 {
	// 		m.delLine()
	// 		m.y--
	// 	}

	// }
}
func (m *Manager) CheckModifier(win *gtk.Window, ev *gdk.Event) {
	keyEvent := &gdk.EventKey{
		Event: ev,
	}
	key := keyEvent.KeyVal()
	switch {
	case key == gdk.KEY_Alt_L, key == gdk.KEY_Alt_R:
		m.modifiers.alt = false
	case key == gdk.KEY_Control_L, key == gdk.KEY_Control_R:
		m.modifiers.ctrl = false
	case key == gdk.KEY_Super_L, key == gdk.KEY_Super_R:
		m.modifiers.super = false
	case key == gdk.KEY_Shift_L, key == gdk.KEY_Shift_R:
		m.modifiers.shift = false
	}
}

func (m *Manager) delLine() {
	m.text = strings.Join(strings.Split(m.text, "\n")[:m.y], "\n")
	m.x = 0
}
func (m *Manager) ReadKey(win *gtk.Window, ev *gdk.Event) {
	keyEvent := &gdk.EventKey{
		Event: ev,
	}
	key := keyEvent.KeyVal()
	switch {
	// Add support to AZERTY keyboard:
	case key >= 32 && key <= 127:
		if m.modifiers.alt || m.modifiers.ctrl || m.modifiers.super {
			m.checkShortcut(win, string(rune(key)))
		}
		m.addChar(string(rune(key)))
	case key == gdk.KEY_BackSpace:
		if m.x == 0 {
			lines := strings.Split(m.text, "\n")
			newx := 0
			if m.y > 0 {
				newx = len(lines[m.y-1])
			}
			m.bkspChar()
			m.x = newx
		} else {
			m.bkspChar()
		}
	case key == gdk.KEY_Return:
		m.addChar("\n")
		m.x = 0
		m.y++
	case key == gdk.KEY_Tab:
		for i := 0; i < 4; i++ {
			m.addChar(" ")
		}
	case key == gdk.KEY_Up:
		m.cursorUp()
	case key == gdk.KEY_Down:
		m.cursorDown()
	case key == gdk.KEY_Right:
		m.cursorRight()
	case key == gdk.KEY_Left:
		m.cursorLeft()
	case key == gdk.KEY_Alt_L, key == gdk.KEY_Alt_R:
		m.modifiers.alt = true
	case key == gdk.KEY_Control_L, key == gdk.KEY_Control_R:
		m.modifiers.ctrl = true
	case key == gdk.KEY_Super_L, key == gdk.KEY_Super_R:
		m.modifiers.super = true
	case key == gdk.KEY_Shift_L, key == gdk.KEY_Shift_R:
		m.modifiers.shift = true
	default:
		println(m.modifiersToString(string(key)))
	}
	win.QueueDraw()
}

func (m *Manager) cursorRight() {
	lines := strings.Split(m.text, "\n")
	if len(lines) < m.y+1 {
		return
	}
	if len(lines[m.y]) == m.x {
		if len(lines) != m.y+1 {
			m.y++
			m.x = 0
		}
	} else {
		m.x++
	}
}

func (m *Manager) cursorUp() {
	lines := strings.Split(m.text, "\n")
	if m.y > 0 {
		m.y--
		if len(lines[m.y]) < m.x {
			m.x = len(lines[m.y])
		}
	}
}

func (m *Manager) cursorDown() {
	lines := strings.Split(m.text, "\n")
	if m.y+1 < len(lines) {
		m.y++
		if len(lines[m.y]) < m.x {
			m.x = len(lines[m.y])
		}
	}
}
func (m *Manager) addChar(c string) {
	offset := m.getOffset()
	m.text = m.text[:offset] + c + m.text[offset:]
	m.x += len(c)
}
func (m *Manager) bkspChar() {
	offset := m.getOffset()
	if offset > 0 {
		m.text = m.text[:offset-1] + m.text[offset:]
		m.cursorLeft()
	}
}
func (m *Manager) cursorLeft() {
	if m.x != 0 {
		m.x--
		return
	}
	if m.y == 0 {
		m.x = 0
		return
	}
	m.y--
	lines := strings.Split(m.text, "\n")
	if len(lines) >= m.y+1 {
		m.x = len(lines[m.y])
	}
}

func (m Manager) getOffset() int {
	return m.getOffsetFor(m.x, m.y)
}

func (m Manager) getOffsetFor(x, y int) int {
	lines := strings.Split(m.text, "\n")
	offset := 0
	for i, val := range lines {
		if y > i {
			offset += len(val) + 1
		}
		if y == i {
			offset += x
			break
		}
	}
	return offset
}

func (m *Manager) GetText(w, h int) string {
	return m.getNLines(int(h/20.0) + 1)
}

func (m *Manager) getNLines(n int) string {
	lines := strings.Split(m.text, "\n")
	text := ""
	for i, line := range lines {
		text += line + "\n"
		if i+1 == n {
			break
		}
	}
	text = text[:len(text)-1]
	return text
}

func (m Manager) GetCursorXY() (int, int) {
	return m.x, m.y
}
