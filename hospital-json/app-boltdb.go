package main

import (
    "encoding/json"
    "fmt"
    bolt "go.etcd.io/bbolt"
    "log"
    "strconv"
)

// SQL: nro_paciente
// Go: NroPaciente
type Paciente struct {
	NroPaciente int // fixme
	nombre  string
	apellido string `json:apellido`
	dni_paciente int
	f_nacimiento string
	nro_obra_social int
	nro_afiliade int
	domicilio string
	telefono string
	email string
}

type Medique struct {
	dni_medique int
	nombre string
	apellido string
	especialidad string
	monto_consulta_privada float64
	telefono string
}

type Consultorio struct {
	nro_consultorio int
	nombre string
	domicilio string
	codigo_postal int
	telefono string
}

type Turno struct {
	nro_turno int
	fecha string
	nro_consultorio int
	dni_medique int
	nro_paciente int
	nro_obra_social_consulta int
	nro_afiliade_consulta int
	monto_paciente float64
	monto_obra_social float64
	f_reserva string
	estado string	
}

type Obra_social struct {
	nro_obra_social int
	nombre string
	contacto_nombre string
	contacto_apellido string
	contacto_telefono string
	contacto_email string
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
		{nro_paciente: 1, nombre: "Martin", apellido: "Galvarini", dni_paciente: 42660991, f_nacimiento: "2000-06-30", nro_obra_social: 1, nro_afiliade: 1, domicilio: "Carlos Pellgrini 2436, Martinez", telefono: "11 4416-3214", email: "13.martingalva@gmail.com"},
		{nro_paciente: 2, nombre: "Lucas", apellido: "Hombrefredi", dni_paciente: 43724987, f_nacimiento: "2000-11-05", nro_obra_social: 2, nro_afiliade: 1, domicilio: "Domingo de Acassuso 1121, La Lucila", telefono: "11 4491-3211", email: "luchetti@gmail.com"},
		{nro_paciente: 3, nombre: "Veronica", apellido: "Suarez", dni_paciente: 40547321, f_nacimiento: "1998-05-03", nro_obra_social: 3, nro_afiliade: 1, domicilio: "Rivadavia 626, Belgrano", telefono: "11 4498-4321", email: "verosua@gmail.com"},
	}
    
    mediques :=[]Medique{
		{dni_medique: 398121041, nombre: "Carlos",	apellido: "Bilardo", especialidad: "Ginecologo", monto_consulta_privada: 1500.50, telefono: "11 4312-4574"},
	}
	
	consultorios :=[]Consultorio{
		{nro_consultorio: 1, nombre: "Favaloro", domicilio: "Italia 231, Martinez", codigo_postal: 1640, telefono: "4799-4153"},
	}
    
    turnos :=[]Turno{
		{nro_turno: 1, fecha: "2023-11-23",	nro_consultorio: 1, dni_medique: 398121041,	nro_paciente: 1, nro_obra_social_consulta: 1, nro_afiliade_consulta: 1,	monto_paciente: 15000.00, monto_obra_social: 10000.00, f_reserva: "2023-09-14", estado: "reservado"},
	}
    
    obras_sociales :=[]Obra_social{
		{nro_obra_social: 1, nombre: "Galeno", contacto_nombre: "Gabriel", contacto_apellido: "Galindo",	contacto_telefono: "11 4414-44120",	contacto_email: "galeno@info.com"},
	}
    
    //Escritura de datos
	
	for _, paciente:= range pacientes {
		data, err := json.Marshal(paciente)
		if err != nil {
			log.Fatal(err)
		}
		CreateUpdate(db, "pacientes", []byte(strconv.Itoa(paciente.nro_paciente)), data)
	} 
	
	for _, medique:= range mediques {
		data, err := json.Marshal(medique)
		if err != nil {
			log.Fatal(err)
		}
		CreateUpdate(db, "mediques", []byte(strconv.Itoa(medique.dni_medique)), data)
	} 
	
	for _, consultorio:= range consultorios {
		data, err := json.Marshal(consultorio)
		if err != nil {
			log.Fatal(err)
		}
		CreateUpdate(db, "consultorios", []byte(strconv.Itoa(consultorio.nro_consultorio)), data)
	} 
	
	for _, turno:= range turnos {
		data, err := json.Marshal(turno)
		if err != nil {
			log.Fatal(err)
		}
		CreateUpdate(db, "turnos", []byte(strconv.Itoa(turno.nro_turno)), data)
	} 
	
	for _, obra_social:= range obras_sociales {
		data, err := json.Marshal(obra_social)
		if err != nil {
			log.Fatal(err)
		}
		CreateUpdate(db, "obras sociales", []byte(strconv.Itoa(obra_social.nro_obra_social)), data)
	}
	
	//Lectura de datos
	for _, paciente:= range pacientes {
		resultado, err := ReadUnique(db, "pacientes", []byte(strconv.Itoa(paciente.nro_paciente)))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", resultado)
	} 
	
	for _, medique:= range mediques {
		resultado, err := ReadUnique(db, "mediques", []byte(strconv.Itoa(medique.dni_medique)))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", resultado)
	} 
	
	for _, consultorio:= range consultorios {
		resultado, err := ReadUnique(db, "consultorios", []byte(strconv.Itoa(consultorio.nro_consultorio)))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", resultado)
	} 
	
	for _, turno:= range turnos {
		resultado, err := ReadUnique(db, "turnos", []byte(strconv.Itoa(turno.nro_turno)))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", resultado)
	} 
	
	for _, obra_social:= range obras_sociales {
		resultado, err := ReadUnique(db, "obras sociales", []byte(strconv.Itoa(obra_social.nro_obra_social)))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", resultado)
	}
}
