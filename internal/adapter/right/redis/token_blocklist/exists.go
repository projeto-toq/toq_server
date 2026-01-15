package tokenblocklist

import "context"

// Exists returns true when the provided JTI is blocklisted.
func (a *Adapter) Exists(ctx context.Context, jti string) (bool, error) {
	res, err := a.client.Exists(ctx, key(jti)).Result()
	if err != nil {
		return false, err
	}
	return res == 1, nil
}
