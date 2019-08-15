本包主要是对Java的.class字节码文件进行分析，并加载到内存ClassFile结构。

步骤1：加入包
将包加入Go的src库里
步骤2：引入包
import (
    "gdcore"
}
步骤3:使用
gdcore.ReadFormRC()  读取RC类型 func ReadFormRC(rc io.ReadCloser) (*ClassFile,error);
gdcore.ReadFromJar() 读取jar包  func ReadFromJar(zipFile string, classdest string) (*ClassFile,error)
或者
class := gdcore.ClassFile{}
class.Load(bs  []byte) error  读取class []byte数据
其他函数:
gdcore.AnalysisPower(sum int) []string;//用于分析ClassFile中的Access_flags权限