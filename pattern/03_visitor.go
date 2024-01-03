package main

import "fmt"

type Visitor interface {
	VisitPerson(*Person)
	VisitOrganization(*Organization)
}

type Element interface {
	Accept(Visitor)
}

type Person struct {
	name  string
	email string
}

func (p *Person) Accept(v Visitor) {
	v.VisitPerson(p)
}

type Organization struct {
	name    string
	address string
}

func (o *Organization) Accept(v Visitor) {
	v.VisitOrganization(o)
}

type EmailVisitor struct{}

func (e *EmailVisitor) VisitPerson(p *Person) {
	fmt.Printf("Sending email to %s at %s\n", p.name, p.email)
}

func (e *EmailVisitor) VisitOrganization(o *Organization) {
	fmt.Printf("Sending mail to %s at %s\n", o.name, o.address)
}

func main() {
	elements := []Element{
		&Person{name: "Alice", email: "alices@example.com"},
		&Organization{name: "Acme Inc.", address: "123 Main St."},
		&Person{name: "Bob", email: "bob@example.com"},
	}

	visitor := &EmailVisitor{}

	for _, element := range elements {
		element.Accept(visitor)
	}

}
