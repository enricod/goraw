package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/enricod/goraw/libraw"
	"github.com/gotk3/gotk3/gtk"
	"github.com/nfnt/resize"
	bolt "go.etcd.io/bbolt"
)

type Settings struct {
	ImagesDir string
}

var appSettings Settings
var flowbox *gtk.FlowBox

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

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) > 0 {
		appSettings = Settings{ImagesDir: argsWithoutProg[0]}
		fmt.Printf("selezionata dir %s", appSettings.ImagesDir)
	} else {
		appSettings = Settings{ImagesDir: "."}
	}
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

	DoExtract(appSettings.ImagesDir)
	//fmt.Scanln()
	fmt.Println("done")
}

func DoExtract(dirname string) {
	exportPath := dirname + "/_export"
	if _, err := os.Stat(exportPath); os.IsNotExist(err) {
		os.Mkdir(exportPath, 0777)
	}

	exportPath_t := dirname + "/_export_t"
	if _, err := os.Stat(exportPath_t); os.IsNotExist(err) {
		os.Mkdir(exportPath_t, 0777)
	}

	db, err := bolt.Open(dirname+"/_grbolt.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		//fmt.Printf("%s", f.Name())
		if IsStringInSlice(filepath.Ext(f.Name()), extensions()) {
			exportedImage, _ := libraw.ExportEmbeddedJPEG(dirname, f, exportPath)
			jpegImg, err2 := readJPEG(exportedImage)
			if err2 != nil {
				panic(err)
			}
			newImage := resize.Resize(160, 0, *jpegImg, resize.Lanczos3)

			// Encode uses a Writer, use a Buffer if you need the raw []byte

			thumbfile, err3 := os.Create(fmt.Sprintf("%s/%s", exportPath_t, f.Name()+"_t.jpg"))
			if err3 != nil {
				panic(err)
			}
			defer thumbfile.Close()

			err = jpeg.Encode(thumbfile, newImage, nil)

		}
	}

	processed := []string{}
	files, err = ioutil.ReadDir(exportPath_t)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), "_t.jpg") {
			processed = append(processed, exportPath_t+"/"+f.Name())

		}
	}

	mostraImmagini(processed, flowbox)
}

func readJPEG(filename string) (*image.Image, error) {
	existingImageFile, err := os.Open(filename)
	if err != nil {
		// Handle error
	}
	defer existingImageFile.Close()

	// Calling the generic image.Decode() will tell give us the data
	// and type of image it is as a string. We expect "png"
	//imageData, imageType, err := image.Decode(existingImageFile)
	//if err != nil {
	// Handle error
	//}
	//fmt.Println(imageData)
	//fmt.Println(imageType)

	// We only need this because we already read from the file
	// We have to reset the file pointer back to beginning
	existingImageFile.Seek(0, 0)

	// Alternatively, since we know it is a png already
	// we can call png.Decode() directly
	loadedImage, err := jpeg.Decode(existingImageFile)
	if err != nil {
		// Handle error
	}
	return &loadedImage, nil
}

func mainPanel() *gtk.Widget {

	horBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
	if err != nil {
		log.Fatal("Unable to create horBox:", err)
	}
	horBox.SetHomogeneous(false)

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
	/*
		pb, err := gtk.ProgressBarNew()
		if err != nil {
			log.Fatal("Unable to create progress bar:", err)
		}
	*/
	grid.Add(sb)
	sb.SetHExpand(true)
	//grid.AttachNextTo(pb, sb, gtk.POS_RIGHT, 1, 1)
	label.SetHExpand(true)

	// Pass in a ProgressBar and the target SpinButton as user data rather
	// than using the sb and pb variables scoped to the anonymous func.
	// This can be useful when passing in a closure that has already been
	// generated, but when you still wish to connect the callback with some
	// variables only visible in this scope.
	/*
		sb.Connect("value-changed", func(sb *gtk.SpinButton, pb *gtk.ProgressBar) {
			pb.SetFraction(sb.GetValue() / 1)
		}, pb)
	*/
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

	horBox.PackStart(grid, false, true, 6)

	flowbox, err = gtk.FlowBoxNew()
	if err != nil {
		log.Fatal("Unable to create FileChooserDialogNewWith1Button:", err)

	}

	//popolaFlowbox(flowbox)
	horBox.PackStart(flowbox, true, true, 6)
	return &horBox.Container.Widget
	//return &grid.Container.Widget
}

func mostraImmagini(immagini []string, flowbox *gtk.FlowBox) {
	for _, color := range immagini {
		img, err := gtk.ImageNewFromFile(color)
		if err != nil {
			log.Fatal("Unable to create FileChooserDialogNewWith1Button:", err)

		}
		flowbox.Add(img)
	}
	flowbox.ShowAll()
}

func popolaFlowbox(flowbox *gtk.FlowBox) {
	colors := []string{
		"AliceBlue",
		"AntiqueWhite",
		"AntiqueWhite1",
		"AntiqueWhite2",
		"AntiqueWhite3",
		"AntiqueWhite4",
		"aqua",
		"aquamarine",
		"aquamarine1",
		"aquamarine2",
		"aquamarine3",
		"aquamarine4",
		"azure",
		"azure1",
		"azure2",
		"azure3",
		"azure4",
		"beige",
		"bisque",
		"bisque1",
		"bisque2",
		"bisque3",
		"bisque4",
		"black",
		"BlanchedAlmond",
		"blue",
		"blue1",
		"blue2",
		"blue3",
		"blue4",
		"BlueViolet",
		"brown",
		"brown1",
		"brown2",
		"brown3",
		"brown4",
		"burlywood",
		"burlywood1",
		"burlywood2",
		"burlywood3",
		"burlywood4",
		"CadetBlue",
		"CadetBlue1",
		"CadetBlue2",
		"CadetBlue3",
		"CadetBlue4",
		"chartreuse",
		"chartreuse1",
		"chartreuse2",
		"chartreuse3",
		"chartreuse4",
		"chocolate",
		"chocolate1",
		"chocolate2",
		"chocolate3",
		"chocolate4",
		"coral",
		"coral1",
		"coral2",
		"coral3",
		"coral4",
	}

	for _, color := range colors {
		button, err := gtk.ButtonNew()
		button.SetLabel(color)
		if err != nil {
			log.Fatal("Unable to create FileChooserDialogNewWith1Button:", err)

		}
		flowbox.Add(button)
	}
}

func dirSelectionChanged(widget *gtk.FileChooserButton) {
	fmt.Printf("dir selected %s\n", widget.GetFilename())
	appSettings.ImagesDir = widget.GetFilename()
	DoExtract(widget.GetFilename())
}
