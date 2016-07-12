package install

import (
	"github.com/wowgroup/mrgrgr/keys"

	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"syscall"
)

type authorizedKeysInstaller struct {
	Path string
}

func NewAuthorizedKeysInstaller(path string) Installer {
	return &authorizedKeysInstaller{
		Path: path,
	}
}

func (i *authorizedKeysInstaller) Install(pkeys []keys.PublicKey) error {
	f, err := ioutil.TempFile(path.Dir(i.Path), path.Base(i.Path))
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	defer f.Close()

	err = copyOwnership(i.Path, f)
	if err != nil {
		return err
	}

	err = writeKeys(f, pkeys)
	if err != nil {
		return err
	}

	return os.Rename(f.Name(), i.Path)
}

func writeKeys(f *os.File, pkeys []keys.PublicKey) error {
	w := bufio.NewWriter(f)

	for _, pk := range pkeys {
		_, err := w.WriteString(fmt.Sprintf("%s\n", pk.Body))
		if err != nil {
			return err
		}
	}

	return w.Flush()
}

func copyOwnership(srcPath string, dest *os.File) error {
	fi, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	stat := fi.Sys().(*syscall.Stat_t)
	if stat == nil {
		return fmt.Errorf("could not stat %s (got nil)", srcPath)
	}

	return dest.Chown(int(stat.Uid), int(stat.Gid))
}
