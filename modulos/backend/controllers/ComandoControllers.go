package Controllers

import (
	Models "Backend/models"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func CMD(comando string) (bytes.Buffer, string, error) { //(entrada,salida)

	var salida bytes.Buffer
	var errors bytes.Buffer
	cmd := exec.Command("bash", "-c", comando)
	cmd.Stdout = &salida
	cmd.Stderr = &errors
	err := cmd.Run()
	return salida, errors.String(), err
}

//	func RemoveIndex(s []int, index int) []int {
//		return append(s[:index], s[index+1:]...)
//	}
func getNameUsuario(arr []Models.PROCESOPADRE) []Models.PROCESOPADRE {

	pos := 0

	for _, item := range arr {
		fmt.Printf(strconv.Itoa(item.ID_USUARIO))
		ejecutar := "getent passwd " + strconv.Itoa(item.ID_USUARIO) + " | cut -d: -f1"
		salida, _, verificar := CMD(ejecutar)

		if verificar != nil {
			log.Printf("error: %v\n", verificar)
		} else {
			arr[pos].USUARIO = salida.String()
			//fmt.Println(errout)
		}
		pos++
	}
	return arr
}

func RequestPrincipal() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		salida, _, verificar := CMD("cat /proc/cpu_grupo5")

		if verificar != nil {
			log.Printf("error: %v\n", verificar)
		} else {
			var dataJson Models.DATAJSON
			json.Unmarshal(salida.Bytes(), &dataJson) //json a objeto
			vm := os.Getenv("NOMBREVM")
			dataJson.VM = vm
			// fmt.Println(errout)

			dataJson.PROCESOSPADRE = getNameUsuario(dataJson.PROCESOSPADRE)
			fmt.Println(dataJson.PROCESOSPADRE)
			//Seccion donde mandamos a guardar a la base de datos usando el endpoint de cloud function
			t := time.Now()
			fechad := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
				t.Year(), t.Month(), t.Day(),
				t.Hour(), t.Minute(), t.Second())
			var logcpu Models.Logcpu
			logcpu.VM = vm
			logcpu.ENDPOINT = "/CPU"
			logcpu.DATA = dataJson
			logcpu.FECHA = fechad
			/*url := "https://us-central1-macro-resolver-341523.cloudfunctions.net/Guardar/GuardarLog"
			mapB, _ := json.Marshal(logcpu)
			peticion, err := http.NewRequest("POST", url, bytes.NewBuffer(mapB))
			if err != nil {
				// Maneja el error de acuerdo a tu situación
				log.Fatalf("Error creando petición: %v", err)
			}
			clienteHttp := &http.Client{}
			respuesta, err := clienteHttp.Do(peticion)
			fmt.Println(respuesta)
			if err != nil {
				// Maneja el error de acuerdo a tu situación
				log.Fatalf("Error haciendo petición: %v", err)
			}*/
			json.NewEncoder(rw).Encode(dataJson)
		}
	}
}

func RequestKill() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/Kill" {
			http.NotFound(rw, r)
			return
		}
		switch r.Method {
		case "GET":
			id := r.URL.Query().Get("pid")
			id = strings.TrimSuffix(id, "/")
			_, _, verificar := CMD("sudo kill -9 " + id)

			if verificar != nil {
				log.Printf("error: %v\n", verificar)
			} else {
				fmt.Println("Eliminando Proceso: " + id)
				// fmt.Println(salida)
				// fmt.Println(errout)
			}
			return
		case "POST":

		default:
			rw.WriteHeader(http.StatusNotImplemented)
			rw.Write([]byte(http.StatusText(http.StatusNotImplemented)))
		}
	}
}

func RequestCPU() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ejecutar := "ps -eo pcpu | sort -k 1 -r | head -50"
		salida, _, verificar := CMD(ejecutar)

		if verificar != nil {
			log.Printf("error: %v\n", verificar)
		} else {
			result := strings.Split(salida.String(), "\n")
			// fmt.Println(result)
			var theArray [50]string
			entro := false
			cont := 0
			for _, item := range result {

				if !entro {
					entro = true
				} else {
					theArray[cont] += strings.TrimSpace(item)
					cont++
				}
			}
			// fmt.Println(errout)
			json.NewEncoder(rw).Encode(theArray)
		}
	}
}
func RequestMemory() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		salida, _, verificar := CMD("cat /proc/ram_grupo5")

		if verificar != nil {
			log.Printf("error: %v\n", verificar)
		} else {

			var dataJson Models.DATAJSONMEMORY
			json.Unmarshal(salida.Bytes(), &dataJson) //json a objeto
			vm := os.Getenv("NOMBREVM")
			dataJson.VM = vm
			ejecutar := "free -m"
			salidaMsj, _, verificar := CMD(ejecutar) //para cache

			if verificar != nil {
				log.Printf("error: %v\n", verificar)
			} else {
				result := strings.Split(salidaMsj.String(), "\n")
				memoria := strings.ReplaceAll(result[1], " ", ",") //fila memoria
				cacheBusqueda := strings.Split(memoria, ",")       //buscar valor cache
				cont := 0
				cachestr := "0"
				for _, item := range cacheBusqueda {
					if item != "" {
						cont++
						if cont == 6 {
							cachestr = item
						}
					}
				}

				valor, err := strconv.Atoi(cachestr)
				if err != nil {
					fmt.Println(err)
				} else {
					dataJson.CACHE = valor * 1000000
				}
			}
			//Seccion donde mandamos a guardar a la base de datos usando el endpoint de cloud function
			t := time.Now()
			fechad := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
				t.Year(), t.Month(), t.Day(),
				t.Hour(), t.Minute(), t.Second())

			var logmemoria Models.Logmemoria
			logmemoria.VM = vm
			logmemoria.ENDPOINT = "/Memoria"
			logmemoria.DATA = dataJson
			logmemoria.FECHA = fechad
			/*url := "https://us-central1-macro-resolver-341523.cloudfunctions.net/Guardar/GuardarLog"
			mapB, _ := json.Marshal(logmemoria)
			peticion, err := http.NewRequest("POST", url, bytes.NewBuffer(mapB))
			if err != nil {
				// Maneja el error de acuerdo a tu situación
				log.Fatalf("Error creando petición: %v", err)
			}
			clienteHttp := &http.Client{}
			respuesta, err := clienteHttp.Do(peticion)
			fmt.Println(respuesta)
			if err != nil {
				// Maneja el error de acuerdo a tu situación
				log.Fatalf("Error haciendo petición: %v", err)
			}*/
			json.NewEncoder(rw).Encode(dataJson)
		}
	}
}
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API GO!\n"))
}
