package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"learn/httpserver/model"
	"learn/httpserver/repo"
	"learn/httpserver/setup"
	"learn/httpserver/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"golang.org/x/oauth2"
)

func GetAllDataFromRedis(c *gin.Context) ([]model.User, error) {

	var allData []model.User

	ctxRedisClient, ctxExist := c.Get("redis-client")
	if !ctxExist {
		return allData, errors.New("redis client not available")
	}

	redisClient, ok := ctxRedisClient.(*redis.Client)
	if !ok {
		panic("not a valid redis client")
	}

	getMapData, err := redisClient.HGetAll(c, "user:23").Result()
	if err != nil {
		return allData, err
	}

	for k, v := range getMapData {
		var mapData model.User
		fmt.Println(k)
		err := json.Unmarshal([]byte(v), &mapData)
		if err != nil {
			panic(err)
		}
		allData = append(allData, mapData)
	}
	return allData, nil

}

func Get(c *gin.Context) {
	var getData []model.GetUser
	var err error

	allData, getErr := GetAllDataFromRedis(c)

	if getErr == nil || allData == nil {
		DB := setup.ConnectDB()
		repos := repo.UserRepo(DB)
		getData, err = repos.GetData()
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"getData-database": getData,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"getMapData-redis": allData,
	})
}

func PutAllDataInRedis(redisClient *redis.Client) {
	DB := setup.ConnectDB()
	repos := repo.UserRepo(DB)
	getData, err := repos.GetData()
	if err != nil {
		panic(err)
	}

	for k, v := range getData {
		fmt.Println(k, v)
		redisKey := strconv.Itoa(int(v.Id))
		err := redisClient.HSet(context.Background(), "user:"+redisKey, "Id", v.Id, "Email", v.Email, "Name", v.Name, "Age", v.Age, "Address", v.Address).Err()
		if err != nil {
			fmt.Printf("HSet Error: %s", err)
		}
	}
}

func Create(c *gin.Context) {
	DB := setup.ConnectDB()
	repos := repo.UserRepo(DB)
	//check data
	var user model.User
	err := c.BindJSON(&user)
	if err != nil {
		panic(err)
	}

	tx, err := DB.Begin(c)
	if err != nil {
		// return err
		log.Fatal("Error in transaction begin : ", err)
		return
	}

	err = repos.CreateEmployee(user, tx)
	if err != nil {
		tx.Rollback(c)
		c.JSON(500, gin.H{
			"isCreated": false,
		})
		return
	}

	err = repos.CreateEmployeeServicePair(user.Id, user.Sid, tx)
	if err != nil {
		tx.Rollback(c)
		c.JSON(500, gin.H{
			"isCreated": false,
		})
		return
	}

	err = tx.Commit(c)

	if err != nil {
		log.Fatal(err)
		c.JSON(500, gin.H{
			"isCreated": false,
			"message ":  "transaction failed",
		})
		return
	}

	fmt.Println("-- transaction committed --")

	c.JSON(http.StatusOK, gin.H{
		"isCreated": true,
		"vaid-data": user,
	})
}

// ------------------------AssignNewServiceToUser -------------------------
func AssignNewServiceToUser(c *gin.Context) {
	DB := setup.ConnectDB()
	repos := repo.ServiceRepo(DB)

	var NewService model.Service
	err := c.BindJSON(&NewService)
	if err != nil {
		panic(err)
	}

	isCreated, creationError := repos.CreateNewService(NewService)
	if creationError != nil {
		panic(creationError)
	}

	c.JSON(http.StatusOK, gin.H{
		"isCreatedService": isCreated,
		"created-service":  "successfully created service",
	})
}

func Delete(c *gin.Context) {
	DB := setup.ConnectDB()
	repos := repo.UserRepo(DB)
	var id = c.Param("id")

	isCreated, deletionError := repos.DeleteData(id)
	if deletionError != nil {
		panic(deletionError)
	}

	//redis-client
	ctxRedisClient, redisConnected := c.Get("redis-client")
	redisClient := ctxRedisClient.(*redis.Client)

	if redisConnected {
		delErr := redisClient.Del(c, "user:"+id).Err()
		if delErr != nil {
			fmt.Printf("HSet delete Error: %s", delErr)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"isDeleted":                 isCreated,
		"delete-response-database ": "deleted from database",
	})
}

