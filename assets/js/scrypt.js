function guessLetter(letter) {
    // Change le bouton en rouge pour indiquer qu'il a été utilisé
    document.querySelectorAll('.letter-button').forEach(button => {
        if (button.textContent === letter) {
            button.classList.add('used'); // Ajoute une classe CSS pour changer le style du bouton
            button.disabled = true;       // Désactive le bouton pour éviter de le presser à nouveau
        }
    });
    
    // Envoie la lettre au serveur
    const formData = new FormData();      // Crée un objet FormData pour faciliter l'envoi de données
    formData.append('letter', letter);    // Ajoute la lettre sélectionnée à formData

    // Effectue une requête POST au serveur avec la lettre choisie
    fetch('/jeu', { method: 'POST', body: formData })
        .then(response => response.text())  // Convertit la réponse en texte (HTML)
        .then(html => {
            // Met à jour la page avec la nouvelle réponse du serveur
            document.open();                // Ouvre le document pour la modification
            document.write(html);           // Écrit le nouveau HTML dans le document
            document.close();               // Ferme le document
        });
}
