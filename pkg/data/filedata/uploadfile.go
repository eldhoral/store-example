package filedata

import (
    "bytes"
    "errors"
    "fmt"
    "io"
    "log"
    "math"
    "mime/multipart"
    "net/http"
    "os"
    "path/filepath"
    "regexp"
    "strings"
    "time"

    "store-api/pkg/helper/slicehelper"

    "github.com/google/uuid"
)

const (
    MB            = 1 << 20
    MAX_FILE_SIZE = 10485760 // 10 MB
    OtherDocument = iota
)

var units = []string{"B", "KB", "MB", "GB", "TB"}
var AllowedFile = []string{"pdf", "jpg", "png", "xlsx", "xls", "jpeg", "docx", "doc", "csv", "txt", "ppt", "pptx"}
var ImageExt = []string{"jpg", "png", "jpeg"}
var ImageType = []string{"image/png", "image/jpeg"}
var documentPath = map[int]string{
    OtherDocument: "./assets/upload/image/",
}

func NewMultipleUploadFile(r *http.Request, key string) ([]UploadFile, error) {
    maxSize := int64(50 * MB) // Maximum 10 MB

    err1 := r.ParseMultipartForm(maxSize)
    if err1 != nil {
        log.Println(err1)
        return nil, fmt.Errorf("Error: %s. on field %s", err1.Error(), key)
    }

    files := r.MultipartForm.File[key]

    uploadFiles := []UploadFile{}

    duplicated := make(map[string]int)

    if len(files) != 0 {
        for i := range files {
            file, err := files[i].Open()
            if err != nil {
                return nil, err
            }

            filename := strings.ToLower(files[i].Filename)

            // Handle duplicated
            _, exist := duplicated[filename]
            if exist {
                continue
            } else {
                duplicated[filename] = 1
            }

            ext := strings.ReplaceAll(filepath.Ext(filename), ".", "")

            if !slicehelper.InArray(ext, AllowedFile) {
                return nil, errors.New("File extension must be between the following types: " + strings.Join(AllowedFile, ", "))
            }

            // Get file type. "image/png", "image/jpeg"
            _fileHeader := make([]byte, files[i].Size)
            if _, err = file.Read(_fileHeader); err != nil {
                return nil, err
            }

            // Set position back to start.
            if _, err = file.Seek(0, 0); err != nil {
                return nil, err
            }
            // Avoid path traversal in S3
            // match, _ := regexp.MatchString("^[./].*$", filename)
            // if match {
            // 	return nil, errors.New(fmt.Sprintf("File name is invalid. Please use valid character: %s", filename))
            // }

            fileRe := regexp.MustCompile(`^\.*/|/|\.\..*$`)
            match := fileRe.MatchString(filename)
            if match {
                return nil, errors.New(fmt.Sprintf("File name is invalid. Please use valid character: %s", filename))
            }

            re := regexp.MustCompile(`(^\.)|(^\/)`)
            filename = re.ReplaceAllString(filename, "")

            s3name := fmt.Sprintf("%s_%s_%s",
                time.Now().Format("20060102150405"), strings.ReplaceAll(uuid.New().String(), "-", ""), filename)

            uploadFiles = append(uploadFiles, UploadFile{Name: filename, S3Name: s3name, Size: files[i].Size, Ext: ext,
                Type: http.DetectContentType(_fileHeader), File: file})
        }
    }

    return uploadFiles, nil
}

