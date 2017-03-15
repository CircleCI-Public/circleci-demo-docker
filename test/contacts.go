package test

import (
	"github.com/circleci/cci-demo-docker/service"
	"github.com/stretchr/testify/require"
)

// SetupContact creates a Contact for use in tests. How it creates the contact is an implementation detail, but it
// should record an error if the resulting contact would be invalid.
func (env *Env) SetupContact(email string, name string) *service.Contact {
	contact, err := env.Client.AddContact(service.AddContactRequest{
		Email: email,
		Name:  name,
	})

	// VERIFY: Response contains the contact
	require.NoError(env.T, err, "Unable to get contact via API")
	require.NotEmpty(env.T, contact, "Contact not found")

	return contact
}

// ReadContactWithEmail reads a contact from the test database with the given email. Helpers like this make it easy to
// verify the state of the database as part of a test.
func (env *Env) ReadContactWithEmail(email string) *service.Contact {
	contact, err := env.DB.GetContactByEmail(email)
	require.NoError(env.T, err)

	return contact
}
