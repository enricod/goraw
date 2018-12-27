package libraw

// #cgo LDFLAGS: -lraw
// #include "libraw/libraw.h"
import "C"
import (
	"fmt"
	"os"
)

func handleError(err int) {
	if err != 0 {
		fmt.Printf("ERROR libraw  %v\n", C.libraw_strerror(C.int(err)))
	}
}

func LRInit() *C.libraw_data_t {
	return C.libraw_init(0)
}

func Export(path string, inputfile os.FileInfo, exportPath string) {

	iprc := LRInit()
	C.libraw_open_file(iprc, C.CString(path+"/"+inputfile.Name()))

	//fmt.Printf("Processing %s\n", inputfile.Name())

	ret := C.libraw_unpack(iprc)
	handleError(int(ret))

	ret = C.libraw_dcraw_process(iprc)
	handleError(int(ret))
	//iprc.params.output_tiff = 1
	//outfile := exportPath + "/" + inputfile.Name() + ".tiff"
	outfile := exportPath + "/" + inputfile.Name() + ".ppm"
	fmt.Printf("exporting %s  ->  %s \n", inputfile.Name(), outfile)
	ret = C.libraw_dcraw_ppm_tiff_writer(iprc, C.CString(outfile))
	handleError(int(ret))

	Close(iprc)

}

func Close(iprc *C.libraw_data_t) {
	C.libraw_close(iprc)
}
