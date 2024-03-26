package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/golang-jwt/jwt"
	_ "github.com/microsoft/go-mssqldb"
	"gopkg.in/go-playground/validator.v9"
)

var (
	debug         = flag.Bool("debug", true, "enable debugging")
	password      = flag.String("password", "Perrona69*", "the database password")
	port     *int = flag.Int("port", 1433, "the database port")
	server        = flag.String("server", "localhost", "the database server")
	user          = flag.String("user", "SA", "the database user")
)

//4. Validar que el teléfono sea a 10 dígitos y el correo tenga un formato válido.

type User struct {
	Username string `json:"username"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,len=10"`
	Password string `json:"password" validate:"required,min=6,max=12"`
}

var MySigningKey = []byte(os.Getenv("SECRET_KEY"))

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {

			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf(("Invalid Signing Method"))
				}
				aud := "test_go"
				checkAudience := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
				if !checkAudience {
					return nil, fmt.Errorf(("invalid aud"))
				}
				// verify iss claim
				iss := "test_go"
				checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
				if !checkIss {
					return nil, fmt.Errorf(("invalid iss"))
				}

				return MySigningKey, nil
			})
			if err != nil {
				fmt.Fprintf(w, err.Error())
			}

			if token.Valid {
				endpoint(w, r)
			}

		} else {
			fmt.Fprintf(w, "No Authorization Token provided")
		}
	})
}
func handleRequests() {
	//1. Crear un servicio de registro de usuario que reciba como parámetros usuario, correo, teléfono y contraseña.
	http.Handle("/register", isAuthorized(registerHandler))
	log.Fatal(http.ListenAndServe(":9001", nil))
}

var validate *validator.Validate

func main() {
	handleRequests()
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	validate = validator.New()
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	fmt.Println(err)
	if err != nil {
		http.Error(w, "Cuerpo de solicitud inválido", http.StatusBadRequest)
		return
	}
	//8. En ambos servicios se deberá validar que todos los parámetros solicitados vayan en el cuerpo de la petición, de lo contrario retorna un mensaje con el campo faltante.
	err_valid := validate.Struct(user)
	if err_valid != nil {
		fmt.Println(err_valid)
		str_error := "Error en uno de los campos" + err_valid.Error()
		http.Error(w, str_error, http.StatusBadRequest)
		return
	}
	err_pass := checkPassword(user.Password)
	if err_pass != "OK" {
		http.Error(w, err_pass, http.StatusBadRequest)
		return
	}

	// 2. El servicio deberá validar que el correo y telefono no se encuentren registrados, de lo contrario deberá retornar un mensaje “el correo/telefono ya se encuentra registrado”.
	conn := getConn()
	defer conn.Close()
	isRepeat := getUserRepeat(user.Email, user.Phone, conn)

	if isRepeat {
		fmt.Fprintln(w, "Usuario ya existente, verifique los datos")
	} else {
		conn := getConn()
		defer conn.Close()
		isInsert, err := insertUser(user.Username, user.Email, user.Phone, user.Password, conn)
		if isInsert >= 0 {
			fmt.Printf("Usuario registrado: %+v\n", user)
			// Responder con un mensaje de éxito
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintln(w, "Usuario registrado exitosamente")
		} else {
			fmt.Fprintln(w, "Hubo un error, pruebe de nuevo en unos instantes: ", err)
		}

	}

}

//3. Deberá validar que la contraseña sea de 6 caracteres mínimo y 12 máximo y contener al menos una mayúscula, una minúscula, un carácter especial (@ $ o &) y un número.

func checkPassword(password string) string {
	var patterns = make(map[string]string)
	patterns["[A-Z]+"] = "Password Error falta por lo menos una Mayúscula"
	patterns["[a-z]+"] = "Password Error falta por lo menos una Minúscula"
	patterns["[@$&]+"] = "Password Error falta por lo menos una carácter ( @$& )"
	patterns[`\d+`] = "Password Error falta por lo menos un Dígito"
	for key, value := range patterns {
		res, e := regexp.MatchString(key, password)
		fmt.Println(e, res)
		if !res {
			return value
		}
	}
	return "OK"
}

func getUserRepeat(mail, phone string, conn *sql.DB) bool {
	sql_string := fmt.Sprintf("SELECT id, user_name FROM EXAM_GO.exam_go.Users where email='%s' or phone='%s'", mail, phone)
	fmt.Println(sql_string)
	stmt, err := conn.Prepare(sql_string)
	if err != nil {
		log.Fatal("Prepare failed:", err.Error())
	}
	defer stmt.Close()

	row := stmt.QueryRow()
	var id int64
	var user_name string
	err = row.Scan(&id, &user_name)
	if err != nil {
		return false
	}
	fmt.Printf("id:%d\n", id)
	fmt.Printf("user_name:%s\n", user_name)

	return true
}

func insertUser(name, mail, phone, password string, conn *sql.DB) (int64, error) {
	sql_string := fmt.Sprintf("INSERT INTO EXAM_GO.exam_go.Users (user_name, email, phone, password) VALUES('%s', '%s', '%s', '%s');", name, mail, phone, password)
	fmt.Println(sql_string)
	_, err := conn.Exec(sql_string)
	if err != nil {
		fmt.Println("Error inserting new row: " + err.Error())
		return -1, err
	}
	return 1, err
}

func getConn() *sql.DB {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d", *server, *user, *password, *port)
	if *debug {
		fmt.Printf(" connString:%s\n", connString)
	}
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	return conn
}
