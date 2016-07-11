package keys

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/iam"
)

type iamKeySource struct {
	Service    *iam.IAM
	GroupNames []string
}

func NewIAMKeySource(c client.ConfigProvider, groupNames []string) KeySource {
	return &iamKeySource{
		Service:    iam.New(c),
		GroupNames: groupNames,
	}
}

func (ks *iamKeySource) PublicKeys() ([]PublicKey, error) {
	var pkeys []PublicKey

	for _, groupName := range ks.GroupNames {
		users, err := groupUsers(ks.Service, groupName)
		if err != nil {
			return nil, err
		}

		for _, user := range users {
			pks, err := userPublicKeys(ks.Service, *user.UserName)
			if err != nil {
				return nil, err
			}
			for _, pk := range pks {
				pkeys = append(pkeys, PublicKey{
					Body: fmt.Sprintf("%s %s@IAM", strings.TrimSpace(*pk.SSHPublicKeyBody), *pk.UserName),
				})
			}
		}
	}

	return pkeys, nil
}

func groupUsers(svc *iam.IAM, groupName string) ([]*iam.User, error) {
	var users []*iam.User
	input := &iam.GetGroupInput{
		GroupName: aws.String(groupName),
	}
	err := svc.GetGroupPages(input, func(output *iam.GetGroupOutput, lastPage bool) bool {
		users = append(users, output.Users...)
		return true
	})

	return users, err
}

func userPublicKeys(svc *iam.IAM, userName string) ([]*iam.SSHPublicKey, error) {
	var pks []*iam.SSHPublicKey
	listInput := &iam.ListSSHPublicKeysInput{
		UserName: aws.String(userName),
	}
	var getErr error
	err := svc.ListSSHPublicKeysPages(listInput, func(listOutput *iam.ListSSHPublicKeysOutput, lastPage bool) bool {
		for _, pk := range listOutput.SSHPublicKeys {
			if *pk.Status == "Active" {
				getInput := &iam.GetSSHPublicKeyInput{
					SSHPublicKeyId: pk.SSHPublicKeyId,
					UserName:       pk.UserName,
					Encoding:       aws.String("SSH"),
				}
				var getOutput *iam.GetSSHPublicKeyOutput
				getOutput, getErr = svc.GetSSHPublicKey(getInput)
				if getErr != nil {
					return false
				}
				pks = append(pks, getOutput.SSHPublicKey)
			}
		}

		return true
	})

	if getErr != nil {
		return nil, getErr
	}

	return pks, err
}
