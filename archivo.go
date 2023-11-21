//conectar con bd
package main

import (
	//"fmt"
	"log"
	"database/sql"
	_ "github.com/lib/pq"
)

func main() {
	_ = pq.QuoteIdentifier("some_text")
	
	fmt.println("Elegi una opcion:")
	fmt.println("1 Crear BD")
	fmt.println("2. Crear Tablas")
	fmt.println("3. Crear PKs y FKs")
	fmt.println("4. Cargar Tablas")
	fmt.println("5. Crear Stored Procedures y Triggers")
	fmt.println("6. Eliminar PKs y FKs")
	
	var opcion int
	_, err:=fmt.scanln(&opcion)
		

	switch {
		case opcion == 1:
			crearDB()
		case opcion == 2:
			createTables()
		case opcion == 3:
			createPK()
			createFK()
		case opcion == 4:
			cargarTablas()
		case opcion == 5:
			crearSP()
			crearTriggers()
		case opcion == 6:
			eliminarPK()
			eliminarFK()
	}
	
}

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

	_, err = db.Exec(
		`create table paciente (nro_paciente int, nombre  text, apellido text, dni_paciente int, f_nacimiento date, nro_obra_social int, nro_afiliade int, domicilio text, telefono char(12),  email text); 
		create table medique (dni_medique int, nombre text, apellido text, especialidad varchar(64), monto_consulta_privada decimal(12, 2), telefono char(12)); 
		create table consultorio (nro_consultorio int, nombre text, domicilio text, codigo_postal char(8), telefono char(12));
		create table agenda (dni_medique int, dia int, nro_consultorio int, hora_desde time, hora_hasta time, duracion_turno interval);
		create table turno (nro_turno int, fecha timestamp, nro_consultorio int, dni_medique int, nro_paciente int, nro_obra_social_consulta int, nro_afiliade_consulta int, monto_paciente decimal(12,2), monto_obra_social decimal(12,2), f_reserva timestamp, estado char(10));
		create table reprogramacion (nro_turno int, nombre_paciente text, apellido_paciente text, telefono_paciente char(12), email_paciente text, nombre_medique text, apellido_medique text, estado char(12));
		create table error (nro_error int, f_turno timestamp, nro_consultorio int, dni_medique int, nro_paciente int, operacion char(12), f_error timestamp, motivo varchar(64));
		create table cobertura (dni_medique int, nro_obra_social int, monto_paciente decimal(12,2), monto_obra_social decimal(12,2));
		create table obra_social (nro_obra_social int, nombre text, contacto_nombre text, contacto_apellido text, contacto_telefono char(12), contacto_email text);
		create table liquidacion_cabecera (nro_liquidacion int, nro_obra_social int, desde date, hasta date, total decimal(15,2));
		create table liquidacion_detalle (nro_liquidacion int, nro_linea int, f_atencion date, nro_afiliade int, dni_paciente int, nombre_paciente text, apellido_paciente text, dni_medique int, nombre_medique text, apellido_medique text, especialidad varchar(64), monto decimal(12,2));
		create table envio_email (nro_email int, f_generacion timestamp, email_paciente text, asunto text, cuerpo text, f_envio timestamp, estado char(10));
		create table solicitud_reservas (nro_orden int, nro_paciente int, dni_medique int, fecha date, hora time);`
		)
		

	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()
}

func createPK() {
	db, err := dbConnection()
	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()
	
	_, err = db.Exec (`alter table paciente add constraint paciente_pk primary key(nro_paciente);
		alter table medique add constraint medique_pk primary key (dni_medique);
		alter table consultorio add constraint consultorio_pk primary key (nro_consultorio);
		alter table agenda add constraint agenda_pk primary key (dni_medique, dia);
		alter table turno add constraint turno_pk primary key (nro_turno);
		alter table reprogramacion add constraint reprogramacion_pk primary key (nro_turno);
		alter table error add constraint error_pk primary key (nro_error);
		alter table cobertura add constraint cobertura_pk primary key (dni_medique, nro_obra_social);
		alter table obra_social add constraint obra_social_pk primary key (nro_obra_social);
		alter table liquidacion_cabecera add constraint liquidacion_cabecera_pk primary key (nro_liquidacion);
		alter table liquidacion_detalle add constraint liquidacion_detalle_pk primary key (nro_liquidacion, nro_linea);
		alter table envio_email add constraint envio_emailpk primary key (nro_email);`
	)
	if err != nil{
		log.Fatal(err)
	}
}

