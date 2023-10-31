package main

type ContactDetails struct {
	Person []struct {
		Age  int
		Name string
	}
	Email   string
	Address struct {
		Line1    string
		Country  string
		Postcode string
	}
}
