package main

import "fmt"

type A struct {
	Name string
}

type B struct {
	A
}

type I interface {
	Get()
}

func (b *A) Get() {
	b.Name = "zs"
}

func main2() {

	b := &B{}
	b.A.Get()

	fmt.Println(b.Name)

}
