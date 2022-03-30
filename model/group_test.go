package model_test

import (
	"chat_server/model"
	"encoding/json"
	"fmt"
	"testing"
)

func TestGroupToJson(t *testing.T) {
	g := &model.Group{
		BaseModel: model.BaseModel{
			ID: 1,
		},
		OwnerID: "2",
	}

	j, _ := json.Marshal(g)

	fmt.Println(string(j))
}
