package hashconv

import (
	"errors"

	cid "github.com/ipfs/go-cid"
)

func ConvertQmHashToCid(hash string) (cid.Cid, error) {
	if hash == "" {
		return cid.Cid{}, errors.New("hash is empty")
	}
	if len(hash) != 46 {
		return cid.Cid{}, errors.New("hash is not 46 length")
	}
	if hash[0:2] != "Qm" {
		return cid.Cid{}, errors.New("hash is not Qm")
	}
	return cid.Decode(hash)
}

func ConvertQmHashToBarfy(hash string) (*cid.Cid, error) {
	if hash == "" {
		return nil, errors.New("hash is empty")
	}
	if len(hash) != 46 {
		return nil, errors.New("hash is not 46 length")
	}
	if hash[0:2] != "Qm" {
		return nil, errors.New("hash is not Qm")
	}
	// cidStr := hash[2:]
	c, err := cid.Decode(hash)
	if err != nil {
		return nil, err
	}
	cf := cid.NewCidV1(c.Type(), c.Hash())
	return &cf, nil
}
