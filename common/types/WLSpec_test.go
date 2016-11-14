package types

import (
	"testing"
)

func TestCopy(T *testing.T) {

	var Source, Destination WLSpec

	Source.Default()
	Destination.Copy(Source)

	if Destination.Image != "Stateful:latest" && Destination.CPU == 1.0 {
		T.Fail()
	}
}

func TestFromJson(T *testing.T) {

	var Source WLSpec
	var Destination WLSpec

	Source.Default()

	Destination.FromJson(Source.ToJson())

	if Destination.Image != "Stateful:latest" && Destination.CPU == 1.0 {
		T.Fail()
	}
}
