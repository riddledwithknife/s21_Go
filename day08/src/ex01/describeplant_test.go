package ex01

import (
	"fmt"
	"testing"
)

func TestDescribePlant(t *testing.T) {
	plant1 := UnknownPlant{
		FlowerType: "Rose",
		LeafType:   "pinnate",
		Color:      10,
	}
	plant2 := AnotherUnknownPlant{
		FlowerColor: 10,
		LeafType:    "lanceolate",
		Height:      15,
	}

	describePlant(plant1)
	fmt.Println("")
	describePlant(plant2)
	fmt.Println("")
}
