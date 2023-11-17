package main

import (
    "encoding/json"
    "fmt"
    bolt "go.etcd.io/bbolt"
    "log"
    "strconv"
)

type Paciente struct {
	nro_paciente int
	nombre  string
	apellido string
	dni_paciente int
	f_nacimiento string
	nro_obra_social int
	nro_afiliade int
	domicilio string
	telefono int
	email string
}

type Medique struct {
	dni_medique int
	nombre string
	apellido string
	especialidad string
	monto_consulta_privada float64
	telefono int
}

type Consultorio struct {
	nro_consultorio int
	nombre string
	domicilio string
	codigo_postal int
	telefono int
}

type Turno struct {
	nro_turno int
	fecha string
	nro_consultorio int
	dni_medique int
	nro_paciente int
	nro_obra_social_consulta int
	nro_afiliade_consulta int
	monto_paciente float64
	monto_obra_social float64
	f_reserva string
	estado string	
}

type Obra_social struct {
	nro_obra_social int
	nombre string
	contacto_nombre string
	contacto_apellido string
	contacto_telefono int
	contacto_email string
}


func CreateUpdate(db *bolt.DB, bucketName string, key []byte, val []byte) error {
    // abre transacción de escritura
    tx, err := db.Begin(true)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    b, _ := tx.CreateBucketIfNotExists([]byte(bucketName))

    err = b.Put(key, val)
    if err != nil {
        return err
    }

    // cierra transacción
    if err := tx.Commit(); err != nil {
        return err
    }

    return nil
}

func ReadUnique(db *bolt.DB, bucketName string, key []byte) ([]byte, error) {
    var buf []byte

    // abre una transacción de lectura
    err := db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(bucketName))
        buf = b.Get(key)
        return nil
    })

    return buf, err
}

func main() {
    db, err := bolt.Open("hospital.db", 0600, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

	/*Chequear/Adaptar
    cristina := Alumne{1, "Cristina", "Kirchner"}
    data, err := json.Marshal(cristina)
    if err != nil {
        log.Fatal(err)
    }
	
	
    CreateUpdate(db, "alumne", []byte(strconv.Itoa(cristina.Legajo)), data)

    resultado, err := ReadUnique(db, "alumne", []byte(strconv.Itoa(cristina.Legajo)))
    */    

    fmt.Printf("%s\n", resultado)
}

