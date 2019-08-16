package gdcore
import (
    "archive/zip"
    "io"
    "io/ioutil"
    "errors"
    "strconv"
    "strings"
)

//加载常量
type ConstantInfo struct{
    Type    string
    String  string
    Value   []int
}

//转化可视输出
func (c ConstantInfo) ToString() string{
    switch c.Type{
        case "Utf8":
            return c.Type+" "+c.String
        case "Class":
            return c.Type+" class_index:"+strconv.Itoa(c.Value[0])
        case "Fieldref","Methodref","InterfaceMethodref":
            return c.Type+" class_index:"+strconv.Itoa(c.Value[0])+" name_and_type_index:"+strconv.Itoa(c.Value[1])
        case "String":
            return c.Type+" string_index:"+strconv.Itoa(c.Value[0])
        case "Integer","Float":
            return c.Type+" "+strconv.Itoa(c.Value[0])
        case "Long","Double":
            return c.Type+" high:"+strconv.Itoa(c.Value[0])+" low:"+strconv.Itoa(c.Value[1])
        case "NameAndType":
            return c.Type+" name_index:"+strconv.Itoa(c.Value[0])+" descriptor_index:"+strconv.Itoa(c.Value[1])
        case "MethodHandle":
            return c.Type+" reference_kind:"+strconv.Itoa(c.Value[0])+" reference_index:"+strconv.Itoa(c.Value[1])+" "
        case "MethodType":
            return c.Type+" descriptor_index:"+strconv.Itoa(c.Value[0])
        case "InvokeDynamic":
            return c.Type+" bootstrap_method_attr_index:"+strconv.Itoa(c.Value[0])+" name_and_type_index:"+strconv.Itoa(c.Value[1])+" "
        default:
            return ""
    }
}

//查询reference_kind解释
func getreference_kind(k int) string{
    if (k>0&&k<10){
        return strings.Split(`    1	REF_getField	getfield C.f:T
        2	REF_getStatic	getstatic C.f:T
        3	REF_putField	putfield C.f:T
        4	REF_putStatic	putstatic C.f:T
        5	REF_invokeVirtual	invokevirtual C.m:(A*)T
        6	REF_invokeStatic	invokestatic C.m:(A*)T
        7	REF_invokeSpecial	invokespecial C.m:(A*)T
        8	REF_newInvokeSpecial	new C; dup; invokespecial C.<init>:(A*)V
        9	REF_invokeInterface	invokeinterface C.m:(A*)T`,"\n")[k-1]
    }else{
        return "Not Found"
    }
}

