package model

type Group struct {
	BaseModel
	Name      string
	PublicKey []byte

	Users []*User `gorm:"many2many:user_groups"`
}

func init() {
	db.SetupJoinTable(&Group{}, "Users", &UserGroups{})
	db.AutoMigrate(Group{})
}
