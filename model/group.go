package model

type Group struct {
	BaseModel
	Name      string `gorm:"size:255" json:"name"`
	PublicKey string `json:"public_key"`
	OwnerID   uint32 `json:"owner_id"`
	Owner     *User  `json:"owner,omitempty"`
}
