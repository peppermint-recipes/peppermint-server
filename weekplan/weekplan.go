package weekplan

import (
	"github.com/peppermint-recipes/peppermint-server/recipe"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// WeekPlanDay is a model of a day for a weekplan
type day struct {
	Breakfast []recipe.Recipe `json:"breakfast,omitempty"`
	Lunch     []recipe.Recipe `json:"lunch,omitempty"`
	Snack     []recipe.Recipe `json:"snack,omitempty"`
	Dinner    []recipe.Recipe `json:"dinner,omitempty"`
}

// weekPlan is a model of a week
type weekPlan struct {
	ID           primitive.ObjectID `json:"id,omitempty"`
	UserID       string             `json:"user_id,omitempty"`
	Year         int                `json:"year,omitempty"`
	CalendarWeek int                `json:"calendar_week,omitempty"`
	Monday       day                `json:"monday,omitempty"`
	Tuesday      day                `json:"tuesday,omitempty"`
	Wednesday    day                `json:"wednesday,omitempty"`
	Thursday     day                `json:"thursday,omitempty"`
	Friday       day                `json:"friday,omitempty"`
	Saturday     day                `json:"saturday,omitempty"`
	Sunday       day                `json:"sunday,omitempty"`
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
