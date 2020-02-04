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
	server.Router.HandleFunc("/users/{userID}", middlewares.SetMiddlewareJSON(server.GetUser)).Methods("GET")
	server.Router.HandleFunc("/users/{userID}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.UpdateUser))).Methods("PUT")
	server.Router.HandleFunc("/users/{userID}", middlewares.SetMiddlewareAuthentication(server.DeleteUser)).Methods("DELETE")

	//Posts routes
	server.Router.HandleFunc("/posts", middlewares.SetMiddlewareJSON(server.CreatePost)).Methods("POST")
	server.Router.HandleFunc("/posts", middlewares.SetMiddlewareJSON(server.GetPosts)).Methods("GET")
	server.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareJSON(server.getPost)).Methods("GET")
	server.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.UpdatePost))).Methods("PUT")
	server.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareAuthentication(server.DeletePost)).Methods("DELETE")

	server.Router.HandleFunc("/spaces", middlewares.SetMiddlewareJSON(server.CreateSpace)).Methods("POST")
	server.Router.HandleFunc("/spaces", middlewares.SetMiddlewareJSON(server.GetSpaces)).Methods("GET")
	server.Router.HandleFunc("/spaces/{spaceID}", middlewares.SetMiddlewareJSON(server.getSpace)).Methods("GET")
	server.Router.HandleFunc("/spaces/{spaceID}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.UpdateSpace))).Methods("PUT")
	server.Router.HandleFunc("/spaces/{spaceID}", middlewares.SetMiddlewareAuthentication(server.DeleteSpace)).Methods("DELETE")

	server.Router.HandleFunc("/zones", middlewares.SetMiddlewareJSON(server.CreateZone)).Methods("POST")
	server.Router.HandleFunc("/zones", middlewares.SetMiddlewareJSON(server.getZone)).Methods("GET")
	server.Router.HandleFunc("/zones/{zoneID}", middlewares.SetMiddlewareJSON(server.fetchZoneForUser)).Methods("GET")
	server.Router.HandleFunc("/zones/{zoneID}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.UpdateZone))).Methods("PUT")
	server.Router.HandleFunc("/zones/{zoneID}", middlewares.SetMiddlewareAuthentication(server.DeleteZone)).Methods("DELETE")

	server.Router.HandleFunc("/boxes", middlewares.SetMiddlewareJSON(server.CreateBox)).Methods("POST")
	server.Router.HandleFunc("/boxes", middlewares.SetMiddlewareJSON(server.getBoxes)).Methods("GET")
	server.Router.HandleFunc("/boxes/{boxID}", middlewares.SetMiddlewareJSON(server.getBox)).Methods("GET")
	server.Router.HandleFunc("/boxes/{boxID}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.UpdateBox))).Methods("PUT")
	server.Router.HandleFunc("/boxes/{boxID}", middlewares.SetMiddlewareAuthentication(server.DeleteBox)).Methods("DELETE")

	server.Router.HandleFunc("/items", middlewares.SetMiddlewareJSON(server.CreateItem)).Methods("POST")
	server.Router.HandleFunc("/items", middlewares.SetMiddlewareJSON(server.GetItems)).Methods("GET")
	server.Router.HandleFunc("/items/{itemID}", middlewares.SetMiddlewareJSON(server.getItem)).Methods("GET")
	server.Router.HandleFunc("/items/{itemID}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.UpdateItem))).Methods("PUT")
	server.Router.HandleFunc("/items/{itemID}", middlewares.SetMiddlewareAuthentication(server.DeleteItem)).Methods("DELETE")

	//server.Router.HandleFunc("/user/{userID}/space/{spaceID}", middlewares.SetMiddlewareJSON(server.getSpaceForID)).Methods("GET")
	server.Router.HandleFunc("/user/{userID}/space/{spaceID}/zone/{zoneID}", middlewares.SetMiddlewareJSON(server.getSpaceForID)).Methods("GET")
	server.Router.HandleFunc("/user/{userID}/space/{spaceID}/zone/{zoneID}/box/{boxID}", middlewares.SetMiddlewareJSON(server.getSpaceForID)).Methods("GET")
	server.Router.HandleFunc("/user/{userID}/space/{spaceID}/zone/{zoneID}/box/{boxID}/item/{itemID}", middlewares.SetMiddlewareJSON(server.fetchItem)).Methods("GET")

	server.Router.HandleFunc("/user/{userID}/spaces", middlewares.SetMiddlewareJSON(server.GetSpaces)).Methods("GET")
	server.Router.HandleFunc("/user/{userID}/spaces", middlewares.SetMiddlewareJSON(server.CreateSpace)).Methods("POST")
	server.Router.HandleFunc("/user/{userID}/space/{spaceID}", middlewares.SetMiddlewareJSON(server.getSpaceForID)).Methods("GET")
	server.Router.HandleFunc("/user/{userID}/space/{spaceID}", middlewares.SetMiddlewareJSON(server.UpdateSpace)).Methods("PUT")
	server.Router.HandleFunc("/user/{userID}/space/{spaceID}", middlewares.SetMiddlewareJSON(server.DeleteSpace)).Methods("DELETE")

	server.Router.HandleFunc("/user/{userID}/zones", middlewares.SetMiddlewareJSON(server.fetchZonesForUser)).Methods("GET")
	server.Router.HandleFunc("/user/{userID}/zones", middlewares.SetMiddlewareJSON(server.CreateZone)).Methods("POST")
	server.Router.HandleFunc("/user/{userID}/zone/{zoneID}", middlewares.SetMiddlewareJSON(server.fetchZoneForUser)).Methods("GET")
	server.Router.HandleFunc("/user/{userID}/zone/{zoneID}", middlewares.SetMiddlewareJSON(server.UpdateZone)).Methods("PUT")
	server.Router.HandleFunc("/user/{userID}/zone/{zoneID}", middlewares.SetMiddlewareJSON(server.DeleteZone)).Methods("DELETE")

	server.Router.HandleFunc("/user/{userID}/boxes", middlewares.SetMiddlewareJSON(server.fetchBoxes)).Methods("GET")
	server.Router.HandleFunc("/user/{userID}/boxes", middlewares.SetMiddlewareJSON(server.CreateBox)).Methods("POST")
	server.Router.HandleFunc("/user/{userID}/box/{boxID}", middlewares.SetMiddlewareJSON(server.fetchBox)).Methods("GET")
	server.Router.HandleFunc("/user/{userID}/box/{boxID}", middlewares.SetMiddlewareJSON(server.UpdateBox)).Methods("PUT")
	server.Router.HandleFunc("/user/{userID}/box/{boxID}", middlewares.SetMiddlewareJSON(server.DeleteBox)).Methods("DELETE")

	server.Router.HandleFunc("/user/{userID}/items", middlewares.SetMiddlewareJSON(server.GetItems)).Methods("GET")
	server.Router.HandleFunc("/user/{userID}/items", middlewares.SetMiddlewareJSON(server.CreateItem)).Methods("POST")
	server.Router.HandleFunc("/user/{userID}/item/{itemID}", middlewares.SetMiddlewareJSON(server.getItem)).Methods("GET")
	server.Router.HandleFunc("/user/{userID}/item/{itemID}", middlewares.SetMiddlewareJSON(server.UpdateItem)).Methods("PUT")
	server.Router.HandleFunc("/user/{userID}/item/{itemID}", middlewares.SetMiddlewareJSON(server.DeleteItem)).Methods("DELETE")

}
