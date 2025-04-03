package delivery

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"certificate/internal/ports"
)

type HTTPServer struct {
	svc ports.RegistrationService
}

func NewHTTPServer(svc ports.RegistrationService) *HTTPServer {
	return &HTTPServer{svc: svc}
}

// Запуск сервера
func (s *HTTPServer) Start(port string) {
	http.HandleFunc("GET /register", s.HandleRegister)
	http.HandleFunc("POST /submit", s.HandleSubmit)

	log.Println("Запуск HTTP-сервера на", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// Загрузка статических файлов
func (s *HTTPServer) ServeStaticFiles() {
	fs := http.FileServer(http.Dir("templates/styles"))
	http.Handle("/styles/", http.StripPrefix("/styles/", fs))
	ft := http.FileServer(http.Dir("templates/fonts"))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", ft))
}

// Обработчик регистрации(когда перешли по ссылке)
func (s *HTTPServer) HandleRegister(w http.ResponseWriter, r *http.Request) {
	encryptedToken := r.URL.Query().Get("token")
	if encryptedToken == "" {
		http.Error(w, "Token is missing", http.StatusBadRequest)
		return
	}

	token, err := s.svc.ValidateAndDecode(encryptedToken)
	if err != nil {
		s.renderPage(w, "error.html", map[string]string{"Message": "Invalid token"})
		return
	}

	s.renderPage(w, "register.html", map[string]string{"Token": token})
}

// Вспомогательный метод для рендера HTML-шаблонов
func (s *HTTPServer) renderPage(w http.ResponseWriter, templateName string, data map[string]string) {
	tmpl, err := template.ParseFiles("templates/" + templateName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading template: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, data)
}

// Обработчик регистрации(после того, как нажали сабмит)
func (s *HTTPServer) HandleSubmit(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	name := r.FormValue("name")
	phone := r.FormValue("phone")
	birthday := r.FormValue("birthday")

	if token == "" || name == "" || phone == "" || birthday == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	err := s.svc.RegisterUser(token, name, phone, birthday)
	if err != nil {
		http.Error(w, "Registration failed", http.StatusBadRequest)
		return
	}

	s.renderPage(w, "success.html", map[string]string{"Message": "Registration successful!"})
}
