package db

import (
	"database/sql"
	"social-network/pkg/models"
)

type GroupRepository struct {
	DB *sql.DB
}

func (repo *GroupRepository) GetAllAndRelations(userID string) ([]models.Group, error) {
	var groups []models.Group
	rows, err := repo.DB.Query(`
		SELECT 
			group_id, 
			name, 
			privacy, 
			(SELECT COUNT(*) FROM group_users WHERE group_users.group_id = groups.group_id AND group_users.user_id = ?) as member, 
			administrator = ? as admin,
			(SELECT COUNT(*) FROM group_users WHERE group_users.group_id = groups.group_id) + 1 as member_count
		FROM groups;
	`, userID, userID)
	if err != nil {
		return groups, err
	}
	for rows.Next() {
		var group models.Group
		var member int
		var admin int
		rows.Scan(&group.ID, &group.Name, &group.Privacy, &member, &admin, &group.MemberCount)
		if member != 0 {
			group.Member = true
		}
		if admin != 0 {
			group.Administrator = true
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func (repo *GroupRepository) GetUserGroups(userID string) ([]models.Group, error) {
	var groups []models.Group
	rows, err := repo.DB.Query(`
		SELECT 
			group_id, 
			name, 
			privacy, 
			administrator = ? as admin,
			(SELECT COUNT(*) FROM group_users WHERE group_users.group_id = groups.group_id) + 1 as member_count
		FROM groups 
		WHERE (SELECT COUNT(*) FROM group_users WHERE group_users.group_id = groups.group_id AND group_users.user_id = ?) = 1 
		   OR administrator = ?;
	`, userID, userID, userID)
	if err != nil {
		return groups, err
	}
	for rows.Next() {
		var group models.Group
		var admin int
		rows.Scan(&group.ID, &group.Name, &group.Privacy, &admin, &group.MemberCount)
		if admin != 0 {
			group.Administrator = true
		} else {
			group.Member = true
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func (repo *GroupRepository) New(group models.Group) error {
	stmt, err := repo.DB.Prepare("INSERT INTO groups (group_id, name, description, administrator, privacy) values (?,?,?,?,?)")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(group.ID, group.Name, group.Description, group.AdminID, group.Privacy); err != nil {
		return err
	}
	return nil
}

func (repo *GroupRepository) GetData(groupId string) (models.Group, error) {
	row := repo.DB.QueryRow(`
		SELECT 
			name, 
			description, 
			administrator, 
			privacy,
			(SELECT COUNT(*) FROM group_users WHERE group_users.group_id = ?) + 1 as member_count
		FROM groups 
		WHERE group_id = ?
	`, groupId, groupId)
	var group models.Group
	if err := row.Scan(&group.Name, &group.Description, &group.AdminID, &group.Privacy, &group.MemberCount); err != nil {
		return group, err
	}
	group.ID = groupId
	return group, nil
}

func (repo *GroupRepository) GetMembers(groupId string) ([]models.User, error) {
	var members []models.User
	rows, err := repo.DB.Query(`
		SELECT 
			users.user_id, 
			users.first_name, 
			users.last_name, 
			users.nickname, 
			users.image,
			CASE WHEN users.user_id = groups.administrator THEN 1 ELSE 0 END as is_admin
		FROM users 
		LEFT JOIN groups ON groups.group_id = ?
		WHERE (users.user_id = groups.administrator) 
		   OR (users.user_id IN (SELECT user_id FROM group_users WHERE group_id = ?))
	`, groupId, groupId)
	if err != nil {
		return members, err
	}
	for rows.Next() {
		var member models.User
		var isAdmin int
		rows.Scan(&member.ID, &member.FirstName, &member.LastName, &member.Nickname, &member.ImagePath, &isAdmin)
		// Set admin status (we'll use a custom field for this)
		member.Following = (isAdmin == 1) // Reusing this field to indicate admin status for group context
		members = append(members, member)
	}
	return members, nil
}

func (repo *GroupRepository) IsMember(groupId, userId string) (bool, error) {
	row := repo.DB.QueryRow("SELECT COUNT() FROM group_users WHERE group_id = ? AND  user_id = ?", groupId, userId)
	var member int
	if err := row.Scan(&member); err != nil {
		return false, err
	}
	if member == 1 {
		return true, nil
	} else {
		return false, nil
	}
}

func (repo *GroupRepository) IsAdmin(groupId, userId string) (bool, error) {
	row := repo.DB.QueryRow("SELECT COUNT() FROM groups WHERE group_id = ? AND  administrator = ?", groupId, userId)
	var admin int
	if err := row.Scan(&admin); err != nil {
		return false, err
	}
	if admin == 1 {
		return true, nil
	} else {
		return false, nil
	}
}

func (repo *GroupRepository) GetAdmin(groupId string) (string, error) {
	row := repo.DB.QueryRow("SELECT administrator FROM groups WHERE group_id = ? ", groupId)
	var admin string
	if err := row.Scan(&admin); err != nil {
		return admin, err
	}
	return admin, nil
}

func (repo *GroupRepository) SaveMember(userId, groupId string) error {
	stmt, err := repo.DB.Prepare("INSERT INTO group_users (group_id, user_id) values (?,?)")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(groupId, userId); err != nil {
		return err
	}
	return nil
}
