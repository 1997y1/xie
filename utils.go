// Source code file, created by Developer@YAN_YING_SONG.
//
// 😈 邪修の手段です！
//
// XIE is an evil form of cultivation method.

package xie

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"unsafe"

	"xie/go/errcause"
)

const KB = 1024
const MB = KB * KB
const GB = MB * KB

func SizeFormat(bytes float64) string {
	if bytes >= GB {
		return fmt.Sprintf("%.1fG", bytes/GB)
	} else if bytes >= MB {
		return fmt.Sprintf("%.1fM", bytes/MB)
	} else if bytes >= KB {
		return fmt.Sprintf("%.1fK", bytes/KB)
	} else {
		return fmt.Sprintf("%.0fB", bytes)
	}
}

func MapFor[k comparable, v any](data map[k]v, item func(k, v) (_ bool)) {
	for key, value := range data {
		if item(key, value) {
			break
		}
	}
}

func CopyBytes(b []byte) []byte {
	buf := make([]byte, len(b))
	copy(buf, b)
	return buf
}

func ExeFp() string {
	exe := os.Args[0]
	if !filepath.IsAbs(exe) {
		var err error
		exe, err = filepath.Abs(exe)
		if err != nil {
			panic(err) // Must succeed.
		}
	}

	return exe
}

func JoinFp(elem ...string) string {
	values := make([]string, 0, 100)
	values = append(values, filepath.Dir(ExeFp()))
	values = append(values, elem...)
	return filepath.Join(values...)
}

func CoverFile(fp string, data []byte, dump ...*error) {
	// Cover data to file.

	dir := filepath.Dir(fp)
	if !FileExist(dir) {
		_ = os.MkdirAll(dir, 0o777)
	}

	if err := os.WriteFile(fp, data, 0o777); err != nil {
		stderr(err)
		if len(dump) == 0 {
			panic(err) // Must succeed
		} else {
			*dump[0] = err
		}
	}
}

func CmdEnter(dir, cmdline, shell string) ([]byte, error) {
	// Execute and read the command-line results.

	if dir == "" {
		dir = filepath.Dir(ExeFp())
	}

	// Open shell session.
	cmd := exec.Command(shell)
	cmd.Dir = dir

	// Read stdout & stderr.
	buf := &bytes.Buffer{}
	cmd.Stdout, cmd.Stderr = buf, buf

	// Write commandline.
	in, err := cmd.StdinPipe()
	if err != nil {
		return nil, errcause.LinkErr(err, "cmd.StdinPipe")
	}
	if err = cmd.Start(); err != nil {
		return nil, errcause.LinkErr(err, "cmd.Start")
	}
	if _, err = in.Write([]byte(cmdline)); err != nil {
		return nil, errcause.LinkErr(err, "cmd.Write")
	}
	_, _ = in.Write([]byte("\nexit\n"))
	_ = in.Close()

	// Wait command to finish.
	if err = cmd.Wait(); err != nil {
		stderr(err)
	}

	return buf.Bytes(), nil
}

func CatString(fp string) string {
	// Cat file data to string.

	b := CatBytes(fp)
	return Ts(b)
}

func CatBytes(fp string) []byte {
	// Cat file data to bytes.

	source, err := os.Open(fp)
	if source != nil {
		defer func() { _ = source.Close() }()
	}
	if err != nil {
		stderr(err)
		return nil
	}
	info, _ := source.Stat()

	// Read All.
	data := make([]byte, info.Size())
	if len(data) > 0 {
		if _, err = io.ReadFull(source, data); err != nil {
			stderr(err)
		}
	} else {
		data, err = io.ReadAll(source)
		if err != nil {
			stderr(err)
		}
	}

	return data
}

func FileExist(fp string) bool {
	// Do file exist?

	_, err := os.Lstat(fp)
	return !os.IsNotExist(err)
}

func Ts(a []byte) string {
	// Bytes value to string with zero copy.

	return *(*string)(unsafe.Pointer(&a))
}

func stderr(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "stderr: %s\n", err)
}
