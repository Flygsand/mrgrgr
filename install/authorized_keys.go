package install

import (
	"github.com/wowgroup/mrgrgr/keys"

	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
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
