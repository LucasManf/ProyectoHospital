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
	
	var opcion int
	
	for opcion != 15 {
		fmt.println("Elegi una opcion:")
		fmt.println("1 Crear BD")
		fmt.println("2. Crear Tablas")
		fmt.println("3. Crear PKs y FKs")
		fmt.println("4. Cargar Tablas")
		fmt.println("5. Cargar base de datos no relacional")
		fmt.println("6. Mostrar base de datos no relacioanl")
		fmt.println("7. Eliminar PKs y FKs")
		fmt.println("8. Crear sp y triggers")
		fmt.println("9. Ejecutar sp")
		fmt.println("10. Ejecutar sp")
		fmt.println("11. Ejecutar sp")
		fmt.println("12. Ejecutar sp")
		fmt.println("13. Ejecutar sp")
		fmt.println("14. Ejecutar sp")
		fmt.println("15. Salir")
		
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
				cargarDatosJson()
			case opcion == 6:
				mostrarDatosJson()
			case opcion == 7:
				eliminarPK()
				eliminarFK()
			case opcion == 8:
				crearSP()
				crearTriggers()
			case opcion == 15:
				fmt.println("Adios!")
			default:
				fmt.println("La opciòn ingresada no es vàlida, por favor ingrese ingrese otro numero.")
		}
	}
	fmt.println("El programa se finalizo correctamente.")
	
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

	db, err := dbConnection()

	if err != nil {
		log.Fatal(err)

     }
	 defer db.Close()

	reservarTurno()
	cancelarTurno()
	atencionTurno()
    emails()
	
	}



func crearTriggers() {
	
	}

func eliminarPK() {
	
	}

func eliminarFK() {
	
	}

//stored Procedures
func reservarTurno() {

     db, err := dbConnection()

	 if err != nil {

		 log.Fatal(err)
	 }
	 defer db.Close()

	_, err = dbExec(`
		create or replace function reservar_turno(_nro_paciente int, _dni_medique int, _fecha_turno DATE, _hora_turno time) returns boolean as $$

		declare
		paciente_obrasocial int;
		medique_obrasocial int;
		cantidad_turnos_reservados int;

		begin
		// verificar si el paciente tiene una obra social
		select nro_obra_social into paciente_obrasocial
		from paciente
		where nro_paciente = _nro_paciente;

		// verificar su el dni del medique existe
		if not exists (select 1 from medique where dni_medique = _dni_medique) then
			raise exception 'Dni del medique no vàlido';
		end if;

		// verificar si el numero de historia clinica existe
		if not exists (select 1 from paciente where nro_paciente = _nro_paciente) then
			raise exception 'Nùmero de historia clìnica no vàlido';
		end if;

		// verificar si el paciente tiene una obra social y obtener la obra social
		if not exists (

			select 1
			from medique_obrasocial
			where dni_medique = _dni_medique and nro_obrasocial = paciente_obrasocial
		) then
		raise exception 'Obra social de paciente no atendida por el medique';
		end if;

		// obtener la obra social del medique
		select nro_obrasocial into medique_obrasocial
		from medique_obrasocial
		where dni_medique = _dni_medique

		// verificar si el turno esta disponible
		if not exists (

			select 1
			from turno
			where dni_medique = _dni_medique
			and fecha = _fecha_turno
			and hora = _hora_turno
			and estado = 'Disponible'
		) then
			raise exception 'Turno inexistente o no disponible'
		end if;

		// realizar la reserva del turno
		update turno
		set
		nro_paciente = _nro_paciente,
		nro_obra_social_consulta = paciente_obrasocial,
		estado = 'Reservado'
		where
		dni_medique = _dni_medique
		and fecha = _fecha_turno
		and hora = _hora_turno
		and estado = 'Disponible';
		return true;
		end;

		// verificar si el paciente ha superado el limite de 5 turnos en estado reservado
		select count(*) into cantidad_turnos_reservados
		from turno
		where nro_paciente = _nro_paciente and estado = 'Reservado';

		if cantidad_turnos_reservados >= 5 then
			raise exception 'Supera el lìmite de reserva de turnos';
		end if;

		$$ language plpgsql;
	`)

}

func cancelarTurno() {
	db, err := dbConnection()

	 if err != nil {

		 log.Fatal(err)
	 }
	 defer db.Close()

	_, err = dbExec(`
		create or replace function cancelacion_turnos(_dni_medique int, _fdesde timestamp, _fhasta timestamp) returns int as $$
		declare
			reprog record;
			cantidad_turnos_cancelados int;
		
		begin
			select t.nro_turno, p.nombre, p.apellido, p.telefono, p.email, m.nombre, m.apellido from turno t, paciente p, medique m into reprog 
			where m.dni_medique = t.dni_medique and p.nro_paciente = t.nro_paciente and t.dni_medique = _dni_medique and t.fecha >= _fdesde and t.fecha <= _fhasta and estado = 'reservado';

			update turno
			set estado = 'cancelado' where dni_medique = _dni_medique and fecha >= _fdesde and fecha <= _fhasta and estado ='reservado';
			
			if found then
				insert into reprogramacion (nro_turno, nombre_paciente, apellido_paciente, telefono_paciente, email_paciente, nombre_medique, apellido_medique, estado) values (reprog.nro_turno, reprog.p.nombre, reprog.p.apellido, reprog.p.telefono, reprog.p.email, reprog.m.nombre, reprog.m.apellido, 'pendiente');
				cantidad_turnos_cancelados := ROW_COUNT;
			end if;
		
			return cantidad_turnos_cancelados;
		end;
		$$ language plpgsql;
	`)
}

