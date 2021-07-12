package shared

import (
	"fmt"
	"log"
	"sync"
	"database/sql"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nuid"
)

// Component is contains reusable logic related to handling
// of the connection to NATS and Database in the system.
type Component struct {
	// cmu is the lock from the component.
	cmu sync.Mutex

	// id is a unique identifier used for this component.
	id string

	// nc is the connection to NATS.
	nc *nats.Conn

	// db is the connection to DB.
	db *sql.DB

	// kind is the type of component.
	kind string
}

// NewComponent creates a
func NewComponent(kind string) *Component {
	id := nuid.Next()
	return &Component{
		id:   id,
		kind: kind,
	}
}

// SetupConnectionToDB creates a connection to the database
func (c *Component) SetupConnectionToDB(dbDriver string, connectionString string) error {
	c.cmu.Lock()
	defer c.cmu.Unlock()
	db, err := sql.Open(dbDriver, connectionString)
	if err != nil {
		panic(err.Error())
	}
	c.db = db
	return err
}

// Component returns the current Database connection.
func (c *Component) DB() *sql.DB {
	c.cmu.Lock()
	defer c.cmu.Unlock()
	return c.db
}

// SetupConnectionToNATS connects to NATS and registers the event
// callbacks and makes it available for discovery requests as well.
func (c *Component) SetupConnectionToNATS(servers string, options ...nats.Option) error {
	// Label the connection with the kind and id from component.
	options = append(options, nats.Name(c.Name()))

	c.cmu.Lock()
	defer c.cmu.Unlock()

	// Connect to NATS with customized options.
	nc, err := nats.Connect(servers, options...)
	if err != nil {
		return err
	}
	c.nc = nc

	// Setup NATS event callbacks
	//
	// Handle protocol errors and slow consumer cases.
	nc.SetErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
		log.Printf("NATS error: %s\n", err)
	})
	nc.SetReconnectHandler(func(_ *nats.Conn) {
		log.Println("Reconnected to NATS!")
	})
	nc.SetDisconnectHandler(func(_ *nats.Conn) {
		log.Println("Disconnected from NATS!")
	})
	nc.SetClosedHandler(func(_ *nats.Conn) {
		panic("Connection to NATS is closed!")
	})

	return err
}

// NATS returns the current NATS connection.
func (c *Component) NATS() *nats.Conn {
	c.cmu.Lock()
	defer c.cmu.Unlock()
	return c.nc
}

// ID returns the ID from the component.
func (c *Component) ID() string {
	c.cmu.Lock()
	defer c.cmu.Unlock()
	return c.id
}

// Name is the label used to identify the NATS connection.
func (c *Component) Name() string {
	c.cmu.Lock()
	defer c.cmu.Unlock()
	return fmt.Sprintf("%s:%s", c.kind, c.id)
}

// Shutdown makes the component go away.
func (c *Component) Shutdown() error {
	c.NATS().Close()
	defer c.DB().Close()
	return nil
}
