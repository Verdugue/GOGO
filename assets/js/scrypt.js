function guessLetter(letter) {
    // Change le bouton en rouge pour indiquer qu'il a été utilisé
    document.querySelectorAll('.letter-button').forEach(button => {
        if(button.textContent === letter) {
            button.classList.add('used');
            button.disabled = true; // Désactive le bouton pour éviter de le presser à nouveau
        }
    });
    
    // Envoie la lettre au serveur
    const formData = new FormData();
    formData.append('letter', letter);
    fetch('/jeu', { method: 'POST', body: formData })
        .then(response => response.text())
        .then(html => {
            // Mettre à jour la page avec la nouvelle réponse du serveur
            document.open();
            document.write(html);
            document.close();
        });
}