func atencionTurno() {
	db, err := dbConnection()

	 if err != nil {

		 log.Fatal(err)
	 }
	 defer db.Close()

	_, err = dbExec(`
		create or replace function atencion_turnos(_nro_turno int) returns bool as $$
		
		declare
			turno_atendido record;
			
		begin
			select t.nro_turno from turno t into turno_atendido where t.nro_turno = _nro_turno;
			
			if not found then
				insert into error (f_turno, nro_consultorio, dni_medique, nro_paciente, operacion, f_error, motivo)
				values (now(), null, null, null, 'atencion de turnos', now(), 'nro de turno no valido');
				return false;
			end if;
			
			
			select t.estado from turno t into turno_atendido where t.nro_turno = _nro_turno and t.estado = 'reservado';
			
			if not found then
				insert into error (f_turno, nro_consultorio, dni_medique, nro_paciente, operacion, f_error, motivo)
				values (now(), null, null, null, 'atencion de turnos', now(), 'turno no reservado');
				return false;
			end if;
			
			
			select t.fecha from turno t into turno_atendido where t.nro_turno = _nro_turno and t.fecha = current_date;
			
			if not found then
				insert into error (f_turno, nro_consultorio, dni_medique, nro_paciente, operacion, f_error, motivo)
				values (now(), null, null, null, 'atencion de turnos', now(), 'turno no corresponde a la fecha del dia');
				return false;
			end if;
			
			update turno
			set estado = 'atendido' where nro_turno = _nro_turno;
			
			return true;			
		end;
		$$ language plpgsql;

	`)
}


func emailsSP() {
	
	db, err := dbConnection()
	if err != nil {

		log.Fatal(err)
	}
	defer db.Close()

	_, err = dbExec(`
	
		create or replace function email_recordatorio() returns void as $$
		declare
			datos_turno record;
			cuerpo_email text;
			fecha_email date := current_date = interval '2 days';
			
		begin
			select p.email, m.nombre, m.apellido, t.fecha from paciente p, turno t, medique m where t.nro_paciente = p.nro_paciente and t.dni_medique = m.dni_medique and t.fecha = fecha_email into datos_turno;
			cuerpo_email := 'Recordatorio de turno del dia: ' || datos_turno.fecha ||', con el Dr. ' || datos_turno.nombre || ' ' || datos_turno.apellido;
			
			insert into envio_email (f_generacion, email_paciente, asunto, cuerpo, f_envio, estado)
			values (now(), datos_paciente.email, 'Reserva de turno', cuerpo_email, null, 'pendiente');

		end;
		$$ language plpgsql;
	
------------------------------------------------------------------------
	
	
		create or replace function email_perdida_turno() returns void as $$
			declare
				datos_turno record;
				cuerpo_email text;
				
			begin
				select p.email, m.nombre, m.apellido, t.fecha from paciente p, turno t, medique m where t.nro_paciente = p.nro_paciente and t.dni_medique = m.dni_medique and t.fecha <= now() into datos_turno;
				cuerpo_email := 'Perdiste tu turno del dia: ' || datos_turno.fecha ||', con el Dr. ' || datos_turno.nombre || ' ' || datos_turno.apellido;
				
				insert into envio_email (f_generacion, email_paciente, asunto, cuerpo, f_envio, estado)
				values (now(), datos_paciente.email, 'Reserva de turno', cuerpo_email, null, 'pendiente');

			end;
			$$ language plpgsql;
		
	`)	
}


//Triggers
func emailsTrigger() {

	db, err := dbConnection()
	if err != nil {

		log.Fatal(err)
	}
	defer db.Close()

	_, err = dbExec(`

	create or replace function email_reserva() returns trigger as $$
	declare
		datos_turno record;
		cuerpo_email text;
		
	begin
		select p.email, m.nombre, m.apellido, t.fecha from paciente p, turno t, medique m where new.nro_turno = t.nro_turno and t.nro_paciente = p.nro_paciente and t.dni_medique = m.dni_medique into datos_turno;
		
		cuerpo_email := 'se reservo su turno para el dia: ' || datos_turno.fecha ||', con el Dr. ' || datos_turno.nombre || ' ' || datos_turno.apellido;
		
		insert into envio_email (f_generacion, email_paciente, asunto, cuerpo, f_envio, estado)
		values (now(), datos_paciente.email, 'Reserva de turno', cuerpo_email, null, 'pendiente');

	end;
	$$ language plpgsql;
	
	create trigger trigger_email_reserva
	after insert on turno
	for each row
	execute function email_reserva();
	
------------------------------------------------------------------------
	
	create or replace function email_cancelacion() returns trigger as $$
	declare
		datos_turno record;
		cuerpo_email text;
		
	begin
		select p.email, m.nombre, m.apellido, t.fecha from paciente p, medique m, turno t, reprogramacion r where new.nro_turno = r.nro_turno and t.nro_turno = r.nro_turno and t.nro_paciente = p.nro_paciente and t.dni_medique = m.dni_medique into datos_paciente;		
		cuerpo_email := 'se cancelo su turno del dia: ' || datos_turno.fecha ||', con el Dr. ' || datos_turno.nombre || ' ' || datos_turno.apellido;
		
		insert into envio_email (f_generacion, email_paciente, asunto, cuerpo, f_envio, estado)
		values (now(), datos_paciente.email, 'Reserva de turno', cuerpo_email, null, 'pendiente');

	end;
	$$ language plpgsql;
	
	create trigger trigger_email_cancelacion
	after insert on reprogramacion
	for each row
	execute function email_cancelacion();
	

	

	 `)

	if err != nil {

		log.Fatal(err)

	}

}