func NewUploadFileInstance(r *http.Request, name string) (*UploadFile, error) {
    maxSize := int64(10 * MB) // Maximum 10 MB

    err := r.ParseMultipartForm(maxSize)
    if err != nil {
        return nil, fmt.Errorf("%s.  on field %s", err.Error(), name)
    }
    file, fileHeader, err := r.FormFile(name)
    if err != nil {
        log.Println(err)
        return nil, err
    }
    defer file.Close()

    filename := strings.ToLower(fileHeader.Filename)
    ext := strings.ReplaceAll(filepath.Ext(filename), ".", "")

    // Get file type. "image/png", "image/jpeg"
    _fileHeader := make([]byte, fileHeader.Size)
    if _, err = file.Read(_fileHeader); err != nil {
        if err.Error() == "EOF" {
            return nil, fmt.Errorf("UploadFile: Empty content. on field %s", name)
        }
        return nil, err
    }

    // Set position back to start.
    if _, err = file.Seek(0, 0); err != nil {
        return nil, err
    }
    // Avoid path traversal in S3
    // match, _ := regexp.MatchString("^[./].*$", fileHeader.Filename)
    fileRe := regexp.MustCompile(`^\.*/|/|\.\..*$`)
    match := fileRe.MatchString(fileHeader.Filename)
    if match {
        return nil, fmt.Errorf("File name is invalid. Please use valid character: %s", fileHeader.Filename)
    }

    re := regexp.MustCompile(`(^\.)|(^\/)`)
    filename = re.ReplaceAllString(filename, "")

    s3name := fmt.Sprintf("%s_%s_%s",
        time.Now().Format("20060102150405"), strings.ReplaceAll(uuid.New().String(), "-", ""), filename)

    uniqueSuffix := uuid.New()
    destinationFile := uniqueSuffix.String() + "_" + fileHeader.Filename

    //create file and copy
    out, err := os.Create(destinationFile)
    if err != nil {
        return nil, err
    }
    defer out.Close()

    _, err = io.Copy(out, file)
    if err != nil {
        return nil, err
    }

    return &UploadFile{Name: filename,
        S3Name:     s3name,
        Size:       fileHeader.Size,
        Ext:        ext,
        Path:       destinationFile,
        Type:       http.DetectContentType(_fileHeader),
        File:       file,
        FileHeader: fileHeader}, nil
}

type UploadFile struct {
    Name   string
    S3Name string
    Size   int64
    Type   string
    Ext    string
    Path   string

    File       multipart.File
    FileHeader *multipart.FileHeader
}

func (u UploadFile) GetFile() multipart.File {
    return u.File
}

// GetFileBody get buffer for S3
func (u UploadFile) GetFileBody() *bytes.Reader {
    buffer := make([]byte, u.Size)
    u.File.Read(buffer)
    return bytes.NewReader(buffer)
}

func (u UploadFile) GetName() string {
    return u.Name
}

func (u UploadFile) GetS3Name() string {
    return u.S3Name
}

func (u UploadFile) GetSize() int64 {
    return u.Size
}

// GetType ex:  "image/png", "image/jpeg"
func (u UploadFile) GetType() string {
    return u.Type
}

func (u UploadFile) GetExtension() string {
    return u.Ext
}

func (u UploadFile) IsAllowedFile() bool {
    ext := strings.ToLower(u.Ext)
    return slicehelper.StringInSlice(ext, AllowedFile)
}

func (u UploadFile) IsAllowedImage() bool {
    return slicehelper.StringInSlice(u.Ext, ImageExt) && slicehelper.StringInSlice(u.Type, ImageType)
}

func (u UploadFile) Println() {
    fmt.Printf("Name: %s, Type: %s, Ext: %s, S3_File: %s, Size: %d\n", u.Name, u.Type, u.Ext, u.GetS3Name(), u.Size)
}

// Static public functions
func GetFileSizeText(size int64) string {
    var result float64
    var text string
    var pow int64

    bytes := math.Max(float64(size), 0)
    if bytes == 0 {
        return "0 KB"
    }

    pow = int64(math.Floor(math.Log(bytes) / math.Log(1024)))
    pow = int64(math.Min(float64(pow), float64(len(units)-1)))

    result = bytes / math.Pow(1024, float64(pow))

    if int64(result) == 0 {
        text = "0 KB"
    } else {
        text = fmt.Sprintf(`%0.2f %v`, result, units[pow])
    }
    return text
}

func IsFileExtAllowed(fileExt string) bool {
    ext := strings.ToLower(fileExt)
    ext = strings.ReplaceAll(ext, ".", "") // replace all dots to empty string in file extension
    return slicehelper.StringInSlice(ext, AllowedFile)
}

func IsImageTypes(mimeType string) bool {
    return slicehelper.StringInSlice(mimeType, ImageType)
}

func IsImageExt(filename string) bool {
    ext := strings.ReplaceAll(filepath.Ext(filename), ".", "")
    return slicehelper.StringInSlice(ext, ImageExt)
}
