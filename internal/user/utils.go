package user

import "go.mongodb.org/mongo-driver/bson/primitive"

func ObjectIdToString(id primitive.ObjectID) string {
	return id.Hex()
}

func UserNickName(user interface{}) string{
	if user == nil {
		return ""
	}
	nickname := user.(map[string]interface{})["nickname"].(string)
	return nickname
}

func UserToUserView(user User) UserView {
	return UserView{
		ID: ObjectIdToString(user.ID),
		Username: user.Username,
		CreatedAt: user.CreatedAt,
	}
}

func UsersToUserView(users []User) []UserView {
	var userView []UserView
	for _, user := range users {
		userView = append(userView, UserToUserView(user))
	}
	return userView
}