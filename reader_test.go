package gdcore

import(
    "archive/zip"
    "testing"
    "io/ioutil"
)

var testc *ClassFile

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
                    _,err=ReadFormRC(rc)
                    if err != nil{
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
    _,err := ReadFormFile("./testdata/Lambda.class")
    if err != nil {
        t.Fatal(err)
    }
}
//读取jar包函数测试
func TestReaderjar(t *testing.T) {
    _,err := ReadFromJar("./testdata/data-java-jdk-1.8.0.zip","org/jd/core/test/Lambda.class")
    if err != nil {
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