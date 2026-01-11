package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	guacamole "github.com/seknox/guacamole"
)

// Config Guacamole 客户端配置
type Config struct {
	// GuacdHost guacd 服务器地址
	GuacdHost string
	// GuacdPort guacd 服务器端口（默认 4822）
	GuacdPort int
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		GuacdHost: "localhost",
		GuacdPort: 4822,
	}
}

// Tunnel Guacamole协议隧道
// 使用 seknox/guacamole 库的 NewInstruction 来构建指令，但手动控制握手流程
// 握手后，直接使用底层连接读取原始字节，避免使用 Stream.ReadSome() 的解析器
type Tunnel struct {
	conn          net.Conn
	stream        *guacamole.Stream // 仅用于库的工具函数（如 NewInstruction），不用于 Handshake
	reader        *bufio.Reader     // 用于读取原始字节
	writer        *bufio.Writer     // 用于写入指令
	config        *Config
	handshakeDone bool // 标记握手是否完成
}

// NewTunnel 创建Guacamole隧道
func NewTunnel(config *Config) *Tunnel {
	if config == nil {
		config = DefaultConfig()
	}
	return &Tunnel{
		config: config,
	}
}

// newInstruction 创建一个 Guacamole 协议指令（自动计算长度）
// 使用 seknox/guacamole 库的实现，自动处理长度计算
func newInstruction(opcode string, args ...string) string {
	inst := guacamole.NewInstruction(opcode, args...)
	return inst.String()
}

// InternalOpcode 内部操作码前缀（以这个开头的消息不应发送到WebSocket或guacd）
// 使用 seknox/guacamole 库的常量
var InternalOpcode = guacamole.InternalOpcodeIns

// Connect 连接到guacd并建立与目标服务器的连接
// 参数：
//   - hostname: 目标服务器地址
//   - username: 用户名
//   - password: 密码
//   - port: 目标服务器端口
//   - protocol: 协议类型（如 "ssh", "rdp", "vnc"）
//   - width: 显示宽度（像素）
//   - height: 显示高度（像素）
func (t *Tunnel) Connect(hostname, username, password string, port int, protocol string, width, height int) error {
	// 连接到guacd
	guacdAddr := net.JoinHostPort(t.config.GuacdHost, fmt.Sprintf("%d", t.config.GuacdPort))

	conn, err := net.Dial("tcp", guacdAddr)
	if err != nil {
		return fmt.Errorf("连接guacd失败: %w", err)
	}

	t.conn = conn
	t.reader = bufio.NewReader(conn)
	t.writer = bufio.NewWriter(conn)

	// 创建一个 Stream 对象用于库的工具函数（但不用它的 Handshake）
	t.stream = guacamole.NewStream(conn, guacamole.SocketTimeout)

	// 参考 Java 官方库的实现，按照正确的握手顺序：
	// 1. 首先发送 select 指令（选择协议）
	// 2. 等待并解析 args/required 响应
	// 3. 发送客户端能力信息（size, audio, video, image, timezone）
	// 4. 发送 connect 指令（建立连接）

	// 步骤1：发送select指令（选择协议）
	instruction := newInstruction("select", protocol)
	if err := t.writeInstruction(instruction); err != nil {
		conn.Close()
		return fmt.Errorf("发送协议选择指令失败: %w", err)
	}

	// 步骤2：等待服务器回复 args/required 指令（参数列表）
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	argsBuffer := make([]byte, 2048)
	argsN, argsErr := t.reader.Read(argsBuffer)
	conn.SetReadDeadline(time.Time{}) // 清除超时

	if argsErr != nil {
		conn.Close()
		return fmt.Errorf("读取服务器required/args响应失败: %w", argsErr)
	}

	if argsN == 0 {
		conn.Close()
		return fmt.Errorf("服务器required/args响应为空")
	}

	argsResponse := string(argsBuffer[:argsN])

	// 解析 required/args 响应，提取参数列表和版本号
	var argsParams []string
	var protocolVersion string
	if strings.HasPrefix(argsResponse, "8.required,") || strings.HasPrefix(argsResponse, "4.required,") {
		// 解析 required 指令
		argvIndex := strings.Index(argsResponse, "3.argv,")
		if argvIndex == -1 {
			argvIndex = strings.Index(argsResponse, "4.args,")
		}

		if argvIndex != -1 {
			argsContent := strings.TrimSuffix(argsResponse[argvIndex+7:], ";")
			parts := strings.Split(argsContent, ",")
			for i, part := range parts {
				part = strings.TrimSpace(part)
				if part == "" {
					continue
				}
				if idx := strings.Index(part, "."); idx > 0 {
					paramName := part[idx+1:]
					if i == 0 && strings.HasPrefix(paramName, "VERSION_") {
						protocolVersion = paramName
						continue
					}
					argsParams = append(argsParams, paramName)
				}
			}
		}
	} else if strings.HasPrefix(argsResponse, "4.args,") {
		// 解析 args 指令
		argsContent := strings.TrimSuffix(argsResponse[7:], ";")
		parts := strings.Split(argsContent, ",")
		for i, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			if idx := strings.Index(part, "."); idx > 0 {
				paramName := part[idx+1:]
				if i == 0 && strings.HasPrefix(paramName, "VERSION_") {
					protocolVersion = paramName
				} else {
					argsParams = append(argsParams, paramName)
				}
			}
		}
	} else {
		conn.Close()
		return fmt.Errorf("无法识别服务器响应格式: %s", argsResponse)
	}

	// 如果没有解析到版本号，使用默认版本
	if protocolVersion == "" {
		protocolVersion = "VERSION_1_5_0"
	}

	if len(argsParams) == 0 {
		argsParams = []string{"hostname", "port", "username", "password"}
	}

	// 步骤3：发送客户端能力信息（size, audio, video, image, timezone）
	instruction = newInstruction("size", fmt.Sprintf("%d", width), fmt.Sprintf("%d", height), "96")
	if err := t.writeInstruction(instruction); err != nil {
		conn.Close()
		return fmt.Errorf("发送尺寸指令失败: %w", err)
	}

	instruction = newInstruction("audio")
	if err := t.writeInstruction(instruction); err != nil {
		conn.Close()
		return fmt.Errorf("发送audio指令失败: %w", err)
	}

	instruction = newInstruction("video")
	if err := t.writeInstruction(instruction); err != nil {
		conn.Close()
		return fmt.Errorf("发送video指令失败: %w", err)
	}

	instruction = newInstruction("image", "image/png", "image/jpeg")
	if err := t.writeInstruction(instruction); err != nil {
		conn.Close()
		return fmt.Errorf("发送image指令失败: %w", err)
	}

	instruction = newInstruction("timezone", "UTC")
	if err := t.writeInstruction(instruction); err != nil {
		conn.Close()
		return fmt.Errorf("发送timezone指令失败: %w", err)
	}

	// 步骤4：构建并发送connect指令
	// connect 指令的第一个参数必须是版本号，然后是参数值
	connectValues := []string{protocolVersion}
	for _, paramName := range argsParams {
		switch paramName {
		case "hostname":
			connectValues = append(connectValues, hostname)
		case "port":
			connectValues = append(connectValues, fmt.Sprintf("%d", port))
		case "username":
			connectValues = append(connectValues, username)
		case "password":
			connectValues = append(connectValues, password)
		default:
			connectValues = append(connectValues, "") // 其他参数使用空字符串
		}
	}

	instruction = newInstruction("connect", connectValues...)
	if err := t.writeInstruction(instruction); err != nil {
		conn.Close()
		return fmt.Errorf("发送connect指令失败: %w", err)
	}

	// 步骤5：等待guacd的connect响应
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	connectBuffer := make([]byte, 1024)
	connectN, connectErr := t.reader.Read(connectBuffer)
	conn.SetReadDeadline(time.Time{})

	if connectErr != nil && connectErr != io.EOF {
		conn.Close()
		return fmt.Errorf("读取guacd connect响应失败: %w", connectErr)
	}

	if connectN > 0 {
		connectResponse := string(connectBuffer[:connectN])
		if strings.Contains(connectResponse, ".error,") {
			conn.Close()
			return fmt.Errorf("guacd返回错误: %s", connectResponse)
		}
	}

	// 握手完成后，标记为完成
	t.handshakeDone = true
	return nil
}

