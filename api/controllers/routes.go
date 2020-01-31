package controllers

import "github.com/robbie-thomas/fullstack/api/middlewares"

func (server *Server) initialiseRoutes() {

	// Home Route
	server.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(server.Home)).Methods("GET")

	// Login Route
	server.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(server.Login)).Methods("POST")

	//Users routes
	server.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(server.CreateUser)).Methods("POST")
	server.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(server.GetUsers)).Methods("GET")
	server.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(server.GetUser)).Methods("GET")
	server.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.UpdateUser))).Methods("PUT")
	server.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(server.DeleteUser)).Methods("DELETE")

	//Posts routes
	server.Router.HandleFunc("/posts", middlewares.SetMiddlewareJSON(server.CreatePost)).Methods("POST")
	server.Router.HandleFunc("/posts", middlewares.SetMiddlewareJSON(server.GetPosts)).Methods("GET")
	server.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareJSON(server.getPost)).Methods("GET")
	server.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.UpdatePost))).Methods("PUT")
	server.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareAuthentication(server.DeletePost)).Methods("DELETE")

	server.Router.HandleFunc("/spaces", middlewares.SetMiddlewareJSON(server.CreateSpace)).Methods("POST")
	server.Router.HandleFunc("/spaces", middlewares.SetMiddlewareJSON(server.GetSpaces)).Methods("GET")
	server.Router.HandleFunc("/spaces/{id}", middlewares.SetMiddlewareJSON(server.getSpace)).Methods("GET")
	server.Router.HandleFunc("/spaces/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.UpdateSpace))).Methods("PUT")
	server.Router.HandleFunc("/spaces/{id}", middlewares.SetMiddlewareAuthentication(server.DeleteSpace)).Methods("DELETE")
}
