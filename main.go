package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func main() {
	gtk.Init(nil)

	// Create the main window
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("Text Editor")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	manager := NewManager()
	drawingArea, _ := gtk.DrawingAreaNew()
	drawingArea.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) {
		FONT_SIZE := 20.0
		LINE_HEIGHT := float64(FONT_SIZE + 3)
		cr.SetSourceRGB(0.15625, 0.1640625, 0.2109375)
		cr.Paint()
		cr.SetAntialias(cairo.ANTIALIAS_SUBPIXEL)
		cr.SetFontSize(FONT_SIZE)
		text := manager.GetText(da.GetAllocatedWidth(), da.GetAllocatedHeight())
		lines := strings.Split(text, "\n")
		x, y := manager.GetCursorXY()
		// Highlight current line
		cr.SetSourceRGB(0.396, 0.412, 0.525)
		cr.Rectangle(0.0, (float64((y))*LINE_HEIGHT)+5.0, float64(da.GetAllocatedWidth()), 20.0)
		cr.StrokePreserve()
		cr.Fill()
		// Print the cursor.
		cr.SetSourceRGB(1, 1, 1)
		cr.Rectangle(float64(50.0+float64(x)*11), (float64((y))*LINE_HEIGHT)+5.0, 0.05, 20.0)
		cr.StrokePreserve()
		cr.Fill()
		if len(lines) != 0 {
			for i, line := range lines {
				if y == i {
					cr.SelectFontFace("Mono", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_BOLD)
				} else {
					cr.SelectFontFace("Mono", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_NORMAL)
				}
				cr.MoveTo(5, float64(i+1)*(LINE_HEIGHT))
				cr.SetSourceRGB(0.265625, 0.27734375, 0.3515625)
				cr.ShowText(fmt.Sprint(i + 1))
				j := 0
				for _, char := range line {
					cr.SelectFontFace("Mono", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_NORMAL)
					cr.MoveTo(50+(float64(j)*11), float64(i+1)*(LINE_HEIGHT))
					cr.SetSourceRGB(0.96875, 0.96875, 0.9453125)
					if string(char) == "\t" {
						cr.SetSourceRGB(0.265625, 0.27734375, 0.3515625)
						cr.ShowText("»")
						j++
					} else {
						cr.ShowText(string(char))
						j++
					}
				}
			}
		}
		cr.Fill()
	})
	win.Connect("key-press-event", func(win *gtk.Window, ev *gdk.Event) {
		manager.ReadKey(win, ev)
	})
	win.Connect("key-release-event", func(win *gtk.Window, ev *gdk.Event) {
		manager.CheckModifier(win, ev)
	})
	// Add the scrolled window to the main window
	win.Add(drawingArea)

	win.ShowAll()

	// Start the GTK main loop
	gtk.Main()
}
