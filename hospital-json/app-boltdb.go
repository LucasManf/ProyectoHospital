package main

import (
    "encoding/json"
    "fmt"
    bolt "go.etcd.io/bbolt"
    "log"
    "strconv"
)

type Paciente struct {
	NroPaciente int
	Nombre  string
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

func main() {
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
	
	//Lectura de datos
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
