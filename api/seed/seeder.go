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

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Post{}, &models.User{}, &models.Space{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Post{}, &models.Space{}).Error
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

	for i := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		posts[i].AuthorID = users[i].ID

		err = db.Debug().Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}

		err = db.Debug().Model(&models.Space{}).Create(&spaces[i]).Error
		if err != nil {
			log.Fatalf("cannot seed spaces table: %v", err)
		}
	}
}
