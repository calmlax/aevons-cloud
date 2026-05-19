package handler

import (
	"context"
	"fmt"
	"net"
	"runtime"
	"strings"
	"time"

	authstore "github.com/calmlax/aevons-framework/auth/store"
	"github.com/calmlax/aevons-framework/config"
	"github.com/calmlax/aevons-framework/consts"
	apperr "github.com/calmlax/aevons-framework/errors"
	"github.com/calmlax/aevons-framework/response"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"gorm.io/gorm"
)

type MonitorHandler struct {
	db          *gorm.DB
	redisClient *redis.Client
	cfg         *config.Config
}

type CacheEntry struct {
	Key      string `json:"key"`
	Type     string `json:"type"`
	TTL      int64  `json:"ttl"`
	Size     int64  `json:"size"`
	ExpireAt string `json:"expireAt"`
}

type OnlineUser struct {
	Token        string `json:"token"`
	Username     string `json:"username"`
	Nickname     string `json:"nickname"`
	ClientId     string `json:"clientId"`
	RefreshToken string `json:"refreshToken"`
	LoginAt      string `json:"loginAt"`
	ExpireAt     string `json:"expireAt"`
	TTL          int64  `json:"ttl"`
}

func NewMonitorHandler(db *gorm.DB, redisClient *redis.Client, cfg *config.Config) *MonitorHandler {
	return &MonitorHandler{
		db:          db,
		redisClient: redisClient,
		cfg:         cfg,
	}
}

func (h *MonitorHandler) GetServerInfo(c *gin.Context) {
	ctx := context.Background()
	response.Success(c, gin.H{
		"sys":      h.sysInfo(),
		"cpu":      h.cpuInfo(),
		"mem":      h.memInfo(),
		"disk":     h.diskInfo(),
		"go":       h.goInfo(),
		"redis":    h.redisInfo(ctx),
		"db":       h.dbInfo(),
		"rocketmq": h.rocketmqInfo(),
	})
}

func (h *MonitorHandler) List(c *gin.Context) {
	ctx := context.Background()

	groups := []struct {
		Name   string
		Prefix string
	}{
		{"访问令牌", consts.RedisKeyAccessToken},
		{"刷新令牌", consts.RedisKeyRefreshToken},
		{"会话信息", consts.RedisKeySession},
		{"会话访问令牌索引", consts.RedisKeySessionAccess},
		{"会话刷新令牌索引", consts.RedisKeySessionRefresh},
		{"授权码", consts.RedisKeyAuthCode},
		{"邮箱验证码", consts.RedisKeyEmailCode},
		{"用户会话", consts.RedisKeyUserSessions},
		{"OAuth2 State", consts.RedisKeyOAuthState},
		{"RSA 私钥", consts.RedisKeyRSAPrivateKey},
		{"系统配置缓存", consts.ConfCacheKeyPrefix},
		{"字典数据缓存", consts.DictCacheKeyPrefix},
		{"网关限流", "gateway-service:rate-limit:"},
		{"网关OAuth客户端规则", "gateway:oauth-client-rules:"},
	}

	type GroupInfo struct {
		Name   string       `json:"name"`
		Prefix string       `json:"prefix"`
		Count  int          `json:"count"`
		Keys   []CacheEntry `json:"keys"`
	}

	var result []GroupInfo
	knownKeys := map[string]bool{}

	for _, g := range groups {
		var entries []CacheEntry
		iter := h.redisClient.Scan(ctx, 0, g.Prefix+"*", 500).Iterator()
		for iter.Next(ctx) {
			key := iter.Val()
			knownKeys[key] = true
			entries = append(entries, h.buildEntry(ctx, key))
		}
		if entries == nil {
			entries = []CacheEntry{}
		}
		result = append(result, GroupInfo{
			Name:   g.Name,
			Prefix: g.Prefix,
			Count:  len(entries),
			Keys:   entries,
		})
	}

	var others []CacheEntry
	iter := h.redisClient.Scan(ctx, 0, "*", 500).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		if !knownKeys[key] {
			others = append(others, h.buildEntry(ctx, key))
		}
	}
	if others == nil {
		others = []CacheEntry{}
	}
	result = append(result, GroupInfo{
		Name:   "其他",
		Prefix: "",
		Count:  len(others),
		Keys:   others,
	})

	response.Success(c, result)
}