func Update(c *gin.Context) {
	DB := setup.ConnectDB()
	repos := repo.UserRepo(DB)

	id := c.Param("id")

	var user model.User
	err := c.BindJSON(&user)
	if err != nil {
		panic(err)
	}

	isUpdated, updationError := repos.UpdateData(user, id)
	if updationError != nil {
		panic(updationError)
	}

	ctxRedisClient, redisConnected := c.Get("redis-client")
	redisClient := ctxRedisClient.(*redis.Client)

	if redisConnected {
		createErr := redisClient.HSet(c, "user:"+id, "Id", id, "Email", user.Email, "Name", user.Name, "Age", user.Age, "Address", user.Address, "Sid", user.Sid).Err()
		if createErr != nil {
			fmt.Printf("HSet create Error: %s", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"isUpdated": isUpdated,
	})

}

func RefreshToken(c *gin.Context) {
	token, isAvailable := c.Get("token")
	if !isAvailable {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Token not found",
		})
		return
	}
	DB := setup.ConnectDB()
	repos := repo.UserRepo(DB)
	refreshToken, err := repos.RefreshToken(token.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Refresh token Generation failed",
		})
		return
	}
	accessToken, err := utils.GenerateJwtToken(refreshToken, time.Now().Add(time.Second*30))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Access token Generation failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"refreshToken": refreshToken,
		"accessToken":  accessToken,
	})
}

// HydraPublicPortCall
func HydraPublicPortCall(c *gin.Context) {
	clientId := os.Getenv("HYDRA_CLIENT_ID") // Replace with your client ID
	redirectURI := os.Getenv("REDIRECT_URL")
	scope := os.Getenv("HYDRA_SCOPE")
	state := os.Getenv("HYDRA_STATE") // Replace with a secure random state
	oauthUrl := os.Getenv("HYDRA_AUTH_URL")

	params := url.Values{
		"client_id":     {clientId},
		"client_secret": {"demosecret"},
		"prompt":        {"consent"},
		"response_type": {"code"},
		"redirect_uri":  {redirectURI},
		"scope":         {scope},
		"state":         {state},
	}

	encodedUrl := oauthUrl + "?" + params.Encode()

	fmt.Println("encodedUrl: ", encodedUrl)

	// Redirect user to Hydra authorization endpoint
	c.Redirect(http.StatusFound, encodedUrl)

	//for api test purpose
	// c.JSON(http.StatusOK, gin.H{
	// 	"message": "successfully redirected to hydra authorization endpoint",
	// 	"encodedUrl": encodedUrl,
	// })
}

func (h Handler) AuthGetLogin(c *gin.Context) {

	login_challenge := c.Query("login_challenge")
	if login_challenge == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "login_challenge not found",
		})
		return
	}

	fmt.Println("login_challenge: ", login_challenge)

	ctx := c.Request.Context()

	// Using Hydra Admin to get the login challenge info
	loginGetParams := admin.NewGetLoginRequestParams()
	loginGetParams.WithContext(ctx)
	loginGetParams.SetLoginChallenge(login_challenge)

	// get login request from hydra
	respLoginGet, err := h.HydraAdmin.GetLoginRequest(loginGetParams)

	if err != nil {
		//redirect to login page with error
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error in getting login request",
		})
		return
	}

	fmt.Println("respLoginGet login: ", respLoginGet)

	skip := false
	if respLoginGet.GetPayload().Skip != nil {
		skip = *respLoginGet.GetPayload().Skip
	}

	// If hydra was already able to authenticate the user, skip will be true and we do not need to re-authenticate
	// the user.
	if skip {
		// Using Hydra Admin to accept login request!
		loginAcceptParams := admin.NewAcceptLoginRequestParams()
		loginAcceptParams.WithContext(ctx)
		loginAcceptParams.SetLoginChallenge(login_challenge)
		loginAcceptParams.SetBody(&models.AcceptLoginRequest{
			Subject: respLoginGet.GetPayload().Subject,
		})

		// accept login request
		respLoginAccept, err := h.HydraAdmin.AcceptLoginRequest(loginAcceptParams)
		if err != nil {
			//redirect to login page with error
			c.JSON(http.StatusBadRequest, gin.H{
				//Redirect to login page with error
				"message": "error in accepting login request",
			})
			return
		}

		// Redirect user to hydra consent endpoint
		c.Redirect(http.StatusFound, *respLoginAccept.Payload.RedirectTo)
		return
	}

	loginURL := os.Getenv("CLIENT_LOGIN_URL")

	redirectURL := loginURL + "?login_challenge=" + login_challenge
	c.Redirect(http.StatusFound, redirectURL)
}

