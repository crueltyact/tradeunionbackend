package s3

type Config struct {
	Bucket string `validate:"required"`
	Key    string `validate:"required"`
	Secret string `validate:"required"`
	Region string `validate:"required"`
}