func (h *MonitorHandler) Detail(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		response.Success(c, nil)
		return
	}
	ctx := context.Background()

	entry := h.buildEntry(ctx, key)

	var value string
	if entry.Type == "string" {
		value, _ = h.redisClient.Get(ctx, key).Result()
	}

	response.Success(c, gin.H{
		"key":      entry.Key,
		"type":     entry.Type,
		"ttl":      entry.TTL,
		"expireAt": entry.ExpireAt,
		"size":     entry.Size,
		"value":    value,
	})
}

func (h *MonitorHandler) Delete(c *gin.Context) {
	var body struct {
		Keys []string `json:"keys" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return
	}
	ctx := context.Background()
	if err := h.redisClient.Del(ctx, body.Keys...).Err(); err != nil {
		response.FailServerError(c, "err.sys.server_error", nil)
		return
	}
	response.Success(c, nil)
}

func (h *MonitorHandler) DeleteByPrefix(c *gin.Context) {
	prefix := c.Query("prefix")
	if prefix == "" {
		response.FailBy(c, apperr.ErrInvalidQuery)
		return
	}
	ctx := context.Background()
	var keys []string
	iter := h.redisClient.Scan(ctx, 0, prefix+"*", 500).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if len(keys) > 0 {
		if err := h.redisClient.Del(ctx, keys...).Err(); err != nil {
			response.FailServerError(c, "err.sys.server_error", nil)
			return
		}
	}
	response.Success(c, gin.H{"deleted": len(keys)})
}

func (h *MonitorHandler) ListOnline(c *gin.Context) {
	ctx := context.Background()
	prefix := consts.RedisKeyAccessToken
	store := authstore.NewRedisTokenStore(h.redisClient)

	var users []OnlineUser
	iter := h.redisClient.Scan(ctx, 0, prefix+"*", 500).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		token := strings.TrimPrefix(key, prefix)

		user, err := store.GetLoginUser(ctx, token)
		if err != nil {
			continue
		}
		if user.UserId == 0 {
			continue
		}
		sessionID, _ := store.GetSessionIDByAccessToken(ctx, token)
		refreshToken := ""
		if sessionID != "" {
			refreshToken, _ = store.GetRefreshTokenBySessionID(ctx, sessionID)
		}

		ttl, _ := h.redisClient.TTL(ctx, key).Result()
		ttlSec := int64(ttl.Seconds())

		var expireAt string
		if ttlSec > 0 {
			expireAt = time.Now().Add(ttl).Format("2006-01-02 15:04:05")
		}

		users = append(users, OnlineUser{
			Token:        token,
			Username:     user.Username,
			Nickname:     user.Nickname,
			ClientId:     user.ClientId,
			RefreshToken: refreshToken,
			ExpireAt:     expireAt,
			TTL:          ttlSec,
		})
	}

	if users == nil {
		users = []OnlineUser{}
	}
	response.Success(c, users)
}

func (h *MonitorHandler) ForceLogout(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		response.FailBy(c, apperr.ErrInvalidId)
		return
	}
	ctx := context.Background()
	store := authstore.NewRedisTokenStore(h.redisClient)

	if sessionID, err := store.GetSessionIDByAccessToken(ctx, token); err == nil && sessionID != "" {
		if refreshToken, e := store.GetRefreshTokenBySessionID(ctx, sessionID); e == nil && refreshToken != "" {
			_ = h.redisClient.Del(ctx, consts.RedisKeyRefreshToken+refreshToken).Err()
		}
		_ = h.redisClient.Del(ctx, consts.RedisKeySessionRefresh+sessionID).Err()
		_ = h.redisClient.Del(ctx, consts.RedisKeySessionAccess+sessionID).Err()
		_ = h.redisClient.Del(ctx, consts.RedisKeySession+sessionID).Err()
	}
	_ = h.redisClient.Del(ctx, consts.RedisKeyAccessToken+token).Err()
	response.Success(c, nil)
}

func (h *MonitorHandler) sysInfo() gin.H {
	info, err := host.Info()
	if err != nil {
		return gin.H{"error": err.Error()}
	}
	return gin.H{
		"hostname":        info.Hostname,
		"os":              info.OS,
		"platform":        info.Platform,
		"platformVersion": info.PlatformVersion,
		"kernelVersion":   info.KernelVersion,
		"kernelArch":      info.KernelArch,
		"uptime":          formatDuration(time.Duration(info.Uptime) * time.Second),
		"bootTime":        time.Unix(int64(info.BootTime), 0),
	}
}

func (h *MonitorHandler) cpuInfo() gin.H {
	percents, err := cpu.Percent(200*time.Millisecond, false)
	usedPercent := 0.0
	if err == nil && len(percents) > 0 {
		usedPercent = percents[0]
	}
	infos, _ := cpu.Info()
	modelName := ""
	physicalCores := 0
	if len(infos) > 0 {
		modelName = infos[0].ModelName
		physicalCores = int(infos[0].Cores)
	}
	logicalCores, _ := cpu.Counts(true)
	return gin.H{
		"modelName":     modelName,
		"physicalCores": physicalCores,
		"logicalCores":  logicalCores,
		"usedPercent":   fmt.Sprintf("%.1f", usedPercent),
	}
}

func (h *MonitorHandler) memInfo() gin.H {
	v, err := mem.VirtualMemory()
	if err != nil {
		return gin.H{"error": err.Error()}
	}
	return gin.H{
		"total":       v.Total,
		"used":        v.Used,
		"free":        v.Free,
		"available":   v.Available,
		"usedPercent": fmt.Sprintf("%.1f", v.UsedPercent),
		"buffers":     v.Buffers,
		"cached":      v.Cached,
	}
}

func (h *MonitorHandler) diskInfo() gin.H {
	parts, err := disk.Partitions(false)
	if err != nil {
		return gin.H{"error": err.Error()}
	}
	skipFsTypes := map[string]bool{
		"tmpfs": true, "devtmpfs": true, "squashfs": true,
		"overlay": true, "proc": true, "sysfs": true,
		"cgroup": true, "cgroup2": true, "pstore": true,
		"debugfs": true, "tracefs": true, "securityfs": true,
		"hugetlbfs": true, "mqueue": true, "fusectl": true,
		"bpf": true, "ramfs": true,
	}
	var disks []gin.H
	for _, p := range parts {
		if skipFsTypes[p.Fstype] {
			continue
		}
		skip := false
		for _, prefix := range []string{"/snap/", "/sys", "/proc", "/dev", "/run"} {
			if len(p.Mountpoint) >= len(prefix) && p.Mountpoint[:len(prefix)] == prefix {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		u, err := disk.Usage(p.Mountpoint)
		if err != nil || u.Total == 0 {
			continue
		}
		disks = append(disks, gin.H{
			"path":        u.Path,
			"fstype":      u.Fstype,
			"total":       u.Total,
			"used":        u.Used,
			"free":        u.Free,
			"usedPercent": fmt.Sprintf("%.1f", u.UsedPercent),
		})
	}
	return gin.H{"list": disks}
}

func (h *MonitorHandler) goInfo() gin.H {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return gin.H{
		"version":      runtime.Version(),
		"os":           runtime.GOOS,
		"arch":         runtime.GOARCH,
		"numCPU":       runtime.NumCPU(),
		"numGoroutine": runtime.NumGoroutine(),
		"heapAlloc":    m.HeapAlloc,
		"heapSys":      m.HeapSys,
		"heapInuse":    m.HeapInuse,
		"stackInuse":   m.StackInuse,
		"totalAlloc":   m.TotalAlloc,
		"sys":          m.Sys,
		"gcNum":        m.NumGC,
		"gcPauseTotal": m.PauseTotalNs,
	}
}

func (h *MonitorHandler) redisInfo(ctx context.Context) gin.H {
	if h.redisClient == nil {
		return gin.H{"error": "redis not connected"}
	}
	info, err := h.redisClient.Info(ctx, "server", "clients", "memory", "stats").Result()
	if err != nil {
		return gin.H{"error": err.Error()}
	}
	parsed := parseRedisInfo(info)

	var tokenCount int64
	iter := h.redisClient.Scan(ctx, 0, consts.RedisKeyAccessToken+"*", 0).Iterator()
	for iter.Next(ctx) {
		tokenCount++
	}
	dbSize, _ := h.redisClient.DBSize(ctx).Result()
	maxMemory := parsed["maxmemory_human"]
	if maxMemory == "0B" || maxMemory == "" {
		maxMemory = "无限制"
	}
	return gin.H{
		"version":          parsed["redis_version"],
		"mode":             parsed["redis_mode"],
		"uptimeSeconds":    parsed["uptime_in_seconds"],
		"connectedClients": parsed["connected_clients"],
		"usedMemory":       parsed["used_memory_human"],
		"usedMemoryPeak":   parsed["used_memory_peak_human"],
		"maxMemory":        maxMemory,
		"memFragRatio":     parsed["mem_fragmentation_ratio"],
		"totalCommands":    parsed["total_commands_processed"],
		"totalConnections": parsed["total_connections_received"],
		"keyspaceHits":     parsed["keyspace_hits"],
		"keyspaceMisses":   parsed["keyspace_misses"],
		"dbSize":           dbSize,
		"onlineTokens":     tokenCount,
	}
}

func (h *MonitorHandler) dbInfo() gin.H {
	sqlDB, err := h.db.DB()
	if err != nil {
		return gin.H{"error": err.Error()}
	}
	var version string
	h.db.Raw("SELECT VERSION()").Scan(&version)
	stats := sqlDB.Stats()
	var dbSizeMB string
	h.db.Raw("SELECT ROUND(SUM(data_length + index_length) / 1024 / 1024, 2) FROM information_schema.tables").Scan(&dbSizeMB)
	var processCount int
	h.db.Raw("SELECT COUNT(*) FROM information_schema.processlist").Scan(&processCount)
	var uptimeSeconds int64
	h.db.Raw("SHOW GLOBAL STATUS LIKE 'Uptime'").Row().Scan(nil, &uptimeSeconds)
	return gin.H{
		"version":      version,
		"maxOpenConns": stats.MaxOpenConnections,
		"openConns":    stats.OpenConnections,
		"inUseConns":   stats.InUse,
		"idleConns":    stats.Idle,
		"waitCount":    stats.WaitCount,
		"dbSizeMB":     dbSizeMB,
		"processCount": processCount,
		"uptime":       formatDuration(time.Duration(uptimeSeconds) * time.Second),
	}
}

func (h *MonitorHandler) rocketmqInfo() gin.H {
	if h.cfg == nil || !h.cfg.RocketMQ.Enabled {
		return gin.H{"enabled": false}
	}

	cfg := &h.cfg.RocketMQ
	type nsStatus struct {
		Addr      string `json:"addr"`
		Reachable bool   `json:"reachable"`
	}

	endpoints := cfg.NameServerEndpoints()
	var nameServers []nsStatus
	allReachable := true
	for _, addr := range endpoints {
		conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
		reachable := err == nil
		if reachable {
			conn.Close()
		} else {
			allReachable = false
		}
		nameServers = append(nameServers, nsStatus{Addr: addr, Reachable: reachable})
	}

	status := "online"
	if !allReachable {
		status = "degraded"
	}
	if len(nameServers) == 0 || !nameServers[0].Reachable {
		status = "offline"
	}

	return gin.H{
		"enabled":     true,
		"status":      status,
		"nameServers": nameServers,
		"producer": gin.H{
			"group":       cfg.Producer.Group,
			"retryTimes":  cfg.Producer.RetryTimes,
			"sendTimeout": cfg.Producer.SendTimeout,
		},
		"consumer": gin.H{
			"group":       cfg.Consumer.Group,
			"concurrency": cfg.Consumer.Concurrency,
		},
	}
}

func (h *MonitorHandler) buildEntry(ctx context.Context, key string) CacheEntry {
	keyType, _ := h.redisClient.Type(ctx, key).Result()
	ttlDur, _ := h.redisClient.TTL(ctx, key).Result()
	ttlSec := int64(ttlDur.Seconds())

	var expireAt string
	if ttlSec > 0 {
		expireAt = time.Now().Add(ttlDur).Format("2006-01-02 15:04:05")
	}

	var size int64
	if keyType == "string" {
		size, _ = h.redisClient.StrLen(ctx, key).Result()
	}

	return CacheEntry{
		Key:      key,
		Type:     keyType,
		TTL:      ttlSec,
		Size:     size,
		ExpireAt: expireAt,
	}
}

func parseRedisInfo(info string) map[string]string {
	result := make(map[string]string)
	for _, line := range splitLines(info) {
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		for i, ch := range line {
			if ch == ':' {
				result[line[:i]] = line[i+1:]
				break
			}
		}
	}
	return result
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			line := s[start:i]
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			lines = append(lines, line)
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return d.Truncate(time.Second).String()
	}
	if d < 24*time.Hour {
		return d.Truncate(time.Minute).String()
	}
	days := int(d / (24 * time.Hour))
	remain := d % (24 * time.Hour)
	hours := int(remain / time.Hour)
	minutes := int((remain % time.Hour) / time.Minute)
	return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
}
