package weekplan

import (
	"time"

	"github.com/peppermint-recipes/peppermint-server/recipe"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// WeekPlanDay is a model of a Day for a weekplan
type Day struct {
	Breakfast []recipe.Recipe `json:"breakfast"`
	Lunch     []recipe.Recipe `json:"lunch"`
	Snack     []recipe.Recipe `json:"snack"`
	Dinner    []recipe.Recipe `json:"dinner"`
}

// WeekPlan is a model of a week
type WeekPlan struct {
	ID           primitive.ObjectID `json:"id"`
	UserID       string             `json:"userId"`
	Year         int                `json:"year"`
	CalendarWeek int                `json:"calendarWeek"`
	Monday       Day                `json:"monday"`
	Tuesday      Day                `json:"tuesday"`
	Wednesday    Day                `json:"wednesday"`
	Thursday     Day                `json:"thursday"`
	Friday       Day                `json:"friday"`
	Saturday     Day                `json:"saturday"`
	Sunday       Day                `json:"sunday"`
	Deleted      bool               `json:"deleted"`
	LastUpdated  time.Time          `json:"lastUpdated"`
}

func (wp *WeekPlan) isValid() bool {
	if wp.UserID == "" {
		return false
	}

	if wp.CalendarWeek <= 0 || wp.CalendarWeek >= 53 {
		return false
	}

	if wp.Year <= 2020 || wp.Year >= 2100 {
		return false
	}

	return true
}
