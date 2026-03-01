package main

// importing the project controller and service packages
import (
	"context" // 
	"strconv" // strconv for converting string to integer
	
	"net/http" // http for making HTTP requests to Supabase API
	"io" // io for reading response body from Supabase API
	"os" //os for file handling

	"github.com/lengzuo/supa" // supa for interacting with Supabase services
	"fmt" // fmt for printing error messages

	"github.com/gin-gonic/gin" // gin framework

	"gitlab.com/pragmaticreviews/golang-gin-poc/controller"
	"gitlab.com/pragmaticreviews/golang-gin-poc/service"
	"gitlab.com/pragmaticreviews/golang-gin-poc/entity"

	"log" // log for logging error messages
	"encoding/json" // encoding/json for marshaling user data
	"github.com/lengzuo/supa/dto" // dto for sign up and sign in request body
	"github.com/joho/godotenv" // godotenv for loading environment variables
	"github.com/jackc/pgx/v5" // pgx for connecting to Supabase database (if needed in the future)
)

// importing project controller and service
var (
	projectService    service.ProjectService // projectService is the service layer for handling project-related business logic
	projectController controller.ProjectController // projectController is the controller layer for handling HTTP requests related to projects
	supaClient        *supabase.Client // supaClient is the Supabase client for interacting with Supabase services
)

func main() {

	// load environment variables from .env file
	err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

	// Connect to the Supabase database using pgx
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	// Ensure the database connection is closed when the main function exits
	defer conn.Close(context.Background())

	// Query the database version to verify the connection
	var version string
	if err := conn.QueryRow(context.Background(), "SELECT version()").Scan(&version); err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	log.Println("Connected to:", version)

	// Initialize the project service and controller with database connection
	projectService = service.NewProjectService(conn)
	projectController = controller.New(projectService)

	// Supabase configuration
	conf := supabase.Config{
		ApiKey: os.Getenv("SUPABASE_API_KEY"),
		ProjectRef: os.Getenv("SUPABASE_PROJECT_REF"),
		// set to false in production environment to disable debug logs
		Debug: false,
	}

	// create a new Supabase client using the provided configuration
	supaClient, err = supabase.New(conf)
	if err != nil {
		fmt.Println("Error creating Supabase client:", err)
	} else {
		fmt.Println("Supabase client created successfully")
	}

	// create a gin server instance
	// Change Default() to New() karena kita ingin menggunakan middleware yang kita buat sendiri, yaitu logger middleware
	server := gin.New()
	server.Use(gin.Recovery())

	// Define project routes with authentication middleware
	// The AuthMiddleware will check if the user is authenticated before allowing access to the project routes
	projectRoutes := server.Group("/projects")
	projectRoutes.Use(AuthMiddleware(supaClient)) 
    {
        projectRoutes.GET("/", getProjects)
        projectRoutes.POST("/", createProject)
        projectRoutes.DELETE("/:id", deleteProject)
        projectRoutes.PUT("/:id", updateProject)
		projectRoutes.GET("/:id", getSpesificProject)
    }

	// Define authentication routes without authentication middleware, since these routes are for logging in and signing up
	authRoutes := server.Group("/auth") // Use Gin's built-in logger for auth routes
	{
		authRoutes.POST("/login", Login)
		authRoutes.POST("/signup", SignUp)
		authRoutes.GET("/user", GetAuthUser)
		authRoutes.POST("/logout", Logout)
	}

	// run the server 
	server.Run(":8080")
}

func Login(ctx *gin.Context) {
	// The body should be "email" and "password", so we need to define a struct to bind the JSON request body
	// bind the JSON request body to a struct and pass it to the sign up method of the Supabase client
	var body struct {
		Email	string `json:"email"`
		Password string `json:"password"`
	}
	ctx.BindJSON(&body)

	// call the sign up method of the Supabase client and return the result as JSON response
	resp, err := supaClient.Auth.SignInWithPassword(ctx, dto.SignInRequest{
		Email: body.Email,
		Password: body.Password,
	})

	if err != nil {
		log.Printf("Error signing in: %v", err)
		ctx.JSON(500, gin.H{"error": "Failed to sign in"})
		return
	}

	ctx.JSON(200, gin.H{"message": "User signed in successfully", "user": resp.User, "token": resp.AccessToken})
}


func SignUp(ctx *gin.Context) {
	// bind the JSON request body to a struct and pass it to the sign up method of the Supabase client
	var body struct {
		Email	string `json:"email"`
		Password string `json:"password"`
	}
	ctx.BindJSON(&body)
	// call the sign up method of the Supabase client and return the result as JSON response
	resp, err := supaClient.Auth.SignUp(ctx, dto.SignUpRequest{
		Email: body.Email,
		Password: body.Password,
	})

	if err != nil {
		log.Printf("Error signing up: %v", err)
		ctx.JSON(500, gin.H{"error": "Failed to sign up"})
		return
	} else {
		ctx.JSON(201, gin.H{"message": "User signed up successfully", "user": resp.User})
		// confirmation email will be sent to the user's email address, so we don't need to return the user data in the response
	}
}

func GetAuthUser(ctx *gin.Context) {
	// Get the token from the Authorization header
	token := ctx.GetHeader("Authorization")
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:] // Remove "Bearer " prefix
	}
	if token == "" {
		ctx.JSON(400, gin.H{"error": "Authorization header required"})
		return
	} else {
		user, err := supaClient.Auth.User(ctx, token)
		if err != nil {
			log.Printf("Error getting user: %v", err)
			ctx.JSON(401, gin.H{"error": "Invalid token"})
			// ctx.JSON(500, gin.H{"error": "Failed to get user"})
			return
		} else {
			// Marshal the user data to JSON and return it in the response
			ctx.JSON(200, gin.H{"user": user})
		}
	}
}

