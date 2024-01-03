package main

import "fmt"

type Computer struct {
	CPU       string
	VideoCard string
	RAM       int
}

type BuilderI interface {
	CPU(val string) BuilderI
	VideoCard(val string) BuilderI
	RAM(val int) BuilderI

	Build() Computer
}

type Builder struct {
	cpu       string
	videoCard string
	ram       int
}

func NewBuilder() *Builder {
	return &Builder{
		cpu:       "Ryzen 5",
		videoCard: "GTX 2060",
		ram:       8,
	}
}

func (b *Builder) CPU(val string) BuilderI {
	b.cpu = val
	return b
}

func (b *Builder) RAM(val int) BuilderI {
	b.ram = val
	return b
}

func (b *Builder) VideoCard(val string) BuilderI {
	b.videoCard = val
	return b
}

func (b *Builder) Build() Computer {
	return Computer{
		CPU:       b.cpu,
		VideoCard: b.videoCard,
		RAM:       b.ram,
	}
}

func main() {
	builder := NewBuilder()
	computer := builder.CPU("Intel i5 10400f").Build()
	fmt.Println(computer)
}
