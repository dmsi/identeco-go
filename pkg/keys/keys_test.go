package keys

import (
	"testing"

	"github.com/dmsi/identeco/pkg/s3helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetJWKS(t *testing.T) {
	keyService := KeyService{
		S3:             s3helper.NewMockSession(),
		Bucket:         "test-bucket",
		JWKSObjectName: "jwks.json",
	}
	path := keyService.Bucket + "/" + keyService.JWKSObjectName
	storage := keyService.S3.Downloader.(*s3helper.S3Mock).Data
	storage[path] = []byte(
		`{
			"keys": [
				{
					"e": "AQAB",
					"kid": "XooolbD0BPGABjHzSDRfQ4YBg8H87zwTJVmmP8I81OA",
					"kty": "RSA",
					"n": "tRXzVqY51HMCh-iK2K0YmGF044P2qM_42MDBZuk6CpqUg1Vm7ylBHLm41QWNIwvzyVtBiibjSPtT_Ua2-_6v5dz2bwZqUzxYU_yq5sacv3yfOpwe8mYej2wyaC0fBcKSigrpFj3nDHTXEUGIiR0Vptd7ja7vjOcj_8raGjaR7zGF_5P42OA-UUDmRmyU1PG_d4fV-bagip1byEcPM4GSxqOnWkJdNX9da82S9QxYSofFq9t8MYH2texM5ImcqZ0FmdUXb8k1DeBXv0dqg1ZbhaDvCzNWfgoMjhPeB5lpnCP0gR-X_3dLJDPI1lU0ddnjepCWuh48WuImxfilaoQCcw",
					"alg": "RS256",
					"use": "sig"
				}
			]
		}`)

	j, err := keyService.GetJWKS()
	require.ErrorIs(t, err, nil)
	require.Equal(t, 1, len(j.Keys))
	assert.Equal(t, "AQAB", j.Keys[0].E)
	assert.Equal(t, "XooolbD0BPGABjHzSDRfQ4YBg8H87zwTJVmmP8I81OA", j.Keys[0].Kid)
	assert.Equal(t, "RSA", j.Keys[0].Kty)
	assert.Equal(t, "tRXzVqY51HMCh-iK2K0YmGF044P2qM_42MDBZuk6CpqUg1Vm7ylBHLm41QWNIwvzyVtBiibjSPtT_Ua2-_6v5dz2bwZqUzxYU_yq5sacv3yfOpwe8mYej2wyaC0fBcKSigrpFj3nDHTXEUGIiR0Vptd7ja7vjOcj_8raGjaR7zGF_5P42OA-UUDmRmyU1PG_d4fV-bagip1byEcPM4GSxqOnWkJdNX9da82S9QxYSofFq9t8MYH2texM5ImcqZ0FmdUXb8k1DeBXv0dqg1ZbhaDvCzNWfgoMjhPeB5lpnCP0gR-X_3dLJDPI1lU0ddnjepCWuh48WuImxfilaoQCcw", j.Keys[0].N)
	assert.Equal(t, "RS256", j.Keys[0].Alg)
	assert.Equal(t, "sig", j.Keys[0].Use)
}