// AuthPostLogin
func (h Handler) AuthPostLogin(c *gin.Context) {

	// login_challenge := c.Query("login_challenge")
	login_challenge := c.PostForm("login_challenge")
	if login_challenge == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "login_challenge not found",
		})
		return
	}

	fmt.Println("login_challenge post login : ", login_challenge)

	ctx := c.Request.Context()

	// Using Hydra Admin to get the login challenge info
	loginGetParams := admin.NewGetLoginRequestParams()
	loginGetParams.WithContext(ctx)
	loginGetParams.SetLoginChallenge(login_challenge)

	// get login request from hydra
	respLoginGet, err := h.HydraAdmin.GetLoginRequest(loginGetParams)

	if err != nil {
		//redirect to login page with error
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error in getting login request",
		})
		return
	}

	fmt.Println("respLoginGet post login: ", respLoginGet)

	subject := "static-userid-set-from-login-in-subject"

	loginAcceptParams := admin.NewAcceptLoginRequestParams()
	loginAcceptParams.WithContext(ctx)
	loginAcceptParams.SetLoginChallenge(login_challenge)
	loginAcceptParams.SetBody(&models.AcceptLoginRequest{
		Subject:  &subject,
		Remember: true,
	})

	respLoginAccept, err := h.HydraAdmin.AcceptLoginRequest(loginAcceptParams)
	if err != nil {
		//redirect to login page with error
		c.JSON(http.StatusBadRequest, gin.H{
			//Redirect to login page with error
			"message": "error in accepting login request :" + err.Error(),
		})
		return
	}

	// Redirect user to hydra consent endpoint
	c.Redirect(http.StatusFound, *respLoginAccept.GetPayload().RedirectTo)

	//for api test purpose
	// c.JSON(http.StatusOK, gin.H{
	// 	"message": "successfully login",
	// 	"respLoginAccept": respLoginAccept,
	// })
}

// AuthGetConsent
func (h Handler) AuthGetConsent(c *gin.Context) {

	consent_challenge := c.Query("consent_challenge")
	if consent_challenge == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "consent_challenge not found",
		})
		return
	}

	// fmt.Println("consent_challenge: ", consent_challenge)

	ctx := c.Request.Context()

	// Using Hydra Admin to get the consent challenge info
	consentGetParams := admin.NewGetConsentRequestParams()
	consentGetParams.SetContext(c)
	consentGetParams.SetConsentChallenge(consent_challenge)

	consentGetResp, err := h.HydraAdmin.GetConsentRequest(consentGetParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error in getting consent request",
		})
		return
	}

	// If a user has granted this application the requested scope, hydra will tell us to not show the UI.
	if consentGetResp.GetPayload().Skip {
		// Now it's time to grant the consent request.
		// You could also deny the request if something went terribly wrong
		consentAcceptBody := &models.AcceptConsentRequest{
			GrantAccessTokenAudience: consentGetResp.GetPayload().RequestedAccessTokenAudience,
			GrantScope:               consentGetResp.GetPayload().RequestedScope,
		}

		// Using Hydra Admin to accept consent request
		consentAcceptParams := admin.NewAcceptConsentRequestParams()
		consentAcceptParams.WithContext(ctx)
		consentAcceptParams.SetConsentChallenge(consent_challenge)
		consentAcceptParams.WithBody(consentAcceptBody)

		consentAcceptResp, err := h.HydraAdmin.AcceptConsentRequest(consentAcceptParams)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "error in accepting consent request",
			})
			return
		}

		c.Redirect(http.StatusFound, *consentAcceptResp.GetPayload().RedirectTo)
		return
	}

	consentMessage := fmt.Sprintf("Application %s wants access resources on your behalf and to:",
		consentGetResp.GetPayload().Client.ClientName,
	)

	fmt.Println("consentMessage: ", consentMessage)

	// Redirect user to hydra consent endpoint
	redirectURL := os.Getenv("CLIENT_CONSENT_URL") + "?consent_challenge=" + consent_challenge + "&scope=" + consentGetResp.GetPayload().Client.Scope + "&appName" + consentGetResp.GetPayload().Client.ClientName
	c.Redirect(http.StatusFound, redirectURL)
}

