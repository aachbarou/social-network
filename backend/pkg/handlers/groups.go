package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"social-network/pkg/models"
	"social-network/pkg/utils"
	ws "social-network/pkg/wsServer"
	"strings"
)

/* -------------------------------------------------------------------------- */
/*                                  get data                                  */
/* -------------------------------------------------------------------------- */

func (handler *Handler) AllGroups(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access user id
	userId := r.Context().Value(utils.UserKey).(string)
	// request all groups + relations
	groups, errGroups := handler.repos.GroupRepo.GetAllAndRelations(userId)
	if errGroups != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	utils.RespondWithGroups(w, groups, 200)
}

// returns all groups that current user is a member of or admin
func (handler *Handler) UserGroups(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access user id
	userId := r.Context().Value(utils.UserKey).(string)
	// request user Groups
	groups, errGroups := handler.repos.GroupRepo.GetUserGroups(userId)
	if errGroups != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	utils.RespondWithGroups(w, groups, 200)
}

// returns all groups that specified user is a member of or admin
func (handler *Handler) OtherUserGroups(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access user id
	query := r.URL.Query()
	userId := query.Get("userId")
	if userId == "" { //check if user id provided in request
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	// request user Groups
	groups, errGroups := handler.repos.GroupRepo.GetUserGroups(userId)
	if errGroups != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	utils.RespondWithGroups(w, groups, 200)
}

// returns info about group - > name, description, id and administrator id
// also includes group status for current user -> admin / member or pending member request
func (handler *Handler) GroupInfo(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// get group id from request
	query := r.URL.Query()
	groupId := query.Get("groupId")
	if groupId == "" { //check if group id exists in request
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	group, err := handler.repos.GroupRepo.GetData(groupId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	/* --------------- // add additional data about member status --------------- */
	// check if admin, or member or member request pending
	// access current user id
	userId := r.Context().Value(utils.UserKey).(string)
	if userId == group.AdminID {
		group.Administrator = true
	} else {
		group.Member, err = handler.repos.GroupRepo.IsMember(group.ID, userId)
		if err != nil {
			utils.RespondWithError(w, "Error on getting data", 200)
			return
		}
		if !group.Member {
			notification := models.Notification{
				TargetID: group.ID,
				Type:     "GROUP_REQUEST",
				Content:  userId,
			}
			group.RequestPending, err = handler.repos.NotifRepo.CheckIfExists(notification)
			if err != nil {
				utils.RespondWithError(w, "Error on getting data", 200)
				return
			}
		}

	}
	utils.RespondWithGroups(w, []models.Group{group}, 200)
}

// returns list of all group members and administrator
func (handler *Handler) GroupMembers(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// get group id from request
	query := r.URL.Query()
	groupId := query.Get("groupId")
	if groupId == "" { //check if group id exists in request
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	members, err := handler.repos.GroupRepo.GetMembers(groupId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	utils.RespondWithUsers(w, members, 200)
}

// returns list of events for group
// GET request with group_id included in query
// in case if user is not a member of the group, returns error
func (handler *Handler) GroupEvents(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access current user id
	userId := r.Context().Value(utils.UserKey).(string)
	// get group id from request
	query := r.URL.Query()
	groupId := query.Get("groupId")
	if groupId == "" { //check if group id exists in request
		utils.RespondWithError(w, "Error on reading group ID", 200)
		return
	}
	// check if current user is a member or admin  of the group
	var isMember = false
	isAdmin, err := handler.repos.GroupRepo.IsAdmin(groupId, userId)
	if err != nil {
		utils.RespondWithError(w, "Error on reading role", 200)
		return
	}
	if !isAdmin {
		isMember, err = handler.repos.GroupRepo.IsMember(groupId, userId)
		if err != nil {
			utils.RespondWithError(w, "Error on checking if is group member", 200)
			return
		}
	}
	if !isMember && !isAdmin {
		utils.RespondWithError(w, "Not a member", 200)
		return
	}
	// current user is a member or admin -> get events
	events, err := handler.repos.EventRepo.GetAll(groupId)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(w, "Error on getting event data", 200)
		return
	}
	/* ----------------------- attach author to each event ---------------------- */
	for i := 0; i < len(events); i++ {
		events[i].Author, _ = handler.repos.UserRepo.GetDataMin(events[i].AuthorID)
	}
	/* -------------------- attach participation to each event ------------------- */
	for i := 0; i < len(events); i++ {
		going, _ := handler.repos.EventRepo.IsParticipating(events[i].ID, userId)
		if going {
			events[i].Going = "YES"
		} else {
			events[i].Going = "NO"
		}
	}
	utils.RespondWithEvents(w, events, 200)
}

// returns all posts that belongs to group
func (handler *Handler) GroupPosts(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access current user id
	userId := r.Context().Value(utils.UserKey).(string)
	// get group id from request
	query := r.URL.Query()
	groupId := query.Get("groupId")
	if groupId == "" { //check if group id exists in request
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	/* -------- check if current user is a member or admin  of the group -------- */
	isAdmin, err := handler.repos.GroupRepo.IsAdmin(groupId, userId)
	isMember, err := handler.repos.GroupRepo.IsMember(groupId, userId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	if !isMember && !isAdmin {
		utils.RespondWithError(w, "Not a member", 200)
		return
	}
	/* ------------- current user is a member or admin -> get posts ------------- */
	posts, err := handler.repos.PostRepo.GetGroupPosts(groupId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	// Get post author info attached
	if err = AttachAuthors(handler, &posts); err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	// Get comment info for each post
	if err = AttachComments(handler, &posts); err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	utils.RespondWithPosts(w, posts, 200)
}

// returns pending requests to join to group, only for admin
// for others respond with error
func (handler *Handler) GroupRequests(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access current user id
	userId := r.Context().Value(utils.UserKey).(string)
	// get group id from request
	query := r.URL.Query()
	groupId := query.Get("groupId")
	if groupId == "" { //check if group id exists in request
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	/* ------------------------- check if user is admin ------------------------- */
	isAdmin, err := handler.repos.GroupRepo.IsAdmin(groupId, userId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	if !isAdmin {
		utils.RespondWithError(w, "Unauthorized", 200)
		return
	}
	/* ---------------------- get pending requests form db ---------------------- */
	notifications, err := handler.repos.NotifRepo.GetGroupRequests(groupId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	for i := 0; i < len(notifications); i++ {
		notifications[i].User, err = handler.repos.UserRepo.GetDataMin(notifications[i].Content)
		if err != nil {
			utils.RespondWithError(w, "Error on getting data", 200)
			return
		}
	}
	utils.RespondWithNotifications(w, notifications, 200)
}

func (handler *Handler) CancelGroupRequests(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access current user id
	currentUserId := r.Context().Value(utils.UserKey).(string)
	// get group id from request
	query := r.URL.Query()
	groupId := query.Get("groupId")
	if groupId == "" { //check if group id exists in request
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	// delete notification
	notif := models.Notification{
		Type:     "GROUP_REQUEST",
		TargetID: groupId,
		Content:  currentUserId,
	}
	if err := handler.repos.NotifRepo.DeleteByType(notif); err != nil {
		utils.RespondWithError(w, "Error on canceling request", 200)
		return
	}
	utils.RespondWithSuccess(w, "gROUP request canceled successfuly", 200)
}

/* -------------------------------------------------------------------------- */
/*                                save new data                               */
/* -------------------------------------------------------------------------- */

func (handler *Handler) NewGroup(wsServer *ws.Server, w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	if r.Method != "POST" {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	
	var newGroup models.Group
	
	// Check if this is a multipart form (with image) or JSON
	contentType := r.Header.Get("Content-Type")
	fmt.Printf("Content-Type: %s\n", contentType)
	
	if strings.Contains(contentType, "multipart/form-data") {
		// Handle multipart form data (with image)
		fmt.Println("Processing multipart form data...")
		err := r.ParseMultipartForm(10 << 20) // 10MB max
		if err != nil {
			fmt.Printf("Error parsing multipart form: %v\n", err)
			utils.RespondWithError(w, "Error parsing form data", 200)
			return
		}
		
		newGroup.Name = r.FormValue("name")
		newGroup.Description = r.FormValue("description")
		fmt.Printf("Form values - Name: %s, Description: %s\n", newGroup.Name, newGroup.Description)
		
		// Save image if provided
		newGroup.ImagePath = utils.SaveImage(r)
		fmt.Printf("Saved image path: %s\n", newGroup.ImagePath)
		
	} else {
		// Handle JSON data (no image)
		fmt.Println("Processing JSON data...")
		err := json.NewDecoder(r.Body).Decode(&newGroup)
		if err != nil {
			fmt.Printf("Error decoding JSON: %v\n", err)
			utils.RespondWithError(w, "Error on form submittion", 200)
			return
		}
		fmt.Printf("JSON data - Name: %s, Description: %s\n", newGroup.Name, newGroup.Description)
	}
	
	/* ------------------------- attach additional data ------------------------ */
	// access user id
	newGroup.AdminID = r.Context().Value(utils.UserKey).(string)
	newGroup.ID = utils.UniqueId()
	
	fmt.Printf("Final group data - ID: %s, Name: %s, Description: %s, ImagePath: %s, AdminID: %s\n", 
		newGroup.ID, newGroup.Name, newGroup.Description, newGroup.ImagePath, newGroup.AdminID)

	/* ------------------------------- save in db ------------------------------- */
	err := handler.repos.GroupRepo.New(newGroup)
	if err != nil {
		utils.RespondWithError(w, "Error on saving group", 200)
		return
	}
	/* ------------------------- save invitations in db ------------------------- */
	for i := 0; i < len(newGroup.Invitations); i++ {
		// save each follower in db
		newNotif := models.Notification{
			ID:       utils.UniqueId(),
			TargetID: newGroup.Invitations[i],
			Type:     "GROUP_INVITE",
			Content:  newGroup.ID,
			Sender:   newGroup.AdminID,
		}
		err = handler.repos.NotifRepo.Save(newNotif)
		if err != nil {
			utils.RespondWithError(w, "Internal server error", 200)
			return
		}

		for client := range wsServer.Clients {
			if client.ID == newNotif.TargetID {
				client.SendNotification(newNotif)
			}
		}

		// save as a new member of group
		// if err = handler.repos.GroupRepo.SaveMember(newGroup.Invitations[i], newGroup.ID); err != nil {
		// 	utils.RespondWithError(w, "Internal server error", 200)
		// 	return
		// }
	}
	newGroup.Administrator = true
	//NOTIFY WEBSOCKET ABOUT NEW NOTIFICATION
	fmt.Printf("Returning group data: %+v\n", newGroup)
	utils.RespondWithGroups(w, []models.Group{newGroup}, 200)
}

// NOT TESTED
func (handler *Handler) NewGroupPost(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	if r.Method != "POST" {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	/* ---------------------------- read incoming data --------------------------- */
	err := r.ParseMultipartForm(3145728) // 3MB
	if err != nil {
		utils.RespondWithError(w, "Error in form validation", 200)
		return
	}
	userId := r.Context().Value(utils.UserKey).(string)
	/* ------------------------ create new post instance ------------------------ */
	newPost := models.Post{
		ID:       utils.UniqueId(),
		Content:  r.PostFormValue("body"),
		GroupID:  r.PostFormValue("groupId"),
		AuthorID: userId,
	}
	/* -------- check if current user is a member or admin  of the group -------- */
	isAdmin, err := handler.repos.GroupRepo.IsAdmin(newPost.GroupID, userId)
	isMember, err := handler.repos.GroupRepo.IsMember(newPost.GroupID, userId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	if !isMember && !isAdmin {
		utils.RespondWithError(w, "Not a member", 200)
		return
	}
	/* ------------------------ save image in filesystem ------------------------ */
	newPost.ImagePath = utils.SaveImage(r)
	/* -------------------------- save post in database ------------------------- */
	errDB := handler.repos.PostRepo.New(newPost)
	if errDB != nil {
		utils.RespondWithError(w, "Error on saving post", 200)
		return
	}
	utils.RespondWithSuccess(w, "New post created", 200)
}

// handle when new user wants to join the group
func (handler *Handler) NewGroupRequest(wsServer *ws.Server, w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access current user id
	userId := r.Context().Value(utils.UserKey).(string)
	// get group id from request
	query := r.URL.Query()
	groupId := query.Get("groupId")
	if groupId == "" { //check if group id exists in request
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	/* --------------------------- check if not admin --------------------------- */
	isAdmin, err := handler.repos.GroupRepo.IsAdmin(groupId, userId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	if isAdmin {
		utils.RespondWithError(w, "Invalid request", 200)
		return
	}
	/* -------------------- create new notification instance -------------------- */
	notification := models.Notification{
		ID:       utils.UniqueId(),
		TargetID: groupId,
		Type:     "GROUP_REQUEST",
		Content:  userId,
		Sender:   userId,
	}
	/* --------------------- check if request already exists -------------------- */
	exists, err := handler.repos.NotifRepo.CheckIfExists(notification)
	if err != nil {
		utils.RespondWithError(w, "Error on saving request", 200)
		return
	}
	if exists {
		utils.RespondWithSuccess(w, "Request already saved", 200)
		return
	}
	/* ------------------------- save notification in db ------------------------ */
	err = handler.repos.NotifRepo.Save(notification)
	if err != nil {
		utils.RespondWithError(w, "Error on saving request", 200)
	}

	//SEND MESSAGE TO TO GROUP ADMIN IF IS ONLINE
	admin, err := handler.repos.GroupRepo.GetAdmin(groupId)
	if err != nil {
		utils.RespondWithError(w, "Error on finding group admin", 200)
		return
	}
	for client := range wsServer.Clients {
		if client.ID == admin {
			client.SendNotification(notification)
		}
	}
	utils.RespondWithSuccess(w, "Request saved successfuly", 200)
}

// NOT TESTED
// handle response from group administrator for requests to join group
// waits for requestId and response -accept/decline
func (handler *Handler) ResponseGroupRequest(wsServer *ws.Server, w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	if r.Method != "POST" {
		utils.RespondWithError(w, "Error on data submittion", 200)
		return
	}
	/* --------------------------- read incoming data --------------------------- */
	type Response struct {
		GroupID   string `json:"groupId"`
		RequestID string `json:"requestId"`
		Response  string `json:"response"`
	}
	var response Response
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	/* ---------------------- check if all fields provided ---------------------- */
	if response.Response != "accept" && response.Response != "decline" {
		utils.RespondWithError(w, "Response unvalid", 200)
		return
	}
	if response.GroupID == "" || response.RequestID == "" {
		utils.RespondWithError(w, "Response incomplete", 200)
		return
	}
	/* ------------------- chack if curren user is group admin ------------------ */
	// access user id
	userId := r.Context().Value(utils.UserKey).(string)
	isAdmin, err := handler.repos.GroupRepo.IsAdmin(response.GroupID, userId)
	if !isAdmin {
		utils.RespondWithError(w, "Unauthorized", 200)
		return
	}
	/* ----------------------------- handle response ---------------------------- */
	// if accepted -> save as new member
	if response.Response == "accept" {
		// get id of member that requests to join
		joinerId, err := handler.repos.NotifRepo.GetUserFromRequest(response.RequestID)
		if err != nil {
			utils.RespondWithError(w, "Internal server error", 200)
			return
		}
		// save as a new member of group
		if err = handler.repos.GroupRepo.SaveMember(joinerId, response.GroupID); err != nil {
			utils.RespondWithError(w, "Internal server error", 200)
			return
		}
		// if joiner online, send updated group status
		for client := range wsServer.Clients {
			if client.ID == joinerId {
				client.SendGroupRequestAccept(response.GroupID)
			}
		}
	}
	// delete from pending notification table
	if err = handler.repos.NotifRepo.Delete(response.RequestID); err != nil {
		utils.RespondWithError(w, "Internal server error", 200)
		return
	}
	utils.RespondWithSuccess(w, "Response was successful", 200)
}

// NOT TESTED
// waits for POST request with groupId and Invitation list included
func (handler *Handler) NewGroupInvite(wsServer *ws.Server, w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	if r.Method != "POST" {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	/* ---------------------------- read incoming data --------------------------- */
	// Try to decode the JSON request to a new Group
	var group models.Group
	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	userId := r.Context().Value(utils.UserKey).(string)
	/* -------- check if current user is a member or admin  of the group -------- */
	isAdmin, err := handler.repos.GroupRepo.IsAdmin(group.ID, userId)
	isMember, err := handler.repos.GroupRepo.IsMember(group.ID, userId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	if !isMember && !isAdmin {
		utils.RespondWithError(w, "Not a member", 200)
		return
	}
	for i := 0; i < len(group.Invitations); i++ {
		// save each invitation in db
		newNotif := models.Notification{
			ID:       utils.UniqueId(),
			TargetID: group.Invitations[i],
			Type:     "GROUP_INVITE",
			Content:  group.ID,
			Sender:   userId,
		}
		err = handler.repos.NotifRepo.Save(newNotif)
		if err != nil {
			utils.RespondWithError(w, "Internal server error", 200)
			return
		}

		// search if user has open ws connection and send notification
		for client := range wsServer.Clients {
			if client.ID == group.Invitations[i] {
				client.SendNotification(newNotif)
			}
		}
	}
	utils.RespondWithSuccess(w, "Invitations saved", 200)
}

// NOT TESTED
func (handler *Handler) ResponseInviteRequest(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	if r.Method != "POST" {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	/* ---------------------------- read incoming data --------------------------- */
	// Try to decode the JSON request to a new response
	type Response struct {
		RequestID string `json:"requestId"`
		Response  string `json:"response"` // ACCEPT or DECLINE
	}
	var resp Response
	err := json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	// get group id from notification
	groupId, err := handler.repos.NotifRepo.GetGroupId(resp.RequestID)
	userId := r.Context().Value(utils.UserKey).(string)
	if err != nil {
		utils.RespondWithError(w, "Internal server error", 200)
		return
	}
	if strings.ToUpper(resp.Response) == "ACCEPT" {
		err = handler.repos.GroupRepo.SaveMember(userId, groupId)
		if err != nil {
			utils.RespondWithError(w, "Internal server error", 200)
			return
		}
	}
	/* ----------------------- delete pending notification ---------------------- */
	err = handler.repos.NotifRepo.Delete(resp.RequestID)
	if err != nil {
		utils.RespondWithError(w, "Internal server error", 200)
		return
	}
	// notify websocket about notification changes
	utils.RespondWithSuccess(w, "Response successful", 200)
}

/* -------------------------------------------------------------------------- */
/*                              group management                              */
/* -------------------------------------------------------------------------- */

// DeleteGroup allows group creators to delete their groups
func (handler *Handler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	if r.Method != "DELETE" {
		utils.RespondWithError(w, "Method not allowed", 405)
		return
	}
	
	// Get group ID from URL parameters
	query := r.URL.Query()
	groupId := query.Get("groupId")
	if groupId == "" {
		utils.RespondWithError(w, "Group ID is required", 400)
		return
	}
	
	// Get current user ID
	userId := r.Context().Value(utils.UserKey).(string)
	
	// Delete the group
	err := handler.repos.GroupRepo.DeleteGroup(groupId, userId)
	if err != nil {
		fmt.Printf("Error deleting group: %v\n", err)
		if err.Error() == "user is not authorized to delete this group" {
			utils.RespondWithError(w, "Unauthorized to delete this group", 403)
		} else {
			utils.RespondWithError(w, "Failed to delete group", 500)
		}
		return
	}
	
	utils.RespondWithSuccess(w, "Group deleted successfully", 200)
}

// UpdateGroup allows group creators to update their groups
func (handler *Handler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	if r.Method != "PUT" && r.Method != "POST" {
		utils.RespondWithError(w, "Method not allowed", 405)
		return
	}
	
	// Get current user ID
	userId := r.Context().Value(utils.UserKey).(string)
	
	// Parse multipart form data for potential image upload
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		// If multipart fails, try regular form parsing
		err = r.ParseForm()
		if err != nil {
			utils.RespondWithError(w, "Error parsing form data", 400)
			return
		}
	}
	
	// Get group data from form
	groupId := r.FormValue("groupId")
	name := r.FormValue("name")
	description := r.FormValue("description")
	
	if groupId == "" || name == "" {
		utils.RespondWithError(w, "Group ID and name are required", 400)
		return
	}
	
	// Create group object
	group := models.Group{
		ID:          groupId,
		Name:        name,
		Description: description,
	}
	
	// Handle image upload if present
	file, fileHeader, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		
		// Save the uploaded image
		imagePath, err := utils.SaveUploadedFile(file, fileHeader)
		if err != nil {
			fmt.Printf("Error saving image: %v\n", err)
			utils.RespondWithError(w, "Error saving image", 500)
			return
		}
		group.ImagePath = imagePath
		fmt.Printf("Image saved to: %s\n", imagePath)
	} else {
		// If no new image, keep the existing one
		existingGroup, err := handler.repos.GroupRepo.GetData(groupId)
		if err == nil {
			group.ImagePath = existingGroup.ImagePath
		}
	}
	
	// Update the group
	err = handler.repos.GroupRepo.UpdateGroup(group, userId)
	if err != nil {
		fmt.Printf("Error updating group: %v\n", err)
		if err.Error() == "user is not authorized to update this group" {
			utils.RespondWithError(w, "Unauthorized to update this group", 403)
		} else {
			utils.RespondWithError(w, "Failed to update group", 500)
		}
		return
	}
	
	fmt.Printf("Group updated successfully: %+v\n", group)
	utils.RespondWithSuccess(w, "Group updated successfully", 200)
}
