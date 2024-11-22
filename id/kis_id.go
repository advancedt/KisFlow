package id

import (
	"KisFlow/common"
	"github.com/google/uuid"
	"strings"
)

// KisID 获取一个中随机实例ID
// 格式为  "prefix1-[prefix2-][prefix3-]ID"
// 如：flow-1234567890
// 如：func-1234567890
// 如: conn-1234567890
// 如: func-1-1234567890

func KisID(prefix ...string) (kisID string) {
	// 生成一个不含连接符的uuid
	idStr := strings.Replace(uuid.New().String(), "-", "", -1)
	kisID = formatKisID(idStr, prefix...)

	return
}

// 添加可选前缀
func formatKisID(idStr string, prefix ...string) string {
	var kisID string
	// 将每个前缀添加到kisID
	for _, fix := range prefix {
		kisID += fix
		kisID += common.KisIdJoinChar
	}

	// 将uuid添加到kisID
	kisID += idStr

	return kisID
}
