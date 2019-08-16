# 欢迎使用 gdcore

**本包主要是go语言的第三方工具包，实现对Java的.class字节码文件进行分析，并加载到内存ClassFile结构。**



## 获取
### 方法1 从github获取
#### 步骤1 下载库

`go get github.com/LY1806620741/gdcore`

>Go1.8版本以后GOPATH会设置成当前用户文件夹下的go文件夹，但是之前的版本比较混乱，如果Windows下不会配置Go变量，可以使用我的GoBat工具
>https://github.com/LY1806620741/GoBat

#### 步骤2 导入包

在你的项目中

    import (
        "github.com/LY1806620741/gdcore"
		//"其他的包"
    )
### 方法2

####步骤1 下载库

   将各种方式下载的gdcore包放入***go语言根目录的src***或者***自定的GOPATH文件夹***里

####步骤二 导入包

    import (
		"gdcore"
	)

##使用
### 步骤1 使用提供的方法读取class文件

    gdcore.ReadFormRC()  func ReadFormRC(rc io.ReadCloser) (*ClassFile,error);                   //读取RC类型

[例子 点我](#ReadFormRC)

    gdcore.ReadFromJar() func ReadFromJar(zipFile string, classdest string) (*ClassFile,error)； //读取jar包里的类

[例子 点我](#ReadFromJar)

	gdcore.ReadFormFile(file string) (*ClassFile,error)；                                        //加载文件

[例子 点我](#ReadFormFile)

或者自己读取文件byte

    class := gdcore.ClassFile{}
    class.Load(bs  []byte) error  //读取class []byte数据

在库的reader_test.go里的测试方法也有用法可以参考

### 步骤2 进行你自己的处理，以下是数据结构和例子

数据的结构

class文件结构在Class File format.txt,go 数据结构在 struct

    //class文件结构
	type ClassFile struct{
		Magic                []byte             //魔数
		Minor_version        int                //次要版本
		Major_version        int                //主要版本
		Constant_pool_count  int                //常量
		Constant_pool        []Cp_Info          //常量池，通过索引导向常量
		Access_flags         int                //访问权限
		This_class           int                //class文件自身索引
		Supe_class           int                //父类索引
		Interfaces_count     int                //实现的接口个数 ps:java一个类可以实现多个接口
		Interfaces           []int              //每个接口的索引
		Fields_count         int                //字段个数
		Fields               []Field_Info       //字段信息
		Methods_count        int                //方法数
		Methods              []Method_Info      //方法信息
		Attributes_count     int                //类属性数
		Attributes           []Attribute_Info   //类属性信息
	}
	//常量
	type Cp_Info struct{
		Tag                  int                 //常量标志(代表数据类型)
		Info                 []byte              //常量数据(各种类型的结构不一样，暂时不关心所有数据类型)
	}
	//字段
	type Field_Info struct{
		Access_flags         int                 //访问权限
		Name_index           int                 //名字索引，指向常量池UTF8类型
		Descriptor_index     int                 //描述符索引
		Attributes_count     int                 //字段属性的个数
		Attributes           []Attribute_Info    //字段属性详细
	}
	//方法
	type Method_Info struct{
		Access_flags         int                 //权限
		Name_index           int                 //方法名常量池UTF8索引
		Descriptor_index     int                 //描述符索引
		Attributes_count     int                 //方法属性的个数
		Attributes           []Attribute_Info    //方法属性详细
	}
	//属性信息
	type Attribute_Info struct{
		Attribute_name_index int                 //指向常量库utf8
		Attribute_length     int                 //长度
		Info                 []byte              //属性信息
	}

<tag id="ReadFormRC">打开RC流例子</tag>

    rc,_=io.Open("./Lambda.class")                               //io.ReadCloser类型
	if err != nil{                                               //错误检查
	    Log.Fatalln("文件打开失败")
	}
	class,_ := gdcore.ReadFormRC(rc)                             //错误不检查就_丢掉，class就是Lambda.class的内存映像

<tag id="ReadFromJar">打开jar包里的class例子</tag>

	class,err := gdcore.ReadFromJar("./data-java-jdk-1.8.0.zip","org/jd/core/test/Lambda.class")
	if err != nil{                                               //错误检查
	    Log.Fatalln("文件打开失败")
	}

<tag id="ReadFormFile">打开class文件例子</tag>

	class,_ := gdcore.ReadFormFile("./Lambda.class")             //错误不检查就_丢掉，class就是Lambda.class的内存映像

<tag id="AnalysisPower">分析权限例子<tag>

    powerlist := AnalysisPower(class.Access_flags)              //分析class的标志
    fmt.Printf("This class's power is %s\n",powerlist)          //输出这个类拥有的属性Public等

<tag id="AnalysisConstant">分析常量池例子<tag>

    c,_ := AnalysisConstant(class.Constant_pool[0])             //分析class的常量池的第一个
    fmt.Printf("%s\n",c.ToString())                             //输出所翻译的数据

<tag id="lookClassMethodName">查看所有方法的名字例子</tag>

	class,_ := gdcore.ReadFromJar("./data-java-jdk-1.8.0.zip","org/jd/core/test/Lambda.class")    //加载class文件
	fmt.Printf("方法总共%d个\n",class.Methods_count)                                               //输出方法总数
	for i:=0;i<class.Methods_count;i++ {                                                          //循环每个方法
        var methodname=class.Constant_pool[class.Methods[i].Name_index].Info                      //查找其常量索引
		fmt.Printf("%s\n",methodname)
	}
	//因为Method中只存了索引，名字是在常量池中，Method.Name_index标志了其位置，class.Constant_pool.Tag是其常量类型
	//class.Constant_pool.Info是其常量[]byte数据，所有的方法和变量开头字母大写才能给其他包访问
	
其余例子可以自行查看reader_test.go,这是测试文件,进入gdcore包`go test`可以跑测试
	
## 其他函数:

	gdcore.AnalysisPower(sum int) []string;                  //用于分析Access_flags权限（Public，Private等）

[例子 点我](#AnalysisPower)

	gdcore.AnalysisConstant(cp Cp_Info) (ConstantInfo,error);//用于翻译常量
	type ConstantInfo struct{                                //返回的类型
		Type    string
		String  string
		Value   []int
		ToString() func string;                              //输出字符串
	}
[例子 点我](#AnalysisConstant)

## 其他的话

1. 用于测试的jar包来自jd-core
2. 本包是借鉴了jd-gui和jvm文档所开发的，目前只实现了读取class文件到内存，以后看看有时间再实现反编译（也就是把内存里的数据转成string或文件），其他功能还在完善中
3. 文档只有中文版，除了程序怕字符兼容问题，能中文肯定要中文
