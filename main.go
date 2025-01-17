package main

import (
	"log"
	"splay/utils"
)

func main() {
	client := utils.Client()
	id := "https://tlawanvhwwoubgyspgqv.supabase.co/storage/v1/s3"
	resp, err := client.Storage.GetBucket(id)
	if err != nil {
		log.Fatalf("Error al consultar la base de datos: %v", err)
	}

	log.Println("Respuesta obtenida:", resp)
}
