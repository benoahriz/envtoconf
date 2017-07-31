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

func writeFile(data *bytes.Buffer, outfile string) {
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		_, err := io.Copy(pw, data)
		if err != nil {
			log.Fatalf("%s\n", err)
		}
	}()
	out, err := os.Create(outfile)
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	defer func() {
		cerr := out.Close()
		if cerr != nil {
			log.Fatalf("%s\n", cerr)
		}
	}()
	if _, err = io.Copy(out, pr); err != nil {
		return
	}
	err = out.Sync()
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	return
}

func parsetemplate(infile string) (b *bytes.Buffer, err error) {
	if _, err := os.Stat(infile); !os.IsNotExist(err) {
		if err != nil {
			log.Fatalf("%s\n", err)
		}
	}
	tpl, err := template.New(infile).Funcs(sprig.TxtFuncMap()).ParseFiles(infile)
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	b = new(bytes.Buffer)
	err = tpl.ExecuteTemplate(b, filepath.Base(infile), nil)
	log.Debugln(b.String())
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	return b, err
}

func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("0.0.1").Author("Benjamin Rizkowsky")
	kingpin.CommandLine.Help = "A simple environment variable config template writer."
	kingpin.Parse()
	if *verbose == true {
		log.SetLevel(log.DebugLevel)
	}
	data, err := parsetemplate(*intemplate)
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	writeFile(data, *outfile)
}
