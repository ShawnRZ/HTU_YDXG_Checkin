package main

import "net/http"

type Config struct {
	User []User `toml:"User"`
}
type Cookie struct {
	Yxktml               string `toml:"yxktml"`
	RememberStudentKey   string `toml:"remember_student_key"`
	RememberStudentValue string `toml:"remember_student_value"`
}
type User struct {
	Cookie     Cookie            `toml:"Cookie"`
	Postdata   map[string]string `toml:"Postdata"`
	FormID     string
	Name       string
	CheckinUrl string
	Mail       Mail `toml:"Mail"`
	Client     *http.Client
}

type Mail struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
}
