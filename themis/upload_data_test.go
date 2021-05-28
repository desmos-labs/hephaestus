package themis_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/desmos-labs/hephaestus/themis"
	"github.com/desmos-labs/hephaestus/types"
	"github.com/desmos-labs/hephaestus/utils"
)

func TestUploadData(t *testing.T) {
	privKey, err := utils.ReadPrivateKeyFromFile("/home/riccardo/Coding/Forbole/Hephaestus/.data/hephaestus.priv")
	require.NoError(t, err)

	data := types.NewConnectionData("8902A4822B87C1ADED60AE947044E614BD4CAEE2",
		"033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d",
		"Riccardo Montagnin#5414",
		"d10db146bb4d234c5c1d2bc088e045f4f05837c690bce4101e2c0f0c6c96e1232d8516884b0a694ee85e9c9da51be74966886cbb12af4ad87e5336da76d75cfb",
	)

	err = themis.UploadData(data, "http://localhost:5000", privKey)
	require.NoError(t, err)
}
