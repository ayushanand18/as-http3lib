package utils

import (
	"context"
	"fmt"

	"github.com/ayushanand18/as-http3lib/internal/config"
	"github.com/ayushanand18/as-http3lib/internal/constants"
)

func GetListeningAddress(ctx context.Context) string {
	ipAddress := config.GetString(ctx, "service.address.ip", constants.DEFAULT_SERVER_IP_ADDRESS)
	port := config.GetInt(ctx, "service.address.port", constants.DEFAULT_SERVER_PORT)

	return fmt.Sprintf("%s:%d", ipAddress, port)
}
