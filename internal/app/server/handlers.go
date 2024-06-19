package server

import (
	"SPORTALK/internal/model"
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
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
		execTmpl(w, templates.Lookup("registerPage.html"), nil)
	}
}

func (s *server) saveRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := model.RegisterPageData{}

		userName := r.FormValue("userName")
		email := r.FormValue("email")
		password := r.FormValue("password")
		rePassword := r.FormValue("rePassword")

		// Check if passwords match
		if password != rePassword {
			s.logger.Println("Passwords don't match")
			data.ErrorMsg = "Passwords don't match"
			execTmpl(w, templates.Lookup("registerPage.html"), data)
			return
		}

		err := s.store.User().ExistingUser(userName, email)
		if err != nil {
			s.logger.Println("error:", err)
			data.UserExistsErrorMsg = "User already exists in the system"
			execTmpl(w, templates.Lookup("registerPage.html"), data)
			return
		}

		user, err := model.NewUser(userName, email, password)
		if err != nil {
			s.logger.Println("NewUser() error: ", err)
			data.ErrorMsg = "Failed to create the user"
			execTmpl(w, templates.Lookup("registerPage.html"), data)
			return
		}

		if err = s.store.User().Register(user); err != nil {
			s.logger.Println("Register() error: ", err)
			data.ErrorMsg = "Failed to register the user"
			execTmpl(w, templates.Lookup("registerPage.html"), data)
			return
		}

		execTmpl(w, templates.Lookup("main.html"), nil)
	}
}
func (s *server) loginPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errorMessage := ""
		if errMsg := r.URL.Query().Get("error"); errMsg == "notfound" {
			errorMessage = "User not found. Please try again."
		}

		execTmpl(w, templates.Lookup("login.html"), map[string]string{"ErrorMessage": errorMessage})
	}
}

func (s *server) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Println("@ login page")
		email := r.FormValue("email")
		password := r.FormValue("password")
		s.logger.Println(email)

		user := &model.User{
			Email:    email,
			Password: password,
		}

		// Authenticate the user
		if err := s.store.User().Login(user); err != nil {
			// we changed here to include an error message in the query parameters
			s.logger.Println("Login() error: ", err)
			http.Redirect(w, r, "/loginPage?error=notfound", http.StatusSeeOther)
			return
		}

		// Create a new session for the user
		session, err := model.NewSession(user.UUID)
		if err != nil {
			s.logger.Println("NewSession() error: ", err)
			http.Redirect(w, r, "/loginPage", http.StatusInternalServerError)
			return
		}

		// Store the session in the database
		if err := s.store.Session().Create(session); err != nil {
			s.logger.Println("CreateSession() error: ", err)
			http.Redirect(w, r, "/loginPage", http.StatusInternalServerError)
			return
		}

		// Set a session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_uuid",
			Value:    session.SessionID,
			Expires:  session.ExpiresAt,
			HttpOnly: true,
			Secure:   false, // Set to true if you have HTTPS
		})

		// Redirect the user to their profile
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

func (s *server) serveStatic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))).ServeHTTP(w, r)
	}
}

func (s *server) home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get current user if exists
		var user *model.User
		if sessionCookie, err := r.Cookie("session_uuid"); err == nil {
			session, err := s.store.Session().GetByUUID(sessionCookie.Value)
			if err == nil {
				user, _ = s.store.User().GetByUUID(session.UserUUID)
			}
		}

		// Fetching all posts.
		posts, err := s.store.Post().GetAll()
		if err != nil {
			s.logger.Println("error fetching posts:", err)
			http.Error(w, "error fetching posts", http.StatusInternalServerError)
			return
		}

		for _, post := range posts {
			// Fetch user who created the post
			fetchedUser, _ := s.store.User().GetByUUID(post.UserID)
			post.User = fetchedUser

			// Fetch categories for each post
			categories, err := s.store.Post().GetCategories(post.ID)
			if err != nil {
				s.logger.Println("error fetching categories for post:", err)
				http.Error(w, "error fetching post categories", http.StatusInternalServerError)
				return
			}
			post.Categories = categories

			// Fetch comments with reactions for each post using the updated repository method
			comments, err := s.store.Comment().GetCommentsWithReactionsByPostID(post.ID)
			if err != nil {
				s.logger.Println("error fetching comments for post:", err)
				http.Error(w, "error fetching post comments", http.StatusInternalServerError)
				return
			}
			for _, comment := range comments {
				fetchedUser, _ := s.store.User().GetByUUID(comment.UserID)
				comment.User = fetchedUser
			}
			post.Comments = comments
		}

		// Fetching all categories.
		allCategories, err := s.store.Category().GetAll()
		if err != nil {
			s.logger.Println("error fetching categories:", err)
			http.Error(w, "error fetching categories", http.StatusInternalServerError)
			return
		}

		// Struct to pass into template.
		data := &model.PageData{
			User:       user,
			Posts:      posts,
			Categories: allCategories,
		}

		execTmpl(w, templates.Lookup("main.html"), data)
	}
}

