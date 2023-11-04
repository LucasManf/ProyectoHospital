drop database if exists hospital;
create database hospital;

/c hospital

create table paciente(
	nro_paciente int,
	nombre  text,
	apellido text,
	dni_paciente int,
	f_nacimiento date,
	nro_obra_social int,
	nro_afiliade int,
	domicilio text,
	telefono char(12),
	email text
);

create table medique(
	dni_medique int,
	nombre text,
	apellido text,
	especialidad varchar(64),
	monto_consulta_privada decimal(12, 2),
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
	dia int,
	nro_consultorio int,
	dni_medique int,
	nro_paciente int,
	nro_obra_social int,
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
):

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

alter table paciente add constraint paciente_pk primary key(nro_paciente);
alter table medique add constraint medique_pk primary key (dni_medique);
alter table consultorio add constraint consultorio_pk primary key (nro_consultorio);
alter table agenda add constraint agenda_pk primary key (dni_medique, dni);
alter table turno add constraint turno_pk primary key (nro_turno);
alter table reprogramacion add constraint reprogramacion_pk primary key (nro_turno);
alter table error add constraint error_pk primary key (nro_error);
alter table cobertura add constraint cobertura_pk primary key (dni_medique, nro_obra_social);
alter table obra_social add constraint obra_social_pk primary key (nro_obra_social);
alter table liquidacion_cabecera add constraint liquidacion_cabecera_pk primary key (nro_liquidacion);
alter table liquidacion_detalle add constraint liquidacion_detalle_pk primary key (nro_liquidacion, nro_linea);
alter table envio_email add constraint envio_emailpk primary key (nro_email);











