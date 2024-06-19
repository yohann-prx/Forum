package server

import (
	"SPORTALK/internal/model"
	"fmt"
	"html/template"
	"net/http"
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
