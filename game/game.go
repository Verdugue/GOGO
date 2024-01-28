package jeux

import (
	"bufio"
	"fmt"
	Temps "hangman/temp"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type GameState struct {
	Word           string
	DisplayWord    string
	Tries          int
	Difficulty     string
	GameOver       bool
	MaxTries       int
	Alphabet       []string
	GuessedLetters map[string]bool // Pour enregistrer les lettres déjà devinées
	Errors         int             // Pour compter le nombre d'erreurs
}

var currentGameState *GameState

func Acceuil(w http.ResponseWriter, r *http.Request) {
	//j'initie le template du acceuil et renvoie une erreur si le template n'est pas trouve
	err := Temps.Temp.ExecuteTemplate(w, "acceuil", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func Regle(w http.ResponseWriter, r *http.Request) {
	//j'initie le template du acceuil et renvoie une erreur si le template n'est pas trouve
	err := Temps.Temp.ExecuteTemplate(w, "regle", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func Mention(w http.ResponseWriter, r *http.Request) {
	//j'initie le template du acceuil et renvoie une erreur si le template n'est pas trouve
	err := Temps.Temp.ExecuteTemplate(w, "mention", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}



func Jeux(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic détecté :", r)
			http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		}
	}()

	if r.Method == "GET" {
		difficulty := r.URL.Query().Get("difficulty")
		currentGameState, _ = NewGame(difficulty)
	} else if r.Method == "POST" {
		letter := r.FormValue("letter")
		ProcessGuess(currentGameState, letter)
	}

	if currentGameState == nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}

	err := Temps.Temp.ExecuteTemplate(w, "jeux", currentGameState)
	if err != nil {
		log.Println("Erreur lors de l'exécution du template:", err)
		http.Error(w, "Erreur lors de l'exécution du template: "+err.Error(), http.StatusInternalServerError)
	}
}

func StartNewGame(w http.ResponseWriter, r *http.Request) {
	difficulty := r.URL.Query().Get("difficulty")
	var err error
	currentGameState, err = NewGame(difficulty)
	if err != nil {
		log.Println("Erreur lors de la création du jeu:", err)
		http.Error(w, "Erreur lors de la création du jeu: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func NewGame(difficulty string) (*GameState, error) {
	// Obtenir un mot aléatoire en fonction de la difficulté
	word, err := GetRandomWord(difficulty)
	if err != nil {
		return nil, err
	}
	alphabet := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

	maxTries := DetermineMaxTries(difficulty)

	return &GameState{
		Word:           word,
		DisplayWord:    strings.Repeat("_", len(word)),
		Tries:          0,
		MaxTries:       maxTries,
		GuessedLetters: make(map[string]bool),
		Alphabet:       alphabet,
		Difficulty:     difficulty,
	}, nil
}

func ProcessGuess(gameState *GameState, guess string) {
	if gameState == nil {
		return
	}

	if gameState.Errors >= gameState.MaxTries {
		return
	}

	if _, guessed := gameState.GuessedLetters[guess]; guessed {
		// La lettre a déjà été devinée, ne rien faire
		return
	}

	// Logique pour mettre à jour DisplayWord et Tries en fonction de la lettre devinée
	// Exemple simplifié :
	gameState.GuessedLetters[guess] = true

	correctGuess := false
	for i, char := range gameState.Word {
		if string(char) == guess {
			gameState.DisplayWord = gameState.DisplayWord[:i] + guess + gameState.DisplayWord[i+1:]
			correctGuess = true
		}
	}

	if !correctGuess {
		gameState.Errors++
	}

	if gameState.Errors >= gameState.MaxTries {
		gameState.GameOver = true
	}

}

func ProcessUserGuess(w http.ResponseWriter, r *http.Request) {
	if currentGameState == nil {
		return // ou initialiser un nouveau jeu si cela est approprié
	}

	letter := r.FormValue("letter")
	log.Println("Lettre devinée:", letter)
	ProcessGuess(currentGameState, letter)
}

func GetRandomWord(difficulty string) (string, error) {
	var filename string
	switch difficulty {
	case "facile":
		filename = "game/facile.txt"
	case "moyen":
		filename = "game/moyen.txt"
	case "difficile":
		filename = "game/difficile.txt"
	case "pays":
		filename = "game/pays.txt"
	case "capitales":
		filename = "game/capitales.txt" // Changer le nom du fichier pour la catégorie "capitales"
	case "animaux":
		filename = "game/animaux.txt" // Changer le nom du fichier pour la catégorie "animaux"
	case "musique":
		filename = "game/musique.txt" // Changer le nom du fichier pour la catégorie "musique"
	case "sport":
		filename = "game/sport.txt"
	default:
		return "", fmt.Errorf("invalid difficulty")
	}

	return GetRandomWordLowercase(filename)
}

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

func DetermineMaxTries(difficulty string) int {
	switch difficulty {
	case "facile":
		return 10
	case "moyen":
		return 7
	case "difficile":
		return 5
	case "pays":
		return 10
	case "capitales":
		return 10 // Changer la valeur si nécessaire pour la catégorie "capitales"
	case "animaux":
		return 10 // Changer la valeur si nécessaire pour la catégorie "animaux"
	case "musique":
		return 10 // Changer la valeur si nécessaire pour la catégorie "musique"
	case "sport":
		return 10 // Changer la valeur si nécessaire pour la catégorie "sport"
	default:
		return 10
	}
}
