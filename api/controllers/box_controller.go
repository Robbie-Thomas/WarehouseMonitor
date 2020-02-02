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

func (server *Server) CreateBox(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	box := models.Box{}
	err = json.Unmarshal(body, &box)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	box.Prepare()
	err = box.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	zid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorised"))
		return
	}
	if zid != box.ID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	boxCreated, err := box.SaveBox(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, boxCreated.ID))
	responses.JSON(w, http.StatusCreated, boxCreated)
}

func (server *Server) getBoxes(w http.ResponseWriter, r *http.Request) {
	box := models.Box{}

	boxs, err := box.FindAllBoxes(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, boxs)
}

func (server *Server) fetchBoxes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["userID"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	user := models.User{}
	UserReceived, err := user.FindUserByID(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	box := models.Box{}
	boxes, err := box.FindAllBoxes(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	for i := range *boxes {
		if (*boxes)[i].Zone.Space.OwnerID == uint32(uid) {
			(*boxes)[i].Zone.Space.User = *UserReceived
		}
	}
	responses.JSON(w, http.StatusOK, boxes)
}

func (server *Server) getBox(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bid, err := strconv.ParseUint(vars["boxID"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	box := models.Box{}
	boxReceived, err := box.FindBoxByID(server.DB, bid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, boxReceived)
}

func (server *Server) fetchBox(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bid, err := strconv.ParseUint(vars["boxID"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	box := models.Box{}
	boxReceived, err := box.FindBoxByID(server.DB, bid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, boxReceived)
}

func (server *Server) UpdateBox(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	bid, err := strconv.ParseUint(vars["boxID"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorised"))
		return
	}

	zid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorised"))
		return
	}

	box := models.Box{}
	err = server.DB.Debug().Model(models.Box{}).Where("id = ?", bid).Take(&box).Error
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Box not found"))
		return
	}

	if zid != box.ID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorised"))
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	boxUpdate := models.Box{}
	err = json.Unmarshal(body, &boxUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	if zid != boxUpdate.ID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorised"))
		return
	}
	boxUpdate.Prepare()
	err = boxUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	boxUpdate.ID = box.ID
	UpdatedBox, err := boxUpdate.UpdateABox(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, UpdatedBox)

}

func (server *Server) DeleteBox(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bid, err := strconv.ParseUint(vars["boxID"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	zid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	box := models.Box{}
	err = server.DB.Debug().Model(models.Box{}).Where("id = ?", bid).Take(&box).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorised"))
		return
	}

	if zid != box.ZoneID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorised"))
		return
	}
	_, err = box.DeleteABox(server.DB, bid, zid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", zid))
	responses.JSON(w, http.StatusNoContent, "")

}
