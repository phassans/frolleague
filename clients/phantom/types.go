package phantom

const (
	apiPath       = "api/v1/agent/34811/launch"
	sessionCookie = "AQEDAQTnvfEBUzSWAAABabh-Tu8AAAFrPgSLhVYAyGCswM-1OjPAjAbX4QNag2tPtNSplUEiWWiOBPpUSrFN_3yL4TG0aGMyWWMooBiUyYeCOGxvoG6PNETHmwJLxTFNznPAZiKgxk28ZAiQKjMMn0DW"
	output        = "result-object-with-output"
	userDataPath  = "user_data/"
)

type (
	Argument struct {
		SessionCookie string   `json:"sessionCookie"`
		ProfileUrls   []string `json:"profileUrls"`
		NoDatabase    bool     `json:"noDatabase"`
	}

	CrawlRequest struct {
		Output   string `json:"output"`
		Argument `json:"argument"`
	}

	Response struct {
		ContainerID   int `json:"containerId"`
		ExecutionTime int `json:"executionTime"`
		ExitCode      int `json:"exitCode"`
		ResultObject  []struct {
			General struct {
				ImgURL           string `json:"imgUrl"`
				FullName         string `json:"fullName"`
				Headline         string `json:"headline"`
				Company          string `json:"company"`
				School           string `json:"school"`
				Location         string `json:"location"`
				Connections      string `json:"connections"`
				ProfileURL       string `json:"profileUrl"`
				ConnectionDegree string `json:"connectionDegree"`
				Vmid             string `json:"vmid"`
				Description      string `json:"description"`
				FirstName        string `json:"firstName"`
				LastName         string `json:"lastName"`
			} `json:"general"`
			Jobs []struct {
				CompanyName string      `json:"companyName"`
				CompanyURL  string      `json:"companyUrl"`
				JobTitle    string      `json:"jobTitle"`
				DateRange   string      `json:"dateRange"`
				Location    string      `json:"location"`
				Description interface{} `json:"description"`
			} `json:"jobs"`
			Schools []struct {
				SchoolURL   string `json:"schoolUrl"`
				SchoolName  string `json:"schoolName"`
				Degree      string `json:"degree"`
				DegreeSpec  string `json:"degreeSpec"`
				DateRange   string `json:"dateRange"`
				Description string `json:"description,omitempty"`
			} `json:"schools"`
			Details struct {
				LinkedinProfile string `json:"linkedinProfile"`
				Twitter         string `json:"twitter"`
				Mail            string `json:"mail"`
			} `json:"details"`
			Skills []struct {
				Name         string `json:"name"`
				Endorsements string `json:"endorsements"`
			} `json:"skills"`
			AllSkills string `json:"allSkills"`
		} `json:"resultObject"`
		Output string `json:"output"`
	}

	CrawlResponse struct {
		Status  string   `json:"status"`
		Message string   `json:"message"`
		Data    Response `json:"data"`
	}

	FirstName string

	LastName string

	Username string

	Password string

	LinkedInURL string

	SchoolName string

	Degree string

	FieldOfStudy string

	FromYear int

	ToYear int

	CompanyName string

	Title string

	Location string

	Group string

	FileName string

	User struct {
		Firstname FirstName
		LastName  LastName
	}

	Company struct {
		CompanyName CompanyName
		FromYear    FromYear
		ToYear      ToYear
		Title       Title
		Location    Location
	}

	School struct {
		SchoolName   SchoolName
		Degree       Degree
		FieldOfStudy FieldOfStudy
		FromYear     FromYear
		ToYear       ToYear
	}

	Profile struct {
		User      User
		Companies []Company
		Schools   []School
		FileName  FileName
	}
)
