package db

import (
	"errors"
	"github.com/codemicro/spacetraders/internal/config"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Conn *gorm.DB

func init() {
	db, err := gorm.Open(sqlite.Open(config.C.DatabaseFile), &gorm.Config{
		Logger: nil,
	})
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	Conn = db

	if err = Conn.AutoMigrate(&Ship{}, &Market{}); err != nil {
		log.Fatal().Err(err).Msg("Failed to automigrate")
	}

	log.Info().Msg("Database connected")
}

type Ship struct {
	ID   string `gorm:"primarykey"`
	Type int
	Data string
}

func GetShip(id string) (*Ship, bool, error) {
	var sh Ship
	sh.ID = id
	err := Conn.Take(&sh).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}
	return &sh, true, nil
}

func (s *Ship) Create() error {
	return Conn.Create(s).Error
}

func CountShips() (int, error) {
	var n int64
	return int(n), Conn.Model(&Ship{}).Count(&n).Error
}

type Market struct {
	gorm.Model
	Location string
	Data     string
}

func RecordMarketData(location string, data string) error {
	return Conn.Create(&Market{
		Location: location,
		Data:     data,
	}).Error
}

func GetLatestDataForLocation(location string) (*Market, bool, error) {
	var ma Market
	ma.Location = location
	err := Conn.Where(&ma).Order("created_at desc").Take(&ma).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}
	return &ma, true, nil
}

func GetMarketLocations() ([]string, error) {
	var locations []string
	return locations, Conn.Model(&Market{}).Distinct().Pluck("location", &locations).Error
}