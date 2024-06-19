package server

import (
	"SPORTALK/internal/model"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var templates = template.Must(template.ParseGlob("./web/templates/*.html"))

func (s *server) HandlePaths() {
	s.router.Handle("/static/", s.serveStatic())
	s.router.HandleFunc("/", s.home())
	s.router.HandleFunc("/registerPage", s.registerPage())
	s.router.HandleFunc("/saveUser", s.saveRegister())
	s.router.HandleFunc("/loginPage", s.loginPage())
	s.router.HandleFunc("/login", s.login())
	s.router.HandleFunc("/createPost", s.createPost())
	s.router.HandleFunc("/createPostPage", s.createPostPage())
	s.router.HandleFunc("/createCategory", s.createCategory())
	s.router.HandleFunc("/createCategoryPage", s.createCategoryPage())
	s.router.HandleFunc("/category/", s.categoryPosts())
	s.router.HandleFunc("/userProfilePage", s.serveUserProfile())
	s.router.HandleFunc("/logout", s.logout())
	s.router.HandleFunc("/createComment", s.createComment())
	s.router.HandleFunc("/createPostReaction", s.handleCreatePostReaction())
	s.router.HandleFunc("/reactComment", s.handleCreateCommentReaction())
}

func (s *server) registerPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Vérifier la méthode HTTP
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Exécuter le template
		err := execTmpl(w, templates.Lookup("registerPage.html"), nil)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

// Fonction utilitaire pour exécuter le template avec gestion des erreurs
func execTmpl(w http.ResponseWriter, tmpl *template.Template, data interface{}) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tmpl.Execute(w, data)
	if err != nil {
		return fmt.Errorf("template execution failed: %v", err)
	}
	return nil
}

func (s *server) saveRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := model.RegisterPageData{}

		// Vérification de la méthode HTTP
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Récupération des valeurs des champs du formulaire
		userName := r.FormValue("userName")
		email := r.FormValue("email")
		password := r.FormValue("password")
		rePassword := r.FormValue("rePassword")

		// Vérification si les mots de passe correspondent
		if password != rePassword {
			s.logger.Println("Passwords don't match")
			data.ErrorMsg = "Passwords don't match"
			execTmpl(w, templates.Lookup("registerPage.html"), data)
			return
		}

		// Vérification si l'utilisateur existe déjà
		err := s.store.User().ExistingUser(userName, email)
		if err != nil {
			s.logger.Println("ExistingUser() error:", err)
			data.UserExistsErrorMsg = "User already exists in the system"
			execTmpl(w, templates.Lookup("registerPage.html"), data)
			return
		}

		// Création d'un nouvel utilisateur
		user, err := model.NewUser(userName, email, password)
		if err != nil {
			s.logger.Println("NewUser() error:", err)
			data.ErrorMsg = "Failed to create the user"
			execTmpl(w, templates.Lookup("registerPage.html"), data)
			return
		}

		// Enregistrement de l'utilisateur
		if err = s.store.User().Register(user); err != nil {
			s.logger.Println("Register() error:", err)
			data.ErrorMsg = "Failed to register the user"
			execTmpl(w, templates.Lookup("registerPage.html"), data)
			return
		}

		// Redirection vers une page principale après l'inscription réussie
		http.Redirect(w, r, "/main", http.StatusSeeOther)
	}
}

func (s *server) loginPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Vérification de la méthode HTTP
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Déterminer le message d'erreur à afficher
		errorMessage := ""
		errMsg := r.URL.Query().Get("error")
		switch errMsg {
		case "notfound":
			errorMessage = "User not found. Please try again."
		case "invalid":
			errorMessage = "Invalid username or password. Please try again."
		default:
			// Aucune erreur spécifiée dans l'URL
		}

		// Exécuter le template avec les données nécessaires
		data := map[string]string{"ErrorMessage": errorMessage}
		execTmpl(w, templates.Lookup("login.html"), data)
	}
}

