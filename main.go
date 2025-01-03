package poker

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/d1nch8g/poker/gen/database"
	"github.com/d1nch8g/poker/gen/migr"
	"github.com/d1nch8g/poker/mail"
	"github.com/davecgh/go-spew/spew"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jessevdk/go-flags"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
)

//	@title			Self-hosted poker app
//	@version		1.0
//	@description	Self-hosted poker.
//	@termsOfService	http://github.com/d1nch8g/poker

//	@contact.name	Swap Support
//	@contact.url	http://github.com/d1nch8g/poker
//	@contact.email	d1nch8g@gmail.com

//	@license.name	MIT
//	@license.url	https://github.com/d1nch8g/poker/src/branch/main/LICENSE

//	@host		localhost:8080
//	@BasePath	/api
//	@schemes	http https

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Token authorization for internal operations

var opts struct {
	Port         string `long:"port" env:"PORT" default:"8080"`
	Host         string `long:"host" env:"HOST"`
	Database     string `long:"database" env:"DATABASE" default:"postgresql://user:password@localhost:5432/db?sslmode=disable"`
	AdminCreds   string `long:"admin-creds" env:"ADMIN" default:"support@chx.su:password"`
	EmailAddress string `long:"email-address" env:"EMAIL_ADDRESS" default:"mail.hosting.reg.ru"`
	EmailPort    int    `long:"email-port" env:"EMAIL_PORT" default:"587"`
	EmailCreds   string `long:"email-creds" env:"EMAIL_CREDS" default:"support@chx.su:password"`
	CertFile     string `long:"cert-file" env:"CERT_FILE"`
	KeyFile      string `long:"key-file" env:"KEY_FILE"`
	Help         bool   `short:"h" long:"help"`
}

var help = `Server parameters:
--port          - Port on which to run application on
--host          - Hostname, should be inintialized on production runs
--database      - database connection string
--admin-creds   - admin creds "email:password"
--admin-wallet  - admin wallet passphrase for bots  
--email-address - email client address
--email-port    - email client port
--email-creds   - email "login:password"
--cert-file     - Cert file path (should be used for TLS)
--key-file      - Key file path (should be used for TLS)
-h --help       - Show this help message and exit`

func main() {
	_, err := flags.NewParser(&opts, flags.IgnoreUnknown).Parse()
	if err != nil {
		panic(err)
	}
	if opts.Help {
		fmt.Println(help)
		return
	}

	spew.Dump(opts)

	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:               true,
		ForceQuote:                true,
		EnvironmentOverrideColors: true,
	})

	s := bindata.Resource(migr.AssetNames(), migr.Asset)

	d, err := bindata.WithInstance(s)
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithSourceInstance(
		"go-bindata", d, opts.Database,
	)
	if err != nil {
		panic(err)
	}

	err = m.Up()
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			panic(err)
		}
	}

	conn, err := pgxpool.New(context.Background(), opts.Database)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	e := echo.New()
	sqlc := database.New(conn)
	mail := mail.New(
		opts.EmailAddress,
		strings.Split(opts.EmailCreds, ":")[0],
		strings.Split(opts.EmailCreds, ":")[1],
		"https://"+opts.Host,
		opts.EmailPort,
	)

	hasher := sha512.New()
	hasher.Write([]byte(strings.Split(opts.AdminCreds, ":")[1]))
	passhash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	_, err = sqlc.CreateUser(context.Background(), database.CreateUserParams{
		Name:            "Main Admin",
		Email:           strings.Split(opts.AdminCreds, ":")[0],
		Verified:        true,
		IsAdmin:         true,
		EncryptedWallet: "",
		Passwhash:       passhash,
	})
	if err != nil && !strings.Contains(err.Error(), "duplicate key value violates unique constraint ") {
		panic(err)
	}

	
}
