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

func (server *Server) CreateItem(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	item := models.Item{}
	err = json.Unmarshal(body, &item)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	item.Prepare()
	err = item.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	bix, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorised"))
		return
	}
	if bix != item.ID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	itemCreated, err := item.SaveItem(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, itemCreated.ID))
	responses.JSON(w, http.StatusCreated, itemCreated)
}

func (server *Server) GetItems(w http.ResponseWriter, r *http.Request) {
	item := models.Item{}

	items, err := item.FindAllItems(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, items)
}

func (server *Server) fetchItemsForUser(w http.ResponseWriter, r *http.Request) {
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
	item := models.Item{}
	items, err := item.FindAllItems(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	for i := range *items {
		if (*items)[i].Box.Zone.Space.User.ID == uint32(uid) {
			(*items)[i].Box.Zone.Space.User = *UserReceived
		}
	}
	responses.JSON(w, http.StatusOK, items)
}

func (server *Server) getItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	iid, err := strconv.ParseUint(vars["itemID"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	item := models.Item{}
	itemReceived, err := item.FindItemByID(server.DB, iid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	box := models.Box{}
	box = itemReceived.Box
	zone := models.Zone{}
	zone = box.Zone
	space := models.Space{}
	space = zone.Space
	_ = space.User

	responses.JSON(w, http.StatusOK, itemReceived)
}

func (server *Server) fetchItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	iid, err := strconv.ParseUint(vars["itemID"], 10, 64)
	bid, err := strconv.ParseUint(vars["boxID"], 10, 64)
	zid, err := strconv.ParseUint(vars["zoneID"], 10, 64)
	sid, err := strconv.ParseUint(vars["spaceID"], 10, 64)
	uid, err := strconv.ParseUint(vars["userID"], 10, 64)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	item := models.Item{}
	itemReceived, err := item.FindItemByIDAndBoxID(server.DB, iid, bid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	box := models.Box{}
	BoxReceived, err := box.FindBoxByIDAndZoneID(server.DB, bid, zid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	zone := models.Zone{}
	ZoneReceived, err := zone.FindZoneBySpaceID(server.DB, zid, sid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	space := models.Space{}
	spaceReceived, err := space.FindSpaceByIDAndUserID(server.DB, sid, uid)
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

	user = *UserReceived
	ZoneReceived.Space = *spaceReceived
	BoxReceived.Zone = *ZoneReceived
	itemReceived.Box = *BoxReceived
	responses.JSON(w, http.StatusOK, itemReceived)
}

func (server *Server) fetchItemForUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	iid, err := strconv.ParseUint(vars["itemID"], 10, 64)
	uid, err := strconv.ParseUint(vars["userID"], 10, 64)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	item := models.Item{}
	itemReceived, err := item.FindItemByID(server.DB, iid)
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
	user = *UserReceived
	box := models.Box{}
	BoxReceived, err := box.FindBoxByIDAndZoneID(server.DB, uint64(itemReceived.Box.ID), uint64(itemReceived.Box.ZoneID))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	zone := models.Zone{}
	ZoneReceived, err := zone.FindZoneBySpaceID(server.DB, uint64(BoxReceived.Zone.ID), uint64(BoxReceived.Zone.SpaceID))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	space := models.Space{}
	spaceReceived, err := space.FindSpaceByIDAndUserID(server.DB, uint64(ZoneReceived.Space.ID), uint64(ZoneReceived.Space.OwnerID))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	ZoneReceived.Space = *spaceReceived
	BoxReceived.Zone = *ZoneReceived
	itemReceived.Box = *BoxReceived
	responses.JSON(w, http.StatusOK, itemReceived)
}

func (server *Server) UpdateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	iid, err := strconv.ParseUint(vars["itemID"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorised"))
		return
	}

	bix, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorised"))
		return
	}

	item := models.Item{}
	err = server.DB.Debug().Model(models.Item{}).Where("id = ?", iid).Take(&item).Error
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Item not found"))
		return
	}

	if bix != item.ID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorised"))
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	itemUpdate := models.Item{}
	err = json.Unmarshal(body, &itemUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	if bix != itemUpdate.ID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorised"))
		return
	}
	itemUpdate.Prepare()
	err = itemUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	itemUpdate.ID = item.ID
	UpdatedItem, err := itemUpdate.UpdateAItem(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, UpdatedItem)

}

func (server *Server) DeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	iid, err := strconv.ParseUint(vars["itemID"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	bix, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	item := models.Item{}
	err = server.DB.Debug().Model(models.Item{}).Where("id = ?", iid).Take(&item).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorised"))
		return
	}

	if bix != item.BoxID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorised"))
		return
	}
	_, err = item.DeleteAItem(server.DB, iid, bix)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", bix))
	responses.JSON(w, http.StatusNoContent, "")

}
