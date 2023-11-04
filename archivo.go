//conectar con bd
package main

import (
	//"fmt"
	"log"
	"database/sql"
	_ "github.com/lib/pq"
)

func dbConnection()(*sql.DB, error){
       db, err := sql.Open("postgres", "user=postgres host=localhost dbname=Hospital sslmode=disable")
       if err != nil{
       
       log.Fatal(err)
}
return db, nil
}

// crear base de datos

func crearDB() {
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	_, err = db.Exec("create DATABASE Hospital")

	if err != nil {
		log.Fatal(err)
	}
}


// crear tablas

func createTables() {
	db, err := dbConnection()

	_, err = db.Exec("DROP SCHEMA public CASCADE")

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("create SCHEMA public")

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("create table paciente (nro_paciente int, nombre text, apellido text, dni paciente int, f_nac date, nro_obra_social int, nro afiliade int, domicilio text, telefono char(12), email text)")

	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()
}

func main() {
	_ = pq.QuoteIdentifier("some_text")
	
	dbConnection()
	crearDB()
	createTables()
	
	//conexion a base de datos postgres
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=Hospital sslmode=disable")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()
}
