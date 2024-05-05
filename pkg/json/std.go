//go:build stdjson

package json

import "encoding/json"

/*
Marshal 返回 v 的 JSON 编码。

Marshal 递归地遍历值 v。
如果遇到的值实现了 Marshaler 接口并且不是 nil 指针，Marshal 会调用其 MarshalJSON 方法来生成 JSON。
如果不存在 MarshalJSON 方法，但该值实现了encoding.TextMarshaler，则 Marshal 调用其 MarshalText 方法并将结果编码为 JSON 字符串。
nil 指针异常并不是严格必要的，但模仿了 UnmarshalJSON 行为中类似的必要异常。

否则，Marshal 使用以下与类型相关的默认编码：

布尔值编码为 JSON 布尔值。

浮点、整数和数字值编码为 JSON 数字。
NaN 和 +/-Inf 值将返回 [UnsupportedValueError]。

字符串值编码为强制转换为有效 UTF-8 的 JSON 字符串，用 Unicode 替换符文替换无效字节。
为了使 JSON 能够安全地嵌入 HTML <script> 标记中，该字符串使用 HTMLEscape 进行编码，它将替换“<”、“">”、“&”，U+2028 和 U+2029 转义为“\” u003c”、“\u003e”、“\u0026”、“\u2028”和“\u2029”。
使用编码器时，可以通过调用 SetEscapeHTML(false) 来禁用此替换。

数组和切片值编码为 JSON 数组，但 []byte 编码为 Base64 编码字符串，而 nil 切片编码为空 JSON 值。

结构体值编码为 JSON 对象。
每个导出的结构体字段都成为对象的成员，使用字段名称作为对象键，除非由于下面给出的原因之一而省略该字段。

每个结构体字段的编码可以通过存储在结构体字段标记中的“json”键下的格式字符串进行自定义。
格式字符串给出字段的名称，后面可能是逗号分隔的选项列表。
该名称可以为空，以便指定选项而不覆盖默认字段名称。

“omitempty”选项指定如果字段具有空值（定义为 false、0、nil 指针、nil 接口值以及任何空数组、切片、映射或字符串），则应从编码中省略该字段。

作为一种特殊情况，如果字段标记为“-”，则该字段始终被省略。
请注意，名称为“-”的字段仍然可以使用标签“-,”生成。

结构体字段标签示例及其含义：

	// 字段在 JSON 中显示为键“myName”。
	Field int `json:"myName"`

	// 字段在 JSON 中显示为键“myName”并且
	// 如果该字段的值为空，则从对象中省略该字段，
	// 如上所定义。
	Field int `json:"myName,omitempty"`

	// Field 在 JSON 中显示为键“Field”（默认值），但是
	// 如果该字段为空，则跳过该字段。
	// 注意前面的逗号。
	Field int `json:",omitempty"`

	// 该包忽略该字段。
	Field int `json:"-"`

	// 字段在 JSON 中显示为键“-”。
	Field int `json:"-,"`

“string”选项表示字段以 JSON 形式存储在 JSON 编码的字符串中。
它仅适用于字符串、浮点、整数或布尔类型的字段。
与 JavaScript 程序通信时有时会使用这种额外的编码级别：

	Int64String int64 `json:",string"`

如果键名称是仅由 Unicode 字母、数字和 ASCII 标点符号（引号、反斜杠和逗号除外）组成的非空字符串，则将使用该键名称。

匿名结构体字段通常被编组，就好像它们的内部导出字段是外部结构体中的字段一样，遵循下一段中描述的修改后的常见 Go 可见性规则。
具有在其 JSON 标记中给出的名称的匿名结构字段被视为具有该名称，而不是匿名的。
接口类型的匿名结构体字段被视为与将该类型作为其名称相同，而不是匿名。

在决定封送或解封哪个字段时，针对 JSON 修改了结构体字段的 Go 可见性规则。
如果同一级别有多个字段，并且该级别是嵌套最少的（因此将是通常 Go 规则选择的嵌套级别），则适用以下额外规则：

1) 在这些字段中，如果有 JSON 标记的字段，则仅考虑标记字段，即使存在多个可能发生冲突的未标记字段。

2) 如果恰好有一个字段（根据第一条规则标记或未标记），则选择该字段。

3）否则有多个字段，全部忽略； 没有错误发生。

处理匿名结构体字段是 Go 1.1 中的新功能。
在 Go 1.1 之前，匿名结构体字段被忽略。
要强制忽略当前版本和早期版本中的匿名结构字段，请为该字段指定 JSON 标记“-”。

映射值编码为 JSON 对象。
映射的键类型必须是字符串、整数类型或实现encoding.TextMarshaler。
通过应用以下规则对映射键进行排序并用作 JSON 对象键，并遵守上面针对字符串值描述的 UTF-8 强制转换：

  - 任何字符串类型的键都可以直接使用
  - coding.TextMarshalers 被编组
  - 整数键转换为字符串

指针值编码为指向的值。
nil 指针编码为 null JSON 值。

接口值编码为接口中包含的值。
nil 接口值编码为 null JSON 值。

通道、复数和函数值无法以 JSON 进行编码。
尝试对此类值进行编码会导致 Marshal 返回 UnsupportedTypeError。

JSON 无法表示循环数据结构，Marshal 不处理它们。
将循环结构传递给 Marshal 将导致错误。
*/
var Marshal = json.Marshal

