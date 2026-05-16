package apisixadmin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	BaseURL string
	APIKey  string
	Client  *http.Client
}

func New(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		APIKey:  apiKey,
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) PutRoute(ctx context.Context, id string, resource RouteResource) error {
	return c.putJSON(ctx, "/apisix/admin/routes/"+id, resource)
}

func (c *Client) PutUpstream(ctx context.Context, id string, resource UpstreamResource) error {
	return c.putJSON(ctx, "/apisix/admin/upstreams/"+id, resource)
}

func (c *Client) PutConsumer(ctx context.Context, username string, resource ConsumerResource) error {
	return c.putJSON(ctx, "/apisix/admin/consumers/"+username, resource)
}

func (c *Client) PutPluginConfig(ctx context.Context, id string, resource PluginConfigResource) error {
	return c.putJSON(ctx, "/apisix/admin/plugin_configs/"+id, resource)
}

func (c *Client) PutGlobalRule(ctx context.Context, id string, resource GlobalRuleResource) error {
	return c.putJSON(ctx, "/apisix/admin/global_rules/"+id, resource)
}

func (c *Client) putJSON(ctx context.Context, path string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.BaseURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		req.Header.Set("X-API-KEY", c.APIKey)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("apisix admin returned status %d for %s", resp.StatusCode, path)
	}
	return nil
}
