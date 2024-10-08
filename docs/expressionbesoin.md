# EXPRESSION DE BESOINS - ORCHESTRATOR v0.3.0 Alpha

## 1. Commun à tout
- Tout doit être loggé dans un logger qui exporte ses logs dans des fichiers textes pour info, debug, warning, error.
- Le logger s'affiche aussi sur la console pour le serveur ou la CLI go.
- Tout le texte du logiciel est en anglais.
- Découpe proprement tout le logiciel pour qu'un fichier n'excède jamais 3000 tokens.
- Le logiciel est versionné, on commencera à la version "0.3.0 Alpha".
- Journalisation détaillée (debug, info, warning, error) pour toutes les actions.
- Possibilité de changer le niveau de log en cours d'exécution.
- Structure du projet séparée en serveur et client CLI.
- Utilisation de Cobra pour la gestion des lignes de commande.
- Mode silencieux (--silent) pour désactiver la sortie console des logs.
- Mode debug (--debug) pour activer des logs très verbeux.
- Configuration centralisée via fichier YAML, avec possibilité de surcharge par ligne de commande.

## 2. Fonctionnalités principales

### 2.1 Structure du projet
- Séparation claire entre le serveur (orchestrator-server) et le client CLI (orchestrator-cli).
- Structure modulaire avec des fichiers ne dépassant pas 3000 tokens.
- Utilisation de Cobra pour la gestion des commandes CLI.
- Configuration centralisée via fichier YAML, avec possibilité de surcharge par ligne de commande.

### 2.2 Gestion des jobs
- Création de jobs avec paramètres (nom, commande, répertoire de travail, timeout, nombre max de tentatives).
- Listage de tous les jobs.
- Récupération des détails d'un job spécifique.
- Exécution d'un job spécifique.
- Sauvegarde du résultat dans un contexte associé à l'ID du job.
- Une ID de job commence toujours par J.
- Durée d'exécution maximale paramétrable par job (timeout).
- Gestion des erreurs et reprises automatiques.
- Exécution de commandes système avec paramètres; lorsqu'une commande est exécutée, le répertoire d'exécution doit être donné dans la définition d'un job. Un change dir doit être fait pour exécuter le job.
- Possibilité d'attribuer un nom aux jobs pour une meilleure identification.

### 2.3 Gestion des pipelines
- Définition de séquences de jobs.
- Une ID de Pipeline commence toujours par P.
- Planification de l'exécution à une date et heure spécifiées.
- Agrégation des contextes des jobs dans le contexte du pipeline.
- Le dernier job produit le contexte final du pipeline.
- Capacité de gérer plusieurs dizaines de jobs et pipelines simultanément.
- Gestion des dépendances entre les jobs.
- Possibilité d'attribuer un nom aux pipelines pour une meilleure identification.

### 2.4 Gestion des schedules
- Le but est d'ordonnancer les jobs et les pipelines.
- Sur le principe, fonctionnement comme un cron pour la définition du scheduling.
- Un schedule peut s'occuper de plusieurs jobs ou pipeline qui du coup devront être lancés en même temps. Prévoir une exécution parallèle au niveau du serveur.
- Le schedule log ses actions dans le logger.
- Possibilité d'attribuer un nom au schedule pour une meilleure identification.
- Un ID de schedule commence toujours par S.
- Il faut qu'on puisse affecter des jobs ou des pipelines à des schedules.

### 2.5 Contextes
- Un contexte est le résultat de l'exécution d'un job ou d'un pipeline.
- C'est littéralement la sortie console de l'exécution du job ou du pipeline.
- Un contexte est attaché à l'exécution d'un job ou d'un pipeline.
- Par défaut il est vide tant que le job ou le pipeline n'a jamais été exécuté.
- Lorsque son job a été exécuté, le contexte associé à un job contient le résultat de la sortie standard de la dernière exécution.
- Lorsque son pipeline a été exécuté, le contexte associé à un pipeline contient l'agrégation des contextes des jobs.
- Le dernier job produit le contexte final du pipeline.
- Il a un ID propre. Cet ID commence par CJ si c'est un contexte lié à un job ou CP si c'est un contexte lié à un pipeline.
- Tous les pipelines ou les jobs peuvent utiliser un contexte dans leur ligne d'exécution, en utilisant la commande $$(nomducontexte).
- L'API doit prévoir que l'on puisse consulter, réinitialiser un contexte.
- Associé à un contexte, on a les détails de sa génération (l'id du job ou du pipeline, la date et l'heure, la durée, ... tout ce qui est important).

### 2.6 Interface utilisateur
- Interface en mode HTML5 / JS utilisant Vue.js avec Tailwind CSS.
- Affichage des logs en temps réel avec possibilité de filtrage par niveau.
- Visualisation de l'état des jobs et pipelines.
- Visualisation des contextes.
- Visualisation des schedules.
- Navigation intuitive entre les différentes vues (jobs, pipelines, logs, détails).
- Champ de saisie de commandes occupant toute la largeur de l'écran.
- Design responsive et élégant.
- Considérations d'accessibilité intégrées dans le développement.
- Visualisation de données simple et élégante (graphiques, tableaux de bord).
- Interface exclusivement en anglais.

