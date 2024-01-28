package route

import (
	"fmt"
	jeux "hangman/game" // Importe le package "game" sous l'alias "jeux"
	"net/http"
	"os"
)

// InitRoute configure les routes de l'application et démarre le serveur web.
func InitRoute() {
	// Définition des gestionnaires pour les différentes routes.
	// Chaque route est associée à une fonction dans le package "jeux".
	http.HandleFunc("/acceuil", jeux.Acceuil) // Gestionnaire pour la page d'accueil
	http.HandleFunc("/jeu", jeux.Jeux)         // Gestionnaire pour la page du jeu
	http.HandleFunc("/regle", jeux.Regle)     // Gestionnaire pour la page des règles
	http.HandleFunc("/mention", jeux.Mention) // Gestionnaire pour la page des mentions légales
	http.HandleFunc("/", jeux.Erreur)
	// Serveur de fichiers statiques pour servir les ressources comme les CSS, images, etc.
	rootDoc, _ := os.Getwd() // Obtient le répertoire de travail actuel
	fileserver := http.FileServer(http.Dir(rootDoc + "/assets")) // Crée un serveur de fichiers statiques
	// Configure le serveur de fichiers statiques pour servir les fichiers sous le chemin "/static/"
	http.Handle("/static/", http.StripPrefix("/static/", fileserver))

	// Affiche un message indiquant que le serveur a démarré.
	fmt.Println("(http://localhost:8080/acceuil) - Server started on port:8080")

	// Démarre le serveur web sur le port 8080 et écoute les requêtes entrantes.
	http.ListenAndServe("localhost:8080", nil)

	// Message indiquant que le serveur est fermé (ce message ne s'affichera jamais car http.ListenAndServe est bloquant).
	fmt.Println("Server closed")
}
