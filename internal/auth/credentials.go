package auth

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/api/analyticsadmin/v1beta"
	"google.golang.org/api/analyticsdata/v1beta"
	"google.golang.org/api/option"
)

func credentialOption() (option.ClientOption, error) {
	path := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if path != "" {
		return option.WithCredentialsFile(path), nil
	}
	// Fall back to Application Default Credentials
	return nil, nil
}

func NewAnalyticsDataService(ctx context.Context) (*analyticsdata.Service, error) {
	opt, err := credentialOption()
	if err != nil {
		return nil, fmt.Errorf("credentials: %w", err)
	}
	var opts []option.ClientOption
	if opt != nil {
		opts = append(opts, opt)
	}
	svc, err := analyticsdata.NewService(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("analytics data service: %w", err)
	}
	return svc, nil
}

func NewAnalyticsAdminService(ctx context.Context) (*analyticsadmin.Service, error) {
	opt, err := credentialOption()
	if err != nil {
		return nil, fmt.Errorf("credentials: %w", err)
	}
	var opts []option.ClientOption
	if opt != nil {
		opts = append(opts, opt)
	}
	svc, err := analyticsadmin.NewService(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("analytics admin service: %w", err)
	}
	return svc, nil
}
