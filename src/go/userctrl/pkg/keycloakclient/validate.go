package keycloakclient

import (
	"context"
)

func (kc *KeycloakClient) ValidateToken(ctx context.Context,
	accessToken string) (bool, error) {
	res, err := kc.client.RetrospectToken(ctx, accessToken, kc.clientID, kc.clientSecret, kc.realm)
	if err != nil {
		return false, err
	}
	return *res.Active, nil
}