// AuthPostConsent
func (h Handler) AuthPostConsent(c *gin.Context) {

	consent_challenge := c.PostForm("consent_challenge")
	if consent_challenge == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "consent_challenge not found",
		})
		return
	}

	ctx := c.Request.Context()

	// Using Hydra Admin to get the consent challenge info
	consentGetParams := admin.NewGetConsentRequestParams()
	consentGetParams.WithContext(ctx)
	consentGetParams.SetConsentChallenge(consent_challenge)

	consentGetResp, err := h.HydraAdmin.GetConsentRequest(consentGetParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error in getting consent request",
		})
		return
	}

	// If a user has granted this application the requested scope, hydra will tell us to not show the UI.
	// if consentGetResp.GetPayload().Skip {
	// Now it's time to grant the consent request.
	// You could also deny the request if something went terribly wrong
	consentAcceptBody := &models.AcceptConsentRequest{
		GrantAccessTokenAudience: consentGetResp.GetPayload().RequestedAccessTokenAudience,
		GrantScope:               consentGetResp.GetPayload().RequestedScope,
		Session: &models.ConsentRequestSession{
			AccessToken: map[string]string{
				"permission": "static-permission-set-from-consent",
			},
		},
	}

	// Using Hydra Admin to accept consent request
	consentAcceptParams := admin.NewAcceptConsentRequestParams()
	consentAcceptParams.WithContext(ctx)
	consentAcceptParams.SetConsentChallenge(consent_challenge)
	consentAcceptParams.WithBody(consentAcceptBody)

	consentAcceptResp, err := h.HydraAdmin.AcceptConsentRequest(consentAcceptParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error in accepting consent request",
		})
		return
	}

	//for ui flow redirect
	c.Redirect(http.StatusFound, *consentAcceptResp.GetPayload().RedirectTo)
	return

	//for api test purpose
	// c.JSON(http.StatusOK, gin.H{
	// 	"message": "successfully consent",
	// 	"consentAcceptResp": consentAcceptResp,
	// })
}

// ----------------- token endpoint ----------------------------

// HydraTokenEndpoint
func (h Handler) HydraTokenEndpoint(c *gin.Context) {
	// Endpoint is OAuth 2.0 endpoint.
	 Endpoint := oauth2.Endpoint{
		AuthURL:  os.Getenv("HYDRA_AUTH_URL"),
		TokenURL: os.Getenv("HYDRA_TOKEN_URL"),
	}

	 redirect_uri := os.Getenv("REDIRECT_URL")

	 hydra_client_id := os.Getenv("HYDRA_CLIENT_ID")
	 hydra_client_secret := os.Getenv("HYDRA_CLIENT_SECRET")

	// Scopes: OAuth 2.0 scopes provide a way to limit the amount of access that is granted to an access token.
	  OAuthConf := &oauth2.Config{
		RedirectURL:  redirect_uri,
		ClientID:     hydra_client_id,
		ClientSecret: hydra_client_secret,
		Scopes:       []string{"users.write", "users.read", "users.edit", "users.delete", "offline"},
		Endpoint:     Endpoint,
	}

	code := c.PostForm("code")
	refreshToken := c.PostForm("refresh_token")
	// if refresh token is available then generate new token
	if refreshToken != "" {
		tokenRespWithRefresh, err := OAuthConf.TokenSource(context.Background(), &oauth2.Token{RefreshToken: refreshToken}).Token()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "error in sending request to token endpoint" + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":   "successfully token generated with refresh token",
			"tokenResp": tokenRespWithRefresh,
		})
		return
	}

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "authorization code not found",
		})
		return
	}

	// if refresh token is not available then generate new token
	tokenResp, err := OAuthConf.Exchange(context.Background(), code)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error in sending request to token endpoint" + err.Error(),
		})
		return
	}

	fmt.Println("tokenResp from token endpoint: ", tokenResp)

	//set in cookie
	// c.SetCookie("access_token", tokenResp.AccessToken, 3600, "/", "localhost", false, true)
	// c.SetCookie("refresh_token", tokenResp.RefreshToken, 3600, "/", "localhost", false, true)

	// for ui flow redirect
	redirectURL := os.Getenv("REDIRECT_URL") + "?access_token=" + tokenResp.AccessToken + "&refresh_token=" + tokenResp.RefreshToken + "&token_type=" + tokenResp.TokenType + "&expires_in=" + tokenResp.Expiry.String()
	c.Redirect(http.StatusFound, redirectURL)

	// for api test purpose
	// c.JSON(http.StatusOK, gin.H{
	// 	"message":   "successfully new token generated",
	// 	"tokenResp": tokenResp,
	// })

}

