package main

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Resource struct {
	Id      int64
	Pid     int64
	Rname   string
	Rtype   string
	Request bool
	Dsc     string
	Zipcode string
}

type Charity struct {
	Id int64
	Pid int64
	Cname string
	CURL string
	Ccity string
	Cstate string
}

type Person struct {
	Id int64
	Username string
	Email  string
	Password string
	Zipcode string
}

func (s *Server) Health(w http.ResponseWriter, r *http.Request) {
	enableCors(&w);
	err := s.ds.db.Ping()

	var dbStatus string
	if dbStatus = "ok"; err!=nil {
		dbStatus = "not ok"
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"healthy": "ok",
		"db": dbStatus,
	})
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {

	if r.Method!=http.MethodPost{
		http.Error(w, ("This is meant to be posted to."), http.StatusBadRequest)
		return
	}
	// Enable response for all access
	enableCors(&w);

	// Declare a new User struct.
	var p Person;

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p, err = s.ds.GetPerson(p.Username, p.Password)
	if err != nil{
		log.Info().Msg("Wasn't able to login " + err.Error())
	}
	log.Info().Msg(fmt.Sprintf("Person: %+v", p))

	// ! Do need to send back response in order for information to work right
	_ = json.NewEncoder(w).Encode(p)
}

func (s *Server) PersonCreate(w http.ResponseWriter, r *http.Request){
	// Enable response for all access
	enableCors(&w)

	// Declare a new Person struct.
	var p Person

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pid, err := s.ds.SavePerson(p)
	if err != nil{
		log.Info().Msg("Wasn't able to save person " + err.Error())
	}

	p.Id = pid

	log.Info().Msg(fmt.Sprintf("Person: %+v", p))

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"personCreate": pid,
	})
}

func (s *Server) GetResources(w http.ResponseWriter, r *http.Request) {
	// Enable response for all access
	enableCors(&w)
	// Declare a new Resource struct.
	var p Person
	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := s.ds.GetResourceByZip(p.Zipcode)
	if err != nil {
		log.Info().Msg("Wasn't able to read resources " + err.Error())
	}
	// _ = json.NewEncoder(w).Encode(res)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"resources": res,
	})
}

func (s *Server) PostResource(w http.ResponseWriter, r *http.Request) {
	// Enable response for all access
	enableCors(&w)
	// Declare a new Resource struct.
	var resource Resource
	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&resource)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := s.ds.SaveResource(resource)
	if err != nil {
		log.Info().Msg("Wasn't able to save resource " + err.Error())
	}
	resource.Id = id
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"resourceId": id,
	})
}

func (s *Server) GetCharities(w http.ResponseWriter, r *http.Request){

	// Enable response for all access
	enableCors(&w)

	// Declare a new Person struct.
	var p Person

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := s.ds.GetCharitiesByUser(int(p.Id))
	if err != nil {
		log.Info().Msg("Wasn't able to read charities " + err.Error())
	}

	_ = json.NewEncoder(w).Encode(res)


}

func (s *Server) PostCharity(w http.ResponseWriter, r *http.Request) {

	// Enable response for all access
	enableCors(&w)

	// Declare a new Charity struct.
	var charity Charity

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&charity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.Logger.Info().Msg(fmt.Sprintf("We're adding this charity to the db: %+v", charity))
	id, err := s.ds.SaveCharity(charity)
	if err != nil {
		log.Info().Msg("Wasn't able to save charity " + err.Error())
	}
	charity.Id = id


}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
	(*w).Header().Set("Content-Type", "application/json")
}
