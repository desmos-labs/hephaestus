package limitations_test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/desmos-labs/hephaestus/limitations"

	"github.com/andersfylling/disgord"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/libs/json"
)

func TestLimitationTestSuite(t *testing.T) {
	suite.Run(t, new(LimitationsTestSuite))
}

type LimitationsTestSuite struct {
	suite.Suite

	tempFile string
}

func (suite *LimitationsTestSuite) SetupTest() {
	file, err := ioutil.TempFile(os.TempDir(), "test_")
	suite.Require().NoError(err)
	suite.tempFile = file.Name()
}

func (suite *LimitationsTestSuite) TearDownTest() {
	err := os.Remove(suite.tempFile)
	suite.Require().NoError(err)
}

func (suite *LimitationsTestSuite) TestGetLimitationExpiration() {
	date := time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC)
	usecases := []struct {
		name              string
		storedLimitations map[string]limitations.UserLimitations
		userID            disgord.Snowflake
		command           string
		expErr            bool
		expTime           *time.Time
	}{
		{
			name:              "empty limitations returns nil",
			storedLimitations: nil,
			userID:            disgord.Snowflake(1),
			command:           "send",
			expErr:            false,
			expTime:           nil,
		},
		{
			name: "non existing user returns nil",
			storedLimitations: map[string]limitations.UserLimitations{
				"2": {
					CommandsLimitations: map[string]time.Time{
						"send": date,
					},
				},
			},
			userID:  disgord.Snowflake(1),
			command: "send",
			expErr:  false,
			expTime: nil,
		},
		{
			name: "non existing command returns nil",
			storedLimitations: map[string]limitations.UserLimitations{
				"1": {
					CommandsLimitations: map[string]time.Time{
						"test": date,
					},
				},
			},
			userID:  disgord.Snowflake(1),
			command: "send",
			expErr:  false,
			expTime: nil,
		},
		{
			name: "existing user and command returs correct date",
			storedLimitations: map[string]limitations.UserLimitations{
				"1": {
					CommandsLimitations: map[string]time.Time{
						"send": date,
					},
				},
			},
			userID:  disgord.Snowflake(1),
			command: "send",
			expErr:  false,
			expTime: &date,
		},
	}

	for _, uc := range usecases {
		uc := uc
		suite.Run(uc.name, func() {
			suite.SetupTest()

			// Setup the temp file
			limitations.SetLimitationsFile(suite.tempFile)

			if uc.storedLimitations != nil {
				bz, err := json.Marshal(uc.storedLimitations)
				suite.Require().NoError(err)

				err = ioutil.WriteFile(suite.tempFile, bz, os.ModePerm)
				suite.Require().NoError(err)
			}

			result, err := limitations.GetLimitationExpiration(uc.userID.String(), uc.command)
			if uc.expErr {
				suite.Require().Error(err)
				suite.Require().Nil(result)
			} else {
				suite.Require().NoError(err)

				if uc.expTime == nil {
					suite.Require().Nil(uc.expTime)
				} else {
					suite.Require().True(result.Equal(*uc.expTime))
				}
			}
		})
	}
}

func (suite *LimitationsTestSuite) TestRefreshLimitation() {
	usecases := []struct {
		name              string
		storedLimitations map[string]*limitations.UserLimitations
		userID            disgord.Snowflake
		command           string
		expiration        time.Time
		expErr            bool
		expLimitations    map[string]*limitations.UserLimitations
	}{
		{
			name:              "refreshing the limitation when limitations are empty",
			storedLimitations: nil,
			userID:            disgord.Snowflake(1),
			command:           "send",
			expiration:        time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
			expErr:            false,
			expLimitations: map[string]*limitations.UserLimitations{
				"1": {
					CommandsLimitations: map[string]time.Time{
						"send": time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
					},
				},
			},
		},
		{
			name: "refreshing the limitation when limitations are not empty but the account does not exist",
			storedLimitations: map[string]*limitations.UserLimitations{
				"2": {
					CommandsLimitations: map[string]time.Time{
						"send": time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
					},
				},
			},
			userID:     disgord.Snowflake(1),
			command:    "send",
			expiration: time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
			expErr:     false,
			expLimitations: map[string]*limitations.UserLimitations{
				"1": {
					CommandsLimitations: map[string]time.Time{
						"send": time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
					},
				},
				"2": {
					CommandsLimitations: map[string]time.Time{
						"send": time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
					},
				},
			},
		},
		{
			name: "refreshing the limitation when limitation for the account exists, but not the command",
			storedLimitations: map[string]*limitations.UserLimitations{
				"1": {
					CommandsLimitations: map[string]time.Time{
						"test": time.Date(2020, 1, 2, 00, 00, 00, 000, time.UTC),
					},
				},
			},
			userID:     disgord.Snowflake(1),
			command:    "send",
			expiration: time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
			expErr:     false,
			expLimitations: map[string]*limitations.UserLimitations{
				"1": {
					CommandsLimitations: map[string]time.Time{
						"test": time.Date(2020, 1, 2, 00, 00, 00, 000, time.UTC),
						"send": time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
					},
				},
			},
		},
		{
			name: "refreshing the limitation when limitation for the account and command exists",
			storedLimitations: map[string]*limitations.UserLimitations{
				"1": {
					CommandsLimitations: map[string]time.Time{
						"send": time.Date(2020, 1, 2, 00, 00, 00, 000, time.UTC),
					},
				},
			},
			userID:     disgord.Snowflake(1),
			command:    "send",
			expiration: time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
			expErr:     false,
			expLimitations: map[string]*limitations.UserLimitations{
				"1": {
					CommandsLimitations: map[string]time.Time{
						"send": time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
					},
				},
			},
		},
	}

	for _, uc := range usecases {
		uc := uc
		suite.Run(uc.name, func() {
			suite.SetupTest()
			limitations.SetLimitationsFile(suite.tempFile)

			if uc.storedLimitations != nil {
				bz, err := json.Marshal(uc.storedLimitations)
				suite.Require().NoError(err)
				err = ioutil.WriteFile(suite.tempFile, bz, os.ModePerm)
				suite.Require().NoError(err)
			}

			err := limitations.SetLimitationExpiration(uc.userID.String(), uc.command, uc.expiration)
			suite.Require().NoError(err)

			if uc.expErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)

				stored, err := limitations.ReadLimitations(suite.tempFile)
				suite.Require().NoError(err)

				suite.Len(stored, len(uc.expLimitations))
				for key, value := range uc.expLimitations {
					storedValue := stored[key]
					suite.True(storedValue.Equal(value))
				}
			}
		})
	}
}
