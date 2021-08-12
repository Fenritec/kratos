package session_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/ory/kratos/identity"
	"github.com/ory/kratos/internal"
	"github.com/ory/kratos/session"
)

func TestSession(t *testing.T) {
	conf, _ := internal.NewFastRegistryWithMocks(t)
	authAt := time.Now()

	t.Run("case=active session", func(t *testing.T) {
		i := new(identity.Identity)
		i.State = identity.StateActive
		s, _ := session.NewActiveSession(i, conf, authAt, identity.CredentialsTypePassword)
		assert.True(t, s.IsActive())
		require.NotEmpty(t, s.Token)
		require.NotEmpty(t, s.LogoutToken)
		assert.EqualValues(t, identity.CredentialsTypePassword, s.AMR[0].Method)

		i = new(identity.Identity)
		s, err := session.NewActiveSession(i, conf, authAt, identity.CredentialsTypePassword)
		assert.Nil(t, s)
		assert.ErrorIs(t, err, session.ErrIdentityDisabled)
	})

	t.Run("case=expired", func(t *testing.T) {
		assert.False(t, (&session.Session{ExpiresAt: time.Now().Add(time.Hour)}).IsActive())
		assert.False(t, (&session.Session{Active: true}).IsActive())
	})

	t.Run("case=amr", func(t *testing.T) {
		s := session.NewInactiveSession()
		s.CompletedLoginFor(identity.CredentialsTypeOIDC)
		assert.EqualValues(t, identity.CredentialsTypeOIDC, s.AMR[0].Method)
		s.CompletedLoginFor(identity.CredentialsTypeRecoveryLink)
		assert.EqualValues(t, identity.CredentialsTypeOIDC, s.AMR[0].Method)
		assert.EqualValues(t, identity.CredentialsTypeRecoveryLink, s.AMR[1].Method)
	})
}