func Logout(ctx *gin.Context) {
	// Get the token from the Authorization header
	token := ctx.GetHeader("Authorization")

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:] // Remove "Bearer " prefix
	}

	if token == "" {
		ctx.JSON(400, gin.H{"error": "Authorization header required"})
		return
	}

	// Supabase doesn't have a built-in logout method, so we need to make a POST request to the Supabase API to log out the user
	projectRef := os.Getenv("SUPABASE_PROJECT_REF")
	apiKey := os.Getenv("SUPABASE_API_KEY")
	supaBaseUrl := fmt.Sprintf("https://%s.supabase.co/auth/v1/logout", projectRef)

	req, err := http.NewRequestWithContext(ctx.Request.Context(), "POST", supaBaseUrl, nil)

	if err != nil {
		// log.Printf("Error creating logout request: %v", err)
		ctx.JSON(500, gin.H{"error": "Failed to create logout request"})
		return
	}

	// Set the Authorization header with the Bearer token and the API key
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("apikey", apiKey)

	// Send the request to the Supabase API
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending logout request: %v", err)
		ctx.JSON(500, gin.H{"error": "Failed to send logout request"})
		return
	}

	// Close the response body after reading it to prevent memory leaks
	defer resp.Body.Close()

	// Check the response status code to determine if the logout was successful
	if resp.StatusCode != 204 {
		body, _ := io.ReadAll(resp.Body)
		ctx.JSON(resp.StatusCode, gin.H{"error": "Failed to log out", "details": string(body)})
		return
	}

	ctx.JSON(200, gin.H{"message": "User signed out successfully"})
}

func getProjects(ctx *gin.Context) {
	// call the find all method of the project controller and return the result as JSON response
	projects := projectController.FindAll()
	ctx.JSON(200, projects)
}

func getSpesificProject(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Error converting id to integer: %v", err)
		ctx.JSON(400, gin.H{"error": "Invalid project ID"})
		return
	} 
	projects := projectController.FindAll()
	for _, project := range projects {
		if project.Id == idInt {
			ctx.JSON(200, project)
			return
		}
	}
	ctx.JSON(404, gin.H{"error": "Project not found"})
}

func createProject(ctx *gin.Context) {
	// bind the JSON request body to a project entity and pass it to the save method of the project controller
	savedProject := projectController.Save(ctx)
	ctx.JSON(201, savedProject)
}

func deleteProject(ctx *gin.Context) {
	// get the title parameter from the URL and pass it to the delete method of the project controller
	id := ctx.Param("id")
	// Convert the id string to an integer
	idInt, err := strconv.Atoi(id)

	if err != nil {
		log.Printf("Error converting id to integer: %v", err)
		ctx.JSON(400, gin.H{"error": "Invalid project ID"})
		return
	} else {
		ctx.JSON(200, gin.H{"message": "Project deleted successfully"})
	}

	projectController.Delete(idInt)

	// Cek apakah project dengan title tersebut ada atau tidak, jika tidak ada maka return 404
	ctx.Status(204)
}

func updateProject(ctx *gin.Context) {
	// get the title parameter from the URL and pass it to the delete method of the project controller
	idInt, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Printf("Error converting id to integer: %v", err)
		ctx.JSON(400, gin.H{"error": "Invalid project ID"})
		return
	}
	// bind the JSON request body to a project entity and pass it to the update method of the project controller
	var updateData entity.Project
	ctx.BindJSON(&updateData)
	project := projectController.Update(idInt, updateData)
	ctx.JSON(200, project)
}

func AuthMiddleware(supaClient *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:] // Remove "Bearer " prefix
		} 

		// In Gin, the context is 'c', but Supabase needs 'context.Background()' or 'c.Request.Context()'
		user, err := supaClient.Auth.User(c.Request.Context(), token)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set the user information in the Gin context so that it can be accessed in the handlers
		c.Set("user", user)
		c.Next()
	}
}

func signUp(supaClient *supabase.Client) {
	ctx := context.Background()
	// Sign up a new user with email and password
	// user, err := supaClient.Auth.SignUpWithEmail(email, password)
	body := dto.SignUpRequest{
		Email: "faradisy20@gmail.com",
		Password: "20Januari",
	}

	resp, err := supaClient.Auth.SignUp(ctx, body)
	if err != nil {
		log.Printf("Error signing up: %v", err)
		return
	}

    bytes, _ := json.Marshal(resp)
    fmt.Printf("sign up success: %s", bytes)
}

// Create a local struct to replace the missing 'dto' package
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func signIn(supaClient *supabase.Client) {
	// 1. You must define ctx (use context.Background() for standalone functions)
	ctx := context.Background()

	body := dto.SignInRequest{
		Email: "faradisy20@gmail.com",
		Password: "20Januari",
	}

	// 2. The library uses .Login(), not .SignIn()
	// 3. Use := to declare resp and err for the first time
	resp, err := supaClient.Auth.SignInWithPassword(ctx, body)

	if err != nil {
		log.Printf("Error signing in: %v", err)
		return
	}

	fmt.Printf("User signed in: %v\n", resp.AccessToken)
}


