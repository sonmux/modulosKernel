package Models

type PROCESOPADRE struct {
	PROCESO       string        `json:"proceso"`
	PID           int           `json:"pid"`
	ESTADO        int           `json:"estado"`
	MEMORIA_USO   int           `json:"memoria_uso"`
	MEMORIA_TOTAL int           `json:"memoria_total"`
	ID_USUARIO    int           `json:"id_usuario"`
	USUARIO       string        `json:"usuario"`
	PROCESOHIJO   []PROCESOHIJO `json:"procesos_hijo"`
}

type PROCESOHIJO struct {
	PROCESO   string `json:"proceso"`
	PID       int    `json:"pid"`
	ESTADO    int    `json:"estado"`
	PID_PADRE int    `json:"pid_padre"`
}

type DATAJSON struct {
	PROCESOSPADRE []PROCESOPADRE
	PROCESOSHIJO  []PROCESOHIJO
	VM            string `json:"VM"`
}

type DATAJSONCPU struct {
	PORCENTAJE_PROCESO []string `json:"porcentaje_proceso"`
}

type DATAJSONMEMORY struct {
	MEMORIA_TOTAL int    `json:"memoria_total"`
	MEMORIA_LIBRE int    `json:"memoria_libre"`
	BUFFER        int    `json:"buffer"`
	CACHE         int    `json:"cache"`
	MEM_UNIT      int    `json:"mem_unit"`
	VM            string `json:"VM"`
}

type Logcpu struct {
	VM       string   `json:"vm"`
	ENDPOINT string   `json:"endpoint"`
	DATA     DATAJSON `json:"data"`
	FECHA    string   `json:"fecha"`
}
type Logmemoria struct {
	VM       string         `json:"vm"`
	ENDPOINT string         `json:"endpoint"`
	DATA     DATAJSONMEMORY `json:"data"`
	FECHA    string         `json:"fecha"`
}
