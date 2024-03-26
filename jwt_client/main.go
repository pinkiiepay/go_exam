package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt"
	_ "github.com/microsoft/go-mssqldb"
	"gopkg.in/go-playground/validator.v9"
)

var (
	password      = flag.String("password", "Perrona69*", "the database password")
	port     *int = flag.Int("port", 1433, "the database port")
	server        = flag.String("server", "localhost", "the database server")
	user          = flag.String("user", "SA", "the database user")
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=12"`
}

var validate *validator.Validate
var mySigningKey = []byte(os.Getenv("SECRET_KEY"))

func GetJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = "test_go"
	claims["aud"] = "test_go"
	claims["iss"] = "test_go"
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func Index(w http.ResponseWriter, r *http.Request) {
	validate = validator.New()
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
	conn := getConn()
	defer conn.Close()
	verifyUser := getUser(user.Username, user.Email, user.Password, conn)
	//7. Deberá validar que el usuario o correo y contraseña sean válidos, de lo contrario retorna un mensaje “usuario / contraseña incorrectos”.
	if verifyUser {
		validToken, err := GetJWT()
		fmt.Println(validToken)
		if err != nil {
			fmt.Println("Failed to generate token")
		}
		//6. El servicio debe devolver un token jwt.
		fmt.Fprintf(w, string(validToken))
	} else {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "Error, Usuario no existe, verifique sus datos.")
	}

}

func handleRequests() {
	//5. Crear un servicio login que reciba como parámetros usuario o correo y contraseña.
	http.HandleFunc("/login", Index)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	handleRequests()
}

func getUser(user, mail, password string, conn *sql.DB) bool {
	sql_string := fmt.Sprintf("SELECT id, user_name FROM EXAM_GO.exam_go.Users where user_name='%s' and email='%s' and password='%s'", user, mail, password)
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

func getConn() *sql.DB {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d", *server, *user, *password, *port)
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	return conn
}
