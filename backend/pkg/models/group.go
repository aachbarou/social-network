package models

type Group struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	AdminID     string `json:"adminId"`
	ImagePath   string `json:"image"`
	Creator     *User  `json:"creator,omitempty"` // Creator information
	CreatedAt   string `json:"createdAt,omitempty"` // Creation date
	MemberCount int    `json:"memberCount,omitempty"` // Number of members

	Invitations []string `json:"invitations"`

	Member        bool `json:"member"` // true if current user is a member
	Administrator bool `json:"admin"`  // true if current user is admin
	RequestPending bool `json:"requestPending"` // true if request to join is pending
}

type GroupRepository interface {
	GetAllAndRelations(userId string) ([]Group, error)
	GetUserGroups(userId string) ([]Group, error)
	New(Group) error                               //create new group
	GetData(groupId string) (Group, error)         //get info- name and desc
	GetMembers(groupId string) ([]User, error)     // get all group members and admin
	GetAdmin(groupId string) (string, error)       //get admin id
	IsMember(groupId, userId string) (bool, error) //checks if user is a member
	IsAdmin(groupId, userId string) (bool, error)  //checks if user is admin
	DeleteGroup(groupId, adminId string) error     //delete group (admin only)
	UpdateGroup(group Group, adminId string) error //update group (admin only)

	SaveMember(userId, groupId string)error 
}
