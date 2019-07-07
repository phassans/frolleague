package engines

type (
	LinkedInId                string
	LinkedInFirstName         string
	LinkedInLastName          string
	LinkedInProfilePicture    string
	LinkedInAccessToken       string
	LinkedInAuthorizationCode string

	linkedInUser struct {
		LinkedInID             LinkedInId             `json:"linkedInID"`
		LinkedInFirstName      LinkedInFirstName      `json:"linkedInFirstName"`
		LinkedInLastName       LinkedInLastName       `json:"linkedInLastName"`
		LinkedInProfilePicture LinkedInProfilePicture `json:"linkedInProfilePicture,omitempty"`
	}
)
