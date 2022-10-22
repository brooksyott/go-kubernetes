package http

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/websocket/v2"
	"go-kubernetes/logger"
	"math"
	"strconv"
	"sync"
)

var Server *fiber.App
var log = *logger.Logger

type JsonResponse struct {
	HelloMessage string `json:"helloMessage"`
}

func readiness(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusOK)
}

func health(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusOK)
}

func prime(ctx *fiber.Ctx) error {
	minStr := ctx.Query("min")
	maxStr := ctx.Query("max")

	if minStr == "" {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	if maxStr == "" {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	min, _ := strconv.Atoi(minStr)
	max, _ := strconv.Atoi(maxStr)

	for min <= max {
		isPrime := true
		for i := 2; i <= int(math.Sqrt(float64(min))); i++ {
			if min%i == 0 {
				isPrime = false
				break
			}
		}
		if isPrime {
			fmt.Printf("%d ", min)
		}
		min++
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func echo(ctx *fiber.Ctx) error {
	name := ctx.Query("name")
	if name == "" {
		name = "Anonymous"
	}

	hw := &JsonResponse{
		HelloMessage: "Hi " + name,
	}

	helloBytes, err := json.MarshalIndent(hw, "", "    ")
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.Status(fiber.StatusOK).SendString(string(helloBytes))
}

func InitServer() {
	if Server != nil {
		return
	}

	Server = fiber.New()
}

func InitWebsocketControllers() {
	Server.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		// c.Locals is added to the *websocket.Conn
		//allowedString := fmt.Sprintln(c.Locals("allowed"))

		//Logger.Info(allowedString)        // true
		//Logger.Info(c.Params("id"))       // 123
		//Logger.Info(c.Query("v"))         // 1.0
		//Logger.Info(c.Cookies("session")) // ""

		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Info(fmt.Sprintf("Read error: %s", err))
				break
			}

			log.Info(fmt.Sprintf("recv: %s", msg))

			if err = c.WriteMessage(mt, msg); err != nil {
				log.Info(fmt.Sprintf("write: %s", err))
				break
			}
		}

	}))
}

func InitWebsocket() {
	Server.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
}

func InitControllers() {
	if Server == nil {
		return
	}

	Server.Get("/prime", func(ctx *fiber.Ctx) error {
		return prime(ctx)
	})

	Server.Get("/hi", func(ctx *fiber.Ctx) error {
		return echo(ctx)
	})
}

func InitMonitors() {
	if Server == nil {
		return
	}

	Server.Use(pprof.New())
	Server.Get("/monitor", monitor.New(monitor.Config{Title: "GO KUBE Metrics Page"}))
}

func InitHealth() {

	Server.Get("/readiness", func(ctx *fiber.Ctx) error {
		return readiness(ctx)
	})

	Server.Get("/health", func(ctx *fiber.Ctx) error {
		return health(ctx)
	})

}

func Start(addr string, group *sync.WaitGroup) {
	defer group.Done()
	Logger := logger.Logger
	Logger.Info("Http server starting")

	err := Server.Listen(addr)
	if err != nil {
		Logger.Info(fmt.Sprintf("Server exited with message: %s", err))
	}
}
