package main

import (
    "fmt"
    "time"
	"log"
	"database/sql"
	_"github.com/lib/pq"
)

func main() {	
	var opcion int
	
	for opcion != 14 {
		fmt.Println("Elegi una opcion:")
		fmt.Println("1 Crear BD")
		fmt.Println("2. Crear Tablas")
		fmt.Println("3. Crear PKs y FKs")
		fmt.Println("4. Cargar Tablas")
		fmt.Println("5. Eliminar PKs y FKs")
		fmt.Println("6. Crear sp y triggers")
		fmt.Println("7. Generar turnos por mes")
		fmt.Println("8. Reservar turno")
		fmt.Println("9. Cancelacion de turno")
		fmt.Println("10. Atencion turno")
		fmt.Println("11. Email recordatorio")
		fmt.Println("12. Email perdida de turno")
		fmt.Println("13. Generar liquidacion de obras sociales")
		fmt.Println("14. Salir")
		fmt.Print("Ingrese una opcion: ")
		fmt.Scanf("%d", &opcion)

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
				eliminarFk()
				eliminarPK()
			case opcion == 6:
				crearSP()
				crearTriggers()
			case opcion == 7:
				var anio int
				var mes int
				
				fmt.Print("Ingrese el año a generar turnos: ")
				fmt.Scanf("%d", &anio)
				fmt.Print("Ingrese el mes a generar turnos: ")
				fmt.Scanf("%d", &mes)
				
				generarTurnos(anio, mes)
			case opcion == 8:
				var (
					nro_paciente int
					dni_medique int
					fecha_turno string
					hora_turno string
				)
				fmt.Print("Ingrese el numero de historia medica del paciente: ")
				fmt.Scanf("%d", &nro_paciente)
				fmt.Print("Ingrese el DNI del medique: ")
				fmt.Scanf("%d", &dni_medique)
				fmt.Print("Ingresa una fecha para el turno (formato: yyyy-mm-dd): ")
				fmt.Scan(&fecha_turno)
				fmt.Print("Ingresa una hora para el turno (formato: HH:MM): ")
				fmt.Scan(&hora_turno)
				
				fecha, err := time.Parse("2006-01-02", fecha_turno)
				if err != nil {
					fmt.Println("Error al parsear la fecha: ", err)
					return
				}		
				
				hora, err := time.Parse("15:04", hora_turno)
				if err != nil {
					fmt.Println("Error al parsear la hora: ", err)
					return
				}
					
			    reservarTurno(nro_paciente, dni_medique, fecha, hora)
			case opcion == 9:
				var (
					dni_medique int
					f_desde string
					f_hasta string
				)
				
				fmt.Print("Ingrese el DNI del medique: ")
				fmt.Scanf("%d", &dni_medique)
				fmt.Print("Ingrese la fecha de inicio para cancelar (formato: yyyy-mm-dd): ")
				fmt.Scan(&f_desde)				
				fmt.Printf("Ingrese la fecha final para cancelar (formato: yyyy-mm-dd): ")
				fmt.Scan(&f_hasta)
						
				td, err := time.Parse("2006-01-02", f_desde)
				if err != nil {
					fmt.Println("Error al parsear la fecha desde: ", err)
					return
				}			
					
				th, err := time.Parse("2006-01-02", f_hasta)
				if err != nil {
					fmt.Println("Error al parsear la fecha hasta: ", err)
					return
				}		
								
			    cancelacionTurnos(dni_medique, td, th)
			case opcion == 10:
				var nro_turno int
				fmt.Print("Ingrese el numero del turno: ")
				fmt.Scanf("%d", &nro_turno)
			
			    atencionTurnos(nro_turno)
			case opcion == 11:			
			    emailRecordatorio()
			case opcion == 12:
			    emailPerdidaTurno()
			case opcion == 13:
				var (
					nro_obra_social int
					f_desde string
					f_hasta string
				)
				
				fmt.Print("Ingrese el numero de obra social: ")
				fmt.Scanf("%d", &nro_obra_social)				
				fmt.Print("Ingrese la fecha de inicio para liquidar (formato: yyyy-mm-dd): ")
				fmt.Scan(&f_desde)				
				fmt.Printf("Ingrese la fecha final para liquidar (formato: yyyy-mm-dd): ")
				fmt.Scan(&f_hasta)
						
				td, err := time.Parse("2006-01-02", f_desde)
				if err != nil {
					fmt.Println("Error al parsear la fecha inicial:", err)
					return
				}			
					
				th, err := time.Parse("2006-01-02", f_hasta)
				if err != nil {
					fmt.Println("Error al parsear la fecha final:", err)
					return
				}		
				
			    generarLiquidacionObrasSociales(nro_obra_social, td, th)                 	
			case opcion == 14:
				fmt.Println("Adios!")
			default:
				fmt.Println("La opciòn ingresada no es vàlida, por favor ingrese ingrese otro numero.")
			}
		}
	fmt.Println("El programa se finalizo correctamente.")
}