// execTmpl renders a template with the given data or returns an internal server error.
func execTmpl(w http.ResponseWriter, tmpl *template.Template, data interface{}) {
	if err := tmpl.Execute(w, data); err != nil {
		log.Println("Error executing template:", err)
	}
}
func (s *server) createPostPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categories, err := s.store.Category().GetAll()
		if err != nil {
			http.Error(w, "Failed to load categories", http.StatusInternalServerError)
			return
		}

		errMsg := ""
		if errParam := r.URL.Query().Get("error"); errParam == "atleast_one_category_required" {
			errMsg = "At least one category must be selected."
		}

		data := struct {
			Categories   []*model.Category
			ErrorMessage string
		}{
			Categories:   categories,
			ErrorMessage: errMsg,
		}

		if err := templates.Lookup("createPostPage.html").Execute(w, data); err != nil {
			http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		}
	}
}

func (s *server) createPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get session cookie
		sessionCookie, err := r.Cookie("session_uuid")
		if err != nil {
			http.Redirect(w, r, "/loginPage", http.StatusSeeOther)
			return
		}

		// Fetch the session
		session, err := s.store.Session().GetByUUID(sessionCookie.Value)
		if err != nil {
			http.Redirect(w, r, "/loginPage", http.StatusSeeOther)
			return
		}

		// Get the user UUID from session
		userUUID := session.UserUUID

		subject := r.FormValue("postTitle")
		content := r.FormValue("postText")

		// Parse form for category checkboxes
		r.ParseForm()
		categoryIDs := r.PostForm["categoryIDs"]

		// validation for post category
		if len(categoryIDs) == 0 {
			http.Redirect(w, r, "/createPostPage?error=atleast_one_category_required", http.StatusSeeOther)
			return
		}

		post, err := model.NewPost(userUUID, subject, content)
		if err != nil {
			s.logger.Println("NewPost() error: ", err)
			http.Redirect(w, r, "/createPostPage", http.StatusSeeOther)
			return
		}

		if err = s.store.Post().Create(post); err != nil {
			s.logger.Println("Create() error: ", err)
			http.Redirect(w, r, "/createPostPage", http.StatusSeeOther)
			return
		}

		// add the category to the post
		for _, categoryIDStr := range categoryIDs {
			categoryID, err := strconv.Atoi(categoryIDStr)
			if err != nil {
				s.logger.Println("Error converting categoryID to int: ", err)
				http.Redirect(w, r, "/createPostPage", http.StatusSeeOther)
				return
			}

			if err := s.store.Post().AddCategoryToPost(post.ID, categoryID); err != nil {
				s.logger.Println("Error adding category to post: ", err)
				http.Redirect(w, r, "/createPostPage", http.StatusSeeOther)
				return
			}
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (s *server) createCategoryPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errorMessage := ""
		switch r.URL.Query().Get("error") {
		case "checkError":
			errorMessage = "Failed to check category existence."
		case "categoryExists":
			errorMessage = "This category already exists."
		}
		execTmpl(w, templates.Lookup("createCategory.html"), map[string]string{"ErrorMessage": errorMessage})
	}
}

func (s *server) createCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categoryName := r.FormValue("categoryName")
		exists, err := s.store.Category().Exists(categoryName)
		if err != nil {
			s.logger.Println("Check category existence error: ", err)
			http.Redirect(w, r, "/createCategoryPage?error=checkError", http.StatusSeeOther)
			return
		}

		if exists {
			s.logger.Println("Category already exists: ", categoryName)
			http.Redirect(w, r, "/createCategoryPage?error=categoryExists", http.StatusSeeOther)
			return
		}

		category := &model.Category{Name: categoryName}
		if err := s.store.Category().Create(category); err != nil {
			s.logger.Println("Create category error: ", err)
			http.Redirect(w, r, "/createCategoryPage", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (s *server) categoryPosts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get current user if exists
		var user *model.User
		if sessionCookie, err := r.Cookie("session_uuid"); err == nil {
			session, err := s.store.Session().GetByUUID(sessionCookie.Value)
			if err == nil {
				user, _ = s.store.User().GetByUUID(session.UserUUID)
			}
		}

		categoryIDStr := strings.TrimPrefix(r.URL.Path, "/category/")
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			s.logger.Println("Error converting category ID:", err)
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		// Get all posts in the category.
		posts, err := s.store.Post().GetByCategory(categoryID)
		if err != nil {
			s.logger.Println("Error fetching posts:", err)
			http.Error(w, "error fetching posts by category", http.StatusInternalServerError)
			return
		}

		for _, post := range posts {
			// Fetch user who created the post
			fetchedUser, _ := s.store.User().GetByUUID(post.UserID)
			post.User = fetchedUser

			// Fetch categories for each post
			categories, err := s.store.Post().GetCategories(post.ID)
			if err != nil {
				s.logger.Println("error fetching categories for post:", err)
				http.Error(w, "error fetching post categories", http.StatusInternalServerError)
				return
			}
			post.Categories = categories

			// Fetch comments with reactions for each post using the updated repository method
			commentsWithReactions, err := s.store.Comment().GetCommentsWithReactionsByPostID(post.ID)
			if err != nil {
				s.logger.Println("error fetching comments with reactions for post:", err)
				http.Error(w, "error fetching post comments with reactions", http.StatusInternalServerError)
				return
			}
			for _, comment := range commentsWithReactions {
				// Fetch user who created the comment
				fetchedUser, _ := s.store.User().GetByUUID(comment.UserID)
				comment.User = fetchedUser
			}

			post.Comments = commentsWithReactions
		}

		// Fetch all categories except currently used
		allCategories, err := s.store.Category().GetAll()
		if err != nil {
			s.logger.Println("Error fetching categories:", err)
			http.Error(w, "error fetching all categories", http.StatusInternalServerError)
			return
		}

		// Pass data to template
		data := &model.PageData{
			User:       user,
			Posts:      posts,
			Categories: allCategories,
		}

		execTmpl(w, templates.Lookup("home.html"), data)
	}
}

