package main

import (
	"log"

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

	// Create a text view
	textView, err := gtk.TextViewNew()
	if err != nil {
		log.Fatal("Unable to create text view:", err)
	}

	// Create a scrollable container for the text view
	scrolledWin, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatal("Unable to create scrolled window:", err)
	}
	scrolledWin.Add(textView)

	// Add the scrolled window to the main window
	win.Add(scrolledWin)

	// Show all the widgets
	win.ShowAll()

	// Start the GTK main loop
	gtk.Main()
}
