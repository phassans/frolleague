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

func TestDatabaseEngine_AddDeleteUser(t *testing.T) {
	newDataBaseEngine(t)
	{
		userID, err := engine.AddUser(testUserName, testPassword, testLinkedInURL)
		require.NoError(t, err)
		require.NotEmpty(t, userID)

		err = engine.UpdateUserWithNameAndReference(testUserFirstName, testUserLastName, testFileName, userID)
		require.NoError(t, err)

		user, err := engine.GetUserByUserNameAndPassword(testUserName, testPassword)
		require.NoError(t, err)
		require.NotEqual(t, UserID(0), user.UserID)

		user, err = engine.GetUserByUserNameAndPassword(testUserNameInvalid, testPassword)
		require.Error(t, err)
		require.Equal(t, UserID(0), user.UserID)

		user, err = engine.GetUserByLinkedInURL(testLinkedInURL)
		require.NoError(t, err)
		require.NotEqual(t, UserID(0), user.UserID)

		err = engine.DeleteUser(testUserName)
		require.NoError(t, err)
	}
}
func TestDatabaseEngine_AddUserDuplicates(t *testing.T) {
	newDataBaseEngine(t)
	{
		userID, err := engine.AddUser(testUserName, testPassword, testLinkedInURL)
		require.NoError(t, err)
		require.NotEmpty(t, userID)

		userID, err = engine.AddUser(testUserName, testPassword, testLinkedInURL)
		require.Error(t, err)
		require.Equal(t, UserID(0), userID)

		err = engine.DeleteUser(testUserName)
		require.NoError(t, err)
	}
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
		userID, err := engine.AddUser(testUserName, testPassword, testLinkedInURL)
		require.NoError(t, err)
		require.NotEmpty(t, userID)

		schoolID, err := engine.AddSchoolIfNotPresent(testSchool, testDegree, testFieldOfStudy)
		require.NoError(t, err)
		require.NotEmpty(t, schoolID)

		err = engine.AddUserToSchool(userID, schoolID, testFromYear, testToYear)
		require.NoError(t, err)

		err = engine.RemoveUserFromSchool(userID, schoolID)
		require.NoError(t, err)

		err = engine.DeleteSchool(testSchool, testDegree, testFieldOfStudy)
		require.NoError(t, err)

		err = engine.DeleteUser(testUserName)
		require.NoError(t, err)
	}
}

func TestDatabaseEngine_AddRemoveUserCompany(t *testing.T) {
	newDataBaseEngine(t)
	{
		userID, err := engine.AddUser(testUserName, testPassword, testLinkedInURL)
		require.NoError(t, err)
		require.NotEmpty(t, userID)

		companyID, err := engine.AddCompanyIfNotPresent(testCompany, testLocation)
		require.NoError(t, err)
		require.NotEmpty(t, companyID)

		err = engine.AddUserToCompany(userID, companyID, testTitle, testFromYear, testToYear)
		require.NoError(t, err)

		err = engine.RemoveUserFromCompany(userID, companyID)
		require.NoError(t, err)

		err = engine.DeleteCompany(testCompany, testLocation)
		require.NoError(t, err)

		err = engine.DeleteUser(testUserName)
		require.NoError(t, err)
	}
}

func TestDatabaseEngine_AddGroupsToUser(t *testing.T) {
	newDataBaseEngine(t)
	{
		userID, err := engine.AddUser(testUserName, testPassword, testLinkedInURL)
		require.NoError(t, err)
		require.NotEmpty(t, userID)

		companyID, err := engine.AddCompanyIfNotPresent(testCompany, testLocation)
		require.NoError(t, err)
		require.NotEmpty(t, companyID)

		err = engine.AddUserToCompany(userID, companyID, testTitle, testFromYear, testToYear)
		require.NoError(t, err)

		schoolID, err := engine.AddSchoolIfNotPresent(testSchool, testDegree, testFieldOfStudy)
		require.NoError(t, err)
		require.NotEmpty(t, schoolID)

		err = engine.AddUserToSchool(userID, schoolID, testFromYear, testToYear)
		require.NoError(t, err)

		groups, err := engine.AddGroupsToUser(userID)
		require.NoError(t, err)
		require.Equal(t, 4, len(groups))

		err = engine.RemoveUserFromSchool(userID, schoolID)
		require.NoError(t, err)

		err = engine.DeleteSchool(testSchool, testDegree, testFieldOfStudy)
		require.NoError(t, err)

		err = engine.RemoveUserFromCompany(userID, companyID)
		require.NoError(t, err)

		err = engine.DeleteCompany(testCompany, testLocation)
		require.NoError(t, err)

		err = engine.DeleteUser(testUserName)
		require.NoError(t, err)
	}
}
