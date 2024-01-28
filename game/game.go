package jeux

import (
	//Importation des packages nécaissaires
	Temps "hangman/temp"
	"fmt"
	"bufio"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

//Structure pour le jeu et l'etat actuel du jeu
type GameState struct {
	Word           string			// Mot à deviner
	DisplayWord    string			// Représentation du mot avec des _ pour les lettres non devinées
	Tries          int				// Nombre d'essais effectués
	Difficulty     string			// Difficulté de la partie (facile, moyen, difficile...)
	GameOver       bool	    		// Indique si la partie est terminée
	MaxTries       int	   			// Nombre maximum d'essais autorisés
	Alphabet       []string			// Liste des lettres de l'alphabet
	GuessedLetters map[string]bool  // Pour enregistrer les lettres déjà devinées
	Errors         int              // Pour compter le nombre d'erreurs
}


// Pointeur vers l'état actuel du jeu (GameState) pour chaque partie.
var currentGameState *GameState


//Fonction de la page principale du jeu.
func Acceuil(w http.ResponseWriter, r *http.Request) {
	//j'initie le template du acceuil et renvoie une erreur si le template n'est pas trouve
	err := Temps.Temp.ExecuteTemplate(w, "acceuil", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func Regle(w http.ResponseWriter, r *http.Request) {
	//j'initie le template de la page Regle/aide et renvoie une erreur si le template n'est pas trouve
	err := Temps.Temp.ExecuteTemplate(w, "regle", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func Mention(w http.ResponseWriter, r *http.Request) {
	//j'initie le template de la page Mention et renvoie une erreur si le template n'est pas trouve
	err := Temps.Temp.ExecuteTemplate(w, "mention", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}



func Jeux(w http.ResponseWriter, r *http.Request) {
	//fonction de Gestion des erreurs internes du serveur.
	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic détecté :", r)
			http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		}
	}()

	//Recupere la difficulté / catégorie choisis et crée une nouvelle partie avec NewGame
	if r.Method == "GET" {
		difficulty := r.URL.Query().Get("difficulty")
		currentGameState, _ = NewGame(difficulty)
	} else if r.Method == "POST" {
		letter := r.FormValue("letter")
		ProcessGuess(currentGameState, letter)
	}

	//fonction qui me met une page d'erreur
	if currentGameState == nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}

	//j'initie le template du jeux et renvoie une erreur si le template n'est pas trouve
	err := Temps.Temp.ExecuteTemplate(w, "jeux", currentGameState)
	if err != nil {
		log.Println("Erreur lors de l'exécution du template:", err)
		http.Error(w, "Erreur lors de l'exécution du template: "+err.Error(), http.StatusInternalServerError)
	}
}

// StartNewGame crée une nouvelle partie du jeu.
func StartNewGame(w http.ResponseWriter, r *http.Request) {
	// Récupération de la difficulté choisie et création d'une nouvelle partie.
	difficulty := r.URL.Query().Get("difficulty")
	var err error
	currentGameState, err = NewGame(difficulty)
	if err != nil {
		log.Println("Erreur lors de la création du jeu:", err)
		http.Error(w, "Erreur lors de la création du jeu: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
// NewGame initialise une nouvelle partie avec une difficulté donnée.
func NewGame(difficulty string) (*GameState, error) {
	// Obtenir un mot aléatoire en fonction de la difficulté
	word, err := GetRandomWord(difficulty)
	if err != nil {
		return nil, err
	}
	//map avec toute les lettre pour la page de jeux.
	alphabet := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

	// Détermination du nombre maximum d'essais en fonction de la difficulté.
	maxTries := DetermineMaxTries(difficulty)

	// Affichage initial des lettres en fonction de la difficulté.
	displayWord := strings.Repeat("_", len(word))
    switch difficulty {
    case "facile", "moyen", "pays", "animaux", "capitales", "musique", "sport":
        displayWord = revealLetter(word, displayWord, 1)
    case "difficile":
        displayWord = revealLetter(word, displayWord, 2)
    }

	//j'inistialise le jeux avec les infos données
	return &GameState{
		Word:           word,
        DisplayWord:    displayWord,
		Tries:          0,
		MaxTries:       maxTries,
		GuessedLetters: make(map[string]bool),
		Alphabet:       alphabet,
		Difficulty:     difficulty,
	}, nil
}

// revealLetter révèle un nombre spécifié de lettres au début du jeu
func revealLetter(word, displayWord string, count int) string {
    for i := 0; i < count; i++ {
        randIndex := rand.Intn(len(word))
        displayWord = displayWord[:randIndex] + string(word[randIndex]) + displayWord[randIndex+1:]
    }
    return displayWord
}


// ProcessGuess traite les propositions de lettres ou de mots du joueur.
func ProcessGuess(gameState *GameState, guess string) {
	if gameState == nil {
		return
	}

	// Vérifie si le nombre d'erreurs a atteint le maximum.
	if gameState.Errors >= gameState.MaxTries {
		return
	}

	// Mise à jour de l'état du jeu en fonction de la lettre devinée.
	if len(guess) > 1 {
        if strings.ToLower(guess) == gameState.Word {
            gameState.DisplayWord = gameState.Word // Le joueur a deviné le mot.
        } else {
            gameState.Errors++ // Le joueur perd une tentative
        }
        return
    }

	if _, guessed := gameState.GuessedLetters[guess]; guessed {
		// La lettre a déjà été devinée, ne rien faire
		return
	}

	// Logique pour mettre à jour DisplayWord et Tries en fonction de la lettre devinée
	gameState.GuessedLetters[guess] = true

	correctGuess := false
	for i, char := range gameState.Word {
		if string(char) == guess {
			gameState.DisplayWord = gameState.DisplayWord[:i] + guess + gameState.DisplayWord[i+1:]
			correctGuess = true
		}
	}

	
	// Gestion des erreurs.
	if !correctGuess {
		gameState.Errors++
	}

	// Vérification si le jeu est terminé.
	if gameState.Errors >= gameState.MaxTries {
		gameState.GameOver = true
	}

}


// ProcessUserGuess traite les propositions de lettres de l'utilisateur.
func ProcessUserGuess(w http.ResponseWriter, r *http.Request) {
	if currentGameState == nil {
		return // Initialiser un nouveau jeu si nécessaire.
	}

	letter := r.FormValue("letter")
	log.Println("Lettre devinée:", letter)
	ProcessGuess(currentGameState, letter)
}


// GetRandomWord obtient un mot aléatoire en fonction de la difficulté.
func GetRandomWord(difficulty string) (string, error) {
	var filename string
	switch difficulty {
	case "facile":
		filename = "game/facile.txt" // Changer le nom du fichier pour la catégorie "facile"
	case "moyen":
		filename = "game/moyen.txt" // Changer le nom du fichier pour la catégorie "moyen"
	case "difficile":
		filename = "game/difficile.txt" // Changer le nom du fichier pour la catégorie "difficile"
	case "pays":
		filename = "game/pays.txt" // Changer le nom du fichier pour la catégorie "pays"
	case "capitales":
		filename = "game/capitales.txt" // Changer le nom du fichier pour la catégorie "capitales"
	case "animaux":
		filename = "game/animaux.txt" // Changer le nom du fichier pour la catégorie "animaux"
	case "musique":
		filename = "game/musique.txt" // Changer le nom du fichier pour la catégorie "musique"
	case "sport":
		filename = "game/sport.txt" // Changer le nom du fichier pour la catégorie "sport"
	default:
		return "", fmt.Errorf("invalid difficulty")
	}

	return GetRandomWordLowercase(filename)
}


// GetRandomWordLowercase lit un fichier de mots et en sélectionne un aléatoirement.
func GetRandomWordLowercase(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		// Convertir le mot en minuscules
		word = strings.ToLower(word)
		words = append(words, word)
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	rand.Seed(time.Now().UnixNano())
	return words[rand.Intn(len(words))], nil
}

// DetermineMaxTries détermine le nombre maximal d'essais en fonction de la difficulté.
func DetermineMaxTries(difficulty string) int {
	switch difficulty {
	case "facile":
		return 10 // Plus d'essais pour les niveaux faciles
	case "moyen":
		return 7  // Moins d'essais pour les niveaux moyens
	case "difficile":
		return 5 // Encore moins d'essais pour les niveaux difficiles
	case "pays", "capitales", "animaux", "musique", "sport":
        return 10 // Nombre d'essais standard pour les catégories spéciales
	default:
		return 10 // Nombre par défaut d'essais si la difficulté n'est pas reconnue
	}
}
