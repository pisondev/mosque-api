package storage

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Client wraps S3-compatible object storage operations.
type Client struct {
	client *minio.Client
	bucket string
	host   string
	secure bool
}

func (c *Client) ensurePublicReadPolicy(ctx context.Context) error {
	if c.client == nil {
		return fmt.Errorf("storage client is not initialized")
	}

	policy := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Sid":"PublicReadGetObject","Effect":"Allow","Principal":"*","Action":["s3:GetObject"],"Resource":["arn:aws:s3:::%s/*"]}]}`,
		c.bucket,
	)

	if err := c.client.SetBucketPolicy(ctx, c.bucket, policy); err != nil {
		return err
	}

	return nil
}

func normalizeHost(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "https://")
	raw = strings.TrimPrefix(raw, "http://")
	return strings.TrimSuffix(raw, "/")
}

func isSecureHost(raw string) bool {
	raw = strings.TrimSpace(strings.ToLower(raw))
	return !strings.HasPrefix(raw, "http://")
}

// New creates a storage client from env vars.
// HOST may be either "is3.cloudhost.id" or "https://is3.cloudhost.id".
func New() *Client {
	rawHost := os.Getenv("HOST")
	host := normalizeHost(rawHost)
	secure := isSecureHost(rawHost)

	accessKey := os.Getenv("ACCESS_KEY_ID")
	secretKey := os.Getenv("SECRET_ACCESS_KEY")
	bucket := os.Getenv("BUCKET")

	mc, err := minio.New(host, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: secure,
		Region: "auto",
	})
	if err != nil {
		// Return client with nil SDK client; call sites will return explicit errors.
		return &Client{bucket: bucket, host: host, secure: secure}
	}

	return &Client{client: mc, bucket: bucket, host: host, secure: secure}
}

func (c *Client) publicURL(key string) string {
	scheme := "https"
	if !c.secure {
		scheme = "http"
	}
	return fmt.Sprintf("%s://%s/%s/%s", scheme, c.host, c.bucket, key)
}

// UploadStream uploads object content and returns public URL.
func (c *Client) UploadStream(ctx context.Context, key string, body io.Reader, contentType string) (string, error) {
	if c.client == nil {
		return "", fmt.Errorf("storage client is not initialized")
	}

	if err := c.ensurePublicReadPolicy(ctx); err != nil {
		return "", fmt.Errorf("failed to set bucket public-read policy: %w", err)
	}

	raw, err := io.ReadAll(body)
	if err != nil {
		return "", fmt.Errorf("failed to read upload body: %w", err)
	}

	_, err = c.client.PutObject(ctx, c.bucket, key, bytes.NewReader(raw), int64(len(raw)), minio.PutObjectOptions{
		ContentType:          contentType,
		DisableMultipart:     true,
		DisableContentSha256: true,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	return c.publicURL(key), nil
}

// UploadDataURL uploads a base64 data URL and returns public URL.
func (c *Client) UploadDataURL(ctx context.Context, key string, dataURL string) (string, error) {
	if !strings.HasPrefix(dataURL, "data:") {
		return "", fmt.Errorf("not a data URL")
	}

	semi := strings.Index(dataURL, ";base64,")
	if semi == -1 {
		return "", fmt.Errorf("invalid data URL format")
	}

	mime := dataURL[5:semi]
	b64data := dataURL[semi+8:]

	raw, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	return c.UploadStream(ctx, key, bytes.NewReader(raw), mime)
}

// Delete removes object at key.
func (c *Client) Delete(ctx context.Context, key string) error {
	if c.client == nil {
		return fmt.Errorf("storage client is not initialized")
	}
	return c.client.RemoveObject(ctx, c.bucket, key, minio.RemoveObjectOptions{})
}

// GetBucketSizeBytes sums size for all objects under prefix.
func (c *Client) GetBucketSizeBytes(ctx context.Context, prefix string) (int64, error) {
	if c.client == nil {
		return 0, fmt.Errorf("storage client is not initialized")
	}

	var total int64
	for obj := range c.client.ListObjects(ctx, c.bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	}) {
		if obj.Err != nil {
			return 0, fmt.Errorf("failed to list objects: %w", obj.Err)
		}
		total += obj.Size
	}

	return total, nil
}

// GetObjectSizeBytes returns object size for a key. Returns (0, nil) if not found.
func (c *Client) GetObjectSizeBytes(ctx context.Context, key string) (int64, error) {
	if c.client == nil {
		return 0, fmt.Errorf("storage client is not initialized")
	}

	info, err := c.client.StatObject(ctx, c.bucket, key, minio.StatObjectOptions{})
	if err != nil {
		// Treat missing object as size 0
		errMsg := strings.ToLower(err.Error())
		if strings.Contains(errMsg, "not found") || strings.Contains(errMsg, "no such key") || strings.Contains(errMsg, "404") {
			return 0, nil
		}
		return 0, err
	}

	return info.Size, nil
}

// KeyFromURL extracts object key from public URL.
func KeyFromURL(publicURL string) string {
	host := normalizeHost(os.Getenv("HOST"))
	bucket := os.Getenv("BUCKET")
	if host == "" || bucket == "" || publicURL == "" {
		return ""
	}

	httpsPrefix := fmt.Sprintf("https://%s/%s/", host, bucket)
	httpPrefix := fmt.Sprintf("http://%s/%s/", host, bucket)

	if strings.HasPrefix(publicURL, httpsPrefix) {
		return strings.TrimPrefix(publicURL, httpsPrefix)
	}
	if strings.HasPrefix(publicURL, httpPrefix) {
		return strings.TrimPrefix(publicURL, httpPrefix)
	}

	return ""
}

// TenantFolder returns base prefix for tenant subdomain.
func TenantFolder(subdomain string) string {
	return subdomain + "/"
}

// HeaderImageKey returns key path for profile header image.
func HeaderImageKey(subdomain, filename string) string {
	return path.Join(subdomain, "header", filename)
}

// ManagementPhotoKey returns key path for management profile image.
func ManagementPhotoKey(subdomain, filename string) string {
	return path.Join(subdomain, "management", filename)
}

// EventPosterKey returns key path for event poster image.
func EventPosterKey(subdomain, filename string) string {
	return path.Join(subdomain, "event", filename)
}

// QrisImageKey returns key path for static account QRIS image.
func QrisImageKey(subdomain, filename string) string {
	return path.Join(subdomain, "qris", filename)
}