/*
Unmarshal 解析 JSON 编码的数据并将结果存储在 v 指向的值中。如果 v 为 nil 或不是指针，Unmarshal 将返回 InvalidUnmarshalError。

Unmarshal 使用与 Marshal 使用的编码相反的编码，根据需要分配映射、切片和指针，并具有以下附加规则：

要将 JSON 解组为指针，Unmarshal 首先处理 JSON 为 JSON 文字 null 的情况。
在这种情况下，Unmarshal 将指针设置为 nil。
否则，Unmarshal 将 JSON 解组为指针指向的值。
如果指针为 nil，Unmarshal 会为其分配一个新值来指向。

要将 JSON 解组为实现 Unmarshaler 接口的值，Unmarshal 会调用该值的 UnmarshalJSON 方法，包括当输入为 JSON null 时。
否则，如果该值实现 encoding.TextUnmarshaler 并且输入是 JSON 带引号的字符串，则 Unmarshal 会使用该字符串的不带引号形式调用该值的 UnmarshalText 方法。

要将 JSON 解组为结构，Unmarshal 会将传入对象键与 Marshal 使用的键（结构体字段名称或其标记）进行匹配，首选完全匹配，但也接受不区分大小写的匹配。
默认情况下，没有相应结构字段的对象键将被忽略（请参阅 Decoder.DisallowUnknownFields 了解替代方案）。

要将 JSON 解组为接口值，Unmarshal 将以下内容之一存储在接口值中：

	bool, for JSON booleans
	float64, for JSON numbers
	string, for JSON strings
	[]interface{}, for JSON arrays
	map[string]interface{}, for JSON objects
	nil for JSON null

要将 JSON 数组解组为切片，Unmarshal 将切片长度重置为零，然后将每个元素附加到切片。
作为一种特殊情况，要将空 JSON 数组解组为切片，Unmarshal 会用新的空切片替换该切片。

要将 JSON 数组解组为 Go 数组，Unmarshal 将 JSON 数组元素解码为相应的 Go 数组元素。
如果 Go 数组小于 JSON 数组，则多余的 JSON 数组元素将被丢弃。
如果 JSON 数组小于 Go 数组，则附加的 Go 数组元素将设置为零值。

要将 JSON 对象解组为映射，Unmarshal 首先建立要使用的映射。
如果映射为零，Unmarshal 会分配一个新映射。
否则，Unmarshal 会重用现有的映射，并保留现有的条目。
然后，Unmarshal 将 JSON 对象中的键值对存储到映射中。
映射的键类型必须是任何字符串类型、整数、实现 json.Unmarshaler 或实现 encoding.TextUnmarshaler。

如果 JSON 编码的数据包含语法错误，Unmarshal 将返回 SyntaxError。

如果 JSON 值不适合给定的目标类型，或者 JSON 数字溢出目标类型，Unmarshal 会跳过该字段并尽可能完成解组。
如果没有遇到更严重的错误，Unmarshal 将返回一个 UnmarshalTypeError 来描述最早的此类错误。
无论如何，都不能保证有问题的字段后面的所有剩余字段都将被解组到目标对象中。

通过将该 Go 值设置为 nil，JSON null 值将解组为接口、映射、指针或切片。
由于 null 通常在 JSON 中表示“不存在”，因此将 JSON null 解组为任何其他 Go 类型不会对该值产生任何影响，也不会产生错误。

解组带引号的字符串时，无效的 UTF-8 或无效的 UTF-16 代理项对不会被视为错误。
相反，它们被 Unicode 替换字符 U+FFFD 替换。
*/
var Unmarshal = json.Unmarshal

/*
MarshalIndent 与 Marshal 类似，但应用 Indent 来格式化输出。
输出中的每个 JSON 元素都将从新行开始，以前缀开头，后跟根据缩进嵌套的一个或多个缩进副本。
*/
var MarshalIndent = json.MarshalIndent

/*
NewDecoder 返回一个从 r 读取的新解码器。

解码器引入了自己的缓冲，并且可能从 r 读取超出请求的 JSON 值的数据。
*/
var NewDecoder = json.NewDecoder

/*
NewEncoder 返回一个写入 w 的新编码器。
*/
var NewEncoder = json.NewEncoder
