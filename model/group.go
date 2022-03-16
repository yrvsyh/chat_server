package model

type Group struct {
	BaseModel
	Name      string
	PublicKey []byte
	OwnerID   string
	Owner     *User `json:"owner,omitempty"`

	Members []*User `gorm:"many2many:user_groups"`
}

func init() {
	db.SetupJoinTable(&Group{}, "Members", &UserGroups{})
	db.AutoMigrate(Group{})
}
