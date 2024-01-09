package main

import (
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/kkyr/fig"
	"golang.org/x/sys/unix"
	"golang.org/x/term"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strings"
)

func die(err error) {
	if err != nil {
		// Only print stack if DEBUG variable has been defined
		if os.Getenv("DEBUG") != "" {
			debug.PrintStack()
		}
		log.Fatal(err)
	}
}

type Filetype struct {
	Prefix  string
	Suffix  string
	Regex   *regexp.Regexp
	Mime    *regexp.Regexp
	Exec    string
	Exectty string
	// TODO: something that opens the binary in a terminal (for running non-GUI programs in X)
}

type Config struct {
	Filetypes []Filetype
}

func main() {
	configDir, err := os.UserConfigDir()
	die(err)
	configDir = filepath.Join(configDir, "justopen")

	var cfg Config
	err = fig.Load(&cfg, fig.Dirs(configDir))
	die(err)

	if len(os.Args) <= 1 {
		die(fmt.Errorf("Gimme a file or something"))
	}

	var mime *mimetype.MIME = nil
	fileinfo, _ := os.Stat(os.Args[1])

	if fileinfo != nil && fileinfo.IsDir() {
		mimetype.Extend(func(raw []byte, limit uint32) bool { return false }, "inode/directory", "")
		mime = mimetype.Lookup("inode/directory")
	}

	for _, item := range cfg.Filetypes {
		if item.Prefix != "" && !strings.HasPrefix(os.Args[1], item.Prefix) {
			continue
		}

		if item.Suffix != "" && !strings.HasSuffix(os.Args[1], item.Suffix) {
			continue
		}

		if item.Regex != nil && !item.Regex.MatchString(os.Args[1]) {
			continue
		}

		if item.Mime != nil {
			if fileinfo == nil {
				continue
			}

			if mime == nil {
				mime, err = mimetype.DetectFile(os.Args[1])
				die(err)
			}

			if !item.Mime.MatchString(mime.String()) {
				continue
			}
		}

		var cmd string
		if item.Exectty != "" && term.IsTerminal(int(os.Stdin.Fd())) {
			cmd = item.Exectty
		} else {
			cmd = item.Exec
		}

		if strings.Index(cmd, "%f") == -1 {
			cmd = cmd + " %f"
		}

		argv := strings.Split(cmd, " ")
		for i := range argv {
			if argv[i] == "%f" {
				argv[i] = os.Args[1]
			}
		}
		bin, err := exec.LookPath(argv[0])
		die(err)
		unix.Exec(bin, argv, os.Environ())
		break
	}

	fmt.Println("Sorry, no idea what kind of file this is")
}
