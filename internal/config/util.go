package config

import (
	"context"
	"fmt"

	"github.com/ayushanand18/as-http3lib/internal/constants"
)

func GetListeningAddress(ctx context.Context) string {
	ipAddress := GetString(ctx, "service.address.ip", constants.DEFAULT_SERVER_IP_ADDRESS)
	port := GetString(ctx, "service.address.port", constants.DEFAULT_SERVER_PORT)

	return fmt.Sprintf("%s:%s", ipAddress, port)
}
