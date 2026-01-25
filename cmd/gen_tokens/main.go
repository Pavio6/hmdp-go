package main

import (
	"context"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"hmdp-backend/internal/utils"
)

// This helper reads a CSV of users and writes login tokens to Redis + an output CSV.
// Input CSV requirements:
//   - Contains at least a phone column. Header is detected if present (field name "phone" case-insensitive).
//   - Optionally contains id/nickName/icon columns; if missing we generate reasonable defaults.
//
// Output CSV contains: token,phone
//
// Usage:
//
//	go run cmd/gen_tokens/main.go -in tb_user.csv -out tokens.csv -redis 127.0.0.1:6379
func main() {
	in := flag.String("in", "hmdp_tb_user.csv", "input CSV file (must contain phone column)")
	out := flag.String("out", "tokens.csv", "output CSV file")
	redisAddr := flag.String("redis", "127.0.0.1:6379", "redis address")
	redisDB := flag.Int("db", 0, "redis db index")
	flag.Parse()

	if *in == "" {
		log.Fatal("input file is required, use -in")
	}

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: *redisAddr,
		DB:   *redisDB,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("ping redis: %v", err)
	}

	f, err := os.Open(*in)
	if err != nil {
		log.Fatalf("open input: %v", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	header, err := reader.Read()
	if err != nil {
		log.Fatalf("read header/data: %v", err)
	}

	hasHeader := containsHeader(header)
	var (
		idIdx    = -1
		phoneIdx = -1
		nickIdx  = -1
		iconIdx  = -1
		firstRow []string
	)

	if hasHeader {
		for i, h := range header {
			switch strings.ToLower(strings.TrimSpace(h)) {
			case "id":
				idIdx = i
			case "phone":
				phoneIdx = i
			case "nickname", "nick", "nick_name":
				nickIdx = i
			case "icon":
				iconIdx = i
			}
		}
	} else {
		// Treat first read row as data
		firstRow = header
		// Assume schema: [phone] or [id,phone]
		if len(firstRow) >= 2 {
			idIdx = 0
			phoneIdx = 1
		} else {
			phoneIdx = 0
		}
	}

	if phoneIdx == -1 {
		log.Fatal("phone column not found; ensure CSV has a 'phone' header or phone is first/second column")
	}

	outFile, err := os.Create(*out)
	if err != nil {
		log.Fatalf("create output: %v", err)
	}
	defer outFile.Close()
	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	processRow := func(row []string) error {
		if len(row) == 0 {
			return nil
		}
		if phoneIdx >= len(row) {
			return errors.New("phone column out of range")
		}
		phone := strings.TrimSpace(row[phoneIdx])
		if phone == "" {
			return nil
		}

		var idVal string
		if idIdx >= 0 && idIdx < len(row) {
			idVal = strings.TrimSpace(row[idIdx])
		}
		if idVal == "" {
			idVal = phone
		}

		nick := ""
		if nickIdx >= 0 && nickIdx < len(row) {
			nick = strings.TrimSpace(row[nickIdx])
		}
		if nick == "" {
			if len(phone) >= 4 {
				nick = utils.USER_NICK_NAME_PREFIX + phone[len(phone)-4:]
			} else {
				nick = utils.USER_NICK_NAME_PREFIX + phone
			}
		}

		icon := ""
		if iconIdx >= 0 && iconIdx < len(row) {
			icon = strings.TrimSpace(row[iconIdx])
		}

		token := uuid.NewString()
		tokenKey := utils.LOGIN_USER_KEY + token
		data := map[string]string{
			"id":       idVal,
			"nickName": nick,
			"icon":     icon,
		}

		if err := rdb.HSet(ctx, tokenKey, data).Err(); err != nil {
			return fmt.Errorf("hset: %w", err)
		}
		if err := rdb.Expire(ctx, tokenKey, time.Duration(utils.LOGIN_USER_TTL)*time.Second).Err(); err != nil {
			return fmt.Errorf("expire: %w", err)
		}

		if err := writer.Write([]string{token}); err != nil {
			return fmt.Errorf("write csv: %w", err)
		}
		return nil
	}

	count := 0
	if !hasHeader {
		if err := processRow(firstRow); err != nil {
			log.Printf("process row error: %v", err)
		} else {
			count++
		}
	}

	for {
		row, err := reader.Read()
		if err != nil {
			break
		}
		if err := processRow(row); err != nil {
			log.Printf("process row error: %v", err)
			continue
		}
		count++
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		log.Fatalf("flush output: %v", err)
	}
	log.Printf("generated %d tokens to %s and wrote to Redis at %s", count, *out, *redisAddr)
}

func containsHeader(row []string) bool {
	for _, v := range row {
		switch strings.ToLower(strings.TrimSpace(v)) {
		case "phone", "id", "nickname", "nick_name", "icon":
			return true
		}
	}
	return false
}