func (s *server) registerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userName := r.FormValue("username")
		password := r.FormValue("password")
		email := r.FormValue("email")

		if userName == "" || password == "" || email == "" {
			data := struct {
				ErrorMsg string
			}{
				ErrorMsg: "All fields must be provided",
			}
			execTmpl(w, templates.Lookup("/registerPage.html"), data)
			return
		}

		err := s.store.User().ExistingUser(userName, email)
		if err != nil {
			s.logger.Println("redirect - error:", err)
			data := struct {
				UserExistsErrorMsg string
			}{
				UserExistsErrorMsg: "User already exists in the system",
			}
			execTmpl(w, templates.Lookup("/registerPage.html"), data)
			return
		}

		user, err := model.NewUser(userName, email, password)
		if err != nil {
			s.logger.Println("NewUser() error: ", err)
			http.Redirect(w, r, "/registerPage", http.StatusSeeOther)
			return
		}

		if err = s.store.User().Register(user); err != nil {
			s.logger.Println("Register() error: ", err)
			http.Redirect(w, r, "/registerPage", http.StatusSeeOther)
			return
		}

		// Create a new session for the user
		session, err := model.NewSession(user.UUID)
		if err != nil {
			http.Error(w, "failed to create session", http.StatusInternalServerError)
			return
		}
		if err := s.store.Session().Create(session); err != nil {
			http.Error(w, "failed to store session", http.StatusInternalServerError)
			return
		}

		// Set a cookie for the session
		cookie := &http.Cookie{Name: "session_uuid", Value: session.SessionID,
			Expires: session.ExpiresAt, HttpOnly: true}
		http.SetCookie(w, cookie)

		// Redirect the user to the their profile
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

func (s *server) serveUserProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract session cookie
		sessionCookie, err := r.Cookie("session_uuid")
		if err != nil {
			http.Redirect(w, r, "/loginPage", http.StatusSeeOther)
			return
		}

		// Fetch the session
		session, err := s.store.Session().GetByUUID(sessionCookie.Value)
		if err != nil {
			http.Redirect(w, r, "/loginPage", http.StatusSeeOther)
			return
		}

		// Fetch the user
		user, err := s.store.User().GetByUUID(session.UserUUID)
		if err != nil {
			http.Redirect(w, r, "/loginPage", http.StatusSeeOther)
			return
		}

		// Render template with user data
		execTmpl(w, templates.Lookup("userProfilePage.html"), user)
	}
}

