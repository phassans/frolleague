package phantom

import (
	"testing"

	"github.com/phassans/frolleague/common"
	"github.com/stretchr/testify/require"
)

const (
	phantomURL  = "https://phantombuster.com"
	linkedInURL = "https://www.linkedin.com/in/pramod-shashidhara-21568923/"
)

var (
	pClient Client
)

func newPhantomClient(t *testing.T) {
	common.InitLogger()
	pClient = NewPhantomClient(phantomURL, common.GetLogger())
}

func TestClient_CrawlUrl(t *testing.T) {
	newPhantomClient(t)
	{
		resp, err := pClient.CrawlUrl(linkedInURL, false)
		require.NoError(t, err)
		require.NotNil(t, resp)
	}
}

func TestClient_GetUserFromResponse(t *testing.T) {
	newPhantomClient(t)
	{
		ifFile := false
		crawlResponse, err := pClient.CrawlUrl(linkedInURL, ifFile)
		require.NoError(t, err)
		user := pClient.GetUserFromResponse(crawlResponse)
		require.NoError(t, err)
		require.Equal(t, FirstName("Pramod"), user.Firstname)
		require.Equal(t, LastName("Shashidhara"), user.LastName)
	}
}

func TestClient_GetSchoolsFromResponse(t *testing.T) {
	newPhantomClient(t)
	{
		crawlResponse, err := pClient.CrawlUrl(jsonFile, true)
		require.NoError(t, err)
		schools, err := pClient.GetSchoolsFromResponse(crawlResponse)
		require.NoError(t, err)
		require.Equal(t, 2, len(schools))
	}
}

func TestClient_GetCompaniesFromResponse(t *testing.T) {
	newPhantomClient(t)
	{
		crawlResponse, err := pClient.CrawlUrl(jsonFile, true)
		require.NoError(t, err)
		companies, err := pClient.GetCompaniesFromResponse(crawlResponse)
		require.NoError(t, err)
		require.Equal(t, 12, len(companies))
	}
}

func TestClient_GetUserProfile(t *testing.T) {
	newPhantomClient(t)
	{
		profile, err := pClient.GetUserProfile(linkedInURL, true)
		require.NoError(t, err)
		require.NotNil(t, profile)
	}
}

func TestClient_SaveResponse(t *testing.T) {
	newPhantomClient(t)
	{
		crawlResponse, err := pClient.CrawlUrl(jsonFile, true)
		require.NoError(t, err)
		filename, err := pClient.SaveUserProfile(crawlResponse)
		require.NoError(t, err)
		require.NotNil(t, filename)
	}
}
