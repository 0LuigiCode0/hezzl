package dclickhouse

import "time"

type LogEventGood struct {
	Id          int       `json:"id" ch:"id"`
	ProjectId   int       `json:"project_id" ch:"project_id"`
	Name        string    `json:"name" ch:"name"`
	Description string    `json:"description" ch:"description"`
	Priority    int       `json:"priority" ch:"priority"`
	Removed     bool      `json:"removed" ch:"removed"`
	EventTime   time.Time `json:"event_time" ch:"event_time"`
}

func (leg *LogEventGood) PrepareBatch() []any {
	return []any{leg.Id, leg.ProjectId, leg.Name, leg.Description, leg.Priority, leg.Removed, leg.EventTime}
}
