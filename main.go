package main

import (
	r "hangman/route" // Importe le package 'route' avec l'alias 'r'
	t "hangman/temp"  // Importe le package 'temp' avec l'alias 't'
)

// La fonction main est le point d'entrée de l'application.
func main() {
	t.IniTemps() // Appelle la fonction IniTemps du package 'temp' pour initialiser les templates HTML.
	r.InitRoute() // Appelle la fonction InitRoute du package 'route' pour configurer et démarrer le serveur web.
}