func dbConnection()(*sql.DB, error){
    db, err := sql.Open("postgres", "user=postgres host=localhost dbname=hospital sslmode=disable")
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
	_, err = db.Exec("drop database if exists hospital")
	if err != nil {
		log.Fatal(err)
	}
	
	_, err = db.Exec("create database hospital")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Base de datos creada exitosamente.")
}

// crear tablas
func createTables() {
	db, err := dbConnection()

	_, err = db.Exec("drop schema public cascade")

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("create schema public")

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		create table paciente (nro_paciente int, nombre  text, apellido text, dni_paciente int, f_nacimiento date, nro_obra_social int, nro_afiliade int, domicilio text, telefono char(12),  email text); 
		create table medique (dni_medique int, nombre text, apellido text, especialidad varchar(64), monto_consulta_privada decimal(12, 2), telefono char(12)); 
		create table consultorio (nro_consultorio int, nombre text, domicilio text, codigo_postal char(8), telefono char(12));
		create table agenda (dni_medique int, dia int, nro_consultorio int, hora_desde time, hora_hasta time, duracion_turno interval);
		create table turno (nro_turno serial, fecha timestamp, nro_consultorio int, dni_medique int, nro_paciente int, nro_obra_social_consulta int, nro_afiliade_consulta int, monto_paciente decimal(12,2), monto_obra_social decimal(12,2), f_reserva timestamp, estado char(10));
		create table reprogramacion (nro_turno int, nombre_paciente text, apellido_paciente text, telefono_paciente char(12), email_paciente text, nombre_medique text, apellido_medique text, estado char(12));
		create table error (nro_error serial, f_turno timestamp, nro_consultorio int, dni_medique int, nro_paciente int, operacion char(12), f_error timestamp, motivo varchar(64));
		create table cobertura (dni_medique int, nro_obra_social int, monto_paciente decimal(12,2), monto_obra_social decimal(12,2));
		create table obra_social (nro_obra_social int, nombre text, contacto_nombre text, contacto_apellido text, contacto_telefono char(12), contacto_email text);
		create table liquidacion_cabecera (nro_liquidacion serial, nro_obra_social int, desde date, hasta date, total decimal(15,2));
		create table liquidacion_detalle (nro_liquidacion serial, nro_linea serial, f_atencion date, nro_afiliade int, dni_paciente int, nombre_paciente text, apellido_paciente text, dni_medique int, nombre_medique text, apellido_medique text, especialidad varchar(64), monto decimal(12,2));
		create table envio_email (nro_email serial, f_generacion timestamp, email_paciente text, asunto text, cuerpo text, f_envio timestamp, estado char(10));
		create table solicitud_reservas (nro_orden int, nro_paciente int, dni_medique int, fecha date, hora time);
		`)
		
	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()
	
	fmt.Println("Tablas creadas exitosamente.")
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
		alter table envio_email add constraint envio_email_pk primary key (nro_email);
		`)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println("Pk creadas exitosamente.")
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
		alter table liquidacion_detalle add constraint fk_liq_det2 foreign key (dni_medique) references medique(dni_medique);
		`)
	
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println("Fk creadas exitosamente.")
}


func cargarTablas() {
	db, err := dbConnection()
	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()
	
	_, err = db.Exec (`
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
		
		insert into medique (dni_medique, nombre, apellido, especialidad, monto_consulta_privada, telefono)
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
		
		insert into paciente values (1, 'Martin', 'Galvarini', 42660991, '2000-06-30', null, null, 'Carlos Pellegrini 2436, Martinez', '0114416-3214', '13.martingalva@gmail.com');
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
		(10456789,1,10,'08:00','10:40','40 minutes'),
		(40342233,2,2,'10:00','14:00','15 minutes'),
		(56565656,3,3,'7:00','15:00','25 minutes'),
		(11442233,4,4,'7:00','16:00','1 hour'),
		(23959693,5,5,'14:00','22:00','30 minutes'),
		(30506070,6,6,'8:00','16:00','50 minutes'),
		(12094587,7,7,'8:00','16:00','45 minutes'),
		(23233412,1,8, '7:00','16:00', '1 hour'),
		(11990922,2,9, '14:00','22:00', '35 minutes'),
		(99349349,3,10, '10:00','14:00', '40 minutes'),
		(87654321,4,1, '15:00','19:00', '45 minutes'),
		(11122334,5,2, '14:00','20:00', '1 hour'),
		(23545427,6,3, '15:00', '19:00', '40 minutes'),
		(47351623,7,4, '16:00', '22:00', '35 minutes'),
		(41553273,1,5, '8:00', '14:00', '15 minutes'),
		(35213523,2,6, '16:00', '20:00', '20 minutes'),
		(23125673,3,7, '16:00', '22:00', '30 minutes'),
		(12536732,4,8, '16:00', '22:00', '25 minutes'),
		(14353515,5,9, '8:00', '14:00', '1 hour'),
		(12459135,6,10, '14:00' , '20:00', '55 minutes'),
		(10456789,2,1, '19:00', '23:00', '45 minutes'),
		(40342233,3,2, '6:00', '10:00' , '45 minutes'),
		(56565656,4,3, '19:00', '23:00', '45 minutes' );
		`)
	
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println("Tablas cargadas")
}
	
func crearSP() {
    sp_generarTurnos()
	sp_reservarTurno()
	sp_cancelarTurno()
	sp_atencionTurno()
    sp_emailRecordatorioSP()
    sp_emailPerdidaSP()
	sp_liquidacionObrasSociales()
	fmt.Println("Stored procedures generadas.")
	}

func crearTriggers() {
	emailsTrigger()
	fmt.Println("Triggers generadas.")
}

func eliminarFk() {
	db, err := dbConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	_,err = db.Exec(`
		alter table paciente drop constraint fk_paciente restrict;
		alter table agenda drop constraint fk_agenda restrict;
		alter table agenda drop constraint fk_agenda2 restrict;
		alter table turno drop constraint fk_turno restrict;
		alter table turno drop constraint fk_turno2 restrict;
		alter table turno drop constraint fk_turno3 restrict;
		alter table turno drop constraint fk_turno4 restrict;
		alter table reprogramacion drop constraint fk_reprogramacion restrict;
		alter table error drop constraint fk_error restrict;
		alter table error drop constraint fk_error2 restrict;
		alter table error drop constraint fk_error3 restrict;
		alter table cobertura drop constraint fk_cobertura restrict;
		alter table cobertura drop constraint fk_cobertura2 restrict;
		alter table liquidacion_cabecera drop constraint fk_liq_cab restrict;
		alter table liquidacion_detalle drop constraint fk_liq_det restrict;
		alter table liquidacion_detalle drop constraint fk_liq_det2 restrict;`)

		if err != nil{
			log.Fatal(err)
		}
		fmt.Println("Fk eliminadas.")
}

func eliminarPK() {
	db, err := dbConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	_,err = db.Exec(`
		alter table paciente drop constraint paciente_pk restrict;
		alter table medique drop constraint medique_pk restrict;
		alter table consultorio drop constraint consultorio_pk restrict;
		alter table agenda drop constraint agenda_pk restrict;
		alter table turno drop constraint turno_pk restrict;
		alter table reprogramacion drop constraint reprogramacion_pk restrict;
		alter table error drop constraint error_pk restrict;
		alter table cobertura drop constraint cobertura_pk restrict;
		alter table obra_social drop constraint obra_social_pk restrict;
		alter table liquidacion_cabecera drop constraint liquidacion_cabecera_pk restrict;
		alter table liquidacion_detalle drop constraint liquidacion_detalle_pk restrict;
		alter table envio_email drop constraint envio_email_pk restrict;
		`)

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Pk eliminadas.")

}

//stored Procedures
func sp_generarTurnos() {
	db, err := dbConnection()

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		create or replace function generar_turnos(_anio int, _mes int) returns bool as $$
		declare
				aux record;
				_fecha date;
				_agenda record;
				fecha_completa timestamp;
		begin

			--verificar si los turnos ya existen
			select * from turno t where extract(year from t.fecha) = _anio and extract(month from t.fecha) = _mes into aux;
			
			if found then
				raise notice 'Los turnos para esas fechas ya fueron generados';
				return false;
			end if;
			
			for _agenda in select * from agenda a loop
				_fecha = to_date(_anio || '-' || _mes || '-' || '01', 'yyyy-mm-dd');
								
				while (extract(month from _fecha))::int = _mes loop
					if extract(isodow from _fecha)::int = _agenda.dia then
						fecha_completa = _fecha + _agenda.hora_desde;
				
						while fecha_completa::time < _agenda.hora_hasta loop
							insert into turno (fecha, nro_consultorio, dni_medique, estado)
							select fecha_completa, _agenda.nro_consultorio, _agenda.dni_medique, 'disponible'
							from agenda
							where agenda.dni_medique = _agenda.dni_medique and extract(isodow from fecha_completa) = agenda.dia;
							
							fecha_completa = fecha_completa + _agenda.duracion_turno;

						end loop;
					end if;			
					
					_fecha = _fecha + 1;
				end loop;
			end loop;
			raise notice 'Los turnos fueron generados';
			return true;

		end;
		$$ language plpgsql;

	`)
}

