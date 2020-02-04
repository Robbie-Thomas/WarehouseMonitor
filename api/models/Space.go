package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"html"
	"strings"
	"time"
)

type Space struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	SpaceName string    `gorm:"size:255;not null;unique" json:"spacename"`
	User      User      `json:"user"`
	OwnerID   uint32    `gorm:"not null" json:"owner_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (s *Space) Prepare() {
	s.ID = 0
	s.SpaceName = html.EscapeString(strings.TrimSpace(s.SpaceName))
	s.User = User{}
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
}

func (s *Space) Validate() error {

	if s.SpaceName == "" {
		return errors.New("Required Space Name")
	}
	if s.ID < 1 {
		return errors.New("Requires User")
	}
	return nil
}

func (s *Space) SaveSpace(db *gorm.DB) (*Space, error) {
	var err error
	err = db.Debug().Model(&Space{}).Create(&s).Error
	if err != nil {
		return &Space{}, err
	}
	space, err2 := s.FetchUser(err, db)
	if err2 != nil {
		return space, err2
	}
	return s, nil
}

func (s *Space) FetchUser(err error, db *gorm.DB) (*Space, error) {
	if s.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", s.OwnerID).Take(&s.User).Error
		if err != nil {
			return &Space{}, err
		}
	}
	return nil, nil
}

func (s *Space) FindAllSpaces(db *gorm.DB) (*[]Space, error) {
	var err error
	var spaces []Space
	err = db.Debug().Model(&Space{}).Limit(100).Find(&spaces).Error
	if err != nil {
		return &[]Space{}, err
	}
	if len(spaces) > 0 {
		for i := range spaces {
			err := db.Debug().Model(&User{}).Where("id = ?", spaces[i].OwnerID).Take(&spaces[i].User).Error
			if err != nil {
				return &[]Space{}, err
			}
		}
	}
	return &spaces, nil
}

func (s *Space) FindSpaceByID(db *gorm.DB, uid uint64) (*Space, error) {
	var err error
	err = db.Debug().Model(&Space{}).Where("id = ?", uid).Take(&s).Error
	if err != nil {
		return &Space{}, err
	}
	space, err2 := s.FetchUser(err, db)
	if err2 != nil {
		return space, err2
	}
	return s, nil
}

func (s *Space) FindSpaceByIDAndUserID(db *gorm.DB, sid uint64, uid uint64) (*Space, error) {
	var err error
	err = db.Debug().Model(&Space{}).Where("id = ? and owner_id = ?", sid, uid).Take(&s).Error
	if err != nil {
		return &Space{}, err
	}
	user, err2 := s.FetchUser(err, db)
	if err2 != nil {
		return user, err2
	}
	return s, nil
}

func (s *Space) UpdateASpace(db *gorm.DB) (*Space, error) {
	var err error
	err = db.Debug().Model(&Space{}).Where("id = ?", s.ID).Updates(Space{
		SpaceName: s.SpaceName,
		UpdatedAt: time.Now(),
	}).Error
	if err != nil {
		return &Space{}, err
	}
	space, err2 := s.FetchUser(err, db)
	if err2 != nil {
		return space, err2
	}
	return s, nil
}

func (s *Space) DeleteASpace(db *gorm.DB, pid uint64, uid uint32) (int64, error) {
	db = db.Debug().Model(&Space{}).Where("id = ? and user_id = ?", pid, uid).Take(&Space{}).Delete(&Space{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Space not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (s *Space) FetchUserForSpace(db *gorm.DB) (User, error) {
	var err error
	s.FetchUser(err, db)
	if err != nil {
		return s.User, err
	}
	return s.User, nil
}
