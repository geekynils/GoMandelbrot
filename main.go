package main

import (
	// "flag"
	"github.com/nightlifelover/GoMandelbrot/inputoutput"
	// "log"
	"github.com/nightlifelover/GoMandelbrot/mandelbrot"
	// "os"
	"runtime"
	//"runtime/pprof"
)

// var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {

	/*
		flag.Parse()

		if *cpuprofile != "" {
			f, err := os.Create(*cpuprofile)
			if err != nil {
				log.Fatal(err)
			}
			pprof.StartCPUProfile(f)

			defer pprof.StopCPUProfile()
		}
	*/

	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	inputoutput.Init()

	mandelbrot.LinkScreenOutput(inputoutput.ScreenOutputChan)
	mandelbrot.LinkInput(inputoutput.InputChan)
	mandelbrot.LinkNThreads(inputoutput.NThreadsChan)

	go mandelbrot.DrawMandelbrot()

	inputoutput.Run()
}