// ----------------------- token endpoint ----------------------------

// HydraIntroSpectEndpoint
func (h Handler) HydraIntroSpectEndpoint(c *gin.Context) {
	token := c.Request.Header.Get("token")
	fmt.Println("access-token: ", token)

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "token not found",
		})
		return
	}

	ctx := c.Request.Context()

	introspectParams := admin.NewIntrospectOAuth2TokenParams()
	introspectParams.WithContext(ctx)
	introspectParams.SetToken(token)

	introspectResp, err := h.HydraAdmin.IntrospectOAuth2Token(introspectParams)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error in introspect token",
		})
		return
	}

	// check valid token and authorize user
	if *introspectResp.Payload.Active {
		c.JSON(http.StatusOK, gin.H{
			"message":         "token is valid",
			"hydraintrospect": introspectResp.Payload,
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"message": "token is not valid, please login again",
	})
}

// Protect Test
func ProtectTest(c *gin.Context) {
	auth_user := c.Request.Header.Get("auth-user")
	auth_permission := c.Request.Header.Get("auth-permission")
	c.JSON(http.StatusOK, gin.H{
		"message":                           "successfully authenticated and authorized by oathkeeper",
		"auth-user-set-by-oathkeeper":       auth_user,
		"auth-permission-set-by-oathkeeper": auth_permission,
	})
}

// public test
func PublicTest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "This is response from public api which is not protected by oathkeeper",
	})
}

// Login
func Login(c *gin.Context) {

	login_challenge := c.Query("login_challenge")
	if login_challenge == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "login_challenge not found",
		})
		return
	}

	fmt.Println("login_challenge: ", login_challenge)

	// Using Hydra Admin to get the login challenge info
	loginGetParams := admin.NewGetLoginRequestParams()
	loginGetParams.SetContext(c)
	loginGetParams.SetLoginChallenge(login_challenge)

	var loginUserData model.Login
	err := c.BindJSON(&loginUserData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}
	//database operation
	DB := setup.ConnectDB()
	repos := repo.UserRepo(DB)

	//if data is not in redis hit database
	loggedInStatus, loggedInError := repos.CheckUserExist(loginUserData)
	if loggedInError != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": loggedInError.Error(),
		})
		return
	}

	// get refresh token
	// check in db if refresh token is expired then generate new refresh token and store in db
	loginUser, refreshToken, err := repos.UserLogin(loginUserData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	accessToken, err := utils.GenerateJwtToken(loginUser.Email, time.Now().Add(time.Second*30))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":                loggedInStatus,
		"refreshToken":           refreshToken,
		"accessToken":            accessToken,
		"message-db-login-check": "successfully loggedIn",
	})
}

// GetEmployeeData
func GetEmployeeData(c *gin.Context) {
	DB := setup.ConnectDB()
	repos := repo.UserRepo(DB)
	getData, err := repos.GetData()
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message":   "successfully authenticated and authorized",
		"Auth Data": getData,
	})
}

// Logout
func Logout(c *gin.Context) {
	var userData model.User
	err := c.BindJSON(&userData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}
	DB := setup.ConnectDB()
	repos := repo.UserRepo(DB)
	_, err = repos.DeleteData(userData.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	// delete token from db
	_, err = DB.Exec(context.Background(), `DELETE FROM refreshtokens WHERE email=$1`, userData.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "successfully logout",
	})
}

// AuthData
func AuthData(c *gin.Context) {
	DB := setup.ConnectDB()
	//repositories initialization
	repos := repo.UserRepo(DB)
	getData, err := repos.GetData()
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message":   "successfully authenticated and authorized",
		"Auth Data": getData,
	})
}

// SessionTest
func SessionTest(c *gin.Context) {
	DB := setup.ConnectDB()
	//repositories initialization
	repos := repo.UserRepo(DB)
	getData, err := repos.GetData()
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"session-test-message": "successfully test session",
		"AllData":              getData,
	})
}