func sp_reservarTurno() {

     db, err := dbConnection()
	 if err != nil {
		 log.Fatal(err)
	 }
	 defer db.Close()

	_, err = db.Exec(`
		create or replace function reservar_turno(_nro_paciente int, _dni_medique int, _fecha_turno date, _hora_turno time) returns boolean as $$

		declare
			obrasocial_datos record;
			paciente_datos record;
			cantidad_turnos_reservados int;
			condicion boolean;
			fecha_completa timestamp;
			aux record;
			datos_ob record;

		begin
			
			-- verificar su el dni del medique existe
			select m.dni_medique from medique m where m.dni_medique = _dni_medique into aux;
			
			if not found then
				insert into error (f_turno, nro_consultorio, dni_medique, nro_paciente, operacion, f_error, motivo)
				values (now(), null, null, null, 'reserva', now(), 'dni de medique no valido');
				raise notice 'dni de medique no valido';
				return false;
			end if;
			

			-- verificar si el numero de historia clinica existe
			select p.nro_paciente from paciente p where p.nro_paciente = _nro_paciente into aux;
			
			if not found then
				insert into error (f_turno, nro_consultorio, dni_medique, nro_paciente, operacion, f_error, motivo)
				values (now(), null, null, null, 'reserva', now(), 'nro de historia clinica no valido');
				raise notice 'nro de historia clinica no valido';
				return false;
			end if;
			
			-- verificar si el paciente tiene una obra social y obtener la obra social
			condicion := false;
			select p.nro_obra_social from paciente p, obra_social o where p.nro_obra_social = o.nro_obra_social and p.nro_paciente = _nro_paciente into paciente_datos;
			
			if found then
				condicion := true;
			end if;
			
			-- verificar si el medique atiende esa obra social
			case
				when condicion then
					select * from cobertura c, obra_social o where c.dni_medique = _dni_medique and c.nro_obra_social = paciente_datos.nro_obra_social into aux;
					
					if not found then
						insert into error (f_turno, nro_consultorio, dni_medique, nro_paciente, operacion, f_error, motivo)
						values (now(), null, null, null, 'reserva', now(), 'obra social no atendida por el medique');
						raise notice 'obra social no atendida por el medique';
						return false;
					end if;
					
					select c.monto_paciente, c.monto_obra_social, p.nro_afiliade from medique m, cobertura c, paciente p where p.nro_obra_social = c.nro_obra_social and m.dni_medique = c.dni_medique and _dni_medique = m.dni_medique and _nro_paciente = p.nro_paciente into datos_ob;
				
				else
					select m.monto_consulta_privada from medique m, paciente p where _dni_medique = m.dni_medique into datos_ob;
					
			end case;
			
			-- verificar si el turno esta disponible
			fecha_completa = _fecha_turno + _hora_turno;
			
			select nro_turno from turno t where t.fecha = fecha_completa and t.dni_medique = _dni_medique and estado = 'disponible' into aux;
			
			if not found then
				insert into error (f_turno, nro_consultorio, dni_medique, nro_paciente, operacion, f_error, motivo)
				values (now(), null, null, null, 'reserva', now(), 'turno inexistente o no disponible');
				raise notice 'turno inexistente o no disponible';
				return false;
			end if;
			
			-- verificar si el paciente ha superado el limite de 5 turnos en estado reservado
			select count(*) from turno where nro_paciente = _nro_paciente and estado = 'Reservado' into cantidad_turnos_reservados ;

			if cantidad_turnos_reservados > 5 then
				insert into error (f_turno, nro_consultorio, dni_medique, nro_paciente, operacion, f_error, motivo)
				values (now(), null, null, null, 'reserva', now(), 'Supera el lìmite de reserva de turnos');
				raise notice 'Supera el lìmite de reserva de turnos';
				return false;
			end if;
			
			
			-- realizar la reserva del turno
			if condicion then
				update turno
				set
				nro_paciente = _nro_paciente,
				nro_obra_social_consulta = paciente_datos.nro_obra_social,
				nro_afiliade_consulta = datos_ob.nro_afiliade,
				estado = 'reservado',
				f_reserva = now(),
				monto_paciente = datos_ob.monto_paciente,
				monto_obra_social = datos_ob.monto_obra_social
				where
				turno.dni_medique = _dni_medique
				and turno.fecha = fecha_completa;
			else
				update turno
				set
				nro_paciente = _nro_paciente,
				estado = 'reservado',
				f_reserva = now(),
				monto_paciente = datos_ob.monto_consulta_privada
				where
				turno.dni_medique = _dni_medique
				and turno.fecha = fecha_completa;
			end if;			
			
			return true;
		end;
		$$ language plpgsql;

	`)
}

