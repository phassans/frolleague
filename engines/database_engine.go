package engines

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/phassans/frolleague/common"
	"github.com/rs/zerolog"
)

type (
	databaseEngine struct {
		sql    *sql.DB
		logger zerolog.Logger
	}

	DatabaseEngine interface {
		// DBMethods
		SaveUser(LinkedInUserID, FirstName, LastName, LinkedInImage) error
		GetUserByID(UserID LinkedInUserID) (LinkedInUser, error)
		UpdateUserWithLinkedInURL(LinkedInUserID, LinkedInURL) error
		SaveToken(userID LinkedInUserID, accessToken AccessToken) error
		GetTokenByUserID(userID LinkedInUserID) (AccessToken, error)
		UpdateUserWithToken(userID LinkedInUserID, token AccessToken) error
		GetSchoolsByUserID(userID UserID) ([]School, error)
		GetCompaniesByUserID(userID UserID) ([]Company, error)

		// School Methods
		AddSchoolIfNotPresent(school SchoolName, degree Degree, fieldOfStudy FieldOfStudy) (SchoolID, error)
		DeleteSchool(school SchoolName, degree Degree, fieldOfStudy FieldOfStudy) error
		GetSchoolID(SchoolName, Degree, FieldOfStudy) (SchoolID, error)

		// Company Methods
		AddCompanyIfNotPresent(company CompanyName, location Location) (CompanyID, error)
		DeleteCompany(company CompanyName, location Location) error
		GetCompanyID(CompanyName, Location) (CompanyID, error)

		// UserToSchool
		AddUserToSchool(userID UserID, schoolID SchoolID, fromYear FromYear, toYear ToYear) error
		RemoveUserFromSchool(userID UserID, schoolID SchoolID) error
		UpdateUserStatusForAllSchools(userID UserID) error

		// UserToCompany
		AddUserToCompany(userID UserID, companyID CompanyID, title Title, fromYear FromYear, toYear ToYear) error
		RemoveUserFromCompany(userID UserID, companyID CompanyID) error
		UpdateUserStatusForAllCompanies(userID UserID) error

		// UserGroups
		AddGroupsToUser(userID UserID) ([]GroupInfo, error)
		GetGroupsByUserID(userID UserID) ([]GroupInfo, error)
		GetGroupsWithStatusByUserID(id UserID) ([]GroupWithStatus, error)
		ToggleUserGroup(userID UserID, group Group, status bool) error
	}
)

// NewDatabaseEngine returns an instance of userEngine
func NewDatabaseEngine(psql *sql.DB, logger zerolog.Logger) DatabaseEngine {
	return &databaseEngine{psql, logger}
}

func (l *databaseEngine) SaveUser(UserID LinkedInUserID, firstName FirstName, lastName LastName, linkedInImage LinkedInImage) error {
	user, err := l.GetUserByID(UserID)
	if err != nil {
		if _, ok := err.(common.ErrorUserNotExist); !ok {
			return err
		}
	}
	if user.UserID != "" {
		return common.DuplicateLinkedInUser{LinkedInUserID: string(UserID), Message: fmt.Sprintf("user with userID: %v already exists", UserID)}
	}

	return l.doSaveUser(UserID, firstName, lastName, linkedInImage)
}

func (l *databaseEngine) GetUserByID(UserID LinkedInUserID) (LinkedInUser, error) {
	var user LinkedInUser
	var sqlLinkedInURL sql.NullString
	rows := l.sql.QueryRow("SELECT user_id,url FROM linkedin_user WHERE user_id = $1", UserID)

	switch err := rows.Scan(&user.UserID, &sqlLinkedInURL); err {
	case sql.ErrNoRows:
		return LinkedInUser{}, common.ErrorUserNotExist{Message: fmt.Sprintf("user doesnt exist")}
	case nil:
		user.LinkedInURL = LinkedInURL(sqlLinkedInURL.String)
		return user, nil
	default:
		return LinkedInUser{}, common.DatabaseError{DBError: err.Error()}
	}
}