//分析权限
func AnalysisPower(sum int) []string{
    var power []string
    if sum-0x4000>=0{
        power=append(power,"ACC_ENUM")
        sum-=0x4000
    }
    if sum-0x2000>=0{
        power=append(power,"ACC_ANNOTATION")
        sum-=0x2000
    }
    if sum-0x1000>=0{
        power=append(power,"ACC_SYNTHETIC")
        sum-=0x1000
    }
    if sum-0x0400>=0{
        power=append(power,"ACC_ABSTRACT")
        sum-=0x0400
    }
    if sum-0x0200>=0{
        power=append(power,"ACC_INTERFACE")
        sum-=0x0200
    }
    if sum-0x0020>=0{
        power=append(power,"ACC_SUPER")//用于兼容早期编译器
        sum-=0x0020
    }
    if sum-0x0010>=0{
        power=append(power,"ACC_FINAL")//不能有子类
        sum-=0x0010
    }
    if sum-0x0001>=0{
        power=append(power,"ACC_PUBLIC")//包外可访问
        sum-=0x0001
    }
    return power
}
//查看jar包目录
func ListJar(zipFile string) ([]string,error){
    //读取zip
    reader, err := zip.OpenReader(zipFile)
    if err != nil {
        return nil,err
    }
    defer reader.Close()
    //循环每个文件
    var s []string
    for _, file := range reader.File {
        if !file.FileInfo().IsDir(){
            s=append(s,file.Name)
        }
    }
    return s,nil
}
//读取jar("*.jar","com/jieshao/demo/main.class")
func ReadFromJar(zipFile string, classdest string) (*ClassFile,error) {
    //读取zip
    reader, err := zip.OpenReader(zipFile)
    if err != nil {
        return nil,err
    }
    defer reader.Close()
    //循环每个文件
    for _, file := range reader.File {
        if (file.Name==classdest){
            rc, err := file.Open()
            if err != nil {
                return nil,err
            }
            defer rc.Close()
            class,err := ReadFormRC(rc)
            if err != nil {
                return nil,err
            }
            return class,nil
        }
    }
    return nil,errors.New("Error. Can't find the class")
}
//加载reader
func ReadFormRC(rc io.ReadCloser) (*ClassFile,error){
    c := ClassFile{}                //初始化内存数据结构
    bs,err:= ioutil.ReadAll(rc)     //一次读文件
    if err != nil{
       return nil,err
    }
    err = c.Load(bs)                //从[]byte读取class内容
    if err != nil{
       return nil,err
    }
    return &c,nil                   //返回类地址
}
//加载文件
func ReadFormFile(file string) (*ClassFile,error){
    c := ClassFile{}                //初始化内存数据结构
    bs,err:= ioutil.ReadFile(file)  //一次读文件
    if err != nil{
       return nil,err
    }
    err = c.Load(bs)                //从[]byte读取class内容
    if err != nil{
       return nil,err
    }
    return &c,nil                   //返回类地址
}
//翻译常量
func AnalysisConstant(cp Cp_Info) (ConstantInfo,error){ 
    switch cp.Tag{                                                     //对常量池进行分类，但不关心其余类型
        case 1:                                                        //CONSTANT_Utf8_info
            return ConstantInfo{"Utf8",string(cp.Info),[]int{0}},nil
        case 3:                                                        //CONSTANT_Integer_info format
            return ConstantInfo{"Integer","",[]int{byte2int(cp.Info)}},nil
        case 4:                                                        //CONSTANT_Float_info
            return ConstantInfo{"Float","",[]int{byte2int(cp.Info)}},nil
        case 5:                                                        //CONSTANT_Long_info  All 8-byte constants take up two entries in the constant_pool table of the class file. If a CONSTANT_Long_info or CONSTANT_Double_info structure is the item in the constant_pool table at index n, then the next usable item in the pool is located at index n+2. The constant_pool index n+1 must be valid but is considered unusable.
            return ConstantInfo{"Long","",[]int{byte2int(cp.Info[:2]),byte2int(cp.Info[2:])}},nil
        case 6:                                                        //CONSTANT_Double_info
            return ConstantInfo{"Double","",[]int{byte2int(cp.Info[:2]),byte2int(cp.Info[2:])}},nil
        case 7:                                                        //CONSTANT_Class_info format
            return ConstantInfo{"Class","",[]int{byte2int(cp.Info)}},nil
        case 8:                                                        //CONSTANT_String_info
            return ConstantInfo{"String","",[]int{byte2int(cp.Info)}},nil
        case 9:                                                        //CONSTANT_Fieldref_info
           return ConstantInfo{"Fieldref","",[]int{byte2int(cp.Info[:2]),byte2int(cp.Info[2:])}},nil
        case 10:                                                        //CONSTANT_Methodref_info
            return ConstantInfo{"Methodref","",[]int{byte2int(cp.Info[:2]),byte2int(cp.Info[2:])}},nil
        case 11:                                                        //CONSTANT_InterfaceMethodref_info
            return ConstantInfo{"InterfaceMethodref","",[]int{byte2int(cp.Info[:2]),byte2int(cp.Info[2:])}},nil
        case 12:                                                        //CONSTANT_NameAndType_info
            return ConstantInfo{"NameAndType","",[]int{byte2int(cp.Info[:2]),byte2int(cp.Info[2:])}},nil
        case 15:                                                        //CONSTANT_MethodHandle_info
            return ConstantInfo{"MethodHandle","",[]int{byte2int(cp.Info[:1]),byte2int(cp.Info[1:])}},nil
        case 16:                                                        //CONSTANT_MethodType_info
            return ConstantInfo{"MethodType","",[]int{byte2int(cp.Info)}},nil
        case 18:                                                        //CONSTANT_InvokeDynamic_info
            return ConstantInfo{"InvokeDynamic","",[]int{byte2int(cp.Info[:2]),byte2int(cp.Info[2:])}},nil
        default:                                                        //未知
    }
    return ConstantInfo{},errors.New("Error. Unknown Constant type")
}
//内函数[]byte转int
func byte2int(data []byte)int{
       var ret int = 0
       var len uint = uint(len(data))
       var i uint
       for i=0; i<len; i++{
              ret = ret | (int(data[i]) << ((len-1-i)*8))
       }
       return ret
}