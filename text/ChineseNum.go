package text

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const DOT string = ".点块元"
const Negative string = "负-－"
const Ge = "0０零1１一壹2２二两貳贰3３三叁參4４四肆5５五伍6６六陸陆7７七柒8８八捌9９九玖" + DOT

var ValidChar = []string{Ge + "十拾百佰千仟万萬", Ge + "十拾百佰千仟", Ge + "十拾百佰", Ge + "十拾", Ge}
var Numbers = []string{"0０零", "1１一壹", "2２二两貳贰", "3３三叁參", "4４四肆", "5５五伍", "6６六陸陆", "7７七柒", "8８八捌", "9９九玖"}
var Levels = []string{"億亿", "万萬", "千仟", "佰百", "十拾"}
var Scale = []int64{100000000, 10000, 1000, 100, 10, 1}

type ChineseNumber struct {
	Input   *[]rune
	Begin   int
	End     int
	Value   int64
	Decimal string //为了避免精度问题，直接用字符串表示小数
}

func (cn *ChineseNumber) GetBegin() int {
	return cn.Begin
}
func (cn *ChineseNumber) SetBegin(begin int) {
	cn.Begin = begin
}
func (cn *ChineseNumber) GetEnd() int {
	return cn.End
}
func (cn *ChineseNumber) SetEnd(end int) {
	cn.End = end
}
func (cn *ChineseNumber) SetPosition(position int) {
	cn.SetBegin(position)
	cn.SetEnd(position)
}
func (cn *ChineseNumber) GetValue() int64 {
	return cn.Value
}
func (cn *ChineseNumber) SetValue(value int64) {
	cn.Value = value
}
func (cn *ChineseNumber) GetDecimal() string {
	return cn.Decimal
}
func (cn *ChineseNumber) SetDecimal(decimal string) {
	cn.Decimal = decimal
}
func (cn *ChineseNumber) ToFloat() float64 {
	if cn.Decimal == "" {
		return float64(cn.Value)
	} else {
		if d, err := strconv.ParseFloat("0."+cn.Decimal, 64); err == nil {
			return float64(cn.Value) + d
		} else {
			panic(err)
		}
	}
}

func (cn *ChineseNumber) ToString() string {
	str := strconv.FormatInt(cn.Value, 10)
	if cn.Decimal != "" {
		str = str + "." + cn.Decimal
	}
	return str
}

// ------------------ 实现排序接口 -----------------
type Nums []*ChineseNumber

// 实现sort.Interface接口取元素数量方法
func (s Nums) Len() int {
	return len(s)
}

// 实现sort.Interface接口比较元素方法
func (s Nums) Less(i, j int) bool {
	return s[i].GetBegin() < s[j].GetBegin() //升序排列
}

// 实现sort.Interface接口交换元素方法
func (s Nums) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// 添加一个业务方法
func (s *Nums) Add(result *ChineseNumber) {
	// 处理尾部的小数
	if result.Decimal == "" {
		result.ReadDecimal()
	}
	// 检查是否有负数情况
	if result.Begin > 0 && strings.ContainsRune(Negative, (*result.Input)[result.Begin-1]) {
		result.SetValue(-result.Value)
		result.SetBegin(result.Begin - 1)
	}
	// 销毁已识别的数字，防止重复识别。
	for k := result.Begin; k < result.End; k++ {
		(*result.Input)[k] = ' '
	}
	*s = append(*s, result)
}

