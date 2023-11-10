drop database if exists hospital;
create database hospital;

\c hospital

--creación tablas

create table paciente(
	nro_paciente int, --numero de historia clinica
	nombre  text,
	apellido text,
	dni_paciente int,
	f_nacimiento date,
	nro_obra_social int,
	nro_afiliade int,
	domicilio text,
	telefono char(12),
	email text --valido
);

create table medique(
	dni_medique int,
	nombre text,
	apellido text,
	especialidad varchar(64),
	monto_consulta_privada decimal(12, 2), --para pacientes sin obra social
	telefono char(12)
);

create table consultorio(
	nro_consultorio int,
	domicilio text,
	codigo_postal char(8),
	telefono char(12)
);

create table agenda(
	dni_medique int,
	dia int, --lunes:1, martes:2, etc.
	nro_consultorio int,
	hora_desde time, -- 08:30,11:45, etc.
	hora_hasta time,
	duracion_turno interval
);	
	
create table turno(	
	nro_turno int,
	fecha timestamp,
	nro_consultorio int,
	dni_medique int,
	nro_paciente int,
	nro_obra_social_consulta int, --para las liquidaciones
	nro_afiliade_consulta int,
	monto_paciente decimal(12,2),
	monto_obra_social decimal(12,2), --para las liquidaciones
	f_reserva timestamp,
	estado char(10) --`disponible',`reservado',`cancelado',`atendido',`liquidado'
);

create table reprogramacion(
	nro_turno int,
	nombre_paciente text,
	apellido_paciente text,
	telefono_paciente char(12),
	email_paciente text,
	nombre_medique text,
	apellido_medique text,
	estado char(12) --`pendiente', `reprogramado', `desistido'
);

create table error(
	nro_error int,
	f_turno timestamp,
	nro_consultorio int,
	dni_medique int,
	nro_paciente int,
	operacion char(12), --`reserva', `cancelación', `atención', `liquidación'
	f_error timestamp,
	motivo varchar(64)
);

create table cobertura(
	dni_medique int,
	nro_obra_social int,
	monto_paciente decimal(12,2), --monto a abonar por el paciente
	monto_obra_social decimal(12,2) --monto a liquidar a la obra social
);

create table obra_social (
	nro_obra_social int,
	nombre text,
	contacto_nombre text,
	contacto_apellido text,
	contacto_telefono char(12),
	contacto_email text
);

create table liquidacion_cabecera(
	nro_liquidacion int,
	nro_obra_social int,
	desde date,
	hasta date,
	total decimal(15,2)
);

create table liquidacion_detalle(
	nro_liquidacion int,
	nro_linea int,
	f_atencion date,
	nro_afiliade int,
	dni_paciente int,
	nombre_paciente text,
	apellido_paciente text,
	dni_medique int,
	nombre_medique text,
	apellido_medique text,
	especialidad varchar(64),
	monto decimal(12,2)
);

create table envio_email(
	nro_email int,
	f_generacion timestamp,
	email_paciente text,
	asunto text,
	cuerpo text,
	f_envio timestamp,
	estado char(10) --`pendiente', `enviado'
);

-- Esta tabla *no* es parte del modelo de datos, pero se incluye para
-- poder probar las funciones.
create table solicitud_reservas(
	nro_orden int, --en qué orden se ejecutarán las reservas
	nro_paciente int,
	dni_medique int,
	fecha date,
	hora time
);

--Primary Keys
alter table paciente add constraint paciente_pk primary key(nro_paciente);
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
alter table envio_email add constraint envio_emailpk primary key (nro_email);


--Foreign Keys
alter table paciente add constraint fk_paciente foreign key (nro_obra_social) references obra_social(nro_obra_social);

alter table agenda add constraint fk_agenda foreign key (dni_medique) references medique(dni_medique);
alter table agenda add constraint fk_agenda2 foreign key (nro_consultorio) references consultorio(nro_consultorio);

alter table turno add constraint fk_turno foreign key (dni_medique) references medique(dni_medique);
alter table turno add constraint fk_turno2 foreign key (nro_consultorio) references consultorio(nro_consultorio);
alter table turno add constraint fk_turno3 foreign key (nro_paciente) references paciente(nro_paciente);
alter table turno add constraint fk_turno4 foreign key (nro_obra_social_consulta) references obra_social(nro_obra_social);
--alter table turno add constraint fk_turno5 foreign key (nro_afiliade_consulta) references paciente(nro_afiliade);

alter table reprogramacion add constraint fk_reprogramacion foreign key (nro_turno) references turno(nro_turno);

alter table error add constraint fk_error foreign key (dni_medique) references medique(dni_medique);
alter table error add constraint fk_error2 foreign key (nro_consultorio) references consultorio(nro_consultorio);
alter table error add constraint fk_error3 foreign key (nro_paciente) references paciente(nro_paciente);

alter table cobertura add constraint fk_cobertura foreign key (dni_medique) references medique(dni_medique);
alter table cobertura add constraint fk_cobertura2 foreign key (nro_obra_social) references obra_social(nro_obra_social);

alter table liquidacion_cabecera add constraint fk_liq_cab foreign key (nro_obra_social) references obra_social(nro_obra_social);

alter table liquidacion_detalle add constraint fk_liq_det foreign key (nro_liquidacion) references liquidacion_cabecera(nro_liquidacion);
alter table liquidacion_detalle add constraint fk_liq_det2 foreign key (dni_medique) references medique(dni_medique);
--alter table liquidacion_detalle add constraint fk_liq_det3 foreign key (nro_afiliade) references paciente(nro_afiliade);
--alter table liquidacion_detalle add constraint fk_liq_det4 foreign key (dni_paciente) references paciente(dni_paciente);

