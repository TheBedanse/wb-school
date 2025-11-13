package main

import "fmt"

type Human struct {
	Name   string
	Age    int
	Gender string
}

func (h *Human) Speak() {
	fmt.Printf("Привет. Меня зовут %s, мне %d лет, я %s\n", h.Name, h.Age, h.Gender)
}

func (h *Human) Jump() {
	fmt.Printf("%s Прыгает\n", h.Name)
}

type Action struct {
	Human
	Work    string
	Married bool
}

func (a *Action) GoWork() {
	fmt.Printf("%s идет на работу в %s\n", a.Name, a.Work)
}

func (a *Action) FreeMan() {
	if a.Married {
		fmt.Printf("Мне %d и я Занят\n", a.Age)
	} else {
		fmt.Printf("Мне %d и я Свободен\n", a.Age)
	}

}

func main() {
	human := Human{
		Name:   "Том",
		Age:    20,
		Gender: "Парень",
	}

	action := Action{
		Human: Human{
			Name:   "Ника",
			Age:    21,
			Gender: "Девушка",
		},
		Work:    "Сбербанк",
		Married: false,
	}
	fmt.Println("Human")
	human.Speak()
	human.Jump()

	fmt.Println("Action")
	action.FreeMan()
	action.GoWork()

	fmt.Println("Action all methods")
	action.Speak()
	action.Jump()
	action.GoWork()
	action.FreeMan()
}
