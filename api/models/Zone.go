package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"html"
	"strings"
	"time"
)

type Zone struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	ZoneName  string    `gorm:"size:255;not null;unique" json:"zonename"`
	Space     Space     `json:"space"`
	SpaceID   uint32    `gorm:"not null" json:"space_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (z *Zone) Prepare() {
	z.ID = 0
	z.ZoneName = html.EscapeString(strings.TrimSpace(z.ZoneName))
	z.Space = Space{}
	z.CreatedAt = time.Now()
	z.UpdatedAt = time.Now()
}

func (z *Zone) Validate() error {

	if z.ZoneName == "" {
		return errors.New("Required Zone Name")
	}
	if z.SpaceID < 1 {
		return errors.New("requires Space")
	}
	return nil
}

func (z *Zone) SaveZone(db *gorm.DB) (*Zone, error) {
	var err error
	err = db.Debug().Model(&Zone{}).Create(&z).Error
	if err != nil {
		return &Zone{}, err
	}
	zone, err2 := z.FetchSpace(err, db)
	if err2 != nil {
		return zone, err2
	}
	return z, nil
}

func (z *Zone) FetchSpace(err error, db *gorm.DB) (*Zone, error) {
	if z.ID != 0 {
		err = db.Debug().Model(&Space{}).Where("id = ?", z.SpaceID).Take(&z.Space).Error
		if err != nil {
			return &Zone{}, err
		}
	}
	return nil, nil
}

func (z *Zone) FindAllZones(db *gorm.DB) (*[]Zone, error) {
	var err error
	var zones []Zone
	err = db.Debug().Model(&Zone{}).Limit(100).Find(&zones).Error
	if err != nil {
		return &[]Zone{}, err
	}
	if len(zones) > 0 {
		for i := range zones {
			err := db.Debug().Model(&Space{}).Where("id = ?", zones[i].SpaceID).Take(&zones[i].Space).Error
			if err != nil {
				return &[]Zone{}, err
			}
		}
	}
	return &zones, nil
}

func (z *Zone) FindZoneBySpaceID(db *gorm.DB, zid uint64, sid uint64) (*Zone, error) {
	var err error
	err = db.Debug().Model(&Space{}).Where("id = ? and space_id = ?", zid, sid).Take(&z).Error
	if err != nil {
		return &Zone{}, err
	}
	zone, err2 := z.FetchSpace(err, db)
	if err2 != nil {
		return zone, err2
	}
	return z, nil
}

func (z *Zone) FindZoneByID(db *gorm.DB, pid uint64) (*Zone, error) {
	var err error
	err = db.Debug().Model(&Zone{}).Where("id = ?", pid).Take(&z).Error
	if err != nil {
		return &Zone{}, err
	}
	zone, err2 := z.FetchSpace(err, db)
	if err2 != nil {
		return zone, err2
	}
	return z, nil
}

func (z *Zone) UpdateAZone(db *gorm.DB) (*Zone, error) {
	var err error
	err = db.Debug().Model(&Zone{}).Where("id = ?", z.ID).Updates(Zone{
		ZoneName:  z.ZoneName,
		UpdatedAt: time.Now(),
	}).Error
	if err != nil {
		return &Zone{}, err
	}
	zone, err2 := z.FetchSpace(err, db)
	if err2 != nil {
		return zone, err2
	}
	return z, nil
}

func (z *Zone) DeleteAZone(db *gorm.DB, pid uint64, uid uint32) (int64, error) {
	db = db.Debug().Model(&Zone{}).Where("id = ? and space_id = ?", pid, uid).Take(&Zone{}).Delete(&Zone{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Zone not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (z *Zone) FindZoneByIDAndSpaceID(db *gorm.DB, zid uint64, sid uint64) (*Zone, error) {
	var err error
	err = db.Debug().Model(&Zone{}).Where("id = ? and space_id = ?", zid, sid).Take(&z).Error
	if err != nil {
		return &Zone{}, err
	}
	space, err2 := z.FetchSpace(err, db)
	if err2 != nil {
		return space, err2
	}
	return z, nil
}