--Ingreso de obras sociales
insert into obra_social (nro_obra_social, nombre, contacto_nombre, contacto_apellido, contacto_telefono, contacto_email) 
values 
(1, 'Galeno', 'Juan', 'Galeno', '4798-2345', 'galeno@gmail.com'),
(2, 'Swiss Medical', 'Maria', 'Suazo', '4545-6532', 'swissmedical@gmail.com'),
(3, 'OSDE', 'Orlando', 'Debarro', '4891-3214', 'osde@gmail.com'),
(4, 'Omint', 'Omar', 'Lopez', '4671-5634', 'omint@gmail.com'),
(5, 'Medicus', 'Mercedes', 'Costa', '4761-8799', 'medicus@gmail.com'),
(6, 'Sancor Seguros', 'Sandra', 'Corleone', '4531-1927', 'sancor@gmail.com');

--Ingreso 20 pacientes
insert into paciente values (1, 'Martin', 'Galvarini', 42660991, '2000-06-30', 1, 1000, 'Carlos Pellegrini 2436, Martinez', '0114416-3214', '13.martingalva@gmail.com');
insert into paciente values (2, 'Pascual', 'Galvarini', 60123321, '2015-01-15', 2, 2000, 'Carlos Pellegrini 2436, Martinez', '0113395-2194', 'pascualgalva@gmail.com');
insert into paciente values (3, 'Lucas', 'Manfredi', 43021777, '2000-11-02', 3, 3000, 'Domingo de Acassuso 150, La Lucila' , '0117483-2745', 'luquita@gmail.com');
insert into paciente values (4, 'Juan Ignacio', 'Mussino', 41345123, '1999-07-23', 4, 4000, 'Centenario 100, San Isidro', '0115467-1234', 'juanimu@aol.com');
insert into paciente values (5, 'Lorenzo', 'Paparo', 44736192, '2001-04-10', 5, 5000, 'Calle Cerrada 500, Villa Adelina', '0118432-8326', 'lorepapa@yahoo.com');
insert into paciente values (6, 'Gianluca', 'Zeolla', 46932721, '2004-10-19', 6, 6000, 'Calle Lejana 1500, Muy muy lejano', '0115674-2341', 'gianze@live.com');
insert into paciente values (7, 'Juan', 'Pérez', 12345678, '1980-01-01', 1, 5678, 'Calle Falsa 123, Victoria', '0111234-5678', 'juanperez@example.com');
insert into paciente values (8, 'María', 'González', 23456789, '1985-02-02', 2, 6789, 'Calle Falsa 234, La Matanza', '0112345-6789', 'mariagonzalez@example.com');
insert into paciente values (9, 'Pedro', 'Rodríguez', 34567890, '1990-03-03', 3, 7890, 'Calle Falsa 345, Tigre', '0113456-7890', 'pedrorodriguez@example.com');
insert into paciente values (10, 'Lucía', 'Fernández', 45678901, '1995-04-04', 4, 8901, 'Calle Falsa 456, Olivos', '0114567-8901', 'luciafernandez@example.com');
insert into paciente values (11, 'Jorge', 'Gómez', 56789012, '2000-05-05', 5, 9012, 'Calle Falsa 567, Florida', '0115678-9012', 'jorgegomez@example.com');
insert into paciente values (12, 'Ana', 'Díaz', 67890123, '2005-06-06', 6, 1234, 'Calle Falsa 678, Villa Crespo', '0116789-1234', 'anadiaz@example.com');
insert into paciente values (13, 'Diego', 'Martínez', 78901234, '2010-07-07', 3, 2345, 'Calle Falsa 789, Caballito', '0117890-2345', 'diegomartinez@example.com');
insert into paciente values (14, 'Carla', 'Pérez', 89012345, '2015-08-08', 4, 3456, 'Calle Falsa 890, Belgrano', '0118901-3456', 'carlaperez@example.com');
insert into paciente values (15, 'Lucas', 'González', 90123456, '2020-09-09', 2, 4567, 'Calle Falsa 901, Palermo', '0119012-4567', 'lucasgonzalez@example.com');
insert into paciente values (16, 'Sofía', 'Rodríguez', 12345679, '1980-05-10', 1, 5778, 'Calle Falsa 1341, Pinamar', '0111234-4312', 'sofiarodriguez@example.com');
insert into paciente values (17, 'Carlos', 'Bianchi', 2543765, '1949-04-26', 1, 17000, 'Campeones 2000, La Boca', '0114637-7584', 'virrey@hotmail.com');
insert into paciente values (18, 'Manuel', 'Belgrano', 1, '1778-08-15', 4, 18000, 'Santa Fe 1812, Rosario', '0114531-1234', 'manubel@gmail.com');
insert into paciente values (19, 'Lionel', 'Messi', 35094577, '1987-06-24', 3, 10, 'La Pampa 2133, Belgrano', '0114126-7789', 'lio.d10s@yopmail.com');
insert into paciente values (20, 'Carlos Saul', 'Menem', 11365578, '1930-07-02', 4, 1312, 'Dos Metros Bajo Tierra, La Rioja', '0119012-1243', 'charly.menem10@outlook.com');
	
--Ingreso 20 mediques	
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
