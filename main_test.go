package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParsetemplate(t *testing.T) {
	Convey("Check that parsetemplate returns a bytes.buffer of rendered template", t, func() {
		os.Setenv("FOO", "bar")
		tmpDir, _ := ioutil.TempDir("/tmp/", "envtoconf")
		tmpfile, _ := ioutil.TempFile(tmpDir, "template")
		defer os.RemoveAll(tmpDir)
		file := tmpfile.Name()
		testdata := `home path is: {{env "HOME"}} and foo: {{env "FOO"}}`
		if err := ioutil.WriteFile(file, []byte(testdata), 0644); err != nil {
			t.Fatalf("WriteFile %s: %v", file, err)
		}
		log.Printf("File is: %s\n", file)

		data, err := parsetemplate(file)
		check(err)
		So(data, ShouldHaveSameTypeAs, new(bytes.Buffer))
		So(data.String(), ShouldContainSubstring, `bar`)
	})

}
