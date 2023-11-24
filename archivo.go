//conectar con bd
package main

import (
	"encoding/json"
    "fmt"
    "time"
    bolt "go.etcd.io/bbolt"
    "strconv"
	"log"
	"database/sql"
	_ "github.com/lib/pq"
)

func main() {
	_ = pq.QuoteIdentifier("some_text")
	
	var opcion int
	
	for opcion != 16 {
		fmt.Println("Elegi una opcion:")
		fmt.Println("1 Crear BD")
		fmt.Println("2. Crear Tablas")
		fmt.Println("3. Crear PKs y FKs")
		fmt.Println("4. Cargar Tablas")
		fmt.Println("5. Cargar base de datos no relacional")
		fmt.Println("6. Mostrar base de datos no relacional")
		fmt.Println("7. Eliminar PKs y FKs")
		fmt.Println("8. Crear sp y triggers")
		fmt.Println("9. Generar turnos por mes")
		fmt.Println("10. Reservar turno")
		fmt.Println("11. Cancelacion de turno")
		fmt.Println("12. Atencion turno")
		fmt.Println("13. Email recordatorio")
		fmt.Println("14. Email perdida de turno")
		fmt.Println("15. Generar liquidacion de obras sociales")
		fmt.Println("16. Salir")
		
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
				eliminarFk()
				eliminarPK()
			case opcion == 8:
				crearSP()
				crearTriggers()
			case opcion == 9:
				var anio int
				var mes int
				
				fmt.Print("Ingrese el año a generar turnos: ")
				fmt.Scanf("%d", &anio)
				fmt.Print("Ingrese el mes a generar turnos: ")
				fmt.Scanf("%d", &mes)
				
				generarTurnos(int anio, int mes)
			case opcion == 10:
				var (
					nro_paciente int
					dni_medique int
					fecha_turno string
				)
				fmt.Print("Ingrese el numero de historia medica del paciente: ")
				fmt.Scanf("%d", &nro_paciente)
				fmt.Print("Ingrese el DNI del medique: ")
				fmt.Scanf("%d", &dni_medique)
				fmt.Print("Ingresa una fecha y hora para el turno (formato: yyyy-mm-dd HH:MM:SS): ")
				fmt.Scanf("%s", &fecha_turno)	
				
				t, err := time.Parse("2006-01-02 15:04:05", fecha_turno)
				if err != nil {
					fmt.Println("Error al parsear la fecha y hora:", err)
					return
				}			
				
			    reservarTurno(nro_paciente, dni_medique, fecha_turno)
			case opcion == 11:
				var (
					dni_medique int
					f_desde string
					f_hasta string
				)
				
				fmt.Print("Ingrese el DNI del medique: ")
				fmt.Scanf("%d", &dni_medique)
				fmt.Print("Ingrese la fecha de inicio para cancelar (formato: yyyy-mm-dd HH:MM:SS): ")
				fmt.Scanf("%s", &f_desde)
				fmt.Printf("Ingrese la fecha final para cancelar (formato: yyyy-mm-dd HH:MM:SS): ")
				fmt.Scanf("%s", &f_hasta)
						
				t, err := time.Parse("2006-01-02 15:04:05", f_desde)
				if err != nil {
					fmt.Println("Error al parsear la fecha y hora:", err)
					return
				}			
					
				t, err := time.Parse("2006-01-02 15:04:05", f_hasta)
				if err != nil {
					fmt.Println("Error al parsear la fecha y hora:", err)
					return
				}		
								
			    cancelacionTurnos(dni_medique, f_desde, f_hasta)
			case opcion == 12:
				var nro_turno int
				fmt.Print("Ingrese el numero del turno: ")
				fmt.Scanf("%d", &nro_turno)
			
			    atencionTurnos(nro_turno)
			case opcion == 13:			
			    emailRecordatorio()
			case opcion == 14:
			    emailPerdidaTurno()
			case opcion == 15:
				var (
					nro_obra_social int
					f_desde string
					f_hasta string
				)
				
				fmt.Print("Ingrese el numero de obra social: ")
				fmt.Scanf("%d", &dni_medique)				
				fmt.Print("Ingrese la fecha de inicio para liquidar (formato: yyyy-mm-dd HH:MM:SS): ")
				fmt.Scanf("%s", &f_desde)
				fmt.Printf("Ingrese la fecha final para liquidar (formato: yyyy-mm-dd HH:MM:SS): ")
				fmt.Scanf("%s", &f_hasta)
						
				t, err := time.Parse("2006-01-02 15:04:05", f_desde)
				if err != nil {
					fmt.Println("Error al parsear la fecha y hora:", err)
					return
				}			
				t, err := time.Parse("2006-01-02 15:04:05", f_hasta)
				if err != nil {
					fmt.Println("Error al parsear la fecha y hora:", err)
					return
			
			    generarLiquidacionObrasSociales(nro_obra_social, f_desde, f_hasta)                 	
				case opcion == 16:
					fmt.Println("Adios!")
				default:
					fmt.Println("La opciòn ingresada no es vàlida, por favor ingrese ingrese otro numero.")
			}
		}
		fmt.Println("El programa se finalizo correctamente.")
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
		create table error (nro_error int, f_turno timestamp, nro_consultorio int, dni_medique int, nro_paciente int, operacion char(12), f_error timestamp, motivo varchar(64));
		create table cobertura (dni_medique int, nro_obra_social int, monto_paciente decimal(12,2), monto_obra_social decimal(12,2));
		create table obra_social (nro_obra_social int, nombre text, contacto_nombre text, contacto_apellido text, contacto_telefono char(12), contacto_email text);
		create table liquidacion_cabecera (nro_liquidacion int, nro_obra_social int, desde date, hasta date, total decimal(15,2));
		create table liquidacion_detalle (nro_liquidacion int, nro_linea int, f_atencion date, nro_afiliade int, dni_paciente int, nombre_paciente text, apellido_paciente text, dni_medique int, nombre_medique text, apellido_medique text, especialidad varchar(64), monto decimal(12,2));
		create table envio_email (nro_email int, f_generacion timestamp, email_paciente text, asunto text, cuerpo text, f_envio timestamp, estado char(10));
		create table solicitud_reservas (nro_orden int, nro_paciente int, dni_medique int, fecha date, hora time);
		`)
		
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
	if err != nil {
		log.Fatal(err)

    }
	reservarTurno()
	cancelarTurno()
	atencionTurno()
    emailRecordatorioSP()
    emailPerdidaSP()
	liquidacionObrasSociales()
	}

func crearTriggers() {
	if err != nil {
		log.Fatal(err)
	}
	
	emailsTrigger()
}

func eliminarFk() {
	db, err := dbConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.close()
	
	_,err = db = db.Exec(`
		alter table paciente drop constraint fk_paciente restrict;
		alter table agenda drop constraint fk_agenda restrict;
		alter table agenda drop contraint fk_agenda2 restrict;
		alter table turno drop constraint fk_turno restric;
		alter table turno drop constraint fk_turno2 restrict;
		alter table turno drop constraint fk_turno3 restrict;
		alter table turno drop constraint fk_turno4 restrict;
		alter table reprogramacion drop constraint fk_reprogramacion restrict;
		alter table error drop constraint fk_error restrict;
		alter table error drop constraint fk_error2 restrict;
		alter table error drop constraint fk_error3 restrict;
		alter table cobertura drop constraint fk_cobertura restrict;
		alter table cobertura drop constraint fk_cobertura2 restrict;
		alter table liqudiacion_cabecera drop constraint fk_liq_cab restrict;
		alter table liquidacion_detalle drop constraint fk_liq_det restrict;
		alter table liquidacion_detalle drop constraint fk_liq_det2 restrict;`)

		if err != nil{
			log.Fatal(err)
		}
}

func eliminarPK() {
	db, err := dbConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.close()
	
	_,err = db = db.Exec(`
		alter table paciente drop constraint paciente_pk restrict;
		alter table medique drop constraint medique_pk restrict;
		alter table consultorio drop constraint consultorio_pk restrict;
		alter table agenda drop constraint agenda_pk restrict;
		alter table turno drop constraint turno_pk restrict;
		alter table reporgramacion drop constraint reprogramacion_pk restrict;
		alter table error drop constraint error_pk restrict;
		alter table cobertura drop constraint cobertura_pk restrict;
		alter table obra_social drop cosntraint obra_social_pk restrict;
		alter table liquidacion_cabecera drop constraint liquidacion_cabecera_pk restrict;
		alter table liquidacion_detalle drop constraint liquidacion_detalle_pk restrict;
		alter table envio_mail drop constraint envio_mail_pk restrict;
		`)

		if err != nil {
			log.Fatal(err)
		}
}

//stored Procedures
func generarTurnos() {
	db, err := dbConnection()

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = dbExec(`
		create or replace function generar_turnos(_anio int, _mes int) returns boolean as $$
		declare
			fechas_generadas timestamp[];
			horas_generadas time[];
			
		begin
			//verificar si los turnos ya existen
			select * from turno t where extract(year from t.fecha) = _anio and extract(month from t.fecha) = _mes;
			
			if found then
				raise notice 'Los turnos para esas fechas ya fueron generados';
				return false;
			end if;
			
			select array (
				generate_series(
					make_date(_anio, _mes, 01),
					make_date(_anio, _mes + 1, 01) - interval '1 day', interval '1 day'
				) :: timestamp as dias
			) into fechas_generadas;
			
			for i in select dia from agenda loop
				for j in array_length(fechas_generadas) loop
					
				end loop;
			end loop;
			
			
			
			
				for dia in select distinct(a.dia) from agenda a loop
					select date_add(hora_desde) into horas_generadas;
					select date_add(hora_desde, interval 1 hour) into horas_generadas;
					select date_add(hora_desde, interval 2 hour) into horas_generadas;
					select date_add(hora_desde, interval 3 hour) into horas_generadas;
					select date_add(hora_desde, interval 4 hour) into horas_generadas;
					select date_add(hora_desde, interval 5 hour) into horas_generadas;
					select date_add(hora_desde, interval 6 hour) into horas_generadas;
					select date_add(hora_desde, interval 7 hour) into horas_generadas;
				end loop;
			
			
			insert into turno (fecha, dni_medique)
			select fechas_generadas.dias + horas_generadas, a.dni_medique from agenda a where extract(isodow from dias.fechas_generadas) = a.dia and a.dia = 1 where horas_generadas = a.hora_desde;
			
			insert into turno (fecha, dni_medique)
			select fechas_generadas.dias + horas_generadas, a.dni_medique from agenda a where extract(isodow from dias.fechas_generadas) = a.dia and a.dia = 2 where horas_generadas = a.hora_desde;
			
			insert into turno (fecha, dni_medique)
			select fechas_generadas.dias + horas_generadas, a.dni_medique from agenda a where extract(isodow from dias.fechas_generadas) = a.dia and a.dia = 3 where horas_generadas = a.hora_desde;
			
			insert into turno (fecha, dni_medique)
			select fechas_generadas.dias + horas_generadas, a.dni_medique from agenda a where extract(isodow from dias.fechas_generadas) = a.dia and a.dia = 4 where horas_generadas = a.hora_desde;
			
			insert into turno (fecha, dni_medique)
			select fechas_generadas.dias + horas_generadas, a.dni_medique from agenda a where extract(isodow from dias.fechas_generadas) = a.dia and a.dia = 5 where horas_generadas = a.hora_desde;
			
			insert into turno (fecha, dni_medique)
			select fechas_generadas.dias + horas_generadas, a.dni_medique from agenda a where extract(isodow from dias.fechas_generadas) = a.dia and a.dia = 6 where horas_generadas = a.hora_desde;
			
			insert into turno (fecha, dni_medique)
			select fechas_generadas.dias + horas_generadas, a.dni_medique from agenda a where extract(isodow from dias.fechas_generadas) = a.dia and a.dia = 7 where horas_generadas = a.hora_desde;
			
			return true;
		end;
		$$ language plpgsql;
	`)
}

func reservarTurno() {

     db, err := dbConnection()
	 if err != nil {
		 log.Fatal(err)
	 }
	 defer db.Close()

	_, err = dbExec(`
		create or replace function reservar_turno(_nro_paciente int, _dni_medique int, _fecha_hora_turno timestamp) returns boolean as $$

		declare
			obrasocial_datos record;
			paciente_datos record;
			cantidad_turnos_reservados int;
			condicion boolean;

		begin
			
			// verificar su el dni del medique existe
			select m.dni_medique from medique m where m.dni_medique = _dni_medique;
			
			if not found then
				insert into error (f_turno, nro_consultorio, dni_medique, nro_paciente, operacion, f_error, motivo)
				values (now(), null, null, null, 'reserva de turnos', now(), 'dni de medique no valido');
				raise notice 'dni de medique no valido';
				return false;
			end if;
			

			// verificar si el numero de historia clinica existe
			select p.nro_paciente from paciente p where p.nro_paciente = _nro_paciente;
			
			if not found then
				insert into error (f_turno, nro_consultorio, dni_medique, nro_paciente, operacion, f_error, motivo)
				values (now(), null, null, null, 'reserva de turnos', now(), 'nro de historia clinica no valido');
				raise notice 'nro de historia clinica no valido';
				return false;
			end if;
			
			// verificar si el paciente tiene una obra social y obtener la obra social
			condicion := false;
			select p.nro_obra_social from paciente p, obra social o where p.nro_obra_social = o.nro_obra_social into paciente_datos;
			
			if found then
				condicion := true;
			end if;
			
			// verificar si el medique atiende esa obra social
			case
				when condicion then
					select * from cobertura p, obra_social o where c.dni_medique = _dni_medique and c.nro_obra_social = paciente_datos.nro_obra_social;
					
					if not found then
						insert into error (f_turno, nro_consultorio, dni_medique, nro_paciente, operacion, f_error, motivo)
						values (now(), null, null, null, 'reserva de turnos', now(), 'obra social no atendida por el medique');
						raise notice 'obra social no atendida por el medique';
						return false;
					end if;
			end case;
			
			// verificar si el turno esta disponible
			
			select nro_turno from turno t where t.fecha = _fecha_hora_turno and t.dni_medique = _dni_medique and estado = 'disponible'
			
			if not found then
				insert into error (f_turno, nro_consultorio, dni_medique, nro_paciente, operacion, f_error, motivo)
				values (now(), null, null, null, 'reserva de turnos', now(), 'turno inexistente o no disponible');
				raise notice 'turno inexistente o no disponible';
				return false;
			end if;
			
			// verificar si el paciente ha superado el limite de 5 turnos en estado reservado
			select count(*) from turno where nro_paciente = _nro_paciente and estado = 'Reservado' into cantidad_turnos_reservados ;

			if cantidad_turnos_reservados >= 5 then
				insert into error (f_turno, nro_consultorio, dni_medique, nro_paciente, operacion, f_error, motivo)
				values (now(), null, null, null, 'reserva de turnos', now(), 'Supera el lìmite de reserva de turnos');
				raise notice 'Supera el lìmite de reserva de turnos';
				return false;
			end if;
			
			// realizar la reserva del turno
			update turno
			set
			nro_paciente = _nro_paciente,
			nro_obra_social_consulta = paciente_obrasocial,
			estado = 'Reservado',
			f_reserva = datetime
			where
			dni_medique = _dni_medique
			and fecha = _fecha_turno
			and hora = _hora_turno
			and estado = 'Disponible';
			return true;
			end;

			
		end;
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
				raise notice 'nro de turno no valido';
				return false;
			end if;
			
			
			select t.estado from turno t into turno_atendido where t.nro_turno = _nro_turno and t.estado = 'reservado';
			
			if not found then
				insert into error (f_turno, nro_consultorio, dni_medique, nro_paciente, operacion, f_error, motivo)
				values (now(), null, null, null, 'atencion de turnos', now(), 'turno no reservado');
				raise notice 'turno no reservado';
				return false;
			end if;
			
			
			select t.fecha from turno t into turno_atendido where t.nro_turno = _nro_turno and t.fecha = current_date;
			
			if not found then
				insert into error (f_turno, nro_consultorio, dni_medique, nro_paciente, operacion, f_error, motivo)
				values (now(), null, null, null, 'atencion de turnos', now(), 'turno no corresponde a la fecha del dia');
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

func emailRecordatorioSP() {
	
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
		
		`)	
}
	
func emailPerdidaSP() {

	db, err := dbConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	_, err = dbExec(`
	
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

func liquidacionObrasSociales() { 

    db, err :=dbConnection()
	if err != nil {

		los.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
	create or replace function generar_liquidacion_obras_sociales(_nro_obra_social int, _desde DATE, _hasta DATE) returns void as $$
	declare 
		total_liquidacion float := 0.0;
		nro_liquidacion int;
	begin
		insert into liquidacion_cabecera (nro_obra_social, desde, hasta, total)
		values(_nro_obra_social, _desde, _hasta, total_liquidacion)
		returning nro_liquidacion into nro_liquidacion;

		update turno
		set estado = 'liquidado'
		where nro_obra_social_consulta = _nro_obra_social
		and fecha between _desde and _hasta;

		insert into liquidacion_detalle (nro_liquidacion, f_atencion, nro_afiliade, dni_paciente, nombre_paciente, apellido_paciente, dni_medique, apellido_medique, especialidad, monto)
		select
		nro_liquidacion,
		t.fecha,
		t.nro_afiliade_consulta,
		t.dni_paciente,
		p.nombre,
		p.apellido,
		t.dni_medique,
		m.nombre,
		m.apellido,
		m.especialidad,
		t.monto-obra_social
		from
		turno t,
		paciente p,
		medique m
		where
		t.nro_obra_social_consulta = _nro_obra_social 
		and t.dni_paciente = p.dni_paciente
		and t.dni_medique = m.dni_medique
		and t.fecha between _desde and _hasta;

		select sum(monto) into total_liquidacion
		from liquidacion_detalle
		where nro_liquidacion = nro_liquidacion;

		update liquidacion_cabecera
		set total = total_liquidacion
		where nro_liquidacion = nro_liquidacion;

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
		select p.email, m.nombre, m.apellido, t.fecha from paciente p, medique m, turno t, reprogramacion r where new.nro_turno = r.nro_turno and t.nro_turno = r.nro_turno and t.nro_paciente = p.nro_paciente and t.dni_medique = m.dni_medique into datos_turno;		
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

func generarTurnos(_anio int, _mes int) {
	db, err := dbConnection()
	if err != nil {

		log.Fatal(err)
	}
	defer db.Close()

	_, err = dbExec(`
		select generarTurnos($1,$2);
	`, _anio, _mes)
}
	
func reservarTurno(_nro_paciente int, _dni_medique int, _fecha_hora_turno time.Time)  {
	db, err := dbConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = dbExec(`
		select reservar_turno($1, $2, $3);
	`, _nro_paciente, _dni_medique, _fecha_hora_turno)
}
	
func cancelacionTurnos(_dni_medique int, _fdesde time.Time, _fhasta time.Time)	{
	db, err := dbConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = dbExec(`
		select cancelacion_turnos($1, $2, $3);
	`, _dni_medique, _fdesde, _fhasta)
}
func atencionTurnos(_nro_turno int) {
	db, err := dbConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = dbExec(`
		select atencion_turnos($1);
	`, _nro_turno)	
}
	
func emailRecordatorio() {
	db, err := dbConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = dbExec(`
		select email_recordatorio();
	`)	
}	

func emailPerdidaTurno() {
	db, err := dbConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = dbExec(`
		select email_perdida_turno();
	`)	
}

func generarLiquidacionObrasSociales(_nro_obra_social int, _desde time.Time, _hasta time.Time)	{
	db, err := dbConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = dbExec(`
		select generar_liquidacion_obras_sociales($1, $2, $3);
	`, _nro_obra_social, _desde, _hasta)	
	
}

//Comienzo json
type Paciente struct {
	NroPaciente int
	Nombre string
	Apellido string
	DniPaciente int
	FechaNacimiento string
	NroObraSocial int
	NroAfiliade int
	Domicilio string
	Telefono string
	Email string
}

type Medique struct {
	DniMedique int
	Nombre string
	Apellido string
	Especialidad string
	MontoConsultaPrivada float64
	Telefono string
}

type Consultorio struct {
	NroConsultorio int
	Nombre string
	Domicilio string
	CodigoPostal int
	Telefono string
}

type Turno struct {
	NroTurno int
	Fecha string
	NroConsultorio int
	DniMedique int
	NroPaciente int
	NroObraSocialConsulta int
	NroAfiliadeConsulta int
	MontoPaciente float64
	MontoObraSocial float64
	FechaReserva string
	Estado string	
}

type ObraSocial struct {
	NroObraSocial int
	Nombre string
	ContactoNombre string
	ContactoApellido string
	ContactoTelefono string
	ContactoEmail string
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

func cargarDatosJson() {
	//abrir db
    db, err := bolt.Open("hospital.db", 0600, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

	//Ingreso datos
	pacientes :=[]Paciente{
		{NroPaciente: 1, Nombre: "Martin", Apellido: "Galvarini", DniPaciente: 42660991, FechaNacimiento: "2000-06-30", NroObraSocial: 1, NroAfiliade: 1, Domicilio: "Carlos Pellgrini 2436, Martinez", Telefono: "11 4416-3214", Email: "13.martingalva@gmail.com"},
		{NroPaciente: 2, Nombre: "Lucas", Apellido: "Hombrefredi", DniPaciente: 43724987, FechaNacimiento: "2000-11-05", NroObraSocial: 2, NroAfiliade: 1, Domicilio: "Domingo de Acassuso 1121, La Lucila", Telefono: "11 4491-3211", Email: "luchetti@gmail.com"},
		{NroPaciente: 3, Nombre: "Veronica", Apellido: "Suarez", DniPaciente: 40547321, FechaNacimiento: "1998-05-03", NroObraSocial: 3, NroAfiliade: 1, Domicilio: "Rivadavia 626, Belgrano", Telefono: "11 4498-4321", Email: "verosua@gmail.com"},
	}
    
    mediques :=[]Medique{
		{DniMedique: 398121041, Nombre: "Carlos", Apellido: "Bilardo", Especialidad: "Ginecologo", MontoConsultaPrivada: 1500.50, Telefono: "11 4312-4574"},
		{DniMedique: 10456789, Nombre: "Juan", Apellido: "Gomez", Especialidad: "Cardiología", MontoConsultaPrivada: 15000.00, Telefono: "112533-1234"},
		{DniMedique: 40342233, Nombre: "Maria", Apellido: "Lopez", Especialidad: "Dermatología", MontoConsultaPrivada: 12000.00, Telefono: "113453-5678"},
		{DniMedique: 11442233, Nombre: "Laura", Apellido: "Martinez", Especialidad: "Neurología", MontoConsultaPrivada: 20000.00, Telefono: "118987-3456"},
		{DniMedique:12459135, Nombre: "Alberto", Apellido: "Fuentes", Especialidad: "Pediatría", MontoConsultaPrivada: 16000.00, Telefono: "119374-8901"},
	}
	
	consultorios :=[]Consultorio{
		{NroConsultorio: 1, Nombre: "Favaloro", Domicilio: "Italia 231, Martinez", CodigoPostal: 1640, Telefono: "4799-4153"},
		{NroConsultorio: 2, Nombre: "Santa Catalina", Domicilio: "Belgrano 2809, Benavidez", CodigoPostal: 1621, Telefono: "5952-1897"},
		{NroConsultorio: 3, Nombre: "Las Acacias", Domicilio: "Marcelo 1393, Don Torcuato", CodigoPostal: 1611, Telefono: "6826-7027"},
	}
    
    turnos :=[]Turno{
		{NroTurno: 1, Fecha: "2023-11-23",	NroConsultorio: 1, DniMedique: 398121041, NroPaciente: 1, NroObraSocialConsulta: 1, NroAfiliadeConsulta: 1,	MontoPaciente: 15000.00, MontoObraSocial: 10000.00, FechaReserva: "2023-09-14", Estado: "reservado"},
		{NroTurno: 2, Fecha: "2023-11-24",	NroConsultorio: 2, DniMedique: 10456789, NroPaciente: 2, NroObraSocialConsulta: 2, NroAfiliadeConsulta: 1,	MontoPaciente: 12000.00, MontoObraSocial: 9700.00, FechaReserva: "2023-10-09", Estado: "reservado"},
		{NroTurno: 3, Fecha: "2023-11-25",	NroConsultorio: 3, DniMedique: 40342233, NroPaciente: 3, NroObraSocialConsulta: 3, NroAfiliadeConsulta: 1,	MontoPaciente: 16000.00, MontoObraSocial: 10100.00, FechaReserva: "2023-09-01", Estado: "reservado"},
		{NroTurno: 4, Fecha: "2023-11-23",	NroConsultorio: 1, DniMedique: 398121041, NroPaciente: 1, NroObraSocialConsulta: 1, NroAfiliadeConsulta: 1,	MontoPaciente: 15000.00, MontoObraSocial: 10000.00, FechaReserva: "2023-08-13", Estado: "reservado"},
		{NroTurno: 5, Fecha: "2023-11-24",	NroConsultorio: 2, DniMedique: 11442233, NroPaciente: 2, NroObraSocialConsulta: 2, NroAfiliadeConsulta: 1,	MontoPaciente: 20000.00, MontoObraSocial: 11000.00, FechaReserva: "2023-04-14", Estado: "reservado"},
		{NroTurno: 6, Fecha: "2023-11-25",	NroConsultorio: 3, DniMedique: 398121041, NroPaciente: 3, NroObraSocialConsulta: 3, NroAfiliadeConsulta: 1,	MontoPaciente: 15000.00, MontoObraSocial: 10000.00, FechaReserva: "2023-09-20", Estado: "reservado"},
		{NroTurno: 7, Fecha: "2023-11-23",	NroConsultorio: 1, DniMedique: 12459135, NroPaciente: 1, NroObraSocialConsulta: 1, NroAfiliadeConsulta: 1,	MontoPaciente: 14000.00, MontoObraSocial: 9900.00, FechaReserva: "2023-10-22", Estado: "reservado"},
		{NroTurno: 8, Fecha: "2023-11-24",	NroConsultorio: 2, DniMedique: 11442233, NroPaciente: 2, NroObraSocialConsulta: 2, NroAfiliadeConsulta: 1,	MontoPaciente: 20000.00, MontoObraSocial: 11000.00, FechaReserva: "2023-11-14", Estado: "reservado"},
		{NroTurno: 9, Fecha: "2023-11-25",	NroConsultorio: 3, DniMedique: 10456789, NroPaciente: 3, NroObraSocialConsulta: 3, NroAfiliadeConsulta: 1,	MontoPaciente: 12000.00, MontoObraSocial: 9700.00, FechaReserva: "2023-08-25", Estado: "reservado"},
		{NroTurno: 10, Fecha: "2023-11-23",	NroConsultorio: 1, DniMedique: 11442233, NroPaciente: 1, NroObraSocialConsulta: 1, NroAfiliadeConsulta: 1,	MontoPaciente: 20000.00, MontoObraSocial: 11000.00, FechaReserva: "2023-09-30", Estado: "reservado"},
		{NroTurno: 11, Fecha: "2023-11-24",	NroConsultorio: 2, DniMedique: 12459135, NroPaciente: 2, NroObraSocialConsulta: 2, NroAfiliadeConsulta: 1,	MontoPaciente: 14000.00, MontoObraSocial: 9900.00, FechaReserva: "2023-07-07", Estado: "reservado"},
		{NroTurno: 12, Fecha: "2023-11-25",	NroConsultorio: 3, DniMedique: 398121041, NroPaciente: 3, NroObraSocialConsulta: 3, NroAfiliadeConsulta: 1,	MontoPaciente: 15000.00, MontoObraSocial: 10000.00, FechaReserva: "2023-07-30", Estado: "reservado"},
		{NroTurno: 13, Fecha: "2023-11-23",	NroConsultorio: 1, DniMedique: 398121041, NroPaciente: 1, NroObraSocialConsulta: 1, NroAfiliadeConsulta: 1,	MontoPaciente: 15000.00, MontoObraSocial: 10000.00, FechaReserva: "2023-09-12", Estado: "reservado"},
		{NroTurno: 14, Fecha: "2023-11-24",	NroConsultorio: 2, DniMedique: 398121041, NroPaciente: 2, NroObraSocialConsulta: 2, NroAfiliadeConsulta: 1,	MontoPaciente: 15000.00, MontoObraSocial: 10000.00, FechaReserva: "2023-09-11", Estado: "reservado"},
		{NroTurno: 15, Fecha: "2023-11-25",	NroConsultorio: 3, DniMedique: 10456789, NroPaciente: 3, NroObraSocialConsulta: 3, NroAfiliadeConsulta: 1,	MontoPaciente: 12000.00, MontoObraSocial: 9700.00, FechaReserva: "2023-09-18", Estado: "reservado"},
		
	}
    
    obras_sociales :=[]ObraSocial{
		{NroObraSocial: 1, Nombre: "Galeno", ContactoNombre: "Gabriel", ContactoApellido: "Galindo", ContactoTelefono: "11 4414-44120", ContactoEmail: "galeno@info.com"},
		{NroObraSocial: 2, Nombre: "Omint", ContactoNombre: "Felipe", ContactoApellido: "Luna", ContactoTelefono: "11 5624-1120", ContactoEmail: "omint@info.com"},
		{NroObraSocial: 3, Nombre: "Hospital Italiano", ContactoNombre: "Roberto", ContactoApellido: "Perez", ContactoTelefono: "11 2855-9291", ContactoEmail: "hitaliano@info.com"},
	}
    
    //Escritura de datos
	
	for _, paciente:= range pacientes {
		data, err := json.Marshal(paciente)
		if err != nil {
			log.Fatal(err)
		}
		CreateUpdate(db, "pacientes", []byte(strconv.Itoa(paciente.NroPaciente)), data)
	}
	
	fmt.Printf("Pacientes cargados correctamente.\n")
	
	for _, medique:= range mediques {
		data, err := json.Marshal(medique)
		if err != nil {
			log.Fatal(err)
		}
		CreateUpdate(db, "mediques", []byte(strconv.Itoa(medique.DniMedique)), data)
	}
	
	fmt.Printf("Mediques cargados correctamente.\n")
	
	for _, consultorio:= range consultorios {
		data, err := json.Marshal(consultorio)
		if err != nil {
			log.Fatal(err)
		}
		CreateUpdate(db, "consultorios", []byte(strconv.Itoa(consultorio.NroConsultorio)), data)
	}
	
	fmt.Printf("Consultorios cargados correctamente.\n")
	
	for _, turno:= range turnos {
		data, err := json.Marshal(turno)
		if err != nil {
			log.Fatal(err)
		}
		CreateUpdate(db, "turnos", []byte(strconv.Itoa(turno.NroTurno)), data)
	}
	
	fmt.Printf("Turnos cargados correctamente.\n")
	
	for _, obra_social:= range obras_sociales {
		data, err := json.Marshal(obra_social)
		if err != nil {
			log.Fatal(err)
		}
		CreateUpdate(db, "obras sociales", []byte(strconv.Itoa(obra_social.NroObraSocial)), data)
	}
	
	fmt.Printf("Obras sociales cargados correctamente.\n")
}

func mostrarDatosJson() {
	//abrir db
    db, err := bolt.Open("hospital.db", 0600, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
	
	fmt.Printf("Pacientes:\n")
	for _, paciente:= range pacientes {
		resultado, err := ReadUnique(db, "pacientes", []byte(strconv.Itoa(paciente.NroPaciente)))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", resultado)
	} 
	fmt.Printf("Mediques:\n")
	for _, medique:= range mediques {
		resultado, err := ReadUnique(db, "mediques", []byte(strconv.Itoa(medique.DniMedique)))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", resultado)
	} 
	fmt.Printf("Consultorios:\n")
	for _, consultorio:= range consultorios {
		resultado, err := ReadUnique(db, "consultorios", []byte(strconv.Itoa(consultorio.NroConsultorio)))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", resultado)
	} 
	fmt.Printf("Turnos:\n")
	for _, turno:= range turnos {
		resultado, err := ReadUnique(db, "turnos", []byte(strconv.Itoa(turno.NroTurno)))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", resultado)
	} 
	fmt.Printf("Obras sociales::\n")
	for _, obra_social:= range obras_sociales {
		resultado, err := ReadUnique(db, "obras sociales", []byte(strconv.Itoa(obra_social.NroObraSocial)))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", resultado)
	}
}
