package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/enricod/goraw/libraw"
	"github.com/gotk3/gotk3/gtk"
)

type Settings struct {
	ImagesDir string
}

var appSettings Settings

func extensions() []string {
	return []string{".ORF", ".CR2", ".RAF", ".ARW"}
}

// IsStringInSlice true if the slice contains the string a
func IsStringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func main() {

	appSettings = Settings{ImagesDir: "."}
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	win.Add(mainPanel())
	win.ShowAll()

	gtk.Main()

	//fmt.Scanln()
	fmt.Println("done")
}

func DoExtract(appSettings Settings) {
	exportPath := appSettings.ImagesDir + "/_export"
	if _, err := os.Stat(exportPath); os.IsNotExist(err) {
		os.Mkdir(exportPath, 0777)
	}

	db, err := bolt.Open(exportPath+"/goraw.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	files, err := ioutil.ReadDir(appSettings.ImagesDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		//fmt.Printf("%s", f.Name())
		if IsStringInSlice(filepath.Ext(f.Name()), extensions()) {
			libraw.ExportThumb(appSettings.ImagesDir, f, exportPath)
		}
	}
}

func mainPanel() *gtk.Widget {
	grid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create grid:", err)
	}
	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	grid.SetRowSpacing(6)
	grid.SetMarginStart(6)
	grid.SetMarginTop(6)

	entry, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create entry:", err)
	}
	s, _ := entry.GetText()
	label, err := gtk.LabelNew(s)
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}
	grid.Add(entry)
	entry.SetHExpand(true)
	grid.AttachNextTo(label, entry, gtk.POS_RIGHT, 1, 1)
	label.SetHExpand(true)

	// Connects this entry's "activate" signal (which is emitted whenever
	// Enter is pressed when the Entry is activated) to an anonymous
	// function that gets the current text of the entry and sets the text of
	// the label beside it with it.  Unlike with native GTK callbacks,
	// (*glib.Object).Connect() supports closures.  In this example, this is
	// demonstrated by using the label variable.  Without closures, a
	// pointer to the label would need to be passed in as user data
	// (demonstrated in the next example).
	entry.Connect("activate", func() {
		s, _ := entry.GetText()
		label.SetText(s)
	})

	sb, err := gtk.SpinButtonNewWithRange(0, 1, 0.1)
	if err != nil {
		log.Fatal("Unable to create spin button:", err)
	}
	pb, err := gtk.ProgressBarNew()
	if err != nil {
		log.Fatal("Unable to create progress bar:", err)
	}
	grid.Add(sb)
	sb.SetHExpand(true)
	grid.AttachNextTo(pb, sb, gtk.POS_RIGHT, 1, 1)
	label.SetHExpand(true)

	// Pass in a ProgressBar and the target SpinButton as user data rather
	// than using the sb and pb variables scoped to the anonymous func.
	// This can be useful when passing in a closure that has already been
	// generated, but when you still wish to connect the callback with some
	// variables only visible in this scope.
	sb.Connect("value-changed", func(sb *gtk.SpinButton, pb *gtk.ProgressBar) {
		pb.SetFraction(sb.GetValue() / 1)
	}, pb)

	label, err = gtk.LabelNew("")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}
	s = "Hyperlink to <a href=\"https://www.cyphertite.com/\">Cyphertite</a> for your clicking pleasure"
	label.SetMarkup(s)
	grid.AttachNextTo(label, sb, gtk.POS_BOTTOM, 2, 1)

	dirChooserBtn, err := gtk.FileChooserButtonNew("Dir selection", gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER)
	if err != nil {
		log.Fatal("Unable to create FileChooserDialogNewWith1Button:", err)

	}
	dirChooserBtn.Connect("selection-changed", dirSelectionChanged)
	grid.Add(dirChooserBtn)

	// Some GTK callback functions require arguments, such as the
	// 'gchar *uri' argument of GtkLabel's "activate-link" signal.
	// These can be used by using the equivalent go type (in this case,
	// a string) as a closure argument.
	label.Connect("activate-link", func(_ *gtk.Label, uri string) {
		fmt.Println("you clicked a link to:", uri)
	})

	return &grid.Container.Widget
}

func dirSelectionChanged(widget *gtk.FileChooserButton) {
	fmt.Printf("dir selected %s\n", widget.GetFilename())
	appSettings.ImagesDir = widget.GetFilename()
	DoExtract(appSettings)
}
