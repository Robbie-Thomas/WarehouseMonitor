package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"html"
	"strings"
	"time"
)

type Item struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	ItemName  string    `gorm:"size:255;not null;unique" json:"itemname"`
	Box       Box       `json:"box"`
	BoxID     uint32    `gorm:"not null" json:"box_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (i *Item) Prepare() {
	i.ID = 0
	i.ItemName = html.EscapeString(strings.TrimSpace(i.ItemName))
	i.Box = Box{}
	i.CreatedAt = time.Now()
	i.UpdatedAt = time.Now()
}

func (i *Item) Validate() error {

	if i.ItemName == "" {
		return errors.New("Required Item Name")
	}
	if i.BoxID < 1 {
		return errors.New("Requires Box")
	}
	return nil
}

func (i *Item) SaveItem(db *gorm.DB) (*Item, error) {
	var err error
	err = db.Debug().Model(&Item{}).Create(&i).Error
	if err != nil {
		return &Item{}, err
	}
	item, err2 := i.FetchBox(err, db)
	if err2 != nil {
		return item, err2
	}
	return i, nil
}

func (i *Item) FetchBox(err error, db *gorm.DB) (*Item, error) {
	if i.ID != 0 {
		err = db.Debug().Model(&Box{}).Where("id = ?", i.BoxID).Take(&i.Box).Error
		if err != nil {
			return &Item{}, err
		}
	}
	return nil, nil
}

func (i *Item) FindAllItems(db *gorm.DB) (*[]Item, error) {
	var err error
	items := []Item{}
	err = db.Debug().Model(&Item{}).Limit(100).Find(&items).Error
	if err != nil {
		return &[]Item{}, err
	}
	if len(items) > 0 {
		for i := range items {
			err := db.Debug().Model(&Box{}).Where("id = ?", items[i].BoxID).Take(&items[i].Box).Error
			if err != nil {
				return &[]Item{}, err
			}
		}
	}
	return &items, nil
}

func (i *Item) FindItemByID(db *gorm.DB, pid uint64) (*Item, error) {
	var err error
	err = db.Debug().Model(&Item{}).Where("id = ?", pid).Take(&i).Error
	if err != nil {
		return &Item{}, err
	}
	item, err2 := i.FetchBox(err, db)
	if err2 != nil {
		return item, err2
	}
	return i, nil
}

func (i *Item) UpdateAItem(db *gorm.DB) (*Item, error) {
	var err error
	err = db.Debug().Model(&Item{}).Where("id = ?", i.ID).Updates(Item{
		ItemName:  i.ItemName,
		UpdatedAt: time.Now(),
	}).Error
	if err != nil {
		return &Item{}, err
	}
	item, err2 := i.FetchBox(err, db)
	if err2 != nil {
		return item, err2
	}
	return i, nil
}

func (i *Item) DeleteAItem(db *gorm.DB, pid uint64, uid uint32) (int64, error) {
	db = db.Debug().Model(&Item{}).Where("id = ? and box_id = ?", pid, uid).Take(&Item{}).Delete(&Item{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Item not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
