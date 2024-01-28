package temps

import (
	"fmt"
	"html/template" // Importation du package pour la gestion des templates HTML
	"os"            // Importation du package pour les fonctionnalités du système d'exploitation
)

// Temp est une variable globale qui stockera le template compilé.
var Temp *template.Template

// IniTemps est la fonction d'initialisation pour charger les templates HTML.
func IniTemps() {
	// ParseGlob lit tous les fichiers .html dans le dossier './temp/' et compile les templates.
	temp, errTemp := template.ParseGlob("./temp/*.html")

	// Gestion d'erreur si le chargement des templates échoue.
	if errTemp != nil {
		fmt.Println("Error template:", errTemp) // Affiche l'erreur dans la console.
		os.Exit(1)                              // Quitte le programme avec un code d'erreur 1.
	}

	// Si les templates sont correctement chargés, ils sont assignés à la variable globale Temp.
	Temp = temp
}