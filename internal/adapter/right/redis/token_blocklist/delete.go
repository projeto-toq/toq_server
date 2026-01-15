package tokenblocklist

import "context"

// Delete removes a JTI from the blocklist.
func (a *Adapter) Delete(ctx context.Context, jti string) error {
	_, err := a.client.Del(ctx, key(jti)).Result()
	return err
}
