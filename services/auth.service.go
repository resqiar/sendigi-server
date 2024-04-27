package services

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sendigi-server/configs"
	"sendigi-server/dtos"
	"sendigi-server/repos"
	"sendigi-server/utils"

	"github.com/gofiber/fiber/v2"
)

func GoogleLoginService(c *fiber.Ctx) error {
	var conf = configs.GoogleConfig

	// generate random id for state identification
	generated := utils.GenerateRandomString(32)

	sess, err := configs.StateStore.Get(c)
	if err != nil {
		log.Printf("Failed to initiate session: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to initiate session")
	}
	sess.Set("session_state", generated)

	if err := sess.Save(); err != nil {
		log.Printf("Failed to save session: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to save session")
	}

	URL := conf.AuthCodeURL(generated)

	return c.Redirect(URL)
}

func GoogleLoginCallbackService(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	sess, sessErr := configs.SessionStore.Get(c)
	stateSess, stateErr := configs.StateStore.Get(c)
	if sessErr != nil || stateErr != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Session expired or invalid")
	}

	// state from the session storage
	savedState := stateSess.Get("session_state")

	// compare the state that is coming from the callback
	// with the one that is stored inside the session storage.
	if state != savedState {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid state")
	}

	// exchange token retrieved from google with valid token
	token, err := configs.GoogleConfig.Exchange(c.Context(), code)
	if err != nil {
		log.Printf("Failed to exchange token: %v", err)
		return c.SendStatus(http.StatusUnauthorized)
	}

	// convert token to user profile data
	payload, err := utils.ConvertGoogleToken(token.AccessToken)
	if err != nil {
		return c.SendStatus(http.StatusUnauthorized)
	}

	// check if there is a user recorded with the same creds
	existingUser, _ := repos.FindUserByEmail(payload.Email)
	if existingUser == nil {
		newUser, err := CreateUser(payload)
		if err != nil {
			log.Printf("Failed to register user: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to register user")
		}

		// Store the user's id in the session
		sess.Set("ID", newUser.ID)

		// Save into memory session and.
		// saving also set a session cookie containing session_id
		if err := sess.Save(); err != nil {
			log.Printf("Failed to initiate session: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to initiate session")
		}

		return c.Redirect(os.Getenv("AUTH_REDIRECT_URL"))
	}

	// Store the user's id in the session
	sess.Set("ID", existingUser.ID)

	// Save into memory session and.
	// saving also set a session cookie containing session_id
	if err := sess.Save(); err != nil {
		log.Printf("Failed to initiate session: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to initiate session")
	}

	// Clean up state session as it's no longer needed.
	// Even though the lifetime of the stateSession is only 10 minutes,
	// cleaning up manually like this ensures there is no leak of state session, making it more secure.
	if err := stateSess.Destroy(); err != nil {
		log.Printf("Failed to delete state session: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to clean up session")
	}

	return c.Redirect(os.Getenv("AUTH_REDIRECT_URL"))
}

func SendLogout(c *fiber.Ctx) error {
	sess, err := configs.SessionStore.Get(c)
	if err != nil {
		log.Println(err.Error())
	}

	// destroy current user session
	sess.Destroy()

	return c.SendStatus(fiber.StatusOK)
}

func MobileLoginCallbackService(c *fiber.Ctx) error {
	rawPayload := c.Body()

	var payload dtos.GooglePayload

	err := json.Unmarshal([]byte(rawPayload), &payload)
	if err != nil {
		log.Println("Unmarshal Error:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// check if there is a user recorded with the same creds
	existingUser, _ := repos.FindUserByEmail(payload.Email)
	if existingUser == nil {
		newUser, err := CreateUser(&payload)
		if err != nil {
			log.Printf("Failed to register user: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to register user")
		}

		// generate JWT token
		token := utils.GenerateJWT(newUser.ID)
		return c.JSON(fiber.Map{"token": token})
	}

	token := utils.GenerateJWT(existingUser.ID)
	return c.JSON(fiber.Map{"token": token})
}
