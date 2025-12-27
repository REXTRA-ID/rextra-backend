package export

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type ExportService interface {
	GenerateCSV(data []map[string]interface{}, filename string) (string, error)
	GenerateExcel(data []map[string]interface{}, filename, sheetName string) (string, error)
	GeneratePDF(data []map[string]interface{}, filename, title string) (string, error)
}

type exportService struct {
	csvExporter   CSVExporter
	excelExporter ExcelExporter
	pdfExporter   PDFExporter
}

// New creates a default export service with CSV, Excel, and PDF exporters.
func New() ExportService {
	return &exportService{
		csvExporter:   NewCSVExporter(),
		excelExporter: NewExcelExporter(),
		pdfExporter:   NewPDFExporter(),
	}
}

func (s *exportService) GenerateCSV(data []map[string]interface{}, filename string) (string, error) {
	return s.generateAndUpload(func() (string, error) {
		return s.csvExporter.Generate(data, filename)
	})
}

func (s *exportService) GenerateExcel(data []map[string]interface{}, filename, sheetName string) (string, error) {
	return s.generateAndUpload(func() (string, error) {
		return s.excelExporter.Generate(data, filename, sheetName)
	})
}

func (s *exportService) GeneratePDF(data []map[string]interface{}, filename, title string) (string, error) {
	return s.generateAndUpload(func() (string, error) {
		return s.pdfExporter.Generate(data, filename, title)
	})
}

func (s *exportService) generateAndUpload(generate func() (string, error)) (string, error) {
	filePath, err := generate()
	if err != nil {
		return "", err
	}

	uploadFn := s.buildUploader()
	if uploadFn == nil {
		return filePath, nil
	}

	url, err := uploadFn(filePath)
	_ = os.Remove(filePath) // cleanup best-effort
	if err != nil {
		return "", err
	}

	return url, nil
}

func (s *exportService) buildUploader() func(string) (string, error) {
	bucket := os.Getenv("S3_BUCKET")
	region := os.Getenv("AWS_REGION")
	if bucket == "" || region == "" {
		return nil
	}

	accessKey := os.Getenv("AWS_ACCESS_KEY")
	secretKey := os.Getenv("AWS_SECRET_KEY")

	return func(filePath string) (string, error) {
		cfg, err := config.LoadDefaultConfig(context.Background(),
			config.WithRegion(region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		)
		if err != nil {
			return "", err
		}

		client := s3.NewFromConfig(cfg)

		file, err := os.Open(filePath)
		if err != nil {
			return "", err
		}
		defer file.Close()

		objectKey := filepath.Join("exports", fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(filePath)))

		_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(objectKey),
			Body:   file,
		})
		if err != nil {
			return "", err
		}

		publicURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, objectKey)
		return publicURL, nil
	}
}