func sp_cancelarTurno() {
	db, err := dbConnection()

	 if err != nil {

		 log.Fatal(err)
	 }
	 defer db.Close()

	_, err = db.Exec(`
		create or replace function cancelacion_turnos(_dni_medique int, _fdesde timestamp, _fhasta timestamp) returns int as $$
		declare
			aux record;
			reprog record;
			reprog2 record;
			cantidad_turnos_cancelados int;
			resultado int;
		
		begin
			
			select m.apellido, m.nombre from medique m into aux;
			
			select t.nro_turno, p.nombre, p.apellido, p.telefono, p.email from turno t, paciente p, medique m into reprog 
			where m.dni_medique = t.dni_medique and p.nro_paciente = t.nro_paciente and t.dni_medique = _dni_medique and t.fecha >= _fdesde and t.fecha <= _fhasta and estado = 'reservado' ;
			
			if found then
				insert into reprogramacion (nro_turno, nombre_paciente, apellido_paciente, telefono_paciente, email_paciente, nombre_medique, apellido_medique, estado) values (reprog.nro_turno, reprog.nombre, reprog.apellido, reprog.telefono, reprog.email, aux.nombre, aux.apellido, 'pendiente');

				update turno
				set estado = 'cancelado' where dni_medique = _dni_medique and fecha >= _fdesde and fecha <= _fhasta and estado ='reservado';
			end if;
			
			select t.nro_turno, m.nombre, m.apellido from turno t, medique m into reprog2
			where m.dni_medique = t.dni_medique and t.dni_medique = _dni_medique and t.fecha >= _fdesde and t.fecha <= _fhasta and estado = 'disponible' ;
			
			if found then
				update turno
				set estado = 'cancelado' where dni_medique = _dni_medique and fecha >= _fdesde and fecha <= _fhasta and estado ='disponible';

			end if;
			
			resultado = 1; --despues ver lo del row_count para contar la cantidad de turnos cancelados
			
			return resultado;
		end;
		$$ language plpgsql;


	`)
}