// writeInstruction 写入指令
func (t *Tunnel) writeInstruction(instruction string) error {
	if t.writer == nil {
		return fmt.Errorf("writer未初始化")
	}
	if _, err := t.writer.WriteString(instruction); err != nil {
		return err
	}
	return t.writer.Flush()
}

// Read 读取数据（io.Reader接口）
func (t *Tunnel) Read(p []byte) (n int, err error) {
	if t.conn == nil {
		return 0, fmt.Errorf("连接未建立")
	}
	return t.conn.Read(p)
}

// ReadSome 读取一些数据（用于流式读取）
// 握手后直接使用底层连接读取原始字节
func (t *Tunnel) ReadSome() ([]byte, error) {
	if !t.handshakeDone || t.reader == nil {
		return nil, fmt.Errorf("连接未建立或握手未完成")
	}

	// 设置读取超时，避免无限阻塞（30秒超时）
	if t.conn != nil {
		if err := t.conn.SetReadDeadline(time.Now().Add(30 * time.Second)); err != nil {
			return nil, fmt.Errorf("设置读取超时失败: %w", err)
		}
		defer t.conn.SetReadDeadline(time.Time{}) // 清除超时
	}

	// 首先检查缓冲区是否有数据
	if t.reader.Buffered() > 0 {
		buffered := make([]byte, t.reader.Buffered())
		if n, err := t.reader.Read(buffered); err == nil && n > 0 {
			return buffered[:n], nil
		}
	}

	// 从连接读取数据
	buffer := make([]byte, 4096)
	n, err := t.reader.Read(buffer)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			if t.reader.Buffered() > 0 {
				buffered := make([]byte, t.reader.Buffered())
				if bufferedN, bufferedErr := t.reader.Read(buffered); bufferedErr == nil && bufferedN > 0 {
					return buffered[:bufferedN], nil
				}
			}
			return nil, nil // 超时且无缓冲数据，返回nil而不是错误
		}
		if err == io.EOF {
			return nil, io.EOF
		}
		return nil, err
	}
	if n == 0 {
		return nil, nil
	}
	return buffer[:n], nil
}

// Available 检查是否有更多数据可读
func (t *Tunnel) Available() bool {
	if !t.handshakeDone {
		if t.stream != nil {
			return t.stream.Available()
		}
		return false
	}
	if t.reader != nil {
		return t.reader.Buffered() > 0
	}
	return false
}

// Write 写入数据（io.Writer接口）
func (t *Tunnel) Write(p []byte) (n int, err error) {
	if !t.handshakeDone || t.stream == nil {
		return 0, fmt.Errorf("连接未建立或握手未完成")
	}
	if len(p) == 0 {
		return 0, nil
	}
	n, err = t.stream.Write(p)
	if err != nil {
		return n, err
	}
	t.stream.Flush()
	return n, nil
}

// Close 关闭连接
func (t *Tunnel) Close() error {
	if t.stream != nil {
		return t.stream.Close()
	}
	if t.conn != nil {
		return t.conn.Close()
	}
	return nil
}
