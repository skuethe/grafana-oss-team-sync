package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/singleflight"
)

type GroupRepresentation struct {
	ID            string                `json:"id,omitempty"`
	Name          string                `json:"name,omitempty"`
	Description   string                `json:"description,omitempty"`
	Path          string                `json:"path,omitempty"`
	ParentID      string                `json:"parentId,omitempty"`
	SubGroupCount int64                 `json:"subGroupCount,omitempty"`
	SubGroups     []GroupRepresentation `json:"subGroups,omitempty"`
	Attributes    map[string][]string   `json:"attributes,omitempty"`
	RealmRoles    []string              `json:"realmRoles,omitempty"`
	ClientRoles   map[string][]string   `json:"clientRoles,omitempty"`
	Access        map[string]bool       `json:"access,omitempty"`
}

type UserRepresentation struct {
	ID                         string              `json:"id,omitempty"`
	Username                   string              `json:"username,omitempty"`
	FirstName                  string              `json:"firstName,omitempty"`
	LastName                   string              `json:"lastName,omitempty"`
	Email                      string              `json:"email,omitempty"`
	EmailVerified              bool                `json:"emailVerified,omitempty"`
	Attributes                 map[string][]string `json:"attributes,omitempty"`
	Enabled                    bool                `json:"enabled,omitempty"`
	Self                       string              `json:"self,omitempty"`
	Origin                     string              `json:"origin,omitempty"`
	CreatedTimestamp           int64               `json:"createdTimestamp,omitempty"`
	Totp                       bool                `json:"totp,omitempty"`
	FederationLink             string              `json:"federationLink,omitempty"`
	ServiceAccountClientID     string              `json:"serviceAccountClientId,omitempty"`
	DisableableCredentialTypes []string            `json:"disableableCredentialTypes,omitempty"`
	RequiredActions            []string            `json:"requiredActions,omitempty"`
	RealmRoles                 []string            `json:"realmRoles,omitempty"`
	ClientRoles                map[string][]string `json:"clientRoles,omitempty"`
	NotBefore                  int32               `json:"notBefore,omitempty"`
	ApplicationRoles           map[string][]string `json:"applicationRoles,omitempty"`
	Groups                     []string            `json:"groups,omitempty"`
	Access                     map[string]bool     `json:"access,omitempty"`
}

// AdminClient provides thread-safe authenticated access to Keycloak Admin API.
type AdminClient struct {
	baseURL  string
	realm    string
	username string
	password string
	clientID string

	mu           sync.RWMutex
	token        string
	refreshToken string
	tokenExpiry  time.Time

	sf         singleflight.Group
	httpClient *http.Client
}

// tokenResponse represents Keycloak’s token endpoint response.
type tokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
}

// NewAdminClient initializes a new Keycloak AdminClient.
func NewAdminClient(baseURL, realm, username, password, clientID string) *AdminClient {
	return &AdminClient{
		baseURL:    strings.TrimRight(baseURL, "/"),
		realm:      realm,
		username:   username,
		password:   password,
		clientID:   clientID,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// tokenValid checks if the current access token is still valid.
func (a *AdminClient) tokenValid() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.token != "" && time.Now().Before(a.tokenExpiry)
}

// Login performs password-grant authentication to obtain a new access token.
func (a *AdminClient) Login(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	form := url.Values{}
	form.Set("grant_type", "password")
	form.Set("client_id", "admin-cli")
	form.Set("username", a.username)
	form.Set("password", a.password)

	endpoint := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", a.baseURL, a.realm)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed: %s", string(body))
	}

	var t tokenResponse
	if err := json.Unmarshal(body, &t); err != nil {
		return fmt.Errorf("parse token response: %w", err)
	}

	a.token = t.AccessToken
	a.refreshToken = t.RefreshToken
	a.tokenExpiry = time.Now().Add(time.Duration(t.ExpiresIn-10) * time.Second) // 10s buffer
	return nil
}

