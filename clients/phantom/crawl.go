package phantom

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"

	"github.com/nu7hatch/gouuid"
)

const (
	jsonFile = "response.json"
)

func (c *client) CrawlUrl(linkedInURL string, file bool) (CrawlResponse, error) {
	arg := Argument{
		SessionCookie: sessionCookie,
		ProfileUrls:   []string{linkedInURL},
		NoDatabase:    true,
	}

	if file {
		return c.doCrawlUrlFile()
	}

	return c.doCrawlUrl(
		CrawlRequest{Output: output,
			Argument: arg,
		},
	)
}

func (c *client) doCrawlUrlFile() (CrawlResponse, error) {
	time.Sleep(20 * time.Second)

	fileName := userDataPath + jsonFile
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		fileName = jsonFile
	}

	fileBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Errorf("error reading file: %s", fileName)
		return CrawlResponse{}, nil
	}

	if len(fileBytes) == 0 {
		fmt.Errorf("error reading file: %s and length: %d", fileName, len(fileBytes))
	}

	var resp Response
	err = json.Unmarshal(fileBytes, &resp)
	if err != nil {
		return CrawlResponse{}, nil
	}
	return CrawlResponse{Data: resp}, nil
}

func (c *client) doCrawlUrl(request CrawlRequest) (CrawlResponse, error) {
	logger := c.logger

	response, err := c.DoPost(request)
	if err != nil {
		logger.Error().Msgf("CrawlUrl returned with error")
		return CrawlResponse{}, fmt.Errorf("CrawlUrl returned with error: %s", err)
	}

	// read response.json
	var resp CrawlResponse
	err = json.Unmarshal(response, &resp)
	if err != nil {
		logger = logger.With().Str("error", err.Error()).Logger()
		logger.Error().Msgf("unmarshal error on CrawlResponse")
		return CrawlResponse{}, err
	}

	return resp, nil
}

func (c *client) GetUserProfile(linkedInURL string, isFile bool) (Profile, error) {
	resp, err := c.CrawlUrl(string(linkedInURL), isFile)
	if err != nil {
		return Profile{}, err
	}
	if resp.Status == "error" {
		return Profile{}, fmt.Errorf("%s", resp.Message)
	}

	fileName, err := c.SaveUserProfile(resp)
	if err != nil {
		return Profile{}, err
	}

	schools, err := c.GetSchoolsFromResponse(resp)
	if err != nil {
		return Profile{}, err
	}
	companies, err := c.GetCompaniesFromResponse(resp)
	if err != nil {
		return Profile{}, err
	}

	user := c.GetUserFromResponse(resp)
	return Profile{user, companies, schools, fileName}, nil
}

func (c *client) GetUserFromResponse(resp CrawlResponse) User {
	var u User
	for _, obj := range resp.Data.ResultObject {
		u = User{FirstName(obj.General.FirstName), LastName(obj.General.LastName)}
	}
	return u
}

func (c *client) GetSchoolsFromResponse(resp CrawlResponse) ([]School, error) {
	var schools []School
	for _, obj := range resp.Data.ResultObject {
		for _, school := range obj.Schools {
			from, to, err := parseDateRangeForSchool(school.DateRange)
			if err != nil {
				return nil, err
			}
			s := School{SchoolName(removeSpecialChars(school.SchoolName)), Degree(removeSpecialChars(school.Degree)), FieldOfStudy(removeSpecialChars(school.DegreeSpec)), from, to}
			schools = append(schools, s)
		}
	}
	return schools, nil
}

func (c *client) GetCompaniesFromResponse(resp CrawlResponse) ([]Company, error) {
	var companies []Company
	for _, obj := range resp.Data.ResultObject {
		for _, jobs := range obj.Jobs {
			from, to, err := parseDateRangeForCompanies(jobs.DateRange)
			if err != nil {
				return nil, err
			}

			c := Company{CompanyName(removeSpecialChars(jobs.CompanyName)), from, to, Title(removeSpecialChars(jobs.JobTitle)), Location(removeSpecialChars(jobs.Location))}
			companies = append(companies, c)
		}
	}
	return companies, nil
}

func (c *client) SaveUserProfile(resp CrawlResponse) (FileName, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	// get user name
	user := c.GetUserFromResponse(resp)
	if user.Firstname == "" {
		return "", fmt.Errorf("cannot save file. firstName is empty")
	}
	fileName := fmt.Sprintf("%s.%s.%s.json", user.Firstname, user.LastName, u)

	// marshall the resp
	b, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	// write to file
	err = ioutil.WriteFile(userDataPath+fileName, b, 0644)
	if err != nil {
		return "", err
	}

	return FileName(fileName), nil
}

func removeSpecialChars(str string) string {
	// Make a Regex to say we only want letters and numbers
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return ""
	}
	return reg.ReplaceAllString(str, "")
}
