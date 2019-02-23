package identity

type Key struct {
	// The GPG fingerprint of the key
	Fingerprint string `json:"fingerprint"`
	PubKey      string `json:"pub_key"`
}
