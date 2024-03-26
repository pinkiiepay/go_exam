package main

// import (
// 	"fmt"
// 	"regexp"

// 	"gopkg.in/go-playground/validator.v9"
// )

// type User struct {
// 	Username string `json:"username"`
// 	Email    string `json:"email" validate:"required,email"`
// 	Phone    string `json:"phone" validate:"required,len=10"`
// 	Password string `json:"password" validate:"required,min=6,max=12"`
// }

// func validatePassword(password string) bool {
// 	// Define el patrón regex
// 	pattern := `(([A-Z]+)([a-z]+)([@$&]+)(\d+))`
// 	re := regexp.MustCompile(pattern)
// 	return re.MatchString(password)
// }

// var validate *validator.Validate

// func main() {
// 	validate = validator.New()
// 	// password1 := "45F@de"
// 	// password2 := "Perrona69$"
// 	// password3 := "Short1!"

// 	// fmt.Printf("¿Es '%s' una contraseña válida? %v\n", password1, checkPassword(password1))
// 	// fmt.Printf("¿Es '%s' una contraseña válida? %v\n", password2, checkPassword(password2))
// 	// fmt.Printf("¿Es '%s' una contraseña válida? %v\n", password3, checkPassword(password3))
// 	nur := User{Username: "pedro", Email: "pedro@gmail.com", Phone: "fefe", Password: "pipoD$123"}
// 	fmt.Println(nur)
// 	if errs := validate.Struct(nur); errs != nil {
// 		fmt.Println(errs)
// 	}
// }

// func checksPassword(password string) bool {
// 	patterns := [4]string{"[A-Z]+", "[a-z]+", "[@$&]+", `\d+`}

// 	for i, s := range patterns {
// 		res, e := regexp.MatchString(s, password)
// 		fmt.Println(i, e, res)
// 		if res == false {
// 			return false
// 		}
// 	}
// 	return true
// }

// func checkPassword(password string) string {
// 	var patterns = make(map[string]string)
// 	patterns["[A-Z]+"] = "Error falta por lo menos una Mayúscula"
// 	patterns["[a-z]+"] = "Error falta por lo menos una Minúscula"
// 	patterns["[@$&]+"] = "Error falta por lo menos una carácter ( @$& )"
// 	patterns[`\d+`] = "Error falta por lo menos un Dígito"
// 	for key, value := range patterns {
// 		res, e := regexp.MatchString(key, password)
// 		fmt.Println(e, res)
// 		if res == false {
// 			return value
// 		}
// 	}
// 	return "OK"
// }

import (
	"database/sql"

	// Import the Azure AD driver module (also imports the regular driver package)
	"github.com/denisenkom/go-mssqldb/azuread"
)

func ConnectWithMSI() (*sql.DB, error) {
	return sql.Open(azuread.DriverName, "sqlserver://azuresql.database.windows.net?database=yourdb&fedauth=ActiveDirectoryMSI")
}