func sp_atencionTurno() {
	db, err := dbConnection()

	 if err != nil {

		 log.Fatal(err)
	 }
	 defer db.Close()

	_, err = db.Exec(`
		create or replace function atencion_turnos(_nro_turno int) returns bool as $$
		
		declare
			turno_atendido record;
			
		begin
			select t.nro_turno from turno t where t.nro_turno = _nro_turno into turno_atendido;
			
			if not found then
				insert into error (f_turno, nro_consultorio, dni_medique, nro_paciente, operacion, f_error, motivo)
				values (now(), null, null, null, 'atencion', now(), 'nro de turno no valido');
				raise notice 'nro de turno no valido';
				return false;
			end if;
			
			
			select t.estado from turno t where t.nro_turno = _nro_turno and t.estado = 'reservado' into turno_atendido;
			
			if not found then
				insert into error (f_turno, nro_consultorio, dni_medique, nro_paciente, operacion, f_error, motivo)
				values (now(), null, null, null, 'atencion', now(), 'turno no reservado');
				raise notice 'turno no reservado';
				return false;
			end if;
			
			
			select t.fecha from turno t where t.nro_turno = _nro_turno and t.fecha::date = current_date into turno_atendido; --tengo que separar la fecha del timestamp t.fecha para compararlo con current_date?
			
			if not found then
				insert into error (f_turno, nro_consultorio, dni_medique, nro_paciente, operacion, f_error, motivo)
				values (now(), null, null, null, 'atencion', now(), 'turno no corresponde a la fecha del dia');
				raise notice 'turno no corresponde a la fecha del dia';
				return false;
			end if;
			
			update turno
			set estado = 'atendido' where nro_turno = _nro_turno;
			
			return true;
			
		end;
		$$ language plpgsql;


	`)
}

