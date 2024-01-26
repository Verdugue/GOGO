package route

import (
	"fmt"
	jeux "hangman/game"
	"net/http"
	"os"
)

func InitRoute() {
	http.HandleFunc("/acceuil", jeux.Acceuil)
	http.HandleFunc("/jeu", jeux.Jeux)
	http.HandleFunc("/regle", jeux.Regle)
	http.HandleFunc("/mention", jeux.Mention)

	rootDoc, _ := os.Getwd()
	fileserver := http.FileServer(http.Dir(rootDoc + "/assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fileserver))

	fmt.Println("(http://localhost:8080/acceuil) - Server started on port:8080")
	http.ListenAndServe("localhost:8080", nil)
	fmt.Println("Server closed")
}
