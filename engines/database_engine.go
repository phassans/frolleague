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
		// User Methods
		AddUser(username Username, password Password, linkedInURL LinkedInURL) (UserID, error)
		DeleteUser(username Username) error
		UpdateUserWithNameAndReference(name FirstName, lastName LastName, fileName FileName, id UserID) error
		UpdateUserWithImage(id UserID, imageName ImageLink) error
		GetUserByUserNameAndPassword(Username, Password) (User, error)
		GetUserByLinkedInURL(LinkedInURL) (User, error)
		GetUserByUserID(UserID) (User, error)
		UpdateUserPassword(UserID, Password) error

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

		// UserToCompany
		AddUserToCompany(userID UserID, companyID CompanyID, title Title, fromYear FromYear, toYear ToYear) error
		RemoveUserFromCompany(userID UserID, companyID CompanyID) error

		// UserGroups
		AddGroupsToUser(userID UserID) ([]Group, error)
		GetGroupsByUserID(userID UserID) ([]Group, error)
		GetGroupsWithStatusByUserID(id UserID) ([]GroupWithStatus, error)
		ToggleUserGroup(userID UserID, group Group, status bool) error
	}
)

// NewDatabaseEngine returns an instance of userEngine
func NewDatabaseEngine(psql *sql.DB, logger zerolog.Logger) DatabaseEngine {
	return &databaseEngine{psql, logger}
}

func (d *databaseEngine) AddUser(username Username, password Password, linkedInURL LinkedInURL) (UserID, error) {
	user, err := d.GetUserByLinkedInURL(linkedInURL)
	if err != nil {
		if _, ok := err.(common.ErrorUserNotExist); !ok {
			return 0, err
		}
	}

	if user.UserID != 0 {
		return 0, common.DuplicateSignUp{Username: string(user.Username), LinkedInURL: string(linkedInURL), Message: fmt.Sprintf("user with linkedingURL already exists")}
	}
	return d.doAddUser(username, password, linkedInURL)
}

func (d *databaseEngine) doAddUser(username Username, password Password, linkedInURL LinkedInURL) (UserID, error) {
	var userID UserID
	err := d.sql.QueryRow("INSERT INTO viraagh_user(username,password,linkedIn_URL,insert_time) "+
		"VALUES($1,$2,$3,$4) returning user_id;", username, password, linkedInURL, time.Now()).Scan(&userID)
	if err != nil {
		return 0, common.DatabaseError{DBError: err.Error()}
	}

	d.logger.Info().Msgf("successfully added a user with ID: %d", userID)

	return userID, nil
}

func (d *databaseEngine) DeleteUser(username Username) error {
	_, err := d.sql.Exec("DELETE FROM viraagh_user WHERE username=$1", username)
	if err != nil {
		return common.DatabaseError{DBError: err.Error()}
	}

	d.logger.Info().Msgf("successfully delete user: %s", username)
	return nil
}

func (d *databaseEngine) UpdateUserWithNameAndReference(firstName FirstName, lastName LastName, fileName FileName, id UserID) error {
	updateUserWithNameAndReferenceSQL := `UPDATE viraagh_user SET first_name = $1, last_name = $2, filename = $3 WHERE user_id=$4;`

	_, err := d.sql.Exec(updateUserWithNameAndReferenceSQL, firstName, lastName, fileName, id)
	if err != nil {
		return err
	}
	return nil
}

func (d *databaseEngine) UpdateUserWithImage(id UserID, imageLink ImageLink) error {
	updateUserWithImageSQL := `UPDATE viraagh_user SET image_link = $1 WHERE user_id=$2;`
	_, err := d.sql.Exec(updateUserWithImageSQL, imageLink, id)
	if err != nil {
		return err
	}
	return nil
}

func (d *databaseEngine) UpdateUserPassword(id UserID, password Password) error {
	updateUserPassword := `UPDATE viraagh_user SET password = $1 WHERE user_id=$2;`

	_, err := d.sql.Exec(updateUserPassword, password, id)
	if err != nil {
		return err
	}
	return nil
}

func (d *databaseEngine) GetUserByUserNameAndPassword(userName Username, password Password) (User, error) {
	var user User
	rows := d.sql.QueryRow("SELECT user_id, first_name, last_name, username, linkedin_url, filename FROM viraagh_user WHERE username = $1 AND password = $2", userName, password)
	err := rows.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Username, &user.LinkedInURL, &user.FileName)

	if err == sql.ErrNoRows {
		return User{}, common.UserError{Message: fmt.Sprintf("user doesnt exist")}
	} else if err != nil {
		return User{}, common.DatabaseError{DBError: err.Error()}
	}

	return user, nil
}

