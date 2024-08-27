package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	atr "gitlab.com/parking-finder/parking-finder-api/internal/automate_test/routers"
	ats "gitlab.com/parking-finder/parking-finder-api/internal/automate_test/services"

	"github.com/joho/godotenv"
	ar "gitlab.com/parking-finder/parking-finder-api/internal/admin/routers"
	as "gitlab.com/parking-finder/parking-finder-api/internal/admin/services"
	cr "gitlab.com/parking-finder/parking-finder-api/internal/customer/routers"
	cs "gitlab.com/parking-finder/parking-finder-api/internal/customer/services"
	"gitlab.com/parking-finder/parking-finder-api/internal/httpclient"
	ms "gitlab.com/parking-finder/parking-finder-api/internal/message"
	ns "gitlab.com/parking-finder/parking-finder-api/internal/notification"
	pms "gitlab.com/parking-finder/parking-finder-api/internal/payment"
	rss "gitlab.com/parking-finder/parking-finder-api/internal/reserve"

	pr "gitlab.com/parking-finder/parking-finder-api/internal/provider/routers"
	ps "gitlab.com/parking-finder/parking-finder-api/internal/provider/services"
	ss "gitlab.com/parking-finder/parking-finder-api/internal/search"
	wr "gitlab.com/parking-finder/parking-finder-api/internal/webhook"

	_ "gitlab.com/parking-finder/parking-finder-api/pkg/utility"

	"gitlab.com/parking-finder/parking-finder-api/pkg/connector"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Get args
	var port string
	wd, _ := os.Getwd()
	if err := godotenv.Load(wd + "/.env"); err != nil {
		fmt.Println("Error loading .env file.")
	}

	flag.StringVar(&port, "p", "", "a string")
	flag.Parse()
	fmt.Println("Start api.")
	// Config echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper: func(c echo.Context) bool {
			path := c.Request().URL.Path
			return strings.Contains(path, "/bots")
		},
		AllowOrigins: []string{"http://localhost:3100", "http://" + os.Getenv("WEB_HOST")},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPut,
			http.MethodPatch,
			http.MethodPost,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	// Load ENV

	redis := connector.NewRedisClient(
		"",
		os.Getenv("REDIS_ADDRESS"),
		os.Getenv("REDIS_PASSWORD"),
	)
	if _, err := redis.Ping().Result(); err != nil {
		log.Fatalln("Error connecting to Redis")
	}

	// Connect and test DB connection
	mongodb := connector.NewMongoDBClient(os.Getenv("MONGODB_URI"), 100)
	if err := mongodb.Ping(); err != nil {
		log.Fatalln("Error connecting to MongoDB cluster")
	}

	// Select DB
	db := mongodb.SelectDB(os.Getenv("DATABASE_NAME"))
	h := httpclient.NewHTTPClient(http.Client{})

	// Setup new service
	customerServices := cs.NewCustomerServices(db, redis, h)
	searchServices := ss.NewSearchServices(db)
	notificationServices := ns.NewNotificationServices(db)
	paymentServices := pms.NewPaymentService(db, h, notificationServices)
	messageServices := ms.NewMessageServices(db)
	reserveServices := rss.NewReserveServices(db, h, notificationServices)
	automateTestServices := ats.NewAutomateTeserServices(db)

	providerServices := ps.NewProviderServices(db, redis, h)

	adminServices := as.NewAdminServices(db)

	// Setup api group
	gCustomer := e.Group("/customer")
	gProvider := e.Group("/provider")
	gAdmin := e.Group("/admin")
	gWebhook := e.Group("/webhook")
	gAutomateTest := e.Group("/automate_test")

	// Setup new router
	cr.NewCustomerRouter(gCustomer, customerServices, searchServices, paymentServices, reserveServices, messageServices, notificationServices)
	pr.NewProviderRouter(gProvider, providerServices, notificationServices, reserveServices, customerServices, messageServices, h)
	ar.NewAdminRouter(gAdmin, adminServices, notificationServices)
	wr.NewWebhookRouter(gWebhook, customerServices, paymentServices, reserveServices, notificationServices, h, redis)
	atr.NewAutomateTestRouter(gAutomateTest, automateTestServices)

	// Show all routes api
	fmt.Println("All registered routes.")
	data := e.Routes()
	for i := 0; i < len(data); i++ {
		fmt.Printf("Method: %s, Path: %s\n", data[i].Method, data[i].Path)
	}

	// Setup port
	if port == "" {
		port = os.Getenv("PORT")
	} else {
		fmt.Println("User selected port: " + port)
	}

	// Start scheduler
	fmt.Println("Start scheduler.")

	// Start server
	e.Logger.Fatal(e.Start(":" + port))
}
