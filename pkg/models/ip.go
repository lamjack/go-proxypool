package models

type Ip struct {
	Ip   string `bson:"ip"`
	Port int    `bson:"port"`
}
