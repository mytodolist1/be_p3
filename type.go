package bep3

type User struct {
	Email    string `bson:"email,omitempty" json:"email,omitempty"`
	Username string `bson:"username,omitempty" json:"username,omitempty"`
	Password string `bson:"password,omitempty" json:"password,omitempty"`
	Role     string `json:"role,omitempty" bson:"role,omitempty"`
}

type Credential struct {
	Status  bool   `json:"status" bson:"status"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

type Task struct {
	ID          int    `json:"id,omitempty" bson:"id,omitempty"`
	Name        string `json:"name,omitempty" bson:"name,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	Deadline    string `json:"deadline,omitempty" bson:"deadline,omitempty"`
	Done        bool   `json:"done,omitempty" bson:"done,omitempty"`
}
