package ex01

import (
	"fmt"
	"reflect"
)

type UnknownPlant struct {
	FlowerType string
	LeafType   string
	Color      int `color_scheme:"rgb"`
}

type AnotherUnknownPlant struct {
	FlowerColor int
	LeafType    string
	Height      int `unit:"inches"`
}

func describePlant(plant interface{}) {
	v := reflect.ValueOf(plant)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag

		switch tag.Get("unit") {
		case "":
			fmt.Printf("%s:%v\n", t.Field(i).Name, field.Interface())
		default:
			fmt.Printf("%s(unit=%s):%v\n", t.Field(i).Name, tag.Get("unit"), field.Interface())
		}
	}
}
