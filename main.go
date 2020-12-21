package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var app = &cli.App{
	Name:  "covid-alert",
	Usage: "Runs the covid-alert set of tools",
	Commands: []*cli.Command{
		reportCommand,
	},
}

var reportCommand = &cli.Command{
	Name:  "report",
	Usage: "Run covid-alert report",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "db",
			Usage:   "database file",
			EnvVars: []string{"DB", "DB_FILE"},
		},
		&cli.StringFlag{
			Name:    "twilio_sid",
			Usage:   "Twilio account SID",
			EnvVars: []string{"TWILIO_SID"},
		},
		&cli.StringFlag{
			Name:    "twilio_auth_token",
			Usage:   "Twilio auth token",
			EnvVars: []string{"TWILIO_AUTH_TOKEN"},
		},
		&cli.StringFlag{
			Name:    "twilio_from",
			Usage:   "Twilio sender phone number",
			EnvVars: []string{"TWILIO_FROM"},
		},
		&cli.StringSliceFlag{
			Name:    "twilio_to",
			Usage:   "Twilio recipient numbers",
			EnvVars: []string{"TWILIO_TO"},
		},
	},
	Action: func(c *cli.Context) error {
		db, err := getDB(c.String("db"))
		if err != nil {
			return fmt.Errorf("get db: %w", err)
		}

		twilioCfg := &TwilioConfig{
			AccountSID: c.String("twilio_sid"),
			AuthToken:  c.String("twilio_auth_token"),
			From:       c.String("twilio_from"),
			To:         c.StringSlice("twilio_to"),
		}
		twilio := NewTwilio(twilioCfg)

		err = runReporter(db, twilio)
		if err != nil {
			return fmt.Errorf("run reporter: %w", err)
		}

		return nil
	},
}

func main() {
	log.Logger = log.With().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("app run failure")
	}
}

func runReporter(db *sqlx.DB, twilio *Twilio) error {
	firstTime := false

	// Get latest report
	latest, err := getLatestReport(db)
	if errors.Is(err, sql.ErrNoRows) {
		firstTime = true
	} else if err != nil {
		return fmt.Errorf("get latest report: %w", err)
	} else {
		log.Info().Msgf("last report (%s) - total: %d, ICU: %d, inpatient: %d",
			latest.Time.Format("2006/01/02 03:04PM MST"), latest.Total(), latest.ICUUnits, latest.InpatientUnits)
	}

	// Get new report
	report, err := getPageReport()
	if err != nil {
		return fmt.Errorf("get page report: %w", err)
	}

	log.Info().Msgf("new report (%s) - total: %d, ICU: %d, inpatient: %d",
		report.Time.Format("2006/01/02 03:04PM MST"), report.Total(), report.ICUUnits, report.InpatientUnits)

	// Save new report
	err = saveReport(report, db)
	if err != nil {
		return fmt.Errorf("save report: %w", err)
	}

	if firstTime {
		return nil
	}

	// Check changes
	if latest.HasChangedFrom(report) {
		err := twilio.PublishReportChanges(latest, report)
		if err != nil {
			return fmt.Errorf("publish report changes: %w", err)
		}
	}

	return nil
}
