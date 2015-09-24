package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"text/tabwriter"
)

type byteSize int64

const (
	_           = iota
	KB byteSize = 1 << (10 * iota)
	MB
	GB
)

func (b byteSize) String() string {
	switch {
	case b >= GB:
		return b.fmtString("GB", GB)
	case b >= MB:
		return b.fmtString("MB", MB)
	case b >= KB:
		return b.fmtString("KB", KB)
	}
	return b.fmtString("B", b)
}

func (b byteSize) fmtString(name string, divisor byteSize) string {
	return fmt.Sprintf("%.2f"+name, float64(b)/float64(divisor))
}

type fileSize struct {
	path string
	byteSize
}

func (f fileSize) String() string {
	return fmt.Sprintf("%s\t%s\t%s\t", f.byteSize.String(), filepath.Ext(f.path), f.path)
}

type fileSizes []fileSize

func (f fileSizes) Len() int           { return len(f) }
func (f fileSizes) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f fileSizes) Less(i, j int) bool { return f[i].byteSize > f[j].byteSize }

var allFileSizes fileSizes
var tabw *tabwriter.Writer

func main() {
	args := struct {
		root         string
		num          int
		rightjustify bool
	}{}
	flag.StringVar(&args.root, "root", "/", "The root directory to run from.")
	flag.IntVar(&args.num, "top", 10, "The top number of files to output.")
	flag.BoolVar(&args.rightjustify, "rightjustify", false, "Align file paths to the right in output")
	flag.Parse()

	tabw = new(tabwriter.Writer)
	tabw.Init(os.Stdout, 8, 0, 1, ' ', tabwriter.AlignRight)
	allFileSizes = make(fileSizes, 0)

	err := filepath.Walk(args.root, mark)
	if err != nil {
		log.Fatal(err)
	}
	sort.Sort(allFileSizes)
	for i := 0; i < args.num; i++ {
		fmt.Fprintln(tabw, allFileSizes[i])
		if !args.rightjustify {
			tabw.Flush()
		}
	}
	if args.rightjustify {
		tabw.Flush()
	}
}

func mark(path string, info os.FileInfo, err error) error {
	if info.IsDir() || err == os.ErrPermission {
		return nil
	}
	if err != nil {
		return err
	}
	allFileSizes = append(allFileSizes, fileSize{path + info.Name(), byteSize(info.Size())})
	return nil
}