func (l *databaseEngine) doSaveUser(UserID LinkedInUserID, firstName FirstName, lastName LastName, linkedInImage LinkedInImage) error {
	_, err := l.sql.Exec("INSERT INTO linkedin_user(user_id, first_name, last_name, picture, insert_time) "+
		"VALUES($1,$2,$3,$4,$5)", UserID, firstName, lastName, linkedInImage, time.Now())
	if err != nil {
		return common.DatabaseError{DBError: err.Error()}
	}

	l.logger.Info().Msgf("successfully saved a linkedIn user with UserID: %s", UserID)
	return nil
}

func (l *databaseEngine) UpdateUserWithLinkedInURL(UserID LinkedInUserID, url LinkedInURL) error {
	updateWithURL := `UPDATE linkedin_user SET url = $1 WHERE user_id=$2;`

	_, err := l.sql.Exec(updateWithURL, url, UserID)
	if err != nil {
		return err
	}
	return nil
}

func (l *databaseEngine) SaveToken(userID LinkedInUserID, accessToken AccessToken) error {
	dbToken, err := l.GetTokenByUserID(userID)
	if err != nil {
		switch err.(type) {
		case common.ErrorNotExist:
			if err := l.doSaveToken(userID, accessToken); err != nil {
				return err
			}
		default:
			return err
		}
		return nil
	}

	if dbToken != "" {
		if err := l.UpdateUserWithToken(userID, accessToken); err != nil {
			return err
		}
		l.logger.Info().Msgf("saved token for userID: %s", userID)
	}

	return nil
}

func (l *databaseEngine) doSaveToken(userID LinkedInUserID, token AccessToken) error {
	_, err := l.sql.Exec("INSERT INTO linkedin_user_token(user_id, token, insert_time) "+
		"VALUES($1,$2,$3)", userID, token, time.Now())
	if err != nil {
		return common.DatabaseError{DBError: err.Error()}
	}

	l.logger.Info().Msgf("successfully saved token for user with UserID: %s", userID)
	return nil
}

func (l *databaseEngine) GetTokenByUserID(userID LinkedInUserID) (AccessToken, error) {
	var token AccessToken
	rows := l.sql.QueryRow("SELECT token FROM linkedin_user_token WHERE user_id = $1", userID)

	switch err := rows.Scan(&token); err {
	case sql.ErrNoRows:
		return AccessToken(""), common.ErrorNotExist{Message: fmt.Sprintf("user token doesnt exist")}
	case nil:
		return token, nil
	default:
		return AccessToken(""), common.DatabaseError{DBError: err.Error()}
	}
}

func (l *databaseEngine) UpdateUserWithToken(userID LinkedInUserID, token AccessToken) error {
	updateWithToken := `UPDATE linkedin_user_token SET token = $1 WHERE user_id=$2;`

	_, err := l.sql.Exec(updateWithToken, token, userID)
	if err != nil {
		return err
	}
	return nil
}

func (d *databaseEngine) AddSchoolIfNotPresent(schoolName SchoolName, degree Degree, fieldOfStudy FieldOfStudy) (SchoolID, error) {
	var schoolID SchoolID
	var err error

	// check if school is present
	schoolID, err = d.GetSchoolID(schoolName, degree, fieldOfStudy)
	if err != nil {
		return 0, common.DatabaseError{DBError: err.Error()}
	}

	if schoolID != 0 {
		d.logger.Info().Msgf("school:%s is already added with ID: %d", schoolName, schoolID)
		return schoolID, nil
	}

	// insert into school
	if err = d.sql.QueryRow("INSERT INTO school(school_name,degree,field_of_study,insert_time) "+
		"VALUES($1,$2,$3,$4) returning school_id;",
		schoolName, degree, fieldOfStudy, time.Now()).Scan(&schoolID); err != nil {
		return 0, common.DatabaseError{DBError: err.Error()}
	}

	d.logger.Info().Msgf("successfully added a school:%s with ID: %d", schoolName, schoolID)
	return schoolID, nil
}