// Refresh obtains a new token using the refresh_token grant.
func (a *AdminClient) Refresh(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.refreshToken == "" {
		return errors.New("no refresh token available")
	}

	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("client_id", a.clientID)
	form.Set("refresh_token", a.refreshToken)

	endpoint := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", a.baseURL, a.realm)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("refresh failed: %s", string(body))
	}

	var t tokenResponse
	if err := json.Unmarshal(body, &t); err != nil {
		return fmt.Errorf("parse refresh token: %w", err)
	}

	a.token = t.AccessToken
	a.refreshToken = t.RefreshToken
	a.tokenExpiry = time.Now().Add(time.Duration(t.ExpiresIn-10) * time.Second)
	return nil
}

// Do executes an authenticated request, handling automatic login and refresh.
func (a *AdminClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	if a.tokenValid() {
		a.mu.RLock()
		token := a.token
		a.mu.RUnlock()
		req2 := req.Clone(ctx)
		req2.Header.Set("Authorization", "Bearer "+token)
		return a.httpClient.Do(req2)
	}

	_, err, _ := a.sf.Do("refresh", func() (any, error) {
		if a.tokenValid() {
			return nil, nil
		}
		if err := a.Refresh(ctx); err != nil {
			return nil, a.Login(ctx)
		}
		return nil, nil
	})
	if err != nil {
		return nil, fmt.Errorf("auth refresh/login failed: %w", err)
	}

	// Use the updated token
	a.mu.RLock()
	token := a.token
	a.mu.RUnlock()

	req2 := req.Clone(ctx)
	req2.Header.Set("Authorization", "Bearer "+token)
	return a.httpClient.Do(req2)
}

func (a *AdminClient) GetGroupSubCategory(ctx context.Context, groupID string) ([]GroupRepresentation, error) {
	endpoint := fmt.Sprintf("%s/admin/realms/%s/groups/%s/children", a.baseURL, a.realm, groupID)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)

	resp, err := a.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get sub groups: %s", string(b))
	}

	var groups []GroupRepresentation
	if err := json.NewDecoder(resp.Body).Decode(&groups); err != nil {
		return nil, fmt.Errorf("parse sub groups response: %w", err)
	}
	return groups, nil
}

func (a *AdminClient) GetGroupMember(ctx context.Context, groupID string) ([]UserRepresentation, error) {
	endpoint := fmt.Sprintf("%s/admin/realms/%s/groups/%s/members", a.baseURL, a.realm, groupID)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)

	resp, err := a.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get groups: %s", string(b))
	}

	var user []UserRepresentation
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("parse groups response: %w", err)
	}
	return user, nil
}

func (a *AdminClient) GetGroups(ctx context.Context) ([]GroupRepresentation, error) {
	endpoint := fmt.Sprintf("%s/admin/realms/%s/groups", a.baseURL, a.realm)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)

	resp, err := a.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get groups: %s", string(b))
	}

	var groups []GroupRepresentation
	if err := json.NewDecoder(resp.Body).Decode(&groups); err != nil {
		return nil, fmt.Errorf("parse groups response: %w", err)
	}
	return groups, nil
}

func (a *AdminClient) BFSGroupTraversal(ctx context.Context) (map[string]string, error) {
	result := make(map[string]string)
	mu := sync.Mutex{}

	topGroups, err := a.GetGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("get top groups: %w", err)
	}

	type node struct {
		group GroupRepresentation
		path  string
	}
	queue := make([]node, 0, len(topGroups))
	for _, g := range topGroups {
		queue = append(queue, node{g, g.Name})
		mu.Lock()
		result[g.ID] = g.Name
		mu.Unlock()
	}

	// BFS traversal
	for len(queue) > 0 {
		next := make([]node, 0)
		eg, ctx := errgroup.WithContext(ctx)
		eg.SetLimit(8) // limit concurrency to avoid overloading Keycloak

		for _, n := range queue {
			n := n // capture loop variable
			eg.Go(func() error {
				subGroups, err := a.GetGroupSubCategory(ctx, n.group.ID)
				if err != nil {
					return err
				}
				for _, sub := range subGroups {
					subPath := n.path + "/" + sub.Name
					mu.Lock()
					result[sub.ID] = subPath
					mu.Unlock()
					next = append(next, node{sub, subPath})
				}
				return nil
			})
		}

		// Wait for this level to finish
		if err := eg.Wait(); err != nil {
			return nil, err
		}

		// Move to next level
		queue = next
	}

	return result, nil
}