func (s *server) logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie("session_uuid")
		if err != nil {
			http.Redirect(w, r, "/loginPage", http.StatusSeeOther)
			return
		}

		// Delete the session from the DB
		err = s.store.Session().Delete(session.Value)
		if err != nil {
			http.Error(w, "Failed to end session", http.StatusInternalServerError)
			return
		}

		// Delete the session cookie
		session.MaxAge = -1
		http.SetCookie(w, session)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (s *server) createComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get session cookie
		sessionCookie, err := r.Cookie("session_uuid")
		if err != nil {
			http.Redirect(w, r, "/loginPage", http.StatusSeeOther)
			return
		}

		// Fetch the session
		session, err := s.store.Session().GetByUUID(sessionCookie.Value)
		if err != nil {
			http.Redirect(w, r, "/loginPage", http.StatusSeeOther)
			return
		}

		// Get the user UUID from session
		userUUID := session.UserUUID

		// Get post ID from form
		postID := r.FormValue("postID")

		// Get comment text from form
		commentTxt := r.FormValue("commentText")

		// Create new comment
		comment, err := model.NewComment(postID, userUUID, commentTxt)
		if err != nil {
			s.logger.Println("NewComment() error: ", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Save the comment
		if err = s.store.Comment().Create(comment); err != nil {
			s.logger.Println("CreateComment() error: ", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Redirect back to homepage
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (s *server) handleCreatePostReaction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		sessionCookie, err := r.Cookie("session_uuid")
		if err != nil {
			http.Error(w, "Session error", http.StatusBadRequest)
			return
		}
		session, err := s.store.Session().GetByUUID(sessionCookie.Value)
		if err != nil {
			http.Error(w, "Session retrieval error", http.StatusInternalServerError)
			return
		}
		userUUID := session.UserUUID
		postID := r.FormValue("postID")
		reactionTypeStr := r.FormValue("reactionType")

		var reactionID int
		switch reactionTypeStr {
		case "like":
			reactionID = 1
		case "dislike":
			reactionID = 2
		default:
			http.Error(w, "Invalid reaction type", http.StatusBadRequest)
			return
		}

		existingReaction, err := s.store.Reaction().GetUserPostReaction(userUUID, postID)
		if err != nil && err != sql.ErrNoRows {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if existingReaction != nil {
			if existingReaction.ReactionID == reactionID {
				if err := s.store.Reaction().DeletePostReaction(userUUID, postID); err != nil {
					http.Error(w, "Failed to delete reaction", http.StatusInternalServerError)
					return
				}
			} else {
				if err := s.store.Reaction().DeletePostReaction(userUUID, postID); err != nil {
					http.Error(w, "Failed to delete existing reaction", http.StatusInternalServerError)
					return
				}

				reaction := &model.Reaction{
					UserID:     userUUID,
					PostID:     postID,
					ReactionID: reactionID,
				}
				if err := s.store.Reaction().CreatePostReaction(reaction); err != nil {
					http.Error(w, "Failed to create reaction", http.StatusInternalServerError)
					return
				}
			}
		} else {
			reaction := &model.Reaction{
				UserID:     userUUID,
				PostID:     postID,
				ReactionID: reactionID,
			}
			if err := s.store.Reaction().CreatePostReaction(reaction); err != nil {
				http.Error(w, "Failed to create reaction", http.StatusInternalServerError)
				return
			}
		}

		http.Redirect(w, r, "/home?postID="+postID, http.StatusSeeOther)
	}
}

func (s *server) handleCreateCommentReaction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		sessionCookie, err := r.Cookie("session_uuid")
		if err != nil {
			http.Error(w, "Session error", http.StatusBadRequest)
			return
		}

		session, err := s.store.Session().GetByUUID(sessionCookie.Value)
		if err != nil {
			http.Error(w, "Session retrieval error", http.StatusInternalServerError)
			return
		}

		userUUID := session.UserUUID
		commentID := r.FormValue("commentID")
		reactionTypeStr := r.FormValue("reactionType")

		var reactionID int
		switch reactionTypeStr {
		case "like":
			reactionID = 1
		case "dislike":
			reactionID = 2
		default:
			http.Error(w, "Invalid reaction type", http.StatusBadRequest)
			return
		}

		existingReaction, err := s.store.Reaction().GetUserCommentReaction(userUUID, commentID)
		if err != nil && err != sql.ErrNoRows {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if existingReaction != nil {
			if existingReaction.ReactionID == reactionID {
				if err := s.store.Reaction().DeleteCommentReaction(userUUID, commentID); err != nil {
					http.Error(w, "Failed to delete reaction", http.StatusInternalServerError)
					return
				}
			} else {
				if err := s.store.Reaction().DeleteCommentReaction(userUUID, commentID); err != nil {
					http.Error(w, "Failed to delete existing reaction", http.StatusInternalServerError)
					return
				}

				reaction := &model.Reaction{
					UserID:     userUUID,
					CommentID:  commentID,
					ReactionID: reactionID,
				}
				if err := s.store.Reaction().CreateCommentReaction(reaction); err != nil {
					http.Error(w, "Failed to create reaction", http.StatusInternalServerError)
					return
				}
			}
		} else {
			reaction := &model.Reaction{
				UserID:     userUUID,
				CommentID:  commentID,
				ReactionID: reactionID,
			}
			if err := s.store.Reaction().CreateCommentReaction(reaction); err != nil {
				http.Error(w, "Failed to create reaction", http.StatusInternalServerError)
				return
			}
		}

		// Redirect to referer or to a default route if not available
		referer := r.Header.Get("Referer")
		if referer == "" {
			referer = "/home"
		}
		http.Redirect(w, r, referer, http.StatusSeeOther)
	}
}
