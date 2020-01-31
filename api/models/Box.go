package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"html"
	"strings"
	"time"
)

type Box struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	BoxName   string    `gorm:"size:255;not null;unique" json:"boxname"`
	Zone      Zone      `json:"zone"`
	ZoneID    uint32    `gorm:"not null" json:"zone_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (b *Box) Prepare() {
	b.ID = 0
	b.BoxName = html.EscapeString(strings.TrimSpace(b.BoxName))
	b.Zone = Zone{}
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()
}

func (b *Box) Validate() error {

	if b.BoxName == "" {
		return errors.New("Required Box Name")
	}
	if b.ZoneID < 1 {
		return errors.New("Requires Zone")
	}
	return nil
}

func (b *Box) SaveBox(db *gorm.DB) (*Box, error) {
	var err error
	err = db.Debug().Model(&Box{}).Create(&b).Error
	if err != nil {
		return &Box{}, err
	}
	box, err2 := b.FetchZone(err, db)
	if err2 != nil {
		return box, err2
	}
	return b, nil
}

func (b *Box) FetchZone(err error, db *gorm.DB) (*Box, error) {
	if b.ID != 0 {
		err = db.Debug().Model(&Zone{}).Where("id = ?", b.ZoneID).Take(&b.Zone).Error
		if err != nil {
			return &Box{}, err
		}
	}
	return nil, nil
}

func (b *Box) FindAllBoxes(db *gorm.DB) (*[]Box, error) {
	var err error
	boxes := []Box{}
	err = db.Debug().Model(&Box{}).Limit(100).Find(&boxes).Error
	if err != nil {
		return &[]Box{}, err
	}
	if len(boxes) > 0 {
		for i := range boxes {
			err := db.Debug().Model(&Zone{}).Where("id = ?", boxes[i].ZoneID).Take(&boxes[i].Zone).Error
			if err != nil {
				return &[]Box{}, err
			}
		}
	}
	return &boxes, nil
}

func (b *Box) FindBoxByID(db *gorm.DB, pid uint64) (*Box, error) {
	var err error
	err = db.Debug().Model(&Box{}).Where("id = ?", pid).Take(&b).Error
	if err != nil {
		return &Box{}, err
	}
	box, err2 := b.FetchZone(err, db)
	if err2 != nil {
		return box, err2
	}
	return b, nil
}

func (b *Box) UpdateABox(db *gorm.DB) (*Box, error) {
	var err error
	err = db.Debug().Model(&Box{}).Where("id = ?", b.ID).Updates(Box{
		BoxName:   b.BoxName,
		UpdatedAt: time.Now(),
	}).Error
	if err != nil {
		return &Box{}, err
	}
	box, err2 := b.FetchZone(err, db)
	if err2 != nil {
		return box, err2
	}
	return b, nil
}

func (b *Box) DeleteABox(db *gorm.DB, pid uint64, uid uint32) (int64, error) {
	db = db.Debug().Model(&Box{}).Where("id = ? and zone_id = ?", pid, uid).Take(&Box{}).Delete(&Box{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Box not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
