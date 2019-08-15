package gdcore
import (
    "archive/zip"
    "io"
    "io/ioutil"
    "errors"
)
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
        s=append(s,file.Name)
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
    err = c.Load(bs)               //从[]byte读取class内容
    if err != nil{
       return nil,err
    }
    return &c,nil                   //返回类地址
}