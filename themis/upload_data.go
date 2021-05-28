package themis

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"net/http"

	themisdiscord "github.com/desmos-labs/themis/apis/discord"
	"github.com/tendermint/tendermint/libs/json"

	"github.com/desmos-labs/hephaestus/types"
)

// UploadData uploads the given data using the Themis APIs hosted at the provided host,
// after signing the data using the given private key
func UploadData(data *types.ConnectionData, host string, privKey *rsa.PrivateKey) error {
	verData := themisdiscord.VerificationData{
		Address:   data.Address,
		PubKey:    data.PubKey,
		Username:  data.Username,
		Signature: data.Signature,
	}

	bz, err := verData.ToSignBytes()
	if err != nil {
		return err
	}

	signature, err := rsa.SignPSS(rand.Reader, privKey, crypto.SHA256, bz, nil)
	if err != nil {
		return err
	}

	bodyBz, err := json.Marshal(&themisdiscord.SaveDataReq{
		Signature:        hex.EncodeToString(signature),
		VerificationData: verData,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/discord/data", host), bytes.NewBuffer(bodyBz))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid Themis request: %s", resp.Status)
	}

	return nil
}
