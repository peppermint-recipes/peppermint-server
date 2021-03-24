package weekplan

import (
	"time"

	"github.com/peppermint-recipes/peppermint-server/recipe"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// WeekPlanDay is a model of a day for a weekplan
type day struct {
	Breakfast []recipe.Recipe `json:"breakfast"`
	Lunch     []recipe.Recipe `json:"lunch"`
	Snack     []recipe.Recipe `json:"snack"`
	Dinner    []recipe.Recipe `json:"dinner"`
}

// weekPlan is a model of a week
type weekPlan struct {
	ID           primitive.ObjectID `json:"id"`
	UserID       string             `json:"userId"`
	Year         int                `json:"year"`
	CalendarWeek int                `json:"calendar_week"`
	Monday       day                `json:"monday"`
	Tuesday      day                `json:"tuesday"`
	Wednesday    day                `json:"wednesday"`
	Thursday     day                `json:"thursday"`
	Friday       day                `json:"friday"`
	Saturday     day                `json:"saturday"`
	Sunday       day                `json:"sunday"`
	Deleted      bool               `json:"deleted"`
	LastUpdated  time.Time          `json:"lastUpdated"`
}

func (wp *weekPlan) isValid() bool {
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
