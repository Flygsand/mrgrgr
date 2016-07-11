package keys

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
)

type ec2KeySource struct {
	Service *ec2metadata.EC2Metadata
}

func NewEC2KeySource(c client.ConfigProvider) KeySource {
	return &ec2KeySource{
		Service: ec2metadata.New(c),
	}
}

func (ks *ec2KeySource) PublicKeys() ([]PublicKey, error) {
	if !ks.Service.Available() {
		return nil, fmt.Errorf("EC2 metadata is not available")
	}

	var pks []PublicKey
	keyList, err := ks.Service.GetMetadata("public-keys/")
	if err != nil {
		return nil, err
	}

	for _, l := range strings.Split(keyList, "\n") {
		keyId := strings.SplitN(l, "=", 2)[0]
		body, err := ks.Service.GetMetadata(fmt.Sprintf("public-keys/%s/openssh-key", keyId))
		if err != nil {
			return nil, err
		}
		pks = append(pks, PublicKey{
			Body: strings.TrimSpace(body),
		})
	}

	return pks, nil
}