func createFK() {
		db, err := dbConnection()
	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()
	
	_, err = db.Exec (`alter table paciente add constraint fk_paciente foreign key (nro_obra_social) references obra_social(nro_obra_social);
		alter table agenda add constraint fk_agenda foreign key (dni_medique) references medique(dni_medique);
		alter table agenda add constraint fk_agenda2 foreign key (nro_consultorio) references consultorio(nro_consultorio);
		alter table turno add constraint fk_turno foreign key (dni_medique) references medique(dni_medique);
		alter table turno add constraint fk_turno2 foreign key (nro_consultorio) references consultorio(nro_consultorio);
		alter table turno add constraint fk_turno3 foreign key (nro_paciente) references paciente(nro_paciente);
		alter table turno add constraint fk_turno4 foreign key (nro_obra_social_consulta) references obra_social(nro_obra_social);
		alter table reprogramacion add constraint fk_reprogramacion foreign key (nro_turno) references turno(nro_turno);
		alter table error add constraint fk_error foreign key (dni_medique) references medique(dni_medique);
		alter table error add constraint fk_error2 foreign key (nro_consultorio) references consultorio(nro_consultorio);
		alter table error add constraint fk_error3 foreign key (nro_paciente) references paciente(nro_paciente);
		alter table cobertura add constraint fk_cobertura foreign key (dni_medique) references medique(dni_medique);
		alter table cobertura add constraint fk_cobertura2 foreign key (nro_obra_social) references obra_social(nro_obra_social);
		alter table liquidacion_cabecera add constraint fk_liq_cab foreign key (nro_obra_social) references obra_social(nro_obra_social);
		alter table liquidacion_detalle add constraint fk_liq_det foreign key (nro_liquidacion) references liquidacion_cabecera(nro_liquidacion);
		alter table liquidacion_detalle add constraint fk_liq_det2 foreign key (dni_medique) references medique(dni_medique);`
	)
	
	if err != nil{
		log.Fatal(err)
	}
}