func (d *databaseEngine) DeleteSchool(schoolName SchoolName, degree Degree, fieldOfStudy FieldOfStudy) error {
	_, err := d.sql.Exec("DELETE FROM school WHERE school_name=$1 AND degree = $2 AND field_of_study = $3", schoolName, degree, fieldOfStudy)
	if err != nil {
		return common.DatabaseError{DBError: err.Error()}
	}

	d.logger.Info().Msgf("successfully delete school: %s", schoolName)
	return nil
}

func (d *databaseEngine) GetSchoolID(schoolName SchoolName, degree Degree, fieldOfStudy FieldOfStudy) (SchoolID, error) {
	var id SchoolID
	rows := d.sql.QueryRow("SELECT school_id FROM school where school_name = $1  AND degree = $2 AND field_of_study = $3;", schoolName, degree, fieldOfStudy)

	err := rows.Scan(&id)

	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, common.DatabaseError{DBError: err.Error()}
	}

	return id, nil
}

func (d *databaseEngine) AddCompanyIfNotPresent(companyName CompanyName, location Location) (CompanyID, error) {
	var companyID CompanyID
	var err error

	// check if company is present
	companyID, err = d.GetCompanyID(companyName, location)
	if err != nil {
		return 0, common.DatabaseError{DBError: err.Error()}
	}

	if companyID != 0 {
		d.logger.Info().Msgf("company:%s is added with ID: %d", companyName, companyID)
		return companyID, nil
	}

	if err = d.sql.QueryRow("INSERT INTO company(company_name,location,insert_time) "+
		"VALUES($1,$2,$3) returning company_id;",
		companyName, location, time.Now()).Scan(&companyID); err != nil {
		return 0, common.DatabaseError{DBError: err.Error()}
	}

	d.logger.Info().Msgf("successfully added a company with ID: %d", companyID)
	return companyID, nil
}

func (d *databaseEngine) DeleteCompany(companyName CompanyName, location Location) error {
	_, err := d.sql.Exec("DELETE FROM company WHERE company_name=$1 AND location=$2", companyName, location)
	if err != nil {
		return common.DatabaseError{DBError: err.Error()}
	}

	d.logger.Info().Msgf("successfully delete company: %s", companyName)
	return nil
}

func (d *databaseEngine) GetCompanyID(companyName CompanyName, location Location) (CompanyID, error) {
	var id CompanyID
	rows := d.sql.QueryRow("SELECT company_id FROM company where company_name = $1  AND location = $2;", companyName, location)

	err := rows.Scan(&id)

	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, common.DatabaseError{DBError: err.Error()}
	}

	return id, nil
}

func (d *databaseEngine) AddUserToSchool(userID UserID, schoolID SchoolID, fromYear FromYear, toYear ToYear) error {
	schoolIDWithStatus, err := d.GetUserToSchoolByUserAndSchool(userID, schoolID, fromYear, toYear)
	if err != nil {
		return err
	}

	if schoolIDWithStatus.SchoolID == 0 {
		_, err := d.sql.Exec("INSERT INTO user_to_school(user_id,school_id,from_year,to_year,status,insert_time) "+
			"VALUES($1,$2,$3,$4,$5,$6)",
			userID, schoolID, fromYear, toYear, true, time.Now())
		if err != nil {
			return common.DatabaseError{DBError: err.Error()}
		}
		d.logger.Info().Msgf("successfully added a user: %s to school: %d", userID, schoolID)
	} else if schoolIDWithStatus.Status != true {
		if err := d.UpdateUserStatusForSchool(userID, schoolID, true); err != nil {
			return err
		}
	}

	return nil
}

