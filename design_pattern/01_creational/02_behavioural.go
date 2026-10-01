package main

import "fmt"

type Builder struct {
	name, email string
	age         int
}

func NewBuilder(name string) *Builder {
	return &Builder{
		name: name,
	}
}

func (b *Builder) SetAge(age int) *Builder {
	b.age = age
	return b
}

func (b *Builder) SetEmail(email string) *Builder {
	b.email = email
	return b
}

func (b *Builder) Build() Builder {
	return *b
}

func (b Builder) String() string {
	return fmt.Sprintf("name=%s age=%d email=%s", b.name, b.age, b.email)
}

func main() {
	user := NewBuilder("test").
		SetAge(25).
		SetEmail("test@example.com").
		Build()

	fmt.Println("built:", user)

	// Build() returns a copy — later builder changes don't affect `user`
	b := NewBuilder("alex").SetAge(20)
	first := b.Build()
	b.SetAge(30)
	second := b.Build()

	fmt.Println("first: ", first)
	fmt.Println("second:", second)
}
