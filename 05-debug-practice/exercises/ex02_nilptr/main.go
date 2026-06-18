package main

import "fmt"

type User struct {
	Name string
	Age  int
}

// findUser は名前が一致するユーザーを探して返す。
func findUser(users []*User, name string) *User {
	for _, u := range users {
		if u.Name == name {
			return u
		}
	}
	return nil
}

func main() {
	users := []*User{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
	}

	u := findUser(users, "Carol")
	fmt.Printf("%s さんは %d 歳です\n", u.Name, u.Age)
}
