// Package pgtest starts and stops a postgres server, quickly
// and conveniently, for gocheck unit tests.
package pgtest

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"

	gc "gopkg.in/check.v1"
)

var conf = template.Must(template.New("t").Parse(`
fsync = off
listen_addresses = ''

{{if .Plural}}
unix_socket_directories = '{{.ConfDir}}'
{{else}}
unix_socket_directory = '{{.ConfDir}}'
{{end}}

`))

type PGSuite struct {
	// initdbDir will contain the path to the default file set produced by
	// initdb.
	initdbDir string
	// postgres is the discovered path to the postgres binary.
	postgres string

	URL string // Connection URL for sql.Open.
	Dir string

	cmd *exec.Cmd
}

func (s *PGSuite) SetUpSuite(c *gc.C) {
	s.initdbDir = c.MkDir()
	out, err := exec.Command("pg_config", "--bindir").Output()
	c.Assert(err, gc.IsNil, gc.Commentf("pg_config"))

	bindir := string(bytes.TrimSpace(out))
	s.postgres = filepath.Join(bindir, "postgres")
	initdb := filepath.Join(bindir, "initdb")
	err = exec.Command(initdb, "-D", s.initdbDir).Run()
	c.Assert(err, gc.IsNil, gc.Commentf("initdb"))
}

// SetUpTest runs postgres in a temporary directory,
// with a default file set produced by initdb.
// If an error occurs, the test will fail.
func (s *PGSuite) SetUpTest(c *gc.C) {
	s.Dir = c.MkDir()

	err := exec.Command("cp", "-a", s.initdbDir+"/.", s.Dir).Run()
	c.Assert(err, gc.IsNil)

	path := filepath.Join(s.Dir, "postgresql.conf")
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0666)
	c.Assert(err, gc.IsNil)

	plural := !contains("unix_socket_directory", path)
	err = conf.Execute(f, struct {
		ConfDir string
		Plural  bool
	}{s.Dir, plural})
	c.Assert(err, gc.IsNil)

	err = f.Close()
	c.Assert(err, gc.IsNil)

	s.URL = "host=" + s.Dir + " dbname=postgres sslmode=disable"
	s.cmd = exec.Command(s.postgres, "-D", s.Dir)
	err = s.cmd.Start()
	c.Assert(err, gc.IsNil, gc.Commentf("starting postgres"))

	c.Log("starting postgres in", s.Dir)
	sock := filepath.Join(s.Dir, ".s.PGSQL.5432")
	for n := 0; n < 20; n++ {
		if _, err := os.Stat(sock); err == nil {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	c.Fatal("timeout waiting for postgres to start")
	panic("unreached")
}

// TearDownTest stops the running postgres process and removes its
// temporary data directory.
// If an error occurs, the test will fail.
func (s *PGSuite) TearDownTest(c *gc.C) {
	err := s.cmd.Process.Signal(os.Interrupt)
	c.Assert(err, gc.IsNil)
	err = s.cmd.Wait()
	c.Assert(err, gc.IsNil)
	err = os.RemoveAll(s.Dir)
	c.Assert(err, gc.IsNil)
}

func contains(substr, name string) bool {
	b, err := ioutil.ReadFile(name)
	return err == nil && bytes.Contains(b, []byte(substr))
}
