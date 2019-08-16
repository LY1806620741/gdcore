package gdcore

import(
    "archive/zip"
    "testing"
    "io/ioutil"
)

var testc ClassFile

func init(){
    testc,_ = ReadFormFile("./testdata/Lambda.class")
}

func TestReaderzip(t *testing.T) {
    //列出testData
    files, _ := ioutil.ReadDir("./testData")
    for _, f := range files {
        //打开zip
        if f.Name()[len(f.Name())-4:]==".zip"{
            t.Logf("Test testdata/"+f.Name())
            r, err := zip.OpenReader("testdata/"+f.Name())
            if err != nil {
                t.Fatal(err)
            }
            defer r.Close()
            //读取文件
            for _, f := range r.File {
                t.Logf("Contents of %s:\n", f.Name)
                if !f.FileInfo().IsDir(){
                    rc, err := f.Open()
                    if err != nil {
                        t.Fatal(err)
                    }
                    c,err:=ReadFormRC(rc)
                    if err != nil{
                        t.Fatal(err)
                    }
                    if c.Access_flags==0{
                        t.Fatal(err)
                    }
                    rc.Close()
                }
            }
        }
    }
}
//jar包列表函数测试
func TestListJar(t *testing.T){
    listjar,err := ListJar("./testdata/data-java-jdk-1.8.0.zip")
    if err != nil {
        t.Fatal(err)
    }
    if len(listjar) != 53{
        t.Fatal(err)
    }
}
//读取文件函数测试
func TestReaderfile(t *testing.T) {
    c,err := ReadFormFile("./testdata/Lambda.class")
    if err != nil {
        t.Fatal(err)
    }
    if c.Access_flags==0{
        t.Fatal(err)
    }
}
//读取jar包函数测试
func TestReaderjar(t *testing.T) {
    c,err := ReadFromJar("./testdata/data-java-jdk-1.8.0.zip","org/jd/core/test/Lambda.class")
    if err != nil {
        t.Fatal(err)
    }
    if c.Access_flags==0{
        t.Fatal(err)
    }
}
//权限分析测试
func TestAnalysisPower(t *testing.T){
    powerlist := AnalysisPower(testc.Access_flags)
    t.Logf("This class's power is %s",powerlist)
    if len(AnalysisPower(testc.Access_flags))!=2{
        t.Fatal()
    }

}
//测试byte转int
func TestByte2int(t *testing.T){
   if byte2int([]byte{0,1})!=1{
       t.Fatal()
   }
   if byte2int(testc.Magic)!=3405691582{
       t.Fatal()
   }
   if byte2int([]byte{0,1,1})!=257{
       t.Fatal()
   }
   if byte2int([]byte{0,159})!=159{
       t.Fatal()
   }
}
//翻译常量
func TestAnalysisConstant(t *testing.T){
    c,_:=AnalysisConstant(testc.Constant_pool[0])
    if (c.Type!="Methodref"){
        t.Fatal()
    }
    if (c.String!=""){
        t.Fatal()
    }
    if (c.Value[0]!=32||c.Value[1]!=126){
        t.Fatal()
    }
    c,_=AnalysisConstant(testc.Constant_pool[49])
    if (c.Type!="Utf8"){
        t.Fatal()
    }
    if (c.String!="index"){
        t.Fatal()
    }
    if (c.Value[0]!=0){
        t.Fatal()
    }
    c,_=AnalysisConstant(testc.Constant_pool[16])                       //实际上是常量池第17个
    if (c.Type!="Class"){
        t.Fatal()
    }
    if (c.String!=""){
        t.Fatal()
    }
    if (c.Value[0]!=159){
        t.Fatal()
    }
    if string(testc.Constant_pool[c.Value[0]-1].Info)!="java/util/Map"{    //class规范常量池是从1开始的,数组从0开始,)为什么不直接减一,因为读取力求真实，文件写的是多少就是多少)(为什么读取到内存不废弃0从1开始，因为节省内存)
        t.Fatal()
    }
    //读取测试class所有常量
    for _,s:=range testc.Constant_pool{
        cs,_ := AnalysisConstant(s)
        cs.ToString()
    }
}