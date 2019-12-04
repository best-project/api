package storage

import (
	"fmt"
	"github.com/best-project/api/internal"
	"github.com/best-project/api/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Database struct {
	User         *User
	Task         *TaskDB
	Course       *CourseDB
	CourseResult *CourseResultDB
}

func NewDatabase(cfg *config.Config, entry *logrus.Logger) (*Database, error) {
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
		cfg.DbUser, cfg.DbPass, cfg.DbHost, cfg.DbPort, cfg.DbName)

	entry.Info("Starting database connection")
	db, err := gorm.Open("mysql", url)
	if err != nil {
		return nil, errors.Wrap(err, "unable to connect to database")
	}

	entry.Info("Clearing database")

	userDB := &User{db}
	taskDB := &TaskDB{db}
	courseDB := &CourseDB{db}
	courseResultsDB := &CourseResultDB{db}

	// for development
	tables := []interface{}{&internal.Course{}, &internal.User{}, &internal.Task{}, &internal.CourseResult{}}
	db.DropTableIfExists(tables...)
	db.CreateTable(tables...)
	pass, err := bcrypt.GenerateFromPassword([]byte("root1234"), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "while hashing pass")
	}
	userDB.SaveUser(&internal.User{Model: gorm.Model{ID: uint(1)}, Username: "JanNowak23141", Email: "root@o2.pl", Password: string(pass), Level: 1, Points: 0, Avatar: "https://www.pngtube.com/myfile/detail/479-4792237_gopher-dragon-clipart-go-gopher-logos.png"})
	tasks := []internal.Task{
		{
			Word:      "aisle",
			Translate: "przejście, nawa boczna",
			Image:     "https://pixabay.com/get/52e9dd464857b108f5d08460962034761d3cdfe04e50744f722c72d4904cc2_1280.jpg",
		},
		{
			Word:      "art",
			Translate: "sztuka",
			Image:     "https://pixabay.com/get/54e6dc454356ab14f6da8c7dda793f7d1d37dde4564c704c732878dc9f4ec25d_640.jpg",
		},
		{
			Word:      "artist",
			Translate: "artysta",
			Image:     "https://cdn.pixabay.com/photo/2019/09/29/09/06/horse-4512673_960_720.jpg",
		},
		{
			Word:      "artistic",
			Translate: "artystyczny",
			Image:     "https://pixabay.com/get/54e4d34b4255a814f6da8c7dda793f7d1d37dde4564c704c732878dd9748c65d_1280.jpg",
		},
		{
			Word:      "band",
			Translate: "zespół, kapela",
			Image:     "https://pixabay.com/get/50e9d4414856b108f5d08460962034761d3cdfe04e50744f722c72d4924ac2_1280.jpg",
		},
		{
			Word:      "brush",
			Translate: "szczotka, pędzel",
			Image:     "https://pixabay.com/get/55e4d1464955aa14f6da8c7dda793f7d1d37dde4564c704c732878dd974bc45d_1280.jpg",
		},
		{
			Word:      "camera",
			Translate: "kamera, aparat fotograficzny",
			Image:     "https://pixabay.com/get/54e5dc4b4f52ab14f6da8c7dda793f7d1d37dde4564c704c732878dd974bc05e_1280.jpg",
		},
		{
			Word:      "canvas",
			Translate: "płótno",
			Image:     "https://pixabay.com/get/57e8d3454d50a914f6da8c7dda793f7d1d37dde4564c704c732878dd9644c55b_1280.jpg",
		},
	}
	courseDB.SaveCourse(&internal.Course{Name: "Kultura", UserID: 1, DifficultyLevel: "normal", MaxPoints: 80, Language: "en", Rate: 4.0,
		Description: "Słowo kultura ma wiele znaczeń. Interpretuje się ją w różny sposób przez przedstawicieli wielu dziedzin. Kulturę można określić jako ogół wytworów ludzi, zarówno fizycznych, materialnych, jak i duchowych, symbolicznych.",
		Image:       "https://culture360.asef.org/media/2018/5/european_commission_shutterstock_584963080.jpg",
		Task:        tasks})
	courseDB.SaveCourse(&internal.Course{Name: "Nauka", UserID: 1, DifficultyLevel: "normal", MaxPoints: 80, Language: "en", Rate: 4.0,
		Description: "Słowo Nauka ma wiele znaczeń. Interpretuje się ją w różny sposób przez przedstawicieli wielu dziedzin. Kulturę można określić jako ogół wytworów ludzi, zarówno fizycznych, materialnych, jak i duchowych, symbolicznych.",
		Image:       "https://culture360.asef.org/media/2018/5/european_commission_shutterstock_584963080.jpg",
		Task:        tasks})
	courseDB.SaveCourse(&internal.Course{Name: "Religia", UserID: 1, DifficultyLevel: "normal", MaxPoints: 80, Language: "en", Rate: 4.0,
		Description: "Słowo Religia ma wiele znaczeń. Interpretuje się ją w różny sposób przez przedstawicieli wielu dziedzin. Kulturę można określić jako ogół wytworów ludzi, zarówno fizycznych, materialnych, jak i duchowych, symbolicznych.",
		Image:       "https://culture360.asef.org/media/2018/5/european_commission_shutterstock_584963080.jpg",
		Task:        tasks})
	courseDB.SaveCourse(&internal.Course{Name: "Kultura2", UserID: 1, DifficultyLevel: "normal", MaxPoints: 80, Language: "en", Rate: 4.0,
		Description: "Słowo kultura ma wiele znaczeń. Interpretuje się ją w różny sposób przez przedstawicieli wielu dziedzin. Kulturę można określić jako ogół wytworów ludzi, zarówno fizycznych, materialnych, jak i duchowych, symbolicznych.",
		Image:       "https://culture360.asef.org/media/2018/5/european_commission_shutterstock_584963080.jpg",
		Task:        tasks})
	courseDB.SaveCourse(&internal.Course{Name: "Nauka2", UserID: 1, DifficultyLevel: "normal", MaxPoints: 80, Language: "en", Rate: 4.0,
		Description: "Słowo Nauka ma wiele znaczeń. Interpretuje się ją w różny sposób przez przedstawicieli wielu dziedzin. Kulturę można określić jako ogół wytworów ludzi, zarówno fizycznych, materialnych, jak i duchowych, symbolicznych.",
		Image:       "https://culture360.asef.org/media/2018/5/european_commission_shutterstock_584963080.jpg",
		Task:        tasks})
	courseDB.SaveCourse(&internal.Course{Name: "Religia2", UserID: 1, DifficultyLevel: "normal", MaxPoints: 80, Language: "en", Rate: 4.0,
		Description: "Słowo Religia ma wiele znaczeń. Interpretuje się ją w różny sposób przez przedstawicieli wielu dziedzin. Kulturę można określić jako ogół wytworów ludzi, zarówno fizycznych, materialnych, jak i duchowych, symbolicznych.",
		Image:       "https://culture360.asef.org/media/2018/5/european_commission_shutterstock_584963080.jpg",
		Task:        tasks})

	return &Database{userDB, taskDB, courseDB, courseResultsDB}, nil
}
