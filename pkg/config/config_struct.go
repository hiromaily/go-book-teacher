package config

import (
	"github.com/hiromaily/go-book-teacher/pkg/notifier"
	"github.com/hiromaily/go-book-teacher/pkg/storage"
)

// Root is root config
type Root struct {
	Site         *Site
	Storage      *Storage
	Notification *Notification
}

// Site is site information
type Site struct {
	Type string `toml:"type" validate:"oneof=dmm"`
	URL  string `toml:"url" validate:"required"`
}

type Storage struct {
	Mode  storage.Mode `toml:"mode" validate:"oneof=text redis dummy"`
	Text  *Text        `toml:"text" validate:"required"`
	Redis *Redis       `toml:"redis" validate:"required"`
}

// Text is text storage
type Text struct {
	Path string `toml:"path" validate:"required"`
}

// Redis is redis storage
type Redis struct {
	Encrypted bool   `toml:"encrypted"`
	URL       string `toml:"url" validate:"required"`
}

type Notification struct {
	Mode    notifier.Mode `toml:"mode" validate:"required"`
	Console *Console      `toml:"console" validate:"required"`
	Browser *Browser      `toml:"browser" validate:"required"`
	Slack   *Slack        `toml:"slack" validate:"required"`
	Mail    *Mail         `toml:"mail" validate:"required"`
}

// Console is command line notification
type Console struct {
	Enabled bool `toml:"enabled"`
}

// Browser is browser notification
type Browser struct {
	Enabled bool `toml:"enabled"`
}

// Slack is slack notification
type Slack struct {
	Enabled   bool   `toml:"enabled"`
	Encrypted bool   `toml:"encrypted"`
	Key       string `toml:"key" validate:"required"`
}

// Mail is mail notification
type Mail struct {
	Enabled   bool   `toml:"enabled"`
	Encrypted bool   `toml:"encrypted"`
	MailTo    string `toml:"mail_to" validate:"required"`
	MailFrom  string `toml:"mail_from" validate:"required"`
	SMTP      *SMTP  `toml:"smtp" validate:"required"`
}

// SMTP is smtp server
type SMTP struct {
	Address string `toml:"address" validate:"required"`
	Pass    string `toml:"pass" validate:"required"`
	Server  string `toml:"server" validate:"required"`
	Port    int    `toml:"port" validate:"required"`
}
