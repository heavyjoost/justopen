package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/kkyr/fig"
	"golang.org/x/sys/unix"
	"golang.org/x/term"
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
	Filetypes     []Filetype
	CaseSensitive bool
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
	target_path := os.Args[1]
	if !cfg.CaseSensitive {
		target_path = strings.ToLower(target_path)
	}

	var mime *mimetype.MIME = nil
	fileinfo, _ := os.Stat(os.Args[1])

	if fileinfo != nil && fileinfo.IsDir() {
		mimetype.Extend(func(raw []byte, limit uint32) bool { return false }, "inode/directory", "")
		mime = mimetype.Lookup("inode/directory")
	}

	for _, item := range cfg.Filetypes {
		if item.Prefix != "" && !strings.HasPrefix(target_path, item.Prefix) {
			continue
		}

		if item.Suffix != "" && !strings.HasSuffix(target_path, item.Suffix) {
			continue
		}

		if item.Regex != nil && !item.Regex.MatchString(target_path) {
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
			if strings.Contains(argv[i], "%f") {
				argv[i] = strings.ReplaceAll(argv[i], "%f", os.Args[1])
			}
		}
		bin, err := exec.LookPath(argv[0])
		die(err)
		unix.Exec(bin, argv, os.Environ())
		break
	}

	fmt.Println("Sorry, no idea what kind of file this is")
}