func sp_emailRecordatorioSP() {
	
	db, err := dbConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
	
		create or replace function email_recordatorio() returns void as $$
		declare
			datos_turno record;
			cuerpo_email text;
			fecha_email date;
			
		begin
			fecha_email := current_date + interval '2 days';
			select p.email, m.nombre, m.apellido, t.fecha from paciente p, turno t, medique m where t.nro_paciente = p.nro_paciente and t.dni_medique = m.dni_medique and t.fecha::date = fecha_email and t.estado = 'reservado' into datos_turno;
			
			if found then
			cuerpo_email := 'Recordatorio de turno del dia: ' || datos_turno.fecha ||', con el Dr. ' || datos_turno.nombre || ' ' || datos_turno.apellido;
			insert into envio_email (f_generacion, email_paciente, asunto, cuerpo, f_envio, estado)
			values (now(), datos_turno.email, 'Recordatorio de turno', cuerpo_email, now(), 'enviado');
			end if;

		end;
		$$ language plpgsql;
		
		`)	
}
	
func sp_emailPerdidaSP() {

	db, err := dbConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	_, err = db.Exec(`
	
	create or replace function email_perdida_turno() returns void as $$
	declare
		datos_turno record;
		cuerpo_email text;
		
	begin
		select p.email, m.nombre, m.apellido, t.fecha from paciente p, turno t, medique m where t.nro_paciente = p.nro_paciente and t.dni_medique = m.dni_medique and t.fecha <= now() and estado = 'reservado' into datos_turno;
		cuerpo_email := 'Perdiste tu turno del dia: ' || datos_turno.fecha ||', con el Dr. ' || datos_turno.nombre || ' ' || datos_turno.apellido;
		
		insert into envio_email (f_generacion, email_paciente, asunto, cuerpo, f_envio, estado)
		values (now(), datos_turno.email, 'Pérdida de turno reservado', cuerpo_email, now(), 'enviado');

	end;
	$$ language plpgsql;
		
	`)	

}

func sp_liquidacionObrasSociales() { 

    db, err :=dbConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
	create or replace function generar_liquidacion_obras_sociales(_nro_obra_social int, _desde DATE, _hasta DATE) returns void as $$
	declare 
		total_liquidacion float;
		_nro_liquidacion int;
		aux record;
	begin
		total_liquidacion := 0.0;
		
		
		
		insert into liquidacion_cabecera (nro_obra_social, desde, hasta, total)
		values(_nro_obra_social, _desde, _hasta, total_liquidacion)
		returning nro_liquidacion into _nro_liquidacion;

		update turno
		set estado = 'liquidado'
		where nro_obra_social_consulta = _nro_obra_social and estado = 'atendido'
		and fecha between _desde and _hasta;

		insert into liquidacion_detalle (nro_liquidacion, f_atencion, nro_afiliade, dni_paciente, nombre_paciente, apellido_paciente, dni_medique, nombre_medique, apellido_medique, especialidad, monto)
		select
		_nro_liquidacion, t.fecha, t.nro_afiliade_consulta, p.dni_paciente, p.nombre, p.apellido, t.dni_medique, m.nombre, m.apellido, m.especialidad, t.monto_obra_social
		from turno t, paciente p, medique m
		where t.nro_obra_social_consulta = _nro_obra_social  and t.nro_paciente = p.nro_paciente and t.dni_medique = m.dni_medique and t.estado = 'liquidado' and t.fecha between _desde and _hasta;

		select t.estado from turno t, paciente p, medique m
		where t.nro_obra_social_consulta = _nro_obra_social  and t.nro_paciente = p.nro_paciente and t.dni_medique = m.dni_medique and t.estado = 'liquidado' and t.fecha between _desde and _hasta into aux;

		select sum(monto) into total_liquidacion
		from liquidacion_detalle
		where nro_liquidacion = _nro_liquidacion;

		update liquidacion_cabecera
		set total = total_liquidacion
		where nro_liquidacion = _nro_liquidacion and nro_obra_social = _nro_obra_social and estado = 'liquidado';

	end;
	$$ language plpgsql;

	`)
	if err != nil {
		log.Fatal(err)
	}
}

