package awsutil

import (
    "bytes"
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "time"

    "store-api/pkg/helper/slicehelper"

    "github.com/h2non/filetype"

    "github.com/google/uuid"

    "store-api/pkg/data/filedata"

    "github.com/aws/aws-sdk-go/aws/awserr"
    "github.com/aws/aws-sdk-go/aws/credentials"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
    "github.com/sirupsen/logrus"

    awstrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/aws/aws-sdk-go/aws"
)

type AWSService struct {
    client *s3.S3
    bucket string
}

func NewAWSService(awsAccessKey, awsSecretKey, awsRegion, bucket string) (*AWSService, error) {
    sess, err := session.NewSession(&aws.Config{
        Region:      aws.String(awsRegion),
        Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
    })

    if err != nil {
        return nil, err
    }

    // Create S3 service client
    svc := s3.New(awstrace.WrapSession(sess))

    return &AWSService{client: svc, bucket: bucket}, nil
}

func (s AWSService) Head(folder string, filename string) (*s3.HeadObjectOutput, error) {
    result, err := s.client.HeadObject(&s3.HeadObjectInput{
        Bucket: aws.String(s.bucket),
        Key:    aws.String(folder + "/" + filename),
    })
    if err != nil {
        return nil, err
    }
    return result, nil
}

func (s AWSService) Get(folder string, filename string) (string, error) {
    req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
        Bucket: aws.String(s.bucket),
        Key:    aws.String(folder + "/" + filename),
    })

    // Create the pre-signed url with an expiry
    url, err := req.Presign(15 * time.Minute)
    if err != nil {
        return err.Error(), err
    }

    return url, nil
}

// PutRawFile   ContentType: "binary/octet-stream",
func (s AWSService) Put(folder string, filename string, body *bytes.Reader, filetype string) error {

    _, err := s.client.PutObject(&s3.PutObjectInput{
        Bucket:      aws.String(s.bucket),
        Key:         aws.String(folder + "/" + filename),
        Body:        body,
        ContentType: aws.String(filetype),
    })
    if err != nil {
        logrus.Error(fmt.Errorf("S3.Put - %v", err))
        return err
    }

    return nil
}

// PutOctetStream   ContentType: "binary/octet-stream", use for Excel
func (s AWSService) PutOctetStream(folder string, filename string, body *bytes.Reader) (string, error) {

    s3name := fmt.Sprintf("%s_%s_%s",
        time.Now().Format("20060102150405"), strings.ReplaceAll(uuid.New().String(), "-", ""), filename)

    _, err := s.client.PutObject(&s3.PutObjectInput{
        Bucket:             aws.String(s.bucket),
        Key:                aws.String(folder + "/" + s3name),
        Body:               body,
        ContentType:        aws.String("application/octet-stream"),
        ContentDisposition: aws.String(fmt.Sprintf("attachment; filename=%s", filename)),
    })
    if err != nil {
        logrus.Error(fmt.Errorf("S3.PutOctetStream - %v", err))
        return "", err
    }

    filePath, err := s.Get("tmp", s3name)
    if err != nil {
        logrus.Error(fmt.Errorf("S3.PutOctetStream Get - %v", err))
        return "", err
    }
    return filePath, nil
}

func (s AWSService) PutUploadFile(folder string, file *filedata.UploadFile) error {

    _, err := s.client.PutObject(&s3.PutObjectInput{
        Bucket:        aws.String(s.bucket),
        Key:           aws.String(folder + "/" + file.GetS3Name()),
        Body:          file.GetFileBody(),
        ContentLength: aws.Int64(file.Size),
        ContentType:   aws.String(file.Type),
    })
    if err != nil {
        logrus.Error(fmt.Errorf("S3.PutUploadFile - %v", err))
        return err
    }

    return nil
}

func (s AWSService) PutMultipleUploadFile(folder string, files []filedata.UploadFile) error {

    for _, file := range files {
        _, err := s.client.PutObject(&s3.PutObjectInput{
            Bucket:        aws.String(s.bucket),
            Key:           aws.String(folder + "/" + file.GetS3Name()),
            Body:          file.GetFileBody(),
            ContentLength: aws.Int64(file.Size),
            ContentType:   aws.String(file.Type),
        })

        if err != nil {
            logrus.Error(fmt.Errorf("S3.PutMultipleUploadFile - %v", err))
            return err
        }
    }

    return nil
}

func (s AWSService) Delete(folder string, filename string) error {
    _, err := s.client.DeleteObject(&s3.DeleteObjectInput{
        Bucket: aws.String(s.bucket),
        Key:    aws.String(folder + "/" + filename),
    })
    if err != nil {
        logrus.Error(fmt.Errorf("S3.Delete - %v", err))
        return err
    }

    return nil
}

// GetAvatar will return blank picture if avatar url is nul
func (s AWSService) GetAvatar(avatar *string, usingDefault bool) (*string, error) {
    defaultProfile := "https://talenta.s3-ap-southeast-1.amazonaws.com/avatar/blank.jpg"
    if avatar == nil {
        if usingDefault {
            return &defaultProfile, nil
        } else {
            return nil, nil
        }
    }

    avatarUrl, err := s.Get("avatar", *avatar)
    if err != nil {
        if usingDefault {
            return &defaultProfile, err
        } else {
            return nil, err
        }
    }

    return &avatarUrl, nil
}

// Check file exsits or not
func (s AWSService) Check(folder, fileName string) (bool, error) {
    key := fmt.Sprintf("%s/%s", folder, fileName)
    _, err := s.client.HeadObject(&s3.HeadObjectInput{
        Bucket: aws.String(s.bucket),
        Key:    aws.String(key),
    })
    if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            case "NotFound":
                return false, nil
            default:
                return false, err

            }
        }
        return false, err
    }
    return true, nil
}

func (s AWSService) Read(folder, fileName string) (localFileName string, fileType string, err error) {

    item := fmt.Sprintf("%s/%s", folder, fileName)

    file, err := os.Create("/tmp/" + fileName)

    if err != nil {
        return "", "", err
    }

    defer file.Close()

    sess, _ := session.NewSession(&s.client.Config)
    downloader := s3manager.NewDownloader(sess)
    _, err = downloader.Download(file,
        &s3.GetObjectInput{
            Bucket: aws.String(s.bucket),
            Key:    aws.String(item),
        },
    )

    if err != nil {
        return "", "", err
    }

    var contentType string
    specialExts := []string{".xlsx", ".xls", ".docx", ".doc", ".ppt", ".pptx"}
    ext := filepath.Ext(fileName)
    ext = strings.ToLower(ext)

    if slicehelper.StringInSlice(ext, specialExts) {
        buffer := make([]byte, 8192)
        _, err = file.Read(buffer)
        if err != nil && err != io.EOF {
            return "", "", err
        }

        kind, err := filetype.Match(buffer)
        if err != nil {
            return "", "", err
        }
        contentType = kind.MIME.Value

    } else {
        buffer := make([]byte, 512)
        _, err = file.Read(buffer)
        if err != nil && err != io.EOF {
            return "", "", err
        }

        contentType = http.DetectContentType(buffer)
    }
    return file.Name(), contentType, nil
}
