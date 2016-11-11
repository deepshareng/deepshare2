package backup

type MongoBackupCompressStorageFormat struct {
	Key string   `bson:"key"`
	Val []string `bson:"val"`
}
