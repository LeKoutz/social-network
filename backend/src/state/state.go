package state

import (
	"encoding/json"
	forum "forum/src"
	"forum/src/ferror"
	"forum/src/models"
	"forum/src/utils"
	"net/http"
)

type State struct {
	WebsiteName      string
	User             models.UserType             `json:"User,omitzero"`
	Users            models.UsersType            `json:"Users,omitzero"`
	Posts            models.PostsType            `json:"Posts,omitzero"`
	Categories       models.CategoriesType       `json:"Categories,omitzero"`
	Groups           models.GroupsType           `json:"Groups,omitzero"`
	Profiles         models.UserProfilesType     `json:"Profiles,omitzero"`
	EditingPost      bool                        `json:"EditingPost,omitzero"`
	EditCommentId    int64                       `json:"EditCommentId,omitzero"`
	ChatOffset       int64                       `json:"ChatOffset,omitzero"`
	Error            ferror.Error                `json:"Error,omitzero"`
	Request          *http.Request               `json:"-"`
	Response         http.ResponseWriter         `json:"-"`
	Message          models.Message              `json:"Message,omitzero"`
	FollowRequest    models.FollowRequestType    `json:"FollowRequest,omitzero"`
	Follows          models.InvitationsType      `json:"Follows,omitzero"`
	GroupInvitation  models.GroupInvitationType  `json:"GroupInvitation,omitzero"`
	GroupInvitations models.GroupInvitationsType `json:"GroupInvitations,omitzero"`
	Version          string                      `json:"Version,omitzero"`
}

type StateController interface {
	EditCategories() *models.CategoriesType
	GetCategories() models.CategoriesType
	SetCategories(models.CategoriesType) *State

	EditCategory() *models.CategoryType
	GetCategory() models.CategoryType
	SetCategory(models.CategoryType) *State

	EditGroup() *models.GroupType
	GetGroup() models.GroupType
	SetGroup(models.GroupType) *State

	EditGroups() *models.GroupsType
	GetGroups() models.GroupsType
	SetGroups(models.GroupsType) *State

	EditProfile() *models.UserProfileType

	EditPost() *models.PostType
	GetPost() models.PostType
	SetPost(models.PostType) *State
	SetEditPost(bool)

	SetPosts(models.PostsType) *State
	EditPosts() *models.PostsType

	EditComment() *models.CommentType
	GetComment() models.CommentType
	SetEditCommentId(int64)

	GetUser() models.UserType
	SetUser(models.UserType) *State
	EditUser() *models.UserType

	GetUsers() models.UsersType
	SetUsers(models.UsersType) *State
	EditUsers() *models.UsersType

	GetRequest() *http.Request
	EditResponse() *http.ResponseWriter

	GetChatOffset() int64

	EditChatMessages() *models.ChatMessagesType
	GetChatMessages() models.ChatMessagesType
	EditChatMessage(index int64) *models.ChatMessageType
	GetChatMessage(index int64) models.ChatMessageType

	EditFollowRequest() *models.FollowRequestType
	EditFollows() *models.InvitationsType

	EditGroupInvitation() *models.GroupInvitationType
	EditGroupInvitations() *models.GroupInvitationsType

	SetMessage(models.Message) *State
	SetErrorConsume(error) *State
}

type StateHandler interface {
	Init() *State
	SetErrorConsume(error) *State
	GetError() *ferror.Error

	GetUser() models.UserType
	EditUser() *models.UserType
	SetCategories(models.CategoriesType) *State
	EditCategories() *models.CategoriesType
	EditCategory() *models.CategoryType
	EditGroup() *models.GroupType
	EditComment() *models.CommentType
	EditPost() *models.PostType
	GetPost() models.PostType
	EditPosts() *models.PostsType
	EditChatMessages() *models.ChatMessagesType
	EditChatMessage(index int64) *models.ChatMessageType
	SetChatOffset(int64) *State
	EditFollowRequest() *models.FollowRequestType
	EditProfile() *models.UserProfileType
	EditFollows() *models.InvitationsType
	EditGroupInvitation() *models.GroupInvitationType

	GetRequest() *http.Request
	EditResponse() *http.ResponseWriter

	WriteResponse()
}

type StateRoute interface {
	Init() *State
	SetErrorConsume(error) *State
	GetError() *ferror.Error
	SetResponse(res http.ResponseWriter) *State

	GetRequest() *http.Request
	SetRequest(req *http.Request) *State

	GetUser() models.UserType
	SetUser(models.UserType) *State
	EditUser() *models.UserType

	WriteResponse()
}

func (r *State) SetChatOffset(offset int64) *State {
	r.ChatOffset = offset
	return r
}

func (r *State) GetChatOffset() int64 {
	return r.ChatOffset
}

func (r State) EditResponse() *http.ResponseWriter {
	return &r.Response
}

func (r *State) SetMessage(m models.Message) *State {
	r.Message = m
	return r
}

func (r *State) GetError() *ferror.Error {
	return &r.Error
}

func (r *State) GetRequest() *http.Request {
	return r.Request
}

func (r *State) WriteResponse() {
	if r.Error.StatusCode != 0 {
		r.Response.WriteHeader(r.Error.StatusCode)
	}
	respondView(*r)
}

func (r *State) Init() *State {
	r.WebsiteName = forum.WebsiteName
	r.InitUser()
	r.Version = utils.GetVersion()
	return r
}

func (r *State) SetWebsiteName(websiteName string) *State {
	r.WebsiteName = websiteName
	return r
}

func (r *State) SetError(err ferror.Error) *State {
	r.Error = err
	return r
}

func (r *State) SetErrorConsume(err error) *State {
	r.Error.Consume(err)
	r.Error.LogError()
	return r
}

func (r *State) SetRequest(req *http.Request) *State {
	r.Request = req
	return r
}

func (r *State) SetResponse(res http.ResponseWriter) *State {
	r.Response = res
	return r
}

func (r *State) GetResponse(res http.ResponseWriter) http.ResponseWriter {
	return r.Response
}

func respondView(data State) {
	bs, err := json.Marshal(data)
	if err != nil {
		data.Response.Header().Add("Content-Type", "application/json")
		data.SetErrorConsume(err)
		return
	}
	data.Response.Header().Add("Content-Type", "application/json")
	_, err = data.Response.Write(bs)
	if err != nil {
		data.SetErrorConsume(err)
		return
	}
}
