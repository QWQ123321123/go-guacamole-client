package client

// FindLastCompleteInstruction 从数据末尾向前查找最后一个完整的 Guacamole 指令
// Guacamole 指令格式：<length1>.<opcode>,<length2>.<value1>,<length3>.<value2>,...;
// 指令以分号 ; 结尾
// 返回：最后一个完整指令的结束位置（包含分号），如果没有完整指令则返回 0
func FindLastCompleteInstruction(data []byte) int {
	if len(data) == 0 {
		return 0
	}

	// 从末尾向前查找分号，找到最后一个完整指令的结束位置
	for i := len(data) - 1; i >= 0; i-- {
		if data[i] == ';' {
			// 找到分号，返回位置+1（包含分号）
			return i + 1
		}
	}

	// 没有找到分号，说明没有完整指令
	return 0
}
