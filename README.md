# Administración de Turnos Médicos - Proyecto Final

Este proyecto tiene como objetivo la gestión y administración de turnos médicos utilizando una base de datos relacional y NoSQL, con una aplicación de línea de comandos (CLI) escrita en Go. El sistema permite gestionar la información de los pacientes, médicos, consultorios, obras sociales, turnos, reprogramaciones, y liquidaciones mensuales. 

## Tecnologías Utilizadas

- **Lenguajes:**
  - **SQL (PostgreSQL)**: Para el manejo de la base de datos relacional, creación de tablas, relaciones, y procedimientos almacenados.
  - **Go**: Para la implementación de una CLI que interactúa con la base de datos, gestionando los turnos, reservas, cancelaciones, entre otros.
  
- **Base de Datos:**
  - **PostgreSQL**: Para la implementación de la base de datos relacional.
  - **BoltDB (NoSQL)**: Para almacenar los datos de pacientes, médicos, consultorios, obras sociales y turnos en un formato no relacional basado en JSON.

## Funcionalidades

1. **Gestión de Pacientes:**
   - Alta de pacientes con datos como número de historia clínica, nombre, apellido, fecha de nacimiento, obra social, contacto, etc.
   
2. **Gestión de Médicos:**
   - Alta de médicos con datos como DNI, especialidad, monto de consulta privada, y contacto.
   
3. **Agenda de Médicos:**
   - Registro de la disponibilidad de los médicos en días específicos y horarios, incluyendo la duración de los turnos.
   
4. **Reserva de Turnos:**
   - Los pacientes pueden reservar turnos con médicos disponibles. El sistema valida varios aspectos antes de confirmar la reserva, como la existencia del médico, la validez del número de historia clínica, y el cumplimiento de los turnos disponibles.
   
5. **Cancelación y Reprogramación de Turnos:**
   - Permite cancelar turnos y registrarlos en una tabla de reprogramación para que el centro de atención pueda contactar a los pacientes afectados.
   
6. **Atención de Turnos:**
   - Marca los turnos como atendidos, asegurándose de que el turno esté reservado y sea del día correspondiente.
   
7. **Liquidación para Obras Sociales:**
   - Genera una liquidación mensual por obra social, incluyendo el monto a abonar por los pacientes y el detalle de las atenciones realizadas.
   
8. **Envío Automático de Correos Electrónicos:**
   - Genera correos electrónicos automáticos para notificar a los pacientes sobre:
     - Reserva de turno.
     - Cancelación de turno.
     - Recordatorio de turno (para los turnos próximos).
     - Pérdida de turno reservado (para turnos no atendidos).
   
9. **Modelo NoSQL:**
   - Además de la base de datos relacional, los datos son almacenados en una base de datos NoSQL (BoltDB), permitiendo una comparación entre el modelo relacional y no relacional.
   
10. **CLI en Go:**
    - La aplicación CLI permite interactuar con las funcionalidades del sistema de manera sencilla, con comandos específicos para la reserva, cancelación, y manejo de turnos.

## Estructura del Proyecto

- **`/cmd`**: Contiene el código principal de la aplicación CLI.
- **`/db`**: Archivos relacionados con la base de datos, como scripts de creación de tablas, relaciones y procedimientos almacenados.
- **`/models`**: Definiciones de estructuras de datos, tanto para la base de datos relacional como NoSQL.
- **`/scripts`**: Scripts adicionales para pruebas o configuraciones.
- **`/tests`**: Archivos de prueba para validar el funcionamiento de la aplicación.

## Procedimientos Almacenados y Triggers

El proyecto incluye varios procedimientos almacenados y triggers para la automatización de la gestión de turnos:

- **Generación de turnos disponibles**: Automatiza la creación de turnos para todos los médicos en el mes y año especificado.
- **Reserva de turno**: Permite la validación de reservas y la actualización de estados.
- **Cancelación de turnos**: Cancela turnos y los registra para reprogramación.
- **Atención de turnos**: Marca los turnos como atendidos y actualiza el estado.
- **Liquidación de obras sociales**: Genera y marca como liquidados los turnos atendidos.
- **Envío de correos**: Genera correos automáticos sobre las distintas acciones que ocurren en el sistema.
