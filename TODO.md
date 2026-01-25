1. 在 `internal/data/kafka.go` 中补充 Kafka Producer 幂等与重试配置，目前仅设置了 `acks=all`。

2. retry和DLQ测试还没通过。

	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
    	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"