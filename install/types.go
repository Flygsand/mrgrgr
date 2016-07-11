package install

import (
	"github.com/wowgroup/mrgrgr/keys"
)

type Installer interface {
	Install(pkeys []keys.PublicKey) error
}
