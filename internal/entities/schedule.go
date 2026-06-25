package entities

type ScheduleMinute struct {
	JobID  uint `gorm:"column:job_id;primaryKey"`
	Minute int  `gorm:"column:minute;primaryKey"`
}

func (ScheduleMinute) TableName() string {
	return "schedule_minute"
}

type ScheduleHour struct {
	JobID uint `gorm:"column:job_id;primaryKey"`
	Hour  int  `gorm:"column:hour;primaryKey"`
}

func (ScheduleHour) TableName() string {
	return "schedule_hour"
}

type ScheduleMday struct {
	JobID uint `gorm:"column:job_id;primaryKey"`
	Mday  int  `gorm:"column:mday;primaryKey"`
}

func (ScheduleMday) TableName() string {
	return "schedule_mday"
}

type ScheduleWday struct {
	JobID uint `gorm:"column:job_id;primaryKey"`
	Wday  int  `gorm:"column:wday;primaryKey"`
}

func (ScheduleWday) TableName() string {
	return "schedule_wday"
}

type ScheduleMonth struct {
	JobID uint `gorm:"column:job_id;primaryKey"`
	Month int  `gorm:"column:month;primaryKey"`
}

func (ScheduleMonth) TableName() string {
	return "schedule_month"
}

type Schedule struct {
	JobID uint

	Minute []ScheduleMinute
	Hour   []ScheduleHour
	Mday   []ScheduleMday
	Month  []ScheduleMonth
	Wday   []ScheduleWday
}
