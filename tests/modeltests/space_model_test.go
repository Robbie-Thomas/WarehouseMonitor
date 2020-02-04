package modeltests

/*import (
	"github.com/robbie-thomas/fullstack/api/models"
	"gopkg.in/go-playground/assert.v1"
	"log"
	"testing"
)


func TestFindAllSpaces(t *testing.T) {
	err := refreshAllTables()
	if err != nil {
		log.Fatalf("Error refreshing user and space table %v\n", err)
	}
	_, _, _, _, _, _, err = seedALL()
	if err != nil {
		log.Fatalf("Error seeding user and space  table %v\n", err)
	}
	spaces, err := spaceInstance.FindAllSpaces(server.DB)
	if err != nil {
		t.Errorf("Error seeding user and space  table %v\n", err)
		return
	}
	assert.Equal(t, len(*spaces), 2)
}

func TestSaveSpace(t *testing.T) {

	err := refreshAllTables()
	if err != nil {
		log.Fatalf("Error refreshing tables %v\n", err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}

	newSpace := models.Space{
		ID:       1,
		SpaceName: "Warehouse 1",
		OwnerID: user.ID,
	}
	savedSpace, err := newSpace.SaveSpace(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the space: %v\n", err)
		return
	}
	assert.Equal(t, newSpace.ID, savedSpace.ID)
	assert.Equal(t, newSpace.SpaceName, savedSpace.SpaceName)
	assert.Equal(t, newSpace.OwnerID, savedSpace.OwnerID)

}

func TestGetSpaceByID(t *testing.T) {

	err := refreshAllTables()
	if err != nil {
		log.Fatalf("Error refreshing user and space table: %v\n", err)
	}
	space, err := seedOneUserAndOneSpace()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	foundSpace, err := spaceInstance.FindSpaceByID(server.DB, uint64(space.ID))
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	assert.Equal(t, foundSpace.ID, space.ID)
	assert.Equal(t, foundSpace.SpaceName, space.SpaceName)
}

func TestUpdateASpace(t *testing.T) {

	err := refreshAllTables()
	if err != nil {
		log.Fatalf("Error refreshing user and space table: %v\n", err)
	}
	user , err := seedOneUserAndOneSpace()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	spaceUpdate := models.Space{
		ID:       1,
		SpaceName: "Warehouse Update",
		OwnerID: user.ID,
	}
	updatedSpace, err := spaceUpdate.UpdateASpace(server.DB)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	assert.Equal(t, updatedSpace.ID, spaceUpdate.ID)
	//assert.Equal(t, updatedSpace.ID, 543) //where 543 is invalid
	assert.Equal(t, updatedSpace.SpaceName, spaceUpdate.SpaceName)
	assert.Equal(t, updatedSpace.OwnerID, spaceUpdate.OwnerID)
}

func TestDeleteASpace(t *testing.T) {

	err := refreshAllTables()
	if err != nil {
		log.Fatalf("Error refreshing user and space table: %v\n", err)
	}
	space, err := seedOneUserAndOneSpace()
	if err != nil {
		log.Fatalf("Error Seeding tables")
	}
	isDeleted, err := spaceInstance.DeleteASpace(server.DB, uint64(space.ID), space.OwnerID)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	assert.Equal(t, isDeleted, int64(1))
}*/
