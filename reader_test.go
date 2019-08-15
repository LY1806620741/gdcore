package gdcore

import(
    "archive/zip"
    "testing"
    "io/ioutil"
)

func TestReaderzip(t *testing.T) {
    //列出testData
    files, _ := ioutil.ReadDir("./testData")
    for _, f := range files {
        //打开zip
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
    /*
	*/
	// Output:
	// Contents of README:
	// This is the source code repository for the Go programming language.
}