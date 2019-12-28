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
	User         User
	Task         Task
	Course       Course
	CourseResult CourseResult
}

func NewDatabase(cfg *config.Config, entry *logrus.Logger) (*Database, error) {
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
		cfg.DbUser, cfg.DbPass, cfg.DbHost, cfg.DbPort, cfg.DbName)

	entry.Info("Starting database connection")
	db, err := gorm.Open("mysql", url)
	if err != nil {
		return nil, errors.Wrap(err, "unable to connect to database")
	}

	userDB := &UserDB{db}
	taskDB := &TaskDB{db}
	courseDB := &CourseDB{db}
	courseResultsDB := &CourseResultDB{db}

	// for development
	if cfg.InitDB {
		tables := []interface{}{&internal.Course{}, &internal.User{}, &internal.Task{}, &internal.CourseResult{}}
		entry.Info("Clearing database")
		db.DropTableIfExists(tables...)
		db.CreateTable(tables...)
		pass, err := bcrypt.GenerateFromPassword([]byte("root1234"), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.Wrap(err, "while hashing pass")
		}
		userDB.SaveUser(&internal.User{Model: gorm.Model{ID: uint(1)}, FirstName: "Jan", LastName: "Nowak", Email: "root@o2.pl", Password: string(pass), Level: 1, Points: 0, Avatar: "https://www.pngtube.com/myfile/detail/479-4792237_gopher-dragon-clipart-go-gopher-logos.png"})

		courseDB.SaveCourse(&internal.Course{Name: "Kultura", UserID: 1, DifficultyLevel: "normal", Language: "en", Rate: 4.0,
			Description: "Słowo kultura ma wiele znaczeń. Interpretuje się ją w różny sposób przez przedstawicieli wielu dziedzin. Kulturę można określić jako ogół wytworów ludzi, zarówno fizycznych, materialnych, jak i duchowych, symbolicznych.",
			Image:       "https://culture360.asef.org/media/2018/5/european_commission_shutterstock_584963080.jpg",
			Task: []internal.Task{
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
				{
					Word:      "character",
					Translate: "charakter, natura",
					Image:     "https://pixabay.com/get/5fe8d1434953b108f5d08460962034761d3cdfe04e50744f722c72d59e4ec5_1280.jpg",
				},
				{
					Word:      "chisel",
					Translate: "dłuto",
					Image:     "https://pixabay.com/get/57e2d54a4d54a814f6da8c7dda793f7d1d37dde4564c704c732878dd9644c05c_1280.jpg",
				},
				{
					Word:      "choir",
					Translate: "chór",
					Image:     "https://pixabay.com/get/5ee4d54a4255b108f5d08460962034761d3cdfe04e50744f722c72d59e4ac4_1280.jpg",
				},
				{
					Word:      "cinema",
					Translate: "kino",
					Image:     "https://pixabay.com/get/5fe1dd454f57b108f5d08460962034761d3cdfe04e50744f722c72d69749c3_1280.jpg",
				},
				{
					Word:      "composer",
					Translate: "kompozytor",
					Image:     "https://pixabay.com/get/57e6d0464f57a414f6da8c7dda793f7d1d37dde4564c704c732878dd954dc25c_1280.jpg",
				},
				{
					Word:      "culture",
					Translate: "kultura",
					Image:     "https://pixabay.com/get/57e8d5444f51af14f6da8c7dda793f7d1d37dde4564c704c732878dd954cc459_1280.jpg",
				},
			},
		}, 10)
		courseDB.SaveCourse(&internal.Course{Name: "Zwierzęta", UserID: 1, DifficultyLevel: "normal", Language: "en", Rate: 4.0,
			Description: "W tym kursie nauczysz się różnych zwierząt po angielsku.",
			Image:       "https://www.imperiumtapet.com/public/uploads/preview/zwierzeta-47jpg-3315339418719ooz4hw2jm.jpg",
			Task: []internal.Task{
				{
					Image:     "https://upload.wikimedia.org/wikipedia/commons/thumb/7/71/2010-kodiak-bear-1.jpg/1200px-2010-kodiak-bear-1.jpg",
					Word:      "bear",
					Translate: "niedźwiedź",
				},
				{
					Image:     "https://i.wpimg.pl/985x0/m.fotoblogia.pl/kardynalek-jeremy-black-6cabc0d6.jpg",
					Word:      "bird",
					Translate: "ptak",
				},
				{
					Image:     "https://static1.s-trojmiasto.pl/zdj/c/n/9/2297/620x0/2297932-Sciety-czubek-lewego-lub-prawego-ucha-oznacza-ze-kot-przebyl-zabieg__c_0_262_1664_1271.jpg",
					Word:      "cat",
					Translate: "kot",
				},
				{
					Image:     "https://www.dw.com/image/19343569_303.jpg",
					Word:      "chicken",
					Translate: "kurczak",
				},
				{
					Image:     "https://cdn.britannica.com/55/174255-050-526314B6/brown-Guernsey-cow.jpg",
					Word:      "cow",
					Translate: "krowa",
				},
				{
					Image:     "https://static.scientificamerican.com/blogs/cache/file/BB6F1FE0-4FDE-4E6E-A986664CE30602E4_source.jpg?w=590&h=800&2F8476C1-DF14-49BA-84FFE94218CC4933",
					Word:      "dog",
					Translate: "pies",
				},
				{
					Image:     "https://www.dolphinproject.com/wp-content/uploads/2019/07/Maya-870x580.jpg",
					Word:      "dolphin",
					Translate: "delfin",
				},
				{
					Image:     "https://images.unsplash.com/photo-1459682687441-7761439a709d?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&w=1000&q=80",
					Word:      "duck",
					Translate: "kaczka",
				},
				{
					Image:     "https://cdn.vox-cdn.com/thumbor/yMAqK_pz6uypAmAFSZT0wJK9EYM=/0x0:1280x720/1200x800/filters:focal(538x258:742x462)/cdn.vox-cdn.com/uploads/chorus_image/image/65609281/The_Elephant_Queen_Unit_Photo_13.0.jpg",
					Word:      "elephant",
					Translate: "słoń",
				},
				{
					Image:     "fish",
					Word:      "ryba",
					Translate: "https://cdn0.wideopenpets.com/wp-content/uploads/2019/10/Fish-Names-770x405.png",
				},
				{
					Image:     "https://scx1.b-cdn.net/csz/news/800/2019/11-scientiststr.jpg",
					Word:      "frog",
					Translate: "żaba",
				},
				{
					Image:     "https://upload.wikimedia.org/wikipedia/commons/thumb/a/a1/Domestic_Goose_%282%29.jpg/1200px-Domestic_Goose_%282%29.jpg",
					Word:      "goose",
					Translate: "gęś",
				},
				{
					Image:     "https://a57.foxnews.com/static.foxnews.com/foxnews.com/content/uploads/2019/12/931/524/horse.jpg?ve=1&tl=1",
					Word:      "horse",
					Translate: "koń",
				},
				{
					Image:     "https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/the-lion-king-mufasa-simba-1554901700.jpg?crop=0.535xw:1.00xh;0.121xw,0&resize=480:*",
					Word:      "lion",
					Translate: "lew",
				},
				{
					Image:     "https://cdn.vox-cdn.com/thumbor/Or0rhkc1ciDqjrKv73IEXGHtna0=/0x0:666x444/1200x800/filters:focal(273x193:379x299)/cdn.vox-cdn.com/uploads/chorus_image/image/59384673/Macaca_nigra_self-portrait__rotated_and_cropped_.0.jpg",
					Word:      "monkey",
					Translate: "małpa",
				},
				{
					Image:     "https://images.unsplash.com/photo-1516632664305-eda5d6a5bb99?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&w=1000&q=8",
					Word:      "rabbit",
					Translate: "królik",
				},
			}}, 10)
		courseDB.SaveCourse(&internal.Course{Name: "Owoce i warzywa", UserID: 1, DifficultyLevel: "normal", Language: "en", Rate: 3.0,
			Description: "W tym kursie nauczysz się nazw owoców i warzyw.",
			Image:       "https://www.ang.pl/img/slownik/fruit.jpg",
			Task: []internal.Task{
				{
					Image:     "https://cdn.pixabay.com/photo/2016/01/05/13/58/apple-1122537_960_720.jpg",
					Word:      "apple",
					Translate: "jabłko",
				},
				{
					Image:     "https://cdn.pixabay.com/photo/2016/09/10/17/47/eggplant-1659784_960_720.jpg",
					Word:      "aubergine",
					Translate: "bakłażan",
				},
				{
					Image:     "https://cdn.pixabay.com/photo/2018/10/29/10/01/bananas-3780761_960_720.jpg",
					Word:      "banana",
					Translate: "banan",
				},
				{
					Image:     "https://cdn.pixabay.com/photo/2018/08/26/10/55/beans-3631986_960_720.jpg",
					Word:      "beans",
					Translate: "fasola",
				},
				{
					Image:     "https://cdn.pixabay.com/photo/2015/03/24/08/52/beetroot-687251_960_720.jpg",
					Word:      "beetroot",
					Translate: "burak",
				},
				{
					Image:     "https://cdn.pixabay.com/photo/2010/12/13/10/05/background-2277_960_720.jpg",
					Word:      "blackberry",
					Translate: "jeżyna",
				},
				{
					Image:     "https://cdn.pixabay.com/photo/2016/03/05/19/02/broccoli-1238250_960_720.jpg",
					Word:      "broccoli",
					Translate: "brokuły",
				},
				{
					Image:     "https://cdn.pixabay.com/photo/2018/10/03/21/57/cabbage-3722498_960_720.jpg",
					Word:      "cabbage",
					Translate: "kapusta",
				},
			}}, 10)
	}

	return &Database{userDB, taskDB, courseDB, courseResultsDB}, nil
}
