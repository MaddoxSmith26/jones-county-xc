package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	"jones-county-xc/backend/db"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var queries *db.Queries

type AthleteResponse struct {
	ID             int32  `json:"id"`
	Name           string `json:"name"`
	Grade          int32  `json:"grade"`
	PersonalRecord string `json:"personalRecord"`
	Events         string `json:"events"`
}

type MeetResponse struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Date        string `json:"date"`
	Location    string `json:"location"`
	Description string `json:"description"`
}

type ResultResponse struct {
	ID           int32  `json:"id"`
	AthleteID    int32  `json:"athleteId"`
	MeetID       int32  `json:"meetId"`
	Time         string `json:"time"`
	Place        int32  `json:"place"`
	AthleteName  string `json:"athleteName,omitempty"`
	AthleteGrade int32  `json:"athleteGrade,omitempty"`
}

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "xcapp:yourpassword@tcp(127.0.0.1:3306)/jones_county_xc?parseTime=true"
	}

	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer conn.Close()

	if err = conn.Ping(); err != nil {
		log.Fatal("Database not available:", err)
	}
	log.Println("Connected to MySQL database")

	queries = db.New(conn)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/health", HealthCheck)
	r.GET("/api/athletes", GetAthletes)
	r.GET("/api/athletes/:id", GetAthleteByID)
	r.POST("/api/athletes", CreateAthlete)
	r.DELETE("/api/athletes/:id", DeleteAthlete)
	r.GET("/api/meets", GetMeets)
	r.GET("/api/meets/:id/results", GetMeetResults)
	r.POST("/api/meets", CreateMeet)
	r.GET("/api/results", GetResults)
	r.POST("/api/results", CreateResult)
	r.GET("/api/top-times", GetTopTimes)

	r.Run(":8080")
}

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}

func GetAthletes(c *gin.Context) {
	athletes, err := queries.GetAllAthletes(context.Background())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	response := make([]AthleteResponse, len(athletes))
	for i, a := range athletes {
		response[i] = AthleteResponse{
			ID:             a.ID,
			Name:           a.Name,
			Grade:          a.Grade,
			PersonalRecord: a.PersonalRecord.String,
			Events:         a.Events.String,
		}
	}
	c.JSON(200, response)
}

func GetAthleteByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid athlete ID"})
		return
	}

	athlete, err := queries.GetAthleteByID(context.Background(), int32(id))
	if err != nil {
		c.JSON(404, gin.H{"error": "Athlete not found"})
		return
	}

	response := AthleteResponse{
		ID:             athlete.ID,
		Name:           athlete.Name,
		Grade:          athlete.Grade,
		PersonalRecord: athlete.PersonalRecord.String,
		Events:         athlete.Events.String,
	}
	c.JSON(200, response)
}

func CreateAthlete(c *gin.Context) {
	var input struct {
		Name           string `json:"name" binding:"required"`
		Grade          int32  `json:"grade" binding:"required"`
		PersonalRecord string `json:"personalRecord"`
		Events         string `json:"events"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result, err := queries.CreateAthlete(context.Background(), db.CreateAthleteParams{
		Name:           input.Name,
		Grade:          input.Grade,
		PersonalRecord: sql.NullString{String: input.PersonalRecord, Valid: input.PersonalRecord != ""},
		Events:         sql.NullString{String: input.Events, Valid: input.Events != ""},
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	c.JSON(201, gin.H{"id": id, "message": "Athlete created"})
}

func DeleteAthlete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid athlete ID"})
		return
	}

	if err := queries.DeleteAthlete(context.Background(), int32(id)); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Athlete deleted"})
}

func GetMeets(c *gin.Context) {
	meets, err := queries.GetAllMeets(context.Background())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	response := make([]MeetResponse, len(meets))
	for i, m := range meets {
		response[i] = MeetResponse{
			ID:          m.ID,
			Name:        m.Name,
			Date:        m.Date.Format("2006-01-02"),
			Location:    m.Location.String,
			Description: m.Description.String,
		}
	}
	c.JSON(200, response)
}

func GetMeetResults(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid meet ID"})
		return
	}

	results, err := queries.GetResultsByMeetID(context.Background(), int32(id))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	response := make([]ResultResponse, len(results))
	for i, r := range results {
		response[i] = ResultResponse{
			ID:           r.ID,
			AthleteID:    r.AthleteID,
			MeetID:       r.MeetID,
			Time:         r.Time,
			Place:        r.Place.Int32,
			AthleteName:  r.AthleteName,
			AthleteGrade: r.AthleteGrade,
		}
	}
	c.JSON(200, response)
}

func CreateMeet(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Date        string `json:"date" binding:"required"`
		Location    string `json:"location"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	date, err := parseDate(input.Date)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	result, err := queries.CreateMeet(context.Background(), db.CreateMeetParams{
		Name:        input.Name,
		Date:        date,
		Location:    sql.NullString{String: input.Location, Valid: input.Location != ""},
		Description: sql.NullString{String: input.Description, Valid: input.Description != ""},
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	c.JSON(201, gin.H{"id": id, "message": "Meet created"})
}

func GetResults(c *gin.Context) {
	meetID := c.Query("meetId")
	if meetID != "" {
		id, err := strconv.Atoi(meetID)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid meet ID"})
			return
		}
		results, err := queries.GetResultsByMeetID(context.Background(), int32(id))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		response := make([]ResultResponse, len(results))
		for i, r := range results {
			response[i] = ResultResponse{
				ID:           r.ID,
				AthleteID:    r.AthleteID,
				MeetID:       r.MeetID,
				Time:         r.Time,
				Place:        r.Place.Int32,
				AthleteName:  r.AthleteName,
				AthleteGrade: r.AthleteGrade,
			}
		}
		c.JSON(200, response)
		return
	}
	c.JSON(400, gin.H{"error": "meetId query parameter required"})
}

func CreateResult(c *gin.Context) {
	var input struct {
		AthleteID int32  `json:"athleteId" binding:"required"`
		MeetID    int32  `json:"meetId" binding:"required"`
		Time      string `json:"time" binding:"required"`
		Place     int32  `json:"place"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result, err := queries.CreateResult(context.Background(), db.CreateResultParams{
		AthleteID: input.AthleteID,
		MeetID:    input.MeetID,
		Time:      input.Time,
		Place:     sql.NullInt32{Int32: input.Place, Valid: input.Place > 0},
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	c.JSON(201, gin.H{"id": id, "message": "Result created"})
}

func GetTopTimes(c *gin.Context) {
	results, err := queries.GetTopTimes(context.Background())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	type TopTimeResponse struct {
		ID           int32  `json:"id"`
		AthleteID    int32  `json:"athleteId"`
		AthleteName  string `json:"athleteName"`
		AthleteGrade int32  `json:"athleteGrade"`
		MeetID       int32  `json:"meetId"`
		MeetName     string `json:"meetName"`
		MeetDate     string `json:"meetDate"`
		Time         string `json:"time"`
		Place        int32  `json:"place"`
	}

	response := make([]TopTimeResponse, len(results))
	for i, r := range results {
		response[i] = TopTimeResponse{
			ID:           r.ID,
			AthleteID:    r.AthleteID,
			AthleteName:  r.AthleteName,
			AthleteGrade: r.AthleteGrade,
			MeetID:       r.MeetID,
			MeetName:     r.MeetName,
			MeetDate:     r.MeetDate.Format("2006-01-02"),
			Time:         r.Time,
			Place:        r.Place.Int32,
		}
	}
	c.JSON(200, response)
}

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}
