package database

type Joueur struct {
	Pseudo        string
	Difficulty    string
	Mot           string
	ATrouver      string
	Win           int
	Essai         int
	LetterTry     []string
	LettreVisible []string
	Display       string
	Img           string
}
