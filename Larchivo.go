 /conectar con bd
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
    db, err := db.Exec(`
	create table paciente(nro_paciente int, nombre  text, apellido text, dni_paciente int, f_nacimiento date, nro_obra_social int, nro_afiliade int, domicilio text, telefono char(12),  email text); 
	
	create table medique(dni_medique int, nombre text, apellido text, especialidad varchar(64), monto_consulta_privada decimal(12, 2), telefono char(12)); 
	
	create table consultorio(nro_consultorio int, nombre text, domicilio text, codigo_postal char(8), telefono char(12));

	create table agenda(dni_medique int, dia int, nro_consultorio int, hora_desde time, hora_hasta time, duracion_turno interval);

	create table turno( nro_turno int, fecha timestamp, nro_consultorio int, dni_medique int, nro_paciente int, nro_obra_social_consulta int, nro_afiliade_consulta int, monto_paciente decimal(12,2), monto_obra_social decimal(12,2), f_reserva timestamp, estado char(10));

	create table reprogramacion(nro_turno int, nombre_paciente text, apellido_paciente text, telefono_paciente char(12), email_paciente text, nombre_medique text, apellido_medique text, estado char(12));

	create table error(nro_error int, f_turno timestamp, nro_consultorio int, dni_medique int, nro_paciente int, operacion char(12), f_error timestamp, motivo varchar(64));

	create table cobertura(dni_medique int, nro_obra_social int, monto_paciente decimal(12,2), monto_obra_social decimal(12,2));

	create table obra_social (nro_obra_social int, nombre text, contacto_nombre text, contacto_apellido text, contacto_telefono char(12), contacto_email text);
	
	create table liquidacion_cabecera(nro_liquidacion int, nro_obra_social int, desde date, hasta date, total decimal(15,2));

	create table liquidacion_detalle(nro_liquidacion int, nro_linea int, f_atencion date, nro_afiliade int, dni_paciente int, nombre_paciente text, apellido_paciente text, dni_medique int, nombre_medique text, apellido_medique text, especialidad varchar(64), monto decimal(12,2));


	create table envio_email(nro_email int, f_generacion timestamp, email_paciente text, asunto text, cuerpo text, f_envio timestamp, estado char(10));

	create table solicitud_reservas(nro_orden int, nro_paciente int, dni_medique int, fecha date, hora time);
	`)
   
	if err != nil {
        log.Fatal(err)
    }

    defer db.Close()
    _, err = db.Exec("create database Hospital")

    if err != nil {
        log.Fatal(err)
    }
}


// crear tablas

func createTables() {
    db, err := dbConnection()

