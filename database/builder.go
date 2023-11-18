package database

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/FlashpointProject/CommunityWebsite/types"
	"github.com/jackc/pgx/v5"
)

type SqlText struct {
	Text string
}

type SqlBuilder struct {
	query           string
	arguments       []interface{}
	whereConditions []string
	limit           *int64
	offset          *int64
	orderBy         *string
	orderReverse    bool
	counter         int
	finalRegex      *regexp.Regexp
	buildableRegex  *regexp.Regexp
}

func NewSqlBuilder(baseQuery string) *SqlBuilder {
	finalRegex := regexp.MustCompile(`\$\d+`)
	buildableRegex := regexp.MustCompile(`\$(\d+)`)

	return &SqlBuilder{
		query:           baseQuery,
		arguments:       make([]interface{}, 0),
		whereConditions: make([]string, 0),
		limit:           nil,
		offset:          nil,
		orderBy:         nil,
		orderReverse:    false,
		counter:         0,
		finalRegex:      finalRegex,
		buildableRegex:  buildableRegex,
	}
}

func (sb *SqlBuilder) SetBase(baseQuery string) {
	sb.query = baseQuery
}

func (sb *SqlBuilder) Where(condition string, args ...interface{}) {
	matches := sb.buildableRegex.FindAllStringSubmatch(condition, -1)
	tempCounter := 0

	// Find the largest number among the placeholders
	for _, match := range matches {
		if num, err := strconv.Atoi(match[1]); err == nil && num > tempCounter {
			tempCounter = num
		}
	}

	parsedCondition := sb.finalRegex.ReplaceAllStringFunc(condition, func(placeholder string) string {
		// Extract the number from the placeholder
		num, err := strconv.Atoi(placeholder[1:])
		if err != nil {
			// This should not happen as the regex ensures the format is correct
			return placeholder
		}
		// Increment the number by the offset and return the new placeholder
		return fmt.Sprintf("$%d", num+sb.counter)
	})
	sb.whereConditions = append(sb.whereConditions, parsedCondition)
	sb.arguments = append(sb.arguments, args...)

	sb.counter += tempCounter
}

func (sb *SqlBuilder) Limit(limit int64) {
	sb.limit = &limit
}

func (sb *SqlBuilder) Offset(offset int64) {
	sb.offset = &offset
}

func (sb *SqlBuilder) OrderBy(column string, direction string, validOptions []string) {
	valid := false
	for _, option := range validOptions {
		if column == option {
			valid = true
			break
		}
	}
	if !valid {
		return
	}
	sb.orderBy = &column
	realDir := strings.ToUpper(direction)
	if realDir != "ASC" && realDir != "DESC" {
		realDir = "ASC"
	}
	if realDir == "DESC" {
		sb.orderReverse = true
	} else {
		sb.orderReverse = false
	}
}

func (sb *SqlBuilder) Count(offset int) string {
	var builder strings.Builder
	builder.WriteString(sb.query)
	if len(sb.whereConditions) > 0 {
		builder.WriteString(" WHERE ")
		builder.WriteString(strings.Join(sb.whereConditions, " AND "))
	}

	finalQuery := builder.String()

	return sb.finalRegex.ReplaceAllStringFunc(finalQuery, func(placeholder string) string {
		// Extract the number from the placeholder
		num, err := strconv.Atoi(placeholder[1:])
		if err != nil {
			// This should not happen as the regex ensures the format is correct
			return placeholder
		}
		// Increment the number by the offset and return the new placeholder
		return fmt.Sprintf("$%d", num+offset)
	})
}

func (sb *SqlBuilder) Build(offset int) string {
	var builder strings.Builder
	tempCounter := sb.counter
	builder.WriteString(sb.query)
	if len(sb.whereConditions) > 0 {
		builder.WriteString(" WHERE ")
		builder.WriteString(strings.Join(sb.whereConditions, " AND "))
	}

	if sb.orderBy != nil {
		order := "ASC"
		if sb.orderReverse {
			order = "DESC"
		}
		builder.WriteString(fmt.Sprintf(" ORDER BY %s %s", *sb.orderBy, order))
	}

	if sb.limit != nil {
		tempCounter += 1
		builder.WriteString(fmt.Sprintf(" LIMIT $%d", tempCounter))
	}

	if sb.offset != nil {
		tempCounter += 1
		builder.WriteString(fmt.Sprintf(" OFFSET $%d", tempCounter))
	}

	finalQuery := builder.String()

	return sb.finalRegex.ReplaceAllStringFunc(finalQuery, func(placeholder string) string {
		// Extract the number from the placeholder
		num, err := strconv.Atoi(placeholder[1:])
		if err != nil {
			// This should not happen as the regex ensures the format is correct
			return placeholder
		}
		// Increment the number by the offset and return the new placeholder
		return fmt.Sprintf("$%d", num+offset)
	})
}

func (sb *SqlBuilder) Arguments() []interface{} {
	allArgs := make([]interface{}, 0)
	allArgs = append(allArgs, sb.arguments...)
	if sb.limit != nil {
		allArgs = append(allArgs, *sb.limit)
	}
	if sb.offset != nil {
		allArgs = append(allArgs, *sb.offset)
	}
	return allArgs
}

func (sb *SqlBuilder) ArgumentsCount() []interface{} {
	allArgs := make([]interface{}, 0)
	allArgs = append(allArgs, sb.arguments...)
	return allArgs
}

func ReadGame(row pgx.Row) (*types.CachedGame, error) {
	var id string
	var title string
	var series string
	var developer string
	var publisher string
	var releaseDate string
	var playMode []string
	var language []string
	var originalDescription string
	var platformName string
	var updatedAt time.Time
	err := row.Scan(&id, &title, &series, &developer, &publisher, &releaseDate, &playMode, &language, &originalDescription, &platformName, &updatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		} else {
			fmt.Println("Not pgx.ErrNowRows", err.Error())
		}
		return nil, err
	}

	return &types.CachedGame{
		ID:                  id,
		Title:               title,
		Series:              series,
		Developer:           developer,
		Publisher:           publisher,
		ReleaseDate:         releaseDate,
		PlayMode:            playMode,
		Language:            language,
		OriginalDescription: originalDescription,
		Platform:            platformName,
		UpdatedAt:           updatedAt,
	}, nil
}