// ------------------ 业务方法 -----------------
/*
读取右边的数字
beginIndex		右边的开始位置
endIndex		右边的结束位置
levelIndex		当前要处理的数字级别
readDecimal		是否需要读取小数
@return		true:右边有数字，否则：右边没有数字
*/
func (cn *ChineseNumber) ReadRight(beginIndex, endIndex, levelIndex int, readDecimal bool) bool {
	limit := len(*cn.Input)
	if beginIndex < endIndex && endIndex <= limit {
		if levelIndex < 0 || levelIndex > len(Levels) {
			fmt.Printf("ReadRight, 参数错误。levelIndex=%d\r\n", levelIndex)
		} else if endIndex-beginIndex == 1 { //处理单个字符
			value := Read1Char((*cn.Input)[beginIndex])
			if value != -1 {
				cn.SetValue(value)
				cn.SetBegin(beginIndex)
				cn.SetEnd(endIndex)
				return true
			}
		} else {
			cn.SetBegin(beginIndex)
			cn.SetEnd(beginIndex)
			var ch rune
			if levelIndex == len(Levels) {
				ch = (*cn.Input)[beginIndex]
				if strings.ContainsRune(DOT, ch) { //例如，一百点二，十点八
					if readDecimal && cn.ReadDecimal() {
						cn.SetValue(0)
						cn.SetBegin(beginIndex)
						return true
					}
				} else {
					for ii, key := range Numbers {
						if strings.ContainsRune(key, ch) {
							cn.SetValue(int64(ii))
							cn.SetBegin(beginIndex)
							cn.SetEnd(beginIndex + 1)
							if ii == 0 && beginIndex+1 < endIndex { //遇到“零”，就解析下一个字符，例如：一百零一
								ch = (*cn.Input)[beginIndex+1]
								for m := len(Numbers) - 1; m > -1; m-- {
									if strings.ContainsRune(Numbers[m], ch) {
										cn.SetValue(int64(m))
										cn.SetEnd(beginIndex + 2)
										break
									}
								}
							}
							if readDecimal {
								cn.ReadDecimal()
							}
							return true
						}
					}
				}
			} else {
				key := Levels[levelIndex]
				pos := -1
				var max = endIndex
				if !readDecimal {
					switch levelIndex {
					case 0:
						max = endIndex
						break
					case 1:
						max = cn.Begin + 8
						break // “找“万”，例如：九千九百九十九万九千九百九十九
					case 2:
						max = cn.Begin + 3
						break // 找“千”，例如：两亿零一千
					case 3:
						max = cn.Begin + 3
						break // 找“百”，例如：三万零一百
					case 4:
						max = cn.Begin + 3
						break // 找“十”，例如：三千零一十
					case 5:
						max = cn.Begin
						break // “十”的后面无需找什么
					default:
						max = cn.Begin
						break
					}
				}
				if max > endIndex {
					max = endIndex
				}
				if max > cn.Begin {
					//if Debug==1 { fmt.Printf(">>>>>>>>>>> 在右边（%s）中寻找“%s”。\r\n", string((*cn.Input)[cn.Begin:max]), key) }
					for x := cn.Begin; x < max; x++ {
						ch = (*cn.Input)[x]
						if strings.ContainsRune(key, ch) {
							pos = x
							break
						} else if !strings.ContainsRune(ValidChar[levelIndex], ch) { //只要出现了无效字符，就退出循环
							break
						}
					}
				}
				if pos == -1 {
					return cn.ReadRight(beginIndex, endIndex, levelIndex+1, readDecimal) //递归
				} else {
					return cn.ReadNum(beginIndex, endIndex, pos, levelIndex, readDecimal)
				}
			}
		}
	}
	return false
}

