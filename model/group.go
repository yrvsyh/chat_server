package model

type Group struct {
	BaseModel
	Name      string `gorm:"size:255" json:"name,omitempty"`
	PublicKey []byte `json:"public_key,omitempty"`
	OwnerID   string `json:"owner_id,omitempty"`
	Owner     *User  `json:"owner,omitempty"`

	Members []*User `gorm:"many2many:user_groups" json:"members,omitempty"`
}

func init() {
	db.SetupJoinTable(&Group{}, "Members", &UserGroups{})
	db.AutoMigrate(Group{})
}
