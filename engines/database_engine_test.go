package engines

import (
	"testing"

	"github.com/phassans/frolleague/common"
	"github.com/phassans/frolleague/db"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

var (
	roach  db.Roach
	logger zerolog.Logger
	engine DatabaseEngine
)

func newDataBaseEngine(t *testing.T) {
	logger = common.GetLogger()
	if logger.Log() == nil {
		{
			common.InitLogger()
		}
	}

	var err error
	if roach.Db == nil {
		roach, err = db.New(db.Config{Host: testDatabaseHost, Port: testDataPort, User: testDatabaseUsername, Password: testDatabasePassword, Database: testDatabase})
		if err != nil {
			logger.Fatal().Msgf("could not connect to db. errpr %s", err)
		}
	}

	engine = NewDatabaseEngine(roach.Db, logger)
	require.NotEmpty(t, engine)
}

func TestDatabaseEngine_AddDeleteSchool(t *testing.T) {
	newDataBaseEngine(t)
	{
		schoolID, err := engine.AddSchoolIfNotPresent(testSchool, testDegree, testFieldOfStudy)
		require.NoError(t, err)
		require.NotEmpty(t, schoolID)

		schoolID, err = engine.AddSchoolIfNotPresent(testSchool, testDegree, testFieldOfStudy)
		require.NoError(t, err)
		require.NotEmpty(t, schoolID)

		err = engine.DeleteSchool(testSchool, testDegree, testFieldOfStudy)
		require.NoError(t, err)
	}
}

func TestDatabaseEngine_AddDeleteCompany(t *testing.T) {
	newDataBaseEngine(t)
	{
		companyID, err := engine.AddCompanyIfNotPresent(testCompany, testLocation)
		require.NoError(t, err)
		require.NotEmpty(t, companyID)

		err = engine.DeleteCompany(testCompany, testLocation)
		require.NoError(t, err)
	}
}

func TestDatabaseEngine_AddRemoveUserSchool(t *testing.T) {
	newDataBaseEngine(t)
	{
		schoolID, err := engine.AddSchoolIfNotPresent(testSchool, testDegree, testFieldOfStudy)
		require.NoError(t, err)
		require.NotEmpty(t, schoolID)

		err = engine.AddUserToSchool(testUserID, schoolID, testFromYear, testToYear)
		require.NoError(t, err)

		err = engine.RemoveUserFromSchool(testUserID, schoolID)
		require.NoError(t, err)

		err = engine.DeleteSchool(testSchool, testDegree, testFieldOfStudy)
		require.NoError(t, err)
	}
}

func TestDatabaseEngine_AddRemoveUserCompany(t *testing.T) {
	newDataBaseEngine(t)
	{
		companyID, err := engine.AddCompanyIfNotPresent(testCompany, testLocation)
		require.NoError(t, err)
		require.NotEmpty(t, companyID)

		err = engine.AddUserToCompany(testUserID, companyID, testTitle, testFromYear, testToYear)
		require.NoError(t, err)

		err = engine.RemoveUserFromCompany(testUserID, companyID)
		require.NoError(t, err)

		err = engine.DeleteCompany(testCompany, testLocation)
		require.NoError(t, err)
	}
}

func TestDatabaseEngine_AddGroupsToUser(t *testing.T) {
	newDataBaseEngine(t)
	{
		companyID, err := engine.AddCompanyIfNotPresent(testCompany, testLocation)
		require.NoError(t, err)
		require.NotEmpty(t, companyID)

		err = engine.AddUserToCompany(testUserID, companyID, testTitle, testFromYear, testToYear)
		require.NoError(t, err)

		schoolID, err := engine.AddSchoolIfNotPresent(testSchool, testDegree, testFieldOfStudy)
		require.NoError(t, err)
		require.NotEmpty(t, schoolID)

		err = engine.AddUserToSchool(testUserID, schoolID, testFromYear, testToYear)
		require.NoError(t, err)

		groups, err := engine.AddGroupsToUser(testUserID)
		require.NoError(t, err)
		require.Equal(t, 4, len(groups))

		err = engine.RemoveUserFromSchool(testUserID, schoolID)
		require.NoError(t, err)

		err = engine.DeleteSchool(testSchool, testDegree, testFieldOfStudy)
		require.NoError(t, err)

		err = engine.RemoveUserFromCompany(testUserID, companyID)
		require.NoError(t, err)

		err = engine.DeleteCompany(testCompany, testLocation)
		require.NoError(t, err)
	}
}

func TestDatabaseEngine_Token(t *testing.T) {
	newDataBaseEngine(t)
	{
		err := engine.SaveToken(LinkedInUserID(testUserID), AccessToken("12345"))
		require.NoError(t, err)

		token, err := engine.GetTokenByUserID(LinkedInUserID(testUserID))
		require.NoError(t, err)
		require.Equal(t, AccessToken("12345"), token)

		err = engine.UpdateUserWithToken(LinkedInUserID(testUserID), AccessToken("78910"))
		require.NoError(t, err)

		token, err = engine.GetTokenByUserID(LinkedInUserID(testUserID))
		require.NoError(t, err)
		require.Equal(t, AccessToken("78910"), token)
	}
}

func TestDatabaseEngine_LinkedInUser(t *testing.T) {
	newDataBaseEngine(t)
	{
		err := engine.SaveUser(LinkedInUserID(testUserID+"1"), testUserFirstName, testUserLastName, LinkedInImage("foobar"))
		//require.NoError(t, err)

		user, err := engine.GetUserByID(LinkedInUserID(testUserID + "1"))
		require.NoError(t, err)
		require.Equal(t, user.UserID, LinkedInUserID(testUserID+"1"))

		err = engine.UpdateUserWithLinkedInURL(LinkedInUserID(testUserID+"1"), LinkedInURL("https://foobar.com"))
		require.NoError(t, err)

		user, err = engine.GetUserByID(LinkedInUserID(testUserID + "1"))
		require.NoError(t, err)
		require.Equal(t, LinkedInUserID(testUserID+"1"), user.UserID)
		require.Equal(t, LinkedInURL("https://foobar.com"), user.LinkedInURL)
	}
}
