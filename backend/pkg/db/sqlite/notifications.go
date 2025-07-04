package db

import (
	"database/sql"
	"social-network/pkg/models"
)

type NotifRepository struct {
	DB *sql.DB
}

func (repo *NotifRepository) Save(notification models.Notification) error {
	stmt, err := repo.DB.Prepare("INSERT INTO notifications (notif_id, user_id,type,content,sender) values (?,?,?,?,?)")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(notification.ID, notification.TargetID, notification.Type, notification.Content, notification.Sender); err != nil {
		return err
	}
	return nil
}

func (repo *NotifRepository) Delete(notificationId string) error {
	_, err := repo.DB.Exec("DELETE FROM notifications WHERE notif_id=?", notificationId)
	if err != nil {
		return err
	}
	return nil
}

func (repo *NotifRepository) DeleteByType(notif models.Notification) error {
	_, err := repo.DB.Exec("DELETE FROM notifications WHERE user_id =? AND type =? AND content = ?", notif.TargetID, notif.Type, notif.Content)
	if err != nil {
		return err
	}
	return nil
}

// NOT TESTED
func (repo *NotifRepository) GetGroupRequests(groupId string) ([]models.Notification, error) {
	var notifications = []models.Notification{}
	rows, err := repo.DB.Query("SELECT content, notif_id, type, sender, user_id FROM notifications WHERE user_id = ? AND type = 'GROUP_REQUEST';", groupId)
	if err != nil {
		return notifications, err
	}
	for rows.Next() {
		var notif models.Notification
		rows.Scan(&notif.Content, &notif.ID, &notif.Type, &notif.Sender, &notif.TargetID)
		notifications = append(notifications, notif)
	}
	return notifications, nil
}

// GetGroupInvites retrieves all pending GROUP_INVITE notifications for a specific group.
func (repo *NotifRepository) GetGroupInvites(groupId string) ([]models.Notification, error) {
	var notifications = []models.Notification{}
	// For GROUP_INVITE, the `content` field of the notification stores the `groupId`.
	rows, err := repo.DB.Query("SELECT notif_id, user_id, type, content, sender FROM notifications WHERE content = ? AND type = 'GROUP_INVITE';", groupId)
	if err != nil {
		return notifications, err
	}
	defer rows.Close()

	for rows.Next() {
		var notif models.Notification
		// Scan into the fields of the Notification struct
		if err := rows.Scan(&notif.ID, &notif.TargetID, &notif.Type, &notif.Content, &notif.Sender); err != nil {
			// Log the error and continue if possible, or return
			continue
		}
		notifications = append(notifications, notif)
	}

	return notifications, nil
}

// DeleteGroupInvite removes a specific group invitation notification.
func (repo *NotifRepository) DeleteGroupInvite(targetId, groupId string) error {
	// A GROUP_INVITE is uniquely identified by the user it targets (`user_id`) 
	// and the group it's for (`content`).
	_, err := repo.DB.Exec("DELETE FROM notifications WHERE user_id = ? AND content = ? AND type = 'GROUP_INVITE'", targetId, groupId)
	if err != nil {
		return err
	}
	return nil
}

func (repo *NotifRepository) GetUserFromRequest(notificationId string) (string, error) {
	row := repo.DB.QueryRow("SELECT content FROM notifications WHERE notif_id = ? ", notificationId)
	var userId string
	if err := row.Scan(&userId); err != nil {
		return userId, err
	}
	return userId, nil
}

// NOT TESTED
func (repo *NotifRepository) CheckIfExists(notif models.Notification) (bool, error) {
	row := repo.DB.QueryRow("SELECT COUNT() FROM notifications WHERE user_id = ? AND content = ? AND type =? ", notif.TargetID, notif.Content, notif.Type)
	var resp int
	if err := row.Scan(&resp); err != nil {
		return false, err
	}
	if resp == 1 {
		return true, nil
	} else {
		return false, nil
	}
}

//NOT TESTED
func (repo *NotifRepository) GetGroupId(notificationId string) (string, error) {
	row := repo.DB.QueryRow("SELECT content FROM notifications WHERE notif_id = ? ", notificationId)
	var groupId string
	if err := row.Scan(&groupId); err != nil {
		return groupId, err
	}
	return groupId, nil
}

func (repo *NotifRepository) GetAll(userId string) ([]models.Notification, error) {
	var notifications = []models.Notification{}
	rows, err := repo.DB.Query("SELECT content, notif_id, type, sender, user_id, COALESCE(read, FALSE) as read FROM notifications WHERE user_id = ? OR (SELECT administrator FROM groups WHERE group_id = notifications.user_id) = ? ORDER BY notif_id DESC;", userId, userId)
	if err != nil {
		return notifications, err
	}
	for rows.Next() {
		var notif models.Notification
		rows.Scan(&notif.Content, &notif.ID, &notif.Type, &notif.Sender, &notif.TargetID, &notif.Read)
		notifications = append(notifications, notif)
	}
	return notifications, nil
}

func (repo *NotifRepository)GetCahtNotifById(notificationId string) (models.Notification, error){
row := repo.DB.QueryRow("SELECT content, user_id, sender FROM notifications WHERE notif_id = ?", notificationId)
	var notif models.Notification
	if err := row.Scan(&notif.Content, &notif.TargetID, &notif.Sender); err != nil {
		return notif, err
	}
	return notif, nil
}

func (repo *NotifRepository)CheckIfChatRequestExists(senderId, receiverId string)(bool, error){
		row := repo.DB.QueryRow("SELECT COUNT() FROM notifications WHERE user_id = ? AND sender = ? AND type = 'CHAT_REQUEST' ", receiverId, senderId)
	var resp int
	if err := row.Scan(&resp); err != nil {
		return false, err
	}
	if resp >= 1 {
		return true, nil
	} else {
		return false, nil
	}
}

func (repo *NotifRepository)GetContentFromChatRequest(senderId, receiverId string)(string, error){
	row := repo.DB.QueryRow("SELECT content FROM notifications WHERE user_id = ? AND sender = ? AND type = 'CHAT_REQUEST' ", receiverId, senderId)
	var resp string
	if err := row.Scan(&resp); err != nil {
		return resp, err
	}
	return resp, nil
}

// MarkAsRead marks a specific notification as read
func (repo *NotifRepository) MarkAsRead(notificationId string) error {
	_, err := repo.DB.Exec("UPDATE notifications SET read = TRUE WHERE notif_id = ?", notificationId)
	return err
}

// MarkAllAsRead marks all notifications for a user as read
func (repo *NotifRepository) MarkAllAsRead(userId string) error {
	_, err := repo.DB.Exec("UPDATE notifications SET read = TRUE WHERE user_id = ?", userId)
	return err
}