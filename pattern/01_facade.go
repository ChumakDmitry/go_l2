package main

import "fmt"

type SubSys1 struct{}

func (sub1 *SubSys1) Method1() {
	fmt.Println("SubSys1 method")
}

type SubSys2 struct{}

func (sub2 *SubSys2) Method2() {
	fmt.Println("SubSys2 method")
}

type SubSys3 struct{}

func (sub3 *SubSys3) Method3() {
	fmt.Println("SubSys3 method")
}

type Facade struct {
	SubSys1 *SubSys1
	SubSys2 *SubSys2
	SubSys3 *SubSys3
}

func NewFacade() *Facade {
	return &Facade{
		SubSys1: &SubSys1{},
		SubSys2: &SubSys2{},
		SubSys3: &SubSys3{},
	}
}

func (f *Facade) FacMethod1() {
	fmt.Println("Facade method 1")
	f.SubSys1.Method1()
	f.SubSys2.Method2()
}
func (f *Facade) FacMethod2() {
	fmt.Println("Facade method 2")
	f.SubSys1.Method1()
	f.SubSys3.Method3()
}

func main() {
	Facade := NewFacade()
	Facade.FacMethod1()
	Facade.FacMethod2()
}
