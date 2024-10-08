# build.ps1

# Définition des variables
$serverBinary = "orchestrator-server.exe"
$clientBinary = "orchestrator-cli.exe"
$configFile = "config.yaml"
$outputDir = "binary"

# Création du dossier de build s'il n'existe pas
if (!(Test-Path -Path $outputDir)) {
    New-Item -ItemType Directory -Force -Path $outputDir
}

# Fonction pour afficher les messages avec un timestamp
function Log-Message {
    param([string]$message)
    Write-Host "$(Get-Date -Format 'yyyy-MM-dd HH:mm:ss') - $message"
}

# Compilation du serveur
Log-Message "Compilation du serveur..."
go build -o "$outputDir\$serverBinary" .\cmd\orchestrator-server\main.go
if ($LASTEXITCODE -ne 0) {
    Log-Message "Erreur lors de la compilation du serveur"
    exit 1
}
Log-Message "Serveur compilé avec succès"

# Compilation du client
Log-Message "Compilation du client..."
go build -o "$outputDir\$clientBinary" .\cmd\orchestrator-cli\main.go
if ($LASTEXITCODE -ne 0) {
    Log-Message "Erreur lors de la compilation du client"
    exit 1
}
Log-Message "Client compilé avec succès"

# Copie du fichier de configuration
Log-Message "Copie du fichier de configuration..."
Copy-Item "config\$configFile" -Destination "$outputDir\$configFile"
Log-Message "Fichier de configuration copié"

# Création du dossier logs dans le dossier de build
New-Item -ItemType Directory -Force -Path "$outputDir\logs"
Log-Message "Dossier logs créé dans le dossier de build"

Log-Message "Build terminé avec succès"
Log-Message "Les binaires et la configuration se trouvent dans le dossier '$outputDir'"