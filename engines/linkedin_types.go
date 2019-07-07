package engines

type (
	LinkedInUserID string

	LinkedInURL string

	LinkedInImage string

	LinkedInUser struct {
		UserID        LinkedInUserID
		FirstName     FirstName
		LastName      LastName
		LinkedInURL   LinkedInURL
		LinkedInImage LinkedInImage
	}

	AccessToken string
)
