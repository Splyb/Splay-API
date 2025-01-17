package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

func Client() *supabase.Client {
	err := godotenv.Load()
	if err != nil {
		log.Println("Advertencia: No se pudo cargar el archivo .env")
	}

	SUPABASE_URL := os.Getenv("SUPABASE_URL")
	SUPABASE_SERVICE_KEY := os.Getenv("SUPABASE_SERVICE_KEY")

	if SUPABASE_URL == "" || SUPABASE_SERVICE_KEY == "" {
		panic("[-] Error: SUPABASE_URL y SUPABASE_SERVICE_KEY son requeridos")
	}

	client, err := supabase.NewClient(SUPABASE_URL, SUPABASE_SERVICE_KEY, &supabase.ClientOptions{})
	if err != nil {
		panic("[-] Error: " + err.Error())
	}

	return client
}
