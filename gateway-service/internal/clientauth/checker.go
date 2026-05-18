package clientauth

import (
	"strings"

	"gateway-service/internal/model"
)

type Checker struct {
	rules map[string]model.ClientRule
}

func NewChecker(configs []model.ClientRuleConfig) *Checker {
	rules := make(map[string]model.ClientRule, len(configs))
	for _, cfg := range configs {
		clientID := strings.TrimSpace(cfg.ClientID)
		if clientID == "" {
			continue
		}

		rule := model.ClientRule{
			ClientID:   clientID,
			Enabled:    cfg.Enabled,
			ExactRules: map[string]struct{}{},
		}

		for _, resource := range cfg.Resources {
			resource = strings.TrimSpace(resource)
			if resource == "" {
				continue
			}
			if resource == "ALL" {
				rule.AllowAll = true
				continue
			}
			if strings.HasSuffix(resource, "**") {
				rule.PrefixRules = append(rule.PrefixRules, normalizeRule(strings.TrimSuffix(resource, "**")))
				continue
			}
			rule.ExactRules[normalizeRule(resource)] = struct{}{}
		}

		rules[clientID] = rule
	}

	return &Checker{rules: rules}
}

func (c *Checker) Allow(clientID, method, path string) bool {
	if c == nil || len(c.rules) == 0 {
		return true
	}

	clientID = strings.TrimSpace(clientID)
	if clientID == "" {
		return false
	}

	rule, ok := c.rules[clientID]
	if !ok || !rule.Enabled {
		return false
	}
	if rule.AllowAll {
		return true
	}

	key := normalizeRule(strings.ToUpper(method) + ":" + path)
	if _, ok := rule.ExactRules[key]; ok {
		return true
	}
	if _, ok := rule.ExactRules["ANY:"+path]; ok {
		return true
	}

	for _, prefix := range rule.PrefixRules {
		if strings.HasPrefix(key, prefix) {
			return true
		}
		if strings.HasPrefix("ANY:"+path, prefix) {
			return true
		}
	}
	return false
}

func (c *Checker) Rule(clientID string) (model.ClientRule, bool) {
	if c == nil {
		return model.ClientRule{}, false
	}
	rule, ok := c.rules[strings.TrimSpace(clientID)]
	return rule, ok
}

func normalizeRule(rule string) string {
	rule = strings.TrimSpace(rule)
	if rule == "" {
		return ""
	}
	parts := strings.SplitN(rule, ":", 2)
	if len(parts) != 2 {
		return rule
	}
	return strings.ToUpper(strings.TrimSpace(parts[0])) + ":" + strings.TrimSpace(parts[1])
}
