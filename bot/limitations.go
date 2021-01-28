package bot

import (
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/rs/zerolog/log"
	"github.com/tendermint/tendermint/libs/json"

	"github.com/desmos-labs/discord-bot/keys"
)

var (
	LimitationsFile = path.Join(keys.DataDir, "limitations.json")
)

// UserLimitations contains the data about the limitations of a single user
type UserLimitations struct {
	CommandsLimitations map[string]time.Time `json:"commands_limitations"` // Map of limitations for each command
}

func NewUserLimitations() *UserLimitations {
	return &UserLimitations{
		CommandsLimitations: map[string]time.Time{},
	}
}

func (u *UserLimitations) Equal(v *UserLimitations) bool {
	if len(u.CommandsLimitations) != len(v.CommandsLimitations) {
		return false
	}

	for key, value := range u.CommandsLimitations {
		vValue := v.CommandsLimitations[key]
		if !value.Equal(vValue) {
			return false
		}
	}

	return true
}

// ReadLimitations reads the user limitations contained inside the given file
func ReadLimitations(file string) (map[string]*UserLimitations, error) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		_, err := os.Create(file)
		if err != nil {
			return nil, err
		}
	}

	bz, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// Empty file, return nil
	if len(bz) == 0 {
		return map[string]*UserLimitations{}, nil
	}

	// Read all the limitations
	var limitations map[string]*UserLimitations
	return limitations, json.Unmarshal(bz, &limitations)
}

func GetLimitationExpiration(userID disgord.Snowflake, command string) (*time.Time, error) {
	limitations, err := ReadLimitations(LimitationsFile)
	if err != nil {
		return nil, err
	}

	userLimit, found := limitations[userID.String()]
	if !found {
		log.Debug().Str(keys.LogCommand, command).Str(keys.LogUser, userID.String()).Msg("has no limitations set")
		return nil, nil
	}

	// Get the limitation expiration for the specific command
	timeLimit, ok := userLimit.CommandsLimitations[command]
	if !ok {
		log.Debug().Str(keys.LogCommand, command).Str(keys.LogUser, userID.String()).Msg("no limitations for the command found")
		return nil, nil
	}
	return &timeLimit, nil
}

func SetLimitationExpiration(userID disgord.Snowflake, command string, expiration time.Time) error {
	usersLimitations, err := ReadLimitations(LimitationsFile)
	if err != nil {
		return err
	}

	// Get the limitations for the user
	userLimits, ok := usersLimitations[userID.String()]
	if !ok {
		userLimits = NewUserLimitations()
	}

	// Update the limitation
	userLimits.CommandsLimitations[command] = expiration
	usersLimitations[userID.String()] = userLimits

	// Serialize the data
	bz, err := json.Marshal(&usersLimitations)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(LimitationsFile, bz, os.ModePerm)
}
