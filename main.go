package main

import (
	"github.com/wowgroup/mrgrgr/install"
	"github.com/wowgroup/mrgrgr/keys"

	"os"
	"os/user"
	"path"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/urfave/cli"
)

const (
	flagAuthorizedKeysPath = "authorized-keys-path"
	flagIamGroup           = "iam-group"
	flagPollInterval       = "poll-interval"
)

func main() {
	app := cli.NewApp()
	app.Name = "mrgrgr"
	app.Version = "0.1.0"
	app.Action = doMain
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  flagAuthorizedKeysPath,
			Usage: "Path to authorized_keys",
			Value: path.Join(currentUser().HomeDir, ".ssh", "authorized_keys"),
		},
		cli.StringSliceFlag{
			Name:  flagIamGroup,
			Usage: "IAM group(s) from which to source keys",
			Value: &cli.StringSlice{},
		},
		cli.IntFlag{
			Name:  flagPollInterval,
			Usage: "Interval (seconds) between polling attempts",
			Value: 900,
		},
	}
	app.Run(os.Args)
}

func doMain(c *cli.Context) {
	authorizedKeysPath := c.String(flagAuthorizedKeysPath)
	iamGroups := c.StringSlice(flagIamGroup)
	pollInterval := c.Int(flagPollInterval)

	installer := install.NewAuthorizedKeysInstaller(authorizedKeysPath)
	session := session.New()
	keySources := []keys.KeySource{
		keys.NewIAMKeySource(session, iamGroups),
		keys.NewEC2KeySource(session),
	}
	tick := time.NewTicker(time.Duration(pollInterval) * time.Second).C

	for {
		var pkeys []keys.PublicKey
		for _, ks := range keySources {
			pks, err := ks.PublicKeys()
			if err != nil {
				log.Fatal(err)
			}
			pkeys = append(pkeys, pks...)
		}

		err := installer.Install(pkeys)
		if err != nil {
			log.Fatal(err)
		}

		log.Infof("refreshed %d key(s)", len(pkeys))

		<-tick
	}
}

func currentUser() *user.User {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return u
}
