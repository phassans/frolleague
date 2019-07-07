package engines

type (
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

	ImageLink string

	UserID int64

	SchoolID int64

	CompanyID int64

	User struct {
		FirstName   FirstName
		LastName    LastName
		UserID      UserID
		Username    Username
		LinkedInURL LinkedInURL
		FileName    FileName
		ImageLink   ImageLink
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

	GroupWithStatus struct {
		Group  Group `json:"group"`
		Status bool  `json:"status"`
	}
)