func (d *databaseEngine) GetUserByLinkedInURL(linkedInURL LinkedInURL) (User, error) {
	var user User
	rows := d.sql.QueryRow("SELECT user_id, username, linkedin_url FROM viraagh_user WHERE linkedin_url = $1", linkedInURL)

	switch err := rows.Scan(&user.UserID, &user.Username, &user.LinkedInURL); err {
	case sql.ErrNoRows:
		return User{}, common.ErrorUserNotExist{Message: fmt.Sprintf("user doesnt exist")}
	case nil:
		return user, nil
	default:
		return User{}, common.DatabaseError{DBError: err.Error()}
	}
}

func (d *databaseEngine) GetUserByUserID(userID UserID) (User, error) {
	var user User
	rows := d.sql.QueryRow("SELECT user_id, username, linkedin_url FROM viraagh_user WHERE user_id = $1", userID)

	switch err := rows.Scan(&user.UserID, &user.Username, &user.LinkedInURL); err {
	case sql.ErrNoRows:
		return User{}, common.ErrorUserNotExist{Message: fmt.Sprintf("user doesnt exist")}
	case nil:
		return user, nil
	default:
		return User{}, common.DatabaseError{DBError: err.Error()}
	}
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
	count, err := d.GetUserToSchoolByUserAndSchool(userID, schoolID, fromYear, toYear)
	if err != nil {
		return err
	}

	if count == 0 {
		_, err := d.sql.Exec("INSERT INTO user_to_school(user_id,school_id,from_year,to_year,insert_time) "+
			"VALUES($1,$2,$3,$4,$5)",
			userID, schoolID, fromYear, toYear, time.Now())
		if err != nil {
			return common.DatabaseError{DBError: err.Error()}
		}
		d.logger.Info().Msgf("successfully added a user: %d to school: %d", userID, schoolID)
	}

	return nil
}

func (d *databaseEngine) GetUserToSchoolByUserAndSchool(userID UserID, schoolID SchoolID, fromYear FromYear, toYear ToYear) (int, error) {
	var count int
	rows := d.sql.QueryRow("SELECT count(*) FROM user_to_school WHERE user_id = $1 AND school_id = $2 AND from_year = $3 AND to_year = $4", userID, schoolID, fromYear, toYear)
	err := rows.Scan(&count)

	if err == sql.ErrNoRows {
		return 0, common.UserError{Message: fmt.Sprintf("user doesnt exist")}
	} else if err != nil {
		return 0, common.DatabaseError{DBError: err.Error()}
	}

	return count, nil
}

func (d *databaseEngine) RemoveUserFromSchool(userID UserID, schoolID SchoolID) error {
	_, err := d.sql.Exec("DELETE FROM user_to_school WHERE user_id=$1 AND school_id=$2", userID, schoolID)
	if err != nil {
		return common.DatabaseError{DBError: err.Error()}
	}

	d.logger.Info().Msgf("successfully removed user: %d from school: %d", userID, schoolID)
	return nil
}

func (d *databaseEngine) GetSchoolsByUserID(userID UserID) ([]School, error) {
	rows, err := d.sql.Query("SELECT school.school_name, school.degree, school.field_of_study,user_to_school.from_year, user_to_school.to_year "+
		"FROM school INNER JOIN user_to_school ON school.school_id = user_to_school.school_id "+
		"WHERE user_to_school.user_id=$1", userID)
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
	count, err := d.GetUserToCompanyByUserAndCompany(userID, companyID, title, fromYear, toYear)
	if err != nil {
		return err
	}

	if count == 0 {
		if _, err := d.sql.Exec("INSERT INTO user_to_company(user_id,company_id,title,from_year,to_year,insert_time) "+
			"VALUES($1,$2,$3,$4,$5,$6)",
			userID, companyID, title, fromYear, toYear, time.Now()); err != nil {
			return common.DatabaseError{DBError: err.Error()}
		}
		d.logger.Info().Msgf("successfully added a user: %d to company: %d", userID, companyID)
	}

	return nil
}

func (d *databaseEngine) GetUserToCompanyByUserAndCompany(userID UserID, companyID CompanyID, title Title, fromYear FromYear, toYear ToYear) (int, error) {
	var count int
	rows := d.sql.QueryRow("SELECT count(*) FROM user_to_company WHERE user_id = $1 AND company_id = $2 AND title = $3 AND from_year = $4 AND to_year = $5", userID, companyID, title, fromYear, toYear)
	err := rows.Scan(&count)

	if err == sql.ErrNoRows {
		return 0, common.UserError{Message: fmt.Sprintf("user doesnt exist")}
	} else if err != nil {
		return 0, common.DatabaseError{DBError: err.Error()}
	}

	return count, nil
}