func (d *databaseEngine) GetUserToSchoolByUserAndSchool(userID UserID, schoolID SchoolID, fromYear FromYear, toYear ToYear) (SchoolIDWithStatus, error) {
	var m SchoolIDWithStatus
	rows := d.sql.QueryRow("SELECT school_id,status FROM user_to_school WHERE user_id = $1 AND school_id = $2 AND from_year = $3 AND to_year = $4", userID, schoolID, fromYear, toYear)
	err := rows.Scan(&m.SchoolID, &m.Status)

	if err == sql.ErrNoRows {
		return SchoolIDWithStatus{}, nil
	} else if err != nil {
		return SchoolIDWithStatus{}, common.DatabaseError{DBError: err.Error()}
	}

	return m, nil
}

func (d *databaseEngine) RemoveUserFromSchool(userID UserID, schoolID SchoolID) error {
	_, err := d.sql.Exec("DELETE FROM user_to_school WHERE user_id=$1 AND school_id=$2", userID, schoolID)
	if err != nil {
		return common.DatabaseError{DBError: err.Error()}
	}

	d.logger.Info().Msgf("successfully removed user: %s from school: %d", userID, schoolID)
	return nil
}

func (d *databaseEngine) UpdateUserStatusForAllSchools(userID UserID) error {
	_, err := d.sql.Exec("UPDATE user_to_school SET status=$1 WHERE user_id=$2", false, userID)
	if err != nil {
		return common.DatabaseError{DBError: err.Error()}
	}

	d.logger.Info().Msgf("successfully removed user: %s from all schools", userID)
	return nil
}

func (d *databaseEngine) UpdateUserStatusForSchool(userID UserID, schoolID SchoolID, status bool) error {
	_, err := d.sql.Exec("UPDATE user_to_school SET status=$1 WHERE user_id=$2 AND school_id=$3", status, userID, schoolID)
	if err != nil {
		return common.DatabaseError{DBError: err.Error()}
	}

	d.logger.Info().Msgf("successfully removed user: %s from all schools", userID)
	return nil
}

