package master

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	"tomato/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	bolt "go.etcd.io/bbolt"
)

// resp struct
type H map[string]interface{}

// 自定义validator
type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := err.Error()
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = fmt.Sprintf("%v", he.Message)
	}
	// c.Logger().Error(err)
	c.JSON(code, H{"code": code, "message": message})
}

func Start() {
	// start nats server
	opts := &server.Options{Username: "nats", Password: "123qwe"}
	ns, err := server.NewServer(opts)
	if err != nil {
		panic(err)
	}
	go ns.Start()

	// start monitoring subscriptions(all agents) // no need anymore
	// copts := &server.ConnzOptions{SubscriptionsDetail: true}
	// for {
	// 	time.Sleep(5 * time.Second)
	// 	connz, _ := ns.Connz(copts)
	// 	subs := []string{}
	// 	if len(connz.Conns) > 0 {
	// 		for _, cz := range connz.Conns {
	// 			if len(cz.SubsDetail) > 0 {
	// 				subs = append(subs, cz.SubsDetail[0].Subject)
	// 			}
	// 		}
	// 	}
	// 	log.Println(subs)
	// }

	// start bolt storage
	db, err := bolt.Open("server.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// create agent bucket
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("AgentBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	// start each http server
	// Echo instance
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// validater
	e.Validator = &customValidator{validator: validator.New()}
	// error handler
	e.HTTPErrorHandler = customHTTPErrorHandler
	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Tomato")
	})

	e.POST("/register", func(c echo.Context) error {
		reg := &utils.Agent{}
		if err := c.Bind(reg); err != nil {
			return c.JSON(http.StatusOK, H{"code": 1, "message": err.Error()})
		}
		if err := c.Validate(reg); err != nil {
			return c.JSON(http.StatusOK, H{"code": 2, "message": err.Error()})
		}
		// check if there is duplicate
		err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("AgentBucket"))
			v := b.Get([]byte(reg.Name))
			if v != nil {
				age := &utils.Agent{}
				json.Unmarshal(v, age)
				if reg.Info.ID == age.Info.ID {
					return errors.New("duplicated agent")
				}
			}
			return nil
		})
		if err != nil {
			return c.JSON(http.StatusOK, H{"code": 3, "message": err.Error()})
		}
		// save data
		val, err := json.Marshal(reg)
		if err != nil {
			return c.JSON(http.StatusOK, H{"code": 4, "message": err.Error()})
		}
		err = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("AgentBucket"))
			err := b.Put([]byte(reg.Name), val)
			return err
		})
		if err != nil {
			return c.JSON(http.StatusOK, H{"code": 5, "message": err.Error()})
		}
		return c.JSON(http.StatusOK, H{"code": 0, "message": "register success"})
	})

	e.POST("/cmd", func(c echo.Context) error {
		cmd := &utils.CMD{}
		if err := c.Bind(cmd); err != nil {
			return c.JSON(http.StatusOK, H{"code": 1, "message": err.Error()})
		}
		if err := c.Validate(cmd); err != nil {
			return c.JSON(http.StatusOK, H{"code": 2, "message": err.Error()})
		}
		// get all agents that available in both db and request
		var agents []string
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("AgentBucket"))
			b.ForEach(func(k, _ []byte) error {
				if utils.FindValInSlice(cmd.Agents, string(k)) {
					agents = append(agents, string(k))
				}
				return nil
			})
			return nil
		})
		// start nats client
		nc, err := nats.Connect(ns.ClientURL(), nats.UserInfo("nats", "123qwe"))
		if err != nil {
			return c.JSON(http.StatusOK, H{"code": 3, "message": err.Error()})
		}
		defer nc.Close()
		// start to send cmd to agents
		var wg sync.WaitGroup
		var lock sync.Mutex
		var results []utils.Response
		for _, a := range agents {
			wg.Add(1)
			go func(agent string) {
				msg, err := nc.Request(agent, []byte(cmd.CMD), 10*time.Second)
				lock.Lock()
				if err != nil {
					results = append(results, utils.Response{Agent: agent, Msg: err.Error()})
					lock.Unlock()
					wg.Done()
				}
				results = append(results, utils.Response{Agent: agent, Msg: string(msg.Data)})
				lock.Unlock()
				wg.Done()
			}(a)
		}
		wg.Wait()
		return c.JSON(http.StatusOK, H{"code": 0, "message": results})
	})

	e.GET("/agents", func(c echo.Context) error {
		var val []string
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("AgentBucket"))
			b.ForEach(func(k, v []byte) error {
				val = append(val, string(v))
				return nil
			})
			return nil
		})
		var agents []utils.Agent
		for _, i := range val {
			a := utils.Agent{}
			json.Unmarshal([]byte(i), &a)
			agents = append(agents, a)
		}
		return c.JSON(http.StatusOK, H{"code": 0, "data": agents})
	})

	// Start server
	e.Logger.Fatal(e.Start(":1323"))

}
