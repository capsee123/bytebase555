package cmd

import (
	"github.com/bytebase/bytebase/backend/component/config"
	api "github.com/bytebase/bytebase/backend/legacyapi"
	"github.com/bytebase/bytebase/backend/plugin/app/feishu"
)

func getBaseProfile() config.Profile {
	backupStorageBackend := api.BackupStorageBackendLocal
	if flags.backupBucket != "" {
		backupStorageBackend = api.BackupStorageBackendS3
	}

	return config.Profile{
		ExternalURL:          flags.externalURL,
		GrpcPort:             flags.port + 1, // Using flags.port + 1 as our gRPC server port.
		DatastorePort:        flags.port + 2, // Using flags.port + 2 as our datastore port.
		Readonly:             flags.readonly,
		Debug:                flags.debug,
		DemoName:             flags.demoName,
		Version:              version,
		GitCommit:            gitcommit,
		PgURL:                flags.pgURL,
		DisableMetric:        flags.disableMetric,
		BackupStorageBackend: backupStorageBackend,
		BackupRegion:         flags.backupRegion,
		BackupBucket:         flags.backupBucket,
		BackupCredentialFile: flags.backupCredential,
		FeishuAPIURL:         feishu.APIPath,
	}
}