/*
读取左边的数字
beginIndex		左边的开始位置
endIndex		左边的结束位置
levelIndex		当前要处理的数字级别
readDecimal		是否需要读取小数
@return		true:左边有数字，否则：左边没有数字
*/
func (cn *ChineseNumber) ReadLeft(beginIndex, endIndex, levelIndex int, readDecimal bool) bool {
	if endIndex > beginIndex && beginIndex >= 0 {
		if levelIndex < 0 || levelIndex > len(Levels) {
			fmt.Printf("ReadLeft, 参数错误。levelIndex=%d\r\n", levelIndex)
		} else if endIndex-beginIndex == 1 { //处理单个字符
			value := Read1Char((*cn.Input)[beginIndex])
			if value != -1 {
				cn.SetValue(value)
				cn.SetBegin(beginIndex)
				cn.SetEnd(endIndex)
				return true
			}
		} else {
			cn.SetBegin(endIndex)
			cn.SetEnd(endIndex)
			var ch rune
			if levelIndex == len(Levels) { //处理个位，该函数的递归出口
				index := endIndex - 1
				if readDecimal {
					var nn string
					var ii, vv int
					var key string
					pp := index //记录数字的开始位置
					for ; index > -1; index-- {
						ch = (*cn.Input)[index]
						if strings.ContainsRune(DOT, ch) {
							if strings.IndexByte(nn, '.') == -1 {
								nn = "." + nn
								pp = index
							} else {
								break
							}
						} else {
							vv = -1
							for ii, key = range Numbers {
								if strings.ContainsRune(key, ch) {
									vv = ii
									break
								}
							}
							if vv == -1 { //遇到无效字符
								break
							} else {
								nn = strconv.Itoa(vv) + nn //fmt.Sprintf("%d",ii)
								pp = index
							}
						}
					}
					if nn != "" {
						point := strings.IndexByte(nn, '.')
						if point == -1 {
							vv = int(nn[len(nn)-1] - '0')
							cn.SetValue(int64(vv))
							cn.SetDecimal("")
							cn.SetBegin(endIndex - 1)
							cn.SetEnd(endIndex)
							if len(nn) > 1 && nn[len(nn)-2] == '0' { //例如：一万零二，“二”的前面有个0
								cn.SetBegin(endIndex - 2)
							}
							return true
						} else {
							if vv, err := strconv.ParseInt(nn[:point], 10, 64); err == nil {
								cn.SetValue(vv)
								cn.SetDecimal(nn[point+1:])
								cn.SetBegin(pp)
								cn.SetEnd(endIndex)
								return true
							}
						}
					}
				} else {
					ch = (*cn.Input)[index]
					for ii, key := range Numbers {
						if strings.ContainsRune(key, ch) {
							cn.SetValue(int64(ii))
							cn.SetBegin(index)
							cn.SetEnd(endIndex)
							//if Debug==1 { fmt.Printf("在左边，个位上的值是“%d”，readDecimal=%v, levelIndex=%d, cn=%+v\r\n", ii, readDecimal, levelIndex, cn) }
							return true
						}
					}
				}
			} else {
				key := Levels[levelIndex]
				pos := -1
				var min int
				if !readDecimal {
					switch levelIndex {
					case 0:
						min = 0
						break
					case 1:
						min = cn.Begin - 8
						break // “亿”前找“万”
					case 2:
						min = cn.Begin - 6
						break // “万”前找“千”
					case 3:
						min = cn.Begin - 3
						break // “千”前找“百”
					case 4:
						min = cn.Begin - 2
						break // “百”前找“十”
					case 5:
						min = cn.Begin
						break // “十”的前面无需找什么
					default:
						min = cn.Begin
						break
					}
				}
				if min < beginIndex {
					min = beginIndex
				}
				if cn.Begin > min {
					for x := cn.Begin - 1; x >= min; x-- {
						ch = (*cn.Input)[x]
						if strings.ContainsRune(key, ch) {
							pos = x
							break
						} else if !strings.ContainsRune(ValidChar[levelIndex], ch) { //只要出现了无效字符，就退出循环
							break
						}
					}
				} /*else{
					if Debug==1 { fmt.Printf("“%s”的前面无需找“%s”。\r\n", Levels[levelIndex-1], key) }
				}*/
				if pos == -1 {
					if cn.Begin == min {
						return cn.ReadLeft(beginIndex, endIndex, 5, false) //递归，直接解析个位
					} else {
						return cn.ReadLeft(beginIndex, endIndex, levelIndex+1, readDecimal) //递归
					}
				} else {
					return cn.ReadNum(beginIndex, endIndex, pos, levelIndex, readDecimal)
				}
			}
		}
	}
	return false
}

