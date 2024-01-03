package main

import "fmt"

const (
	IPhoneType  = "iphone"
	SamsungType = "samsung"
	XiaomiType  = "xiaomi"
)

type Phone interface {
	GetType() string
	PrintDetails()
}

func New(typeName string) Phone {
	switch typeName {
	case IPhoneType:
		return NewIPhone()
	case SamsungType:
		return NewSamsung()
	case XiaomiType:
		return NewXiaomi()
	default:
		fmt.Printf("%s неверный тип\n", typeName)
		return nil
	}
}

type IPhone struct {
	Type            string
	Model           string
	Storage         int
	AlwaysOnDisplay bool
}

func NewIPhone() Phone {
	return IPhone{
		Type:            IPhoneType,
		Model:           "12 mini",
		Storage:         128,
		AlwaysOnDisplay: false,
	}
}

func (i IPhone) GetType() string {
	return i.Type
}

func (i IPhone) PrintDetails() {
	fmt.Printf("%s, Model: %s, RAM: %d, AlwaysOn: %t\n", i.Type, i.Model, i.Storage, i.AlwaysOnDisplay)
}

type Samsung struct {
	Type  string
	Model string
	RAM   int
}

func NewSamsung() Phone {
	return Samsung{
		Type:  SamsungType,
		Model: "A52",
		RAM:   8,
	}
}

func (s Samsung) GetType() string {
	return s.Type
}

func (s Samsung) PrintDetails() {
	fmt.Printf("%s, Model: %s, RAM: %d\n", s.Type, s.Model, s.RAM)
}

type Xiaomi struct {
	Type  string
	Model string
	RAM   int
}

func NewXiaomi() Phone {
	return Xiaomi{
		Type:  XiaomiType,
		Model: "Note 10",
		RAM:   8,
	}
}

func (x Xiaomi) GetType() string {
	return x.Type
}

func (x Xiaomi) PrintDetails() {
	fmt.Printf("%s, Model: %s, RAM: %d\n", x.Type, x.Model, x.RAM)
}

func main() {
	var phones = []string{IPhoneType, "POCO", SamsungType, XiaomiType}

	for _, typeName := range phones {
		phone := New(typeName)
		if phone == nil {
			continue
		}

		phone.PrintDetails()
	}
}
