package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	intemplate = kingpin.Flag("template", "the source template").Default("file.tpl").String()
	outfile    = kingpin.Flag("outfile", "output file").Default("outfile.txt").String()
	verbose    = kingpin.Flag("verbose", "Produce verbose output").Short('v').Default("false").Bool()
)

func check(e error) {
	if e != nil {
		log.Fatalf("%s\n", e)
		panic(e)
	}
}

func writeFile(data *bytes.Buffer, outfile string) {
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		_, err := io.Copy(pw, data)
		check(err)
	}()
	out, err := os.Create(outfile)
	check(err)
	defer func() {
		cerr := out.Close()
		check(cerr)
	}()
	if _, err = io.Copy(out, pr); err != nil {
		return
	}
	err = out.Sync()
	return
}

func parsetemplate(infile string) (b *bytes.Buffer, err error) {
	if _, err := os.Stat(infile); !os.IsNotExist(err) {
		check(err)
	}
	tpl, err := template.New(infile).Funcs(sprig.TxtFuncMap()).ParseFiles(infile)
	check(err)
	b = new(bytes.Buffer)
	err = tpl.ExecuteTemplate(b, filepath.Base(infile), nil)
	log.Debugln(b.String())
	check(err)
	return b, err
}

func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("0.0.1").Author("Benjamin Rizkowsky")
	kingpin.CommandLine.Help = "A simple environment variable config template writer."
	kingpin.Parse()
	if *verbose == true {
		log.SetLevel(log.DebugLevel)
	}
	parsetemplate(*intemplate)
	data, err := parsetemplate(*intemplate)
	check(err)
	writeFile(data, *outfile)
}
