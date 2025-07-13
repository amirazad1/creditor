package model

import (
	"time"
)

type City struct {
	BaseModel
	Name      string `gorm:"not null;unique"`
	Buildings []Building
}

type Building struct {
	BaseModel
	Name     string `gorm:"not null"`
	CityId   int    `gorm:"not null"`
	City     City   `gorm:"foreignKey:CityId;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
	Projects []Project
}

type Project struct {
	BaseModel
	Name       string   `gorm:"not null"`
	BuildingId int      `gorm:"not null"`
	Building   Building `gorm:"foreignKey:BuildingId;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
}

type CostTypeParent struct {
	BaseModel
	Name      string `gorm:"not null;unique"`
	CostTypes []CostType
}

type CostType struct {
	BaseModel
	Name             string         `gorm:"not null"`
	CostTypeParentId int            `gorm:"not null"`
	CostTypeParent   CostTypeParent `gorm:"foreignKey:CostTypeParentId;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
	CostSubTypes     []CostSubType
}

type CostSubType struct {
	BaseModel
	Name       string   `gorm:"not null"`
	CostTypeId int      `gorm:"not null"`
	CostType   CostType `gorm:"foreignKey:CostTypeId;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
}

type Year struct {
	BaseModel
	Name      string    `gorm:"not null;unique"`
	Year      int       `gorm:"type:int;uniqueIndex;not null"`
	StartAt   time.Time `gorm:"type:TIMESTAMP with time zone;not null;unique"`
	EndAt     time.Time `gorm:"type:TIMESTAMP with time zone;not null;unique"`
	Additions []Addition
}

type CostDescription struct {
	BaseModel
	Name string `gorm:"not null;unique"`
}

type Plan struct {
	BaseModel
	Name         string `gorm:"not null;unique"`
	No           string `gorm:"not null"`
	Date         string `gorm:"not null"`
	PlanProjects []PlanProject
}

type PlanProject struct {
	BaseModel
	PlanId    int
	Plan      Plan    `gorm:"foreignKey:PlanId;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
	ProjectId int     `gorm:"not null"`
	Project   Project `gorm:"foreignKey:ProjectId;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
	Amount    int     `gorm:"not null"`
}

type Addition struct {
	BaseModel
	Name   string  `gorm:"not null"`
	Amount float32 `gorm:"not null"`
	YearId int     `gorm:"not null"`
	Year   Year    `gorm:"foreignKey:YearId;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
}
