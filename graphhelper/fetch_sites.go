package graphhelper

import (
	"context"
	"fmt"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/sites"
)

func (g *GraphHelper) GetSiteByName(ctx context.Context, siteName string) ([]models.Siteable, error) {
	quotedSiteName := fmt.Sprintf("\"%s\"", siteName)
	query := &sites.SitesRequestBuilderGetQueryParameters{
		Search: &quotedSiteName,
	}
	config := &sites.SitesRequestBuilderGetRequestConfiguration{
		QueryParameters: query,
	}
	sitesResponse, err := g.appClient.Sites().Get(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("site search error: %w", err)
	}
	siteList := sitesResponse.GetValue()
	if len(siteList) == 0 {
		return nil, fmt.Errorf("no sites found matching the search query '%s'", siteName)
	}
	return siteList, nil
}