### 2.7 API REST
- Contrôle et surveillance à distance du système.
- Système d'authentification interne avec génération de clés API.
- Endpoints pour la gestion des jobs, pipelines, et la récupération des logs.
- Tous les verbes traditionnels : Add, Remove, Copy, ...
- Documentation de l'API avec Swagger.
- Endpoints implémentés :
  - POST /jobs : Créer un nouveau job
  - GET /jobs : Lister tous les jobs
  - GET /jobs/:id : Obtenir les détails d'un job spécifique
  - POST /jobs/:id/run : Exécuter un job spécifique

### 2.8 Stockage et persistance
- Utilisation de BoltDB pour la persistance des données.
- Sauvegarde persistante des contextes et résultats des jobs.
- Gestion de l'historique des exécutions.
- Configuration de la durée de rétention des données.
- Gestion de l'export/import de la configuration pour une sauvegarde efficace.

### 2.9 Configuration
- Utilisation de fichiers YAML pour la configuration globale de l'application.
- Configuration flexible des timeouts, des tentatives de reprise, et autres paramètres.
- Possibilité de surcharger la configuration via des options de ligne de commande.

### 2.10 Monitoring et reporting
- Intégration de métriques pour le suivi des performances (cpu, ram, disque, io, nb de job/pipeline par état) dans l'API.
- Accès par la CLI pour supervision.
- Affichage élégant en HTML des statistiques dans l'interface utilisateur web.
- Intégration d'un système de notifications par webhook pour les erreurs critiques et la fin d'exécution des jobs/pipelines.

### 2.11 Fiabilité
- Mécanisme de reprise sur erreur pour les pipelines et les jobs individuels.
- Gestion robuste des erreurs et des cas limites.

### 2.12 Sécurité
- Authentification et autorisation pour l'accès à l'API REST.
- Chiffrement des données sensibles dans la base de données.
- Système de gestion des utilisateurs avec trois niveaux d'accès :
  - Admin : accès complet, peut tout modifier
  - Writer : peut modifier les jobs, pipelines et schedules
  - Reader : accès en lecture seule, ne peut rien modifier

### 2.13 Gestion des ressources
- Implémentation d'un système de gestion des ressources pour limiter le nombre de jobs simultanés en fonction des ressources disponibles (prévu pour une implémentation future).

### 2.14 Audit et logging
- Mise en place d'un système de logging détaillé pour toutes les actions effectuées dans le système.
- Les logs doivent inclure l'utilisateur, l'action, la date/heure et les détails pertinents.
- Support des modes silencieux et debug pour contrôler le niveau de verbosité des logs.
- Les logs sont exportés dans des fichiers textes et affichés sur la console.
- Mode silencieux (--silent) pour désactiver la sortie console des logs.
- Mode debug (--debug) pour activer des logs très verbeux.
- Possibilité de changer le niveau de log en cours d'exécution.

## 3. Contraintes techniques
- Développement du serveur et de la CLI en Go.
- Utilisation de Cobra pour la gestion des lignes de commande.
- Interface web en Vue.js avec Tailwind CSS.
- Base de données BoltDB pour la persistance.
- Compatibilité multi-plateforme : Windows, Linux, macOS (priorité sur le déploiement serveur Linux).
- Optimisation des performances pour gérer efficacement 1000 jobs et 100 pipelines par instance.

## 4. Livrables
- Code source du logiciel.
- Binaires exécutables pour Windows, Linux, et macOS.
- Documentation technique détaillée générée à partir du code.
- Documentation utilisateur pour l'interface web et la CLI.
- Tests unitaires et d'intégration.
- Scripts de déploiement et de configuration.
- Script PowerShell pour la compilation des binaires et la préparation du package de distribution.

## 5. Critères de succès
- Performance : temps de réponse rapide, gestion efficace de 1000 jobs et 100 pipelines par instance.
- Fiabilité : gestion robuste des erreurs, pas de perte de données.
- Extensibilité : facilité d'ajout de nouvelles fonctionnalités via le système de plugins.
- Utilisabilité : interface utilisateur intuitive, responsive et accessible.
- Sécurité : authentification et autorisation robustes pour l'API REST.
- Flexibilité : capacité à s'adapter à différents environnements et cas d'utilisation.
- Facilité de déploiement et de mise à jour.

## 6. Fonctionnalités futures potentielles (non prioritaires pour la v0.3.0 Alpha)
- Versionnement des jobs et pipelines.
- Système de plugins.
- Gestion multi-environnements (dev, staging, prod).
- Système de quotas et limites par utilisateur.
- Support multilingue de l'interface utilisateur.