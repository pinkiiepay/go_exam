package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "github.com/microsoft/go-mssqldb"
)

var (
	debug         = flag.Bool("debug", true, "enable debugging")
	password      = flag.String("password", "Perrona69*", "the database password")
	port     *int = flag.Int("port", 1433, "the database port")
	server        = flag.String("server", "localhost", "the database server")
	user          = flag.String("user", "SA", "the database user")
)

func main() {
	flag.Parse()

	if *debug {
		// fmt.Printf(" password:%s\n", *password)
		// fmt.Printf(" port:%d\n", *port)
		// fmt.Printf(" server:%s\n", *server)
		// fmt.Printf(" user:%s\n", *user)
	}
	conn := getConn()
	genToken := getUser("pedro", "pedro@gmail.com", "Pipo65$", conn)
	if genToken {
		fmt.Println("fasilito")
	} else {
		fmt.Println("error")
	}
	conn.Close()

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
	if *debug {
		fmt.Printf(" connString:%s\n", connString)
	}
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	return conn
}