//Triggers
func emailsTrigger() {

	db, err := dbConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`

create or replace function email_reserva() returns trigger as $$
	declare
		datos_turno record;
		cuerpo_email text;
		
	begin
		select p.email, m.nombre, m.apellido, t.fecha from paciente p, turno t, medique m where new.nro_turno = t.nro_turno and t.nro_paciente = p.nro_paciente and t.dni_medique = m.dni_medique into datos_turno;
		
		cuerpo_email := 'se reservo su turno para el dia: ' || datos_turno.fecha ||', con el Dr. ' || datos_turno.nombre || ' ' || datos_turno.apellido;
		
		insert into envio_email (f_generacion, email_paciente, asunto, cuerpo, f_envio, estado)
		values (now(), datos_turno.email, 'Reserva de turno', cuerpo_email, null, 'pendiente');
		
		return new;
	end;
	$$ language plpgsql;
	
	create trigger trigger_email_reserva
	after update on turno
	for each row
	execute function email_reserva();
	
------------------------------------------------------------------------
	
create or replace function email_cancelacion() returns trigger as $$
	declare
		datos_turno record;
		cuerpo_email text;
		
	begin
		select p.email, m.nombre, m.apellido, t.fecha from paciente p, medique m, turno t, reprogramacion r where new.nro_turno = r.nro_turno and t.nro_turno = r.nro_turno and t.nro_paciente = p.nro_paciente and t.dni_medique = m.dni_medique into datos_turno;		
		cuerpo_email := 'se cancelo su turno del dia: ' || datos_turno.fecha ||', con el Dr. ' || datos_turno.nombre || ' ' || datos_turno.apellido;
		
		insert into envio_email (f_generacion, email_paciente, asunto, cuerpo, f_envio, estado)
		values (now(), datos_turno.email, 'Reserva de turno', cuerpo_email, null, 'pendiente');
		
		return new;
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

//Ejecuciones de sp
func generarTurnos(_anio int, _mes int) {
	db, err := dbConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var resultado bool
	
	err = db.QueryRow(`
		select generar_turnos($1,$2);
	`, _anio, _mes).Scan(&resultado)
	if err != nil {
		log.Fatal(err)
	}	
	
	fmt.Println("Resultado de la generación: ", resultado)
}
	
func reservarTurno(_nro_paciente int, _dni_medique int, _fecha_turno time.Time, _hora_turno time.Time)  {
	db, err := dbConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var resultado bool

	err = db.QueryRow(`
		select reservar_turno($1, $2, $3,$4);
	`, _nro_paciente, _dni_medique, _fecha_turno, _hora_turno).Scan(&resultado)
	
	if err != nil {
		log.Fatal(err)
	}	
	
	fmt.Println("Resultado de la reserva: ", resultado)
}
	
func cancelacionTurnos(_dni_medique int, _fdesde time.Time, _fhasta time.Time)	{
	db, err := dbConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var turnos_cancelados int

	err = db.QueryRow(`
		select cancelacion_turnos($1, $2, $3);
	`, _dni_medique, _fdesde, _fhasta).Scan(&turnos_cancelados)
	
	if err != nil {
		log.Fatal(err)
	}	
	
	fmt.Println("Cantidad de turnos cancelados: ", turnos_cancelados)
}
func atencionTurnos(_nro_turno int) {
	db, err := dbConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var atendido bool

	err = db.QueryRow(`
		select atencion_turnos($1);
	`, _nro_turno).Scan(&atendido)
	if err != nil {
		log.Fatal(err)
	}	
	
	fmt.Println("Turno atendido: ", atendido)
}
	
func emailRecordatorio() {
	db, err := dbConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Query(`
		select email_recordatorio();
	`)	
}	

func emailPerdidaTurno() {
	db, err := dbConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Query(`
		select email_perdida_turno();
	`)	
}

func generarLiquidacionObrasSociales(_nro_obra_social int, _desde time.Time, _hasta time.Time)	{
	db, err := dbConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Query(`
		select generar_liquidacion_obras_sociales($1, $2, $3);
	`, _nro_obra_social, _desde, _hasta)	
	
}
