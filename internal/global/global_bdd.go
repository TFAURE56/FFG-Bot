package global

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/joho/godotenv"
)

// Se connecter a la base de donnée MariaDb
func ConnectToDatabase() (*sql.DB, error) {

	_ = godotenv.Load(".env")

	BDDConnexion := os.Getenv("ConnectDB")

	// Connexion à la base de données
	db, err := sql.Open("mysql", BDDConnexion) // Logique de connexion à la base de données (user:password@tcp(x.x.x.x:3306)/NomBDD)
	if err != nil {
		log.Printf("❌ Erreur lors de l'ouverture de la base de données : %v", err)
		return nil, err
	}
	// Vérifier la connexion
	err = db.Ping()
	if err != nil {
		log.Printf("❌ Erreur de connexion à la base de données : %v", err)
		return nil, err
	}
	return db, nil
}