func (s *server) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Vérification de la méthode HTTP
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Récupération des valeurs du formulaire
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Création d'une instance utilisateur avec les informations de connexion
		user := &model.User{
			Email:    email,
			Password: password,
		}

		// Authentification de l'utilisateur
		err := s.store.User().Login(user)
		if err != nil {
			s.logger.Println("Login() error:", err)
			// Redirection vers la page de connexion avec un message d'erreur approprié
			http.Redirect(w, r, "/loginPage?error=notfound", http.StatusSeeOther)
			return
		}

		// Création d'une nouvelle session pour l'utilisateur
		session, err := model.NewSession(user.UUID)
		if err != nil {
			s.logger.Println("NewSession() error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Enregistrement de la session dans la base de données
		err = s.store.Session().Create(session)
		if err != nil {
			s.logger.Println("CreateSession() error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Définition d'un cookie de session sécurisé
		http.SetCookie(w, &http.Cookie{
			Name:     "session_uuid",
			Value:    session.SessionID,
			Expires:  session.ExpiresAt,
			HttpOnly: true,
			Secure:   false, // Modifier à true si vous utilisez HTTPS
			Path:     "/",
		})

		// Redirection de l'utilisateur vers sa page d'accueil
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

func (s *server) serveStatic() http.HandlerFunc {
	fs := http.FileServer(http.Dir("./web/static"))

	return func(w http.ResponseWriter, r *http.Request) {
		// Vérification de la méthode HTTP
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Sécurité : Vérifier que le chemin commence bien par "/static/"
		if !strings.HasPrefix(r.URL.Path, "/static/") {
			http.NotFound(w, r)
			return
		}

		// Utilisation de http.StripPrefix pour servir les fichiers statiques
		http.StripPrefix("/static/", fs).ServeHTTP(w, r)
	}
}

func (s *server) home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Vérification de la méthode HTTP
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Récupération de l'utilisateur actuel s'il existe
		var user *model.User
		if sessionCookie, err := r.Cookie("session_uuid"); err == nil {
			session, err := s.store.Session().GetByUUID(sessionCookie.Value)
			if err == nil {
				user, _ = s.store.User().GetByUUID(session.UserUUID)
			}
		}

		// Récupération de tous les posts
		posts, err := s.store.Post().GetAll()
		if err != nil {
			s.logger.Println("error fetching posts:", err)
			http.Error(w, "error fetching posts", http.StatusInternalServerError)
			return
		}

		// Récupération des informations complémentaires pour chaque post
		for _, post := range posts {
			// Récupération de l'utilisateur qui a créé le post
			fetchedUser, err := s.store.User().GetByUUID(post.UserID)
			if err != nil {
				s.logger.Printf("error fetching user for post %s: %v\n", post.ID, err)
				http.Error(w, "error fetching post user", http.StatusInternalServerError)
				return
			}
			post.User = fetchedUser

			// Récupération des catégories pour chaque post
			categories, err := s.store.Post().GetCategories(post.ID)
			if err != nil {
				s.logger.Printf("error fetching categories for post %s: %v\n", post.ID, err)
				http.Error(w, "error fetching post categories", http.StatusInternalServerError)
				return
			}
			post.Categories = categories

			// Récupération des commentaires avec réactions pour chaque post
			comments, err := s.store.Comment().GetCommentsWithReactionsByPostID(post.ID)
			if err != nil {
				s.logger.Printf("error fetching comments for post %s: %v\n", post.ID, err)
				http.Error(w, "error fetching post comments", http.StatusInternalServerError)
				return
			}
			for _, comment := range comments {
				fetchedUser, err := s.store.User().GetByUUID(comment.UserID)
				if err != nil {
					s.logger.Printf("error fetching user for comment %s: %v\n", comment.ID, err)
					http.Error(w, "error fetching comment user", http.StatusInternalServerError)
					return
				}
				comment.User = fetchedUser
			}
			post.Comments = comments
		}

		// Récupération de toutes les catégories
		allCategories, err := s.store.Category().GetAll()
		if err != nil {
			s.logger.Println("error fetching categories:", err)
			http.Error(w, "error fetching categories", http.StatusInternalServerError)
			return
		}

		// Création des données à passer au template
		data := &model.PageData{
			User:       user,
			Posts:      posts,
			Categories: allCategories,
		}

		// Exécution du template avec les données
		execTmpl(w, templates.Lookup("main.html"), data)
	}
}

func execTmpl(w http.ResponseWriter, tmpl *template.Template, data interface{}) {
	// Utilisation de Recover pour capturer les paniques potentielles lors de l'exécution du template
	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic recovered in execTmpl:", r)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}()

	// Exécution du template et gestion des erreurs
	err := tmpl.Execute(w, data)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (s *server) createPostPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Vérification de la méthode HTTP
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Récupération de toutes les catégories
		categories, err := s.store.Category().GetAll()
		if err != nil {
			s.logger.Println("Failed to load categories:", err)
			http.Error(w, "Failed to load categories", http.StatusInternalServerError)
			return
		}

		// Gestion des erreurs basée sur les paramètres d'URL
		errMsg := ""
		if errParam := r.URL.Query().Get("error"); errParam == "atleast_one_category_required" {
			errMsg = "At least one category must be selected."
		}

		// Données à passer au template
		data := struct {
			Categories   []*model.Category
			ErrorMessage string
		}{
			Categories:   categories,
			ErrorMessage: errMsg,
		}

		// Exécution du template avec les données
		if err := templates.Lookup("createPostPage.html").Execute(w, data); err != nil {
			s.logger.Println("Failed to execute template:", err)
			http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		}
	}
}
