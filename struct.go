package main

type Config struct {
	User []User `toml:"User"`
}
type Cookie struct {
	Yxktml               string `toml:"yxktml"`
	RememberStudentKey   string `toml:"remember_student_key"`
	RememberStudentValue string `toml:"remember_student_value"`
}
type Postdata struct {
	V string
	Q string
	Z string
	X string
	W string
	A string
	Y string
	B string
	C string
	D string
	E string
	F string
	G string
	H string
	I string
	J string
	K string
	L string
	M string
	R string
	S string
	T string
	U string
}
type User struct {
	Cookie     Cookie   `toml:"Cookie"`
	Postdata   Postdata `toml:"Postdata"`
	FormID     string
	Name       string
	CheckinUrl string
	Mail       Mail `toml:"Mail"`
}

type Mail struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
}