func (d *databaseEngine) GetSchoolsByUserID(userID UserID) ([]School, error) {
	rows, err := d.sql.Query("SELECT school.school_name, school.degree, school.field_of_study,user_to_school.from_year, user_to_school.to_year "+
		"FROM school INNER JOIN user_to_school ON school.school_id = user_to_school.school_id "+
		"WHERE user_to_school.user_id=$1 AND user_to_school.status=$2", userID, true)
	if err != nil {
		return nil, common.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var schools []School
	for rows.Next() {
		var school School
		err = rows.Scan(&school.SchoolName, &school.Degree, &school.FieldOfStudy, &school.FromYear, &school.ToYear)
		if err != nil {
			return nil, common.DatabaseError{DBError: err.Error()}
		}
		schools = append(schools, school)
	}

	return schools, nil
}

func (d *databaseEngine) AddUserToCompany(userID UserID, companyID CompanyID, title Title, fromYear FromYear, toYear ToYear) error {
	companyIDWithStatus, err := d.GetUserToCompanyByUserAndCompany(userID, companyID)
	if err != nil {
		return err
	}

	if companyIDWithStatus.CompanyID == 0 {
		if _, err := d.sql.Exec("INSERT INTO user_to_company(user_id,company_id,title,from_year,to_year,status,insert_time) "+
			"VALUES($1,$2,$3,$4,$5,$6,$7)",
			userID, companyID, title, fromYear, toYear, true, time.Now()); err != nil {
			return common.DatabaseError{DBError: err.Error()}
		}
		d.logger.Info().Msgf("successfully added a user: %s to company: %d", userID, companyID)
	} else if companyIDWithStatus.Status != true {
		if err := d.UpdateUserStatusForCompany(userID, companyID, true); err != nil {
			return err
		}
	}

	return nil
}

func (d *databaseEngine) GetUserToCompanyByUserAndCompany(userID UserID, companyID CompanyID) (CompanyIDWithStatus, error) {
	var companyIDWithStatus CompanyIDWithStatus
	rows := d.sql.QueryRow("SELECT company_id,status FROM user_to_company "+
		"WHERE user_id = $1 AND company_id = $2",
		userID, companyID)
	err := rows.Scan(&companyIDWithStatus.CompanyID, &companyIDWithStatus.Status)

	if err == sql.ErrNoRows {
		return CompanyIDWithStatus{}, nil
	} else if err != nil {
		return CompanyIDWithStatus{}, common.DatabaseError{DBError: err.Error()}
	}

	return companyIDWithStatus, nil
}

func (d *databaseEngine) RemoveUserFromCompany(userID UserID, companyID CompanyID) error {
	_, err := d.sql.Exec("DELETE FROM user_to_company WHERE user_id=$1 AND company_id=$2", userID, companyID)
	if err != nil {
		return common.DatabaseError{DBError: err.Error()}
	}

	d.logger.Info().Msgf("successfully removed user: %s from company: %d", userID, companyID)
	return nil
}

func (d *databaseEngine) UpdateUserStatusForAllCompanies(userID UserID) error {
	_, err := d.sql.Exec("UPDATE user_to_company SET status=$1 WHERE user_id=$2", false, userID)
	if err != nil {
		return common.DatabaseError{DBError: err.Error()}
	}

	d.logger.Info().Msgf("successfully removed user: %s from all companies", userID)
	return nil
}

func (d *databaseEngine) UpdateUserStatusForCompany(userID UserID, companyID CompanyID, status bool) error {
	_, err := d.sql.Exec("UPDATE user_to_company SET status=$1 WHERE user_id=$2 AND company_id=$3", status, userID, companyID)
	if err != nil {
		return common.DatabaseError{DBError: err.Error()}
	}

	d.logger.Info().Msgf("successfully removed user: %s from all schools", userID)
	return nil
}

func (d *databaseEngine) GetCompaniesByUserID(userID UserID) ([]Company, error) {
	rows, err := d.sql.Query("SELECT company.company_name, company.location "+
		"FROM company INNER JOIN user_to_company ON company.company_id = user_to_company.company_id "+
		"WHERE user_to_company.user_id=$1 AND user_to_company.status=$2", userID, true)
	if err != nil {
		return nil, common.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var companies []Company
	for rows.Next() {
		var company Company
		err = rows.Scan(&company.CompanyName, &company.Location)
		if err != nil {
			return nil, common.DatabaseError{DBError: err.Error()}
		}
		companies = append(companies, company)
	}

	return companies, nil
}

func (d *databaseEngine) AddGroupsToUser(userID UserID) ([]GroupInfo, error) {
	groups, err := d.getGroupsBySchoolsAndCompanies(userID)
	if err != nil {
		return nil, err
	}
	userGroups, err := d.GetGroupsByUserID(userID)
	diffGroups := difference(groups, userGroups)

	// Note: only add unique groups
	grpMap := make(map[Group]bool)
	var uniqGroups []Group
	for _, group := range diffGroups {
		if !grpMap[group.Group] {
			// insert into school
			_, err = d.sql.Exec("INSERT INTO user_to_groups(user_id,group_name,status,group_source) VALUES($1,$2,$3,$4);", userID, group.Group, true, group.GroupSource)
			if err != nil {
				return nil, common.DatabaseError{DBError: err.Error()}
			}
			grpMap[group.Group] = true
			uniqGroups = append(uniqGroups, group.Group)

			d.logger.Info().Msgf("user with ID:%s joined group: %s", userID, group)
		}
	}

	userGroups = append(userGroups, diffGroups...)
	return userGroups, nil
}

func difference(a, b []GroupInfo) []GroupInfo {
	mb := map[Group]bool{}
	for _, x := range b {
		mb[x.Group] = true
	}
	var ab []GroupInfo
	for _, x := range a {
		if _, ok := mb[x.Group]; !ok {
			ab = append(ab, x)
		}
	}
	return ab
}

func (d *databaseEngine) ToggleUserGroup(userID UserID, group Group, status bool) error {
	updateUserGroups := `UPDATE user_to_groups SET status = $1 WHERE user_id=$2 AND group_name=$3;`
	_, err := d.sql.Exec(updateUserGroups, status, userID, group)
	if err != nil {
		return err
	}

	d.logger.Info().Msgf("user with ID:%s removed from group: %s", userID, group)
	return nil
}

func (d *databaseEngine) GetGroupsByUserID(userID UserID) ([]GroupInfo, error) {
	rows, err := d.sql.Query("SELECT group_name, group_source FROM user_to_groups "+
		"WHERE user_id=$1", userID)
	if err != nil {
		return nil, common.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var groups []GroupInfo

	for rows.Next() {
		var group GroupInfo
		err = rows.Scan(&group.Group, &group.GroupSource)
		if err != nil {
			return nil, common.DatabaseError{DBError: err.Error()}
		}
		groups = append(groups, group)

	}

	return groups, nil
}

func (d *databaseEngine) GetGroupsWithStatusByUserID(userID UserID) ([]GroupWithStatus, error) {
	rows, err := d.sql.Query("SELECT group_name, status, group_source FROM user_to_groups "+
		"WHERE user_id=$1", userID)
	if err != nil {
		return nil, common.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var groups []GroupWithStatus

	for rows.Next() {
		var groupWithStatus GroupWithStatus
		err = rows.Scan(&groupWithStatus.Group, &groupWithStatus.Status, &groupWithStatus.GroupSource)
		if err != nil {
			return nil, common.DatabaseError{DBError: err.Error()}
		}
		groups = append(groups, groupWithStatus)

	}

	return groups, nil
}

func (d *databaseEngine) getGroupsBySchoolsAndCompanies(userID UserID) ([]GroupInfo, error) {
	var groups []GroupInfo
	groupsSchools, err := d.GetSchoolsByUserID(userID)
	if err != nil {
		return nil, err
	}
	groups = append(groups, d.schoolsToGroups(groupsSchools)...)

	groupsCompanies, err := d.GetCompaniesByUserID(userID)
	if err != nil {
		return nil, err
	}
	groups = append(groups, d.companiesToGroups(groupsCompanies)...)
	return groups, nil
}

func (d *databaseEngine) schoolsToGroups(schools []School) []GroupInfo {
	var grps []GroupInfo

	for _, school := range schools {
		// add schoolName
		schoolName := strings.Replace(string(school.SchoolName), " ", "", -1)
		grps = append(grps, GroupInfo{Group(schoolName), "school"})

		if school.Degree != "" {
			degree := strings.Replace(string(school.Degree), " ", "", -1)
			fieldOfStudy := strings.Replace(string(school.FieldOfStudy), " ", "", -1)

			// add combination of school, degree & fieldOfStudy
			groupName := fmt.Sprintf("%s-%s-%s-%d-%d", schoolName, degree, fieldOfStudy, school.FromYear, school.ToYear)
			grps = append(grps, GroupInfo{Group(groupName), "school"})
		}

	}
	return grps
}

func (d *databaseEngine) companiesToGroups(companies []Company) []GroupInfo {
	var grps []GroupInfo

	for _, company := range companies {
		// add companyName
		companyName := strings.Replace(string(company.CompanyName), " ", "", -1)
		grps = append(grps, GroupInfo{Group(companyName), "company"})

		if company.Location != "" {
			location := strings.Replace(string(company.Location), " ", "", -1)
			// add combination of companyName & location
			groupName := fmt.Sprintf("%s-%s", companyName, location)
			grps = append(grps, GroupInfo{Group(groupName), "company"})
		}

	}
	return grps
}
