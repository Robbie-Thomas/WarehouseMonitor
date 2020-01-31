package seed

import (
	"github.com/jinzhu/gorm"
	"github.com/robbie-thomas/fullstack/api/models"
	"log"
)

var users = []models.User{
	{
		Nickname: "Jim bob",
		Email:    "Jibbob@gmail.com",
		Password: "password",
	},
	{
		Nickname: "Sambob",
		Email:    "samsambobbob@gmail.com",
		Password: "password",
	},
}

var posts = []models.Post{
	{
		Title:   "Title 1",
		Content: "Hello world 1",
	},
	{
		Title:   "Title 2",
		Content: "Hello world 2",
	},
}

var spaces = []models.Space{
	{
		SpaceName: "Garage",
		OwnerID:   1,
	},

	{
		SpaceName: "Attic",
		OwnerID:   1,
	},
}
var zones = []models.Zone{
	{
		ZoneName: "Zone 1",
		SpaceID:  1,
	},

	{
		ZoneName: "Zone 2",
		SpaceID:  1,
	},
}

var boxes = []models.Box{
	{
		BoxName: "Box 1",
		ZoneID:  1,
	},

	{
		BoxName: "Box 2",
		ZoneID:  1,
	},
}

var items = []models.Item{
	{
		ItemName: "Table legs",
		BoxID:    1,
	},
	{
		ItemName: "Table top",
		BoxID:    1,
	},
	{
		ItemName: "Years supply of chocolate",
		BoxID:    1,
	},
	{
		ItemName: "Dancing Monkey",
		BoxID:    1,
	},
	{
		ItemName: "Spandalbaner recorsd",
		BoxID:    1,
	},
	{
		ItemName: "Mini eggs",
		BoxID:    1,
	},
	{
		ItemName: "20 year old fine ages gouda",
		BoxID:    1,
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Post{}, &models.User{}, &models.Space{}, &models.Zone{}, &models.Box{}, &models.Item{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.Post{}, &models.User{}, &models.Space{}, &models.Zone{}, &models.Box{}, &models.Item{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Post{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	err = db.Debug().Model(&models.Space{}).AddForeignKey("owner_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}
	err = db.Debug().Model(&models.Zone{}).AddForeignKey("space_id", "spaces(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}
	err = db.Debug().Model(&models.Box{}).AddForeignKey("zone_id", "zones(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}
	err = db.Debug().Model(&models.Item{}).AddForeignKey("box_id", "boxes(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		posts[i].AuthorID = users[i].ID
	}
	for i := range posts {
		err = db.Debug().Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
	for i := range spaces {
		err = db.Debug().Model(&models.Space{}).Create(&spaces[i]).Error
		if err != nil {
			log.Fatalf("cannot seed spaces table: %v", err)
		}
	}
	for i := range zones {
		err = db.Debug().Model(&models.Zone{}).Create(&zones[i]).Error
		if err != nil {
			log.Fatalf("cannot seed zones table: %v", err)
		}
	}
	for i := range boxes {
		err = db.Debug().Model(&models.Box{}).Create(&boxes[i]).Error
		if err != nil {
			log.Fatalf("cannot seed boxes table: %v", err)
		}
	}
	for i := range items {
		err = db.Debug().Model(&models.Item{}).Create(&items[i]).Error
		if err != nil {
			log.Fatalf("cannot seed items table: %v", err)
		}
	}
}
