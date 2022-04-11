package model

type Group struct {
	BaseModel
	Name      string `gorm:"size:255" json:"name"`
	PublicKey []byte `json:"public_key,omitempty"`
	OwnerID   uint32 `json:"owner_id"`
	Owner     *User  `json:"owner,omitempty"`
}
