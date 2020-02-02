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

func (server *Server) CreateZone(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	zone := models.Zone{}
	err = json.Unmarshal(body, &zone)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	zone.Prepare()
	err = zone.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	sid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorised"))
		return
	}
	if sid != zone.ID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	zoneCreated, err := zone.SaveZone(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, zoneCreated.ID))
	responses.JSON(w, http.StatusCreated, zoneCreated)
}

func (server *Server) GetZones(w http.ResponseWriter, r *http.Request) {
	zone := models.Zone{}

	zones, err := zone.FindAllZones(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, zones)
}

func (server *Server) fetchZonesForUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["userID"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	zone := models.Zone{}
	zones, err := zone.FindAllZones(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	user := models.User{}
	UserReceived, err := user.FindUserByID(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	for i := range *zones {
		if (*zones)[i].Space.OwnerID == uint32(uid) {
			(*zones)[i].Space.User = *UserReceived
		}
	}
	responses.JSON(w, http.StatusOK, zones)
}

func (server *Server) getZone(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	zid, err := strconv.ParseUint(vars["zoneID"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	zone := models.Zone{}
	zoneReceived, err := zone.FindZoneByID(server.DB, zid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	space := models.Space{}
	space = zone.Space
	user := models.User{}
	user = space.User
	space.FetchUserForSpace(server.DB)
	zone.Space.User = user
	responses.JSON(w, http.StatusOK, zoneReceived)
}

func (server *Server) fetchZoneForUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	zid, err := strconv.ParseUint(vars["zoneID"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
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
	zone := models.Zone{}
	zoneReceived, err := zone.FindZoneByID(server.DB, zid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	zoneReceived.Space.User = *UserReceived
	responses.JSON(w, http.StatusOK, zoneReceived)
}

func (server *Server) UpdateZone(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	zid, err := strconv.ParseUint(vars["zoneID"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorised"))
		return
	}

	sid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorised"))
		return
	}

	zone := models.Zone{}
	err = server.DB.Debug().Model(models.Zone{}).Where("id = ?", zid).Take(&zone).Error
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Zone not found"))
		return
	}

	if sid != zone.ID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorised"))
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	zoneUpdate := models.Zone{}
	err = json.Unmarshal(body, &zoneUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	if sid != zoneUpdate.ID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorised"))
		return
	}
	zoneUpdate.Prepare()
	err = zoneUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	zoneUpdate.ID = zone.ID
	UpdatedZone, err := zoneUpdate.UpdateAZone(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, UpdatedZone)

}

func (server *Server) DeleteZone(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	zid, err := strconv.ParseUint(vars["zoneID"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	sid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	zone := models.Zone{}
	err = server.DB.Debug().Model(models.Zone{}).Where("id = ?", zid).Take(&zone).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorised"))
		return
	}

	if sid != zone.SpaceID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorised"))
		return
	}
	_, err = zone.DeleteAZone(server.DB, zid, sid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", zid))
	responses.JSON(w, http.StatusNoContent, "")

}
