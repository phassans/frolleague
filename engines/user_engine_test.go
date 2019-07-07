package engines

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/phassans/frolleague/clients/phantom"
	"github.com/phassans/frolleague/common"
	"github.com/phassans/frolleague/db"
	"github.com/rs/zerolog"
)

const (
	phantomURL = "https://phantombuster.com"
)

var (
	testRoach  db.Roach
	testLogger zerolog.Logger
	dbEngine   DatabaseEngine
	uEngine    UserEngine
)

func newUserEngine(t *testing.T) {
	testLogger = common.GetLogger()
	if testLogger.Log() == nil {
		{
			common.InitLogger()
		}
	}

	var err error
	if testRoach.Db == nil {
		testRoach, err = db.New(db.Config{Host: testDatabaseHost, Port: testDataPort, User: testDatabaseUsername, Password: testDatabasePassword, Database: testDatabase})
		if err != nil {
			logger.Fatal().Msgf("could not connect to db. errpr %s", err)
		}
	}
	dbEngine = NewDatabaseEngine(testRoach.Db, testLogger)
	rclient := phantom.NewPhantomClient(phantomURL, testLogger)

	uEngine, err = NewUserEngine(nil, rclient, dbEngine, testLogger)
	require.NoError(t, err)
}

func TestUserEngine_SignUp_Login(t *testing.T) {
	newUserEngine(t)
	{
		err := uEngine.SignUp(testUserName, testPassword, testLinkedInURL)
		require.NoError(t, err)

		user, err := uEngine.Login(testUserName, testPassword)
		require.NoError(t, err)
		require.NotEqual(t, UserID(0), user.UserID)
	}
}