/*
在指定位置解析数字
beginIndex	字符数组的开始位置（含）
endIndex	字符数组的结束位置（不含）
pos			表示“数级单位”的字符位置
levelIndex	当前要计算的“数级单位”
readDecimal	是否需要读取小数
@return  是否解析成功
*/
func (cn *ChineseNumber) ReadNum(beginIndex, endIndex, pos, levelIndex int, readDecimal bool) bool {
	hasLeft := cn.ReadLeft(beginIndex, pos, levelIndex+1, readDecimal)
	if hasLeft && cn.End < pos {
		hasLeft = cn.ReadLeft(cn.End, pos, levelIndex+1, readDecimal)
	}
	if levelIndex == len(Levels)-1 && strings.ContainsRune(Levels[levelIndex], (*cn.Input)[pos]) { //因为人们的习惯问题，“十”存在一些特殊情况
		if !hasLeft { //处理“十”前面的“一”被省略的情况，例如：十一，十二，十三，十四，十五，十六，十七，十八，十九
			cn.SetValue(1)
			cn.SetBegin(pos)
			cn.SetEnd(pos)
			cn.SetDecimal("")
			hasLeft = true
		} else if cn.Value == 0 && cn.Decimal == "" {
			cn.SetValue(1) //例如：一万零十一，“十”的左边应该解析成“1”，因为“一”被省略了。
		}
	}
	if hasLeft { // 先要读左边，例如：二百
		if !readDecimal && cn.Value == 0 && cn.Decimal == "" {
			fmt.Printf("左边不能以“0”开头。\r\n")
		} else {
			begin := cn.Begin
			value := cn.Value * Scale[levelIndex]
			dd := cn.Decimal
			hasRight := cn.ReadRight(pos+1, endIndex, levelIndex+1, readDecimal)
			if hasRight && cn.Begin > pos+1 {
				hasRight = cn.ReadRight(pos+1, cn.Begin, levelIndex+1, readDecimal)
			}
			if hasRight {
				if dd != "" {
					fmt.Printf("左边有小数(%s)，右边却有数字(%d)，表达式不合理。\r\n", dd, cn.Value)
					return false
				}
				ch := (*cn.Input)[pos+1]
				if cn.Value < 10 && !strings.ContainsRune(Numbers[0], ch) { //排除：二百零五
					cn.Value = cn.Value * Scale[levelIndex] / 10 //例如：二百五, 两亿五
				}
				value = value + cn.Value
			} else {
				cn.SetEnd(pos + 1)
				if dd != "" {
					zheng := []byte(strconv.FormatInt(Scale[levelIndex], 10)[1:])
					if len(zheng) != 0 {
						for c := 0; c < len(zheng) && c < len(dd); c++ {
							zheng[c] = dd[c]
						}
						if len(zheng) < len(dd) {
							dd = dd[len(zheng):]
						} else {
							dd = ""
						}
						cn.SetDecimal(dd)
						if vv, err := strconv.ParseInt(string(zheng), 10, 64); err == nil {
							value += vv
						} else {
							panic(fmt.Sprintf("非法的小数“0.%s”。%v", cn.Decimal, err))
						}
					}
				}
				//if Debug==1 { fmt.Printf("在右边（%s）找不到“%s”。\r\n", string((*cn.Input)[pos+1:endIndex]), Levels[levelIndex+1]) }
			}
			cn.SetBegin(begin)
			cn.SetValue(value)
			return true
		}
	}
	return false
}

