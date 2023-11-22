package main

import (
    "encoding/json"
    "fmt"
    bolt "go.etcd.io/bbolt"
    "log"
    "strconv"
)

// Go: NroPaciente
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
	}
	
	consultorios :=[]Consultorio{
		{NroConsultorio: 1, Nombre: "Favaloro", Domicilio: "Italia 231, Martinez", CodigoPostal: 1640, Telefono: "4799-4153"},
	}
    
    turnos :=[]Turno{
		{NroTurno: 1, Fecha: "2023-11-23",	NroConsultorio: 1, DniMedique: 398121041, NroPaciente: 1, NroObraSocialConsulta: 1, NroAfiliadeConsulta: 1,	MontoPaciente: 15000.00, MontoObraSocial: 10000.00, FechaReserva: "2023-09-14", Estado: "reservado"},
	}
    
    obras_sociales :=[]ObraSocial{
		{NroObraSocial: 1, Nombre: "Galeno", ContactoNombre: "Gabriel", ContactoApellido: "Galindo", ContactoTelefono: "11 4414-44120", ContactoEmail: "galeno@info.com"},
	}
    
    //Escritura de datos
	
	for _, paciente:= range pacientes {
		data, err := json.Marshal(paciente)
		if err != nil {
			log.Fatal(err)
		}
		CreateUpdate(db, "pacientes", []byte(strconv.Itoa(paciente.NroPaciente)), data)
	} 
	
	for _, medique:= range mediques {
		data, err := json.Marshal(medique)
		if err != nil {
			log.Fatal(err)
		}
		CreateUpdate(db, "mediques", []byte(strconv.Itoa(medique.DniMedique)), data)
	} 
	
	for _, consultorio:= range consultorios {
		data, err := json.Marshal(consultorio)
		if err != nil {
			log.Fatal(err)
		}
		CreateUpdate(db, "consultorios", []byte(strconv.Itoa(consultorio.NroConsultorio)), data)
	} 
	
	for _, turno:= range turnos {
		data, err := json.Marshal(turno)
		if err != nil {
			log.Fatal(err)
		}
		CreateUpdate(db, "turnos", []byte(strconv.Itoa(turno.NroTurno)), data)
	} 
	
	for _, obra_social:= range obras_sociales {
		data, err := json.Marshal(obra_social)
		if err != nil {
			log.Fatal(err)
		}
		CreateUpdate(db, "obras sociales", []byte(strconv.Itoa(obra_social.NroObraSocial)), data)
	}
	
	//Lectura de datos
	for _, paciente:= range pacientes {
		resultado, err := ReadUnique(db, "pacientes", []byte(strconv.Itoa(paciente.NroPaciente)))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", resultado)
	} 
	
	for _, medique:= range mediques {
		resultado, err := ReadUnique(db, "mediques", []byte(strconv.Itoa(medique.DniMedique)))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", resultado)
	} 
	
	for _, consultorio:= range consultorios {
		resultado, err := ReadUnique(db, "consultorios", []byte(strconv.Itoa(consultorio.NroConsultorio)))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", resultado)
	} 
	
	for _, turno:= range turnos {
		resultado, err := ReadUnique(db, "turnos", []byte(strconv.Itoa(turno.NroTurno)))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", resultado)
	} 
	
	for _, obra_social:= range obras_sociales {
		resultado, err := ReadUnique(db, "obras sociales", []byte(strconv.Itoa(obra_social.NroObraSocial)))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", resultado)
	}
}
