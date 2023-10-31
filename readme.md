# Go Denest Structs v0.1

I started this project to help denest structs in golang.  
  
Often find myself using tools like xml to go-struct or json to go-struct. However these tools only spit out a single struct.  

## Command-line usage

Denest all structs within a file

    // example.go
    package main

    type ContactDetails struct {
        Person struct {
            Age  int
            Name string
        }
        Email    string
        Address  struct {
            Line1    string
            Country  string
            Postcode string
        }
    }

>Execute:   
`cat example.go | go run main.go > result.go`

    // result.go
    package main

    type ContactDetails struct {
            Person Person
            Email   string
            Address Address
    }
    type Person struct {
                    Age  int
                    Name string
            }
    type Address struct {
                    Line1    string
                    Country  string
                    Postcode string
            }