/**
将字符串中的所有数字提取出来
input  字符数组
@return  结果数组
*/
func ReadNums(input *[]rune) *Nums {
	var list = &Nums{}
	var pos, end, beginIndex, endIndex, tt int
	var value2 int64
	var limit = len(*input)
	var maxTimes = limit / 2
	for index, level := range Levels {
		end = 0
		for tt = 0; tt < maxTimes; tt++ { //有限循环次数
			pos = find(input, level, end+1)
			if pos == -1 {
				break
			} else {
				//if Debug==2 { fmt.Printf("找到了“%c”。\r\n", (*input)[pos]) }
				end = pos + 1
				beginIndex = pos - 30
				endIndex = pos + 40
				if beginIndex < 0 {
					beginIndex = 0
				}
				if endIndex > limit {
					endIndex = limit
				}
				cn := &ChineseNumber{Input: input}
				if cn.ReadRight(pos+1, endIndex, index+1, true) {
					value2 = cn.Value
					end = cn.GetEnd()
					if cn.ReadLeft(beginIndex, pos, index+1, false) {
						cn.SetValue(cn.Value*Scale[index] + value2)
						cn.SetEnd(end)
						list.Add(cn)
					} else if index == 4 {
						cn.SetValue(Scale[index] + value2)
						cn.SetEnd(end)
						list.Add(cn)
					}
				} else { //没有零头，左边允许小数
					if !cn.ReadLeft(beginIndex, pos, index+1, true) { //左右两边都没有数字
						cn.SetValue(1)
						cn.SetBegin(pos)
					}
					value2 = cn.Value * Scale[index]
					if cn.Decimal != "" {
						zheng := []byte(strconv.FormatInt(Scale[index], 10)[1:])
						if len(zheng) != 0 {
							for c := 0; c < len(zheng) && c < len(cn.Decimal); c++ {
								zheng[c] = cn.Decimal[c]
							}
							if len(zheng) < len(cn.Decimal) {
								cn.SetDecimal(cn.Decimal[len(zheng):])
							} else {
								cn.SetDecimal("")
							}
							if vv, err := strconv.ParseInt(string(zheng), 10, 64); err == nil {
								value2 += vv
							} else {
								panic(fmt.Sprintf("非法的小数“0.%s”。%v", cn.Decimal, err))
							}
						}
					}
					cn.SetValue(value2)
					cn.SetEnd(end)
					list.Add(cn)
					//if Debug==2 { fmt.Printf("---------- 字符串变成了:%s|cn=%+v\r\n", string(*input), cn) }
				}
			}
			(*input)[pos] = ' ' //防止这个关键词再次被找到
		}
	}
	// 处理小于10的数
	var ch rune
	var key string
	for pos = len(*input) - 1; pos > -1; pos-- {
		ch = (*input)[pos]
		if ch > 255 { //只处理汉字
			for tt, key = range Numbers {
				if strings.ContainsRune(key, ch) {
					cn := &ChineseNumber{Input: input}
					cn.SetValue(int64(tt))
					cn.SetBegin(pos)
					cn.SetEnd(pos + 1)
					list.Add(cn)
					break
				}
			}
		}
	}
	sort.Sort(list)
	return list
}

/**
 * 读取尾部的小数，例如：七百零三块五
 * @return  true：解析成功，否则：解析失败。
 */
func (cn *ChineseNumber) ReadDecimal() bool {
	limit := len(*cn.Input)
	if cn.End < limit && strings.ContainsRune(DOT, (*cn.Input)[cn.End]) {
		var ch rune
		var suffix = ""
		index := cn.End + 1
		var ii int
		var key string
		for ; index < limit; index++ {
			ch = (*cn.Input)[index]
			for ii, key = range Numbers {
				if strings.ContainsRune(key, ch) {
					suffix += strconv.Itoa(ii) //fmt.Sprintf("%d",ii)
					break
				}
			}
			if ii == len(Numbers)-1 { //遇到无效字符
				break
			}
		}
		if suffix != "" {
			cn.SetDecimal(suffix)
			cn.SetEnd(index)
			return true
		}
	}
	return false
}

/*
将一个字符转换成阿拉伯数字，如果不是数字，将返回-1
ch  要转换的字符
*/
func Read1Char(ch rune) int64 {
	var ii int
	var key string
	for ii, key = range Numbers {
		if strings.ContainsRune(key, ch) {
			return int64(ii)
		}
	}
	for ii, key = range Levels {
		if strings.ContainsRune(key, ch) {
			return Scale[ii]
		}
	}
	return -1
}

/*
查找关键字
input	字符数组
key		要查找的关键字
begin	查找的开始位置
*/
func find(input *[]rune, key string, begin int) int {
	var ch rune
	var limit = len(*input)
	for x := begin; x < limit; x++ {
		ch = (*input)[x]
		if strings.ContainsRune(key, ch) {
			return x
		}
	}
	return -1
}

/*
将汉语中的数字替换成阿拉伯数字
str  汉语字符串
*/
func Replace(str string) string {
	input := []rune(str)
	list := *ReadNums(&input)
	var num *ChineseNumber
	for x := len(list) - 1; x > -1; x-- {
		num = list[x]
		input = append(input[0:num.Begin], append([]rune(num.ToString()), input[num.End:]...)...)
		// input = append(append(input[0:num.Begin], []rune(num.String())...), input[num.End:]...)	//这句代码是错的，一定要从后往前处理，否则会导致指针错乱。
	}
	return string(input)
}
