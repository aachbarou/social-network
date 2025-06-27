package db

import (
	"database/sql"
	"fmt"
	"social-network/pkg/models"
)

type GroupRepository struct {
	DB *sql.DB
}

func (repo *GroupRepository) GetAllAndRelations(userID string) ([]models.Group, error) {
	var groups []models.Group
	rows, err := repo.DB.Query(`
		SELECT 
			g.group_id, 
			g.name, 
			g.description, 
			IFNULL(g.image, '') as image,
			g.administrator,
			IFNULL(u.nickname, u.first_name || ' ' || u.last_name) as creator_name,
			IFNULL(u.image, '') as creator_image,
			IFNULL(g.created_at, '') as created_at,
			(SELECT COUNT(*) FROM group_users WHERE group_users.group_id = g.group_id AND group_users.user_id = ?) as member,
			g.administrator = ? as admin,
			(SELECT COUNT(*) FROM group_users WHERE group_users.group_id = g.group_id) + 1 as member_count
		FROM groups g
		LEFT JOIN users u ON g.administrator = u.user_id
	`, userID, userID)
	if err != nil {
		return groups, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var group models.Group
		var member int
		var admin int
		var creatorName, creatorImage string
		
		err := rows.Scan(&group.ID, &group.Name, &group.Description, &group.ImagePath, 
			&group.AdminID, &creatorName, &creatorImage, &group.CreatedAt, &member, &admin, &group.MemberCount)
		if err != nil {
			continue
		}
		
		// Set creator information
		group.Creator = &models.User{
			ID:        group.AdminID,
			Nickname:  creatorName,
			ImagePath: creatorImage,
		}
		
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
			g.group_id, 
			g.name, 
			g.description, 
			IFNULL(g.image, '') as image,
			g.administrator,
			IFNULL(u.nickname, u.first_name || ' ' || u.last_name) as creator_name,
			IFNULL(u.image, '') as creator_image,
			IFNULL(g.created_at, '') as created_at,
			g.administrator = ? as admin,
			(SELECT COUNT(*) FROM group_users WHERE group_users.group_id = g.group_id) + 1 as member_count
		FROM groups g
		LEFT JOIN users u ON g.administrator = u.user_id
		WHERE (SELECT COUNT(*) FROM group_users WHERE group_users.group_id = g.group_id AND group_users.user_id = ?) = 1 
		   OR g.administrator = ?
	`, userID, userID, userID)
	if err != nil {
		return groups, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var group models.Group
		var admin int
		var creatorName, creatorImage string
		
		err := rows.Scan(&group.ID, &group.Name, &group.Description, &group.ImagePath, 
			&group.AdminID, &creatorName, &creatorImage, &group.CreatedAt, &admin, &group.MemberCount)
		if err != nil {
			fmt.Printf("Error scanning group row: %v\n", err)
			continue
		}
		
		fmt.Printf("Scanned group: ID=%s, Name=%s, Description=%s, ImagePath=%s, Admin=%d\n", 
			group.ID, group.Name, group.Description, group.ImagePath, admin)
		
		// Set creator information
		group.Creator = &models.User{
			ID:        group.AdminID,
			Nickname:  creatorName,
			ImagePath: creatorImage,
		}
		
		if admin != 0 {
			group.Administrator = true
		} else {
			group.Member = true
		}
		groups = append(groups, group)
	}
	fmt.Printf("Returning %d user groups\n", len(groups))
	return groups, nil
}

func (repo *GroupRepository) New(group models.Group) error {
	stmt, err := repo.DB.Prepare("INSERT INTO groups (group_id, name, description, administrator, image, created_at) values (?,?,?,?,?,CURRENT_TIMESTAMP)")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(group.ID, group.Name, group.Description, group.AdminID, group.ImagePath); err != nil {
		return err
	}
	return nil
}

func (repo *GroupRepository) GetData(groupId string) (models.Group, error) {
	row := repo.DB.QueryRow(`
		SELECT 
			g.name, 
			g.description, 
			g.administrator, 
			IFNULL(g.image, '') as image,
			IFNULL(u.nickname, u.first_name || ' ' || u.last_name) as creator_name,
			IFNULL(u.image, '') as creator_image,
			IFNULL(g.created_at, '') as created_at,
			(SELECT COUNT(*) FROM group_users WHERE group_users.group_id = g.group_id) + 1 as member_count
		FROM groups g
		LEFT JOIN users u ON g.administrator = u.user_id
		WHERE g.group_id = ?
	`, groupId)
	
	var group models.Group
	var creatorName, creatorImage string
	
	if err := row.Scan(&group.Name, &group.Description, &group.AdminID, &group.ImagePath, 
		&creatorName, &creatorImage, &group.CreatedAt, &group.MemberCount); err != nil {
		return group, err
	}
	
	group.ID = groupId
	
	// Set creator information
	group.Creator = &models.User{
		ID:        group.AdminID,
		Nickname:  creatorName,
		ImagePath: creatorImage,
	}
	
	return group, nil
}

func (repo *GroupRepository) GetMembers(groupId string) ([]models.User, error) {
	var members []models.User
	rows, err := repo.DB.Query("SELECT user_id, IFNULL(nickname, first_name || ' ' || last_name), image FROM users WHERE (user_id = (SELECT administrator FROM groups WHERE group_id =?)) OR ((SELECT COUNT() FROM group_users WHERE group_id = ? AND  user_id = users.user_id )=1) ", groupId, groupId)
	if err != nil {
		return members, err
	}
	for rows.Next() {
		var member models.User
		rows.Scan(&member.ID, &member.Nickname, &member.ImagePath)
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

func (repo *GroupRepository) DeleteGroup(groupId, adminId string) error {
	// First verify the user is the admin
	var count int
	err := repo.DB.QueryRow("SELECT COUNT(*) FROM groups WHERE group_id = ? AND administrator = ?", groupId, adminId).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("user is not authorized to delete this group")
	}
	
	// Begin transaction
	tx, err := repo.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// Delete all group members
	if _, err := tx.Exec("DELETE FROM group_users WHERE group_id = ?", groupId); err != nil {
		return err
	}
	
	// Delete group posts if they exist (you may need to add this table)
	// _, err = tx.Exec("DELETE FROM group_posts WHERE group_id = ?", groupId)
	
	// Delete group events if they exist (you may need to add this table)
	// _, err = tx.Exec("DELETE FROM group_events WHERE group_id = ?", groupId)
	
	// Delete notifications related to this group
	// _, err = tx.Exec("DELETE FROM notifications WHERE target_id = ?", groupId)
	
	// Finally delete the group itself
	if _, err := tx.Exec("DELETE FROM groups WHERE group_id = ?", groupId); err != nil {
		return err
	}
	
	return tx.Commit()
}

func (repo *GroupRepository) UpdateGroup(group models.Group, adminId string) error {
	// First verify the user is the admin
	var count int
	err := repo.DB.QueryRow("SELECT COUNT(*) FROM groups WHERE group_id = ? AND administrator = ?", group.ID, adminId).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("user is not authorized to update this group")
	}
	
	// Update the group
	stmt, err := repo.DB.Prepare("UPDATE groups SET name = ?, description = ?, image = ? WHERE group_id = ? AND administrator = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	
	_, err = stmt.Exec(group.Name, group.Description, group.ImagePath, group.ID, adminId)
	return err
}
