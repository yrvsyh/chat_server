package model

type Group struct {
	BaseModel
	Name      string `gorm:"size:255" json:"name,omitempty"`
	PublicKey []byte `json:"public_key,omitempty"`
	OwnerID   string `json:"owner_id,omitempty"`
	Owner     *User  `json:"owner,omitempty"`
}

func init() {
	db.AutoMigrate(Group{})
}
