# Go Denest Structs v0.1

I started this project to help denest structs in golang.  
  
Often find myself using tools like xml to go-struct or json to go-struct. However these tools only spit out a single struct.  

## usage

Denest all structs within a file

    // Example
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
`cat path/to/my_golang_file.go | go run cmd/go-denest-structs/main.go > path/to/my_destructured_golang_file.txt`

The output structs will not be formatted.

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

>You can pipe the output through 'gofmt' to fix formatting:  
`cat path/to/my_golang_file.go | go run cmd/go-denest-structs/main.go | gofmt > path/to/my_destructured_golang_file.txt`

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

# Install as binary