package catcher

const DefaultDomain string = "catcher.mx.ax"

// EmailHandler is a type alias for a function used by consumers
type EmailHandlerFn func(Email) error

// Producer allows us to push Emails to a consumer
type Producer interface {
	Push(Email) error
	Stop() error
}

// Consumer receives published messages and offers a way to handle them
type Consumer interface {
	Handler(EmailHandlerFn)
	Stop() error
}

// Subscription contains a channel for receiving emails and an id for unsubscribing
type Subscription struct {
	C chan Email
}

// Store is a service that accepts emails and makes them readily available
// for a configurable period of time
type Store interface {
	// Add adds Email to the store using key
	Add(string, Email) error

	// Get attempts to retrieve any emails for an email key, it always returns
	// a list of emails as we are not able to say if the key is pending emails
	Get(string) Emails

	// Has returns whether or not we have emails for a given key
	Has(string) bool

	// Subscribe lets a caller be notified whenever a new email for provided key
	// is received
	Subscribe(string) *Subscription

	// Unsubscribe removes a Subscription from the store
	Unsubscribe(string, *Subscription)
}
