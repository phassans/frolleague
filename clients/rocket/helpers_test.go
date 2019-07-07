package rocket

import (
	"testing"

	"github.com/phassans/frolleague/common"
)

const (
	rocketURL         = "https://chat.districtchat.com"
	testAdminUserName = "pramod"
	testAdminPassword = "12345678"

	testName         = "banana"
	testUsername     = "iambanana"
	testUserEmail    = "banana@gmail.com"
	testUserPassword = "banana123"

	testGroupName = "testGroup"
)

var (
	rClient Client
)

func newRocketChatClient(t *testing.T) {
	common.InitLogger()
	rClient = NewRocketClient(rocketURL, common.GetLogger())

	if err := rClient.InitClient(testAdminUserName, testAdminPassword); err != nil {
		panic("cannot initialize test rocket client")
	}
}
