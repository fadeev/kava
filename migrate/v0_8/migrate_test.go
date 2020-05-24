package v0_8

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	v032tendermint "github.com/kava-labs/kava/migrate/v0_8/tendermint/v0_32"
	v033tendermint "github.com/kava-labs/kava/migrate/v0_8/tendermint/v0_33"
	"github.com/stretchr/testify/require"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/kava-labs/kava/app"
)

func TestMain(m *testing.M) {
	config := sdk.GetConfig()
	app.SetBech32AddressPrefixes(config)
	app.SetBip44CoinType(config)
	config.Seal()

	os.Exit(m.Run())
}

func TestMigrate_Auth_BaseAccount(t *testing.T) {
	bz, err := ioutil.ReadFile(filepath.Join("testdata", "auth-base-old.json"))
	require.NoError(t, err)
	oldAppState := genutil.AppMap{"auth": bz}

	newAppState := MigrateAppState(oldAppState)

	bz, err = ioutil.ReadFile(filepath.Join("testdata", "auth-base-new.json"))
	require.NoError(t, err)
	require.JSONEq(t, string(bz), string(newAppState["auth"]))
}
func TestMigrate_Auth_MultiSigAccount(t *testing.T) {
	bz, err := ioutil.ReadFile(filepath.Join("testdata", "auth-base-multisig-old.json"))
	require.NoError(t, err)
	oldAppState := genutil.AppMap{"auth": bz}

	newAppState := MigrateAppState(oldAppState)

	bz, err = ioutil.ReadFile(filepath.Join("testdata", "auth-base-multisig-new.json"))
	require.NoError(t, err)
	require.JSONEq(t, string(bz), string(newAppState["auth"]))
}

func TestMigrate_Auth_ValidatorVestingAccount(t *testing.T) {
	bz, err := ioutil.ReadFile(filepath.Join("testdata", "auth-valvesting-old.json"))
	require.NoError(t, err)
	oldAppState := genutil.AppMap{"auth": bz}

	newAppState := MigrateAppState(oldAppState)

	bz, err = ioutil.ReadFile(filepath.Join("testdata", "auth-valvesting-new.json"))
	require.NoError(t, err)
	require.JSONEq(t, string(bz), string(newAppState["auth"]))
}

func TestMigrate_Auth_ModuleAccount(t *testing.T) {
	bz, err := ioutil.ReadFile(filepath.Join("testdata", "auth-module-old.json"))
	require.NoError(t, err)
	oldAppState := genutil.AppMap{"auth": bz}

	newAppState := MigrateAppState(oldAppState)

	bz, err = ioutil.ReadFile(filepath.Join("testdata", "auth-module-new.json"))
	require.NoError(t, err)
	require.JSONEq(t, string(bz), string(newAppState["auth"]))
}

func TestMigrate_Auth_PeriodicVestingAccount(t *testing.T) {
	bz, err := ioutil.ReadFile(filepath.Join("testdata", "auth-periodic-old.json"))
	require.NoError(t, err)
	oldAppState := genutil.AppMap{"auth": bz}

	newAppState := MigrateAppState(oldAppState)

	bz, err = ioutil.ReadFile(filepath.Join("testdata", "auth-periodic-new.json"))
	require.NoError(t, err)
	require.JSONEq(t, string(bz), string(newAppState["auth"]))
}

func TestMigrateTendermint(t *testing.T) {
	oldGenDoc, err := v032tendermint.GenesisDocFromFile(filepath.Join("testdata", "tendermint-old.json"))
	require.NoError(t, err)

	newGenDoc := v033tendermint.Migrate(*oldGenDoc)

	expectedGenDoc, err := tmtypes.GenesisDocFromFile(filepath.Join("testdata", "tendermint-new.json"))
	require.NoError(t, err)
	require.Equal(t, *expectedGenDoc, newGenDoc)
}

func TestMigrateDistribution(t *testing.T) {
	bz, err := ioutil.ReadFile(filepath.Join("testdata", "distribution-old.json"))
	require.NoError(t, err)
	oldAppState := genutil.AppMap{"distribution": bz}

	newAppState := MigrateSDK(oldAppState)

	bz, err = ioutil.ReadFile(filepath.Join("testdata", "distribution-new.json"))
	require.NoError(t, err)
	require.JSONEq(t, string(bz), string(newAppState["distribution"]))
}

func TestMigrateStaking(t *testing.T) {
	bz, err := ioutil.ReadFile(filepath.Join("testdata", "staking-old.json"))
	require.NoError(t, err)
	oldAppState := genutil.AppMap{"staking": bz}

	newAppState := MigrateSDK(oldAppState)

	bz, err = ioutil.ReadFile(filepath.Join("testdata", "staking-new.json"))
	require.NoError(t, err)
	require.JSONEq(t, string(bz), string(newAppState["staking"]))
}

func TestMigrateSlashing(t *testing.T) {
	bz, err := ioutil.ReadFile(filepath.Join("testdata", "slashing-old.json"))
	require.NoError(t, err)
	oldAppState := genutil.AppMap{"slashing": bz}

	newAppState := MigrateSDK(oldAppState)

	bz, err = ioutil.ReadFile(filepath.Join("testdata", "slashing-new.json"))
	require.NoError(t, err)
	require.JSONEq(t, string(bz), string(newAppState["slashing"]))
}

func TestMigrateEvidence(t *testing.T) {
	bz, err := ioutil.ReadFile(filepath.Join("testdata", "slashing-old.json"))
	require.NoError(t, err)
	oldAppState := genutil.AppMap{"slashing": bz}

	newAppState := MigrateSDK(oldAppState)

	bz, err = ioutil.ReadFile(filepath.Join("testdata", "evidence-new.json"))
	require.NoError(t, err)
	require.JSONEq(t, string(bz), string(newAppState["evidence"]))
}