func (d *databaseEngine) RemoveUserFromCompany(userID UserID, companyID CompanyID) error {
	_, err := d.sql.Exec("DELETE FROM user_to_company WHERE user_id=$1 AND company_id=$2", userID, companyID)
	if err != nil {
		return common.DatabaseError{DBError: err.Error()}
	}

	d.logger.Info().Msgf("successfully removed user: %d from company: %d", userID, companyID)
	return nil
}

func (d *databaseEngine) GetCompaniesByUserID(userID UserID) ([]Company, error) {
	rows, err := d.sql.Query("SELECT company.company_name, company.location "+
		"FROM company INNER JOIN user_to_company ON company.company_id = user_to_company.company_id "+
		"WHERE user_to_company.user_id=$1", userID)
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

func (d *databaseEngine) AddGroupsToUser(userID UserID) ([]Group, error) {
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
		if !grpMap[group] {
			// insert into school
			_, err = d.sql.Exec("INSERT INTO user_to_groups(user_id,group_name,status) VALUES($1,$2,$3);", userID, group, true)
			if err != nil {
				return nil, common.DatabaseError{DBError: err.Error()}
			}
			grpMap[group] = true
			uniqGroups = append(uniqGroups, group)

			d.logger.Info().Msgf("user with ID:%d joined group: %s", userID, group)
		}
	}

	userGroups = append(userGroups, diffGroups...)
	return userGroups, nil
}

func difference(a, b []Group) []Group {
	mb := map[Group]bool{}
	for _, x := range b {
		mb[x] = true
	}
	ab := []Group{}
	for _, x := range a {
		if _, ok := mb[x]; !ok {
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

	d.logger.Info().Msgf("user with ID:%d removed from group: %s", userID, group)
	return nil
}

func (d *databaseEngine) GetGroupsByUserID(userID UserID) ([]Group, error) {
	rows, err := d.sql.Query("SELECT group_name FROM user_to_groups "+
		"WHERE user_id=$1 AND status=$2", userID, true)
	if err != nil {
		return nil, common.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var groups []Group

	for rows.Next() {
		var group Group
		err = rows.Scan(&group)
		if err != nil {
			return nil, common.DatabaseError{DBError: err.Error()}
		}
		groups = append(groups, group)

	}

	return groups, nil
}

func (d *databaseEngine) GetGroupsWithStatusByUserID(userID UserID) ([]GroupWithStatus, error) {
	rows, err := d.sql.Query("SELECT group_name, status FROM user_to_groups "+
		"WHERE user_id=$1", userID)
	if err != nil {
		return nil, common.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var groups []GroupWithStatus

	for rows.Next() {
		var groupWithStatus GroupWithStatus
		err = rows.Scan(&groupWithStatus.Group, &groupWithStatus.Status)
		if err != nil {
			return nil, common.DatabaseError{DBError: err.Error()}
		}
		groups = append(groups, groupWithStatus)

	}

	return groups, nil
}

func (d *databaseEngine) getGroupsBySchoolsAndCompanies(userID UserID) ([]Group, error) {
	var groups []Group
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

func (d *databaseEngine) schoolsToGroups(schools []School) []Group {
	var grps []Group

	for _, school := range schools {
		// add schoolName
		schoolName := strings.Replace(string(school.SchoolName), " ", "", -1)
		grps = append(grps, Group(schoolName))

		degree := strings.Replace(string(school.Degree), " ", "", -1)
		fieldOfStudy := strings.Replace(string(school.FieldOfStudy), " ", "", -1)

		// add combination of school, degree & fieldOfStudy
		groupName := fmt.Sprintf("%s-%s-%s-%d-%d", schoolName, degree, fieldOfStudy, school.FromYear, school.ToYear)
		grps = append(grps, Group(groupName))

	}
	return grps
}

func (d *databaseEngine) companiesToGroups(companies []Company) []Group {
	var grps []Group

	for _, company := range companies {
		// add companyName
		companyName := strings.Replace(string(company.CompanyName), " ", "", -1)
		grps = append(grps, Group(companyName))

		location := strings.Replace(string(company.Location), " ", "", -1)
		// add combination of companyName & location
		groupName := fmt.Sprintf("%s-%s", companyName, location)
		grps = append(grps, Group(groupName))

	}
	return grps
}
