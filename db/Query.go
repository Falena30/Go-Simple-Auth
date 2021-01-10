package db

import "fmt"

//User digunakan untuk untuk menyimpan nilai dari tb user
type User struct {
	ID       int
	Username string
	Firsname string
	Lastname string
	Password string
}

//QueryUser digunakan untuk melakukan pemilihan user
func QueryUser(Username string) User {
	//butlah variabel untuk menampung user
	var users = User{}
	//Panggilah funsi konek
	databse, err := Connect()
	if err != nil {
		fmt.Println(err.Error())
	}
	defer databse.Close()
	//Lakukan query dan masukkan nilai dari user
	err = databse.QueryRow("SELECT * FROM users WHERE Username = ?", Username).Scan(
		&users.ID,
		&users.Username,
		&users.Firsname,
		&users.Lastname,
		&users.Password,
	)
	//kembalikan nilai users
	return users
}
