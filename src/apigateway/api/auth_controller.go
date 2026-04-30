package api

import (
	"context"
	"log"
	"strconv"
	"time"

	"google.golang.org/grpc"

	pb "github.com/keshavbansal015/chirps/src/apigateway/genproto"
	"github.com/labstack/echo/v4"
)

type authController struct {
	addr string
}

func newAuthController(addr string) *authController {
	return &authController{addr}
}

// createSession creates a new session for a user
// It takes the email and password from the request body and sends them to the auth service
// to create a new session. If the session is created successfully, it sets a cookie with the
// session ID and returns the user ID in the response body.
func (ac *authController) createSession(c echo.Context) error {
	conn, err := grpc.NewClient(ac.addr, insecureCredentials())
	if err != nil {
		log.Printf("Connecting to service failed: %v", err)
		return echo.NewHTTPError(500)
	}
	defer conn.Close()
	client := pb.NewAuthServiceClient(conn) // this creates a client for the auth service

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel() // this cancels the context when the function returns

	req := pb.Credentials{
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
	}

	res, err := client.CreateSession(ctx, &req) // this calls the CreateSession method on the auth service
	if err != nil {
		log.Printf("Creating session failed: %v", err)
		return newHTTPError(err)
	}

	createCookie(c, res.Id)
	return c.JSON(200, map[string]int32{"id": res.UserId})
}

func (ac *authController) validateSession(c echo.Context) error {
	cookie, err := c.Cookie("session")
	if err != nil {
		return echo.NewHTTPError(401)
	}

	conn, err := grpc.NewClient(ac.addr, insecureCredentials())
	if err != nil {
		log.Printf("Connecting to service failed: %v", err)
		return echo.NewHTTPError(500)
	}
	defer conn.Close()
	client := pb.NewAuthServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := pb.SessionRequest{SessionId: cookie.Value}

	res, err := client.GetSession(ctx, &req)
	if err != nil {
		log.Printf("Validating session failed: %v", err)
		clearCookie(c)
		return newHTTPError(err)
	}

	createCookie(c, res.Id)
	setUserID(c, strconv.Itoa(int(res.UserId)))

	return nil
}

func (ac *authController) deleteSession(c echo.Context) error {
	cookie, err := c.Cookie("session")
	if err != nil {
		return echo.NewHTTPError(401)
	}

	conn, err := grpc.NewClient(ac.addr, insecureCredentials())
	if err != nil {
		log.Printf("Connecting to service failed: %v", err)
		return echo.NewHTTPError(500)
	}
	defer conn.Close()
	client := pb.NewAuthServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := pb.SessionRequest{SessionId: cookie.Value}

	_, err = client.DeleteSession(ctx, &req)
	if err != nil {
		log.Printf("Deleting session failed: %v", err)
		return newHTTPError(err)
	}

	clearCookie(c)
	return c.NoContent(204)
}
