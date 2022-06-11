package s3

import (
	"context"
	"sync"

	"github.com/yunyandz/tiktok-demo-backend/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
)

type Mys3 struct {
	s3 *s3.Client
}

var (
	s3Client *Mys3
	once     sync.Once
)

func New(config *config.Config, logger *zap.Logger) *Mys3 {
	if !config.S3.Vaild {
		return nil
	}
	once.Do(func() {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           config.S3.Endpoint,
				SigningRegion: config.S3.Region,
			}, nil
		})
		cfg, err := awscfg.LoadDefaultConfig(context.TODO(), awscfg.WithEndpointResolverWithOptions(customResolver))
		if err != nil {
			panic(err)
		}
		cfg.Region = config.S3.Region
		cfg.Credentials = aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(config.S3.AccessKey, config.S3.SecretKey, ""))
		client := s3.NewFromConfig(cfg)
		// 检查bucket是否存在，不存在则创建
		burkets, err := client.ListBuckets(context.TODO(), nil)
		if err != nil {
			panic(err)
		}
		hasbucket := false
		for _, bucket := range burkets.Buckets {
			if *bucket.Name == config.S3.Bucket {
				logger.Info("bucket", zap.String("find burket", *bucket.Name))
				hasbucket = true
			}
		}
		if !hasbucket {
			_, err := client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
				Bucket: aws.String(config.S3.Bucket),
			})
			if err != nil {
				panic(err)
			}
		}
		s3Client = &Mys3{
			s3: client,
		}
	})
	return s3Client
}

func (s *Mys3) PutObject(ctx context.Context, input *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	return s.s3.PutObject(ctx, input, optFns...)
}

func (s *Mys3) PresignGetObject(ctx context.Context, input *s3.GetObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
	ps := s3.NewPresignClient(s.s3)
	return ps.PresignGetObject(ctx, input, optFns...)
}

type S3ObjectAPI interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	PresignGetObject(ctx context.Context, input *s3.GetObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}
