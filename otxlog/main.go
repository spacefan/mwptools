package main

import (
	"os"
	"flag"
	"fmt"
)

var (
	baud    = flag.Int("b", 115200, "Baud rate")
	device  = flag.String("d", "", "LTM to serial device")
	fd      = flag.Int("fd", -1, "LTM to file descriptor")
	output  = flag.String("out", "", "output LTM to file")
	gpxout  = flag.String("gpx", "", "write gpx to file")
	dump    = flag.Bool("dump", false, "dump headers & exit")
	fast    = flag.Bool("fast", false, "fast replay")
	verbose = flag.Bool("verbose", false, "verbose LTM debug")
	armed   = flag.Bool("armed-only", false, "skip not armed")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of otxlog [options] [files ...]\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	files := flag.Args()
	if len(files) == 0 {
		flag.Usage()
		os.Exit(-1)
	}

	o := NewOTX()

	if *dump {
		o.Set_dump()
	} else if *gpxout != "" {
		o.Gpx_init(*gpxout)
	} else {
		var s *MSPSerial
		if *fd > 0 {
			s = NewMSPFd(*fd)
		} else if len(*device) > 0 {
			s = NewMSPSerial(*device, *baud)
		} else if len(*output) > 0 {
			s = NewMSPFile(*output)
		} else {
			fmt.Fprintf(os.Stderr, "No output specified\n")
			os.Exit(-1)
		}
		o.Stream_init(s)
	}
	o.Verbose(*verbose)
	o.Armed(*armed)
	o.Reader(files[0], *fast)
}
