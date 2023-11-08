drop database if exists hospital;
create database hospital;

/c hospital

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
	nro_consultorio,
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
alter table turno add constraint fk_turno4 foreign key (nro_obra_social) references obra_social(nro_obra_social);
alter table turno add constraint fk_turno5 foreign key (nro_afiliade_consulta) references paciente(nro_afiliade);

alter table reprogramacion add constraint fk_reprogramacion foreign key (nro_turno) references turno(nro_turno);

alter table error add constraint fk_error foreign key (dni_medique) references medique(dni_medique);
alter table error add constraint fk_error2 foreign key (nro_consultorio) references consultorio(nro_consultorio);
alter table error add constraint fk_error3 foreign key (nro_paciente) references paciente(nro_paciente);

alter table cobertura add constraint fk_cobertura foreign key (dni_medique) references medique(dni_medique);
alter table cobertura add constraint fk_cobertura2 foreign key (nro_obra_social) references obra_social(nro_obra_social);

alter table liquidacion_cabecera add constraint fk_liq_cab foreign key (nro_obra_social) references obra_social(nro_obra_social);

alter table liquidacion_detalle add constraint fk_liq_det foreign key (nro_liquidacion) references liquidacion_cabecera(nro_liquidacion);
alter table liquidacion_detalle add constraint fk_liq_det2 foreign key (nro_afiliade) references paciente(nro_afiliade);
alter table liquidacion_detalle add constraint fk_liq_det3 foreign key (dni_paciente) references paciente(dni_paciente);
alter table liquidacion_detalle add constraint fk_liq_det4 foreign key (dni_medique) references medique(dni_medique);


--FK :
-- en la tabla agenda: dni_meqique fk con dni_medique de la tabla medique
-- nro_consultorio fk con nro_consultorio de la tabla consultorio

-- en la tabla turno: dni_medique fk con dni_medique en la tabla medique
-- nro_consultorio fk con nro_consultorio de la tabla consultorio
-- nro_paciente fk con nro_paciente de la tabla paciente
-- nro_obra_social_consulta fk con nro_obra_social de la tabla cobertura y nro_obra_social de la tabla obra_social ??
-- nro_afiliade_consulta fk con nro_afiliade de la tabla paciente

-- en la tabla reprogramacion: nro_turno fk con nro_turno de la tabla turno

-- en la tabla error: dni_medique fk con dni_medique de la tabla medique
-- nro_consultorio fk con nro_consultorio de la tabla consultorio
-- nro_paciente fk con nro_paciente de la tabla paciente

-- en la tabla cobertura: dni_medique fk con dni_medique de la tabla medique
-- nro_obra_social fk con nro_obra_social de la tabla obra_social


-- en la tabla liquidacion_cabecera: nro_obra_social fk con nro_obra_social de la tabla obra_social

-- en la tabla liquidacion_detalle: nro_liquidacion fk con nro_liquidacion de la tabla liquidacion_cabecera
-- nro_afiliade fk con nro_afiliade de la tabla paciente
-- dni_paciente fk con dni_paciente de la tabla paciente
-- dni_medique fk con dni_medique de la tabla medique