func cargarTablas() {
	db, err := dbConnection()
	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()
	
	_, err = db.Exec (`insert into medique (dni_medique, nombre, apellido, especialidad, monto_consulta_privada, telefono)
		values
		(10456789, 'Dr. Juan', 'Gomez', 'Cardiología', 15000.00, '0112533-1234'),
		(40342233, 'Dra. Maria', 'Lopez', 'Dermatología', 12000.00, '0113453-5678'),
		(56565656, 'Dr. Carlos', 'Rodriguez', 'Gastroenterología', 18000.00, '0112343-9012'),
		(11442233, 'Dra. Laura', 'Martinez', 'Neurología', 20000.00, '0118987-3456'),
		(23959693, 'Dra. Ana', 'Perez', 'Oftalmología', 16000.00, '0118489-7890'),
		(30506070, 'Dr. Javier', 'Fernandez', 'Pediatría', 13000.00, '0117347-2345'),
		(12094587, 'Dra. Sofia', 'Diaz', 'Psiquiatría', 17000.00, '011555-6789'),
		(23233412, 'Dr. Manuel', 'Garcia', 'Ortopedia', 14000.00, '0118383-0123'),
		(11990922, 'Dra. Marta', 'Sanchez', 'Oncología', 19000.00, '0118374-4567'),
		(99349349, 'Dr. Alejandro', 'Torres', 'Endocrinología', 16000.00, '0119383-8901'),
		(87654321, 'Dra. Patricia', 'Ramirez', 'Urología', 18000.00, '0119383-2345'),
		(11122334, 'Dr. Daniel', 'Gutierrez', 'Ginecología', 15000.00, '0119373-6789'),
		(23545427, 'Dra. Paula', 'Vargas', 'Reumatología', 17000.00, '0118333-0123'),
		(47351623, 'Dr. Sergio', 'Hernandez', 'Nefrología', 20000.00, '0115348-4567'),
		(41553273, 'Dra. Carolina', 'Flores', 'Cardiología', 19000.00, '0110013-8901'),
		(35213523, 'Dr. Luis', 'Cabrera', 'Dermatología', 15000.00, '0118364-2345'),
		(23125673, 'Dra. Silvia', 'Rojas', 'Gastroenterología', 12000.00, '0113673-6789'),
		(12536732, 'Dr. Gonzalo', 'Luna', 'Neurología', 18000.00, '0118743-0123'),
		(14353515, 'Dra. Ana', 'Mendez', 'Oftalmología', 20000.00, '0118738-4567'),
		(12459135, 'Dra. Alberto', 'Fuentes', 'Pediatría', 16000.00, '0119374-8901');
		
		insert into paciente values (1, 'Martin', 'Galvarini', 42660991, '2000-06-30', 1, 1000, 'Carlos Pellegrini 2436, Martinez', '0114416-3214', '13.martingalva@gmail.com');
		insert into paciente values (2, 'Pascual', 'Galvarini', 60123321, '2015-01-15', 2, 2000, 'Carlos Pellegrini 2436, Martinez', '0113395-2194', 'pascualgalva@gmail.com');
		insert into paciente values (3, 'Lucas', 'Manfredi', 43021777, '2000-11-02', 3, 3000, 'Domingo de Acassuso 150, La Lucila' , '0117483-2745', 'luquita@gmail.com');
		insert into paciente values (4, 'Juan Ignacio', 'Mussino', 41345123, '1999-07-23', 4, 4000, 'Centenario 100, San Isidro', '0115467-1234', 'juanimu@aol.com');
		insert into paciente values (5, 'Lorenzo', 'Paparo', 44736192, '2004-04-10', 5, 5000, 'Calle Cerrada 500, Villa Adelina', '0118432-8326', 'lorepapa@yahoo.com');
		insert into paciente values (6, 'Gianluca', 'Zeolla', 45932721, '2004-10-19', 6, 6000, 'Av.Peron 4303, Benavidez', '0115674-2341', 'gianluze14@gmail.com');
		insert into paciente values (7, 'Juan', 'Pérez', 12345678, '1986-01-01', 1, 5678, 'Guido Spano 1133, Victoria', '0111234-5678', 'juanperez@yahoo.com');
		insert into paciente values (8, 'María', 'González', 23456789, '1985-02-02', 2, 6789, 'Almafuerte 2835, La Matanza', '0112345-6789', 'mariagonzalez@hotmail.com');
		insert into paciente values (9, 'Pedro', 'Rodríguez', 38567890, '1989-03-03', 3, 7890, 'Av. Cazon 332, Tigre', '0113456-7890', 'pedrorodriguez@gmail.com');
		insert into paciente values (10, 'Lucía', 'Fernández', 45678901, '1995-04-04', 4, 8901, 'Roma 659, Olivos', '0114567-8901', 'luciafernandez@gmail.com');
		insert into paciente values (11, 'Jorge', 'Gómez', 56789012, '2000-05-05', 5, 9012, 'Jose Marmol 2590, Florida', '0115678-9012', 'jorgegomez@gmial.com');
		insert into paciente values (12, 'Ana', 'Díaz', 67890123, '2005-06-06', 6, 1234, 'Padilla 620, Villa Crespo', '0116789-1234', 'anadiaz@yahoo.com');
		insert into paciente values (13, 'Diego', 'Martínez', 78901234, '2010-07-07', 3, 2345, 'Thompson 565, Caballito', '0117890-2345', 'diegomartinez@gmail.com');
		insert into paciente values (14, 'Carla', 'Pérez', 60012345, '2015-08-08', 4, 3456, 'Av. Cabildo 2188, Belgrano', '0118901-3456', 'carlaperez@hotmail.com');
		insert into paciente values (15, 'Lucas', 'González', 90123456, '2020-09-09', 2, 4567, 'Darwin 1633, Palermo', '0119012-4567', 'lucasgonzalez@aol.com');
		insert into paciente values (16, 'Sofía', 'Rodríguez', 32564990, '1980-05-10', 1, 5778, 'Av. del Mar 649, Pinamar', '0111234-4312', 'sofiarodriguez@yahoo.com');
		insert into paciente values (17, 'Carlos', 'Bianchi', 1243765, '1949-04-26', 1, 17000, 'Brandsen 798, La Boca', '0114637-7584', 'virrey@hotmail.com');
		insert into paciente values (18, 'Manuel', 'Belgrano', 30987217, '1978-08-15', 4, 18000, 'Santa Fe 1812, Rosario', '0114531-1234', 'manubel@gmail.com');
		insert into paciente values (19, 'Lionel', 'Messi', 35094577, '1987-06-24', 3, 10, 'La Pampa 2133, Belgrano', '0114126-7789', 'lio.d10s@yahoo.com');
		insert into paciente values (20, 'Carlos Saul', 'Menem', 11365578, '1930-07-02', 4, 1312, 'La Rosadita 347, La Rioja', '0119012-1243', 'charly.menem10@outlook.com');

		insert into consultorio (nro_consultorio, nombre, domicilio, codigo_postal, telefono)
		values
		(1, 'Rene Favaloro', 'Fleming 2000', '1640', '11-5431-2311'),
		(2, 'Alexander Fleminf', 'Fleming 2020', '1640', '11-5411-2341'),
		(3, 'Edward Jenner', 'Fleming 2050', '1640', '11-5411-2351'),
		(4, 'William Osler', 'Fleming 2080', '1640', '11-5411-4311'),
		(5, 'Louis Pasteur', 'Fleming 2100', '1640', '11-5411-1234'),
		(6, 'Sigmund Freud', 'Andres Rolon 10', '1642', '11-5412-1235'),
		(7, 'Elizabeth Blackwell', 'Andres Rolon 40', '1642', '11-4311-4212'),
		(8, 'Joseph Lister', 'Andres Rolon 70', '1642','11-4312-1215'),
		(9, 'John Snow', 'Andres Rolon 130', '1642', '11-4431-4212'),
		(10, 'Hipocrates','Andres Rolon 162','1642','11-4312-1123');
		
		insert into obra_social (nro_obra_social, nombre, contacto_nombre, contacto_apellido, contacto_telefono, contacto_email) 
		values 
		(1, 'Galeno', 'Juan', 'Galeno', '4798-2345', 'galeno@gmail.com'),
		(2, 'Swiss Medical', 'Maria', 'Suazo', '4545-6532', 'swissmedical@gmail.com'),
		(3, 'OSDE', 'Orlando', 'Debarro', '4891-3214', 'osde@gmail.com'),
		(4, 'Omint', 'Omar', 'Lopez', '4671-5634', 'omint@gmail.com'),
		(5, 'Medicus', 'Mercedes', 'Costa', '4761-8799', 'medicus@gmail.com'),
		(6, 'Sancor Seguros', 'Sandra', 'Corleone', '4531-1927', 'sancor@gmail.com');
		
		insert into cobertura(dni_medique, nro_obra_social, monto_paciente, monto_obra_social)
		values
		(10456789,1, 11000.00, 21000.00),
		(40342233,2, 9000.00, 19000.00),
		(56565656,3, 15000.00, 25000.00),
		(11442233,4, 14000.00, 24000.00),
		(23959693,5, 10000.00, 20000.00),
		(30506070,6, 8000.00, 18000.00),
		(12094587,1, 22000.00, 30000.00),
		(23233412,1, 12000.00, 22000.00),
		(11990922,2, 22000.00, 27000.00),
		(99349349,3, 13000.00, 23000.00),
		(87654321,4, 15000.00, 24000.00),
		(11122334,5, 19000.00, 26000.00),
		(23545427,6, 17000.00, 25000.00),
		(47351623,2, 11000.00, 21000.00),
		(41553273,1, 21000.00, 25000.00),
		(35213523,2, 19000.00, 22000.00),
		(23125673,3, 15000.00, 24000.00),
		(12536732,4, 21000.00, 22000.00),
		(14353515,5, 11000.00, 18000.00),
		(12459135,6, 15000.00, 20000.00),
		(10456789,2, 17000.00, 24000.00),
		(40342233,3, 13000.00, 22000.00);
		
		insert into agenda (dni_medique, dia, nro_consultorio, hora_desde, hora_hasta, duracion_turno)
		values
		(10456789,1,1,'7:00','15:00','8 hour'),
		(40342233,2,2,'10:00','14:00','4 hour'),
		(56565656,3,3,'7:00','15:00','8 hour'),
		(11442233,4,4,'7:00','16:00','9 hour'),
		(23959693,5,5,'14:00','22:00','8 hour'),
		(30506070,6,6,'8:00','16:00','8 hour'),
		(12094587,7,7,'8:00','16:00','8 hour'),
		(23233412,1,8, '7:00','16:00', '9 hour'),
		(11990922,2,9, '14:00','22:00', '8 hour'),
		(99349349,3,10, '10:00','14:00', '4 hour'),
		(87654321,4,1, '15:00','19:00', '4 hour'),
		(11122334,5,2, '14:00','20:00', '6 hour'),
		(23545427,6,3, '15:00', '19:00', '4 hour'),
		(47351623,7,4, '16:00', '22:00', '6 hour'),
		(41553273,1,5, '8:00', '14:00', '6 hour'),
		(35213523,2,6, '16:00', '20:00', '4 hour'),
		(23125673,3,7, '16:00', '22:00', '6 hour'),
		(12536732,4,8, '16:00', '22:00', '6 hour'),
		(14353515,5,9, '8:00', '14:00', '6 hour'),
		(12459135,6,10, '14:00' , '20:00', '6 hour'),
		(10456789,2,1, '19:00', '23:00', '4 hour'),
		(40342233,3,2, '6:00', '10:00' , '4 hour'),
		(56565656,4,3, '19:00', '23:00', '4 hour' );`
		
		
	)
	
	if err != nil{
		log.Fatal(err)
	}
}
	


func crearSP() {
	
	}

func crearTriggers() {
	
	}

func eliminarPK() {
	
	}

func eliminarFK() {
	
	}

