package libraw

// #cgo LDFLAGS: -lraw
// #include "libraw/libraw.h"
import "C"
import (
	"fmt"
	"os"
)

func handleError(msg string, err int) {

	if err != 0 {
		fmt.Printf("ERROR libraw  %v\n", C.libraw_strerror(C.int(err)))
	}
}

func lrInit() *C.libraw_data_t {
	return C.libraw_init(0)
}

func ExportThumb(inputPath string, inputfile os.FileInfo, exportPath string) (string, error) {

	outfile := exportPath + "/" + inputfile.Name() + "_embedded.jpg"
	infile := inputPath + "/" + inputfile.Name()

	if _, err := os.Stat(outfile); os.IsNotExist(err) {
		iprc := lrInit()
		C.libraw_open_file(iprc, C.CString(infile))

		ret := C.libraw_unpack_thumb(iprc)
		handleError("unpack thumb", int(ret))

		//ret = C.libraw_dcraw_process(iprc)
		//handleError("process", int(ret))
		//iprc.params.output_tiff = 1
		//outfile := exportPath + "/" + inputfile.Name() + ".tiff"

		fmt.Printf("exporting %s  ->  %s \n", inputfile.Name(), outfile)
		ret = C.libraw_dcraw_thumb_writer(iprc, C.CString(outfile))

		handleError("save thumb", int(ret))

		lrClose(iprc)
	}
	return outfile, nil
}

func Export(inputPath string, inputfile os.FileInfo, exportPath string) error {

	// FIXME controllare che file input esiste

	// lanciare errore se file input non esiste

	outfile := exportPath + "/" + inputfile.Name() + ".ppm"
	infile := inputPath + "/" + inputfile.Name()

	if _, err := os.Stat(outfile); os.IsNotExist(err) {
		iprc := lrInit()
		C.libraw_open_file(iprc, C.CString(infile))

		ret := C.libraw_unpack(iprc)
		handleError("unpack", int(ret))

		ret = C.libraw_dcraw_process(iprc)

		handleError("dcraw processing", int(ret))
		//iprc.params.output_tiff = 1
		//outfile := exportPath + "/" + inputfile.Name() + ".tiff"

		fmt.Printf("exporting %s  ->  %s \n", inputfile.Name(), outfile)
		ret = C.libraw_dcraw_ppm_tiff_writer(iprc, C.CString(outfile))

		handleError("save ppm", int(ret))

		lrClose(iprc)
	}
	return nil
}

func lrClose(iprc *C.libraw_data_t) {
	C.libraw_close(iprc)
}
