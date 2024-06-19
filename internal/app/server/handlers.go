package server

import (
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
