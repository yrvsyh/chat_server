package model

type Group struct {
	BaseModel
	Name      string `gorm:"size:255" json:"name"`
	Avatar    string `gorm:"size:255" json:"avatar"`
	PublicKey string `json:"public_key"`
	OwnerID   uint32 `json:"owner_id"`
	Owner     *User  `json:"owner,omitempty"`
	Label     string `json:"label"`
	Type      string `json:"type"`
}
