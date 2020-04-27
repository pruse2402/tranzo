package models

import (
	"tranzo/utils"

	"gopkg.in/mgo.v2/bson"
)

type Details struct {
	Id     bson.ObjectId `bson:"_id" json:"id"`
	Name   string        `bson:"name" json:"name"`
	Age    string        `bson:"age" json:"age"`
	Gender string        `bson:"gender" json:"gender"`
}

func (d *Details) Validate() (bool, map[string]interface{}) {

	v := &utils.Validation{}

	v.Required(d.Name).Key("name").Message("Name is missing")

	return v.HasErrors(), v.ErrorMap()
}
