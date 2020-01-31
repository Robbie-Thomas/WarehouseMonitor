package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/robbie-thomas/fullstack/api/auth"
	"github.com/robbie-thomas/fullstack/api/models"
	"github.com/robbie-thomas/fullstack/api/responses"
	"github.com/robbie-thomas/fullstack/api/utils/formaterror"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (server *Server) CreateSpace(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	space := models.Space{}
	err = json.Unmarshal(body, &space)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	space.Prepare()
	err = space.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorised"))
		return
	}
	if uid != space.ID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	spaceCreated, err := space.SaveSpace(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, spaceCreated.ID))
	responses.JSON(w, http.StatusCreated, spaceCreated)
}

func (server *Server) GetSpaces(w http.ResponseWriter, r *http.Request) {
	space := models.Space{}

	spaces, err := space.FindAllSpaces(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, spaces)
}

func (server *Server) getSpace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	space := models.Space{}
	spaceReceived, err := space.FindSpaceByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, spaceReceived)
}

func (server *Server) UpdateSpace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorised"))
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorised"))
		return
	}

	space := models.Space{}
	err = server.DB.Debug().Model(models.Space{}).Where("id = ?", pid).Take(&space).Error
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Space not found"))
		return
	}

	if uid != space.ID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorised"))
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	spaceUpdate := models.Space{}
	err = json.Unmarshal(body, &spaceUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	if uid != spaceUpdate.ID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorised"))
		return
	}
	spaceUpdate.Prepare()
	err = spaceUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	spaceUpdate.ID = space.ID
	UpdatedSpace, err := spaceUpdate.UpdateASpace(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, UpdatedSpace)

}

func (server *Server) DeleteSpace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	space := models.Space{}
	err = server.DB.Debug().Model(models.Space{}).Where("id = ?", pid).Take(&space).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorised"))
		return
	}

	if uid != space.OwnerID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorised"))
		return
	}
	_, err = space.DeleteASpace(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")

}
