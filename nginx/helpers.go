package nginx

import (
	"fmt"
	"github.com/tufanbarisyildirim/gonginx"
	"log"
	"strconv"
)

func BuildWellknownServerReturnContent(serverName string) string {
	return fmt.Sprintf("'{\"m.server\": \"synapse.%s:443\"}'", serverName)
}

func BuildWellknownClientReturnContent(serverName string) string {
	return fmt.Sprintf("'{\"m.homeserver\": {\"base_url\": \"https://synapse.%s\"}}'", serverName)
}

func nginxConfigGetDirective(base gonginx.IBlock, key string) gonginx.IDirective {
	return base.FindDirectives(key)[0]
}

func nginxConfigGetReturnClause(base gonginx.IBlock, key string) MatrixConfLocationWellKnownReturnClause {
	code, err := strconv.Atoi(nginxConfigGetDirective(base, key).GetParameters()[0])
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	return MatrixConfLocationWellKnownReturnClause{
		Code:    code,
		Content: nginxConfigGetDirective(base, key).GetParameters()[1],
	}
}

func nginxConfigGetString(base gonginx.IBlock, key string) string {
	return nginxConfigGetDirective(base, key).GetParameters()[0]
}

func nginxConfigGetInt(base gonginx.IBlock, key string) int {
	val, err := strconv.Atoi(nginxConfigGetString(base, key))
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	return val
}
